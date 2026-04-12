<script lang="ts">
	import { Position } from '@xyflow/svelte';
	import {Edit, Lock, LockOpen, type Play} from 'lucide-svelte';
	import type {AnyNode, NodeType} from "$lib/types";
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
		onmouseleave,
		showProgress = false,
		onToggleUnlock
	}: {
		node: AnyNode
		config: BaseNodeConfig;
		children?: any;
		onEdit?: (nodeID: string) => void;
		onmouseenter?: () => void;
		onmouseleave?: () => void;
		showProgress?: boolean;
		onToggleUnlock?: (nodeID: string) => void;
	} = $props();

	let isUnlocked = $derived(!!node.unlocked_at);
	let isStart = $derived(node.type === 'start' as NodeType);
</script>

<div
	class="base-node"
	class:progress-unlocked={showProgress && isUnlocked}
	class:progress-locked={showProgress && !isUnlocked}
	role="presentation"
	style="border-color: {showProgress ? (isUnlocked ? '#10b981' : '#d1d5db') : config.color}"
	{onmouseenter}
	{onmouseleave}
>
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

	{#if showProgress}
		<div class="progress-badge" class:unlocked={isUnlocked}>
			<div class="badge-info">
				{#if isUnlocked}
					<span class="badge-dot unlocked-dot"></span>
					<span>Unlocked{node.unlocked_at ? ` ${new Date(node.unlocked_at).toLocaleDateString()}` : ''}</span>
				{:else}
					<span class="badge-dot locked-dot"></span>
					<span>Locked</span>
				{/if}
			</div>
			{#if !isStart}
				<button
					class="unlock-toggle"
					class:is-unlocked={isUnlocked}
					onclick={(e) => { e.stopPropagation(); onToggleUnlock?.(node.id); }}
					title={isUnlocked ? 'Lock this node' : 'Unlock this node'}
				>
					{#if isUnlocked}
						<Lock size={12} />
					{:else}
						<LockOpen size={12} />
					{/if}
				</button>
			{/if}
		</div>
	{/if}

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

	.progress-unlocked {
		box-shadow: 0 0 0 2px #10b981, 0 4px 12px rgba(16, 185, 129, 0.3);
	}

	.progress-locked {
		opacity: 0.55;
		filter: grayscale(0.6);
	}

	.progress-badge {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.375rem 0.75rem;
		font-size: 0.6875rem;
		font-weight: 600;
		border-top: 1px solid #e5e7eb;
	}

	.badge-info {
		display: flex;
		align-items: center;
		gap: 0.375rem;
	}

	.unlock-toggle {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.25rem;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		background: white;
		cursor: pointer;
		color: #6b7280;
		transition: all 0.15s ease;
	}

	.unlock-toggle:hover {
		background: #f3f4f6;
		border-color: #9ca3af;
	}

	.unlock-toggle.is-unlocked:hover {
		background: #fef2f2;
		border-color: #f87171;
		color: #dc2626;
	}

	.unlock-toggle:not(.is-unlocked):hover {
		background: #ecfdf5;
		border-color: #10b981;
		color: #059669;
	}

	.progress-badge.unlocked {
		background: #ecfdf5;
		color: #065f46;
	}

	.progress-badge:not(.unlocked) {
		background: #f9fafb;
		color: #6b7280;
	}

	.badge-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		flex-shrink: 0;
	}

	.unlocked-dot {
		background: #10b981;
	}

	.locked-dot {
		background: #9ca3af;
	}
</style>
