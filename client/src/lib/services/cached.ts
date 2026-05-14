import { ServerUnreachableError } from '$lib/api/http';
import { isOnline } from '$lib/online.svelte';

/**
 * Fetch-with-cache-fallback helper. If online, tries `fetchFn` (which should
 * fetch from the API and write through to IndexedDB) and returns the result.
 * If the server is unreachable — either because we already know we're offline,
 * or because the fetch threw a `ServerUnreachableError` — falls back to
 * `cacheFn` (which reads from IndexedDB).
 *
 * Any non-network error (4xx, 5xx, etc.) is re-thrown as-is.
 */
export async function cached<T>(fetchFn: () => Promise<T>, cacheFn: () => Promise<T>): Promise<T> {
	if (!isOnline()) {
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
