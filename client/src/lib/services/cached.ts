import { ServerUnreachableError } from '$lib/api/http';
import { isOnline, isInitialPingDone } from '$lib/online.svelte';

/**
 * Fetch-with-cache-fallback helper.
 *
 * During early startup (before the initial health ping has settled) we don't
 * yet know whether the server is reachable. Rather than blocking the page on a
 * network request that may hang, we return cached data immediately so the UI
 * can render. The background sync cycle will refresh the cache once
 * connectivity is established.
 *
 * After the initial ping has resolved we know the connectivity state, so we
 * use the normal network-first strategy: try `fetchFn`, fall back to `cacheFn`
 * on `ServerUnreachableError`.
 */
export async function cached<T>(fetchFn: () => Promise<T>, cacheFn: () => Promise<T>): Promise<T> {
	if (!isInitialPingDone() || !isOnline()) {
		return cacheFn();
	}

	try {
		return await fetchFn();
	} catch (e) {
		if (e instanceof ServerUnreachableError) {
			return cacheFn();
		}
		throw e;
	}
}
