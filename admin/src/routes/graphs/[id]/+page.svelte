<script lang="ts">
    import { Save } from 'lucide-svelte';
    import Card from "$lib/components/Card.svelte";
    import Button from "$lib/components/Button.svelte";
    import EditableField from "$lib/components/EditableField.svelte";
    import Header from "$lib/components/Header.svelte";
    import DateDisplay from "$lib/components/DateDisplay.svelte";
	import type {APIError, Graph} from '$lib/types';
    import { api } from "$lib/api";
    import GraphRenderer from "./GraphRenderer.svelte";

    let { data } = $props();

    let graph = $state<Graph>(data.graph);
    let description = $derived(() => graph.description ?? '');
    let isSaving = $state(false);

    // Convert UTC datetime to local date (YYYY-MM-DD)
    let startingAtDate = $derived(() => {
        if (!graph.starting_at) return '';
        const date = new Date(graph.starting_at);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        return `${year}-${month}-${day}`;
    });

    async function handleSave() {
        isSaving = true;
        try {
            await api.graph.update(graph.id, graph)
        } catch (error) {
            alert(`Failed to save graph: ${(error as App.Error).body as APIError['internal']}`);
        } finally {
            isSaving = false;
        }
    }

    function handleTitleChange(newValue: string) {
        graph.title = newValue;
    }

    function handleDescriptionChange(newValue: string) {
        graph.description = newValue;
    }

    function handleStartingDateChange(event: Event) {
        const input = event.target as HTMLInputElement;
        const dateValue = input.value;

        if (dateValue) {
            const localDate = new Date(dateValue + 'T00:00:00');
            localDate.setHours(8, 0, 0, 0);
            graph.starting_at = localDate.toISOString();
        }
    }
</script>

<Header title={graph.title} />

<!-- Graph Details Section -->
<Card title="Graph Details">
    {#snippet actions()}
        <Button
                variant="primary"
                onclick={handleSave}
                disabled={isSaving}
        >
            {#snippet icon()}
                <Save size={16} />
            {/snippet}
            {isSaving ? 'Saving...' : 'Save Changes'}
        </Button>
    {/snippet}

    <div class="form-group">
        <label>Title</label>
        <EditableField
                bind:value={graph.title}
                onSave={handleTitleChange}
        />
    </div>

    <div class="form-group">
        <label>Description</label>
        <EditableField
                value={description()}
                onSave={handleDescriptionChange}
                multiline={true} />
    </div>

    <div class="grid-3">
        <div class="form-group">
            <label>Starting At</label>
            <input
                    type="date"
                    value={startingAtDate()}
                    onchange={handleStartingDateChange}
                    class="datetime-input"
            />
        </div>
        <div>
            <label>Created</label>
            <DateDisplay datetime={graph.created_at} inline />
        </div>
        <div>
            <label>Last Updated</label>
            <DateDisplay datetime={graph.updated_at} inline />
        </div>
    </div>
</Card>

<Card title="Graph Canvas">
    <GraphRenderer bind:graph/>
</Card>

<style>
    .form-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        font-size: 0.875rem;
        font-weight: 500;
        color: #4b5563;
        margin-bottom: 0.5rem;
    }

    .grid-3 {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr;
        gap: 1.5rem;
        align-items: end;
    }

    .datetime-input {
        width: 100%;
        padding: 0.5rem 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 4px;
        font-size: 0.875rem;
        color: #1f2937;
        background: white;
        transition: border-color 0.2s;
    }

    .datetime-input:focus {
        outline: none;
        border-color: #3b82f6;
        box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    @media (max-width: 768px) {
        .grid-3 {
            grid-template-columns: 1fr;
        }
    }
</style>
