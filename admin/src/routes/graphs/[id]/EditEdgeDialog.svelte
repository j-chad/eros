<script lang="ts">
	import type { Edge } from '$lib/types';

	let {
		edge,
		allEdges,
		onSave,
		onClose
	}: {
		edge: Edge | null;
		allEdges: Edge[];
		onSave: (edge: Edge) => void;
		onClose: () => void;
	} = $props();

	let dialog: HTMLDialogElement;

	let editForm = $state({
		choice_label: ''
	});

	const neighbourEdges = $derived(() => {
		if (!edge) return [];
		return allEdges.filter(e => e.from === edge.from && e.id !== edge.id);
	});

	// Check if this edge's source has multiple outgoing edges
	const hasMultipleOutgoing = $derived(() => {
		if (!edge) return false;
		return neighbourEdges().length > 0;
	});

	// Check if the label is unique among edges from the same source
	const isLabelUnique = $derived(() => {
		if (!edge) return true;
		return !neighbourEdges().some(e => e.choice_label?.trim() === editForm.choice_label.trim());
	});

	const showError = $derived(hasMultipleOutgoing() && !editForm.choice_label.trim());
	const showDuplicateError = $derived(!!editForm.choice_label.trim() && !isLabelUnique());

	$effect(() => {
		if (edge) {
			editForm = {
				choice_label: edge.choice_label || ''
			};
			requestAnimationFrame(() => {
				dialog?.showModal();
			});
		} else {
			dialog?.close();
		}
	});

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		event.stopPropagation();

		if (!edge) return;

		// Validate if multiple edges exist
		if (hasMultipleOutgoing() && !editForm.choice_label.trim()) {
			return;
		}

		if (!isLabelUnique()) {
			return;
		}

		const updatedEdge: Edge = {
			...edge,
			choice_label: editForm.choice_label.trim()
		};

		onSave(updatedEdge);
		handleClose();
	}

	function handleClose() {
		onClose();
	}

	function handleBackdropClick(event: MouseEvent) {
		const dialogElement = event.currentTarget as HTMLDialogElement;
		const rect = dialogElement.getBoundingClientRect();
		const isInDialog = (
			rect.top <= event.clientY &&
			event.clientY <= rect.top + rect.height &&
			rect.left <= event.clientX &&
			event.clientX <= rect.left + rect.width
		);

		if (!isInDialog) {
			handleClose();
		}
	}
</script>

<dialog bind:this={dialog} class="edit-dialog" onclose={handleClose} onclick={handleBackdropClick}>
	<div class="dialog-content">
		{#if edge}
			{#key edge.id}
				<h2>Edit Connection</h2>

				<form onsubmit={handleSubmit}>
					<div class="form-group">
						<label for="choice_label">
							Choice Label
							{#if hasMultipleOutgoing()}
								<span class="required">*</span>
							{/if}
						</label>
						<input
							id="choice_label"
							type="text"
							bind:value={editForm.choice_label}
							placeholder={hasMultipleOutgoing() ? "e.g., Yes, No, Maybe" : "Optional label"}
							required={hasMultipleOutgoing()}
						/>
						{#if showError}
							<span class="error-text">
								⚠️ Label required when node has multiple outgoing connections
							</span>
						{:else if showDuplicateError}
							<span class="error-text">
								⚠️ Label must be unique among connections from the same node
							</span>
						{:else if hasMultipleOutgoing()}
							<span class="help-text">
								💡 This node has multiple outgoing connections, so each must have a unique label
							</span>
						{:else}
							<span class="help-text">
								💡 Optional: Add a label to describe this connection
							</span>
						{/if}
					</div>

					<div class="dialog-actions">
						<button type="button" class="btn-cancel" onclick={handleClose}>
							Cancel
						</button>
						<button type="submit" class="btn-save" disabled={showError || showDuplicateError}>
							Save Changes
						</button>
					</div>
				</form>
			{/key}
		{/if}
	</div>
</dialog>

<style>
	.edit-dialog {
		border: none;
		border-radius: 12px;
		padding: 0;
		max-width: 500px;
		width: 90%;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
	}

	.edit-dialog::backdrop {
		background: rgba(0, 0, 0, 0.5);
	}

	.dialog-content {
		padding: 1.5rem;
	}

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

	.required {
		color: #dc2626;
	}

	input {
		padding: 0.5rem 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 0.875rem;
		font-family: inherit;
		transition: all 0.2s;
	}

	input::placeholder {
		color: #9ca3af;
	}

	input:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.help-text {
		font-size: 0.75rem;
		color: #6b7280;
		font-style: italic;
	}

	.error-text {
		font-size: 0.75rem;
		color: #dc2626;
		font-weight: 500;
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
		background: #3b82f6;
		border: none;
		color: white;
	}

	.btn-save:hover:not(:disabled) {
		background: #2563eb;
		transform: translateY(-1px);
		box-shadow: 0 4px 6px -1px rgba(59, 130, 246, 0.3);
	}

	.btn-save:active:not(:disabled) {
		transform: translateY(0);
	}

	.btn-save:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
