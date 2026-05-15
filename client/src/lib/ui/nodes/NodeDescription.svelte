<script lang="ts">
	const { description }: { description?: string | null } = $props();
</script>

{#if description}
	<div class="prose prose-sm text-sm opacity-70 leading-relaxed">
		{@html description}
	</div>
{/if}
