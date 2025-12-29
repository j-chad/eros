<script lang="ts">
    import type { Snippet } from 'svelte';

    let {
        onclick,
        title,
        children,
        variant
    }: {
        onclick?: () => void;
        title?: string;
        children: Snippet;
        variant?: 'default' | 'primary'; // default is normal, primary is purple button
    } = $props();
</script>

<button
        class="icon-btn"
        class:primary={variant === 'primary'}
        {onclick}
        {title}
>
    {@render children()}
</button>

<style>
    .icon-btn {
        padding: 0.75rem;
        background: #f3f4f6;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        cursor: pointer;
        transition: background 0.2s;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .icon-btn:hover {
        background: #e5e7eb;
    }

    .icon-btn.primary {
        background: #8b5cf6;
        border-color: #7c3aed;
        color: white;
    }

    .icon-btn.primary:hover {
        background: #7c3aed;
    }
</style>