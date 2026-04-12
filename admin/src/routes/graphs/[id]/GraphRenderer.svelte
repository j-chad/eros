<script lang="ts">
	import {
		Background,
		BackgroundVariant,
		type Connection,
		Controls,
		type EdgeEvents,
		MiniMap,
		type NodeEvents,
		type NodeTypes,
		type PaneEvents,
		SvelteFlow,
		type Viewport
	} from '@xyflow/svelte';

	import type {Edge as FlowEdge, Node as FlowNode} from './nodes/types';

	import '@xyflow/svelte/dist/style.css';
	import {type AnyNode, type Edge, type Graph, NodeType} from "$lib/types";
	import {debounce} from "$lib/utils";
	import StartNode from "./nodes/StartNode.svelte";
	import PaneContextMenu from "./PaneContextMenu.svelte";
	import LocationNode from "./nodes/LocationNode.svelte";
	import EditNodeDialog from "./edit-node-dialog/EditNodeDialog.svelte";
	import CodeNode from "./nodes/CodeNode.svelte";
	import ManualNode from "./nodes/ManualNode.svelte";
	import RewardNode from "./nodes/RewardNode.svelte";
	import EditEdgeDialog from "./EditEdgeDialog.svelte";
	import {Eye, EyeOff} from 'lucide-svelte';
	import {api} from '$lib/api';

	let {graph = $bindable<Graph>()}: { graph: Graph } = $props()

	const nodeTypes: NodeTypes = {
		[NodeType.REWARD]: RewardNode,
		[NodeType.MANUAL]: ManualNode,
		[NodeType.CODE]: CodeNode,
		[NodeType.LOCATION]: LocationNode,
		[NodeType.START]: StartNode
	};

	let contextMenu = $state<'pane' | 'hidden'>('hidden');
	let contextMenuPosition = $state({x: 0, y: 0});
	let editingNode = $state<AnyNode | null>(null);
	let editingEdge = $state<Edge | null>(null);
	let showProgress = $state(false);

	// Non-reactive state for performance
	let nodes = $state.raw<FlowNode[]>([]);
	let edges = $state.raw<FlowEdge[]>([]);

	// prevent graph<->flow feedback loops
	let syncingFromGraph = false;
	let syncingFromFlow = false;

	function findNodeById(id: string, required: true): FlowNode
	function findNodeById(id: string, required = false): FlowNode | undefined {
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
				onEdit: handleEditNode,
				onUpdateData: (data) => handleUpdateNodeData(node.id, data),
				showProgress,
				onToggleUnlock: handleToggleUnlock
			},
			type: node.type,
			deletable: node.type !== NodeType.START,
		};
	}

	async function handleToggleUnlock(nodeId: string) {
		const node = graph.nodes.find(n => n.id === nodeId);
		if (!node) return;

		try {
			if (node.unlocked_at) {
				await api.node.lock(nodeId);
				node.unlocked_at = null;
			} else {
				await api.node.unlock(nodeId);
				node.unlocked_at = new Date().toISOString();
			}
			// Trigger reactivity
			graph = {...graph, nodes: [...graph.nodes]};
		} catch (e) {
			console.error('Failed to toggle node unlock:', e);
			alert('Failed to toggle node unlock state');
		}
	}

	// Build a lookup of unlocked node IDs for edge styling
	let unlockedNodeIds = $derived(new Set(
		graph?.nodes?.filter(n => n.unlocked_at).map(n => n.id) ?? []
	));

	// Sync FROM graph TO flow (when graph changes)
	$effect(() => {
		if (!graph || syncingFromFlow) return;
		// track showProgress so toggling it re-syncs nodes
		void showProgress;

		syncingFromGraph = true;

		nodes = graph.nodes?.map(nodeToFlowNode) ?? [];

		edges = graph.edges.map((edge) => ({
			id: edge.id,
			source: edge.from,
			target: edge.to,
			label: edge.choice_label,
			data: {edge},
			animated: showProgress && unlockedNodeIds.has(edge.from) && unlockedNodeIds.has(edge.to),
			style: showProgress
				? (unlockedNodeIds.has(edge.from) && unlockedNodeIds.has(edge.to)
					? 'stroke: #10b981; stroke-width: 2.5px;'
					: 'stroke: #d1d5db; stroke-width: 1.5px; opacity: 0.5;')
				: undefined
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

		console.log('Committed graph:', graph);

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

	const handleEdgeClick: EdgeEvents<FlowEdge>['onedgeclick'] = ({event, edge}) => {
		event.stopPropagation();
		if (edge.data?.edge) {
			editingEdge = edge.data.edge;
		}
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

		const id = crypto.randomUUID();

		const newEdge: FlowEdge = {
			id,
			source: connection.source,
			target: connection.target,
			label: '',
			data: {
				edge: {
					id,
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

	function handleUpdateNodeData(nodeId: string, newData: AnyNode['data']) {
		const node = findNodeById(nodeId, true);
		const updatedNode = {
			...node.data.node,
			data: newData
		};
		nodes = nodes.map((n) => n.id === nodeId ? nodeToFlowNode(updatedNode as AnyNode) : n);
		commitGraph();
	}

	function handleUpdateEdge(updatedEdge: Edge) {
		edges = edges.map((e) => {
			if (e.id === updatedEdge.id) {
				return {
					...e,
					label: updatedEdge.choice_label,
					data: {edge: updatedEdge}
				};
			}
			return e;
		});
		commitGraph();
	}
</script>

<div class="canvas">
	<SvelteFlow {edges} fitView={graph.viewport === undefined} initialViewport={graph.viewport} {nodeTypes} {nodes}
				onconnect={handleConnect}
				ondelete={handleDeleteItem}
				onedgeclick={handleEdgeClick}
				onmoveend={handleViewportChange}
				onnodedragstop={handleNodeDragStop}
				onpanecontextmenu={handlePaneContextMenu}
	>
		<MiniMap/>
		<Controls/>
		<Background variant={BackgroundVariant.Dots}/>
		<div class="progress-toggle">
			<button
				class="progress-btn"
				class:active={showProgress}
				onclick={() => showProgress = !showProgress}
				title={showProgress ? 'Hide client progress' : 'Show client progress'}
			>
				{#if showProgress}
					<EyeOff size={14} />
				{:else}
					<Eye size={14} />
				{/if}
				<span>Progress</span>
			</button>
		</div>
		{#if contextMenu === 'pane'}
			<PaneContextMenu x={contextMenuPosition.x} y={contextMenuPosition.y} onClose={handleCloseContextMenu}
							 onCreateNode={handleNewNode}/>
		{/if}
	</SvelteFlow>
</div>

<EditNodeDialog node={editingNode} onClose={() => editingNode = null} onSave={handleUpdateNode}></EditNodeDialog>
<EditEdgeDialog
	allEdges={graph.edges}
	edge={editingEdge}
	onClose={() => editingEdge = null}
	onSave={handleUpdateEdge}
/>

<style>
	.canvas {
		width: 100%;
		height: 600px;
		border: 2px dashed #d1d5db;
		border-radius: 8px;
	}

	.progress-toggle {
		position: absolute;
		top: 10px;
		right: 10px;
		z-index: 5;
	}

	.progress-btn {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.5rem 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: white;
		font-size: 0.75rem;
		font-weight: 600;
		color: #374151;
		cursor: pointer;
		transition: all 0.15s ease;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
	}

	.progress-btn:hover {
		background: #f9fafb;
		border-color: #9ca3af;
	}

	.progress-btn.active {
		background: #ecfdf5;
		border-color: #10b981;
		color: #065f46;
	}
</style>
