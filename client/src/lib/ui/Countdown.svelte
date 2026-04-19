<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	const { target, label = 'Unlocks in', onComplete }: {
		target: Date;
		label?: string;
		onComplete?: () => void;
	} = $props();

	let now = $state(Date.now());
	let intervalId: ReturnType<typeof setInterval> | null = null;

	const msLeft = $derived(Math.max(0, target.getTime() - now));
	const isDone = $derived(msLeft === 0);

	const timeLeft = $derived({
		days: Math.floor(msLeft / 86_400_000),
		hours: Math.floor((msLeft % 86_400_000) / 3_600_000),
		minutes: Math.floor((msLeft % 3_600_000) / 60_000),
		seconds: Math.floor((msLeft % 60_000) / 1_000),
	});

	function startTick() {
		if (intervalId) return;
		intervalId = setInterval(() => {
			now = Date.now();
		}, 1000);
	}

	function stopTick() {
		if (intervalId) {
			clearInterval(intervalId);
			intervalId = null;
		}
	}

	onMount(() => {
		if (!isDone) startTick();
	});

	onDestroy(() => {
		stopTick();
	});

	$effect(() => {
		if (isDone) {
			stopTick();
			onComplete?.();
		}
	});
</script>

{#if !isDone}
	<div class="flex flex-col items-center gap-6 py-2">
		<p class="text-sm font-semibold opacity-60 tracking-wide uppercase">{label}</p>

		<div class="flex flex-col sm:grid sm:grid-cols-4 gap-3 w-full">
			{#each [
				{ value: timeLeft.days, label: 'days' },
				{ value: timeLeft.hours, label: 'hours' },
				{ value: timeLeft.minutes, label: 'min' },
				{ value: timeLeft.seconds, label: 'sec' },
			] as unit (unit.label)}
				<div
					class="flex sm:flex-col items-center gap-4 sm:gap-1.5 bg-base-200 rounded-2xl px-5 py-3"
				>
					<span
						class="countdown font-mono text-5xl sm:text-4xl font-bold w-16 sm:w-auto text-center"
					>
						<span
							style="--value:{unit.value};"
							aria-live="polite"
							aria-label="{unit.value} {unit.label}"
						></span>
					</span>
					<span class="text-sm sm:text-xs font-semibold opacity-60">{unit.label}</span>
				</div>
			{/each}
		</div>
	</div>
{/if}
