<script lang="ts">
    import FavourChoiceTable from "./FavourChoiceTable.svelte";
    import Header from "$lib/components/Header.svelte";
    import FavourCounter from "./FavourCounter.svelte";
    import {api} from "$lib/api";
    import FulfilmentTable from "./FulfilmentTable.svelte";
    import Card from "$lib/components/Card.svelte";

    let {data} = $props();

    async function handleFavourUpdate(newCount: number) {
        await api.favours.updateFavourCount(newCount)
    }

    async function handleFulfilmentChange(favourId: string, fulfilled: boolean) {
        let favour = data.requests.find(f => f.id === favourId);
        if (favour) {
            favour.fulfilled = fulfilled;
        }
        // await api.favours.updateFulfilmentStatus(favourId, fulfilled);
    }
</script>

<Header title="Favour Management"/>

<div class="stack">
    <FavourCounter count={data.favourCount} onUpdate={handleFavourUpdate}/>
    <FavourChoiceTable choices={data.choices}/>
    <Card title="Fulfilment Requests">
        <FulfilmentTable onToggleFulfilled={handleFulfilmentChange} favours={data.requests}/>
    </Card>
</div>

<style>
    .stack {
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }
</style>