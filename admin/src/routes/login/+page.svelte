<script lang="ts">
    import { auth } from '$lib/auth.svelte';
    import { api } from '$lib/api';
    import { goto } from '$app/navigation';

    let key = $state('');
    let loading = $state(false);
    let error = $state('');

    async function handleSubmit() {
        error = '';
        loading = true;

        try {
            // Store the key temporarily to test it
            auth.login(key);

            // Test the key by making a request
            await api.registration.get();

            // Success - redirect to dashboard
            goto('/');
        } catch (err) {
            // Invalid key - clear it
            auth.logout();
            error = err instanceof Error ? err.message : 'Invalid API key';
        } finally {
            loading = false;
        }
    }
</script>

<div class="login-container">
    <div class="login-box">
        <h1>Admin Login</h1>
        <p>Enter your admin API key to continue</p>

        <form onsubmit={handleSubmit}>
            <div class="form-group">
                <label for="apiKey">API Key</label>
                <input
                        id="apiKey"
                        type="password"
                        bind:value={key}
                        placeholder="Enter admin API key"
                        required
                        disabled={loading}
                />
            </div>

            {#if error}
                <div class="error">
                    {error}
                </div>
            {/if}

            <button type="submit" disabled={loading || !key}>
                {loading ? 'Verifying...' : 'Login'}
            </button>
        </form>
    </div>
</div>

<style>
    .login-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        background: #f5f5f5;
    }

    .login-box {
        background: white;
        padding: 2rem;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        width: 100%;
        max-width: 400px;
    }

    h1 {
        margin: 0 0 0.5rem;
        font-size: 1.5rem;
    }

    p {
        margin: 0 0 2rem;
        color: #666;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }

    input {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
        box-sizing: border-box;
    }

    input:focus {
        outline: none;
        border-color: #4CAF50;
    }

    button {
        width: 100%;
        padding: 0.75rem;
        background: #4CAF50;
        color: white;
        border: none;
        border-radius: 4px;
        font-size: 1rem;
        cursor: pointer;
    }

    button:hover:not(:disabled) {
        background: #45a049;
    }

    button:disabled {
        background: #ccc;
        cursor: not-allowed;
    }

    .error {
        padding: 0.75rem;
        background: #ffebee;
        color: #c62828;
        border-radius: 4px;
        margin-bottom: 1rem;
    }
</style>