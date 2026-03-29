import { fetchGraph, fetchGraphs } from '$lib/api/graph.api';
import { getAllGraphs, getGraphDetail, putGraphDetail, putGraphs } from '$lib/db/stores/graph';
import type { GraphDetail, GraphSummary } from '$lib/types/graph';

export async function listGraphs(): Promise<GraphSummary[]> {
	if (navigator.onLine) {
		const graphs = await fetchGraphs();
		await putGraphs(graphs);
		return graphs;
	}

	return await getAllGraphs();
}

export async function getGraph(id: string): Promise<GraphDetail | null> {
	if (navigator.onLine) {
		const graph = await fetchGraph(id);
		await putGraphDetail(graph);
		return graph;
	}

	return (await getGraphDetail(id)) ?? null;
}
