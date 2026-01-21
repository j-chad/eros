<script lang="ts">
	import { Position } from '@xyflow/svelte';
	import {Edit, type Play} from 'lucide-svelte';
	import type { AnyNode } from "$lib/types";
	import Handle from "./Handle.svelte";

	interface BaseNodeConfig {
		color: string;
		gradient: string;
		icon: typeof Play;
		label: string;
		hasTargetHandle?: boolean;
		hasSourceHandle?: boolean;
		isSourceConnectable?: boolean;
		isTargetConnectable?: boolean;
	}

	let {
		node,
		config,
		children,
		onEdit,
		onmouseenter,
		onmouseleave
	}: {
		node: AnyNode
		config: BaseNodeConfig;
		children?: any;
		onEdit?: (nodeID: string) => void;
		onmouseenter?: () => void;
		onmouseleave?: () => void;
	} = $props();
</script>

<div class="base-node" role="presentation" style="border-color: {config.color}" {onmouseenter} {onmouseleave}>
	<div class="node-header" style="background: {config.gradient}">
		<div class="icon-wrapper">
			<config.icon size={16} />
		</div>
		<div class="title">
			<div class="type">{config.label}</div>
			<div class="name">{node.title}</div>
		</div>
		{#if onEdit}
			<button class="edit-btn" onclick={() => onEdit?.(node.id)}>
				<Edit size={14} />
			</button>
		{/if}
	</div>

	{#if node.description || children}
		<div class="node-body">
			{#if node.description}
				<div class="description">{node.description}</div>
			{/if}
			{@render children?.()}
		</div>
	{/if}

	{#if config.hasTargetHandle !== false}
		<Handle type="target" position={Position.Left} nodeType={node.type} isConnectable={config.isTargetConnectable ?? true} />
	{/if}
	{#if config.hasSourceHandle !== false}
		<Handle type="source" position={Position.Right} nodeType={node.type} isConnectable={config.isSourceConnectable ?? true} />
	{/if}
</div>

<style>
	.base-node {
		background: white;
		border: 2px solid;
		border-radius: 8px;
		min-width: 220px;
		max-width: 300px;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
	}

	.node-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		color: white;
		border-radius: 6px 6px 0 0;
	}

	.icon-wrapper {
		width: 32px;
		height: 32px;
		border-radius: 6px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.2);
	}

	.title {
		flex: 1;
		min-width: 0;
	}

	.type {
		font-size: 0.625rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		opacity: 0.9;
		font-weight: 600;
	}

	.name {
		font-size: 0.875rem;
		font-weight: 600;
		margin-top: 0.125rem;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.edit-btn {
		background: rgba(255, 255, 255, 0.2);
		border: none;
		border-radius: 4px;
		padding: 0.375rem;
		color: white;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: background 0.2s;
		flex-shrink: 0;
	}

	.edit-btn:hover {
		background: rgba(255, 255, 255, 0.3);
	}

	.node-body {
		padding: 1rem;
	}

	.description {
		font-size: 0.875rem;
		color: #1f2937;
		line-height: 1.5;
		margin-bottom: 0.75rem;
		word-wrap: break-word;
	}
</style>
