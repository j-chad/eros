import { api } from '$lib/api';

export async function load() {
    const subscriptions = await api.push.listSubscriptions();
    return { subscriptions };
}
