import {fetchGraphs} from '$lib/api/graph.api';
import {getAllGraphs, putGraphs} from '$lib/db/stores/graph';
import type {GraphSummary} from '$lib/types/graph';
import {isHttpError} from "@sveltejs/kit";

export async function listGraphs(): Promise<GraphSummary[]> {
	if (navigator.onLine) {
		const graphs = await fetchGraphs();
		await putGraphs(graphs);
		return graphs
	}

	return await getAllGraphs();
}
