import {request} from "$lib/api/http";
import {base64URLDecode} from "$lib/utils/base64";

export async function fetchVapidKey(): Promise<string> {
    const encoded = await request<string>(`/push/vapid-key`, undefined, false);
    return base64URLDecode(encoded);
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