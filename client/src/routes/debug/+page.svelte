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

	let log = $state<string[]>([]);
	let clearing = $state(false);
	let syncing = $state(false);
	let loggingOut = $state(false);

	const online = $derived(isOnline());
	const busy = $derived(clearing || syncing || loggingOut);

	function append(msg: string) {
		log = [...log, `$ ${msg}`];
	}

	function appendResult(msg: string) {
		log = [...log, msg];
	}

	async function clearDatabase() {
		clearing = true;
		log = [];

		try {
			// Save the auth token before wiping
			append('read auth token');
			const kvStore = await db.getStore('kv', 'readonly');
			const authEntry = await promisifyRequest(kvStore.get(KVKey.AuthSession));

			append('close db');
			await db.close();

			append(`delete "${DB_NAME}"`);
			await new Promise<void>((resolve, reject) => {
				const req = indexedDB.deleteDatabase(DB_NAME);
				req.onsuccess = () => resolve();
				req.onerror = () => reject(req.error);
				req.onblocked = () => appendResult('  blocked — close other tabs and retry');
			});
			appendResult('  ok');

			append('init db');
			await db.init();
			appendResult('  ok');

			// Restore the auth token
			if (authEntry) {
				append('restore auth token');
				const newKvStore = await db.getStore('kv', 'readwrite');
				await promisifyRequest(newKvStore.put(authEntry));
				appendResult('  ok');
			}

			appendResult('');
			appendResult('database cleared (auth preserved)');
		} catch (e) {
			appendResult(`err: ${e instanceof Error ? e.message : String(e)}`);
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
			appendResult(`  ${graphs.length} graph(s)`);

			for (const g of graphs) {
				append(`fetch graph "${g.title}"`);
				await getGraph(g.id);
				appendResult('  ok');
			}

			append('fetch favour choices');
			const choices = await listChoices();
			appendResult(`  ${choices.length} choice(s)`);

			append('fetch favour count');
			await getCount();
			appendResult('  ok');

			append('fetch favour requests');
			const requests = await listRequests();
			appendResult(`  ${requests.length} request(s)`);

			appendResult('');
			appendResult('sync complete');
		} catch (e) {
			appendResult(`err: ${e instanceof Error ? e.message : String(e)}`);
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
			appendResult('  ok');
			appendResult('');
			appendResult('logged out — redirecting...');
			await goto('/login');
		} catch (e) {
			appendResult(`err: ${e instanceof Error ? e.message : String(e)}`);
			loggingOut = false;
		}
	}
</script>

<svelte:head>
	<title>Debug</title>
</svelte:head>

<div class="mx-auto min-h-dvh max-w-md px-4 py-6 flex flex-col gap-6">
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

	{#if log.length > 0}
		<div class="terminal rounded-2xl p-4 overflow-auto max-h-80 shadow-inner">
			<pre class="text-xs leading-relaxed whitespace-pre-wrap">{log.join('\n')}</pre>
		</div>
	{/if}
</div>

<style>
	.terminal {
		background: #1a1a2e;
		color: #a0e4a0;
		font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
	}
</style>
