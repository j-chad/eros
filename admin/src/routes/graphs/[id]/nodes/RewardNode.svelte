<script lang="ts">
	import { Gift, Eye, EyeOff, Coins } from 'lucide-svelte';
	import type { RewardNode } from '$lib/types';
	import BaseNode from './BaseNode.svelte';
	import type { NodeProps } from './types';

	let { data }: NodeProps<RewardNode> = $props();
	let node = $derived(data.node);

	let showPayload = $state(false);

	const rewardType = $derived(node.data?.reward_type ?? 'content');
	const payload = $derived(node.data?.payload ?? '');
	const giveFavours = $derived(node.data?.give_favours ?? 0);

	const hasPayload = $derived(!!payload?.trim());
</script>

<BaseNode
	{node}
	onEdit={data.onEdit}
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
						<span class="dot" aria-hidden="true"></span>
						<span class="type-text">{rewardType}</span>
					</div>

					<button
						class="toggle"
						disabled={!hasPayload}
						title={showPayload ? 'Hide payload' : 'Reveal payload'}
						onclick={() => (showPayload = !showPayload)}
					>
						{#if showPayload}
							<EyeOff size={14} />
							<span>Hide</span>
						{:else}
							<Eye size={14} />
							<span>Reveal</span>
						{/if}
					</button>
				</div>

				<div class="payload">
					{#if !hasPayload}
						<div class="payload-empty">
							<span class="muted">No payload set</span>
							<span class="hint">Edit node to add content</span>
						</div>
					{:else if showPayload}
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

				{#if giveFavours > 0}
					<div class="favours">
						<Coins size={14} />
						<span class="favours-label">Favours:</span>
						<span class="favours-value">+{giveFavours}</span>
					</div>
				{/if}
			</div>

			<div class="meta">
				<div class="meta-item">
					<span class="key">🎁 Type</span>
					<span class="value">{rewardType}</span>
				</div>
				<div class="meta-item">
					<span class="key">📦 Payload</span>
					<span class="value">{hasPayload ? `${payload.length} chars` : '—'}</span>
				</div>
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
		padding: 0.375rem 0.5rem;
		border-radius: 999px;
		background: rgba(6, 182, 212, 0.12);
		border: 1px solid rgba(6, 182, 212, 0.25);
		color: #0e7490;
		font-size: 0.75rem;
		font-weight: 700;
		letter-spacing: 0.01em;
		text-transform: lowercase;
	}

	.dot {
		width: 8px;
		height: 8px;
		border-radius: 999px;
		background: linear-gradient(135deg, #06b6d4 0%, #0891b2 100%);
		box-shadow: 0 0 0 3px rgba(6, 182, 212, 0.18);
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

	.toggle:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 8px 14px rgba(0, 0, 0, 0.08);
	}

	.toggle:active:not(:disabled) {
		transform: translateY(0px) scale(0.98);
	}

	.toggle:disabled {
		opacity: 0.5;
		cursor: not-allowed;
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
	}

	.payload-empty {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
		align-items: center;
		text-align: center;
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

	.favours {
		margin-top: 0.6rem;
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		padding: 0.4rem 0.55rem;
		border-radius: 10px;
		background: rgba(16, 185, 129, 0.12);
		border: 1px solid rgba(16, 185, 129, 0.25);
		color: #065f46;
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
	}

	.meta {
		display: flex;
		gap: 0.75rem;
	}

	.meta-item {
		flex: 1;
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
	}
</style>
