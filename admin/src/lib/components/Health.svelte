<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { api } from '$lib/api';

    let status: 'healthy' | 'degraded' | 'down' = $state('healthy');
    let intervalId: number | null = null;

    async function checkHealth() {
        try {
            await api.ping();
            status = 'healthy';
        } catch (error) {
            console.error('Health check failed:', error);
            status = 'down';
        }
    }

    onMount(() => {
        // Check immediately on mount
        checkHealth();

        // Then check every 30 seconds
        intervalId = window.setInterval(checkHealth, 30000);
    });

    onDestroy(() => {
        if (intervalId !== null) {
            clearInterval(intervalId);
        }
    });
</script>

<div class="health" title="API Status: {status}">
    <div class="dot dot-{status}"></div>
    <span class="label">API</span>
</div>

<style>
    .health {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.25rem 0.75rem;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 4px;
        cursor: help;
    }

    .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        animation: pulse 2s ease-in-out infinite;
    }

    .dot-healthy {
        background: #10b981;
        box-shadow: 0 0 8px rgba(16, 185, 129, 0.6);
    }

    .dot-degraded {
        background: #f59e0b;
        box-shadow: 0 0 8px rgba(245, 158, 11, 0.6);
    }

    .dot-down {
        background: #ef4444;
        box-shadow: 0 0 8px rgba(239, 68, 68, 0.6);
        animation: pulse-urgent 1s ease-in-out infinite;
    }

    .label {
        font-size: 0.75rem;
        color: rgba(255, 255, 255, 0.9);
        font-weight: 500;
    }

    @keyframes pulse {
        0%, 100% {
            opacity: 1;
        }
        50% {
            opacity: 0.5;
        }
    }

    @keyframes pulse-urgent {
        0%, 100% {
            opacity: 1;
            transform: scale(1);
        }
        50% {
            opacity: 0.7;
            transform: scale(1.2);
        }
    }
</style>