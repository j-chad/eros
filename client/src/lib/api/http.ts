import {error} from "@sveltejs/kit";
import {readCachedToken} from "$lib/services/session";

const API_BASE: string = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api';

export interface APIError {
	code: string;
	message: string;
	details?: Record<string, unknown>;
	internal?: string;
}

export async function request<T = void>(endpoint: string, options: RequestInit = {}): Promise<T> {
	const token = await readCachedToken()

	if (!token) {
		throw error(401, {
			message: 'Unauthorized: API key is missing',
			base: API_BASE,
			endpoint,
			method: options.method ?? 'GET',
			body: null,
		});
	}

	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		'Authorization': token,
		...options.headers,
	};

	let response: Response;
	try {
		response = await fetch(`${API_BASE}${endpoint}`, {
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

	const contentType = response.headers.get('content-type');
	const isJson = !!contentType?.includes('application/json');

	if (response.status === 401) {

	}

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

	return isJson ? response.json() : null as T;
}
