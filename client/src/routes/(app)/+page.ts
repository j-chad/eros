import { listGraphs } from '$lib/services/graph';
import type { GraphSummary } from '$lib/types/graph';

export async function load(): Promise<{ graphs: GraphSummary[] }> {
	const graphs = await listGraphs();
	return { graphs };
}
