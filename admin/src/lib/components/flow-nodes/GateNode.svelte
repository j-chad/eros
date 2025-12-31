<script lang="ts">
    import { Handle, Position } from '@xyflow/svelte';
    import { Lock, Edit } from 'lucide-svelte';

    let { data } = $props();

    function getGateIcon(type: string) {
        const icons: Record<string, string> = {
            'time_delay': '⏰',
            'task_completion': '✓',
            'location': '📍',
            'favour': '🤝',
            'custom': '⚙️'
        };
        return icons[type] || '🔒';
    }

    function getGateColor(type: string) {
        const colors: Record<string, string> = {
            'time_delay': '#f59e0b',
            'task_completion': '#8b5cf6',
            'location': '#ec4899',
            'favour': '#06b6d4',
            'custom': '#6b7280'
        };
        return colors[type] || '#6b7280';
    }
</script>

<div class="gate-node" style="border-color: {getGateColor(data.type)}">
    <div class="node-header" style="background: {getGateColor(data.type)}">
        <div class="icon-wrapper">
            <Lock size={14} />
        </div>
        <div class="title">
            <div class="type">Gate</div>
            <div class="name">
                {getGateIcon(data.type)}
                {data.type.replace('_', ' ')}
            </div>
        </div>
        <button class="edit-btn" onclick={data.onEdit}>
            <Edit size={12} />
        </button>
    </div>

    <div class="node-body">
        <div class="order-badge">Order: {data.unlockOrder}</div>
        {#if data.config}
            <div class="config">
                {#each Object.entries(data.config) as [key, value]}
                    <div class="config-item">
                        <span class="key">{key}:</span>
                        <span class="value">{value}</span>
                    </div>
                {/each}
            </div>
        {/if}
    </div>

    <Handle type="target" position={Position.Top} />
</div>

<style>
    .gate-node {
        background: white;
        border: 2px solid;
        border-radius: 8px;
        min-width: 180px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .node-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.625rem 0.75rem;
        color: white;
        border-radius: 6px 6px 0 0;
    }

    .icon-wrapper {
        width: 24px;
        height: 24px;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: rgba(255, 255, 255, 0.2);
    }

    .title {
        flex: 1;
    }

    .type {
        font-size: 0.625rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        opacity: 0.9;
        font-weight: 600;
    }

    .name {
        font-size: 0.75rem;
        font-weight: 600;
        margin-top: 0.125rem;
        text-transform: capitalize;
    }

    .edit-btn {
        background: rgba(255, 255, 255, 0.2);
        border: none;
        border-radius: 3px;
        padding: 0.25rem;
        color: white;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background 0.2s;
    }

    .edit-btn:hover {
        background: rgba(255, 255, 255, 0.3);
    }

    .node-body {
        padding: 0.75rem;
    }

    .order-badge {
        display: inline-block;
        padding: 0.25rem 0.5rem;
        background: #f3f4f6;
        border-radius: 4px;
        font-size: 0.625rem;
        font-weight: 600;
        color: #6b7280;
        margin-bottom: 0.5rem;
    }

    .config {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .config-item {
        font-size: 0.625rem;
        display: flex;
        gap: 0.25rem;
    }

    .key {
        color: #6b7280;
        font-weight: 500;
        text-transform: capitalize;
    }

    .value {
        color: #1f2937;
        font-weight: 600;
    }
</style>