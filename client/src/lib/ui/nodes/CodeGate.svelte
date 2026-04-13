<script lang="ts">
	import type { CodeNode } from '$lib/types/graph';
	import type { UnlockResult } from '$lib/api/graph.api';
	import { KeyRound } from 'lucide-svelte';
	import { useOnlineStatus } from '$lib/online.svelte';

	const { node, graphId, onUnlock }: { 
		node: CodeNode; 
		graphId: string; 
		onUnlock: (result: UnlockResult) => void 
	} = $props();

	const isOnline = $derived(useOnlineStatus());
	
	let inputValue = $state('');
	let isSubmitting = $state(false);
	let errorMessage = $state<string | null>(null);
	let showShake = $state(false);

	async function handleSubmit() {
		if (!inputValue.trim() || isSubmitting || !isOnline) return;

		isSubmitting = true;
		errorMessage = null;

		try {
			const { unlockNode } = await import('$lib/services/graph');
			const result = await unlockNode(graphId, node.id, inputValue.trim());
			onUnlock(result);
		} catch (error: any) {
			// Handle specific error types
			if (error.status === 403) {
				errorMessage = "That's not quite right — try again.";
				inputValue = '';
				showShake = true;
				setTimeout(() => showShake = false, 600);
			} else if (error.status === 429) {
				errorMessage = "Too many attempts. Try again shortly.";
			} else {
				errorMessage = "Something went wrong. Please try again.";
			}
		} finally {
			isSubmitting = false;
		}
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			e.preventDefault();
			handleSubmit();
		}
	}
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
		
		{#if errorMessage}
			<div class="alert alert-error rounded-2xl text-sm">
				{errorMessage}
			</div>
		{/if}
		
		<input
			type="text"
			placeholder="••••••"
			bind:value={inputValue}
			onkeydown={handleKeyDown}
			disabled={isSubmitting || !isOnline}
			class="input input-bordered rounded-2xl w-full text-center font-mono tracking-widest text-lg"
			class:animate-shake={showShake}
		/>
		
		<button 
			onclick={handleSubmit}
			disabled={isSubmitting || !isOnline || !inputValue.trim()}
			class="btn btn-primary rounded-2xl w-full"
			class:loading={isSubmitting}
		>
			{#if isSubmitting}
				<span class="loading loading-spinner loading-sm"></span>
				Unlocking...
			{:else}
				Unlock
			{/if}
		</button>
		
		{#if !isOnline}
			<p class="text-xs opacity-50 text-center">You're offline. Connect to try unlocking.</p>
		{/if}
	</div>
</div>

<style>
	@keyframes shake {
		0%, 100% { transform: translateX(0); }
		25% { transform: translateX(-4px); }
		75% { transform: translateX(4px); }
	}
	
	.animate-shake {
		animation: shake 0.6s ease-in-out;
	}
</style>
