import {cached} from "$lib/services/cached";
import {fetchVapidKey, subscribePush} from "$lib/api/push.api";
import {KVKey, KVStore} from "$lib/db/stores/kv";

export function isPushSupported() {
  if (!('serviceWorker' in navigator && 'PushManager' in window && 'Notification' in window)) return false;
  return Notification?.permission !== 'denied';
}

export function isPushEnabled() {
    if (!isPushSupported()) return false;
    return Notification?.permission === "granted";
}

export async function isPushSubscribed() {
    const registration = await navigator.serviceWorker.ready;
    const subscription = await registration.pushManager.getSubscription()
    if (subscription === null) {
        return false
    }

    if (subscription.expirationTime == null) {
        return true
    }

    return subscription.expirationTime > Date.now()
}

export async function subscribe() {
    if (!isPushSupported()) return;
    if (await isPushSubscribed()) return;

    const key = await getVAPIDPublicKey();
    if (key === null) return;

    const registration = await navigator.serviceWorker.ready;
    const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: key
    })

    const payload = subscription.toJSON();
    await subscribePush(payload)
}

async function getVAPIDPublicKey() {
    return cached(async () => {
        const key = await fetchVapidKey();
        await KVStore.set(KVKey.PushVAPIDKey, key);
        return key
    }, () => KVStore.get(KVKey.PushVAPIDKey))
}