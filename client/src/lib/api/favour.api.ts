import { request } from '$lib/api/http';
import type { FavourChoice, FavourCount, FavourRequest } from '$lib/types/favour';

export async function fetchFavourChoices(): Promise<FavourChoice[]> {
	return request<FavourChoice[]>('/favours/choices');
}

export async function fetchFavourCount(): Promise<FavourCount> {
	return request<FavourCount>('/favours');
}

export async function fetchFavourRequests(): Promise<FavourRequest[]> {
	return request<FavourRequest[]>('/favours/requests');
}

export async function requestFavour(
	choiceId: string,
	message?: string
): Promise<FavourRequest> {
	return request<FavourRequest>('/favours/request', {
		method: 'POST',
		body: JSON.stringify({
			choice_id: choiceId,
			...(message ? { message } : {}),
		}),
	});
}
