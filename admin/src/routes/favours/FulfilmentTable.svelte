<script lang="ts">
    import { Check, X, MessageSquare } from 'lucide-svelte';
    import Table from "$lib/components/Table.svelte";
    import Button from "$lib/components/Button.svelte";
    import DateDisplay from "$lib/components/DateDisplay.svelte";
    import type {Favour} from "$lib/types";

    let {
        favours,
        onToggleFulfilled
    }: {
        favours: Favour[];
        onToggleFulfilled: (favourId: string, fulfilled: boolean) => Promise<void>;
    } = $props();

    let showFulfilled = $state(false);
    let loading = $state<string | null>(null);

    let filteredFavours = $derived(
        favours.filter(f => showFulfilled ? f.fulfilled_at !== null : f.fulfilled_at === null)
    );

    async function handleToggle(favourId: string, currentState: boolean) {
        loading = favourId;
        try {
            await onToggleFulfilled(favourId, !currentState);
        } finally {
            loading = null;
        }
    }
</script>

<div class="fulfilment-table">
    <div class="controls">
        <div class="toggle-group">
            <button
                    class="toggle-btn"
                    class:active={!showFulfilled}
                    onclick={() => showFulfilled = false}
            >
                Pending ({favours.filter(f => f.fulfilled_at === null).length})
            </button>
            <button
                    class="toggle-btn"
                    class:active={showFulfilled}
                    onclick={() => showFulfilled = true}
            >
                Fulfilled ({favours.filter(f => f.fulfilled_at !== null).length})
            </button>
        </div>
    </div>

    {#if filteredFavours.length > 0}
        <Table headers={['Choice', 'Message', 'Requested', 'Status']}>
            {#each filteredFavours as favour}
                <tr>
                    <td title={favour.choice_description ?? ''}>
                        <div class="choice">{favour.choice_label}</div>
                    </td>
                    <td>
                        {#if favour.message}
                            <div class="message-container">
                                <MessageSquare size={14} class="message-icon" />
                                <span class="message">{favour.message}</span>
                            </div>
                        {:else}
                            <span class="no-message">No message</span>
                        {/if}
                    </td>
                    <td class="nowrap"><DateDisplay datetime={favour.requested_at}/></td>
                    <td>
                        <Button
                                variant={favour.fulfilled_at !== null ? 'secondary' : 'primary'}
                                size="sm"
                                onclick={() => handleToggle(favour.id, favour.fulfilled_at !== null)}
                        >
                            {#snippet icon()}
                                {#if loading === favour.id}
                                    <div class="spinner"></div>
                                {:else if favour.fulfilled_at !== null}
                                    <X size={14} />
                                {:else}
                                    <Check size={14} />
                                {/if}
                            {/snippet}
                            {loading === favour.id
                                ? 'Updating...'
                                : favour.fulfilled_at !== null
                                    ? 'Mark Pending'
                                    : 'Mark Fulfilled'}
                        </Button>
                    </td>
                </tr>
            {/each}
        </Table>
    {:else}
        <div class="empty">
            {showFulfilled
                ? 'No fulfilled favours yet.'
                : 'No pending favours.'}
        </div>
    {/if}
</div>

<style>
    .fulfilment-table {
        width: 100%;
    }

    .controls {
        margin-bottom: 1.5rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .toggle-group {
        display: inline-flex;
        background: #f3f4f6;
        border-radius: 6px;
        padding: 0.25rem;
    }

    .toggle-btn {
        padding: 0.5rem 1rem;
        border: none;
        background: transparent;
        color: #6b7280;
        font-size: 0.875rem;
        font-weight: 500;
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.2s;
    }

    .toggle-btn:hover {
        color: #374151;
    }

    .toggle-btn.active {
        background: white;
        color: #1f2937;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .choice {
        font-weight: 500;
        color: #1f2937;
    }

    .message-container {
        display: flex;
        align-items: flex-start;
        gap: 0.5rem;
    }

    .message-container :global(.message-icon) {
        color: #6b7280;
        flex-shrink: 0;
        margin-top: 0.125rem;
    }

    .message {
        color: #374151;
        font-style: italic;
        line-height: 1.5;
    }

    .no-message {
        color: #9ca3af;
        font-style: italic;
    }

    :global(td.nowrap) {
        white-space: nowrap;
    }

    .empty {
        text-align: center;
        padding: 3rem 1rem;
        color: #6b7280;
        font-style: italic;
    }

    .spinner {
        width: 14px;
        height: 14px;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>