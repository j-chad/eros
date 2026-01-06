import {browser} from "$app/environment";
import { api } from "$lib/api";

export async function load() {
    if (!browser) {
        return {
            graphs: []
        };
    }

    try {
        const graphs = await api.graph.list();
        return { graphs: graphs };
    } catch (error) {
        console.error('Failed to load graphs:', error);
        return { graphs: [] };
    }
}