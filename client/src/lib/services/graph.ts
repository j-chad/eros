import { fetchGraph, fetchGraphs, unlockNode as unlockNodeAPI, type UnlockResult } from '$lib/api/graph.api';
import { getAllGraphs, getGraphDetail, putGraphDetail, putGraphs } from '$lib/db/stores/graph';
import { cached } from './cached';
import type { GraphDetail, GraphSummary } from '$lib/types/graph';

export async function listGraphs(): Promise<GraphSummary[]> {
	return cached(
		async () => {
			const graphs = await fetchGraphs();
			await putGraphs(graphs);
			return graphs;
		},
		() => getAllGraphs(),
	);
}

export async function getGraph(id: string): Promise<GraphDetail | null> {
	return cached(
		async () => {
			const graph = await fetchGraph(id);
			await putGraphDetail(graph);
			return graph;
		},
		async () => (await getGraphDetail(id)) ?? null,
	);
}

export async function unlockNode(graphId: string, nodeId: string, payload: string): Promise<UnlockResult> {
	const cachedGraph = await getGraphDetail(graphId);
	if (!cachedGraph) {
		throw new Error('Graph not found in cache');
	}

	const result = await unlockNodeAPI(nodeId, payload);

	// Merge the result into the cached graph detail in IndexedDB
	const nodeIndex = cachedGraph.nodes.findIndex(n => n.id === result.unlocked_node.id);
	if (nodeIndex >= 0) {
		cachedGraph.nodes[nodeIndex] = result.unlocked_node;
	}

	cachedGraph.nodes.push(...result.new_nodes);
	cachedGraph.edges.push(...result.new_edges);

	await putGraphDetail(cachedGraph);

	return result;
}
