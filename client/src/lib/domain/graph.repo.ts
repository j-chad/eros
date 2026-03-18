import {request} from "$lib/api/http";
import type {GraphSummary} from "$lib/domain/graph.types";

export async function list(): Promise<GraphSummary[]> {
	return await request<GraphSummary[]>('/graphs');
}
