<script lang="ts">
	import type { LocationHint } from '$lib/types/graph';

	const { hint }: { hint: LocationHint } = $props();

	let mapContainer: HTMLDivElement;
	let isLoaded = $state(false);

	$effect(() => {
		let map: any;
		let cancelled = false;

		(async () => {
			const L = (await import('leaflet')).default;
			await import('leaflet/dist/leaflet.css');

			if (cancelled) return;
			isLoaded = true;

			const isMarkerMode = hint.radius_m === 0;
			const center: [number, number] = [hint.latitude, hint.longitude];

			map = L.map(mapContainer, {
				center,
				zoom: isMarkerMode ? 16 : 13,
				scrollWheelZoom: false,
				attributionControl: false
			});

			L.control.attribution({ position: 'bottomright', prefix: false })
				.addAttribution('&copy; <a href="https://openstreetmap.org/copyright">OSM</a>')
				.addTo(map);

			L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
				maxZoom: 19,
				className: 'hint-map-tiles'
			}).addTo(map);

			if (isMarkerMode) {
				const icon = L.divIcon({
					className: 'hint-map-marker-icon',
					html: '<div class="hint-map-pin"></div>',
					iconSize: [20, 20],
					iconAnchor: [10, 10]
				});
				L.marker(center, { icon, interactive: false }).addTo(map);
			} else {
				const circle = L.circle(center, {
					radius: hint.radius_m,
					color: '#e8457a',
					fillColor: '#e8457a',
					fillOpacity: 0.15,
					weight: 2
				}).addTo(map);
				map.fitBounds(circle.getBounds().pad(0.15));
			}

			// Fix sizing in case container wasn't fully laid out
			setTimeout(() => map?.invalidateSize(), 50);
		})();

		return () => {
			cancelled = true;
			map?.remove();
		};
	});
</script>

<div class="hint-map-wrapper">
	{#if !isLoaded}
		<div class="hint-map-placeholder">
			<span class="loading loading-dots loading-sm text-primary"></span>
		</div>
	{/if}
	<div class="hint-map-container" class:loaded={isLoaded} bind:this={mapContainer}></div>
</div>

<style>
	.hint-map-wrapper {
		position: relative;
		width: 100%;
		aspect-ratio: 4 / 3;
		max-height: 300px;
		border-radius: 1rem;
		overflow: hidden;
	}

	.hint-map-placeholder {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		background: oklch(var(--b2));
		border-radius: 1rem;
	}

	.hint-map-container {
		position: relative;
		width: 100%;
		height: 100%;
		opacity: 0;
		transition: opacity 0.3s ease;
	}

	.hint-map-container.loaded {
		opacity: 1;
	}

	/* Pink-tinted map tiles */
	:global(.hint-map-tiles) {
		filter: saturate(0.5) hue-rotate(330deg) sepia(0.1);
	}

	/* Transparent pink overlay */
	.hint-map-container.loaded::after {
		content: '';
		position: absolute;
		inset: 0;
		background: oklch(0.65 0.2 350 / 0.12);
		pointer-events: none;
		z-index: 401; /* above tiles, below controls (1000) */
	}

	/* Custom marker pin */
	:global(.hint-map-marker-icon) {
		background: none !important;
		border: none !important;
	}

	:global(.hint-map-pin) {
		width: 16px;
		height: 16px;
		border-radius: 50%;
		background: #e8457a;
		border: 2.5px solid white;
		box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
	}

	/* Restyle Leaflet zoom controls */
	:global(.hint-map-wrapper .leaflet-control-zoom) {
		border: none !important;
		border-radius: 0.75rem !important;
		overflow: hidden;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
	}

	:global(.hint-map-wrapper .leaflet-control-zoom a) {
		background: oklch(var(--b1)) !important;
		color: oklch(var(--bc)) !important;
		border-bottom-color: oklch(var(--b2)) !important;
		width: 32px !important;
		height: 32px !important;
		line-height: 32px !important;
		font-size: 16px !important;
	}

	:global(.hint-map-wrapper .leaflet-control-zoom a:hover) {
		background: oklch(var(--b2)) !important;
	}

	/* Attribution styling */
	:global(.hint-map-wrapper .leaflet-control-attribution) {
		background: oklch(var(--b1) / 0.7) !important;
		color: oklch(var(--bc) / 0.5) !important;
		font-size: 0.6rem !important;
		border-radius: 0.5rem 0 0 0 !important;
		padding: 2px 6px !important;
	}

	:global(.hint-map-wrapper .leaflet-control-attribution a) {
		color: oklch(var(--p)) !important;
	}
</style>
