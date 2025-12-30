import {api} from "$lib/api";
import {browser} from "$app/environment";
import type {Favour, FavourChoice, FavourCount} from "$lib/types";

interface FavourPageData {
    choices: Array<FavourChoice>
    favourCount: FavourCount;
    requests: Array<Favour>
}

export async function load(): Promise<FavourPageData> {
    if (!browser) {
        return {
            choices: [],
            favourCount: {total: 0, remaining: 0},
            requests: []
        };
    }

    const [choices, favourCount, requests] = await Promise.all([
        api.favours.listChoices(),
        api.favours.getFavourCount(),
        api.favours.listFavourRequests(),
    ]);

    return {choices, favourCount, requests};
}