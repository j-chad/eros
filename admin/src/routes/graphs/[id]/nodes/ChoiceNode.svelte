<!-- ManualNode.svelte -->
<script lang="ts">
	import {type Node, type NodeProps} from '@xyflow/svelte';
	import { HandMetal, CheckCircle2, Circle } from 'lucide-svelte';
	import type {ManualNode} from "$lib/types";
	import BaseNode from "./BaseNode.svelte";

	let { data }: NodeProps<Node<{node: ManualNode, onEdit?: (nodeId: string) => void}>> = $props();

	const node = $derived(data.node);
	const onEdit = $derived(data.onEdit);

	// This would come from execution state in a real implementation
	const isCompleted = $derived(node.data?.completed ?? false);
	const isPending = $derived(node.data?.pending ?? false);
</script>

<BaseNode
	{node}
	{onEdit}
	config={{
        color: isCompleted ? '#10b981' : (isPending ? '#f59e0b' : '#8b5cf6'),
        gradient: isCompleted
			? 'linear-gradient(135deg, #10b981 0%, #059669 100%)'
			: (isPending
				? 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)'
				: 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)'),
        icon: HandMetal,
        label: 'Manual'
    }}
>
	{#snippet children()}
		<div class="manual-content">
			<div class="status-indicator">
				{#if isCompleted}
					<div class="status completed">
						<CheckCircle2 size={20} />
						<span>Completed</span>
					</div>
				{:else if isPending}
					<div class="status pending">
						<Circle size={20} class="pulse" />
						<span>Awaiting Confirmation</span>
					</div>
				{:else}
					<div class="status idle">
						<Circle size={20} />
						<span>Not Started</span>
					</div>
				{/if}
			</div>

			<button
				class="confirm-button"
				class:completed={isCompleted}
				class:pending={isPending}
				disabled={isCompleted}
				onclick={() => {
					// This would trigger the confirmation in a real implementation
					console.log('Manual node confirmed:', node.id);
				}}
			>
				{#if isCompleted}
					<CheckCircle2 size={24} />
					<span>Confirmed</span>
				{:else}
					<div class="button-inner">
						<div class="button-ring"></div>
						<span>Confirm</span>
					</div>
				{/if}
			</button>

			<div class="manual-hint">
				✋ Manual confirmation required
			</div>
		</div>
	{/snippet}
</BaseNode>

<style>
	.manual-content {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.status-indicator {
		display: flex;
		justify-content: center;
	}

	.status {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		border-radius: 6px;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.status.completed {
		background: #d1fae5;
		color: #065f46;
		border: 1px solid #6ee7b7;
	}

	.status.pending {
		background: #fef3c7;
		color: #92400e;
		border: 1px solid #fbbf24;
	}

	.status.idle {
		background: #f3f4f6;
		color: #6b7280;
		border: 1px solid #d1d5db;
	}

	.confirm-button {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 1rem;
		border: none;
		border-radius: 50%;
		width: 80px;
		height: 80px;
		margin: 0 auto;
		font-size: 0.875rem;
		font-weight: 700;
		cursor: pointer;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
	}

	.confirm-button:not(.completed):not(.pending) {
		background: linear-gradient(135deg, #dc2626 0%, #991b1b 100%);
		color: white;
	}

	.confirm-button.pending {
		background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
		color: white;
		animation: pulse-button 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
	}

	.confirm-button.completed {
		background: #d1fae5;
		color: #065f46;
		cursor: not-allowed;
	}

	.confirm-button:not(.completed):not(:disabled):hover {
		transform: scale(1.1);
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
	}

	.confirm-button:not(.completed):not(:disabled):active {
		transform: scale(0.95);
	}

	.button-inner {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.25rem;
	}

	.button-ring {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: 100px;
		height: 100px;
		border: 3px solid rgba(255, 255, 255, 0.3);
		border-radius: 50%;
		pointer-events: none;
	}

	.confirm-button:not(.completed):hover .button-ring {
		animation: ping 1s cubic-bezier(0, 0, 0.2, 1) infinite;
	}

	.manual-hint {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.375rem;
		padding: 0.5rem;
		background: #f5f3ff;
		border: 1px solid #c4b5fd;
		border-radius: 4px;
		font-size: 0.75rem;
		color: #5b21b6;
		font-weight: 500;
		text-align: center;
	}

	@keyframes pulse-button {
		0%, 100% {
			opacity: 1;
		}
		50% {
			opacity: 0.8;
		}
	}

	@keyframes ping {
		75%, 100% {
			transform: translate(-50%, -50%) scale(1.5);
			opacity: 0;
		}
	}

	:global(.pulse) {
		animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
	}

	@keyframes pulse {
		0%, 100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
	}
</style>
