import { error } from '@sveltejs/kit';
import { api } from '$lib/api';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = async ({ params }) => {
    const id = params.id;

    if (!id) {
        throw error(400, 'Graph ID is required');
    }

    try {
        const graph = await api.graph.get(id);
        return { graph };
    } catch (err) {
        console.error('Failed to load graph:', err);
        throw error(404, 'Graph not found');
    }
};