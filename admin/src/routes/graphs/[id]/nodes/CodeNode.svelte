<!-- CodeNode.svelte -->
<script lang="ts">
    import { Handle, Position } from '@xyflow/svelte';
    import { Key, Edit } from 'lucide-svelte';

    let { data } = $props();
</script>

<div class="code-node">
    <div class="node-header">
        <div class="icon-wrapper">
            <Key size={16} />
        </div>
        <div class="title">
            <div class="type">Code</div>
            <div class="name">{data.title}</div>
        </div>
        <button class="edit-btn" onclick={data.onEdit}>
            <Edit size={14} />
        </button>
    </div>

    <div class="node-body">
        {#if data.description}
            <div class="description">{data.description}</div>
        {/if}
        <div class="code-display">
            <span class="code-label">Code:</span>
            <span class="code-value">{data.data.code}</span>
        </div>
    </div>

    <Handle type="target" position={Position.Top} />
    <Handle type="source" position={Position.Bottom} />
</div>

<style>
    .code-node {
        background: white;
        border: 2px solid #8b5cf6;
        border-radius: 8px;
        min-width: 200px;
        max-width: 280px;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    }

    .node-header {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem 1rem;
        background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
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

    .code-display {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
        padding: 0.625rem 0.75rem;
        background: #f3f4f6;
        border-radius: 4px;
        border: 1px solid #e5e7eb;
    }

    .code-label {
        font-size: 0.625rem;
        color: #6b7280;
        font-weight: 500;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .code-value {
        font-size: 1.125rem;
        color: #1f2937;
        font-weight: 700;
        font-family: 'Courier New', monospace;
        letter-spacing: 0.15em;
    }
</style>