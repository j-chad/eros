import type {
    APIError,
    CreateFavourChoicePayload,
    Device,
    Favour,
    FavourChoice,
    FavourCount,
    RegistrationToken
} from "$lib/types";
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

async function request<T = void>(endpoint: string, options: RequestInit = {}): Promise<T> {
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
    ping: async (): Promise<void> => request('/admin/ping', {method: 'GET'}),
    registration: {
        create: async () => request<RegistrationToken>('/admin/registration-codes', {method: 'POST'}),
        get: async () => request<RegistrationToken>('/admin/registration-codes', {method: 'GET'}),
        deleteAll: async () => request('/admin/registration-codes', {method: 'DELETE'}),
    },
    devices: {
        list: async () => request<Array<Device>>('/admin/devices', {method: 'GET'}),
        delete: async (deviceId: string) => request(`/admin/devices/${deviceId}`, {method: 'DELETE'}),
        updateDeviceInfo: async (deviceId: string, deviceInfo: string) =>
            request(`/admin/devices/${deviceId}`, {
                method: 'PATCH',
                body: deviceInfo,
                headers: {
                    'Content-Type': 'text/plain',
                }
            }),
    },
    favours: {
        listChoices: async () => request<Array<FavourChoice>>('/favours/choices', {method: 'GET'}),
        createChoice: async (choice: CreateFavourChoicePayload) =>
            request<FavourChoice>('/admin/favours/choices', {
                method: 'POST',
                body: JSON.stringify(choice),
            }),
        deleteChoice: async (choiceId: string) => request(`/admin/favours/choices/${choiceId}`, {method: 'DELETE'}),
        updateChoice: async (choiceId: string, choice: FavourChoice) =>
            request<FavourChoice>(`/admin/favours/choices/${choiceId}`, {
                method: 'PUT',
                body: JSON.stringify(choice),
            }),
        getFavourCount: async () => request<FavourCount>('/favours', {method: 'GET'}),
        updateFavourCount: async (total: number) =>
            request('/admin/favours', {
                method: 'PUT',
                body: JSON.stringify(total),
            }),
        updateFavourRequestFulfilment: async (favourId: string, fulfilled: boolean) =>
            request(`/admin/favours/requests/${favourId}`, {
                method: 'PATCH',
                body: JSON.stringify({fulfilled}),
            }),
        listFavourRequests: async () => request<Array<Favour>>('/favours/requests', {method: 'GET'}),
    }
}