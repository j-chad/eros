<script lang="ts">
    import {
        Background,
        BackgroundVariant,
        Controls,
        type Edge as FlowEdge,
        MiniMap,
        type Node as FlowNode,
        type NodeTypes, type PaneEvents,
        SvelteFlow
    } from '@xyflow/svelte';

    import '@xyflow/svelte/dist/style.css';
    import {type AnyNode, type Edge, type Graph, type Node, NodeType} from "$lib/types";
    import {debounce} from "$lib/utils";
    import StartNode from "./nodes/StartNode.svelte";
    import PaneContextMenu from "./PaneContextMenu.svelte";

    let {graph = $bindable<Graph>()}: { graph: Graph } = $props()

    const nodeTypes: NodeTypes = {
        [NodeType.START]: StartNode
    };

    let showContextMenu = $state(false);
    let contextMenuPosition = $state({ x: 0, y: 0 });

    // Non-reactive state for performance
    let nodes = $state.raw<FlowNode<{ node: AnyNode }, NodeType>[]>([]);
    let edges = $state.raw<FlowEdge<{ edge: Edge }>[]>([]);

    // prevent graph<->flow feedback loops
    let syncingFromGraph = false;
    let syncingFromFlow = false;

    function nodeToFlowNode(node: AnyNode): FlowNode<{ node: AnyNode }, NodeType> {
        return {
            id: node.id,
            position: node.ui_position ?? {x: 0, y: 0},
            data: {node},
            type: node.type
        };
    }

    // Sync FROM graph TO flow (when graph changes)
    $effect(() => {
        if (!graph || syncingFromFlow) return;
        console.log("Syncing from graph to flow", {graph, syncingFromGraph, syncingFromFlow});

        syncingFromGraph = true;

        nodes = graph.nodes?.map(nodeToFlowNode) ?? [];

        edges = graph.edges.map((edge) => ({
            id: edge.id,
            source: edge.from,
            target: edge.to,
            label: edge.choice_label,
            data: {edge}
        })) ?? [];

        queueMicrotask(() => {
            syncingFromGraph = false;
        });
    });

    const commitGraph = debounce(() => {
        if (!graph || syncingFromGraph) return;
        console.log("Syncing from flow to graph", {nodes, edges, graph, syncingFromGraph, syncingFromFlow});

        syncingFromFlow = true;

        graph = {
            ...graph,
            nodes: nodes.map((n) => ({
                ...n.data.node,
                ui_position: n.position
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

        queueMicrotask(() => {
            syncingFromFlow = false;
        });
    }, 150);

    const handlePaneContextMenu: PaneEvents['onpanecontextmenu'] = ({event}) => {
        event.preventDefault()
        showContextMenu = true;
        contextMenuPosition = {
            x: event.clientX,
            y: event.clientY
        };
    };

    function handleCloseContextMenu() {
        showContextMenu = false;
    }

    function handleNewNode(node: AnyNode) {
        console.log("New node created, committing graph", {nodes, edges});
        const flowNode = nodeToFlowNode(node);
        nodes = [...nodes, flowNode];
        console.log("Nodes after addition", nodes);
        commitGraph();
    }

    function handleNodeDragStop() {
        console.log("Node drag stopped, committing graph", {nodes, edges});
        commitGraph();
    }
</script>

<div class="canvas">
    <SvelteFlow {nodes} {edges} {nodeTypes} fitView
                onnodedragstop={handleNodeDragStop}
                onpanecontextmenu={handlePaneContextMenu}
    >
        <MiniMap/>
        <Controls/>
        <Background variant={BackgroundVariant.Dots}/>
        {#if showContextMenu}
            <PaneContextMenu x={contextMenuPosition.x} y={contextMenuPosition.y} onClose={handleCloseContextMenu} onCreateNode={handleNewNode}/>
        {/if}
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