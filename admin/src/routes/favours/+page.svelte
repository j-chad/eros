<script lang="ts">
    import FavourChoiceTable from "./FavourChoiceTable.svelte";
    import Header from "$lib/components/Header.svelte";
    import FavourCounter from "./FavourCounter.svelte";
    import {api} from "$lib/api";

    let {data} = $props();

    async function handleFavourUpdate(newCount: number) {
        await api.favours.updateFavourCount(newCount)
    }
</script>

<Header title="Favour Management"/>

<div class="stack">
    <FavourCounter count={data.favourCount} onUpdate={handleFavourUpdate} />
    <FavourChoiceTable choices={data.choices}/>
</div>

<style>
    .stack {
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }
</style>