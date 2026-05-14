<script lang="ts">
	import { db, promisifyRequest } from '$lib/db/db';
	import { DB_NAME } from '$lib/db/schema';
	import { KVKey } from '$lib/db/stores/kv';
	import { isOnline } from '$lib/online.svelte';
	import { listGraphs, getGraph } from '$lib/services/graph';
	import { listChoices, getCount, listRequests } from '$lib/services/favour';
	import { clearToken } from '$lib/api/auth';
	import { goto } from '$app/navigation';
	import { ChevronLeft } from 'lucide-svelte';
	import BrandHeader from '$lib/ui/BrandHeader.svelte';

	type LogType = 'command' | 'result' | 'error' | 'summary';
	type LogEntry = { type: LogType; text: string };

	let log = $state<LogEntry[]>([]);
	let clearing = $state(false);
	let syncing = $state(false);
	let loggingOut = $state(false);
	let terminalEl: HTMLDivElement | undefined = $state();

	const online = $derived(isOnline());
	const busy = $derived(clearing || syncing || loggingOut);

	$effect(() => {
		if (log.length && terminalEl) {
			terminalEl.scrollTop = terminalEl.scrollHeight;
		}
	});

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

	async function clearDatabase() {
		clearing = true;
		log = [];

		try {
			append('read auth token');
			const kvStore = await db.getStore('kv', 'readonly');
			const tx = kvStore.transaction; // grab the underlying IDBTransaction
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

	{#if log.length > 0}
		<div
			bind:this={terminalEl}
			class="bg-base-300 rounded-3xl p-5 font-mono text-xs leading-relaxed
				shadow-[0_2px_12px_0_theme(colors.pink.200/40)]
				lg:sticky lg:top-6 lg:self-start lg:max-h-[calc(100dvh-3rem)] lg:overflow-y-auto"
		>
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
					<p class="font-bold mt-3 pt-3 border-t border-base-300">{entry.text}</p>
				{/if}
			{/each}
		</div>
	{/if}
</div>
