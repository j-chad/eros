<script lang="ts">
	import type {HTMLInputAttributes} from "svelte/elements";

	let {
		label,
		hint,
		error,
		value = $bindable<string>(),
		element = $bindable<HTMLInputElement | null>(null),
		...rest
	}: {
		label?: string;
		hint?: string;
		error?: string;
		value?: string;
		element?: HTMLInputElement | null;
	} & HTMLInputAttributes = $props();
</script>

<label class="form-control w-full">
	{#if label}
		<div class="label">
			<span class="label-text font-semibold">{label}</span>
			{#if hint}
				<span class="label-text-alt opacity-60">{hint}</span>
			{/if}
		</div>
	{/if}

	<input
		class="input input-bordered input-secondary w-full rounded-2xl text-lg"
		bind:value
		bind:this={element}
		{...rest}
	/>

	{#if error}
		<div class="label">
			<span class="label-text-alt text-error">{error}</span>
		</div>
	{/if}
</label>
