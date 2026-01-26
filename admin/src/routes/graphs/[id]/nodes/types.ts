import type {AnyNode, Edge as GraphEdge} from "$lib/types";
import type {Node as FlowNode, Edge as FlowEdge, NodeProps as FlowNodeProps} from "@xyflow/svelte";

interface NodeData<N extends AnyNode = AnyNode> extends Record<string, unknown> {
	node: N;
	onEdit: (nodeID: string) => void;
	onUpdateData: (nodeID: string, data: N['data']) => void;
}

export type Node<N extends AnyNode = AnyNode> = FlowNode<NodeData<N>, N["type"]>;
export type NodeProps<N extends AnyNode = AnyNode> = FlowNodeProps<Node<N>>;

interface EdgeData extends Record<string, unknown> {
	edge: GraphEdge
}
export type Edge = FlowEdge<EdgeData>;
