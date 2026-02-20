<script>
	import { page } from '$app/state'
	import { goto } from '$app/navigation'

	const status = $derived(page.status)
	let showRawError = $state(false)

	const content = $derived(() => {
		if (status === 404) {
			return {
				title: 'Page not found',
				description: 'Sorry, we couldn’t find the page you’re looking for.'
			}
		}

		if (status === 403) {
			return {
				title: 'Access denied',
				description: 'You don’t have permission to access this resource.'
			}
		}

		if (status >= 500) {
			return {
				title: 'Server error',
				description: 'Something went wrong on our end. Please try again.'
			}
		}

		return {
			title: 'Unexpected error',
			description: page.error?.message ?? 'An unexpected error occurred.'
		}
	})

	function goHome() {
		goto('/')
	}
</script>

<svelte:head>
	<title>{status} – {content().title}</title>
</svelte:head>

<main class="min-h-dvh grid place-items-center bg-base-100 px-6 py-24 sm:py-32">
	<div class="text-center max-w-xl">

		<p class="text-base font-semibold text-primary">
			{status}
		</p>

		<h1 class="mt-4 text-5xl font-bold tracking-tight sm:text-7xl text-balance">
			{content().title}
		</h1>

		{#if showRawError && page.error}
			<pre class="mt-6 text-left text-sm bg-neutral text-neutral-content rounded-lg p-4 overflow-x-auto">{JSON.stringify(page.error, null, 2)}</pre>
		{:else}
			<p class="mt-6 text-base-content/70 text-lg sm:text-xl/8 font-medium text-pretty">
				{content().description}
			</p>
		{/if}

		<div class="mt-10 flex items-center justify-center gap-x-6">
			<button
				class="btn btn-primary rounded-xl"
				onclick={goHome}
			>
				Go back home
			</button>

			<button
				class="btn btn-ghost rounded-xl"
				onclick={() => {showRawError = !showRawError}}
			>
				{showRawError ? 'Hide Details' : 'Show Details'}
			</button>
		</div>
	</div>
</main>
