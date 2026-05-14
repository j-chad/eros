<script lang="ts">
	import type {Snippet} from "svelte";
	import { isOnline } from '$lib/online.svelte';

	const {
		title = 'Eros',
		subtitle,
		rightContent,
		compact = false,
	}: {
		title?: string;
		subtitle?: string;
		rightContent?: Snippet;
		compact?: boolean;
	} = $props();

	const online = $derived(isOnline());
</script>

{#if compact}
	<img
		src="/favicon.svg"
		alt="Eros"
		class="h-9 w-9 rounded-xl shadow shadow-pink-200"
		aria-hidden="true"
	/>
{:else}
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
					{#if subtitle}
						<div class="text-xs opacity-70">{subtitle}</div>
					{/if}
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
{/if}
