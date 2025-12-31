<!-- src/lib/components/flow-nodes/RevealNode.svelte -->
<script lang="ts">
    import { Handle, Position } from '@xyflow/svelte';
    import { Eye, Hash, Edit } from 'lucide-svelte';

    let { data } = $props();
</script>

<div class="reveal-node">
    <div class="node-header">
        <div class="icon-wrapper">
            <Eye size={16} />
        </div>
        <div class="title">
            <div class="type">Reveal #{data.order}</div>
            <div class="name">Step {data.order}</div>
        </div>
        <button class="edit-btn" onclick={data.onEdit}>
            <Edit size={14} />
        </button>
    </div>

    <div class="node-body">
        <div class="content">{data.content}</div>
        <div class="gate-indicator">
            <Hash size={12} />
            <span>Gates below must unlock first</span>
        </div>
    </div>

    <!-- Handle for incoming connection from previous reveal or reward -->
    <Handle type="target" position={Position.Top} />

    <!-- Handle for outgoing connection to next reveal -->
    <Handle
            type="source"
            position={Position.Bottom}
            id="next"
            style="left: 50%"
    />

    <!-- Handle for outgoing connections to gates -->
    <Handle
            type="source"
            position={Position.Bottom}
            id="gates"
            style="left: 25%; background: #f59e0b"
    />
</div>

<style>
    .reveal-node {
        background: white;
        border: 2px solid #3b82f6;
        border-radius: 8px;
        min-width: 260px;
        max-width: 320px;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    }

    .node-header {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem 1rem;
        background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
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

    .content {
        font-size: 0.875rem;
        color: #1f2937;
        line-height: 1.5;
        margin-bottom: 0.75rem;
        word-wrap: break-word;
    }

    .gate-indicator {
        display: flex;
        align-items: center;
        gap: 0.375rem;
        padding: 0.5rem;
        background: #fef3c7;
        border: 1px solid #fbbf24;
        border-radius: 4px;
        font-size: 0.75rem;
        color: #92400e;
        font-weight: 500;
    }
</style>