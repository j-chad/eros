<script lang="ts">
	import { Heart, ChevronLeft, Check } from 'lucide-svelte';
	import Card from '$lib/ui/base/Card.svelte';
	import { goto } from '$app/navigation';
	import { isOnline } from '$lib/online.svelte';
	import type { FavourChoice, FavourCount, FavourRequest } from '$lib/types/favour';

	const { data }: {
		data: {
			choices: FavourChoice[];
			count: FavourCount;
			requests: FavourRequest[];
		};
	} = $props();

	const isOnline = $derived(isOnline());

	// Local mutable state - initialised from load data, mutated on actions.
	let count = $state<FavourCount>(data.count);
	let requests = $state<FavourRequest[]>(data.requests);

	// Re-sync when navigating back with fresh load data.
	$effect(() => {
		count = data.count;
		requests = data.requests;
	});

	// Choices sorted by cost (ascending) then label (alphabetical)
	const sortedChoices = $derived.by(() =>
		[...data.choices].sort((a, b) => a.cost - b.cost || a.label.localeCompare(b.label))
	);

	// Tab state - only show tabs if there are requests
	const hasRequests = $derived(requests.length > 0);
	let activeTab = $state<'spend' | 'requests'>('spend');

	// Expand state for the inline request flow
	let expandedChoiceId = $state<string | null>(null);
	let message = $state('');
	let submitting = $state(false);
	let error = $state<string | null>(null);
	let successChoiceId = $state<string | null>(null);

	// Requests split into groups, newest first within each
	const pendingRequests = $derived.by(() =>
		requests
			.filter((r) => !r.fulfilled_at)
			.sort((a, b) => b.requested_at.localeCompare(a.requested_at))
	);
	const fulfilledRequests = $derived.by(() =>
		requests
			.filter((r) => r.fulfilled_at)
			.sort((a, b) => b.fulfilled_at!.localeCompare(a.fulfilled_at!))
	);

	// Collapsible sections - collapsed by default
	let pendingOpen = $state(true);
	let fulfilledOpen = $state(false);

	function handleExpand(choiceId: string) {
		if (expandedChoiceId === choiceId) {
			expandedChoiceId = null;
		} else {
			expandedChoiceId = choiceId;
			message = '';
			error = null;
		}
	}

	async function handleConfirm(choice: FavourChoice) {
		submitting = true;
		error = null;

		try {
			const { requestFavour } = await import('$lib/services/favour');
			const newRequest = await requestFavour(choice.id, message || undefined);

			// Update local state
			requests = [...requests, newRequest];
			count = { ...count, remaining: count.remaining - choice.cost };

			// Show success animation
			expandedChoiceId = null;
			successChoiceId = choice.id;
			setTimeout(() => {
				successChoiceId = null;
			}, 1500);

			// Refresh the count from the server to stay in sync
			const { getCount } = await import('$lib/services/favour');
			count = await getCount();
		} catch (e) {
			if (e && typeof e === 'object' && 'status' in e) {
				// Refresh count - balance may have changed
				try {
					const { getCount } = await import('$lib/services/favour');
					count = await getCount();
				} catch { /* keep going with stale count */ }
			}
			error = 'Not enough favour points. Your balance has been updated.';
		} finally {
			submitting = false;
		}
	}

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString(undefined, {
			day: 'numeric',
			month: 'short',
			year: 'numeric',
		});
	}
</script>

<svelte:head>
	<title>Favours - Eros</title>
</svelte:head>

