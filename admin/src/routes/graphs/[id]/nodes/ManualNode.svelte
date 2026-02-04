<script lang="ts">
	import { HandMetal } from 'lucide-svelte';
	import type {ManualNode} from "$lib/types";
	import BaseNode from "./BaseNode.svelte";
	import type {NodeProps} from "./types";

	let { data }: NodeProps<ManualNode> = $props();
	let node = $derived(data.node);
	let isUnlocked = $derived(!!node.data?.unlocked_at);
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
	config={{
        color: '#10b981',
        gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
        icon: HandMetal,
        label: 'Manual'
    }}
>
	{#snippet children()}
		<div class="manual-content">
			<button
				class:unlocked={isUnlocked}
				class="confirm-button"
				onclick={() => {
					const unlockedAt = isUnlocked ? null : new Date().toISOString();
					data.onUpdateData({
						...node.data ?? {instructions: ""},
						unlocked_at: unlockedAt
					});
				}}
			>
					<div class="button-inner">
						<div class="button-ring"></div>
						{#if isUnlocked}
							<span>Unlocked</span>
						{:else}
							<span>Confirm</span>
						{/if}
					</div>
			</button>
		</div>
	{/snippet}
</BaseNode>

<style>
	.manual-content {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
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

	.confirm-button:not(.unlocked) {
		background: linear-gradient(135deg, #dc2626 0%, #991b1b 100%);
		color: white;
	}

	.confirm-button.unlocked {
		background: #d1fae5;
		color: #065f46;
	}

	.confirm-button:not(.unlocked):hover {
		transform: scale(1.1);
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
	}

	.confirm-button:not(.unlocked):active {
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

	.confirm-button:not(.unlocked):hover .button-ring {
		animation: ping 1s cubic-bezier(0, 0, 0.2, 1) infinite;
	}

	@keyframes ping {
		75%, 100% {
			transform: translate(-50%, -50%) scale(1.5);
			opacity: 0;
		}
	}
</style>
