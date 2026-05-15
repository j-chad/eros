import { db, promisifyRequest } from '$lib/db/db';
import type { GraphDetail, GraphSummary } from '$lib/types/graph';

function hydrateGraphSummary(g: GraphSummary): GraphSummary {
	return {
		...g,
		starting_at: g.starting_at instanceof Date ? g.starting_at : new Date(g.starting_at as unknown as string),
		created_at: g.created_at instanceof Date ? g.created_at : new Date(g.created_at as unknown as string),
		updated_at: g.updated_at instanceof Date ? g.updated_at : new Date(g.updated_at as unknown as string),
	};
}

export async function getAllGraphs(): Promise<GraphSummary[]> {
	const store = await db.getStore('graphs', 'readonly');
	const graphs = await promisifyRequest(store.getAll());
	return graphs.map(hydrateGraphSummary);
}

export async function putGraphs(graphs: GraphSummary[]): Promise<void> {
	const tx = await db.transaction('graphs', 'readwrite');
	const store = tx.objectStore('graphs');
	await Promise.all(graphs.map((g) => promisifyRequest(store.put(g))));
}

export async function getGraphDetail(id: string): Promise<GraphDetail | undefined> {
	const store = await db.getStore('graphDetails', 'readonly');
	return promisifyRequest(store.get(id));
}

export async function putGraphDetail(graph: GraphDetail): Promise<void> {
	const store = await db.getStore('graphDetails', 'readwrite');
	await promisifyRequest(store.put(graph));
}
