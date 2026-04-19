<script lang="ts">
	import type { ManualNode } from '$lib/types/graph';
	import type { UnlockResult } from '$lib/api/graph.api';
	import { Clock } from 'lucide-svelte';
	import { useOnlineStatus } from '$lib/online.svelte';
	import { onMount, onDestroy } from 'svelte';

	const { node, graphId, onUnlock }: {
		node: ManualNode;
		graphId: string;
		onUnlock: (result: UnlockResult) => void
	} = $props();

	const isOnline = $derived(useOnlineStatus());
	const isApproved = $derived(!!node.unlocked_at || !!node.data?.unlocked_at);

	let isSubmitting = $state(false);
	let errorMessage = $state<string | null>(null);
	let pollInterval: number | null = null;

	async function handleContinue() {
		if (isSubmitting || !isOnline || !isApproved) return;

		isSubmitting = true;
		errorMessage = null;

		try {
			const { unlockNode } = await import('$lib/services/graph');
			const result = await unlockNode(graphId, node.id, '');
			onUnlock(result);
		} catch (error: any) {
			if (error.status === 403) {
				errorMessage = "Approval was revoked. Please wait for re-approval.";
				// Restart polling since approval status changed
				startPolling();
			} else if (error.status === 429) {
				errorMessage = "Too many attempts. Try again shortly.";
			} else {
				errorMessage = "Something went wrong. Please try again.";
			}
		} finally {
			isSubmitting = false;
		}
	}

	async function handleRefresh() {
		if (!isOnline) return;

		try {
			const { getGraph } = await import('$lib/services/graph');
			await getGraph(graphId); // Re-fetch graph which updates the cached data
			window.location.reload(); // Simple approach - reload the page to get fresh data
		} catch (error) {
			errorMessage = "Failed to refresh. Try again.";
		}
	}

	function startPolling() {
		if (pollInterval || isApproved || !isOnline) return;

		pollInterval = window.setInterval(() => {
			if (document.visibilityState === 'visible' && isOnline) {
				handleRefresh();
			}
		}, 30000); // 30 seconds
	}

	function stopPolling() {
		if (pollInterval) {
			window.clearInterval(pollInterval);
			pollInterval = null;
		}
	}

	function handleVisibilityChange() {
		if (document.visibilityState === 'visible' && !isApproved) {
			startPolling();
		} else {
			stopPolling();
		}
	}

	onMount(() => {
		if (!isApproved) {
			startPolling();
			document.addEventListener('visibilitychange', handleVisibilityChange);
		}
	});

	onDestroy(() => {
		stopPolling();
		document.removeEventListener('visibilitychange', handleVisibilityChange);
	});

	// Stop polling when approved
	$effect(() => {
		if (isApproved) {
			stopPolling();
		} else if (isOnline && !pollInterval) {
			startPolling();
		}
	});

	// Handle online/offline changes
	$effect(() => {
		if (!isOnline) {
			stopPolling();
		} else if (!isApproved && !pollInterval) {
			startPolling();
		}
	});
</script>

<div class="flex flex-col items-center gap-6 text-center animate-popIn">
	<div class="w-16 h-16 rounded-full flex items-center justify-center {isApproved ? 'bg-success/10' : 'bg-primary/10'}">
		<Clock size={28} class={isApproved ? 'text-success' : 'text-primary'} />
	</div>

	<div class="flex flex-col gap-2">
		<h1 class="text-2xl font-extrabold">{node.title}</h1>
		{#if node.description}
			<p class="text-sm opacity-70 leading-relaxed">{node.description}</p>
		{/if}
	</div>

	{#if node.data?.instructions}
		<div class="w-full bg-base-200 rounded-2xl px-5 py-4 text-left">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide mb-2">Instructions</p>
			<p class="text-sm opacity-80 leading-relaxed">{node.data.instructions}</p>
		</div>
	{/if}

	{#if errorMessage}
		<div class="alert alert-error rounded-2xl text-sm">
			{errorMessage}
		</div>
	{/if}

	{#if isApproved}
		<div class="w-full flex flex-col gap-3">
			<div class="alert rounded-2xl text-sm">
				The content is now unlocked. Click continue to proceed.
			</div>
			<button
				onclick={handleContinue}
				disabled={isSubmitting || !isOnline}
				class="btn btn-primary rounded-2xl w-full"
				class:loading={isSubmitting}
			>
				{#if isSubmitting}
					<span class="loading loading-spinner loading-sm"></span>
					Continuing...
				{:else}
					Continue
				{/if}
			</button>
		</div>
	{:else}
		<div class="w-full flex flex-col gap-3">
			<div class="bg-base-200 rounded-2xl px-5 py-4">
				<p class="text-sm opacity-70">
					{#if !isOnline}
						You're offline. Approval checks will resume when you reconnect.
					{:else}
						Waiting for approval. Check back soon.
					{/if}
				</p>
			</div>
			<button
				onclick={handleRefresh}
				disabled={!isOnline}
				class="btn btn-ghost rounded-2xl w-full btn-sm"
			>
				Refresh
			</button>
		</div>
	{/if}
</div>
