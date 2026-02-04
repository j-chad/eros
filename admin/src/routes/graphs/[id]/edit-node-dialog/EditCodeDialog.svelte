<script lang="ts">
	import type { CodeNode } from '$lib/types';

	let {
		node,
		onSave,
		onCancel
	}: {
		node: CodeNode;
		onSave: (node: CodeNode) => void;
		onCancel: () => void;
	} = $props();

	let editForm = $state({
		title: node.title,
		description: node.description ?? '',
		code: node.data?.code ?? ''
	});

	function handleSubmit(event: Event) {
		event.preventDefault();

		const updatedNode: CodeNode = {
			...node,
			title: editForm.title,
			description: editForm.description,
			data: {
				code: editForm.code
			}
		};

		onSave(updatedNode);
	}
</script>

<h2>Edit Code Node</h2>

<form onsubmit={handleSubmit}>
	<div class="form-group">
		<label for="title">Title</label>
		<input
			id="title"
			type="text"
			bind:value={editForm.title}
			required
			placeholder="e.g., Access Code, PIN, Security Key"
		/>
	</div>

	<div class="form-group">
		<label for="description">Description</label>
		<textarea
			id="description"
			bind:value={editForm.description}
			rows="3"
			placeholder="Optional description for this code"
		/>
	</div>

	<div class="form-group">
		<label for="code">Code</label>
		<input
			id="code"
			type="text"
			bind:value={editForm.code}
			required
			placeholder="e.g., 1234, ABC123, ****"
			class="code-input"
		/>
		<span class="help-text">
			{#if editForm.code}
				Code: {editForm.code}
			{:else}
				Enter the code value
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

	.code-input {
		font-family: 'Courier New', monospace;
		letter-spacing: 0.15em;
		font-weight: 600;
	}

	input::placeholder, textarea::placeholder {
		color: #9ca3af;
	}

	input:focus, textarea:focus {
		outline: none;
		border-color: #8b5cf6;
		box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
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
		background: #8b5cf6;
		border: none;
		color: white;
	}

	.btn-save:hover {
		background: #7c3aed;
		transform: translateY(-1px);
		box-shadow: 0 4px 6px -1px rgba(139, 92, 246, 0.3);
	}

	.btn-save:active {
		transform: translateY(0);
	}
</style>
