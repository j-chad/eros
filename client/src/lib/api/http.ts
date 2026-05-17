import { error, redirect } from '@sveltejs/kit';
import { authToken, loadToken, clearToken } from './auth';
import { get } from 'svelte/store';
import { PUBLIC_SERVER_URL as API_BASE } from '$env/static/public';
import { markOffline } from '$lib/online.svelte';

export interface APIError {
	code: string;
	message: string;
	details?: Record<string, unknown>;
	internal?: string;
}

/**
 * Thrown when a `fetch()` call fails with a network error (connection refused,
 * DNS failure, timeout, etc.). This is NOT a SvelteKit `HttpError` — it won't
 * hit the error boundary. Service-layer code can catch it and retry from cache.
 */
export class ServerUnreachableError extends Error {
	constructor(
		public readonly endpoint: string,
		public readonly method: string,
		cause?: unknown,
	) {
		super(`Server unreachable: ${method} ${endpoint}`);
		this.name = 'ServerUnreachableError';
		this.cause = cause;
	}
}

export async function rawRequest(
	endpoint: string,
	options: RequestInit = {},
	auth = true,
): Promise<Response> {
	const headers = new Headers(options?.headers);
	if (!headers.has('Content-Type')) {
		headers.set('Content-Type', 'application/json');
	}

	if (auth) {
		// Use the cached store value, or fall back to loading from IndexedDB if
		// the store hasn't been hydrated yet (e.g. a child load running concurrently
		// with the layout load that normally does this).
		const token = get(authToken) ?? (await loadToken());
		if (!token) {
			throw redirect(307, '/login');
		}

		headers.set('Authorization', `Bearer ${token}`);
	}

	const method = options.method ?? 'GET';

	try {
		return await fetch(`${API_BASE}${endpoint}`, {
			...options,
			headers,
		});
	} catch (err) {
		markOffline();
		throw new ServerUnreachableError(endpoint, method, err);
	}
}

export async function request<T = void>(
	endpoint: string,
	options: RequestInit = {},
	auth = true,
): Promise<T> {
	const response = await rawRequest(endpoint, options, auth);

	// Token has been revoked or expired — clear local state and send the user
	// back to login. This is the lazy counterpart to the optimistic local-only
	// auth check: we let the user in immediately with a cached token, and boot
	// them here if the server disagrees.
	if (auth && response.status === 401) {
		await clearToken();
		throw redirect(307, '/login');
	}

	const contentType = response.headers.get('content-type');
	const isJson = !!contentType?.includes('application/json');

	if (!response.ok) {
		const body: APIError | string = isJson ? (await response.json()).error : await response.text();

		throw error(response.status, {
			message: typeof body === 'string' ? body : body.message,
			base: API_BASE,
			endpoint,
			method: options.method ?? 'GET',
			body,
		});
	}

	if (response.status === 204) {
		return null as T;
	}

	return await (isJson ? response.json() : response.text());
}
