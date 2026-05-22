<script lang="ts">
	import BrandHeader from '$lib/ui/BrandHeader.svelte';
	import Card from '$lib/ui/base/Card.svelte';
	import Calendar from '$lib/ui/Calendar.svelte';
	import Countdown from '$lib/ui/Countdown.svelte';
	import FavourFab from '$lib/ui/FavourFab.svelte';
	import type { GraphSummary } from '$lib/types/graph';
	import type { FavourCount, FavourRequest } from '$lib/types/favour';
	import {goto} from "$app/navigation";
	import {BellOff} from "lucide-svelte";
	import {isPushSubscribed, isPushSupported, subscribe} from "$lib/services/push";
	import {isOnline} from "$lib/online.svelte";

	const { data }: { data: { graphs: GraphSummary[]; favourCount: FavourCount; favourRequests: FavourRequest[] } } = $props();

	const online = $derived(isOnline());

	// FAB is visible once the user is "in the favour system":
	// they have remaining balance, or they've made at least one request.
	const showFavourFab = $derived(
		data.favourCount.remaining > 0 || data.favourRequests.length > 0
	);

	// --- determine which view to show ---

	let tick = $state(Date.now());

	const pastGraphs = $derived(data.graphs.filter((g) => g.starting_at.getTime() <= tick));
	const hasUnlocked = $derived(pastGraphs.length > 0);

	let notificationsSupported = $state(isPushSupported());
	let notificationsSubscribed = $state(true)

	$effect(() => {
		isPushSubscribed().then((subscribed) => {
			notificationsSubscribed = subscribed;
		});
	})

	const nextGraph = $derived(
		data.graphs
			.filter((g) => g.starting_at.getTime() > tick)
			.sort((a, b) => a.starting_at.getTime() - b.starting_at.getTime())[0] ?? null,
	);

	$effect(() => {
		if (!hasUnlocked && nextGraph) {
			const id = setInterval(() => { tick = Date.now(); }, 1000);
			return () => clearInterval(id);
		}
	});

	async function subscribeToNotifications() {
		if (!isPushSupported()) return;

		const permission = await Notification.requestPermission();
		if (permission !== 'granted') {
			notificationsSupported = isPushSupported();
			return;
		}

		await subscribe();
		notificationsSubscribed = true;
	}
</script>

<svelte:head>
	<title>Eros</title>
	<meta content="width=device-width, initial-scale=1" name="viewport" />
</svelte:head>

<div class="mx-auto min-h-dvh max-w-md px-4 py-6 flex flex-col justify-between gap-8">
	<BrandHeader subtitle={hasUnlocked ? undefined : "Something's coming"}>
		{#snippet rightContent()}
			{#if notificationsSupported && !notificationsSubscribed && online}
				<button class="btn btn-circle btn-secondary" onclick={subscribeToNotifications}>
					<BellOff class="h-4 w-4" />
				</button>
			{/if}
		{/snippet}
	</BrandHeader>

	<div class="flex-1 flex flex-col justify-center gap-6" style="view-transition-name: card-content">
		{#if hasUnlocked}
			<!-- Calendar view: at least one graph has been unlocked -->
		<Card>
			<Calendar graphs={data.graphs} onSelectDate={(_date, graphs) => {
				if (graphs.length > 0) goto(`/graphs/${graphs[0].id}`);
			}} />
		</Card>
		{:else if data.graphs.length === 0}
			<!-- No graphs at all -->
			<Card>
				<p class="text-center text-sm opacity-70 animate-popIn">
					Nothing scheduled yet. Check back soon.
				</p>
			</Card>
		{:else if nextGraph}
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
					<Countdown target={nextGraph.starting_at} label="Starts in" />
				</Card>
			</div>
		{/if}
	</div>

	<div class="text-center text-xs opacity-40">Eros</div>
</div>

{#if showFavourFab}
	<FavourFab count={data.favourCount} requests={data.favourRequests} />
{/if}
