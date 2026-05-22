import { listGraphs, getGraph } from '$lib/services/graph';
import { listChoices, getCount, listRequests } from '$lib/services/favour';
import { isOnline } from '$lib/online.svelte';
import {isPushSubscribed, republishSubscription} from "$lib/services/push";

const SYNC_INTERVAL_MS = 5 * 60 * 1000; // 5 minutes

export type SyncLogger = {
	step: (msg: string) => void;
	result: (msg: string) => void;
	error: (msg: string) => void;
};

/**
 * Re-fetch all cached data from the server and write it back to IndexedDB.
 * Delegates entirely to the existing service functions, which already handle
 * the fetch-and-cache pattern. Does nothing if the server is currently unreachable.
 *
 * Pass a `logger` to receive per-step progress (used by the debug page terminal).
 * Without a logger, errors are written to the console only.
 */
export async function syncAll(logger?: SyncLogger): Promise<void> {
	if (!isOnline()) return;

	try {
		if (await isPushSubscribed()) {
			logger?.step('republishing push subscription')
			await republishSubscription()
			logger?.result('ok');
		}

		logger?.step('fetch graphs');
		const graphs = await listGraphs();
		logger?.result(`${graphs.length} graph(s)`);

		// Only fetch details for graphs that have already started — future graphs
		// have no accessible nodes yet so fetching them is wasted work.
		const pastGraphs = graphs.filter((g) => g.starting_at <= new Date());

		for (const g of pastGraphs) {
			logger?.step(`fetch graph "${g.title}"`);
			await getGraph(g.id);
			logger?.result('ok');
		}

		logger?.step('fetch favour choices');
		const choices = await listChoices();
		logger?.result(`${choices.length} choice(s)`);

		logger?.step('fetch favour count');
		await getCount();
		logger?.result('ok');

		logger?.step('fetch favour requests');
		const requests = await listRequests();
		logger?.result(`${requests.length} request(s)`);

		if (!logger) console.debug('[sync] Background refresh complete');
	} catch (e) {
		if (logger) {
			logger.error(e instanceof Error ? e.message : String(e));
		} else {
			console.warn('[sync] Background refresh error:', e);
		}
		throw e;
	}
}

let syncStarted = false;

/**
 * Start the background sync process. Call once from the root layout on mount.
 * Safe to call multiple times — subsequent calls are no-ops.
 *
 * - Syncs on reconnect: listens for the browser `online` event and fires a sync
 *   1.5 s later (giving the health ping time to confirm the server is reachable).
 * - Syncs periodically: runs every 5 minutes while the app is open.
 */
export function startSync(): void {
	if (typeof window === 'undefined' || syncStarted) return;
	syncStarted = true;

	window.addEventListener('online', () => {
		// Small delay so online.svelte.ts has time to ping and confirm reachability.
		setTimeout(() => syncAll(), 1500);
	});

	setInterval(syncAll, SYNC_INTERVAL_MS);

	void syncAll();
}
