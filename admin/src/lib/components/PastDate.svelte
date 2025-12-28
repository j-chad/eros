<script lang="ts">
    import {formatDate} from "$lib/utils";

    let { datetime }: { datetime: string } = $props();

    let now = $state(new Date());

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

    let formattedDate = $derived(formatDate(datetime));
    let timePassed = $derived(getTimePassed(datetime));

    setInterval(() => {
        now = new Date();
    }, 1000);
</script>

<div class="time">
    {formattedDate}
    <span class="time-passed">{timePassed}</span>
</div>

<style>
    .time {
        font-size: 0.875rem;
        color: #1f2937;
    }

    .time-passed {
        margin-left: 0.5rem;
        color: #1f2937;
        font-weight: 500;
    }
</style>