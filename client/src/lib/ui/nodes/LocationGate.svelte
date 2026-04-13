<script lang="ts">
	import type { LocationNode } from '$lib/types/graph';
	import type { UnlockResult } from '$lib/api/graph.api';
	import { MapPin } from 'lucide-svelte';
	import { useOnlineStatus } from '$lib/online.svelte';

	const { node, graphId, onUnlock }: { 
		node: LocationNode; 
		graphId: string; 
		onUnlock: (result: UnlockResult) => void 
	} = $props();

	const isOnline = $derived(useOnlineStatus());
	
	let isChecking = $state(false);
	let errorMessage = $state<string | null>(null);
	let isPermissionDenied = $state(false);

	async function handleCheckLocation() {
		if (isChecking || !isOnline) return;

		isChecking = true;
		errorMessage = null;
		isPermissionDenied = false;

		try {
			const position = await getCurrentPosition();
			const accuracy = position.coords.accuracy;
			const radius = node.data?.radius_m ?? 10;

			// Warn about poor accuracy
			if (accuracy > radius * 2) {
				errorMessage = `Your location accuracy is low (${Math.round(accuracy)}m). Try moving to an open area.`;
				// Still allow the user to proceed
			}

			const payload = `${position.coords.latitude},${position.coords.longitude}`;
			const { unlockNode } = await import('$lib/services/graph');
			const result = await unlockNode(graphId, node.id, payload);
			onUnlock(result);
		} catch (error: any) {
			if (error.name === 'PermissionDenied') {
				isPermissionDenied = true;
				errorMessage = "Location access is needed to continue. Check your browser settings and try again.";
			} else if (error.status === 403) {
				errorMessage = "You're not at the right spot yet. Get closer and try again.";
			} else if (error.status === 429) {
				errorMessage = "Too many attempts. Try again shortly.";
			} else {
				errorMessage = "Something went wrong getting your location. Try again.";
			}
		} finally {
			isChecking = false;
		}
	}

	function getCurrentPosition(): Promise<GeolocationPosition> {
		return new Promise((resolve, reject) => {
			if (!navigator.geolocation) {
				reject(new Error('Geolocation not supported'));
				return;
			}

			navigator.geolocation.getCurrentPosition(
				resolve,
				(error) => {
					const errorTypes = {
						[error.PERMISSION_DENIED]: 'PermissionDenied',
						[error.POSITION_UNAVAILABLE]: 'PositionUnavailable',
						[error.TIMEOUT]: 'Timeout'
					};
					const wrappedError = new Error(error.message);
					wrappedError.name = errorTypes[error.code] || 'GeolocationError';
					reject(wrappedError);
				},
				{
					enableHighAccuracy: true,
					timeout: 15000,
					maximumAge: 0
				}
			);
		});
	}
</script>

<div class="flex flex-col items-center gap-6 text-center animate-popIn">
	<div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
		<MapPin size={28} class="text-primary" />
	</div>

	<div class="flex flex-col gap-2">
		<h1 class="text-2xl font-extrabold">{node.title}</h1>
		{#if node.description}
			<p class="text-sm opacity-70 leading-relaxed">{node.description}</p>
		{/if}
	</div>

	<div class="w-full flex flex-col gap-3">
		<div class="bg-base-200 rounded-2xl px-5 py-4 flex flex-col gap-1">
			<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">Head to the location</p>
			<p class="text-sm opacity-70">
				You need to be within {node.data?.radius_m ?? 10}&nbsp;metres of the destination.
			</p>
		</div>
		
		{#if errorMessage}
			<div class="alert rounded-2xl text-sm" class:alert-error={!errorMessage.includes('accuracy')} class:alert-warning={errorMessage.includes('accuracy')}>
				{errorMessage}
				{#if isPermissionDenied}
					<button onclick={handleCheckLocation} class="btn btn-sm btn-outline">Try again</button>
				{/if}
			</div>
		{/if}
		
		<button 
			onclick={handleCheckLocation}
			disabled={isChecking || !isOnline}
			class="btn btn-primary rounded-2xl w-full"
			class:loading={isChecking}
		>
			{#if isChecking}
				<span class="loading loading-spinner loading-sm"></span>
				Getting your location...
			{:else}
				Check location
			{/if}
		</button>
		
		{#if !isOnline}
			<p class="text-xs opacity-50 text-center">You're offline. Connect to check your location.</p>
		{/if}
	</div>
</div>
