<script lang="ts">
    import {
        SvelteFlow,
        MiniMap,
        Controls,
        Background,
        BackgroundVariant,
        type Edge as FlowEdge,
        type Node as FlowNode,
        type NodeTypes, type NodeTargetEventWithPointer
    } from '@xyflow/svelte';

    import '@xyflow/svelte/dist/style.css';
    import {type AnyNode, type Edge, type Graph, NodeType} from "$lib/types";
    import {debounce} from "$lib/utils";
    import StartNode from "./nodes/StartNode.svelte";

    let { graph = $bindable<Graph>() }: {graph: Graph} = $props()

    const nodeTypes: NodeTypes = {
        [NodeType.START]: StartNode
    };

    // Non-reactive state for performance
    let nodes = $state.raw<FlowNode<{node: AnyNode}, NodeType>[]>([]);
    let edges = $state.raw<FlowEdge<{edge: Edge}>[]>([]);

    // prevent graph<->flow feedback loops
    let syncingFromGraph = false;
    let syncingFromFlow = false;

    // Sync FROM graph TO flow (when graph changes)
    $effect(() => {
        console.log("Syncing from graph to flow", {graph, syncingFromGraph, syncingFromFlow});

        if (!graph || syncingFromFlow) return;

        syncingFromGraph = true;

        nodes = graph.nodes?.map((node) => ({
            id: node.id,
            position: node.ui_position ?? { x: 0, y: 0 },
            data: {node},
            type: node.type
        })) ?? [];

        edges = graph.edges.map((edge) => ({
            id: edge.id,
            source: edge.from,
            target: edge.to,
            label: edge.choice_label,
            data: {edge}
        })) ?? [];

        syncingFromGraph = false;
    });

    const commitGraph = debounce(() => {
        if (!graph || syncingFromGraph) return;

        syncingFromFlow = true;

        graph = {
            ...graph,
            nodes: nodes.map((n) => ({
                ...n.data.node,
                ui_position: {x: n.position.x, y: n.position.y}
            })),
            edges: edges.map((e) => ({
                from: e.source,
                to: e.target,
                choice_label: e.label,
                created_at: e.data?.edge.created_at ?? new Date().toISOString(),
                updated_at: e.data?.edge.updated_at ?? new Date().toISOString(),
                id: e.id
            }))
        };

        syncingFromFlow = false;
    }, 150);

    function handleNodeDragStop() {
        commitGraph();
    }
</script>

<div class="canvas">
    <SvelteFlow {nodes} {edges} {nodeTypes} fitView onnodedragstop={handleNodeDragStop}>
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