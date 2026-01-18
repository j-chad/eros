<script>
    import {SvelteFlow, MiniMap, Controls, Background, BackgroundVariant} from '@xyflow/svelte';

    import '@xyflow/svelte/dist/style.css';

    let nodes = $state.raw([
        { id: '1', position: { x: 0, y: 0 }, data: { label: '1' } },
        { id: '2', position: { x: 0, y: 100 }, data: { label: '2' } },
    ]);

    let edges = $state.raw([
        { id: 'e1-2', source: '1', target: '2' },
    ]);
</script>

<div class="canvas">
    <SvelteFlow bind:nodes bind:edges>
        <MiniMap/>
        <Controls/>
        <Background variant={BackgroundVariant.Dots}/>
    </SvelteFlow>
</div>

<style>
    .canvas {
        /*min-height: 500px;*/
        width: 100%;
        height: 600px;
        border: 2px dashed #d1d5db;
        border-radius: 8px;
    }
</style>