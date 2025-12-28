<script lang="ts">
	import favicon96 from '$lib/assets/favicon-96x96.png';
	import faviconSVG from '$lib/assets/favicon.svg';
	import faviconICO from '$lib/assets/favicon.ico';
	import appleTouchIcon from '$lib/assets/apple-touch-icon.png';
	import manifest from '$lib/assets/site.webmanifest';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" type="image/png" href={favicon96} sizes="96x96" />
	<link rel="icon" type="image/svg+xml" href={faviconSVG} />
	<link rel="shortcut icon" href={faviconICO} />
	<link rel="apple-touch-icon" sizes="180x180" href={appleTouchIcon} />
	<meta name="apple-mobile-web-app-title" content="Eros Admin" />
	<link rel="manifest" href={manifest} />
	<title>Eros Admin</title>
</svelte:head>

{@render children()}
