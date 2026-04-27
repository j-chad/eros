<script lang="ts">
	import { Heart, Pointer } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { KVStore, KVKey } from '$lib/db/stores/kv';
	import type { FavourCount, FavourRequest } from '$lib/types/favour';

	const { count, requests }: { count: FavourCount; requests: FavourRequest[] } = $props();

	const pendingRequests = $derived(requests.filter((r) => !r.fulfilled_at));

	let panelOpen = $state(false);

	// Tutorial hint — null suppresses render until we know the DB state
	let showCursor: boolean | null = $state(null);

	onMount(async () => {
		const seen = await KVStore.get(KVKey.FavourTipSeen);
		showCursor = !seen;
	});

	function togglePanel() {
		if (showCursor) {
			showCursor = false;
			KVStore.set(KVKey.FavourTipSeen, true);
		}
		panelOpen = !panelOpen;
	}

	function handleNavigate() {
		panelOpen = false;
		goto('/favours');
	}

	function handleBackdropClick() {
		panelOpen = false;
	}
</script>

<!-- Backdrop (dismiss on tap outside) -->
{#if panelOpen}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="fixed inset-0 z-40"
		onclick={handleBackdropClick}
		onkeydown={(e) => { if (e.key === 'Escape') panelOpen = false; }}
	></div>
{/if}

<div class="fixed bottom-6 right-6 z-50 flex flex-col items-end gap-3">
	<!-- Compact panel -->
	{#if panelOpen}
		<div class="card rounded-3xl bg-base-100 shadow-xl shadow-pink-200/40 w-56 animate-popIn">
			<div class="card-body gap-3 p-5">
				<p class="text-sm font-semibold opacity-70">
					{count.remaining} {count.remaining === 1 ? 'favour' : 'favours'} remaining
				</p>
				{#if pendingRequests.length > 0}
					<div class="flex flex-col gap-1">
						<p class="text-xs font-semibold text-primary">Pending</p>
						{#each pendingRequests.slice(0, 3) as req}
							<p class="text-xs opacity-60 truncate">{req.choice_label}</p>
						{/each}
						{#if pendingRequests.length > 3}
							<p class="text-xs opacity-40">+{pendingRequests.length - 3} more</p>
						{/if}
					</div>
				{/if}
				<button class="btn btn-primary btn-sm rounded-2xl" onclick={handleNavigate}>
					View favours
				</button>
			</div>
		</div>
	{/if}

	<!-- Tutorial pointer hint -->
	{#if showCursor === true}
		<div class="absolute -left-8 bottom-10 animate-cursorTap pointer-events-none">
			<Pointer size={28} style="transform: rotate(130deg)" />
		</div>
	{/if}

	<!-- FAB -->
	<button
		class="btn btn-xl btn-primary btn-circle shadow-lg shadow-pink-200/60 text-lg"
		onclick={togglePanel}
		aria-label="Favour menu"
	>
		<Heart size={20} />
		<span class="text-sm font-bold">{count.remaining}</span>
	</button>
</div>
