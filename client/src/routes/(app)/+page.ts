import { browser } from '$app/environment';
import { listGraphs } from '$lib/services/graph';
import type { GraphSummary } from '$lib/types/graph';

export async function load(): Promise<{ graph: GraphSummary | null }> {
	if (!browser) return { graph: null };

	const graphs = await listGraphs();

	if (!graphs.length) return { graph: null };

	// Pick the soonest upcoming graph; fall back to the most recent past one.
	const now = Date.now();
	const upcoming = graphs
		.filter((g) => g.starting_at.getTime() > now)
		.sort((a, b) => a.starting_at.getTime() - b.starting_at.getTime());

	return { graph: upcoming[0] ?? graphs[graphs.length - 1] };
}
