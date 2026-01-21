<script lang="ts">
    import { page } from '$app/state';
    import { dev } from '$app/environment';

    // Access the error from the page store
    $: error = page.error;
    $: status = page.status;
</script>

<div class="error-page">
    <div class="error-container">
        <h1>{status}</h1>
        <h2>{error?.message || 'An error occurred'}</h2>

        {#if dev}
            <!-- Show detailed error info in development -->
            <div class="error-details">
                <h3>Error Details</h3>

                {#if error}
                    <pre class="error-object">{JSON.stringify(error, null, 2)}</pre>
                {/if}
            </div>
        {/if}

        <div class="actions">
            <a href="/">Go Home</a>
            <button onclick={() => window.location.reload()}>Reload Page</button>
        </div>
    </div>
</div>

<style>
    .error-page {
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f9fafb;
        padding: 2rem;
    }

    .error-container {
        max-width: 800px;
        width: 100%;
        height: 100%;
        background: white;
        border-radius: 8px;
        padding: 2rem;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    h1 {
        font-size: 4rem;
        color: #ef4444;
        margin: 0;
    }

    h2 {
        font-size: 1.5rem;
        color: #1f2937;
        margin: 0.5rem 0 1.5rem 0;
    }

    h3 {
        font-size: 1.125rem;
        color: #374151;
        margin: 1.5rem 0 1rem 0;
    }

    .error-details {
        background: #f3f4f6;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        padding: 1rem;
        margin-top: 1.5rem;
    }

    pre {
        background: #1f2937;
        color: #f9fafb;
        padding: 1rem;
        border-radius: 4px;
        overflow-x: auto;
        font-size: 0.875rem;
        line-height: 1.5;
    }

    .actions {
        display: flex;
        gap: 1rem;
        margin-top: 2rem;
    }

    .actions a,
    .actions button {
        padding: 0.75rem 1.5rem;
        border-radius: 4px;
        font-weight: 500;
        text-decoration: none;
        cursor: pointer;
        transition: background 0.2s;
    }

    .actions a {
        background: #3b82f6;
        color: white;
    }

    .actions a:hover {
        background: #2563eb;
    }

    .actions button {
        background: #f3f4f6;
        color: #1f2937;
        border: 1px solid #e5e7eb;
    }

    .actions button:hover {
        background: #e5e7eb;
    }
</style>
