<script lang="ts">
	import {db, promisifyRequest} from '$lib/db/db';
	import {DB_NAME} from '$lib/db/schema';
	import {KVKey} from '$lib/db/stores/kv';
	import {isOnline} from '$lib/online.svelte';
	import {getGraph, listGraphs} from '$lib/services/graph';
	import {getCount, listChoices, listRequests} from '$lib/services/favour';
	import {clearToken} from '$lib/api/auth';
	import {goto} from '$app/navigation';
	import {ChevronLeft, ChevronUp} from 'lucide-svelte';
	import BrandHeader from '$lib/ui/BrandHeader.svelte';

	type LogType = 'command' | 'result' | 'error' | 'summary';
	type LogEntry = { type: LogType; text: string };

	// --- Log state ---
	let log = $state<LogEntry[]>([]);
	let clearing = $state(false);
	let syncing = $state(false);
	let loggingOut = $state(false);
	let desktopTerminalEl: HTMLDivElement | undefined = $state();
	let mobileTerminalEl: HTMLDivElement | undefined = $state();

	const online = $derived(isOnline());
	const busy = $derived(clearing || syncing || loggingOut);

	const lastEntry = $derived(log.length > 0 ? log[log.length - 1] : null);
	const lastEntryHint = $derived(
		lastEntry
			? lastEntry.type === 'command'
				? `> ${lastEntry.text}`
				: lastEntry.text
			: ''
	);

	// Auto-scroll both terminals when log changes
	$effect(() => {
		if (log.length) {
			if (desktopTerminalEl) {
				desktopTerminalEl.scrollTop = desktopTerminalEl.scrollHeight;
			}
			if (mobileTerminalEl) {
				mobileTerminalEl.scrollTop = mobileTerminalEl.scrollHeight;
			}
		}
	});

	// --- Drawer state ---
	const SNAP_COLLAPSED = 48;
	const SNAP_HALF_VH = 40;
	const SNAP_FULL_VH = 80;

	type SnapPoint = 'collapsed' | 'half' | 'full';
	let drawerSnap = $state<SnapPoint>('collapsed');
	let isDragging = $state(false);
	let drawerHeight = $state(SNAP_COLLAPSED);

	// Recalculate drawer height when the viewport resizes (e.g. rotating device,
	// or dragging browser edge between mobile/desktop breakpoints).
	$effect(() => {
		function onResize() {
			// Re-snap to the current snap point with the new viewport size
			drawerHeight = snapHeightPx(drawerSnap);
		}
		window.addEventListener('resize', onResize);
		return () => window.removeEventListener('resize', onResize);
	});

	function snapHeightPx(snap: SnapPoint): number {
		if (typeof window === 'undefined') return SNAP_COLLAPSED;
		const vh = window.innerHeight;
		switch (snap) {
			case 'collapsed':
				return SNAP_COLLAPSED;
			case 'half':
				return vh * (SNAP_HALF_VH / 100);
			case 'full':
				return vh * (SNAP_FULL_VH / 100);
		}
	}

	function snapTo(snap: SnapPoint) {
		drawerSnap = snap;
		drawerHeight = snapHeightPx(snap);
	}

	function handleHandleTap() {
		if (drawerSnap === 'collapsed') {
			snapTo('half');
		} else {
			snapTo('collapsed');
		}
	}

	// --- Touch gesture handling ---
	let touchStartY = 0;
	let touchStartHeight = 0;

	function handleTouchStart(e: TouchEvent) {
		isDragging = true;
		touchStartY = e.touches[0].clientY;
		touchStartHeight = drawerHeight;
	}

	function handleTouchMove(e: TouchEvent) {
		if (!isDragging) return;
		e.preventDefault();
		const deltaY = touchStartY - e.touches[0].clientY;
		const maxHeight = snapHeightPx('full');
		drawerHeight = Math.max(SNAP_COLLAPSED, Math.min(maxHeight, touchStartHeight + deltaY));
	}

	function handleTouchEnd(e: TouchEvent) {
		if (!isDragging) return;
		isDragging = false;

		// Calculate velocity for flick detection
		const endY = e.changedTouches[0].clientY;
		const velocity = touchStartY - endY; // positive = swiped up

		const collapsedH = snapHeightPx('collapsed');
		const halfH = snapHeightPx('half');
		const fullH = snapHeightPx('full');

		const FLICK_THRESHOLD = 60;

		// If a strong flick, bias towards the direction
		if (Math.abs(velocity) > FLICK_THRESHOLD) {
			if (velocity > 0) {
				// Swiped up
				if (drawerSnap === 'collapsed' || drawerHeight < halfH) {
					snapTo('half');
				} else {
					snapTo('full');
				}
			} else {
				// Swiped down
				if (drawerSnap === 'full' || drawerHeight > halfH) {
					snapTo('half');
				} else {
					snapTo('collapsed');
				}
			}
			return;
		}

		// Otherwise snap to nearest
		const dCollapsed = Math.abs(drawerHeight - collapsedH);
		const dHalf = Math.abs(drawerHeight - halfH);
		const dFull = Math.abs(drawerHeight - fullH);
		const minD = Math.min(dCollapsed, dHalf, dFull);

		if (minD === dCollapsed) snapTo('collapsed');
		else if (minD === dHalf) snapTo('half');
		else snapTo('full');
	}

	// --- Log helpers ---
	function append(msg: string) {
		log = [...log, { type: 'command', text: msg }];
	}

	function appendResult(msg: string) {
		log = [...log, { type: 'result', text: msg }];
	}

	function appendError(msg: string) {
		log = [...log, { type: 'error', text: msg }];
	}

	function appendSummary(msg: string) {
		log = [...log, { type: 'summary', text: msg }];
	}

	// --- Actions ---
	async function clearDatabase() {
		clearing = true;
		log = [];

		try {
			append('read auth token');
			const kvStore = await db.getStore('kv', 'readonly');
			const tx = kvStore.transaction;
			const authEntry = await promisifyRequest(kvStore.get(KVKey.AuthSession));
			await new Promise<void>((resolve) => {
				tx.oncomplete = () => resolve();
			});

			append('close db');
			await db.close();

			append(`delete "${DB_NAME}"`);
			await new Promise<void>((resolve, reject) => {
				const req = indexedDB.deleteDatabase(DB_NAME);
				req.onsuccess = () => resolve();
				req.onerror = () => reject(req.error);
				req.onblocked = () => appendError('blocked — close other tabs and retry');
			});
			appendResult('ok');

			append('init db');
			await db.init();
			appendResult('ok');

			if (authEntry) {
				append('restore auth token');
				const newKvStore = await db.getStore('kv', 'readwrite');
				await promisifyRequest(newKvStore.put(authEntry));
				appendResult('ok');
			}

			appendSummary('database cleared (auth preserved)');
		} catch (e) {
			appendError(e instanceof Error ? e.message : String(e));
		} finally {
			clearing = false;
		}
	}

	async function syncAll() {
		syncing = true;
		log = [];

		try {
			append('fetch graphs');
			const graphs = await listGraphs();
			appendResult(`${graphs.length} graph(s)`);

			for (const g of graphs) {
				if (g.starting_at > new Date()) {
					append(`skip future graph "${g.id}"`);
					continue;
				}

				append(`fetch graph "${g.title}"`);
				await getGraph(g.id);
				appendResult('ok');
			}

			append('fetch favour choices');
			const choices = await listChoices();
			appendResult(`${choices.length} choice(s)`);

			append('fetch favour count');
			await getCount();
			appendResult('ok');

			append('fetch favour requests');
			const requests = await listRequests();
			appendResult(`${requests.length} request(s)`);

			appendSummary('sync complete');
		} catch (e) {
			appendError(e instanceof Error ? e.message : String(e));
		} finally {
			syncing = false;
		}
	}

	async function handleLogout() {
		loggingOut = true;
		log = [];
		try {
			append('clear auth token');
			await clearToken();
			appendResult('ok');

			appendSummary('logged out — redirecting...');
			await goto('/login');
		} catch (e) {
			appendError(e instanceof Error ? e.message : String(e));
			loggingOut = false;
		}
	}
