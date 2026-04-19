import { listGraphs } from '$lib/services/graph';
import { getCount, listRequests } from '$lib/services/favour';
import type { GraphSummary } from '$lib/types/graph';
import type { FavourCount, FavourRequest } from '$lib/types/favour';

export async function load(): Promise<{
	graphs: GraphSummary[];
	favourCount: FavourCount;
	favourRequests: FavourRequest[];
}> {
	const [graphs, favourResult] = await Promise.all([
		listGraphs(),
		Promise.all([getCount(), listRequests()]).catch(() => null),
	]);

	return {
		graphs,
		favourCount: favourResult?.[0] ?? { total: 0, remaining: 0 },
		favourRequests: favourResult?.[1] ?? [],
	};
}
