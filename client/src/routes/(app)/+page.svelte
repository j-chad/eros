<script lang="ts">
	import BrandHeader from '$lib/ui/BrandHeader.svelte';
	import Card from '$lib/ui/base/Card.svelte';
	import Calendar from '$lib/ui/Calendar.svelte';
	import type { GraphSummary } from '$lib/types/graph';
	import {goto} from "$app/navigation";

	const { data }: { data: { graphs: GraphSummary[] } } = $props();

	// --- determine which view to show ---

	const pastGraphs = $derived(data.graphs.filter((g) => g.starting_at.getTime() <= Date.now()));
	const hasUnlocked = $derived(pastGraphs.length > 0);

	// For countdown: pick the soonest upcoming graph.
	const nextGraph = $derived(
		data.graphs
			.filter((g) => g.starting_at.getTime() > Date.now())
			.sort((a, b) => a.starting_at.getTime() - b.starting_at.getTime())[0] ?? null,
	);

	// --- countdown logic (only used when no graphs unlocked yet) ---

	let tick = $state(Date.now());
	let intervalId: ReturnType<typeof setInterval> | null = null;

	// Only run the timer when we need the countdown view.
	$effect(() => {
		if (!hasUnlocked && nextGraph) {
			intervalId = setInterval(() => {
				tick = Date.now();
			}, 1000);
		}

		return () => {
			if (intervalId != null) {
				clearInterval(intervalId);
				intervalId = null;
			}
		};
	});

	const msLeft = $derived(nextGraph ? Math.max(0, nextGraph.starting_at.getTime() - tick) : null);

	const timeLeft = $derived(
		msLeft !== null
			? {
					days: Math.floor(msLeft / 86_400_000),
					hours: Math.floor((msLeft % 86_400_000) / 3_600_000),
					minutes: Math.floor((msLeft % 3_600_000) / 60_000),
					seconds: Math.floor((msLeft % 60_000) / 1_000),
					total: msLeft,
				}
			: null,
	);
</script>

<svelte:head>
	<title>Eros</title>
	<meta content="width=device-width, initial-scale=1" name="viewport" />
</svelte:head>

<div class="mx-auto min-h-dvh max-w-md px-4 py-6 flex flex-col justify-between gap-8">
	<BrandHeader subtitle={hasUnlocked ? undefined : "Something's coming"} />

	<div class="flex-1 flex flex-col justify-center gap-6" style="view-transition-name: card-content">
		{#if hasUnlocked}
			<!-- Calendar view: at least one graph has been unlocked -->
			<Card>
				<Calendar graphs={data.graphs} onSelectDate={() => {
					goto('/graphs');
				}} />
			</Card>
		{:else if data.graphs.length === 0}
			<!-- No graphs at all -->
			<Card>
				<p class="text-center text-sm opacity-70 animate-popIn">
					Nothing scheduled yet. Check back soon.
				</p>
			</Card>
		{:else if timeLeft && timeLeft.total === 0}
			<!-- Countdown reached zero -->
			<Card>
				<div class="flex flex-col items-center gap-4 py-4 animate-popIn">
					<div class="text-5xl font-extrabold text-primary">It's time.</div>
					{#if nextGraph?.title}
						<p class="text-center text-base font-semibold">{nextGraph.title}</p>
					{/if}
					{#if nextGraph?.description}
						<p class="text-center text-sm opacity-70">{nextGraph.description}</p>
					{/if}
				</div>
			</Card>
		{:else if timeLeft && nextGraph}
			<!-- Countdown to next graph -->
			<div class="flex flex-col gap-5 animate-popIn">
				{#if nextGraph.title}
					<div class="text-center">
						<h1 class="text-xl font-extrabold">{nextGraph.title}</h1>
						{#if nextGraph.description}
							<p class="mt-1 text-sm opacity-70">{nextGraph.description}</p>
						{/if}
					</div>
				{/if}

				<Card>
					<div class="flex flex-col items-center gap-6 py-2">
						<p class="text-sm font-semibold opacity-60 tracking-wide uppercase">Starts in</p>

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
				</Card>
			</div>
		{/if}
	</div>

	<div class="text-center text-xs opacity-40">Eros</div>
</div>
