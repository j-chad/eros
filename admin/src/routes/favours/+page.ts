import {api} from "$lib/api";
import {browser} from "$app/environment";
import type {Favour} from "$lib/types";

export async function load() {
    if (!browser) {
        return {
            choices: [],
            favourCount: 0,
            requests: []
        };
    }

    const [choices, favourCount] = await Promise.all([
        api.favours.listChoices(),
        api.favours.getFavourCount(),
    ]);

    const requests: Favour[] = [{
        id: '1',
        choice: 'Help with groceries',
        message: 'Can someone help me pick up groceries this weekend?',
        requestedAt: '2024-06-01T10:00:00Z',
        fulfilled: false
    },
        {
            id: '2',
            choice: 'Dog walking',
            message: 'Need a hand walking my dog in the evenings.',
            requestedAt: '2024-05-28T15:30:00Z',
            fulfilled: true
        }]

    return {choices, favourCount, requests};
}