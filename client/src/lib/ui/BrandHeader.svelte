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
		<div
			class="grid h-11 w-11 place-items-center rounded-2xl bg-linear-to-br from-pink-400 to-fuchsia-500 text-white shadow-lg shadow-pink-200"
			aria-hidden="true"
		>
			<span class="text-lg font-black">E</span>
		</div>

		<div>
			<div class="text-lg font-extrabold leading-tight">{title}</div>
			<div class="text-xs opacity-70">{subtitle}</div>
		</div>
	</div>

	<div class="min-w-40 flex justify-end">
		{@render rightContent?.()}
	</div>
</div>
