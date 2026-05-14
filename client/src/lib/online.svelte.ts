import { PUBLIC_SERVER_URL as API_BASE } from '$env/static/public';

const PING_INTERVAL = 30_000;
const PING_TIMEOUT = 5_000;

/**
 * Server-reachability state. `true` means the backend responded to a health
 * check, `false` means a network error occurred (fetch threw). HTTP error
 * responses (4xx/5xx) still count as "online" — the server responded.
 *
 * Starts `true` optimistically so the first load attempts the API. The initial
 * ping will correct this within milliseconds if the server is actually down.
 */
let onlineState = $state(true);

/** Whether the initial health ping has resolved. */
let initialPingDone = $state(false);

async function ping(): Promise<boolean> {
	// No point hitting the network if the browser itself is offline.
	if (typeof navigator !== 'undefined' && !navigator.onLine) {
		onlineState = false;
		return false;
	}

	try {
		const controller = new AbortController();
		const timeout = setTimeout(() => controller.abort(), PING_TIMEOUT);

		const res = await fetch(`${API_BASE}/health`, {
			method: 'GET',
			signal: controller.signal,
		});
		clearTimeout(timeout);

		// Any HTTP response means the server is reachable.
		if (res.ok) {
			onlineState = true;
			return true;
		}

		// Non-OK but still a response — server is reachable, don't flip offline.
		return true;
	} catch {
		// Network error (connection refused, DNS failure, timeout, abort).
		onlineState = false;
		return false;
	}
}

if (typeof window !== 'undefined') {
	// Flip offline immediately when the browser reports no network.
	window.addEventListener('offline', () => {
		onlineState = false;
	});

	// When the browser comes back online, ping to verify the server is up too.
	window.addEventListener('online', async () => {
		await ping();
	});

	// Initial ping — fire immediately, mark done when settled.
	ping().finally(() => {
		initialPingDone = true;
	});

	// Recurring pings.
	setInterval(ping, PING_INTERVAL);
}

/**
 * Whether the server is currently reachable.
 */
export function isOnline(): boolean {
	return onlineState;
}

/**
 * Whether the initial health ping has completed. Useful for gating UI that
 * depends on knowing the real connectivity state (e.g. maintenance banner
 * auto-dismiss).
 */
export function isInitialPingDone(): boolean {
	return initialPingDone;
}

/**
 * Called by `http.ts` when a `fetch()` throws a network error, so the app
 * immediately knows the server is unreachable without waiting for the next
 * scheduled ping.
 */
export function markOffline(): void {
	onlineState = false;
}
