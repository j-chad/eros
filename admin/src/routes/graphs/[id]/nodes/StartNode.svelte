<!-- StartNode.svelte -->
<script lang="ts">
    import {type NodeProps, Position} from '@xyflow/svelte';
    import { Play, Edit } from 'lucide-svelte';
    import Handle from "./Handle.svelte";

    let { data }: NodeProps = $props();
</script>

<div class="start-node">
    <div class="icon-wrapper">
        <Play size={16} />
    </div>
    <div class="title">
        <div class="name">Start</div>
    </div>

    <Handle type="source" nodeType="start" position={Position.Right} class="node-handle start-handle" />
</div>

<style>
    .start-node {
        background: white;
        border: 2px solid #10b981;
        border-radius: 8px;
        min-width: 200px;
        max-width: 280px;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);

        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem 1rem;
        background: linear-gradient(135deg, #10b981 0%, #059669 100%);
        color: white;
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

    .name {
        font-size: 0.875rem;
        font-weight: 600;
        margin-top: 0.125rem;
    }

    .node-body {
        padding: 1rem;
    }

    .description {
        font-size: 0.875rem;
        color: #6b7280;
        line-height: 1.5;
        word-wrap: break-word;
    }
</style>