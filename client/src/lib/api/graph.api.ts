import { request } from '$lib/api/http';
import type {GraphSummary} from "$lib/types/graph";

type GraphSummaryResponse = Omit<GraphSummary, 'starting_at' | 'created_at' | 'updated_at'> & {
	starting_at: string;
	created_at: string;
	updated_at: string;
};

// Coerce stored ISO date strings to Date objects for consumers.
function hydrate(g: GraphSummaryResponse): GraphSummary {
	return {
		...g,
		starting_at: new Date(g.starting_at),
		created_at: new Date(g.created_at),
		updated_at: new Date(g.updated_at),
	};
}

export async function fetchGraphs(): Promise<GraphSummary[]> {
	const graphs = await request<GraphSummaryResponse[]>('/graphs');
	return graphs.map(hydrate);
}
