import {request} from "$lib/api/http";

export async function fetchVapidKey(): Promise<string> {
    return await request<string>(`/push/vapid-key`, undefined, false);
}

export async function subscribePush(subscription: PushSubscriptionJSON): Promise<void> {
    await request(`/push/subscribe`, {
        method: 'POST',
        body: JSON.stringify(subscription),
        headers: {
            'Content-Type': 'application/json'
        }
    });
}

export async function unsubscribePush(): Promise<void> {
    await request(`/push/subscribe`, {
        method: 'DELETE'
    })
}