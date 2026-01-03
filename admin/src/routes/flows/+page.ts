import { api } from "$lib/api";
import { browser } from "$app/environment";

export async function load() {
    if (!browser) {
        return {
            flows: []
        };
    }

    try {
        const flows = [
            {
                id: 'flow1',
                title: 'Sample Flow 1',
                description: 'This is a sample flow description.',
                startingAt: new Date().toISOString(),
            }
        ] as any[]
        return { flows };
    } catch (error) {
        console.error('Failed to load flows:', error);
        return { flows: [] };
    }
}