<script lang="ts">
	import { dev } from '$app/environment'
	import { onMount } from 'svelte';
	import { onNavigate } from '$app/navigation';
	import { PUBLIC_SERVICE_WORKER, PUBLIC_S3_ORIGIN, PUBLIC_MAINTENANCE_MODE } from '$env/static/public';
	import MaintenanceBanner from '$lib/ui/MaintenanceBanner.svelte';
	import { startSync } from '$lib/db/sync';

	import "../app.css";

	let { children } = $props();

	const maintenanceMode = PUBLIC_MAINTENANCE_MODE === 'true';

	const s3Origin = PUBLIC_S3_ORIGIN ?? '';
	const s3Csp = s3Origin ? ` ${s3Origin}` : '';
	const csp = `default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https://*.tile.openstreetmap.org${s3Csp}; connect-src 'self' http://localhost:* https://*; media-src 'self' blob:${s3Csp}; font-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'`;

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});

	onMount(() => {
		startSync();
		navigator.storage?.persist();
	});

	if (PUBLIC_SERVICE_WORKER == "true" && 'serviceWorker' in navigator) {
		const prevController = navigator.serviceWorker.controller;
		navigator.serviceWorker.addEventListener('controllerchange', () => {
			if (prevController) window.location.reload();
		});

		if (document.readyState === 'complete') {
			registerSW();
		} else {
			addEventListener('load', registerSW);
		}
	}

	function registerSW() {
		navigator.serviceWorker.register('/service-worker.js', {
			type: dev ? 'module' : 'classic'
		});
	}
</script>

<svelte:head>
	<meta http-equiv="content-security-policy" content={csp} />
</svelte:head>

{#if maintenanceMode}
	<MaintenanceBanner />
{/if}

{@render children()}
