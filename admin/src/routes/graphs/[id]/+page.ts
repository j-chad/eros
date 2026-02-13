import { error } from '@sveltejs/kit';
import { api } from '$lib/api';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = async ({ params }) => {
    const id = params.id;

    if (!id) {
        throw error(400, 'Graph ID is required');
    }

	return { graph: await api.graph.get(id) };
};
