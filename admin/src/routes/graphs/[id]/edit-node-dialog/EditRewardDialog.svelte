<script lang="ts">
	import type {RewardNode} from '$lib/types';
	import {Calendar, Gift, Image, NotebookPen, Sparkles, Star, Video} from 'lucide-svelte';

	let {
		node,
		onSave,
		onCancel
	}: {
		node: RewardNode;
		onSave: (node: RewardNode) => void;
		onCancel: () => void;
	} = $props();

	let formElement: HTMLFormElement;

	let editForm = $state({
		title: node.title,
		description: node.description || '',
		reward_type: node.data?.reward_type || 'favour',
		payload: node.data?.payload || '',
		give_favours: node.data?.give_favours || 0
	});

	// Parse payload based on type for easier editing
	let payloadForm = $derived.by(() => {
		try {
			const parsed = editForm.payload ? JSON.parse(editForm.payload) : {};
			return {
				url: parsed.url || '',
				coupon_data: parsed.coupon_data || {},
				event_title: parsed.event_title || '',
				event_start: parsed.event_start || '',
				event_end: parsed.event_end || '',
				event_location: parsed.event_location || '',
				event_notes: parsed.event_notes || ''
			};
		} catch {
			return {
				url: '',
				coupon_data: {},
				event_title: '',
				event_start: '',
				event_end: '',
				event_location: '',
				event_notes: ''
			};
		}
	});

	function updatePayload(updates: Partial<typeof payloadForm>) {
		const updated = { ...payloadForm, ...updates };

		// Create clean payload based on type
		let payload: any = {};
		if (editForm.reward_type === 'video' || editForm.reward_type === 'image') {
			payload = { url: updated.url };
		} else if (editForm.reward_type === 'coupon') {
			payload = {
				url: updated.url,
				coupon_data: updated.coupon_data
			};
		} else if (editForm.reward_type === 'calendar') {
			payload = {
				event_title: updated.event_title,
				event_start: updated.event_start,
				event_end: updated.event_end,
				event_location: updated.event_location,
				event_notes: updated.event_notes
			};
		}
		// favour type has no payload, just give_favours

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
				give_favours: editForm.give_favours
			}
		};

		onSave(updatedNode);
	}

	const rewardTypes = [
		{ value: 'favour', label: 'Favour Points', icon: Star, description: 'Award currency to the user' },
		{ value: 'image', label: 'Image', icon: Image, description: 'Show a celebratory image' },
		{ value: 'video', label: 'Video', icon: Video, description: 'Play a reward video' },
		{ value: 'calendar', label: 'Calendar Event', icon: Calendar, description: 'Add event to calendar' },
		{ value: 'coupon', label: 'Coupon', icon: Gift, description: 'Apple Wallet coupon' },
		{ value: 'markdown', label: 'Markdown', icon: NotebookPen, description: 'Custom message with markdown formatting' }
	];
</script>

<h2>🎁 Edit Reward Node</h2>

<form bind:this={formElement} onsubmit={handleSubmit}>
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
				<label for="give_favours">Favour Points</label>
				<input
					id="give_favours"
					type="number"
					min="0"
					bind:value={editForm.give_favours}
					placeholder="0"
				/>
				<span class="help-text">
					💰 Optional: Award favour points with any reward type
				</span>
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
				<div class="info-box">
					<strong>💰 Favour Points Only</strong>
					<p>This reward will award favour points to the user. Configure the amount in the "Favour Points" field on the left.</p>
				</div>

				{#if editForm.give_favours > 0}
					<div class="preview-box favour-preview">
						<div class="favour-display">
							<Star size={48} class="favour-icon" />
							<div class="favour-amount">+{editForm.give_favours}</div>
							<div class="favour-label">Favour Points</div>
						</div>
					</div>
				{/if}

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

				{#if editForm.give_favours > 0}
					<div class="favour-badge">
						<Star size={16} />
						<span>+{editForm.give_favours} Favours</span>
					</div>
				{/if}

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

				{#if editForm.give_favours > 0}
					<div class="favour-badge">
						<Star size={16} />
						<span>+{editForm.give_favours} Favours</span>
					</div>
				{/if}

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

			{:else if editForm.reward_type === 'calendar'}
				<div class="form-group">
					<label for="event_title">Event Title</label>
					<input
						id="event_title"
						type="text"
						value={payloadForm.event_title}
						oninput={(e) => updatePayload({ event_title: e.currentTarget.value })}
						required
						placeholder="Team Meeting"
					/>
				</div>

				<div class="form-row">
					<div class="form-group">
						<label for="event_start">Start Date & Time</label>
						<input
							id="event_start"
							type="datetime-local"
							value={payloadForm.event_start}
							oninput={(e) => updatePayload({ event_start: e.currentTarget.value })}
							required
						/>
					</div>

					<div class="form-group">
						<label for="event_end">End Date & Time</label>
						<input
							id="event_end"
							type="datetime-local"
							value={payloadForm.event_end}
							oninput={(e) => updatePayload({ event_end: e.currentTarget.value })}
							required
						/>
					</div>
				</div>

				<div class="form-group">
					<label for="event_location">Location (Optional)</label>
					<input
						id="event_location"
						type="text"
						value={payloadForm.event_location}
						oninput={(e) => updatePayload({ event_location: e.currentTarget.value })}
						placeholder="Conference Room A"
					/>
				</div>

				<div class="form-group">
					<label for="event_notes">Notes (Optional)</label>
					<textarea
						id="event_notes"
						oninput={(e) => updatePayload({ event_notes: e.currentTarget.value })}
						rows="3"
						placeholder="Bring your laptop"
					>{payloadForm.event_notes}</textarea>
				</div>

				{#if editForm.give_favours > 0}
					<div class="favour-badge">
						<Star size={16} />
						<span>+{editForm.give_favours} Favours</span>
					</div>
				{/if}

				<div class="info-box">
					<strong>📅 Calendar Integration</strong>
					<p>Users will be prompted to add this event to their device calendar.</p>
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

				{#if editForm.give_favours > 0}
					<div class="favour-badge">
						<Star size={16} />
						<span>+{editForm.give_favours} Favours</span>
					</div>
				{/if}

				<div class="info-box">
					<strong>📱 Apple Wallet Integration</strong>
					<p>Users will be able to add this coupon directly to their Apple Wallet. Make sure your .pkpass file is properly signed and hosted with HTTPS.</p>
					<a href="https://developer.apple.com/wallet/" target="_blank" rel="noopener">
						Learn more about Wallet passes →
					</a>
				</div>
			{:else if editForm.reward_type === 'markdown'}
				<div class="form-group">
					<label for="markdown_content">Markdown Content</label>
					<textarea
						id="markdown_content"
						oninput={(e) => updatePayload({ url: e.currentTarget.value })}
						rows="6"
						placeholder="Enter custom message with **markdown** formatting"
					>{payloadForm.url}</textarea>
				</div>

				{#if editForm.give_favours > 0}
					<div class="favour-badge">
						<Star size={16} />
						<span>+{editForm.give_favours} Favours</span>
					</div>
				{/if}
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

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.75rem;
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

	.favour-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.375rem 0.75rem;
		background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
		border: 1px solid #fbbf24;
		border-radius: 6px;
		font-size: 0.875rem;
		font-weight: 600;
		color: #92400e;
		width: fit-content;
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
