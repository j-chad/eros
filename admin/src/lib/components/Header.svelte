<script lang="ts">
    import type { Snippet } from 'svelte';

    let {
        title,
        actions
    }: {
        title: string;
        actions?: Snippet;
    } = $props();
</script>

<div class="header">
    <h2>{title}</h2>
    {#if actions}
        <div class="actions">
            {@render actions()}
        </div>
    {/if}
</div>

<style>
    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
        gap: 1rem;
    }

    h2 {
        font-size: 2rem;
        font-weight: bold;
        color: #1f2937;
        margin: 0;
    }

    .actions {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }

    @media (max-width: 768px) {
        .header {
            flex-direction: column;
            align-items: flex-start;
        }

        .actions {
            width: 100%;
            flex-direction: column;
        }

        .actions :global(button) {
            width: 100%;
        }
    }
</style>