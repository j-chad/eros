<script lang="ts">
    let { expiresAt }: { expiresAt: string } = $props();

    function formatDate(dateString: string) {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function getTimeRemaining(expiresAt: string) {
        const now = new Date();
        const expiry = new Date(expiresAt);
        const diff = expiry.getTime() - now.getTime();

        if (diff < 0) return 'Expired';

        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);

        if (days > 0) return `${days}d ${hours}h`;
        if (hours > 0) return `${hours}h ${minutes}m`;
        if (minutes > 0) return `${minutes}m ${seconds}s`;
        return `${seconds}s`;
    }

    let formattedDate = $derived(formatDate(expiresAt));
    let timeRemaining = $derived(getTimeRemaining(expiresAt));

    setInterval(() => {
        timeRemaining = getTimeRemaining(expiresAt);
    }, 1000);
</script>

<div class="expiry">
    {formattedDate}
    <span class="time-remaining">{timeRemaining}</span>
</div>

<style>
    .expiry {
        font-size: 0.875rem;
        color: #1f2937;
    }

    .time-remaining {
        margin-left: 0.5rem;
        color: #f97316;
        font-weight: 500;
    }
</style>