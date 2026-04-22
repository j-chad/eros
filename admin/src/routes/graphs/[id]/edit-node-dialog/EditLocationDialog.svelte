<script lang="ts">
	import type { LocationNode } from '$lib/types';
	import MapPicker from '$lib/components/MapPicker.svelte';

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
		radius_m: node.data?.radius_m ?? 0,
		show_hint: node.data?.show_hint ?? false,
		hint_latitude: node.data?.hint?.latitude ?? 0,
		hint_longitude: node.data?.hint?.longitude ?? 0,
		hint_radius_m: node.data?.hint?.radius_m ?? 500
	});

	// When hint is enabled and hint coords are still at default (0,0),
	// initialize them near the target location with a small offset.
	let hintInitialized = node.data?.show_hint ?? false;
	$effect(() => {
		if (editForm.show_hint && !hintInitialized) {
			hintInitialized = true;
			if (editForm.hint_latitude === 0 && editForm.hint_longitude === 0
				&& (editForm.latitude !== 0 || editForm.longitude !== 0)) {
				// Offset by ~radius (min 10m) in latitude (~1 degree ≈ 111km)
				const offsetM = Math.max(editForm.radius_m, 10);
				const offsetDeg = offsetM / 111_000;
				editForm.hint_latitude = Math.round((editForm.latitude + offsetDeg) * 1e6) / 1e6;
				editForm.hint_longitude = editForm.longitude;
			}
		}
	});

	function handleSubmit(event: Event) {
		console.log('Form submit triggered');
		event.preventDefault();
		event.stopPropagation();

		const updatedNode: LocationNode = {
			...node,
			title: editForm.title,
			description: editForm.description,
			data: {
				latitude: editForm.latitude,
				longitude: editForm.longitude,
				radius_m: editForm.radius_m,
				show_hint: editForm.show_hint,
				hint: editForm.show_hint ? {
					latitude: editForm.hint_latitude,
					longitude: editForm.hint_longitude,
					radius_m: editForm.hint_radius_m
				} : null
			}
		};

		console.log('Calling onSave with:', updatedNode);
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

	<div class="form-group">
		<label>Location</label>
		<span class="help-text">Click the map to set the target location. Drag the pin to reposition.</span>
		<MapPicker
			bind:latitude={editForm.latitude}
			bind:longitude={editForm.longitude}
			bind:radiusM={editForm.radius_m}
			showHint={editForm.show_hint}
			bind:hintLatitude={editForm.hint_latitude}
			bind:hintLongitude={editForm.hint_longitude}
			bind:hintRadiusM={editForm.hint_radius_m}
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
			/>
		</div>
		<div class="form-group">
			<label for="radius">Radius (m)</label>
			<input
				id="radius"
				type="number"
				min="0"
				step="1"
				bind:value={editForm.radius_m}
				required
			/>
		</div>
	</div>

	<div class="form-group">
		<label class="toggle-label">
			<input type="checkbox" bind:checked={editForm.show_hint} />
			Show approximate location to user
		</label>
		<span class="help-text">
			Displays a shaded circle on a map as a hint — the center can be offset from the real location.
		</span>
	</div>

	{#if editForm.show_hint}
		<div class="hint-section">
			<h3>Hint Configuration</h3>
			<span class="help-text">Drag the purple pin to offset the hint center from the real location.</span>
			<div class="form-row">
				<div class="form-group">
					<label for="hint-latitude">Hint Latitude</label>
					<input
						id="hint-latitude"
						type="number"
						step="0.000001"
						bind:value={editForm.hint_latitude}
						required
					/>
				</div>
				<div class="form-group">
					<label for="hint-longitude">Hint Longitude</label>
					<input
						id="hint-longitude"
						type="number"
						step="0.000001"
						bind:value={editForm.hint_longitude}
						required
					/>
				</div>
				<div class="form-group">
					<label for="hint-radius">Hint Radius (m)</label>
					<input
						id="hint-radius"
						type="number"
						min="1"
						step="1"
						bind:value={editForm.hint_radius_m}
						required
					/>
				</div>
			</div>
		</div>
	{/if}

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
		grid-template-columns: 1fr 1fr 1fr;
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

	.toggle-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
		cursor: pointer;
	}

	.toggle-label input[type="checkbox"] {
		width: 1rem;
		height: 1rem;
		accent-color: #ec4899;
	}

	.hint-section {
		margin-top: 0.5rem;
		padding: 1rem;
		background: #fdf2f8;
		border: 1px solid #fbcfe8;
		border-radius: 8px;
	}

	.hint-section h3 {
		margin: 0 0 0.75rem 0;
		font-size: 0.875rem;
		font-weight: 600;
		color: #be185d;
	}
</style>
