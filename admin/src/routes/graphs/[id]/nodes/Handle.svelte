<script lang="ts">
    import { Handle as FlowHandle, Position, type HandleProps } from '@xyflow/svelte';

    interface Props extends HandleProps {
        nodeType: 'start' | 'location' | 'code' | 'choice' | 'reward';
    }

    let props: Props = $props();

    const colorMap = {
        start: '#10b981',
        location: '#ec4899',
        code: '#8b5cf6',
        choice: '#f59e0b',
        reward: '#06b6d4'
    };

    const handleColor = $derived(colorMap[props.nodeType] ?? '#6b7280');
</script>

<FlowHandle {...props} class="node-handle" style="background: {handleColor}"
/>

<style>
    :global(.node-handle) {
        width: 10px;
        height: 10px;
        border: 3px solid white;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.15);
        transition: box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        z-index: 99999;
    }

    :global(.node-handle:hover) {
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
    }

    :global(.node-handle.connecting) {
        transform: scale(1.4);
        box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.2);
    }

    :global(.node-handle.source) {
        cursor: pointer;
    }

    :global(.node-handle.target) {
        cursor: crosshair;
    }
</style>
