<script lang="ts">
	import type {Snippet} from "svelte";
	import { browser } from '$app/environment';

	const {
		title = 'Eros',
		subtitle,
		rightContent,
	}: {
		title?: string;
		subtitle?: string;
		rightContent?: Snippet;
	} = $props();

	let online = $state(browser ? navigator.onLine : true);

	$effect(() => {
		const setOnline = () => { online = true; };
		const setOffline = () => { online = false; };
		window.addEventListener('online', setOnline);
		window.addEventListener('offline', setOffline);
		return () => {
			window.removeEventListener('online', setOnline);
			window.removeEventListener('offline', setOffline);
		};
	});
</script>

<div class="flex items-center justify-between">
	<div class="flex items-center gap-3">
		<img
			src="/favicon.svg"
			alt="Eros"
			class="h-11 w-11 rounded-2xl shadow-lg shadow-pink-200"
			aria-hidden="true"
		/>

		<div>
			<div class="text-lg font-extrabold leading-tight">{title}</div>
			<div class="flex items-center gap-1.5">
				<div class="text-xs opacity-70">{subtitle}</div>
				{#if !online}
					<div
						class="flex items-center gap-1 rounded-full bg-base-200 px-1.5 py-0.5 animate-popIn"
						title="You are offline"
					>
						<span class="inline-block h-1.5 w-1.5 rounded-full bg-warning"></span>
						<span class="text-xs font-semibold opacity-60">offline</span>
					</div>
				{/if}
			</div>
		</div>
	</div>

	<div class="flex shrink-0 justify-end">
		{@render rightContent?.()}
	</div>
</div>
