import { fetchGraphs } from '$lib/api/graph.api';
import { getAllGraphs, putGraphs } from '$lib/db/stores/graph';
import type { GraphSummary } from '$lib/types/graph';

// Coerce stored ISO date strings to Date objects for consumers.
function hydrate(g: { starting_at: string; created_at: string; updated_at: string } & Omit<GraphSummary, 'starting_at' | 'created_at' | 'updated_at'>): GraphSummary {
	return {
		...g,
		starting_at: new Date(g.starting_at),
		created_at: new Date(g.created_at),
		updated_at: new Date(g.updated_at),
	};
}

export async function listGraphs(): Promise<GraphSummary[]> {
	if (navigator.onLine) {
		try {
			const graphs = await fetchGraphs();
			await putGraphs(graphs);
			return graphs.map(hydrate);
		} catch {
			// Fall through to cache on fetch failure.
		}
	}

	const cached = await getAllGraphs();
	return cached.map(hydrate);
}
