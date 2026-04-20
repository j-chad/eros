<script lang="ts">
	import { Key } from 'lucide-svelte';
	import type { CodeNode } from "$lib/types";
	import BaseNode from "./BaseNode.svelte";
	import type { NodeProps } from "./types";

	let { data }: NodeProps<CodeNode> = $props();
	let node = $derived(data.node);
	let isHovered = $state(false);
	let codes = $derived(node.data?.codes ?? []);
	let codeSet = $derived(codes.length > 0);
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
	showProgress={data.showProgress}
	onToggleUnlock={data.onToggleUnlock}
	onmouseenter={() => isHovered = true}
	onmouseleave={() => isHovered = false}
	config={{
        color: '#8b5cf6',
        gradient: 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)',
        icon: Key,
        label: 'Code',
    }}
>
	{#snippet children()}
		<div
			role="presentation"
			class="config"
		>
			<div class="config-item">
				<span class="key">🔑 Code{codes.length > 1 ? 's' : ''}:</span>
				<span class="value" class:revealed={isHovered}>
                    {#if !codeSet}
                        N/A
                    {:else if isHovered}
                        {codes[0]}{#if codes.length > 1}<span class="badge">+{codes.length - 1}</span>{/if}
                    {:else}
                        ••••••••
                    {/if}
                </span>
			</div>
		</div>
	{/snippet}
</BaseNode>

<style>
	.config {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
		cursor: pointer;
	}

	.config-item {
		font-size: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
	}

	.key {
		color: #6b7280;
		font-weight: 500;
	}

	.value {
		color: #1f2937;
		font-weight: 600;
		font-family: 'Courier New', monospace;
		font-size: 0.875rem;
		letter-spacing: 0.15em;
		transition: all 0.2s;
		overflow-x: auto;
		overflow-y: hidden;
		white-space: nowrap;
		max-width: 100%;
		padding-bottom: 0.25rem;
	}

	.value:not(.revealed) {
		color: #9ca3af;
	}

	.badge {
		display: inline-block;
		margin-left: 0.375rem;
		padding: 0 0.35rem;
		background: #ede9fe;
		color: #7c3aed;
		border-radius: 4px;
		font-size: 0.7rem;
		font-family: inherit;
		letter-spacing: 0;
		font-weight: 700;
		vertical-align: middle;
	}

	/* Custom scrollbar styling */
	.value::-webkit-scrollbar {
		height: 4px;
	}

	.value::-webkit-scrollbar-track {
		background: #f3f4f6;
		border-radius: 2px;
	}

	.value::-webkit-scrollbar-thumb {
		background: #d1d5db;
		border-radius: 2px;
	}

	.value::-webkit-scrollbar-thumb:hover {
		background: #9ca3af;
	}
</style>
