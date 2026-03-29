<script lang="ts">
	import type { ManualNode } from '$lib/types/graph';
	import { Clock } from 'lucide-svelte';

	const { node }: { node: ManualNode } = $props();

	const isUnlocked = $derived(!!node.unlocked_at || !!node.data?.unlocked_at);
</script>

<div class="flex flex-col items-center gap-6 text-center animate-popIn">
	<div class="w-16 h-16 rounded-full flex items-center justify-center {isUnlocked ? 'bg-success/10' : 'bg-primary/10'}">
		<Clock size={28} class={isUnlocked ? 'text-success' : 'text-primary'} />
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

	{#if isUnlocked}
		<div class="badge badge-success badge-outline rounded-2xl px-4 py-3 text-xs font-semibold">
			Approved
		</div>
	{:else}
		<div class="bg-base-200 rounded-2xl px-5 py-4 w-full">
			<p class="text-sm opacity-70">Waiting for approval. Check back soon.</p>
		</div>
	{/if}
</div>
