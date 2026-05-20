export function isPushSupported() {
  return 'serviceWorker' in navigator && 'PushManager' in window && 'Notification' in window;
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