<script lang="ts">
    import { Plus, Trash2 } from 'lucide-svelte';
    import { api } from '$lib/api';
    import Card from '$lib/components/Card.svelte';
    import Button from '$lib/components/Button.svelte';
    import Table from '$lib/components/Table.svelte';
    import EditableField from '$lib/components/EditableField.svelte';
    import DateDisplay from "$lib/components/DateDisplay.svelte";
    import type {CreateFavourChoicePayload, FavourChoice} from "$lib/types";

    let { choices = $bindable<FavourChoice[]>([]) }: {
        choices: FavourChoice[];
    } = $props();
    let showAddChoiceForm = $state(false);

    let newChoice = $state({
        label: '',
        description: '',
        cost: 1,
        can_message: false
    } as CreateFavourChoicePayload);

    function resetNewFavourChoiceForm() {
        newChoice = {
            label: '',
            description: '',
            cost: 1,
            can_message: false
        };
        showAddChoiceForm = false;
    }

    async function handleAddFavourChoice() {
        if (!newChoice.label.trim()) {
            alert('Label is required');
            return;
        }

        try {
            const created = await api.favours.createChoice(newChoice);
            choices = [created, ...choices];
            resetNewFavourChoiceForm();
        } catch (err) {
            console.error('Failed to create favour:', err);
            alert('Failed to create favour');
        }
    }

    async function handleUpdateChoiceLabel(choiceID: string, newLabel: string) {
        if (!newLabel.trim()) {
            alert('Label cannot be empty');
            throw new Error('Label required');
        }

        const choice = choices.find(f => f.id === choiceID);
        if (!choice) return;

        await api.favours.updateChoice(choiceID, { ...choice, label: newLabel });
        choices = choices.map(f =>
            f.id === choiceID ? { ...f, label: newLabel, updated_at: new Date().toISOString() } : f
        );
    }

    async function handleUpdateChoiceDescription(choiceID: string, newDescription: string) {
        const choice = choices.find(f => f.id === choiceID);
        if (!choice) return;

        await api.favours.updateChoice(choiceID, { ...choice, description: newDescription });
        choices = choices.map(f =>
            f.id === choiceID ? { ...f, description: newDescription, updated_at: new Date().toISOString() } : f
        );
    }

    async function handleToggleChoiceMessaging(choiceID: string) {
        const choice = choices.find(f => f.id === choiceID);
        if (!choice) return;

        const currentValue = choice.can_message;
        const newValue = !currentValue;

        try {
            await api.favours.updateChoice(choiceID, { ...choice, can_message: newValue });
            choices = choices.map(f =>
                f.id === choiceID ? { ...f, can_message: newValue, updated_at: new Date().toISOString() } : f
            );
        } catch (err) {
            console.error('Failed to update messaging setting:', err);
            alert('Failed to update messaging setting');
        }
    }

    async function handleDeleteChoice(choiceID: string) {
        if (!confirm('Are you sure you want to delete this favour choice? This cannot be undone.')) {
            return;
        }

        try {
            await api.favours.deleteChoice(choiceID);
            choices = choices.filter(f => f.id !== choiceID);
        } catch (err) {
            console.error('Failed to delete favour:', err);
            alert('Failed to delete favour');
        }
    }
