<script lang="ts">
	import { dev } from '$app/environment'
	import { PUBLIC_SERVICE_WORKER } from '$env/static/public';

	import "../app.css";

	let { children } = $props();

	if (PUBLIC_SERVICE_WORKER === 'true' && 'serviceWorker' in navigator) {
		addEventListener('load', function () {
			navigator.serviceWorker.register('./path/to/service-worker.js', {
				type: dev ? 'module' : 'classic'
			});
		});
	}
</script>

{@render children()}
