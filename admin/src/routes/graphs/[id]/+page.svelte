<script lang="ts">
    import { Save } from 'lucide-svelte';
    import Card from "$lib/components/Card.svelte";
    import Button from "$lib/components/Button.svelte";
    import EditableField from "$lib/components/EditableField.svelte";
    import Header from "$lib/components/Header.svelte";
    import DateDisplay from "$lib/components/DateDisplay.svelte";
    import type { Graph } from '$lib/types';
    import { api } from "$lib/api";

    let { data } = $props();

    let graph = $state<Graph>(data.graph);
    let description = $derived(() => graph.description ?? '');
    let hasChanges = $state(false);
    let isSaving = $state(false);

    async function handleSave() {
        isSaving = true;
        try {
            // Assuming you have an API method like this
            // await api.graph.update(graph.id, {
            //     title: graph.title,
            //     description: graph.description,
            //     starting_at: graph.starting_at,
            //     // Add viewport, nodes, edges when ready
            // });
            hasChanges = false;
        } catch (error) {
            console.error('Failed to save graph:', error);
            alert('Failed to save changes');
        } finally {
            isSaving = false;
        }
    }

    function handleTitleChange(newValue: string) {
        graph.title = newValue;
        hasChanges = true;
    }

    function handleDescriptionChange(newValue: string) {
        graph.description = newValue;
        hasChanges = true;
    }

    function handleStartingTimeChange(newValue: string) {
        graph.starting_at = newValue;
        hasChanges = true;
    }
</script>

<Header title={graph.title} />

<!-- Graph Details Section -->
<Card title="Graph Details">
    {#snippet actions()}
        <Button
                variant="primary"
                onclick={handleSave}
                disabled={!hasChanges || isSaving}
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
<!--                placeholder="Add a description..."-->
    </div>

    <div class="grid-3">
        <div class="form-group">
            <label>Starting At</label>
            <input
                    type="datetime-local"
                    bind:value={graph.starting_at}
                    onchange={() => {
                    hasChanges = true;
                }}
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

<!-- Graph Canvas Section - Placeholder for now -->
<Card title="Graph Canvas">
    <div class="canvas-placeholder">
        Graph visualization will go here (SvelteFlow component)
    </div>
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

    .canvas-placeholder {
        min-height: 500px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f9fafb;
        border: 2px dashed #d1d5db;
        border-radius: 8px;
        color: #6b7280;
        font-size: 0.875rem;
    }

    @media (max-width: 768px) {
        .grid-3 {
            grid-template-columns: 1fr;
        }
    }
</style>