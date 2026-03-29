<script lang="ts">
	import { navigating } from '$app/state';

	let { children } = $props();
</script>

{#if navigating.to}
	<!-- Glowing progress bar -->
	<div class="fixed inset-x-0 top-0 z-50 h-1 overflow-hidden" aria-hidden="true">
		<div class="h-full bg-primary animate-navProgress" style="box-shadow: 0 0 8px 2px var(--color-primary);"></div>
	</div>
	<!-- Subtle page dim so the user knows the current page is "leaving" -->
	<div class="fixed inset-0 z-40 bg-base-100/30 animate-fadeIn pointer-events-none" aria-hidden="true"></div>
{/if}

{@render children()}