<div class="mx-auto min-h-dvh max-w-md px-4 py-6 flex flex-col gap-6">
	<!-- Header -->
	<div class="flex items-center gap-3">
		<button
			class="btn btn-ghost btn-sm btn-circle"
			onclick={() => goto('/')}
			aria-label="Back"
		>
			<ChevronLeft size={20} />
		</button>
		<h1 class="text-xl font-extrabold">Favours</h1>
	</div>

	<!-- Balance -->
	<div class="flex items-center justify-center gap-2 py-2">
		<Heart size={20} class="text-primary" />
		<span class="text-2xl font-extrabold">{count.remaining}</span>
		<span class="text-sm opacity-70">{count.remaining === 1 ? 'favour' : 'favours'} remaining</span>
	</div>

	<!-- Tabs (only shown when requests exist) -->
	{#if hasRequests}
		<div role="tablist" class="tabs tabs-bordered">
			<button
				role="tab"
				class="tab"
				class:tab-active={activeTab === 'spend'}
				onclick={() => { activeTab = 'spend'; }}
			>
				Spend
			</button>
			<button
				role="tab"
				class="tab"
				class:tab-active={activeTab === 'requests'}
				onclick={() => { activeTab = 'requests'; }}
			>
				Requests
			</button>
		</div>
	{/if}

	<!-- Spend tab -->
	{#if activeTab === 'spend'}
		<div class="flex flex-col gap-4 animate-popIn">
			<p class="text-sm opacity-60 text-center">
				Pick from the menu, or propose your own - the worst that can happen is a no.
			</p>

			{#if sortedChoices.length === 0}
				<Card>
					<p class="text-center text-sm opacity-70">
						No favours available yet. Check back soon.
					</p>
				</Card>
			{:else}
				{#each sortedChoices as choice (choice.id)}
					<Card>
						<div class="flex items-start justify-between gap-3">
							<div class="flex-1">
								<h3 class="font-bold">{choice.label}</h3>
								{#if choice.description}
									<p class="text-sm opacity-70 mt-1">{choice.description}</p>
								{/if}
							</div>
							<div class="badge badge-primary rounded-2xl px-3 py-2 text-xs font-semibold shrink-0">
								{choice.cost} {choice.cost === 1 ? 'pt' : 'pts'}
							</div>
						</div>

						{#if successChoiceId === choice.id}
							<!-- Success state -->
							<div class="flex items-center justify-center gap-2 py-3 animate-popIn">
								<div class="w-8 h-8 rounded-full bg-success/15 flex items-center justify-center">
									<Check size={16} class="text-success" />
								</div>
								<span class="text-sm font-semibold text-success">Request sent</span>
							</div>
						{:else if expandedChoiceId === choice.id}
							<!-- Expanded request form -->
							<div class="flex flex-col gap-3 mt-3 animate-popIn">
								{#if choice.can_message}
									<textarea
										class="textarea textarea-bordered rounded-2xl w-full text-sm"
										placeholder="Add a message (optional)"
										rows={2}
										bind:value={message}
										disabled={submitting}
									></textarea>
								{:else}
									<p class="text-sm opacity-70">
										Spend {choice.cost} {choice.cost === 1 ? 'point' : 'points'} on {choice.label}?
									</p>
								{/if}

								{#if error}
									<div class="alert alert-error rounded-2xl text-sm">
										{error}
									</div>
								{/if}

								<div class="flex gap-2">
									<button
										class="btn btn-ghost btn-sm rounded-2xl flex-1"
										onclick={() => { expandedChoiceId = null; }}
										disabled={submitting}
									>
										Cancel
									</button>
									<button
										class="btn btn-primary btn-sm rounded-2xl flex-1"
										onclick={() => handleConfirm(choice)}
										disabled={submitting}
									>
										{#if submitting}
											<span class="loading loading-spinner loading-sm"></span>
										{:else}
											Confirm
										{/if}
									</button>
								</div>
							</div>
						{:else}
							<!-- Request button -->
							<button
								class="btn btn-primary btn-sm rounded-2xl w-full mt-3"
								disabled={count.remaining < choice.cost || !isOnline}
								onclick={() => handleExpand(choice.id)}
							>
								Request
							</button>
							{#if !isOnline}
								<p class="text-xs opacity-50 text-center mt-1">You're offline</p>
							{/if}
						{/if}
					</Card>
				{/each}
			{/if}
		</div>
	{/if}

	<!-- Requests tab -->
	{#if activeTab === 'requests' && hasRequests}
		<div class="flex flex-col gap-6 animate-popIn">
			<!-- Pending requests -->
			{#if pendingRequests.length > 0}
				<div>
					<button
						class="flex items-center justify-between w-full py-2"
						onclick={() => { pendingOpen = !pendingOpen; }}
					>
						<span class="text-sm font-semibold">Pending ({pendingRequests.length})</span>
						<ChevronLeft size={16} class="transition-transform {pendingOpen ? '-rotate-90' : 'rotate-180'}" />
					</button>
					{#if pendingOpen}
						<div class="flex flex-col gap-3 mt-2 animate-popIn">
							{#each pendingRequests as req (req.id)}
								<Card>
									<div class="flex items-start justify-between gap-3">
										<div class="flex-1">
											<h3 class="font-bold">{req.choice_label}</h3>
											{#if req.choice_description}
												<p class="text-sm opacity-70 mt-1">{req.choice_description}</p>
											{/if}
										</div>
										<span class="badge badge-outline rounded-2xl text-xs shrink-0">Pending</span>
									</div>

									{#if req.message}
										<div class="bg-base-200 rounded-2xl px-4 py-3 mt-2">
											<p class="text-sm opacity-80 whitespace-pre-wrap">{req.message}</p>
										</div>
									{/if}

									<p class="text-xs opacity-50 mt-2">
										Requested on {formatDate(req.requested_at)}
									</p>
								</Card>
							{/each}
						</div>
					{/if}
				</div>
			{/if}

			<!-- Fulfilled requests -->
			{#if fulfilledRequests.length > 0}
				<div>
					<button
						class="flex items-center justify-between w-full py-2"
						onclick={() => { fulfilledOpen = !fulfilledOpen; }}
					>
						<span class="text-sm font-semibold">Fulfilled ({fulfilledRequests.length})</span>
						<ChevronLeft size={16} class="transition-transform {fulfilledOpen ? '-rotate-90' : 'rotate-180'}" />
					</button>
					{#if fulfilledOpen}
						<div class="flex flex-col gap-3 mt-2 animate-popIn">
							{#each fulfilledRequests as req (req.id)}
								<Card class="opacity-60">
									<div class="flex items-start justify-between gap-3">
										<div class="flex-1">
											<h3 class="font-bold">{req.choice_label}</h3>
											{#if req.choice_description}
												<p class="text-sm opacity-70 mt-1">{req.choice_description}</p>
											{/if}
										</div>
										<div class="flex items-center gap-1 text-success shrink-0">
											<Check size={14} />
											<span class="text-xs font-semibold">Fulfilled</span>
										</div>
									</div>

									{#if req.message}
										<div class="bg-base-200 rounded-2xl px-4 py-3 mt-2">
											<p class="text-sm opacity-80 whitespace-pre-wrap">{req.message}</p>
										</div>
									{/if}

									<p class="text-xs opacity-50 mt-2">
										Fulfilled on {formatDate(req.fulfilled_at!)}
									</p>
								</Card>
							{/each}
						</div>
					{/if}
				</div>
			{/if}
		</div>
	{/if}
</div>
