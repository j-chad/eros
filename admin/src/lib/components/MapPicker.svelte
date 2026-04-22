<script lang="ts">
	import L from 'leaflet';
	import 'leaflet/dist/leaflet.css';

	let {
		latitude = $bindable(0),
		longitude = $bindable(0),
		radiusM = $bindable(0),
		showHint = false,
		hintLatitude = $bindable(0),
		hintLongitude = $bindable(0),
		hintRadiusM = $bindable(500)
	}: {
		latitude: number;
		longitude: number;
		radiusM: number;
		showHint?: boolean;
		hintLatitude?: number;
		hintLongitude?: number;
		hintRadiusM?: number;
	} = $props();

	let mapContainer: HTMLDivElement;
	let map: L.Map | undefined;
	let marker: L.Marker | undefined;
	let radiusCircle: L.Circle | undefined;
	let hintMarker: L.Marker | undefined;
	let hintCircle: L.Circle | undefined;

	const targetIcon = L.divIcon({
		className: 'map-picker-icon map-picker-icon--target',
		html: '<div class="map-picker-pin"></div>',
		iconSize: [20, 20],
		iconAnchor: [10, 10]
	});

	const hintIcon = L.divIcon({
		className: 'map-picker-icon map-picker-icon--hint',
		html: '<div class="map-picker-pin map-picker-pin--hint"></div>',
		iconSize: [20, 20],
		iconAnchor: [10, 10]
	});

	function hasValidCoords(lat: number, lng: number): boolean {
		return lat !== 0 || lng !== 0;
	}

	function updateTargetMarker() {
		if (!map) return;
		const latlng = L.latLng(latitude, longitude);

		if (!marker) {
			marker = L.marker(latlng, { icon: targetIcon, draggable: true }).addTo(map);
			marker.on('dragend', () => {
				const pos = marker!.getLatLng();
				latitude = Math.round(pos.lat * 1e6) / 1e6;
				longitude = Math.round(pos.lng * 1e6) / 1e6;
				updateRadiusCircle();
			});
		} else {
			marker.setLatLng(latlng);
		}
		updateRadiusCircle();
	}

	function updateRadiusCircle() {
		if (!map) return;
		const latlng = L.latLng(latitude, longitude);

		if (radiusM > 0) {
			if (!radiusCircle) {
				radiusCircle = L.circle(latlng, {
					radius: radiusM,
					color: '#ec4899',
					fillColor: '#ec4899',
					fillOpacity: 0.1,
					weight: 1,
					dashArray: '4 4'
				}).addTo(map);
			} else {
				radiusCircle.setLatLng(latlng);
				radiusCircle.setRadius(radiusM);
			}
		} else if (radiusCircle) {
			radiusCircle.remove();
			radiusCircle = undefined;
		}
	}

	function updateHintLayers() {
		if (!map) return;

		if (showHint) {
			const latlng = L.latLng(hintLatitude, hintLongitude);

			if (!hintMarker) {
				hintMarker = L.marker(latlng, { icon: hintIcon, draggable: true }).addTo(map);
				hintMarker.on('dragend', () => {
					const pos = hintMarker!.getLatLng();
					hintLatitude = Math.round(pos.lat * 1e6) / 1e6;
					hintLongitude = Math.round(pos.lng * 1e6) / 1e6;
					updateHintCircle();
				});
			} else {
				hintMarker.setLatLng(latlng);
			}

			updateHintCircle();
		} else {
			if (hintMarker) {
				hintMarker.remove();
				hintMarker = undefined;
			}
			if (hintCircle) {
				hintCircle.remove();
				hintCircle = undefined;
			}
		}
	}

	function updateHintCircle() {
		if (!map) return;
		const latlng = L.latLng(hintLatitude, hintLongitude);

		if (hintRadiusM > 0) {
			if (!hintCircle) {
				hintCircle = L.circle(latlng, {
					radius: hintRadiusM,
					color: '#8b5cf6',
					fillColor: '#8b5cf6',
					fillOpacity: 0.1,
					weight: 2
				}).addTo(map);
			} else {
				hintCircle.setLatLng(latlng);
				hintCircle.setRadius(hintRadiusM);
			}
		} else if (hintCircle) {
			hintCircle.remove();
			hintCircle = undefined;
		}
	}

	function fitMapView() {
		if (!map) return;
		const layers: L.Circle[] = [];
		if (radiusCircle) layers.push(radiusCircle);
		if (hintCircle) layers.push(hintCircle);

		if (layers.length > 0) {
			const group = L.featureGroup(layers);
			map.fitBounds(group.getBounds().pad(0.2));
		} else if (hasValidCoords(latitude, longitude)) {
			map.setView([latitude, longitude], 15);
		}
	}

	$effect(() => {
		map = L.map(mapContainer, {
			center: hasValidCoords(latitude, longitude) ? [latitude, longitude] : [-36.85, 174.76],
			zoom: hasValidCoords(latitude, longitude) ? 15 : 5
		});

		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			attribution: '&copy; OpenStreetMap contributors',
			maxZoom: 19
		}).addTo(map);

		map.on('click', (e: L.LeafletMouseEvent) => {
			latitude = Math.round(e.latlng.lat * 1e6) / 1e6;
			longitude = Math.round(e.latlng.lng * 1e6) / 1e6;
			updateTargetMarker();
		});

		if (hasValidCoords(latitude, longitude)) {
			updateTargetMarker();
		}

		updateHintLayers();

		// Leaflet needs invalidateSize when initialized in a hidden/resizing container (e.g. dialog).
		// Fit the view after so bounds are calculated with the correct dimensions.
		setTimeout(() => {
			map?.invalidateSize();
			fitMapView();
		}, 100);

		return () => {
			map?.remove();
			map = undefined;
			marker = undefined;
			radiusCircle = undefined;
			hintMarker = undefined;
			hintCircle = undefined;
		};
	});

	// React to radius changes
	$effect(() => {
		void radiusM;
		updateRadiusCircle();
	});

	// React to hint toggle and hint params
	$effect(() => {
		void showHint;
		void hintLatitude;
		void hintLongitude;
		void hintRadiusM;
		updateHintLayers();
	});
</script>

<div class="map-wrapper">
	<div class="map-container" bind:this={mapContainer}></div>
	<div class="map-legend">
		<span class="legend-item legend-item--target">
			<span class="legend-dot"></span> Target location
		</span>
		{#if showHint}
			<span class="legend-item legend-item--hint">
				<span class="legend-dot"></span> Hint area (shown to user)
			</span>
		{/if}
	</div>
</div>

<style>
	.map-wrapper {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.map-container {
		height: 350px;
		border-radius: 8px;
		border: 1px solid #d1d5db;
		overflow: hidden;
	}

	.map-legend {
		display: flex;
		gap: 1rem;
		font-size: 0.75rem;
		color: #6b7280;
	}

	.legend-item {
		display: flex;
		align-items: center;
		gap: 0.35rem;
	}

	.legend-dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		display: inline-block;
	}

	.legend-item--target .legend-dot {
		background: #ec4899;
	}

	.legend-item--hint .legend-dot {
		background: #8b5cf6;
	}

	:global(.map-picker-pin) {
		width: 16px;
		height: 16px;
		border-radius: 50%;
		background: #ec4899;
		border: 2px solid white;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	}

	:global(.map-picker-pin--hint) {
		background: #8b5cf6;
	}

	:global(.map-picker-icon) {
		background: none !important;
		border: none !important;
	}
</style>
