<!-- src/routes/rewards/[id]/+page.svelte -->
<script lang="ts">
    import RewardFlowEditor from './FlowEditor.svelte';
    import Header from '$lib/components/Header.svelte';

    let { data } = $props();

    // async function handleSave(nodes, edges) {
    //     console.log('Saving flow:', nodes, edges);
    //     // Implement your save logic here
    //     // await api.rewards.updateFlow(data.reward.id, nodes, edges);
    // }
</script>

<Header title="Reward Flow Editor" />

<div class="flow-page">
    <RewardFlowEditor
            gates={[]}
            reveals={[]}
            reward={{
                id: "reward1",
                title: "Sample Reward",
                notBefore: new Date().toISOString(),
                description: "This is a sample reward description."
            }}
    />
</div>

<style>
    .flow-page {
        height: calc(100vh - 200px);
        min-height: 600px;
    }
</style>