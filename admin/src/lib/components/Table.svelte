<script lang="ts">
    import type { Snippet } from 'svelte';

    let {
        headers,
        children,
        empty
    }: {
        headers: string[];
        children?: Snippet;
        empty?: string;
    } = $props();
</script>

<div class="table-container">
    <table>
        <thead>
        <tr>
            {#each headers as header}
                <th>{header}</th>
            {/each}
        </tr>
        </thead>
        <tbody>
        {#if children}
            {@render children()}
        {:else}
            <tr>
                <td colspan={headers.length} class="empty">
                    {empty || 'No data available'}
                </td>
            </tr>
        {/if}
        </tbody>
    </table>
</div>

<style>
    .table-container {
        overflow-x: auto;
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    thead {
        background: #f9fafb;
        border-bottom: 1px solid #e5e7eb;
    }

    th {
        padding: 0.75rem 1.5rem;
        text-align: left;
        font-size: 0.75rem;
        font-weight: 500;
        color: #4b5563;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    :global(tbody tr) {
        border-bottom: 1px solid #e5e7eb;
        transition: background 0.2s;
    }

    :global(tbody tr:hover) {
        background: #f9fafb;
    }

    :global(td) {
        padding: 1rem 1.5rem;
        font-size: 0.875rem;
        color: #4b5563;
    }

    td.empty {
        text-align: center;
        color: #6b7280;
    }
</style>