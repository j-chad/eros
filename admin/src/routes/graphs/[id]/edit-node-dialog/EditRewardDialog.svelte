<script lang="ts">
	import type { RewardNode } from '$lib/types';
	import { Image, Video, Gift, Star, Sparkles } from 'lucide-svelte';

	let {
		node,
		onSave,
		onCancel
	}: {
		node: RewardNode;
		onSave: (node: RewardNode) => void;
		onCancel: () => void;
	} = $props();

	let editForm = $state({
		title: node.title,
		description: node.description || '',
		reward_type: node.data?.reward_type || 'favour',
		payload: node.data?.payload || ''
	});

	// Parse payload based on type for easier editing
	let payloadForm = $derived.by(() => {
		try {
			const parsed = editForm.payload ? JSON.parse(editForm.payload) : {};
			return {
				url: parsed.url || '',
				amount: parsed.amount || 0,
				coupon_data: parsed.coupon_data || {}
			};
		} catch {
			return {
				url: '',
				amount: 0,
				coupon_data: {}
			};
		}
	});

	function updatePayload(updates: Partial<typeof payloadForm>) {
		const current = payloadForm;
		const updated = { ...current, ...updates };

		// Create clean payload based on type
		let payload: any = {};
		if (editForm.reward_type === 'video' || editForm.reward_type === 'image') {
			payload = { url: updated.url };
		} else if (editForm.reward_type === 'favour') {
			payload = { amount: updated.amount };
		} else if (editForm.reward_type === 'coupon') {
			payload = {
				url: updated.url,
				coupon_data: updated.coupon_data
			};
		}

		editForm.payload = JSON.stringify(payload);
	}

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		event.stopPropagation();

		const updatedNode: RewardNode = {
			...node,
			title: editForm.title,
			description: editForm.description,
			data: {
				reward_type: editForm.reward_type,
				payload: editForm.payload,
				give_favours: 0
			}
		};

		onSave(updatedNode);
	}

	const rewardTypes = [
		{ value: 'favour', label: 'Favour Points', icon: Star, description: 'Award currency to the user' },
		{ value: 'image', label: 'Image', icon: Image, description: 'Show a celebratory image' },
		{ value: 'video', label: 'Video', icon: Video, description: 'Play a reward video' },
		{ value: 'coupon', label: 'Coupon', icon: Gift, description: 'Apple Wallet coupon' }
	];
</script>

<h2>🎁 Edit Reward Node</h2>

