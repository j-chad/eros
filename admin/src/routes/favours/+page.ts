import { api } from "$lib/api";
import { browser } from "$app/environment";

export async function load() {
    if (!browser) {
        return {
            choices: []
        };
    }

    const choices = await api.favours.listChoices();

    return { choices };
}