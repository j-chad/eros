<!-- src/routes/flows/+page.svelte -->
<script lang="ts">
    import Header from '$lib/components/Header.svelte';
    import { api } from '$lib/api';
    import Calendar from "./Calendar.svelte";
    import type {StartNode} from "$lib/types";
    import {goto} from "$app/navigation";

    let { data } = $props();

    let nodes = $state<StartNode[]>(data.nodes);

    $effect(() => {
        nodes = data.nodes;
    });

    async function handleCreateFlow(date: Date) {
        // Navigate to create flow page with date
        console.log('Create flow for:', date);
        // window.location.href = `/flows/create?date=${date.toISOString()}`;
    }

    async function handleEditFlow(flowId: string) {
        await goto(`/flows/${flowId}`);
    }

    async function handleDeleteFlow(flowId: string) {
        if (confirm('Are you sure you want to delete this flow?')) {
            try {
                await api.flows.delete(flowId);
                nodes = nodes.filter(f => f.id !== flowId);
            } catch (error) {
                console.error('Failed to delete flow:', error);
            }
        }
    }
</script>

<Header title="Reward Flows Calendar" />

<Calendar nodes={nodes} onCreateFlow={handleCreateFlow} onEditFlow={handleEditFlow} onDeleteFlow={handleDeleteFlow} />