<script lang="ts">
    import { Handle, Position } from '@xyflow/svelte';
    import { Gift, Calendar, Edit } from 'lucide-svelte';

    let { data } = $props();

    function formatDate(dateString: string) {
        return new Date(dateString).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }
</script>

<div class="reward-node">
    <div class="node-header">
        <div class="icon-wrapper reward">
            <Gift size={16} />
        </div>
        <div class="title">
            <div class="type">Reward</div>
            <div class="name">{data.title}</div>
        </div>
        <button class="edit-btn" onclick={data.onEdit}>
            <Edit size={14} />
        </button>
    </div>

    <div class="node-body">
        <div class="info-row">
            <Calendar size={14} />
            <span class="label">Available:</span>
            <span class="value">{formatDate(data.notBefore)}</span>
        </div>
        {#if data.description}
            <div class="description">{data.description}</div>
        {/if}
    </div>

    <Handle type="source" position={Position.Bottom} />
</div>

<style>
    .reward-node {
        background: white;
        border: 2px solid #10b981;
        border-radius: 8px;
        min-width: 280px;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    }

    .node-header {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem 1rem;
        background: linear-gradient(135deg, #10b981 0%, #059669 100%);
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
    }

    .edit-btn:hover {
        background: rgba(255, 255, 255, 0.3);
    }

    .node-body {
        padding: 1rem;
    }

    .info-row {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.75rem;
        color: #6b7280;
    }

    .label {
        font-weight: 500;
    }

    .value {
        color: #1f2937;
        font-weight: 600;
    }

    .description {
        margin-top: 0.75rem;
        font-size: 0.75rem;
        color: #6b7280;
        line-height: 1.5;
    }
</style>