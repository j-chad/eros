<script lang="ts">
	import type {Snippet} from "svelte";

	const {
		variant = 'secondary',
		ghost = false,
		block = false,
		loading = false,
		disabled = false,
		type = 'button',
		onclick,
		children
	}: {
		children: Snippet;
		variant?: 'primary' | 'secondary' | 'accent' | 'neutral';
		ghost?: boolean;
		block?: boolean;
		loading?: boolean;
		disabled?: boolean;
		type?: 'button' | 'submit';
		onclick?: () => void;
	} = $props();

	const classes = $derived([
		'btn',
		ghost ? 'btn-ghost' : `btn-${variant}`,
		'rounded-2xl',
		block && 'w-full'
	]);
</script>

<button class={classes} {type} disabled={disabled || loading} {onclick}>
	{#if loading}
		<span class="loading loading-spinner loading-sm"></span>
	{/if}
	{@render children()}
</button>
