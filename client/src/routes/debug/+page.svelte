<script lang="ts">
	import { db, promisifyRequest } from '$lib/db/db';
	import { DB_NAME, DB_VERSION } from '$lib/db/schema';
	import { KVKey, KVStore } from '$lib/db/stores/kv';
	import {
		isOnline,
		getOfflineReason,
		isOfflineOverridden,
		setOfflineOverride,
		type OfflineReason,
	} from '$lib/online.svelte';
	import { listGraphs, getGraph } from '$lib/services/graph';
	import { listChoices, getCount, listRequests } from '$lib/services/favour';
	import { clearToken } from '$lib/api/auth';
	import { goto } from '$app/navigation';
	import { ChevronLeft, ChevronUp, RefreshCw } from 'lucide-svelte';
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import BrandHeader from '$lib/ui/BrandHeader.svelte';
	import type { FavourCount } from '$lib/types/favour';

	// @ts-expect-error — daisyui theme objects have no type declarations
	import draculaTheme from 'daisyui/theme/dracula/object.js';
	// @ts-expect-error — daisyui theme objects have no type declarations
	import retroTheme from 'daisyui/theme/retro/object.js';
	// @ts-expect-error — daisyui theme objects have no type declarations
	import forestTheme from 'daisyui/theme/forest/object.js';
	import {env} from "$env/dynamic/public";

	type LogType = 'command' | 'result' | 'error' | 'summary';
	type LogEntry = { type: LogType; text: string };

	// --- Log state ---
	let log = $state<LogEntry[]>([]);
	let clearing = $state(false);
	let syncing = $state(false);
	let loggingOut = $state(false);
	let pinging = $state(false);
	let updatingSW = $state(false);
	let unregisteringSW = $state(false);
	let resettingTips = $state(false);
	let desktopTerminalEl: HTMLDivElement | undefined = $state();
	let mobileTerminalEl: HTMLDivElement | undefined = $state();

	const online = $derived(isOnline());
	const offlineReason = $derived(getOfflineReason());
	const offlineOverride = $derived(isOfflineOverridden());
	const busy = $derived(
		clearing || syncing || loggingOut || pinging || updatingSW || unregisteringSW || resettingTips
	);

	const lastEntry = $derived(log.length > 0 ? log[log.length - 1] : null);
	const lastEntryHint = $derived(
		lastEntry
			? lastEntry.type === 'command'
				? `> ${lastEntry.text}`
				: lastEntry.text
			: ''
	);

	// --- Status display ---
	const REASON_LABELS: Record<NonNullable<OfflineReason>, string> = {
		'browser-offline': 'Browser offline',
		'ping-failed': 'Server unreachable',
		'ping-timeout': 'Server timed out',
		'api-error': 'API request failed',
		forced: 'Forced offline',
	};

	const statusLabel = $derived(
		online ? 'Server reachable' : REASON_LABELS[offlineReason!] ?? 'Offline'
	);

	const statusColor = $derived(
		online ? 'bg-success' : offlineReason === 'forced' ? 'bg-info' : 'bg-warning'
	);

	const THEMES: { name: string; label: string; obj: Record<string, string> | null }[] = [
		{ name: 'valentine', label: 'Valentine', obj: null },
		{ name: 'dracula', label: 'Dracula', obj: draculaTheme },
		{ name: 'retro', label: 'Retro', obj: retroTheme },
		{ name: 'forest', label: 'Forest', obj: forestTheme },
	];

	let activeTheme = $state('valentine');

	function applyTheme(themeName: string) {
		activeTheme = themeName;
		const theme = THEMES.find((t) => t.name === themeName);

		// Remove any previously injected theme style
		const existing = document.getElementById('debug-theme');

		if (!theme || !theme.obj) {
			// Revert to default valentine (built into the bundle)
			existing?.remove();
			document.documentElement.removeAttribute('data-theme');
			return;
		}

		const props = Object.entries(theme.obj)
			.map(([k, v]) => (k === 'color-scheme' ? `color-scheme: ${v}` : `${k}: ${v}`))
			.join(';\n  ');

		const css = `:root, [data-theme="${themeName}"] {\n  ${props};\n}`;

		if (existing) {
			existing.textContent = css;
		} else {
			const style = document.createElement('style');
			style.id = 'debug-theme';
			style.textContent = css;
			document.head.appendChild(style);
		}

		document.documentElement.setAttribute('data-theme', themeName);
	}

	interface Diagnostics {
		serverUrl: string;
		gitSha: string;
		maintenanceMode: boolean;
		dbVersion: number;
		authSet: string | null;
		swStatus: string;
		storeCounts: Record<string, number>;
		favourBalance: string | null;
	}

	let diagnostics = $state<Diagnostics | null>(null);
	let loadingDiagnostics = $state(false);

	function timeAgo(timestamp: number): string {
		const seconds = Math.floor((Date.now() - timestamp) / 1000);
		if (seconds < 60) return `${seconds}s ago`;
		const minutes = Math.floor(seconds / 60);
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		if (hours < 24) return `${hours}h ago`;
		const days = Math.floor(hours / 24);
		return `${days}d ago`;
	}

	async function loadDiagnostics() {
		loadingDiagnostics = true;
		try {
			// Auth timestamp
			let authSet: string | null = null;
			try {
				const kvStore = await db.getStore('kv', 'readonly');
				const authEntry = await promisifyRequest(kvStore.get(KVKey.AuthSession));
				if (authEntry) {
					authSet = timeAgo((authEntry as { timestamp: number }).timestamp);
				}
			} catch {
				/* db might not be initialized */
			}

			// Store counts
			const storeNames = [
				'graphs',
				'graphDetails',
				'favourChoices',
				'favourRequests',
				'kv',
			] as const;
			const storeCounts: Record<string, number> = {};
			for (const name of storeNames) {
				try {
					const store = await db.getStore(name, 'readonly');
					storeCounts[name] = await promisifyRequest(store.count());
				} catch {
					storeCounts[name] = -1;
				}
			}

			// Favour balance
			let favourBalance: string | null = null;
			try {
				const count = (await KVStore.get(KVKey.FavourCount)) as FavourCount | null;
				if (count) {
					favourBalance = `${count.remaining} / ${count.total} remaining`;
				}
			} catch {
				/* ignore */
			}

			// Service worker status
			let swStatus = 'Not supported';
			if ('serviceWorker' in navigator) {
				try {
					const reg = await navigator.serviceWorker.getRegistration();
					if (!reg) {
						swStatus = 'Not registered';
					} else if (reg.waiting) {
						swStatus = 'Update waiting';
					} else if (reg.active) {
						swStatus = 'Active';
					} else if (reg.installing) {
						swStatus = 'Installing';
					}
				} catch {
					swStatus = 'Error reading';
				}
			}

			diagnostics = {
				serverUrl: PUBLIC_SERVER_URL,
				gitSha: __GIT_SHA__,
				maintenanceMode: env.PUBLIC_MAINTENANCE_MODE === 'true',
				dbVersion: DB_VERSION,
				authSet,
				swStatus,
				storeCounts,
				favourBalance,
			};
		} finally {
			loadingDiagnostics = false;
		}
	}

	// Load diagnostics on mount
	$effect(() => {
		loadDiagnostics();
	});

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

	$effect(() => {
		function onResize() {
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

		const endY = e.changedTouches[0].clientY;
		const velocity = touchStartY - endY;

		const collapsedH = snapHeightPx('collapsed');
		const halfH = snapHeightPx('half');
		const fullH = snapHeightPx('full');

		const FLICK_THRESHOLD = 60;

		if (Math.abs(velocity) > FLICK_THRESHOLD) {
			if (velocity > 0) {
				snapTo(drawerSnap === 'collapsed' || drawerHeight < halfH ? 'half' : 'full');
			} else {
				snapTo(drawerSnap === 'full' || drawerHeight > halfH ? 'half' : 'collapsed');
			}
			return;
		}

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
			loadDiagnostics();
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
			loadDiagnostics();
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

	async function pingServer() {
		pinging = true;
		log = [];

		try {
			append(`ping ${PUBLIC_SERVER_URL}/health`);
			const start = performance.now();
			const controller = new AbortController();
			const timeout = setTimeout(() => controller.abort(), 5000);

			const res = await fetch(`${PUBLIC_SERVER_URL}/health`, {
				method: 'GET',
				signal: controller.signal,
			});
			clearTimeout(timeout);

			const latency = Math.round(performance.now() - start);

			if (res.ok) {
				appendResult(`${res.status} OK — ${latency}ms`);
			} else {
				appendResult(`${res.status} ${res.statusText} — ${latency}ms`);
			}
			appendSummary('server reachable');
		} catch (e) {
			if (e instanceof DOMException && e.name === 'AbortError') {
				appendError('timed out after 5s');
			} else {
				appendError(e instanceof Error ? e.message : String(e));
			}
			appendSummary('server unreachable');
		} finally {
			pinging = false;
		}
	}

	function toggleOfflineOverride() {
		setOfflineOverride(!offlineOverride);
	}

	async function forceUpdateSW() {
		updatingSW = true;
		log = [];

		try {
			if (!('serviceWorker' in navigator)) {
				appendError('service workers not supported');
				return;
			}

			append('get registration');
			const reg = await navigator.serviceWorker.getRegistration();
			if (!reg) {
				appendError('no service worker registered');
				return;
			}
			appendResult('ok');

			append('check for update');
			await reg.update();

			if (reg.waiting) {
				append('activate waiting worker');
				reg.waiting.postMessage({ type: 'SKIP_WAITING' });
				appendResult('ok');
				appendSummary('updated — reloading...');
				window.location.reload();
			} else if (reg.installing) {
				appendResult('installing...');
				appendSummary('new worker installing — check back shortly');
			} else {
				appendResult('already up to date');
				appendSummary('no update available');
			}
			loadDiagnostics();
		} catch (e) {
			appendError(e instanceof Error ? e.message : String(e));
		} finally {
			updatingSW = false;
		}
	}

	async function unregisterSW() {
		unregisteringSW = true;
		log = [];

		try {
			if (!('serviceWorker' in navigator)) {
				appendError('service workers not supported');
				return;
			}

			append('get registration');
			const reg = await navigator.serviceWorker.getRegistration();
			if (!reg) {
				appendResult('none registered');
				appendSummary('nothing to unregister');
				return;
			}

			append('unregister service worker');
			const ok = await reg.unregister();
			appendResult(ok ? 'ok' : 'failed');

			append('clear cache storage');
			const keys = await caches.keys();
			for (const key of keys) {
				await caches.delete(key);
				appendResult(`deleted "${key}"`);
			}
			if (keys.length === 0) {
				appendResult('no caches found');
			}

			appendSummary('service worker removed');
			loadDiagnostics();
		} catch (e) {
			appendError(e instanceof Error ? e.message : String(e));
		} finally {
			unregisteringSW = false;
		}
	}

	async function resetTips() {
		resettingTips = true;
		log = [];

		try {
			append('delete CalendarTipSeen');
			await KVStore.delete(KVKey.CalendarTipSeen);
			appendResult('ok');

			append('delete FavourTipSeen');
			await KVStore.delete(KVKey.FavourTipSeen);
			appendResult('ok');

			appendSummary('onboarding tips reset');
			loadDiagnostics();
		} catch (e) {
			appendError(e instanceof Error ? e.message : String(e));
		} finally {
			resettingTips = false;
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

<div
	class="mx-auto min-h-dvh max-w-4xl px-4 py-6 pb-20 lg:pb-6 flex flex-col gap-6 lg:grid lg:grid-cols-2 lg:gap-8"
>
	<div class="flex flex-col gap-6">
		<BrandHeader subtitle="Debug tools">
			{#snippet rightContent()}
				<button class="btn btn-ghost btn-sm rounded-2xl gap-1" onclick={() => goto('/')}>
					<ChevronLeft size={16} />
					Back
				</button>
			{/snippet}
		</BrandHeader>

		<!-- Diagnostics card -->
		{#if diagnostics}
			<div
				class="bg-base-100 rounded-3xl p-5 shadow-[0_2px_12px_0_--theme(--color-pink-200/40)] flex flex-col gap-4"
			>
				<div class="flex items-center justify-between">
					<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Info</p>
					<button
						class="btn btn-ghost btn-xs rounded-xl gap-1 opacity-40 hover:opacity-100"
						disabled={loadingDiagnostics}
						onclick={() => loadDiagnostics()}
					>
						<RefreshCw size={12} class={loadingDiagnostics ? 'animate-spin' : ''} />
					</button>
				</div>

				<div class="flex flex-col gap-2.5 text-sm">
					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Server</span>
						<span class="font-mono text-xs truncate ml-4">{diagnostics.serverUrl}</span>
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Build</span>
						<span class="font-mono text-xs">{diagnostics.gitSha}</span>
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Maintenance Mode</span>
						<span class="font-mono text-xs">{diagnostics.maintenanceMode}</span>
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Auth</span>
						<span class="text-xs">
							{#if diagnostics.authSet}
								<span class="font-semibold">Active</span>
								<span class="opacity-40">({diagnostics.authSet})</span>
							{:else}
								<span class="opacity-40">None</span>
							{/if}
						</span>
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Database</span>
						<span class="font-mono text-xs">{DB_NAME} v{diagnostics.dbVersion}</span>
					</div>

					<div class="flex flex-wrap gap-1.5">
						{#each Object.entries(diagnostics.storeCounts).filter(([k]) => k !== 'kv') as [name, count]}
							<span class="bg-base-200 rounded-xl px-2 py-0.5 text-xs font-mono">
								{name} <span class="opacity-50">{count === -1 ? '?' : count}</span>
							</span>
						{/each}
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Favours</span>
						<span class="font-mono text-xs">{diagnostics.favourBalance ?? 'Not cached'}</span>
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Flag</span>
						<span class="font-mono text-xs">eros&lbrace;n0t_f0r_y0ur_3y3s&rbrace;</span>
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Source Maps</span>
						<span class="font-mono text-xs">hidden</span>
					</div>

					<div class="border-t border-base-content/5"></div>

					<div class="flex items-center justify-between">
						<span class="opacity-50 text-xs">Service Worker</span>
						<span class="font-mono text-xs">{diagnostics.swStatus}</span>
					</div>
				</div>
			</div>
		{/if}

		<!-- Status -->
		<div class="flex flex-col gap-3">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Connectivity</p>
			<div class="flex flex-col gap-1.5">
				<div class="flex items-center gap-2">
					<span class="inline-block h-2 w-2 rounded-full {statusColor}"></span>
					<span class="text-sm font-semibold">{statusLabel}</span>
				</div>
				{#if !online && offlineReason && offlineReason !== 'forced'}
					<p class="text-xs opacity-50 pl-4">
						{#if offlineReason === 'browser-offline'}
							Device has no network connection
						{:else if offlineReason === 'ping-failed'}
							Health check failed — server may be down or URL is wrong
						{:else if offlineReason === 'ping-timeout'}
							Health check timed out after 5s
						{:else if offlineReason === 'api-error'}
							An API request threw a network error
						{/if}
					</p>
				{/if}
				<div class="flex items-center gap-3 text-xs opacity-40 pl-4 font-mono">
					<span>navigator.onLine: {typeof navigator !== 'undefined' ? String(navigator.onLine) : '?'}</span>
				</div>
			</div>

			<div class="flex gap-2">
				<button
					class="btn btn-outline rounded-2xl flex-1"
					disabled={busy}
					onclick={pingServer}
				>
					{#if pinging}
						<span class="loading loading-spinner loading-sm"></span>
					{/if}
					Ping server
				</button>

				<button
					class="btn rounded-2xl flex-1"
					class:btn-info={!offlineOverride}
					class:btn-warning={offlineOverride}
					onclick={toggleOfflineOverride}
				>
					{offlineOverride ? 'Resume connectivity' : 'Force offline'}
				</button>
			</div>
		</div>

		<!-- Data -->
		<div class="flex flex-col gap-3">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Data</p>
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

			<button class="btn btn-outline rounded-2xl" disabled={busy} onclick={resetTips}>
				{#if resettingTips}
					<span class="loading loading-spinner loading-sm"></span>
				{/if}
				Reset onboarding tips
			</button>
		</div>

		<!-- Service Worker -->
		<div class="flex flex-col gap-3">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Service Worker</p>
			<button class="btn btn-outline rounded-2xl" disabled={busy} onclick={forceUpdateSW}>
				{#if updatingSW}
					<span class="loading loading-spinner loading-sm"></span>
				{/if}
				Force update
			</button>

			<button
				class="btn btn-error btn-outline rounded-2xl"
				disabled={busy}
				onclick={unregisterSW}
			>
				{#if unregisteringSW}
					<span class="loading loading-spinner loading-sm"></span>
				{/if}
				Unregister + clear caches
			</button>
		</div>

		<!-- Theme -->
		<div class="flex flex-col gap-3">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Theme</p>
			<p class="text-xs opacity-40">Session only — resets on reload</p>
			<div class="flex flex-wrap gap-2">
				{#each THEMES as theme}
					<button
						class="btn btn-outline rounded-2xl flex-1 min-w-[calc(50%-0.25rem)]"
						class:btn-active={activeTheme === theme.name}
						onclick={() => applyTheme(theme.name)}
					>
						{theme.label}
					</button>
				{/each}
			</div>
		</div>

		<!-- Session -->
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
				shadow-[0_2px_12px_0_theme(colors.pink.200/40)]
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
			onkeydown={(e) => {
				if (e.key === 'Enter' || e.key === ' ') handleHandleTap();
			}}
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
		<div bind:this={mobileTerminalEl} class="drawer-content">
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
		padding: 4px 20px calc(20px + env(safe-area-inset-bottom, 0px));
		font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
		font-size: 0.75rem;
		line-height: 1.625;
	}
</style>
