<script lang="ts">
    let {title, children, actions = undefined} = $props();
</script>

<div class="card">
    <div class="card-header">
        <h3>{title}</h3>
        {#if actions}
            <div class="actions">
                {@render actions()}
            </div>
        {/if}
    </div>

    <div class="card-body">
        {@render children()}
    </div>
</div>

<style>
    .card {
        background: white;
        border-radius: 8px;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        margin-bottom: 2rem;
        overflow: hidden;
    }

    .card-header {
        padding: 1rem 1.5rem;
        background: #f9fafb;
        border-bottom: 1px solid #e5e7eb;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    h3 {
        font-size: 1.125rem;
        font-weight: 600;
        color: #1f2937;
        margin: 0;
    }

    .actions {
        display: flex;
        gap: 0.5rem;
    }

    .card-body {
        padding: 1.5rem;
    }

    @media (max-width: 768px) {
        .card-header {
            flex-direction: column;
            align-items: flex-start;
            gap: 1rem;
        }

        .actions {
            width: 100%;
            flex-direction: column;
        }
    }
</style>