<form onsubmit={handleSubmit}>
	<div class="form-layout">
		<!-- Left Column -->
		<div class="form-column">
			<div class="form-group">
				<label for="title">Title</label>
				<input
					id="title"
					type="text"
					bind:value={editForm.title}
					required
					placeholder="e.g., Quest Complete!, Daily Check-in"
				/>
			</div>

			<div class="form-group">
				<label for="description">Description</label>
				<textarea
					id="description"
					bind:value={editForm.description}
					rows="3"
					placeholder="Describe what the user accomplished"
				></textarea>
			</div>

			<div class="form-group">
				<label>Reward Type</label>
				<div class="reward-type-grid">
					{#each rewardTypes as type}
						<button
							type="button"
							class="reward-type-card"
							class:selected={editForm.reward_type === type.value}
							onclick={() => editForm.reward_type = type.value}
						>
							<type.icon size={24} />
							<div class="type-info">
								<div class="type-label">{type.label}</div>
								<div class="type-description">{type.description}</div>
							</div>
						</button>
					{/each}
				</div>
			</div>
		</div>

		<!-- Right Column - Reward Configuration -->
		<div class="form-column config-column">
			<div class="config-header">
				<Sparkles size={18} />
				<span>Reward Configuration</span>
			</div>

			{#if editForm.reward_type === 'favour'}
				<div class="form-group">
					<label for="amount">Favour Amount</label>
					<input
						id="amount"
						type="number"
						min="1"
						value={payloadForm.amount}
						oninput={(e) => updatePayload({ amount: parseInt(e.currentTarget.value) || 0 })}
						required
						placeholder="100"
					/>
					<span class="help-text">
						💰 How many favour points to award
					</span>
				</div>

				<div class="preview-box favour-preview">
					<div class="favour-display">
						<Star size={48} class="favour-icon" />
						<div class="favour-amount">+{payloadForm.amount || 0}</div>
						<div class="favour-label">Favour Points</div>
					</div>
				</div>

			{:else if editForm.reward_type === 'image'}
				<div class="form-group">
					<label for="image_url">Image URL</label>
					<input
						id="image_url"
						type="url"
						value={payloadForm.url}
						oninput={(e) => updatePayload({ url: e.currentTarget.value })}
						required
						placeholder="https://example.com/celebration.gif"
					/>
					<span class="help-text">
						🖼️ URL to the reward image (GIF, PNG, JPG)
					</span>
				</div>

				<div class="preview-box">
					{#if payloadForm.url}
						<img src={payloadForm.url} alt="Reward preview" class="preview-media" />
					{:else}
						<div class="empty-preview">
							<Image size={48} />
							<p>Enter an image URL to preview</p>
						</div>
					{/if}
				</div>

			{:else if editForm.reward_type === 'video'}
				<div class="form-group">
					<label for="video_url">Video URL</label>
					<input
						id="video_url"
						type="url"
						value={payloadForm.url}
						oninput={(e) => updatePayload({ url: e.currentTarget.value })}
						required
						placeholder="https://example.com/victory.mp4"
					/>
					<span class="help-text">
						🎬 URL to the reward video (MP4, WebM)
					</span>
				</div>

				<div class="preview-box">
					{#if payloadForm.url}
						<video src={payloadForm.url} controls class="preview-media">
							<track kind="captions" />
						</video>
					{:else}
						<div class="empty-preview">
							<Video size={48} />
							<p>Enter a video URL to preview</p>
						</div>
					{/if}
				</div>

			{:else if editForm.reward_type === 'coupon'}
				<div class="form-group">
					<label for="coupon_url">Coupon Pass URL</label>
					<input
						id="coupon_url"
						type="url"
						value={payloadForm.url}
						oninput={(e) => updatePayload({ url: e.currentTarget.value })}
						required
						placeholder="https://example.com/pass.pkpass"
					/>
					<span class="help-text">
						🎫 URL to Apple Wallet .pkpass file
					</span>
				</div>

				<div class="info-box">
					<strong>📱 Apple Wallet Integration</strong>
					<p>Users will be able to add this coupon directly to their Apple Wallet. Make sure your .pkpass file is properly signed and hosted with HTTPS.</p>
					<a href="https://developer.apple.com/wallet/" target="_blank" rel="noopener">
						Learn more about Wallet passes →
					</a>
				</div>
			{/if}
		</div>
	</div>

	<div class="dialog-actions">
		<button type="button" class="btn-cancel" onclick={onCancel}>
			Cancel
		</button>
		<button type="submit" class="btn-save">
			Save Reward
		</button>
	</div>
</form>

<style>
	h2 {
		margin: 0 0 1.5rem 0;
		font-size: 1.5rem;
		font-weight: 600;
		color: #1f2937;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.form-layout {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 2rem;
	}

	.form-column {
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
	}

	.config-column {
		background: #f9fafb;
		padding: 1.25rem;
		border-radius: 8px;
		border: 1px solid #e5e7eb;
	}

	.config-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-weight: 600;
		color: #6b7280;
		font-size: 0.875rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		margin-bottom: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
	}

	input, textarea {
		padding: 0.5rem 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 0.875rem;
		font-family: inherit;
		transition: all 0.2s;
	}

	input::placeholder, textarea::placeholder {
		color: #9ca3af;
	}

	input:focus, textarea:focus {
		outline: none;
		border-color: #10b981;
		box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
	}

	.help-text {
		font-size: 0.75rem;
		color: #6b7280;
		font-style: italic;
	}

	.reward-type-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.75rem;
	}

	.reward-type-card {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.875rem;
		border: 2px solid #e5e7eb;
		border-radius: 8px;
		background: white;
		cursor: pointer;
		transition: all 0.2s;
		text-align: left;
	}

	.reward-type-card:hover {
		border-color: #d1d5db;
		background: #f9fafb;
	}

	.reward-type-card.selected {
		border-color: #10b981;
		background: #ecfdf5;
	}

	.type-info {
		flex: 1;
		min-width: 0;
	}

	.type-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: #1f2937;
	}

	.type-description {
		font-size: 0.75rem;
		color: #6b7280;
		margin-top: 0.125rem;
	}

	.preview-box {
		border: 2px dashed #d1d5db;
		border-radius: 8px;
		padding: 1.5rem;
		min-height: 200px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: white;
	}

	.favour-preview {
		background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
		border-color: #fbbf24;
	}

	.favour-display {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
	}

	:global(.favour-icon) {
		color: #f59e0b;
		fill: #fbbf24;
	}

	.favour-amount {
		font-size: 2.5rem;
		font-weight: 700;
		color: #92400e;
	}

	.favour-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: #92400e;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.preview-media {
		max-width: 100%;
		max-height: 300px;
		border-radius: 6px;
	}

	.empty-preview {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
		color: #9ca3af;
		text-align: center;
	}

	.info-box {
		padding: 1rem;
		background: #eff6ff;
		border: 1px solid #bfdbfe;
		border-radius: 6px;
		font-size: 0.875rem;
	}

	.info-box strong {
		display: block;
		margin-bottom: 0.5rem;
		color: #1e40af;
	}

	.info-box p {
		margin: 0 0 0.5rem 0;
		color: #1e3a8a;
		line-height: 1.5;
	}

	.info-box a {
		color: #2563eb;
		text-decoration: none;
		font-weight: 500;
	}

	.info-box a:hover {
		text-decoration: underline;
	}

	.dialog-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		padding-top: 1rem;
		border-top: 1px solid #e5e7eb;
	}

	.btn-cancel, .btn-save {
		padding: 0.625rem 1.25rem;
		border-radius: 6px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-cancel {
		background: white;
		border: 1px solid #d1d5db;
		color: #374151;
	}

	.btn-cancel:hover {
		background: #f9fafb;
		border-color: #9ca3af;
	}

	.btn-save {
		background: #10b981;
		border: none;
		color: white;
	}

	.btn-save:hover {
		background: #059669;
		transform: translateY(-1px);
		box-shadow: 0 4px 6px -1px rgba(16, 185, 129, 0.3);
	}

	.btn-save:active {
		transform: translateY(0);
	}
</style>
