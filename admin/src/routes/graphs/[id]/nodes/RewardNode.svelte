<script lang="ts">
	import {Gift, Eye, EyeOff, Star, Image, Video, Calendar, Ticket, Files, ExternalLink} from 'lucide-svelte';
	import {type RewardNode, RewardType} from '$lib/types';
	import BaseNode from './BaseNode.svelte';
	import type { NodeProps } from './types';

	let { data }: NodeProps<RewardNode> = $props();
	let node = $derived(data.node);

	let showPayload = $state(false);

	const rewardType = $derived(node.data?.reward_type ?? RewardType.FAVOUR);
	const payload = $derived(node.data?.payload ?? '');
	const giveFavours = $derived(node.data?.give_favours ?? 0);

	const hasPayload = $derived(!!payload?.trim() && rewardType !== 'favour');

	// Get icon for reward type
	const RewardTypeIcon = $derived.by(() => {
		switch (rewardType) {
			case RewardType.IMAGE: return Image;
			case RewardType.VIDEO: return Video;
			case RewardType.CALENDAR: return Calendar;
			case RewardType.WALLET: return Ticket;
			case RewardType.FAVOUR: return Star;
			case RewardType.FILE: return Files;
			case RewardType.URL: return ExternalLink;
			default: return Gift;
		}
	});

	// Format reward type display name
	const rewardTypeLabel = $derived.by(() => {
		switch (rewardType) {
			case RewardType.IMAGE: return 'Image';
			case RewardType.VIDEO: return 'Video';
			case RewardType.CALENDAR: return 'Calendar Event';
			case RewardType.WALLET: return 'Coupon';
			case RewardType.FAVOUR: return 'Favour Points';
			case RewardType.FILE: return 'File';
			case RewardType.URL: return 'URL';
			default: return rewardType;
		}
	});

	// Parse and format payload preview
	const payloadPreview = $derived.by(() => {
		if (!payload) return '—';
		if (rewardType === RewardType.URL) {
			try { return new URL(payload).hostname; } catch { return payload; }
		}
		try {
			const parsed = JSON.parse(payload);
			switch (rewardType) {
				case RewardType.IMAGE:
				case RewardType.VIDEO:
					return parsed.url ? new URL(parsed.url).hostname : '—';
				case RewardType.CALENDAR:
					return parsed.event_title || '—';
				case RewardType.WALLET:
					return parsed.url ? 'Wallet Pass' : '—';
				case RewardType.FILE:
					return parsed.filename || '—';
				default:
					return `${payload.length} chars`;
			}
		} catch {
			return `${payload.length} chars`;
		}
	});
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
	showProgress={data.showProgress}
	onToggleUnlock={data.onToggleUnlock}
	config={{
		color: '#06b6d4',
		gradient: 'linear-gradient(135deg, #06b6d4 0%, #0891b2 100%)',
		icon: Gift,
		label: 'Reward'
	}}
