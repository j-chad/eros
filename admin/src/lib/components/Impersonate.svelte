<script lang="ts">
    import { UserRoundCog } from 'lucide-svelte';
    import IconButton from "$lib/components/IconButton.svelte";

    async function handleImpersonate() {

    }
</script>

<IconButton onclick={handleImpersonate} title="Impersonate" variant="primary">
    <UserRoundCog size={16} />
</IconButton>

<style>
    .impersonate-btn {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: #8b5cf6;
        color: white;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.875rem;
        transition: background 0.2s;
    }

    .impersonate-btn:hover {
        background: #7c3aed;
    }
</style>