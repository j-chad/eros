import { api } from "$lib/api";
import { browser } from "$app/environment";

export async function load() {
    if (!browser) {
        return {
            favours: []
        };
    }

    const favours = await api.favours.listChoices();

    return { favours };
}