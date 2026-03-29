import { db, promisifyRequest } from '$lib/db/db';
import type { GraphDetail, GraphSummary } from '$lib/types/graph';

export async function getAllGraphs(): Promise<GraphSummary[]> {
	const store = await db.getStore('graphs', 'readonly');
	return promisifyRequest(store.getAll());
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
