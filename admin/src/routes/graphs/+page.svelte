<script lang="ts">
    import Header from '$lib/components/Header.svelte';
    import {api} from '$lib/api';
    import Calendar from "./Calendar.svelte";
    import {NodeType, type Graph} from "$lib/types";
    import {goto} from "$app/navigation";

    let { data } = $props();

    let graphs = $state<Graph[]>(data.graphs);

    $effect(() => {
        graphs = data.graphs;
    });

    async function handleCreateGraph(date: Date) {
        const title = `Unnamed Graph`;

		const dateWithTime = new Date(date.getTime())
		dateWithTime.setHours(8, 0, 0, 0)

        try {
            const graphId = await api.graph.create({
                title,
                starting_at: dateWithTime.toISOString(),
            });

            const newGraph: Graph = {
                id: graphId,
                title,
                description: '',
                starting_at: dateWithTime.toISOString(),
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
				nodes: [],
				edges: []
            }

            graphs = [...graphs, newGraph];
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
                graphs = graphs.filter(f => f.id !== graphId);
            } catch (error) {
                console.error('Failed to delete graph:', error);
            }
        }
    }
</script>

<Header title="Reward Graphs Calendar" />

<Calendar graphs={graphs} onCreateGraph={handleCreateGraph} onEditGraph={handleEditGraph} onDeleteGraph={handleDeleteGraph} />
