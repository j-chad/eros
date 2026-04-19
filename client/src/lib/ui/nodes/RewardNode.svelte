<script lang="ts">
	import type { RewardNode } from '$lib/types/graph';
	import { RewardType } from '$lib/types/graph';
	import { Gift, Heart, Image, FileVideo, Calendar, CreditCard, FileText, File } from 'lucide-svelte';

	const { node }: { node: RewardNode } = $props();

	// lucide-svelte exports Svelte 4-style components; cast to any to avoid
	// fighting their internal types with Svelte 5's Component<...> signature.
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	const Icon: any = $derived(
		node.data?.reward_type === RewardType.IMAGE ? Image :
		node.data?.reward_type === RewardType.VIDEO ? FileVideo :
		node.data?.reward_type === RewardType.CALENDAR ? Calendar :
		node.data?.reward_type === RewardType.WALLET ? CreditCard :
		node.data?.reward_type === RewardType.FAVOUR ? Heart :
		node.data?.reward_type === RewardType.MARKDOWN ? FileText :
		node.data?.reward_type === RewardType.FILE ? File :
		Gift
	);

	const rewardLabel = $derived(
		node.data?.reward_type === RewardType.IMAGE ? 'A photo for you' :
		node.data?.reward_type === RewardType.VIDEO ? 'A video for you' :
		node.data?.reward_type === RewardType.CALENDAR ? 'A date to remember' :
		node.data?.reward_type === RewardType.WALLET ? 'Something special' :
		node.data?.reward_type === RewardType.FAVOUR ? 'A favour, just for you' :
		node.data?.reward_type === RewardType.MARKDOWN ? 'A message for you' :
		node.data?.reward_type === RewardType.FILE ? 'Something to keep' :
		'Your reward'
	);
</script>

<div class="flex flex-col items-center gap-6 text-center animate-rewardReveal">
	<div class="w-20 h-20 rounded-full bg-primary/15 flex items-center justify-center shadow-lg shadow-pink-200/60">
		<Icon size={32} class="text-primary" />
	</div>

	<div class="flex flex-col gap-2">
		<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">{rewardLabel}</p>
		<h1 class="text-2xl font-extrabold">{node.title}</h1>
		{#if node.description}
			<p class="text-sm opacity-70 leading-relaxed">{node.description}</p>
		{/if}
	</div>

	{#if node.data?.give_favours && node.data.give_favours > 0}
		<div class="badge badge-primary rounded-2xl px-4 py-3 gap-1.5 text-xs font-semibold">
			<Heart size={12} />
			+{node.data.give_favours} {node.data.give_favours === 1 ? 'favour' : 'favours'}
		</div>
	{/if}

	{#if node.data?.payload}
		<div class="w-full bg-base-200 rounded-2xl px-5 py-4">
			{#if node.data.reward_type === RewardType.MARKDOWN}
				<p class="text-sm opacity-80 leading-relaxed text-left whitespace-pre-wrap">{node.data.payload}</p>
			{:else}
				<p class="text-xs opacity-50 text-center font-mono">{node.data.payload}</p>
			{/if}
		</div>
	{/if}
</div>
