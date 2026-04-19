import {
	fetchFavourChoices,
	fetchFavourCount,
	fetchFavourRequests,
	requestFavour as requestFavourAPI,
} from '$lib/api/favour.api';
import {
	getAllFavourChoices,
	putFavourChoices,
	getAllFavourRequests,
	putFavourRequests,
	getFavourCount,
	putFavourCount,
} from '$lib/db/stores/favour';
import type { FavourChoice, FavourCount, FavourRequest } from '$lib/types/favour';

export async function listChoices(): Promise<FavourChoice[]> {
	if (navigator.onLine) {
		const choices = await fetchFavourChoices();
		await putFavourChoices(choices);
		return choices;
	}

	return await getAllFavourChoices();
}

export async function getCount(): Promise<FavourCount> {
	if (navigator.onLine) {
		const count = await fetchFavourCount();
		await putFavourCount(count);
		return count;
	}

	return (await getFavourCount()) ?? { total: 0, remaining: 0 };
}

export async function listRequests(): Promise<FavourRequest[]> {
	if (navigator.onLine) {
		const requests = await fetchFavourRequests();
		await putFavourRequests(requests);
		return requests;
	}

	return await getAllFavourRequests();
}

export async function requestFavour(
	choiceId: string,
	message?: string
): Promise<FavourRequest> {
	const result = await requestFavourAPI(choiceId, message);

	// Update the cached requests and count after a successful request.
	const [requests, count] = await Promise.all([fetchFavourRequests(), fetchFavourCount()]);
	await Promise.all([putFavourRequests(requests), putFavourCount(count)]);

	return result;
}
