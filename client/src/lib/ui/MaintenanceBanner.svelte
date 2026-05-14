<script lang="ts">
	import { loadToken } from '$lib/api/auth';
	import { isOnline, isInitialPingDone } from '$lib/online.svelte';
	import {FileExclamationPoint, MessageCircleWarningIcon, MessageSquareWarning} from "lucide-svelte";

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
	<div class="fixed inset-0 z-100 flex items-center justify-center p-6 bg-base-200">
		<div
			class="flex max-w-sm flex-col items-center gap-6 rounded-3xl bg-base-100 p-8 text-center shadow-lg shadow-pink-200/40"
		>
			<img
				src="/favicon.svg"
				alt="Eros"
				class="h-16 w-16 rounded-2xl shadow-lg shadow-pink-200"
				aria-hidden="true"
			/>

			<div>
				<h1 class="text-xl font-extrabold">We'll be back soon</h1>
				<p class="mt-2 text-sm opacity-70">
					Servers are currently undergoing maintenance.
				</p>
			</div>

			{#if hasToken}
				<button
					class="btn btn-base btn-sm rounded-2xl opacity-70"
					onclick={() => (dismissed = true)}
				>
					Continue with cached data
				</button>
			{/if}
		</div>
	</div>
{/if}
