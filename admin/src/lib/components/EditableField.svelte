<script lang="ts">
    import {Check, Pencil, X} from 'lucide-svelte';
    import IconButton from './IconButton.svelte';

    let {
        value = $bindable(),
        onSave,
        multiline = false
    }: {
        value: string;
        onSave?: (newValue: string) => (Promise<void> | void);
        multiline?: boolean;
    } = $props();

    let isEditing = $state(false);
    let editValue = $state(value);
    let isSaving = $state(false);

    async function handleSave() {
        if (editValue === value) {
            isEditing = false;
            return;
        }

        isSaving = true;
        try {
            if (onSave) {
                await onSave(editValue);
            }
            value = editValue; // Update the bindable value
            isEditing = false;
        } catch (err) {
            console.error('Failed to save:', err);
        } finally {
            isSaving = false;
        }
    }

    function handleCancel() {
        editValue = value;
        isEditing = false;
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Enter' && !multiline && !e.shiftKey) {
            e.preventDefault();
            handleSave();
        } else if (e.key === 'Escape') {
            handleCancel();
        }
    }

    function startEditing() {
        editValue = value;
        isEditing = true;
    }
</script>

<div class="editable-field">
    {#if isEditing}
        <div class="edit-mode">
            {#if multiline}
                <textarea
                        bind:value={editValue}
                        onkeydown={handleKeydown}
                        disabled={isSaving}
                        rows="3"></textarea>
            {:else}
                <input
                        type="text"
                        bind:value={editValue}
                        onkeydown={handleKeydown}
                        disabled={isSaving}
                />
            {/if}
            <div class="actions">
                <IconButton onclick={handleSave} title="Save">
                    <Check size={16} class="success"/>
                </IconButton>
                <IconButton onclick={handleCancel} title="Cancel">
                    <X size={16} class="danger"/>
                </IconButton>
            </div>
        </div>
    {:else}
        <div class="view-mode">
            <span class="value">{value}</span>
            <button
                    class="edit-btn"
                    onclick={startEditing}
                    title="Edit"
            >
                <Pencil size={14}/>
            </button>
        </div>
    {/if}
</div>

<style>
    /* Same styles as before */
    .editable-field {
        width: 100%;
    }

    .view-mode {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .value {
        flex: 1;
    }

    .edit-btn {
        padding: 0.25rem;
        background: transparent;
        border: none;
        color: #6b7280;
        cursor: pointer;
        border-radius: 4px;
        display: flex;
        align-items: center;
        opacity: 0;
        transition: opacity 0.2s;
    }

    .view-mode:hover .edit-btn {
        opacity: 1;
    }

    .edit-btn:hover {
        background: #f3f4f6;
        color: #1f2937;
    }

    .edit-mode {
        display: flex;
        gap: 0.5rem;
        align-items: flex-start;
    }

    input,
    textarea {
        flex: 1;
        padding: 0.5rem;
        border: 1px solid #3b82f6;
        border-radius: 4px;
        font-size: 0.875rem;
        font-family: inherit;
        outline: none;
    }

    input:disabled,
    textarea:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    textarea {
        resize: vertical;
        min-height: 60px;
        font-family: 'Courier New', monospace;
    }

    .actions {
        display: flex;
        gap: 0.25rem;
    }

    :global(.success) {
        color: #10b981;
    }

    :global(.danger) {
        color: #ef4444;
    }
</style>