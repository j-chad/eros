<script lang="ts">
	import type { TimeNode } from '$lib/types';

	let {
		node,
		onSave,
		onCancel
	}: {
		node: TimeNode;
		onSave: (node: TimeNode) => void;
		onCancel: () => void;
	} = $props();

	// Convert UTC ISO string to local datetime-local input value
	function utcToLocal(utcString: string | undefined): string {
		if (!utcString) return '';
		const date = new Date(utcString);
		if (isNaN(date.getTime())) return '';
		// datetime-local expects YYYY-MM-DDTHH:mm format in local time
		const offset = date.getTimezoneOffset();
		const local = new Date(date.getTime() - offset * 60000);
		return local.toISOString().slice(0, 16);
	}

	// Convert local datetime-local input value to UTC ISO string
	function localToUtc(localString: string): string {
		if (!localString) return '';
		// The input value is in local time, so new Date() parses it correctly
		return new Date(localString).toISOString();
	}

	let editForm = $state({
		title: node.title,
		description: node.description || '',
		unlockAt: utcToLocal(node.data?.unlock_at),
	});

	function handleSubmit(event: Event) {
		event.preventDefault();

		const updatedNode: TimeNode = {
			...node,
			title: editForm.title,
			description: editForm.description,
			data: {
				unlock_at: localToUtc(editForm.unlockAt)
			}
		};

		onSave(updatedNode);
	}
</script>

<h2>Edit Time Gate</h2>

<form onsubmit={handleSubmit}>
	<div class="form-group">
		<label for="title">Title</label>
		<input
			id="title"
			type="text"
			bind:value={editForm.title}
			required
			placeholder="e.g., Dinner Reservation, Surprise Time"
		/>
	</div>

	<div class="form-group">
		<label for="description">Description</label>
		<textarea
			id="description"
			bind:value={editForm.description}
			rows="3"
			placeholder="Optional description shown to the user"
		></textarea>
	</div>

	<div class="form-group">
		<label for="unlock-at">Unlock Time</label>
		<input
			id="unlock-at"
			type="datetime-local"
			bind:value={editForm.unlockAt}
			required
		/>
		<span class="help-text">
			{#if editForm.unlockAt}
				Unlocks at {new Date(editForm.unlockAt).toLocaleString()} (your local time)
			{:else}
				Select when this gate should unlock
			{/if}
		</span>
	</div>

	<div class="dialog-actions">
		<button type="button" class="btn-cancel" onclick={onCancel}>
			Cancel
		</button>
		<button type="submit" class="btn-save">
			Save Changes
		</button>
	</div>
</form>

<style>
	h2 {
		margin: 0 0 1.5rem 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: #1f2937;
	}

	form {
		display: flex;
		flex-direction: column;
	}

	.form-group {
		margin-bottom: 1rem;
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
		border-color: #f59e0b;
		box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.1);
	}

	.help-text {
		font-size: 0.75rem;
		color: #6b7280;
		font-style: italic;
	}

	.dialog-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		margin-top: 1.5rem;
		padding-top: 1rem;
		border-top: 1px solid #e5e7eb;
	}

	.btn-cancel, .btn-save {
		padding: 0.5rem 1rem;
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
		background: #f59e0b;
		border: none;
		color: white;
	}

	.btn-save:hover {
		background: #d97706;
		transform: translateY(-1px);
		box-shadow: 0 4px 6px -1px rgba(245, 158, 11, 0.3);
	}

	.btn-save:active {
		transform: translateY(0);
	}
</style>
