<script lang="ts">
	import type {Snippet} from "svelte";
	import type {ClassValue} from "svelte/elements";

	const { padded = true, children, ...props}: { padded?: boolean, children: Snippet, class: ClassValue } = $props();
</script>

<div class={["card rounded-3xl bg-base-100 shadow-xl shadow-pink-200/40", props.class]}>
	<div class={padded ? 'card-body gap-4' : ''}>
		{@render children()}
	</div>
</div>