</script>

    <Card title="Favour Choices">
        {#snippet actions()}
            {#if !showAddChoiceForm}
                <Button variant="primary" onclick={() => showAddChoiceForm = true}>
                    {#snippet icon()}<Plus size={16} />{/snippet}
                    Add Choice
                </Button>
            {/if}
        {/snippet}

        {#if showAddChoiceForm}
            <div class="add-form">
                <h4>New Favour Choice</h4>

                <div class="form-group">
                    <label for="label">Label <span class="required">*</span></label>
                    <input
                            id="label"
                            type="text"
                            bind:value={newChoice.label}
                            placeholder="e.g., Coffee Date, Dinner, Movie Night"
                    />
                </div>

                <div class="form-group">
                    <label for="description">Description</label>
                    <textarea
                            id="description"
                            bind:value={newChoice.description}
                            placeholder="Optional description..."
                            rows="3"></textarea>
                </div>

                <div class="form-group">
                    <label for="cost">Cost (in favours)</label>
                    <input
                            id="cost"
                            type="number"
                            min="1"
                            bind:value={newChoice.cost}
                    />
                </div>

                <div class="form-group checkbox-group">
                    <label>
                        <input
                                type="checkbox"
                                bind:checked={newChoice.can_message}
                        />
                        Allow messaging after reveal
                    </label>
                </div>

                <div class="form-actions">
                    <Button variant="primary" onclick={handleAddFavourChoice}>
                        Add Favour
                    </Button>
                    <Button variant="secondary" onclick={resetNewFavourChoiceForm}>
                        Cancel
                    </Button>
                </div>
            </div>
            <div class="divider"></div>
        {/if}

        {#if choices.length > 0}
            <Table headers={['Label', 'Description', 'Cost', 'Allow Messaging', 'Actions']}>
                {#each choices as choice}
                    <tr>
                        <td>
                            <EditableField
                                    bind:value={choice.label}
                                    onSave={(newValue) => handleUpdateChoiceLabel(choice.id, newValue)}
                            />
                        </td>
                        <td>
                            <div class="description-cell">
                                <EditableField
                                        bind:value={choice.description}
                                        onSave={(newValue) => handleUpdateChoiceDescription(choice.id, newValue)}
                                        multiline={true}
                                />
                            </div>
                        </td>
                        <td>
                            <strong>{choice.cost}</strong>
                        </td>
                        <td>
                            <button
                                    class="toggle-btn"
                                    class:active={choice.can_message}
                                    onclick={() => handleToggleChoiceMessaging(choice.id)}
                            >
                                {choice.can_message ? 'Yes' : 'No'}
                            </button>
                        </td>
                        <td>
                            <Button
                                    variant="danger"
                                    size="sm"
                                    onclick={() => handleDeleteChoice(choice.id)}
                            >
                                {#snippet icon()}<Trash2 size={14} />{/snippet}
                                Delete
                            </Button>
                        </td>
                    </tr>
                {/each}
            </Table>
        {:else}
            <div class="empty">No favour choices yet. Click "Add Choice" to create one.</div>
        {/if}
    </Card>

<style>

    .add-form {
        padding: 1.5rem;
        background: #f9fafb;
        border-radius: 8px;
        margin-bottom: 1.5rem;
    }

    h4 {
        margin: 0 0 1.5rem 0;
        font-size: 1.125rem;
        font-weight: 600;
        color: #1f2937;
    }

    .form-group {
        margin-bottom: 1rem;
    }

    label {
        display: block;
        font-size: 0.875rem;
        font-weight: 500;
        color: #4b5563;
        margin-bottom: 0.5rem;
    }

    .required {
        color: #ef4444;
    }

    input[type="text"],
    textarea {
        width: 100%;
        padding: 0.5rem 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 4px;
        font-size: 0.875rem;
        font-family: inherit;
    }

    input[type="text"]:focus,
    textarea:focus {
        outline: none;
        border-color: #3b82f6;
        box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    textarea {
        resize: vertical;
    }

    .checkbox-group label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        cursor: pointer;
    }

    input[type="checkbox"] {
        width: 1rem;
        height: 1rem;
        cursor: pointer;
    }

    .form-actions {
        display: flex;
        gap: 0.5rem;
        margin-top: 1.5rem;
    }

    .divider {
        height: 1px;
        background: #e5e7eb;
        margin: 1.5rem 0;
    }

    .description-cell {
        max-width: 20rem;
    }

    .toggle-btn {
        padding: 0.25rem 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 4px;
        background: white;
        color: #6b7280;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
        min-width: 3rem;
    }

    .toggle-btn.active {
        background: #10b981;
        color: white;
        border-color: #10b981;
    }

    .toggle-btn:hover {
        background: #f3f4f6;
    }

    .toggle-btn.active:hover {
        background: #059669;
    }

    .empty {
        text-align: center;
        color: #6b7280;
        padding: 2rem;
    }

    :global(td.nowrap) {
        white-space: nowrap;
    }
</style>