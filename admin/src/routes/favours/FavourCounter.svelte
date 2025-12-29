<script lang="ts">
    import { Plus, Minus, Check, X } from 'lucide-svelte';

    let {
        count = $bindable(0),
        onUpdate,
        min = 0,
        max = 999
    }: {
        count?: number;
        onUpdate?: (newCount: number) => Promise<void>;
        min?: number;
        max?: number;
    } = $props();

    let originalCount = $state(count);
    let pendingCount = $state(count);
    let isUpdating = $state(false);

    // Sync originalCount when count prop changes externally
    $effect(() => {
        originalCount = count;
        pendingCount = count;
    });

    let delta = $derived(pendingCount - originalCount);
    let hasPendingChanges = $derived(delta !== 0);

    function handleIncrement() {
        if (pendingCount >= max) return;
        pendingCount++;
    }

    function handleDecrement() {
        if (pendingCount <= min) return;
        pendingCount--;
    }

    async function handleCommit() {
        if (!hasPendingChanges || isUpdating) return;

        if (!confirm(`Are you sure you want to ${delta > 0 ? 'increase' : 'decrease'} the favour count by ${Math.abs(delta)}?`)) {
            return;
        }

        isUpdating = true;

        try {
            if (onUpdate) {
                await onUpdate(pendingCount);
            }
            count = pendingCount;
            originalCount = pendingCount;
        } catch (err) {
            console.error('Failed to update count:', err);
            // Revert on error
            pendingCount = originalCount;
        } finally {
            isUpdating = false;
        }
    }

    function handleCancel() {
        pendingCount = originalCount;
    }

    let canDecrement = $derived(pendingCount > min && !isUpdating);
    let canIncrement = $derived(pendingCount < max && !isUpdating);
</script>

<div class="favour-counter">
    <button
            class="counter-btn decrement"
            onclick={handleDecrement}
            disabled={!canDecrement}
            title="Decrease"
    >
        <Minus size={24} />
    </button>

    <div class="counter-center">
        <div class="counter-label">Available Favours</div>
        <div class="counter-display">
            <span class="current-count">{count}</span>
            {#if hasPendingChanges}
                <button class="delta" onclick={handleCommit} class:positive={delta > 0} class:negative={delta < 0}>
                    {delta > 0 ? '+' : ''}{delta}
                </button>
            {/if}
        </div>
    </div>

    <button
            class="counter-btn increment"
            onclick={handleIncrement}
            disabled={!canIncrement}
            title="Increase"
    >
        <Plus size={24} />
    </button>
</div>

<style>
    .favour-counter {
        display: flex;
        align-items: center;
        gap: 1rem;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border-radius: 12px;
        padding: 1.5rem 2rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        width: 100%;
        transition: all 0.3s;
    }

    .counter-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 3.5rem;
        height: 3.5rem;
        background: rgba(255, 255, 255, 0.2);
        color: white;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-radius: 10px;
        cursor: pointer;
        transition: all 0.2s;
        flex-shrink: 0;
    }

    .counter-btn:hover:not(:disabled) {
        background: rgba(255, 255, 255, 0.3);
        border-color: rgba(255, 255, 255, 0.5);
        transform: scale(1.05);
    }

    .counter-btn:active:not(:disabled) {
        transform: scale(0.95);
    }

    .counter-btn:disabled {
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(255, 255, 255, 0.1);
        color: rgba(255, 255, 255, 0.4);
        cursor: not-allowed;
    }

    .counter-center {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .counter-label {
        font-size: 0.875rem;
        font-weight: 600;
        color: rgba(255, 255, 255, 0.9);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .counter-display {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .current-count {
        font-size: 3rem;
        font-weight: bold;
        color: white;
        line-height: 1;
        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .delta {
        font-size: 1.125rem;
        font-weight: 600;
        padding: 0.25rem 0.5rem;
        border-radius: 6px;
        animation: slideIn 0.2s ease-out;
        cursor: pointer;
    }

    .delta.positive {
        background: rgba(16, 185, 129, 0.3);
        color: #d1fae5;
        border: 1px solid rgba(16, 185, 129, 0.5);
    }

    .delta.positive:hover {
        background: rgba(16, 185, 129, 0.5);
        border-color: rgba(16, 185, 129, 0.8);
    }

    .delta.negative {
        background: rgba(239, 68, 68, 0.3);
        color: #fecaca;
        border: 1px solid rgba(239, 68, 68, 0.5);
    }

    .delta.negative:hover {
        background: rgba(239, 68, 68, 0.5);
        border-color: rgba(239, 68, 68, 0.8);
    }

    .action-buttons {
        display: flex;
        gap: 0.5rem;
        flex-shrink: 0;
        animation: slideIn 0.2s ease-out;
    }

    @keyframes slideIn {
        from {
            opacity: 0;
            transform: translateX(10px);
        }
        to {
            opacity: 1;
            transform: translateX(0);
        }
    }

    @media (max-width: 768px) {
        .favour-counter {
            padding: 1rem 1.5rem;
            flex-wrap: wrap;
        }

        .counter-btn {
            width: 3rem;
            height: 3rem;
        }

        .current-count {
            font-size: 2.5rem;
        }

        .delta {
            font-size: 1rem;
        }

        .counter-label {
            font-size: 0.75rem;
        }

        .action-buttons {
            width: 100%;
            justify-content: center;
            margin-top: 0.5rem;
        }
    }

    @media (max-width: 480px) {
        .favour-counter {
            padding: 1rem;
        }

        .counter-btn {
            width: 2.5rem;
            height: 2.5rem;
        }

        .counter-btn :global(svg) {
            width: 20px;
            height: 20px;
        }

        .current-count {
            font-size: 2rem;
        }

        .delta {
            font-size: 0.875rem;
        }

        .action-buttons {
            flex-direction: column;
            width: 100%;
        }

        .action-buttons :global(button) {
            width: 100%;
        }
    }
</style>