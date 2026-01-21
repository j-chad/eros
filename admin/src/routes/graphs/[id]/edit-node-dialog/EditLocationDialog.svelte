<script lang="ts">
	import type { LocationNode } from '$lib/types';

	let {
		node,
		onSave,
		onCancel
	}: {
		node: LocationNode;
		onSave: (node: LocationNode) => void;
		onCancel: () => void;
	} = $props();

	let editForm = $state({
		title: node.title,
		description: node.description || '',
		latitude: node.data?.latitude ?? 0,
		longitude: node.data?.longitude ?? 0,
		radius_meters: node.data?.radius_meters ?? 0
	});

	function handleSubmit(event: Event) {
		event.preventDefault();

		const updatedNode: LocationNode = {
			...node,
			title: editForm.title,
			description: editForm.description,
			data: {
				latitude: editForm.latitude,
				longitude: editForm.longitude,
				radius_meters: editForm.radius_meters
			}
		};

		onSave(updatedNode);
	}
</script>

<h2>Edit Location Node</h2>

<form onsubmit={handleSubmit}>
	<div class="form-group">
		<label for="title">Title</label>
		<input
			id="title"
			type="text"
			bind:value={editForm.title}
			required
			placeholder="e.g., Home, Office, Park"
		/>
	</div>

	<div class="form-group">
		<label for="description">Description</label>
		<textarea
			id="description"
			bind:value={editForm.description}
			rows="3"
			placeholder="Optional description for this location"
		/>
	</div>

	<div class="form-row">
		<div class="form-group">
			<label for="latitude">Latitude</label>
			<input
				id="latitude"
				type="number"
				step="0.000001"
				bind:value={editForm.latitude}
				required
				placeholder="e.g., -36.848461"
			/>
		</div>

		<div class="form-group">
			<label for="longitude">Longitude</label>
			<input
				id="longitude"
				type="number"
				step="0.000001"
				bind:value={editForm.longitude}
				required
				placeholder="e.g., 174.763336"
			/>
		</div>
	</div>

	<div class="form-group">
		<label for="radius">Radius (meters)</label>
		<input
			id="radius"
			type="number"
			min="0"
			step="1"
			bind:value={editForm.radius_meters}
			required
			placeholder="e.g., 100"
		/>
		<span class="help-text">
			{#if editForm.radius_meters > 0}
				Triggers when within {editForm.radius_meters}m of the location
			{:else}
				Enter a radius in meters
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

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
		margin-bottom: 1rem;
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
		border-color: #ec4899;
		box-shadow: 0 0 0 3px rgba(236, 72, 153, 0.1);
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
		background: #ec4899;
		border: none;
		color: white;
	}

	.btn-save:hover {
		background: #db2777;
		transform: translateY(-1px);
		box-shadow: 0 4px 6px -1px rgba(236, 72, 153, 0.3);
	}

	.btn-save:active {
		transform: translateY(0);
	}
</style>
