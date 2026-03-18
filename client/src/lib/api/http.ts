import {error} from "@sveltejs/kit";
import {authToken, loadToken} from "./auth";
import {get} from "svelte/store";
import { PUBLIC_SERVER_URL as API_BASE } from '$env/static/public';

export interface APIError {
	code: string;
	message: string;
	details?: Record<string, unknown>;
	internal?: string;
}

export async function rawRequest(endpoint: string, options: RequestInit = {}, auth = true): Promise<Response> {
	const headers = new Headers(options?.headers)
	if (!headers.has('Content-Type')) {
		headers.set('Content-Type', 'application/json');
	}

	if (auth) {
		// Use the cached store value, or fall back to loading from IndexedDB if
		// the store hasn't been hydrated yet (e.g. a child load running concurrently
		// with the layout load that normally does this).
		const token = get(authToken) ?? await loadToken();
		if (!token) {
			throw error(401, {
				message: 'Unauthorized: No authentication token found',
				base: API_BASE,
				endpoint,
				method: options.method ?? 'GET',
				body: null,
			});
		}

		headers.set('Authorization', `Bearer ${token}`);
	}

	try {
		return await fetch(`${API_BASE}${endpoint}`, {
			...options,
			headers,
		});
	} catch (err) {
		throw error(503, {
			message: `Service Unavailable: Could not reach API`,
			base: API_BASE,
			endpoint,
			method: options.method ?? 'GET',
			body: err instanceof Error ? err.message : String(err),
		});
	}
}

export async function request<T = void>(endpoint: string, options: RequestInit = {}, auth = true): Promise<T> {
	const response = await rawRequest(endpoint, options, auth);

	const contentType = response.headers.get('content-type');
	const isJson = !!contentType?.includes('application/json');

	if (!response.ok) {
		const body: APIError | string = isJson ? (await response.json()).error : await response.text()

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

	return await (isJson ? response.json() : response.text())
}
