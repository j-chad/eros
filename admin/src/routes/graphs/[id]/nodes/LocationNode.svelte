<script lang="ts">
	import { MapPin } from 'lucide-svelte';
	import type {LocationNode} from "$lib/types";
	import BaseNode from "./BaseNode.svelte";
	import type {NodeProps} from "./types";

	let { data }: NodeProps<LocationNode> = $props();
	let node = $derived(data.node);
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
	config={{
        color: '#ec4899',
        gradient: 'linear-gradient(135deg, #ec4899 0%, #db2777 100%)',
        icon: MapPin,
        label: 'Location',
    }}
>
	{#snippet children()}
		<div class="config">
			<div class="config-item">
				<span class="key">📍 Coordinates:</span>
				<span class="value">
                    {node.data?.latitude.toFixed(6)}, {node.data?.longitude.toFixed(6)}
                </span>
			</div>
			<div class="config-item">
				<span class="key">📏 Radius:</span>
				<span class="value">{node.data?.radius_m ?? 0}m</span>
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
		font-family: monospace;
		font-size: 0.6875rem;
	}
</style>
