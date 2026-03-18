import { fetchGraphs } from '$lib/api/graph.api';
import { getAllGraphs, putGraphs } from '$lib/db/stores/graph';
import type { GraphSummary } from '$lib/types/graph';

export async function listGraphs(): Promise<GraphSummary[]> {
	if (navigator.onLine) {
		try {
			const graphs = await fetchGraphs();
			await putGraphs(graphs);
		} catch (e) {
			console.warn('Failed to fetch graphs from server, falling back to cache', e);
		}
	}

	return await getAllGraphs();
}
