<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import BrandHeader from '$lib/ui/BrandHeader.svelte';
	import Card from '$lib/ui/base/Card.svelte';
	import Button from '$lib/ui/base/Button.svelte';

	const status = $derived(page.status);
	let showRawError = $state(false);

	const content = $derived(() => {
		if (status === 404) {
			return {
				title: 'Page not found',
				description: "Sorry, we couldn't find the page you're looking for."
			};
		}
		if (status === 403) {
			return {
				title: 'Access denied',
				description: "You don't have permission to access this resource."
			};
		}
		if (status >= 500) {
			return {
				title: 'Server error',
				description: 'Something went wrong on our end. Please try again.'
			};
		}
		return {
			title: 'Unexpected error',
			description: page.error?.message ?? 'An unexpected error occurred.'
		};
	});
</script>

<svelte:head>
	<title>{status} – {content().title} • Eros</title>
</svelte:head>

<div class="min-h-dvh bg-linear-to-br from-pink-50 via-base-100 to-pink-100">
	<div class="mx-auto min-h-dvh max-w-md px-4 py-6 flex flex-col">
		<div class="flex flex-1 items-center justify-center py-12">
			<div class="w-full animate-popIn">
				<Card class="bg-secondary">
					<div class="card-body gap-6 text-center">
						<p class="text-5xl font-black text-primary/30">{status}</p>

						<div class="space-y-2">
							<h1 class="text-xl font-bold">{content().title}</h1>
							{#if showRawError && page.error}
								<pre class="mt-2 text-left text-xs bg-base-200 rounded-2xl p-4 overflow-x-auto text-base-content/70">{JSON.stringify(page.error, null, 2)}</pre>
							{:else}
								<p class="text-sm opacity-70">{content().description}</p>
							{/if}
						</div>

						<div class="flex flex-col gap-3 pt-1">
							<Button block variant="primary" onclick={() => goto('/')}>Go Home</Button>
							{#if page.error}
								<Button ghost onclick={() => { showRawError = !showRawError; }}>
									{showRawError ? 'Hide details' : 'Details'}
								</Button>
							{/if}
						</div>
					</div>
				</Card>
			</div>
		</div>
	</div>
</div>
