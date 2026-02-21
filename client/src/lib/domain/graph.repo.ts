import {request} from "$lib/api/http";
import type {GraphSummary} from "$lib/domain/graph.types";

export async function list(): Promise<GraphSummary[]> {
	const graphs = await request<GraphSummary[]>('/api/graphs');
}
