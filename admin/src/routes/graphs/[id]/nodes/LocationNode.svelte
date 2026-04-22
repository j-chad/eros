<script lang="ts">
	import { MapPin } from 'lucide-svelte';
	import type {LocationNode} from "$lib/types";
	import BaseNode from "./BaseNode.svelte";
	import type {NodeProps} from "./types";

	let { data }: NodeProps<LocationNode> = $props();
	let node = $derived(data.node);

	let hasCoords = $derived(
		node.data && (node.data.latitude !== 0 || node.data.longitude !== 0)
	);

	function tileUrl(lat: number, lng: number, zoom: number): string {
		const n = Math.pow(2, zoom);
		const x = Math.floor(((lng + 180) / 360) * n);
		const latRad = (lat * Math.PI) / 180;
		const y = Math.floor((1 - Math.log(Math.tan(latRad) + 1 / Math.cos(latRad)) / Math.PI) / 2 * n);
		return `https://tile.openstreetmap.org/${zoom}/${x}/${y}.png`;
	}

	let previewTile = $derived(
		hasCoords ? tileUrl(node.data!.latitude, node.data!.longitude, 12) : null
	);
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
	showProgress={data.showProgress}
	onToggleUnlock={data.onToggleUnlock}
	config={{
        color: '#ec4899',
        gradient: 'linear-gradient(135deg, #ec4899 0%, #db2777 100%)',
        icon: MapPin,
        label: 'Location',
    }}
>
	{#snippet children()}
		{#if previewTile}
			<div class="map-preview">
				<img src={previewTile} alt="Map preview" class="map-tile" draggable="false" />
			</div>
		{/if}
		<div class="config">
			<div class="config-item">
				<span class="key">Coordinates</span>
				<span class="value">
                    {node.data?.latitude.toFixed(6)}, {node.data?.longitude.toFixed(6)}
                </span>
			</div>
			<div class="config-item">
				<span class="key">Radius</span>
				<span class="value">{node.data?.radius_m ?? 0}m</span>
			</div>
			{#if node.data?.show_hint}
				<div class="config-item">
					<span class="key hint-badge">Hint enabled</span>
				</div>
			{/if}
		</div>
	{/snippet}
</BaseNode>

<style>
	.map-preview {
		width: 100%;
		height: 100px;
		overflow: hidden;
		border-radius: 4px;
		margin-bottom: 0.75rem;
		background: #f3f4f6;
	}

	.map-tile {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

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

	.hint-badge {
		color: #7c3aed;
		font-size: 0.6875rem;
	}
</style>
