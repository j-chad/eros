import {browser} from "$app/environment";
import { api } from "$lib/api";

export async function load() {
    if (!browser) {
        return {
            nodes: []
        };
    }

    try {
        const nodes = await api.graph.list();
        return { nodes };
    } catch (error) {
        console.error('Failed to load graphs:', error);
        return { nodes: [] };
    }
}