import { PUBLIC_SERVER_URL as API_BASE } from '$env/static/public';
import { loadToken, clearToken } from '$lib/api/auth';
import { goto } from '$app/navigation';

const PING_INTERVAL = 30_000;
const PING_TIMEOUT = 5_000;

/**
 * The reason the app last went offline.
 */
export enum OfflineReason {
	browserOffline = "browser-offline",
	pingFailed = "ping-failed",
	pingTimeout = "ping-timeout",
	apiError = "api-error",
	forced = "forced",
}

/** Whether the initial health ping has resolved. */
let initialPingDone = $state(false);

/** Why the app is currently considered offline. `null` = online. */
let onlineState = $state<OfflineReason | null>(null);

function setOffline(reason: OfflineReason) {
	onlineState = reason;
}

function setOnline() {
	onlineState = null;
}

async function ping(): Promise<boolean> {
	// Debug override — skip real pings entirely.
	if (onlineState === OfflineReason.forced) {
		return false;
	}

	// No point hitting the network if the browser itself is offline.
	if (typeof navigator !== 'undefined' && !navigator.onLine) {
		setOffline(OfflineReason.browserOffline);
		return false;
	}

	// Use the authenticated /ping endpoint when we have a token, so the
	// server can tell us if the session has been revoked. Fall back to the
	// unauthenticated /health endpoint otherwise (e.g. on the login page).
	const token = await loadToken();
	const endpoint = token ? '/ping' : '/health';
	const headers: HeadersInit = {};
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	try {
		const controller = new AbortController();
		const timeout = setTimeout(() => controller.abort(), PING_TIMEOUT);

		const res = await fetch(`${API_BASE}${endpoint}`, {
			method: 'GET',
			headers,
			signal: controller.signal,
		});
		clearTimeout(timeout);

		// Session revoked — clear local state and redirect to login.
		if (res.status === 401) {
			await clearToken();
			await goto('/login?reason=session_expired');
			return false;
		}

		// Any HTTP response means the server is reachable.
		setOnline();
		return true;
	} catch (err) {
		if (err instanceof DOMException && err.name === 'AbortError') {
			setOffline(OfflineReason.pingTimeout);
		} else {
			setOffline(OfflineReason.pingFailed);
		}
		return false;
	}
}

if (typeof window !== 'undefined') {
	// Flip offline immediately when the browser reports no network.
	window.addEventListener('offline', () => {
		setOffline(OfflineReason.browserOffline);
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
	return onlineState === null;
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
	setOffline(OfflineReason.apiError);
}

/**
 * The current reason the app considers itself offline. `null` when online.
 */
export function getOfflineReason(): OfflineReason | null {
	return onlineState;
}

/**
 * Force the app offline (or release the override). When releasing, immediately
 * pings to restore real connectivity state.
 */
export function setOfflineOverride(force: boolean): void {
	if (force) {
		setOffline(OfflineReason.forced);
	} else {
		// Restore real state by pinging immediately.
		setOffline(OfflineReason.pingFailed)
		void ping();
	}
}

/**
 * Whether the debug offline override is currently active.
 */
export function isOfflineOverridden(): boolean {
	return onlineState === OfflineReason.forced;
}
