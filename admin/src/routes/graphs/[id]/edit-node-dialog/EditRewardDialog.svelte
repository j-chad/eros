<script lang="ts">
	import {type FileInfo, type RewardNode, RewardType} from '$lib/types';
	import {api} from '$lib/api';
	import {Calendar, File as FileIcon, Gift, Image, NotebookPen, Sparkles, Star, Upload, Video} from 'lucide-svelte';

	const FILE_BACKED_TYPES = new Set([
		RewardType.IMAGE,
		RewardType.VIDEO,
		RewardType.FILE,
		RewardType.CALENDAR,
		RewardType.WALLET
	]);

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
		description: node.description ?? '',
		reward_type: node.data?.reward_type ?? RewardType.FAVOUR,
		payload: node.data?.payload ?? '',
		give_favours: node.data?.give_favours ?? 0
	});

	let fileInfo: FileInfo | undefined = $state(node.data?.file ?? undefined);
	let uploading = $state(false);
	let uploadError = $state('');
	let dragOver = $state(false);

	let isFileBacked = $derived(FILE_BACKED_TYPES.has(editForm.reward_type));

	async function handleFileUpload(file: globalThis.File) {
		uploading = true;
		uploadError = '';
		try {
			const result = await api.files.upload(node.id, file);
			fileInfo = {
				id: String(result.id),
				filename: result.filename,
				mime_type: result.mime_type,
				size_bytes: result.size_bytes,
				url: `/api/files/${result.id}`
			};
		} catch (err: unknown) {
			uploadError = (err as { body?: { message?: string } })?.body?.message ?? 'Upload failed';
			console.error('File upload failed:', err);
		} finally {
			uploading = false;
		}
	}

	function handleFileInput(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (file) handleFileUpload(file);
		input.value = '';
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		dragOver = false;
		const file = event.dataTransfer?.files?.[0];
		if (file) handleFileUpload(file);
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		dragOver = true;
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
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
				give_favours: editForm.give_favours,
				file: fileInfo
			}
		};

		onSave(updatedNode);
	}

	const rewardTypes = [
		{ value: RewardType.FAVOUR, label: 'Favour Points', icon: Star, description: 'Award currency to the user' },
		{ value: RewardType.IMAGE, label: 'Image', icon: Image, description: 'Show a celebratory image' },
		{ value: RewardType.VIDEO, label: 'Video', icon: Video, description: 'Play a reward video' },
		{ value: RewardType.CALENDAR, label: 'Calendar Event', icon: Calendar, description: 'Add event to calendar' },
		{ value: RewardType.WALLET, label: 'Coupon', icon: Gift, description: 'Apple Wallet coupon' },
		{ value: RewardType.MARKDOWN, label: 'Markdown', icon: NotebookPen, description: 'Custom message with markdown formatting' },
		{ value: RewardType.FILE, label: 'File', icon: FileIcon, description: 'Provide a downloadable file' }
	];

	const fileTypeAccept: Record<string, string> = {
		[RewardType.IMAGE]: 'image/*',
		[RewardType.VIDEO]: 'video/*',
		[RewardType.CALENDAR]: '.ics',
		[RewardType.WALLET]: '.pkpass',
	};
</script>

<h2>Edit Reward Node</h2>

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
					Optional: Award favour points with any reward type
				</span>
			</div>

			<div class="form-group">
				<label>Reward Type</label>
				<div class="reward-type-grid">
					{#each rewardTypes as type (type.value)}
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

			{#if editForm.reward_type === RewardType.FAVOUR}
				<div class="info-box">
					<strong>Favour Points Only</strong>
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

			{:else if editForm.reward_type === RewardType.MARKDOWN}
				<div class="form-group">
					<label for="markdown_content">Markdown Content</label>
					<textarea
						id="markdown_content"
						bind:value={editForm.payload}
						rows="8"
						placeholder="Enter custom message with **markdown** formatting"
					></textarea>
				</div>

				{#if editForm.give_favours > 0}
					<div class="favour-badge">
						<Star size={16} />
						<span>+{editForm.give_favours} Favours</span>
					</div>
				{/if}

			{:else if isFileBacked}
				<!-- File upload UI for IMAGE, VIDEO, FILE, CALENDAR, WALLET -->

				{#if fileInfo}
					<div class="file-info">
						<div class="file-details">
							<div class="file-name">{fileInfo.filename}</div>
							<div class="file-meta">
								{fileInfo.mime_type} &middot; {formatFileSize(fileInfo.size_bytes)}
							</div>
						</div>
					</div>

					{#if editForm.reward_type === RewardType.IMAGE && fileInfo.url}
						<div class="preview-box">
							<img src={fileInfo.url} alt="Reward preview" class="preview-media" />
						</div>
					{:else if editForm.reward_type === RewardType.VIDEO && fileInfo.url}
						<div class="preview-box">
							<video src={fileInfo.url} controls class="preview-media">
								<track kind="captions" />
							</video>
						</div>
					{/if}
				{/if}

				<label
					class="drop-zone"
					class:drag-over={dragOver}
					class:uploading
					ondrop={handleDrop}
					ondragover={handleDragOver}
					ondragleave={() => dragOver = false}
				>
					<input
						type="file"
						accept={fileTypeAccept[editForm.reward_type] ?? '*/*'}
						onchange={handleFileInput}
						disabled={uploading}
						hidden
					/>
					{#if uploading}
						<div class="drop-zone-content">
							<Upload size={32} />
							<p>Uploading...</p>
						</div>
					{:else}
						<div class="drop-zone-content">
							<Upload size={32} />
							<p>{fileInfo ? 'Drop a file to replace, or click to browse' : 'Drop a file here, or click to browse'}</p>
							<span class="help-text">Max 50 MB</span>
						</div>
					{/if}
				</label>

				{#if uploadError}
					<div class="error-text">{uploadError}</div>
				{/if}

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

	/* File upload styles */
	.file-info {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.75rem 1rem;
		background: white;
		border: 1px solid #d1d5db;
		border-radius: 6px;
	}

	.file-details {
		min-width: 0;
	}

	.file-name {
		font-size: 0.875rem;
		font-weight: 600;
		color: #1f2937;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.file-meta {
		font-size: 0.75rem;
		color: #6b7280;
		margin-top: 0.125rem;
	}

	.drop-zone {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1.5rem;
		border: 2px dashed #d1d5db;
		border-radius: 8px;
		background: white;
		cursor: pointer;
		transition: all 0.2s;
		text-align: center;
	}

	.drop-zone:hover {
		border-color: #10b981;
		background: #ecfdf5;
	}

	.drop-zone.drag-over {
		border-color: #10b981;
		background: #d1fae5;
	}

	.drop-zone.uploading {
		opacity: 0.6;
		cursor: wait;
	}

	.drop-zone-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		color: #6b7280;
	}

	.drop-zone-content p {
		margin: 0;
		font-size: 0.875rem;
	}

	.error-text {
		font-size: 0.8125rem;
		color: #dc2626;
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
