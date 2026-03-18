import {fetchGraphs} from '$lib/api/graph.api';
import {getAllGraphs, putGraphs} from '$lib/db/stores/graph';
import type {GraphSummary} from '$lib/types/graph';
import {isHttpError} from "@sveltejs/kit";

export async function listGraphs(): Promise<GraphSummary[]> {
	if (navigator.onLine) {
		try {
			const graphs = await fetchGraphs();
			await putGraphs(graphs);
		} catch (e) {
			if (isHttpError(e) && e.status !== 503) {
				throw e;
			}

			console.warn('Failed to fetch graphs from server, falling back to cache', e);
		}
	}

	return await getAllGraphs();
}
