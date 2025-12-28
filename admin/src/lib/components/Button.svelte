<script lang="ts">
    import type { Snippet } from 'svelte';

    let {
        variant = 'primary',
        size = 'md',
        onclick,
        children,
        icon
    }: {
        variant?: 'primary' | 'danger' | 'secondary';
        size?: 'sm' | 'md';
        onclick?: () => void;
        children: Snippet;
        icon?: Snippet;
    } = $props();
</script>

<button
        class="btn btn-{variant} btn-{size}"
        {onclick}
>
    {#if icon}
        {@render icon()}
    {/if}
    {@render children()}
</button>

<style>
    .btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-weight: 500;
        transition: background 0.2s;
    }

    .btn-md {
        padding: 0.5rem 1rem;
        font-size: 0.875rem;
    }

    .btn-sm {
        padding: 0.375rem 0.75rem;
        font-size: 0.875rem;
    }

    .btn-primary {
        background: #3b82f6;
        color: white;
    }

    .btn-primary:hover {
        background: #2563eb;
    }

    .btn-danger {
        background: #ef4444;
        color: white;
    }

    .btn-danger:hover {
        background: #dc2626;
    }

    .btn-secondary {
        background: #f3f4f6;
        color: #1f2937;
        border: 1px solid #e5e7eb;
    }

    .btn-secondary:hover {
        background: #e5e7eb;
    }
</style>