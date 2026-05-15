<script lang="ts">
	import type { RewardNode } from '$lib/types/graph';
	import { RewardType } from '$lib/types/graph';
	import {
		Gift,
		Heart,
		Image,
		FileVideo,
		Calendar,
		CreditCard,
		FileText,
		File,
		Download,
		ExternalLink,
	} from 'lucide-svelte';

	import { isOnline } from '$lib/online.svelte';

	const { node }: { node: RewardNode } = $props();

	const online = $derived(isOnline());
	const rewardType = $derived(node.data?.reward_type);
	const file = $derived(node.data?.file);

	const FILE_BACKED_TYPES = new Set([
		RewardType.IMAGE,
		RewardType.VIDEO,
		RewardType.FILE,
		RewardType.CALENDAR,
		RewardType.WALLET,
	]);

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	const ICONS: Record<RewardType, any> = {
		[RewardType.IMAGE]: Image,
		[RewardType.VIDEO]: FileVideo,
		[RewardType.CALENDAR]: Calendar,
		[RewardType.WALLET]: CreditCard,
		[RewardType.FAVOUR]: Heart,
		[RewardType.MARKDOWN]: FileText,
		[RewardType.FILE]: File,
		[RewardType.URL]: ExternalLink,
	};

	const isFileBacked = $derived(rewardType != null && FILE_BACKED_TYPES.has(rewardType));
	const urlHostname = $derived.by(() => {
		const payload = node.data?.payload;
		if (!payload) return '';
		try { return new URL(payload).hostname; } catch { return payload; }
	});
	const Icon = $derived(rewardType ? ICONS[rewardType] : Gift);

	function formatSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
	}
</script>

<div class="flex flex-col items-center gap-6 text-center animate-rewardReveal">
	<!-- Icon -->
	<div
		class="w-20 h-20 rounded-full bg-primary/15 flex items-center justify-center shadow-lg shadow-primary/30"
	>
		<Icon size={32} class="text-primary" />
	</div>

	<!-- Title block -->
	<div class="flex flex-col gap-2">
		<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">
			Reward Unlocked
		</p>
		<h1 class="text-2xl font-extrabold">{node.title}</h1>
		{#if node.description}
			<p class="text-sm opacity-70 leading-relaxed">{node.description}</p>
		{/if}
	</div>

	<!-- Favour badge -->
	{#if node.data?.give_favours && node.data.give_favours > 0}
		<div class="badge badge-primary rounded-2xl px-4 py-3 gap-1.5 text-xs font-semibold">
			<Heart size={12} />
			+{node.data.give_favours} {node.data.give_favours === 1 ? 'favour' : 'favours'}
		</div>
	{/if}

	<!-- Reward content -->
	{#if isFileBacked && file && !online}
		<!-- Offline: file-backed rewards can't load -->
		<div class="w-full bg-base-200 rounded-2xl px-5 py-4 flex flex-col items-center gap-2">
			<Icon size={24} class="opacity-40" />
			<p class="text-sm opacity-50 text-center">
				Available when you're back online.
			</p>
			{#if file.filename}
				<p class="text-xs opacity-30">{file.filename}</p>
			{/if}
		</div>

	{:else if rewardType === RewardType.IMAGE && file}
		<div class="w-full rounded-2xl overflow-hidden shadow-md shadow-primary/15">
			<img
				src={file.url}
				alt={node.title}
				class="w-full h-auto object-contain max-h-[60vh]"
				loading="eager"
			/>
		</div>

	{:else if rewardType === RewardType.VIDEO && file}
		<div class="w-full rounded-2xl overflow-hidden shadow-md shadow-primary/15">
			<!-- svelte-ignore a11y_media_has_caption -->
			<video
				src={file.url}
				controls
				playsinline
				preload="metadata"
				class="w-full max-h-[60vh]"
			>
				Your browser does not support video playback.
			</video>
		</div>

	{:else if rewardType === RewardType.CALENDAR && file}
		<a
			href={file.url}
			download={file.filename}
			class="btn btn-primary rounded-2xl gap-2 w-full max-w-xs"
		>
			<Calendar size={18} />
			Add to Calendar
		</a>
		<p class="text-xs opacity-40">{file.filename}</p>

	{:else if rewardType === RewardType.WALLET && file}
		<a
			href={file.url}
			download={file.filename}
			class="btn btn-primary rounded-2xl gap-2 w-full max-w-xs"
		>
			<CreditCard size={18} />
			Add to Wallet
		</a>
		<p class="text-xs opacity-40">{file.filename}</p>

	{:else if rewardType === RewardType.FILE && file}
		<a
			href={file.url}
			download={file.filename}
			class="btn btn-primary rounded-2xl gap-2 w-full max-w-xs"
		>
			<Download size={18} />
			Download
		</a>
		<p class="text-xs opacity-40">{file.filename} ({formatSize(file.size_bytes)})</p>

	{:else if rewardType === RewardType.MARKDOWN && node.data?.payload}
		<div class="w-full bg-base-200 rounded-2xl px-5 py-4">
			<p class="text-sm opacity-80 leading-relaxed text-left whitespace-pre-wrap">
				{node.data.payload}
			</p>
		</div>

	{:else if rewardType === RewardType.URL && node.data?.payload}
		<a
			href={node.data.payload}
			target="_blank"
			rel="noopener noreferrer"
			class="btn btn-primary rounded-2xl gap-2 w-full max-w-xs"
		>
			<ExternalLink size={18} />
			Open Link
		</a>
		<p class="text-xs opacity-40">{urlHostname}</p>

	{:else if rewardType === RewardType.FAVOUR}
		<!-- Favour-only reward: the badge above is the primary content -->
		{#if !node.data?.give_favours}
			<p class="text-sm opacity-50">This favour is waiting for you.</p>
		{/if}

	{:else if isFileBacked && !file}
		<!-- File-backed type but no file attached -->
		<div class="w-full bg-base-200 rounded-2xl px-5 py-4">
			<p class="text-sm opacity-50 text-center">This reward is not available yet.</p>
		</div>
	{/if}
</div>
