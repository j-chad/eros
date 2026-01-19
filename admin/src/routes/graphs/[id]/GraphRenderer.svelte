<script lang="ts">
	import {
		Background,
		BackgroundVariant,
		Controls,
		type Edge as FlowEdge,
		MiniMap,
		type Node as FlowNode, type NodeEvents,
		type NodeTypes, type PaneEvents,
		SvelteFlow, type Viewport
	} from '@xyflow/svelte';

    import '@xyflow/svelte/dist/style.css';
    import {type AnyNode, type Edge, type Graph, type Node, NodeType} from "$lib/types";
    import {debounce} from "$lib/utils";
    import StartNode from "./nodes/StartNode.svelte";
    import PaneContextMenu from "./PaneContextMenu.svelte";
	import NodeContextMenu from "./NodeContextMenu.svelte";
	import LocationNode from "./nodes/LocationNode.svelte";

    let {graph = $bindable<Graph>()}: { graph: Graph } = $props()

    const nodeTypes: NodeTypes = {
		[NodeType.LOCATION]: LocationNode,
        [NodeType.START]: StartNode
    };

    let contextMenu = $state<'node' | 'pane' | 'hidden'>('hidden');
    let contextMenuPosition = $state({ x: 0, y: 0 });
	let selectedNodeId = $state<string | null>(null);

    // Non-reactive state for performance
    let nodes = $state.raw<FlowNode<{ node: AnyNode }, NodeType>[]>([]);
    let edges = $state.raw<FlowEdge<{ edge: Edge }>[]>([]);

    // prevent graph<->flow feedback loops
    let syncingFromGraph = false;
    let syncingFromFlow = false;

	function findNodeById(id: string): FlowNode<{ node: AnyNode }, NodeType> | undefined {
		return nodes.find((n) => n.id === id);
	}

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

	function handleDeleteNode(nodeId: string) {
		nodes = nodes.filter((n) => n.id !== nodeId);
		edges = edges.filter((e) => e.source !== nodeId && e.target !== nodeId);
		commitGraph();
	}

    const handlePaneContextMenu: PaneEvents['onpanecontextmenu'] = ({event}) => {
        event.preventDefault()
        contextMenu = 'pane'
        contextMenuPosition = {
            x: event.clientX,
            y: event.clientY
        };
    };

	const handleNodeContextMenu: NodeEvents['onnodecontextmenu'] = ({event, node}) => {
		event.preventDefault()

		if (node.type === NodeType.START) {
			// Don't show context menu for start node
			return;
		}

		selectedNodeId = node.id;
		contextMenu = 'node'
		contextMenuPosition = {
			x: event.clientX,
			y: event.clientY
		};
	};

    function handleCloseContextMenu() {
        contextMenu = 'hidden';
    }

    function handleNewNode(node: AnyNode) {
        const flowNode = nodeToFlowNode(node);
        nodes = [...nodes, flowNode];
        commitGraph();
    }

    const handleNodeDragStop: NodeEvents['onnodedragstop'] = ({targetNode}) => {
		if (!targetNode) return;
		const node = findNodeById(targetNode.id);
		if (!node) return;
		node.position = targetNode.position;
		commitGraph();
	};

	function handleViewportChange(_: unknown, viewport: Viewport) {
		syncingFromFlow = true;

		graph = {
			...graph,
			viewport: viewport
		};

		queueMicrotask(() => {
			syncingFromFlow = false;
		});
	}
</script>

<div class="canvas">
    <SvelteFlow {nodes} {edges} {nodeTypes} initialViewport={graph.viewport} fitView={graph.viewport === undefined}
                onnodedragstop={handleNodeDragStop}
                onpanecontextmenu={handlePaneContextMenu}
				onnodecontextmenu={handleNodeContextMenu}
				onmoveend={handleViewportChange}
    >
        <MiniMap/>
        <Controls/>
        <Background variant={BackgroundVariant.Dots}/>
        {#if contextMenu === 'pane'}
            <PaneContextMenu x={contextMenuPosition.x} y={contextMenuPosition.y} onClose={handleCloseContextMenu} onCreateNode={handleNewNode}/>
		{:else if contextMenu === 'node' && selectedNodeId !== null}
			<NodeContextMenu nodeId={selectedNodeId} x={contextMenuPosition.x} y={contextMenuPosition.y} onClose={handleCloseContextMenu} onDelete={handleDeleteNode}/>
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
