<script lang="ts">
	import { dev } from '$app/environment'
	import { PUBLIC_SERVICE_WORKER } from '$env/static/public';

	import "../app.css";

	let { children } = $props();

	if (!!PUBLIC_SERVICE_WORKER && 'serviceWorker' in navigator) {
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

{@render children()}
