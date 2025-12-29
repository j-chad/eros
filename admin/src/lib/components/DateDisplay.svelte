<script lang="ts">
    let {datetime, expiry = false, inline = false} = $props();

    let now = $state(new Date());

    export function formatDate(dateString: string) {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function getTimeRemaining(expiresAt: string) {
        const expiry = new Date(expiresAt);
        const diff = expiry.getTime() - now.getTime();

        if (diff < 0) return 'Expired';

        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);

        if (days > 0) return `${days}d ${hours}h`;
        if (hours > 0) return `${hours}h ${minutes}m`;
        if (minutes > 0) return `${minutes}m`;
        return `${seconds}s`;
    }

    function getTimePassed(datetime: string) {
        const pastDate = new Date(datetime);
        const diff = now.getTime() - pastDate.getTime();

        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);

        if (days > 0) return `${days}d ago`;
        if (hours > 0) return `${hours}h ago`;
        if (minutes > 0) return `${minutes}m ago`;
        return `${seconds}s ago`;
    }

    function getRelativeTime(datetime: string) {
        return expiry ? getTimeRemaining(datetime) : getTimePassed(datetime);
    }

    let formattedDate = $derived(formatDate(datetime));
    let relativeTime = $derived(getRelativeTime(datetime));

    setInterval(() => {
        now = new Date();
    }, 1000);
</script>

{#if expiry && !inline && relativeTime === 'Expired' }
    <div class="date expired" title={formattedDate}>
        {relativeTime}
    </div>
{:else}
    <div class="date" class:expiry title={inline ? '' : relativeTime}>
        {formattedDate}
        {#if inline }
            <span class='relative'>{relativeTime}</span>
        {/if}
    </div>
{/if}

<style>
    .date {
        font-size: 0.875rem;
        color: #1f2937;
    }

    .relative {
        margin-left: 0.5rem;
        color: #1f2937;
        font-weight: 500;
    }

    .expiry .relative, .expired {
        color: #f97316;
    }
</style>