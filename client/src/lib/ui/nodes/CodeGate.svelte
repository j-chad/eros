<script lang="ts">
	import type { CodeNode } from '$lib/types/graph';
	import { KeyRound } from 'lucide-svelte';

	const { node }: { node: CodeNode } = $props();
</script>

<div class="flex flex-col items-center gap-6 text-center animate-popIn">
	<div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
		<KeyRound size={28} class="text-primary" />
	</div>

	<div class="flex flex-col gap-2">
		<h1 class="text-2xl font-extrabold">{node.title}</h1>
		{#if node.description}
			<p class="text-sm opacity-70 leading-relaxed">{node.description}</p>
		{/if}
	</div>

	<div class="w-full flex flex-col gap-3">
		<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Enter the code</p>
		<input
			type="text"
			placeholder="••••••"
			disabled
			class="input input-bordered rounded-2xl w-full text-center font-mono tracking-widest text-lg"
		/>
		<button disabled class="btn btn-primary rounded-2xl w-full">
			Unlock
		</button>
	</div>
</div>
