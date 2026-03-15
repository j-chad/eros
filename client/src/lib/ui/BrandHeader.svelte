<script lang="ts">
	import type {Snippet} from "svelte";

	const {
		title = 'Eros',
		subtitle,
		rightContent,
	}: {
		title?: string;
		subtitle?: string;
		rightContent?: Snippet;
	} = $props();
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
			<div class="text-xs opacity-70">{subtitle}</div>
		</div>
	</div>

	<div class="flex shrink-0 justify-end">
		{@render rightContent?.()}
	</div>
</div>
