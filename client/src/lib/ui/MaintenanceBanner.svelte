<script lang="ts">
	import { loadToken } from '$lib/api/auth';
	import { isOnline, isInitialPingDone } from '$lib/online.svelte';
	import { Wrench } from 'lucide-svelte';

	let dismissed = $state(false);
	let hasToken = $state(false);

	const online = $derived(isOnline());
	const pingDone = $derived(isInitialPingDone());

	// Auto-dismiss: if the server is actually up despite the maintenance flag,
	// there's no reason to block the user.
	const autoDismissed = $derived(pingDone && online);

	// Check for a cached auth token so we can decide whether to show the
	// "Continue with cached data" button.
	loadToken().then((token) => {
		hasToken = !!token;
	});
</script>

{#if !dismissed && !autoDismissed}
	<div class="fixed inset-0 z-100 flex flex-col items-center justify-center gap-10 bg-base-200 p-6">
		<div class="flex flex-col items-center gap-8 animate-popIn">
			<div class="flex h-20 w-20 items-center justify-center rounded-full bg-primary/10">
				<Wrench size={36} class="text-primary" strokeWidth={2} />
			</div>

			<div class="flex flex-col items-center gap-3 text-center">
				<h1 class="text-2xl font-extrabold">Taking a break</h1>
				<p class="max-w-xs text-sm leading-relaxed opacity-70">
					The server is offline while some things are being reworked.
					It'll be back when it's ready.
				</p>
			</div>
		</div>

		{#if hasToken}
			<button
				class="text-sm font-medium opacity-50 underline decoration-dotted underline-offset-4 transition-opacity hover:opacity-70"
				onclick={() => (dismissed = true)}
			>
				Browse cached data instead
			</button>
		{/if}
	</div>
{/if}
