<script lang="ts">
	import {
		Background,
		BackgroundVariant,
		Controls,
		MiniMap, type NodeEvents,
		type NodeTypes, type PaneEvents,
		SvelteFlow, type Viewport,
		type Connection
	} from '@xyflow/svelte';

	import type { Node as FlowNode, Edge as FlowEdge } from './nodes/types';

    import '@xyflow/svelte/dist/style.css';
    import {type AnyNode, type Graph, NodeType} from "$lib/types";
    import {debounce} from "$lib/utils";
    import StartNode from "./nodes/StartNode.svelte";
    import PaneContextMenu from "./PaneContextMenu.svelte";
	import NodeContextMenu from "./NodeContextMenu.svelte";
	import LocationNode from "./nodes/LocationNode.svelte";
	import EditNodeDialog from "./edit-node-dialog/EditNodeDialog.svelte";
	import CodeNode from "./nodes/CodeNode.svelte";
	import ManualNode from "./nodes/ManualNode.svelte";

    let {graph = $bindable<Graph>()}: { graph: Graph } = $props()

    const nodeTypes: NodeTypes = {
		[NodeType.MANUAL]: ManualNode,
		[NodeType.CODE]: CodeNode,
		[NodeType.LOCATION]: LocationNode,
        [NodeType.START]: StartNode
    };

    let contextMenu = $state<'node' | 'pane' | 'hidden'>('hidden');
    let contextMenuPosition = $state({ x: 0, y: 0 });
	let selectedNodeId = $state<string | null>(null);
	let editingNode = $state<AnyNode | null>(null);

    // Non-reactive state for performance
    let nodes = $state.raw<FlowNode[]>([]);
    let edges = $state.raw<FlowEdge[]>([]);

    // prevent graph<->flow feedback loops
    let syncingFromGraph = false;
    let syncingFromFlow = false;

	function findNodeById(id: string, required: true): FlowNode
	function findNodeById(id: string, required?: false): FlowNode | undefined
	function findNodeById(id: string, required=false): FlowNode | undefined {
		const node = nodes.find((n) => n.id === id);
		if (required && !node) {
			throw new Error(`Node with ID ${id} not found`);
		}

		return node;
	}

    function nodeToFlowNode(node: AnyNode): FlowNode {
        return {
            id: node.id,
            position: node.ui_position ?? {x: 0, y: 0},
            data: {
				node,
				onEdit: handleEditNode
			},
            type: node.type,
			deletable: node.type !== NodeType.START,
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
		if (!targetNode) throw new Error('No target node on drag stop');

		const node = findNodeById(targetNode.id, true);
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

	function handleDeleteItem({nodes: deletedNodes, edges: deletedEdges}: { nodes: FlowNode[]; edges: FlowEdge[] }) {
		// ensure start node cannot be deleted
		deletedNodes = deletedNodes.filter((n) => n.type !== NodeType.START);

		if (deletedNodes.length > 0) {
			const deletedNodeIds = new Set(deletedNodes.map((n) => n.id));
			nodes = nodes.filter((n) => !deletedNodeIds.has(n.id));
			edges = edges.filter((e) => !deletedNodeIds.has(e.source) && !deletedNodeIds.has(e.target));
		}

		if (deletedEdges.length > 0) {
			const deletedEdgeIds = new Set(deletedEdges.map((e) => e.id));
			edges = edges.filter((e) => !deletedEdgeIds.has(e.id));
		}

		commitGraph();
	}

	function handleConnect(connection: Connection) {
		if (!connection.source || !connection.target) return;

		const newEdge: FlowEdge = {
			id: crypto.randomUUID(),
			source: connection.source,
			target: connection.target,
			label: '',
			data: {
				edge: {
					id: crypto.randomUUID(),
					from: connection.source,
					to: connection.target,
					choice_label: '',
					created_at: new Date().toISOString(),
					updated_at: new Date().toISOString()
				}
			}
		};

		edges = [...edges, newEdge];
		commitGraph();
	}

	function handleEditNode(nodeId: string) {
		const node = findNodeById(nodeId, true);
		editingNode = node.data.node;
	}

	function handleUpdateNode(updatedNode: AnyNode) {
		nodes = nodes.map((node) => {
			return node.id === updatedNode.id ? nodeToFlowNode(updatedNode) : node;
		});
		commitGraph();
	}
</script>

<div class="canvas">
    <SvelteFlow {nodes} {edges} {nodeTypes} initialViewport={graph.viewport} fitView={graph.viewport === undefined}
                onnodedragstop={handleNodeDragStop}
                onpanecontextmenu={handlePaneContextMenu}
				onnodecontextmenu={handleNodeContextMenu}
				onmoveend={handleViewportChange}
				ondelete={handleDeleteItem}
				onconnect={handleConnect}
    >
        <MiniMap/>
        <Controls/>
        <Background variant={BackgroundVariant.Dots}/>
        {#if contextMenu === 'pane'}
            <PaneContextMenu x={contextMenuPosition.x} y={contextMenuPosition.y} onClose={handleCloseContextMenu} onCreateNode={handleNewNode}/>
		{:else if contextMenu === 'node' && selectedNodeId !== null}
			<NodeContextMenu nodeId={selectedNodeId} x={contextMenuPosition.x} y={contextMenuPosition.y} onClose={handleCloseContextMenu}/>
        {/if}
    </SvelteFlow>
</div>

<EditNodeDialog node={editingNode} onClose={() => editingNode = null} onSave={handleUpdateNode}></EditNodeDialog>

<style>
    .canvas {
        width: 100%;
        height: 600px;
        border: 2px dashed #d1d5db;
        border-radius: 8px;
    }
</style>