>
	{#snippet children()}
		<div class="reward-content">
			<div class="reward-card">
				<div class="top-row">
					<div class="type-badge" title="Reward type">
						<RewardTypeIcon size={12} />
						<span class="type-text">{rewardTypeLabel}</span>
					</div>

					{#if hasPayload}
						<button
							class="toggle"
							title={showPayload ? 'Hide payload' : 'Reveal payload'}
							onclick={() => (showPayload = !showPayload)}
						>
							{#if showPayload}
								<EyeOff size={14} />
								<span>Hide</span>
							{:else}
								<Eye size={14} />
								<span>Show</span>
							{/if}
						</button>
					{/if}
				</div>

				{#if hasPayload}
					<div class="payload">
						{#if showPayload}
							<pre class="payload-pre">{payload}</pre>
						{:else}
							<div class="payload-hidden">
								<div class="dots" aria-hidden="true">
									<span></span><span></span><span></span>
								</div>
								<span class="muted">Hidden</span>
							</div>
						{/if}
					</div>
				{/if}

				{#if giveFavours > 0}
					<div class="favours">
						<Star size={14} fill="currentColor" />
						<span class="favours-value">+{giveFavours}</span>
						<span class="favours-label">Favour{giveFavours !== 1 ? 's' : ''}</span>
					</div>
				{/if}
			</div>

			<div class="meta">
				<div class="meta-item">
					<span class="key">🎁 Type</span>
					<span class="value">{rewardTypeLabel}</span>
				</div>
				<div class="meta-item">
					<span class="key">📦 Content</span>
					<span class="value">{payloadPreview}</span>
				</div>
				{#if giveFavours > 0}
					<div class="meta-item">
						<span class="key">⭐ Favours</span>
						<span class="value">+{giveFavours}</span>
					</div>
				{/if}
			</div>
		</div>
	{/snippet}
</BaseNode>

<style>
	.reward-content {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.reward-card {
		border-radius: 12px;
		border: 1px solid #e5e7eb;
		background: linear-gradient(180deg, #ffffff 0%, #f9fafb 100%);
		padding: 0.75rem;
		box-shadow: 0 4px 10px rgba(0, 0, 0, 0.06);
	}

	.top-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		margin-bottom: 0.5rem;
	}

	.type-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.625rem;
		border-radius: 999px;
		background: rgba(6, 182, 212, 0.12);
		border: 1px solid rgba(6, 182, 212, 0.25);
		color: #0e7490;
		font-size: 0.75rem;
		font-weight: 700;
		letter-spacing: 0.01em;
	}

	.type-text {
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		max-width: 140px;
	}

	.toggle {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		border: 1px solid #e5e7eb;
		background: white;
		border-radius: 999px;
		padding: 0.35rem 0.5rem;
		font-size: 0.75rem;
		font-weight: 700;
		color: #374151;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.toggle:hover {
		transform: translateY(-1px);
		box-shadow: 0 8px 14px rgba(0, 0, 0, 0.08);
	}

	.toggle:active {
		transform: translateY(0px) scale(0.98);
	}

	.payload {
		border-radius: 10px;
		border: 1px dashed #e5e7eb;
		background: #ffffff;
		padding: 0.65rem;
		min-height: 58px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.payload-pre {
		width: 100%;
		margin: 0;
		font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
		monospace;
		font-size: 0.72rem;
		font-weight: 600;
		color: #111827;
		white-space: pre-wrap;
		word-break: break-word;
		max-height: 120px;
		overflow-y: auto;
	}

	.payload-hidden {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.dots {
		display: inline-flex;
		gap: 0.25rem;
	}

	.dots span {
		width: 6px;
		height: 6px;
		border-radius: 999px;
		background: #9ca3af;
		animation: bounce 0.9s infinite ease-in-out;
	}

	.dots span:nth-child(2) {
		animation-delay: 0.12s;
	}
	.dots span:nth-child(3) {
		animation-delay: 0.24s;
	}

	@keyframes bounce {
		0%,
		80%,
		100% {
			transform: translateY(0);
			opacity: 0.6;
		}
		40% {
			transform: translateY(-4px);
			opacity: 1;
		}
	}

	.muted {
		color: #9ca3af;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.favours {
		margin-top: 0.6rem;
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		padding: 0.4rem 0.65rem;
		border-radius: 10px;
		background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
		border: 1px solid #fbbf24;
		color: #92400e;
		font-size: 0.75rem;
		font-weight: 800;
		width: fit-content;
	}

	.favours-label {
		font-weight: 700;
		opacity: 0.9;
	}

	.favours-value {
		font-variant-numeric: tabular-nums;
		font-weight: 900;
	}

	.meta {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
		gap: 0.75rem;
	}

	.meta-item {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
		font-size: 0.75rem;
	}

	.key {
		color: #6b7280;
		font-weight: 600;
	}

	.value {
		color: #1f2937;
		font-weight: 700;
		font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
		monospace;
		font-size: 0.6875rem;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>
