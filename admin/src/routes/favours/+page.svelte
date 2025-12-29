<script lang="ts">
    import FavourChoiceTable from "./FavourChoiceTable.svelte";
    import Header from "$lib/components/Header.svelte";

    let {data} = $props();
</script>

<Header title="Favour Management"/>
<FavourChoiceTable choices={data.choices}/>