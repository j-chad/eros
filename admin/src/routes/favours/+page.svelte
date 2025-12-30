<script lang="ts">
    import FavourChoiceTable from "./FavourChoiceTable.svelte";
    import Header from "$lib/components/Header.svelte";
    import FavourCounter from "./FavourCounter.svelte";
    import {api} from "$lib/api";
    import FulfilmentTable from "./FulfilmentTable.svelte";
    import Card from "$lib/components/Card.svelte";

    let {data} = $props();

    let requests = $state(data.requests);

    async function handleFavourUpdate(newCount: number) {
        await api.favours.updateFavourCount(newCount)
    }

    async function handleFulfilmentChange(favourId: string, fulfilled: boolean) {
        await api.favours.updateFavourRequestFulfilment(favourId, fulfilled);
        requests = requests.map(favour =>
            favour.id === favourId ? {...favour, fulfilled_at: fulfilled ? new Date().toISOString() : null} : favour
        );
    }
</script>

<Header title="Favour Management"/>

<div class="stack">
    <FavourCounter count={data.favourCount} onUpdate={handleFavourUpdate}/>
    <FavourChoiceTable choices={data.choices}/>
    <Card title="Fulfilment Requests">
        <FulfilmentTable onToggleFulfilled={handleFulfilmentChange} favours={requests}/>
    </Card>
</div>

<style>
    .stack {
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }
</style>