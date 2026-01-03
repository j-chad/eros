import {browser} from "$app/environment";
import {NodeType, type StartNode} from "$lib/types";

export async function load() {
    if (!browser) {
        return {
            nodes: []
        };
    }

    try {
        const nodes: StartNode[] = [
            {
                type: NodeType.START,
                id: 'flow1',
                title: 'Sample Flow gjgjhgjg1',
                description: 'This is a sample flow description.',
                starting_at: new Date().toISOString(),
                unlocked_at: null,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
            },
            {
                id: 'flow2',
                title: 'Sample Flow 2',
                description: 'This is a sample flow description.',
                starting_at: new Date().toISOString(),
                type: NodeType.START,
                unlocked_at: null,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
            }
        ]
        return { nodes };
    } catch (error) {
        console.error('Failed to load flows:', error);
        return { nodes: [] };
    }
}