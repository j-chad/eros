<script lang="ts">
    import type { Snippet } from 'svelte';

    let {
        onclick,
        title,
        children
    }: {
        onclick?: () => void;
        title?: string;
        children: Snippet;
    } = $props();
</script>

<button
        class="icon-btn"
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
</style>