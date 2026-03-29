<script lang="ts">
	import type { StartNode } from '$lib/types/graph';
	import { MapPin } from 'lucide-svelte';

	const { node }: { node: StartNode } = $props();
</script>

<div class="flex flex-col items-center gap-6 text-center animate-popIn">
	<div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
		<MapPin size={28} class="text-primary" />
	</div>

	<div class="flex flex-col gap-2">
		<h1 class="text-2xl font-extrabold">{node.title}</h1>
		{#if node.description}
			<p class="text-sm opacity-70 leading-relaxed">{node.description}</p>
		{/if}
	</div>

	<div class="badge badge-primary badge-outline rounded-2xl px-4 py-3 text-xs font-semibold">
		Starting point
	</div>
</div>
