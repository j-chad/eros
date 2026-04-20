<script lang="ts">
	import { X, Plus } from 'lucide-svelte';
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
		codes: node.data?.codes?.length ? [...node.data.codes] : ['']
	});

	function addCode() {
		editForm.codes = [...editForm.codes, ''];
	}

	function removeCode(index: number) {
		editForm.codes = editForm.codes.filter((_, i) => i !== index);
	}

	function handleSubmit(event: Event) {
		event.preventDefault();

		const codes = editForm.codes.map((c) => c.trim()).filter(Boolean);
		if (codes.length === 0) return;

		const updatedNode: CodeNode = {
			...node,
			title: editForm.title,
			description: editForm.description,
			data: { codes }
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
		<div class="codes-header">
			<label>Codes</label>
			<button type="button" class="btn-add" onclick={addCode}>
				<Plus size={14} />
				Add code
			</button>
		</div>
		<div class="codes-list">
			{#each editForm.codes as code, i}
				<div class="code-row">
					<input
						type="text"
						bind:value={editForm.codes[i]}
						required={i === 0}
						placeholder="e.g., 1234, ABC123"
						class="code-input"
					/>
					{#if editForm.codes.length > 1}
						<button
							type="button"
							class="btn-remove"
							onclick={() => removeCode(i)}
							aria-label="Remove code"
						>
							<X size={14} />
						</button>
					{/if}
				</div>
			{/each}
		</div>
		<span class="help-text">
			Any one of these codes will unlock the gate (case-insensitive).
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
		border-color: #8b5cf6;
		box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
	}

	.codes-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.codes-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.code-row {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.code-row input {
		flex: 1;
	}

	.code-input {
		font-family: 'Courier New', monospace;
		letter-spacing: 0.15em;
		font-weight: 600;
	}

	.btn-add {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.25rem 0.625rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		background: white;
		font-size: 0.75rem;
		font-weight: 500;
		color: #374151;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-add:hover {
		background: #f9fafb;
		border-color: #8b5cf6;
		color: #8b5cf6;
	}

	.btn-remove {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 1.75rem;
		height: 1.75rem;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		background: white;
		color: #9ca3af;
		cursor: pointer;
		transition: all 0.2s;
		flex-shrink: 0;
	}

	.btn-remove:hover {
		background: #fef2f2;
		border-color: #fca5a5;
		color: #ef4444;
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
