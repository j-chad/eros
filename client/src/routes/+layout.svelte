<script lang="ts">
	import { dev } from '$app/environment'
	import { onNavigate } from '$app/navigation';
	import { PUBLIC_SERVICE_WORKER, PUBLIC_S3_ORIGIN, PUBLIC_MAINTENANCE_MODE } from '$env/static/public';
	import MaintenanceBanner from '$lib/ui/MaintenanceBanner.svelte';

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

	if (PUBLIC_SERVICE_WORKER == "true" && 'serviceWorker' in navigator) {
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
