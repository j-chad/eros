<script lang="ts">
	import type { LocationNode } from '$lib/types/graph';
	import { MapPin } from 'lucide-svelte';

	const { node }: { node: LocationNode } = $props();
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

	<div class="w-full flex flex-col gap-3">
		<div class="bg-base-200 rounded-2xl px-5 py-4 flex flex-col gap-1">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Head to the location</p>
			<p class="text-sm opacity-70">
				You need to be within {node.data?.radius_m ?? 10}&nbsp;metres of the destination.
			</p>
		</div>
		<button disabled class="btn btn-primary rounded-2xl w-full">
			Check location
		</button>
	</div>
</div>
