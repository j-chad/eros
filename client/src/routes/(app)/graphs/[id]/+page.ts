import { error } from '@sveltejs/kit';
import { getGraph } from '$lib/services/graph';
import type { GraphDetail } from '$lib/types/graph';

export async function load({ params }): Promise<{ graph: GraphDetail }> {
	const graph = await getGraph(params.id);
	if (!graph) {
		throw error(404, { message: 'Graph not found' });
	}

	return { graph };
}
