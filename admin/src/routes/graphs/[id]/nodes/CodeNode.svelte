<script lang="ts">
	import { Key } from 'lucide-svelte';
	import type { CodeNode } from "$lib/types";
	import BaseNode from "./BaseNode.svelte";
	import type { NodeProps } from "./types";

	let { data }: NodeProps<CodeNode> = $props();
	let node = $derived(data.node);
	let isHovered = $state(false);
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
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
				<span class="key">🔑 Code:</span>
				<span class="value" class:revealed={isHovered}>
                    {#if isHovered}
                        {node.data?.code}
                    {:else}
                        {'•'.repeat(node.data?.code?.length ?? 8)}
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
		letter-spacing: 0.3em;
		color: #9ca3af;
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
