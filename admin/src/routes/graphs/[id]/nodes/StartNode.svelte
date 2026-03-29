<script lang="ts">
	import {useNodeConnections} from '@xyflow/svelte';
	import { Play } from 'lucide-svelte';
	import type {StartNode, StartNode as StartNodeType} from "$lib/types";
	import BaseNode from "./BaseNode.svelte";
	import type {NodeProps} from "./types";

	let { data }: NodeProps<StartNode> = $props();

	const node = $derived(data.node as StartNodeType);
	const connections = useNodeConnections();
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
    }}
/>
