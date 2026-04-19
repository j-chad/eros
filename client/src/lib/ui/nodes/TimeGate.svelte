<script lang="ts">
	import type { TimeNode } from '$lib/types/graph';
	import type { UnlockResult } from '$lib/api/graph.api';
	import { Clock } from 'lucide-svelte';
	import { useOnlineStatus } from '$lib/online.svelte';
	import Countdown from '$lib/ui/Countdown.svelte';

	const { node, graphId, onUnlock }: {
		node: TimeNode;
		graphId: string;
		onUnlock: (result: UnlockResult) => void
	} = $props();

	const isOnline = $derived(useOnlineStatus());
	const unlockAt = $derived(new Date(node.data!.unlock_at));

	let isSubmitting = $state(false);
	let errorMessage = $state<string | null>(null);
	let countdownDone = $state(Date.now() >= unlockAt.getTime());

	async function handleContinue() {
		if (isSubmitting || !isOnline) return;

		isSubmitting = true;
		errorMessage = null;

		try {
			const { unlockNode } = await import('$lib/services/graph');
			const result = await unlockNode(graphId, node.id, '');
			onUnlock(result);
		} catch (error: any) {
			if (error.status === 403) {
				errorMessage = "Not quite yet. Please wait a moment.";
			} else if (error.status === 429) {
				errorMessage = "Too many attempts. Try again shortly.";
			} else {
				errorMessage = "Something went wrong. Please try again.";
			}
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="flex flex-col items-center gap-6 text-center animate-popIn">
	<div class="w-16 h-16 rounded-full flex items-center justify-center {countdownDone ? 'bg-success/10' : 'bg-warning/10'}">
		<Clock size={28} class={countdownDone ? 'text-success' : 'text-warning'} />
	</div>

	<div class="flex flex-col gap-2">
		<h1 class="text-2xl font-extrabold">{node.title}</h1>
		{#if node.description}
			<p class="text-sm opacity-70 leading-relaxed">{node.description}</p>
		{/if}
	</div>

	{#if errorMessage}
		<div class="alert alert-error rounded-2xl text-sm">
			{errorMessage}
		</div>
	{/if}

	{#if countdownDone}
		<!-- Time has passed, show Continue -->
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
		<!-- Countdown -->
		<div class="w-full">
			<Countdown target={unlockAt} label="Unlocks in" onComplete={() => countdownDone = true} />
		</div>
	{/if}
</div>
