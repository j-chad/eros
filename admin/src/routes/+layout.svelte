<script lang="ts">
    import {auth} from "$lib/auth.svelte";
    import {goto} from "$app/navigation";
    import {page} from '$app/state';
    import Health from "$lib/components/Health.svelte";
    import Impersonate from "$lib/components/Impersonate.svelte";

    let {children} = $props();

    const s3Origin = import.meta.env.VITE_S3_ORIGIN ?? '';
    const s3Csp = s3Origin ? ` ${s3Origin}` : '';
    const csp = `default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https://*.tile.openstreetmap.org https://tile.openstreetmap.org${s3Csp}; connect-src 'self' http://localhost:* https://*; media-src 'self' blob:${s3Csp}; font-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'`;

    function handleLogout() {
        auth.logout();
        goto('/login');
    }
</script>

<svelte:head>
    <meta http-equiv="Content-Security-Policy" content={csp} />
</svelte:head>

{#if auth.isAuthenticated}
    <!-- Authenticated Layout -->
    <div class="app">
        <nav>
            <div class="nav-content">
                <div class="left-section">
                    <h1>Eros Admin</h1>
                    <Impersonate/>
                    <Health/>
                </div>

                <div class="nav-links">
                    <a href="/" class:active={page.url.pathname === '/'}>Registration</a>
                    <a href="/favours" class:active={page.url.pathname === '/favours'}>Favours</a>
                    <a href="/graphs" class:active={page.url.pathname.startsWith('/graphs')}>Graphs</a>
                    <a href="/notifications" class:active={page.url.pathname === '/notifications'}>Notifications</a>
                    <button onclick={handleLogout} class="logout">Logout</button>
                </div>
            </div>
        </nav>

        <main>
            <div class="container">
                {@render children()}
            </div>
        </main>
    </div>
{:else}
    <!-- Unauthenticated - just show the page (login) -->
    {@render children()}
{/if}

<style>
    :global(html) {
        box-sizing: border-box;
    }

    :global(*, *:before, *:after) {
        box-sizing: inherit;
    }

    :global(body) {
        margin: 0;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    }

    .app {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
    }

    nav {
        background: #2c3e50;
        color: white;
        padding: 1rem;
    }

    .nav-content {
        max-width: 1200px;
        margin: 0 auto;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    h1 {
        margin: 0;
        font-size: 1.25rem;
    }

    .left-section {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .nav-links {
        display: flex;
        gap: 1rem;
        align-items: center;
    }

    .nav-links a {
        color: white;
        text-decoration: none;
        padding: 0.5rem 1rem;
        border-radius: 4px;
        transition: background 0.2s;
    }

    .nav-links a:hover,
    .nav-links a.active {
        background: rgba(255, 255, 255, 0.1);
    }

    .logout {
        background: #e74c3c;
        color: white;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.875rem;
    }

    .logout:hover {
        background: #c0392b;
    }

    main {
        flex: 1;
        max-width: 1200px;
        width: 100%;
        margin: 0 auto;
        padding: 2rem;
    }

    .container {
        max-width: 100%;
    }
</style>
