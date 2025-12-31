<!-- src/lib/components/RewardFlowEditor.svelte -->
<script lang="ts">
    import {
        SvelteFlow,
        Controls,
        Background,
        BackgroundVariant,
        MiniMap,
        type Node,
        type Edge
    } from '@xyflow/svelte';
    import '@xyflow/svelte/dist/style.css';

    import RewardNode from '$lib/components/flow-nodes/RewardNode.svelte';
    import RevealNode from '$lib/components/flow-nodes/RevealNode.svelte';
    import GateNode from '$lib/components/flow-nodes/GateNode.svelte';

    type Reward = {
        id: string;
        title: string;
        notBefore: string;
        description?: string;
    };

    type Reveal = {
        id: string;
        rewardId: string;
        order: number;
        content: string;
    };

    type Gate = {
        id: string;
        revealId: string;
        type: string;
        unlockOrder: number;
        config: any;
    };

    let {
        reward,
        reveals,
        gates,
        onSave
    }: {
        reward: Reward;
        reveals: Reveal[];
        gates: Gate[];
        onSave?: (nodes: Node[], edges: Edge[]) => Promise<void>;
    } = $props();

    // Node types registry
    const nodeTypes = {
        reward: RewardNode,
        reveal: RevealNode,
        gate: GateNode
    };

    // Convert data to nodes and edges
    function buildFlowGraph(): { nodes: Node[], edges: Edge[] } {
        const nodes: Node[] = [];
        const edges: Edge[] = [];

        // Add reward node at the top
        nodes.push({
            id: `reward-${reward.id}`,
            type: 'reward',
            position: { x: 400, y: 50 },
            data: {
                ...reward,
                onEdit: () => handleEditReward(reward)
            }
        });

        // Add reveal nodes
        const sortedReveals = [...reveals].sort((a, b) => a.order - b.order);
        sortedReveals.forEach((reveal, index) => {
            const revealNodeId = `reveal-${reveal.id}`;

            nodes.push({
                id: revealNodeId,
                type: 'reveal',
                position: { x: 400, y: 200 + (index * 300) },
                data: {
                    ...reveal,
                    onEdit: () => handleEditReveal(reveal)
                }
            });

            // Connect from reward to first reveal, or from previous reveal
            if (index === 0) {
                edges.push({
                    id: `reward-to-${revealNodeId}`,
                    source: `reward-${reward.id}`,
                    target: revealNodeId,
                    animated: true,
                    label: `After ${new Date(reward.notBefore).toLocaleDateString()}`
                });
            } else {
                const prevRevealId = `reveal-${sortedReveals[index - 1].id}`;
                edges.push({
                    id: `${prevRevealId}-to-${revealNodeId}`,
                    source: prevRevealId,
                    target: revealNodeId,
                    animated: true,
                    label: 'All gates unlocked'
                });
            }

            // Add gate nodes for this reveal
            const revealGates = gates
                .filter(g => g.revealId === reveal.id)
                .sort((a, b) => a.unlockOrder - b.unlockOrder);

            revealGates.forEach((gate, gateIndex) => {
                const gateNodeId = `gate-${gate.id}`;

                nodes.push({
                    id: gateNodeId,
                    type: 'gate',
                    position: {
                        x: 150 + (gateIndex * 200),
                        y: 250 + (index * 300)
                    },
                    data: {
                        ...gate,
                        onEdit: () => handleEditGate(gate)
                    }
                });

                // Connect reveal to gate
                edges.push({
                    id: `${revealNodeId}-to-${gateNodeId}`,
                    source: revealNodeId,
                    target: gateNodeId,
                    label: `Order ${gate.unlockOrder}`
                });
            });
        });

        return { nodes, edges };
    }

    let { nodes: flowNodes, edges: flowEdges } = $state(buildFlowGraph());

    // Rebuild graph when data changes
    $effect(() => {
        const result = buildFlowGraph();
        flowNodes = result.nodes;
        flowEdges = result.edges;
    });

    function handleEditReward(reward: Reward) {
        console.log('Edit reward:', reward);
        // Implement edit modal
    }

    function handleEditReveal(reveal: Reveal) {
        console.log('Edit reveal:', reveal);
        // Implement edit modal
    }

    function handleEditGate(gate: Gate) {
        console.log('Edit gate:', gate);
        // Implement edit modal
    }

    async function handleSave() {
        if (onSave) {
            await onSave(flowNodes, flowEdges);
        }
    }
</script>

<div class="flow-container">
    <div class="flow-header">
        <h3>Reward Flow: {reward.title}</h3>
        <button class="btn-save" onclick={handleSave}>Save Changes</button>
    </div>

    <div class="flow-wrapper">
        <SvelteFlow
                nodes={flowNodes}
                edges={flowEdges}
                {nodeTypes}
                fitView
                minZoom={0.5}
                maxZoom={1.5}
        >
            <Controls />
            <Background variant={BackgroundVariant.Dots} />
            <MiniMap />
        </SvelteFlow>
    </div>
</div>

<style>
    .flow-container {
        height: 100%;
        display: flex;
        flex-direction: column;
        background: white;
        border-radius: 8px;
        overflow: hidden;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .flow-header {
        padding: 1rem 1.5rem;
        background: #f9fafb;
        border-bottom: 1px solid #e5e7eb;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .flow-header h3 {
        margin: 0;
        font-size: 1.125rem;
        font-weight: 600;
        color: #1f2937;
    }

    .btn-save {
        padding: 0.5rem 1rem;
        background: #3b82f6;
        color: white;
        border: none;
        border-radius: 4px;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: background 0.2s;
    }

    .btn-save:hover {
        background: #2563eb;
    }

    .flow-wrapper {
        flex: 1;
        position: relative;
    }
</style>