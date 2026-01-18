<script lang="ts">
    import {SvelteFlow, MiniMap, Controls, Background, BackgroundVariant, type Edge, type Node} from '@xyflow/svelte';

    import '@xyflow/svelte/dist/style.css';
    import type {Graph, Node as ErosNode, Edge as ErosEdge} from "$lib/types";

    let { graph = $bindable<Graph>() } = $props()

    // Non-reactive state for performance
    let nodes = $state.raw<Node[]>([]);
    let edges = $state.raw<Edge[]>([]);

    // Sync FROM graph TO flow (when graph changes)
    $effect(() => {
        if (!graph) return;

        nodes = graph.nodes?.map((node: ErosNode) => ({
            id: node.id,
            position: node.ui_position || { x: 0, y: 0 },
            data: {
                label: node.title,
                ...node
            },
            type: node.type
        })) || [];

        edges = graph.edges?.map((edge: ErosEdge) => ({
            id: edge.id,
            source: edge.from,
            target: edge.to,
            label: edge.choice_label,
            data: edge
        })) || [];
    });
</script>

<div class="canvas">
    <SvelteFlow nodes={nodes} edges={edges} fitView>
        <MiniMap/>
        <Controls/>
        <Background variant={BackgroundVariant.Dots}/>
    </SvelteFlow>
</div>

<style>
    .canvas {
        width: 100%;
        height: 600px;
        border: 2px dashed #d1d5db;
        border-radius: 8px;
    }
</style>