</script>

{#snippet terminalLines()}
	{#each log as entry, i (i)}
		{#if entry.type === 'command'}
			<p class="text-primary font-semibold mt-2 first:mt-0">
				<span class="opacity-50">&#8250;</span> {entry.text}
			</p>
		{:else if entry.type === 'result'}
			<p class="opacity-60 pl-4">{entry.text}</p>
		{:else if entry.type === 'error'}
			<p class="text-error font-semibold pl-4">{entry.text}</p>
		{:else if entry.type === 'summary'}
			<p class="font-bold mt-3 pt-3 border-t border-base-content/10">{entry.text}</p>
		{/if}
	{/each}
{/snippet}

<svelte:head>
	<title>Debug</title>
</svelte:head>

<div class="mx-auto min-h-dvh max-w-4xl px-4 py-6 flex flex-col gap-6 lg:grid lg:grid-cols-2 lg:gap-8">
	<div class="flex flex-col gap-6">
		<BrandHeader subtitle="Debug tools">
			{#snippet rightContent()}
				<button class="btn btn-ghost btn-sm rounded-2xl gap-1" onclick={() => goto('/')}>
					<ChevronLeft size={16} />
					Back
				</button>
			{/snippet}
		</BrandHeader>

		<div class="flex flex-col gap-2">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Status</p>
			<div class="flex items-center gap-2">
				<span class="inline-block h-2 w-2 rounded-full {online ? 'bg-success' : 'bg-warning'}"></span>
				<span class="text-sm">{online ? 'Server reachable' : 'Server offline'}</span>
			</div>
		</div>

		<div class="flex flex-col gap-3">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Data</p>
			<button
				class="btn btn-error btn-outline rounded-2xl"
				disabled={busy}
				onclick={clearDatabase}
			>
				{#if clearing}
					<span class="loading loading-spinner loading-sm"></span>
				{/if}
				Clear cached data
			</button>

			<button
				class="btn btn-primary rounded-2xl"
				disabled={busy || !online}
				onclick={syncAll}
			>
				{#if syncing}
					<span class="loading loading-spinner loading-sm"></span>
				{/if}
				Sync from server
			</button>
		</div>

		<div class="flex flex-col gap-3">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Session</p>
			<button
				class="btn btn-error btn-outline rounded-2xl"
				disabled={busy}
				onclick={handleLogout}
			>
				{#if loggingOut}
					<span class="loading loading-spinner loading-sm"></span>
				{/if}
				Log out
			</button>
		</div>
	</div>

	<!-- Desktop terminal: sticky side column -->
	{#if log.length > 0}
		<div
			bind:this={desktopTerminalEl}
			class="hidden lg:block bg-base-300 rounded-3xl p-5 font-mono text-xs leading-relaxed
				shadow-[0_2px_12px_0_--theme(--color-pink-200/40)]
				lg:sticky lg:top-6 lg:self-start lg:max-h-[calc(100dvh-3rem)] lg:overflow-y-auto"
		>
			{@render terminalLines()}
		</div>
	{/if}
</div>

<!-- Mobile drawer -->
{#if log.length > 0}
	<div
		class="fixed bottom-0 left-0 right-0 z-50 flex flex-col overflow-hidden
			rounded-t-3xl bg-base-300 will-change-[height] lg:hidden"
		style="height: {drawerHeight}px;
			box-shadow: 0 -4px 24px 0 color-mix(in oklch, var(--color-primary) 12%, transparent),
				0 -1px 4px 0 color-mix(in oklch, var(--color-primary) 6%, transparent);
			transition: {isDragging ? 'none' : 'height 300ms cubic-bezier(0.22, 1, 0.36, 1)'};"
	>
		<!-- Handle bar -->
		<div
			class="drawer-handle"
			role="button"
			tabindex="0"
			onclick={handleHandleTap}
			onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') handleHandleTap(); }}
			ontouchstart={handleTouchStart}
			ontouchmove={handleTouchMove}
			ontouchend={handleTouchEnd}
		>
			<div class="handle-pill"></div>
			{#if drawerSnap === 'collapsed'}
				<div class="handle-hint">
					<ChevronUp size={14} class="opacity-40" />
					<span class="truncate">{lastEntryHint}</span>
				</div>
			{/if}
		</div>

		<!-- Scrollable terminal content -->
		<div
			bind:this={mobileTerminalEl}
			class="drawer-content"
		>
			{@render terminalLines()}
		</div>
	</div>
{/if}

<style>
	.drawer-handle {
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 6px;
		padding: 12px 16px 4px;
		cursor: grab;
		touch-action: none;
		user-select: none;
		-webkit-user-select: none;
	}

	.handle-pill {
		width: 2.5rem;
		height: 4px;
		border-radius: 9999px;
		background: color-mix(in oklch, var(--color-base-content) 20%, transparent);
	}

	.handle-hint {
		display: flex;
		align-items: center;
		gap: 6px;
		font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
		font-size: 0.7rem;
		opacity: 0.5;
		max-width: 80%;
		overflow: hidden;
	}

	.drawer-content {
		flex: 1;
		overflow-y: auto;
		overscroll-behavior: contain;
		padding: 4px 20px 20px;
		font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
		font-size: 0.75rem;
		line-height: 1.625;
	}
</style>
