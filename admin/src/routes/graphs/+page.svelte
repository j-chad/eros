<!-- src/routes/graphs/+page.svelte -->
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
            await handleEditGraph(graphId)
        } catch (error) {
            console.error('Failed to create graph:', error);
        }
    }

    async function handleEditGraph(graphId: string) {
        await goto(`/graphs/${graphId}`);
    }

    async function handleDeleteGraph(graphId: string) {
        if (confirm('Are you sure you want to delete this graph?')) {
            try {
                await api.graph.delete(graphId);
                nodes = nodes.filter(f => f.id !== graphId);
            } catch (error) {
                console.error('Failed to delete graph:', error);
            }
        }
    }
</script>

<Header title="Reward Graphs Calendar" />

<Calendar nodes={nodes} onCreateGraph={handleCreateGraph} onEditGraph={handleEditGraph} onDeleteGraph={handleDeleteGraph} />