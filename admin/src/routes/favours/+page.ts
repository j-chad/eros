import { api } from "$lib/api";
import { browser } from "$app/environment";

export async function load() {
    if (!browser) {
        return {
            choices: []
        };
    }

    const [choices, favourCount] = await Promise.all([
        api.favours.listChoices(),
        api.favours.getFavourCount()
    ]);

    return { choices, favourCount };
}