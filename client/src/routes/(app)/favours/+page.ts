import { listChoices, getCount, listRequests } from '$lib/services/favour';
import type { FavourChoice, FavourCount, FavourRequest } from '$lib/types/favour';

export async function load(): Promise<{
	choices: FavourChoice[];
	count: FavourCount;
	requests: FavourRequest[];
}> {
	const [choices, count, requests] = await Promise.all([
		listChoices(),
		getCount(),
		listRequests(),
	]);
	return { choices, count, requests };
}
