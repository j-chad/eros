<!-- StartNode.svelte -->
<script lang="ts">
	import {type Node, type NodeProps, useNodeConnections} from '@xyflow/svelte';
	import { Play } from 'lucide-svelte';
	import type {AnyNode, StartNode, StartNode as StartNodeType} from "$lib/types";
	import BaseNode from "./BaseNode.svelte";

	let { data }: NodeProps<Node<{node: StartNode}>> = $props();

	const node = $derived(data.node as StartNodeType);
	const connections = useNodeConnections();
	const isConnectable = $derived(connections.current.length === 0);
</script>

<BaseNode
	{node}
	config={{
        color: '#10b981',
        gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
        icon: Play,
        label: 'Start',
        hasTargetHandle: false,
        hasSourceHandle: true,
        isSourceConnectable: isConnectable
    }}
/>
