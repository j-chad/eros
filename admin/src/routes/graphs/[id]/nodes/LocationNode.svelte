<!-- LocationNode.svelte -->
<script lang="ts">
	import {type NodeProps, Position} from '@xyflow/svelte';
    import { MapPin, Edit } from 'lucide-svelte';
	import type {LocationNode} from "$lib/types";
	import Handle from "./Handle.svelte";

    let { data }: NodeProps = $props();

	const node = $derived(data.node) as LocationNode
</script>

<div class="location-node">
    <div class="node-header">
        <div class="icon-wrapper">
            <MapPin size={16} />
        </div>
        <div class="title">
            <div class="type">Location</div>
            <div class="name">{node.title}</div>
        </div>
        <button class="edit-btn">
            <Edit size={14} />
        </button>
    </div>

    <div class="node-body">
        {#if node.description}
            <div class="description">{node.description}</div>
        {/if}
        <div class="config">
            <div class="config-item">
                <span class="key">📍 Coordinates:</span>
                <span class="value">{node.data?.latitude.toFixed(6)}, {node.data?.longitude.toFixed(6)}</span>
            </div>
            <div class="config-item">
                <span class="key">📏 Radius:</span>
                <span class="value">{node.data?.radius_meters ?? 0}m</span>
            </div>
        </div>
    </div>

    <Handle type="target" position={Position.Left} nodeType="location"/>
    <Handle type="source" position={Position.Right} nodeType="location"/>
</div>

<style>
    .location-node {
        background: white;
        border: 2px solid #ec4899;
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
        background: linear-gradient(135deg, #ec4899 0%, #db2777 100%);
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
