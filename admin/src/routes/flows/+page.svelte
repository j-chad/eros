<!-- src/routes/flows/+page.svelte -->
<script lang="ts">
    import Header from '$lib/components/Header.svelte';
    import {api} from '$lib/api';
    import Calendar from "./Calendar.svelte";
    import {NodeType, type StartNode} from "$lib/types";
    import {goto} from "$app/navigation";

    let { data } = $props();

    let nodes = $state<StartNode[]>(data.nodes);

    $effect(() => {
        nodes = data.nodes;
    });

    async function handleCreateGraph(date: Date) {
        const title = `Unnamed Graph`;

        try {
            const graphId = await api.graph.create({
                title,
                starting_at: date.toISOString(),
            });

            const newGraph: StartNode = {
                id: graphId,
                type: NodeType.START,
                title,
                description: '',
                start: {
                    starting_at: date.toISOString()
                },
                unlocked_at: null,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
            }

            nodes = [...nodes, newGraph];
            await goto(`/flows/${graphId}`);
        } catch (error) {
            console.error('Failed to create flow:', error);
        }
    }

    async function handleEditFlow(flowId: string) {
        await goto(`/flows/${flowId}`);
    }

    async function handleDeleteFlow(flowId: string) {
        if (confirm('Are you sure you want to delete this flow?')) {
            try {
                await api.graph.delete(flowId);
                nodes = nodes.filter(f => f.id !== flowId);
            } catch (error) {
                console.error('Failed to delete flow:', error);
            }
        }
    }
</script>

<Header title="Reward Flows Calendar" />

<Calendar nodes={nodes} onCreateFlow={handleCreateGraph} onEditFlow={handleEditFlow} onDeleteFlow={handleDeleteFlow} />