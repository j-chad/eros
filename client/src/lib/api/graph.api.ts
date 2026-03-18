import { request } from '$lib/api/http';
import type { StoredGraph } from '$lib/db/schema';

export async function fetchGraphs(): Promise<StoredGraph[]> {
	return request<StoredGraph[]>('/graphs');
}
