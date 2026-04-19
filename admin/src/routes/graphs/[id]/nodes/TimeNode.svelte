<script lang="ts">
	import { Clock } from 'lucide-svelte';
	import type {TimeNode} from "$lib/types";
	import BaseNode from "./BaseNode.svelte";
	import type {NodeProps} from "./types";

	let { data }: NodeProps<TimeNode> = $props();
	let node = $derived(data.node);

	let formattedTime = $derived(() => {
		if (!node.data?.unlock_at) return 'Not set';
		const date = new Date(node.data.unlock_at);
		return date.toLocaleString();
	});
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
	showProgress={data.showProgress}
	onToggleUnlock={data.onToggleUnlock}
	config={{
        color: '#f59e0b',
        gradient: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
        icon: Clock,
        label: 'Time',
    }}
>
	{#snippet children()}
		<div class="config">
			<div class="config-item">
				<span class="key">Unlocks at:</span>
				<span class="value">{formattedTime()}</span>
			</div>
		</div>
	{/snippet}
</BaseNode>

<style>
	.config {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
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
		font-size: 0.8125rem;
	}
</style>
