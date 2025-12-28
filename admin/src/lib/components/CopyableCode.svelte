<script lang="ts">
    import { Eye, EyeOff, Copy, Check } from 'lucide-svelte';
    import IconButton from './IconButton.svelte';

    let { code }: { code: string } = $props();

    let visible = $state(false);
    let copied = $state(false);

    async function handleCopy() {
        await navigator.clipboard.writeText(code);
        copied = true;
        setTimeout(() => copied = false, 2000);
    }
</script>

<div class="copyable-code">
    <div class="code-display">
        {visible ? code : '••••••••••••••••••••••'}
    </div>
    <IconButton
            onclick={() => visible = !visible}
            title={visible ? 'Hide code' : 'Show code'}
    >
        {#if visible}
            <EyeOff size={18} />
        {:else}
            <Eye size={18} />
        {/if}
    </IconButton>
    <IconButton
            onclick={handleCopy}
            title="Copy code"
    >
        {#if copied}
            <Check size={18} class="success" />
        {:else}
            <Copy size={18} />
        {/if}
    </IconButton>
</div>

<style>
    .copyable-code {
        display: flex;
        gap: 0.5rem;
    }

    .code-display {
        flex: 1;
        padding: 0.75rem 1rem;
        background: #f9fafb;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        font-family: 'Courier New', monospace;
        font-size: 0.875rem;
    }

    :global(.success) {
        color: #10b981;
    }
</style>