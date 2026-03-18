<script lang="ts">
	import { onDestroy } from 'svelte';
	import BrandHeader from '$lib/ui/BrandHeader.svelte';
	import Card from '$lib/ui/base/Card.svelte';
	import type { GraphSummary } from '$lib/domain/graph.types';

	const { data }: { data: { graph: GraphSummary | null } } = $props();

	// Tick every second so derived values stay current.
	let now = $state(Date.now());
	let intervalId: ReturnType<typeof setInterval> | null = setInterval(() => {
		now = Date.now();
	}, 1000);

	onDestroy(() => {
		if (intervalId != null) clearInterval(intervalId);
	});

	// --- derived countdown ---
	const target = $derived(data.graph ? new Date(data.graph.starting_at) : null);

	const msLeft = $derived(target ? Math.max(0, target.getTime() - now) : null);

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

	const isReady = $derived(timeLeft !== null && timeLeft.total === 0);
	const hasGraph = $derived(data.graph !== null && timeLeft !== null);
</script>

<svelte:head>
	<title>Eros</title>
	<meta content="width=device-width, initial-scale=1" name="viewport" />
</svelte:head>

<div class="mx-auto min-h-dvh max-w-md px-4 py-6 flex flex-col justify-between gap-8">
	<BrandHeader subtitle="Something's coming" />

	<div class="flex-1 flex flex-col justify-center gap-6">
		{#if !hasGraph}
			<Card>
				<p class="text-center text-sm opacity-70 animate-popIn">
					Nothing scheduled yet. Check back soon.
				</p>
			</Card>
		{:else if isReady}
			<Card>
				<div class="flex flex-col items-center gap-4 py-4 animate-popIn">
					<div class="text-5xl font-extrabold text-primary">It's time.</div>
					{#if data.graph?.title}
						<p class="text-center text-base font-semibold">{data.graph.title}</p>
					{/if}
					{#if data.graph?.description}
						<p class="text-center text-sm opacity-70">{data.graph.description}</p>
					{/if}
				</div>
			</Card>
		{:else if timeLeft}
			<div class="flex flex-col gap-5 animate-popIn">
				{#if data.graph?.title}
					<div class="text-center">
						<h1 class="text-xl font-extrabold">{data.graph.title}</h1>
						{#if data.graph?.description}
							<p class="mt-1 text-sm opacity-70">{data.graph.description}</p>
						{/if}
					</div>
				{/if}

				<Card>
					<div class="flex flex-col items-center gap-6 py-2">
						<p class="text-sm font-semibold opacity-60 tracking-wide uppercase">Starts in</p>

						<div class="grid grid-cols-4 gap-3 w-full">
							{#each [
								{ value: timeLeft.days, label: 'days' },
								{ value: timeLeft.hours, label: 'hours' },
								{ value: timeLeft.minutes, label: 'min' },
								{ value: timeLeft.seconds, label: 'sec' },
							] as unit (unit.label)}
								<div class="flex flex-col items-center gap-1.5">
									<div class="bg-base-200 rounded-2xl w-full aspect-square flex items-center justify-center">
										<span class="countdown font-mono text-3xl sm:text-4xl font-bold">
											<span
												style="--value:{unit.value};"
												aria-live="polite"
												aria-label="{unit.value} {unit.label}"
											></span>
										</span>
									</div>
									<span class="text-xs font-semibold opacity-60">{unit.label}</span>
								</div>
							{/each}
						</div>

						<div class="w-full rounded-2xl bg-base-200/60 px-4 py-3 text-center">
							<span class="text-xs font-semibold opacity-60">Unlocks</span>
							<div class="mt-0.5 text-sm font-semibold">
								{target?.toLocaleDateString(undefined, {
									weekday: 'long',
									month: 'long',
									day: 'numeric',
								})}
								<span class="font-normal opacity-70">
									at {target?.toLocaleTimeString(undefined, {
										hour: 'numeric',
										minute: '2-digit',
									})}
								</span>
							</div>
						</div>
					</div>
				</Card>
			</div>
		{/if}
	</div>

	<div class="text-center text-xs opacity-40">Eros</div>
</div>
