import type {APIError, RegistrationToken} from "$lib/types";
import {auth} from "$lib/auth.svelte";

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api';

class APIException extends Error {
    constructor(
        message: string,
        public status: number,
        public data?: APIError
    ) {
        super(message);
        this.name = 'APIException';
    }
}

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    if (!auth.apiKey) {
        throw new APIException('Not authenticated', 401);
    }

    const headers: HeadersInit = {
        'Content-Type': 'application/json',
        'Authorization': `Admin ${auth.apiKey}`,
        ...options.headers,
    };

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
    });

    const contentType = response.headers.get('content-type');
    const isJson = contentType?.includes('application/json');

    if (!response.ok) {
        const error = isJson ? await response.json() : {message: response.statusText};
        throw new APIException(
            error.error?.message ?? error.message ?? 'Request failed',
            response.status,
            error
        );
    }

    if (response.status === 204) {
        return null as T;
    }

    return isJson ? response.json() : null as T;
}

export const api = {
    registration: {
        refresh: async (): Promise<RegistrationToken> => request<RegistrationToken>('/admin/registration/refresh', {method: 'POST'}),
        get: async (): Promise<RegistrationToken> => request<RegistrationToken>('/admin/registration', {method: 'GET'}),
        deleteAll: async (): Promise<void> => request<void>('/admin/registration', {method: 'DELETE'}),
    }
}