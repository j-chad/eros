import { db, promisifyRequest } from '$lib/db/db';
import type { StoredGraph } from '$lib/db/schema';

export async function getAllGraphs(): Promise<StoredGraph[]> {
	const store = await db.getStore('graphs', 'readonly');
	return promisifyRequest(store.getAll());
}

export async function putGraphs(graphs: StoredGraph[]): Promise<void> {
	const tx = await db.transaction('graphs', 'readwrite');
	const store = tx.objectStore('graphs');
	await Promise.all(graphs.map((g) => promisifyRequest(store.put(g))));
}
