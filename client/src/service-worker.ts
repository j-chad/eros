/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />
/// <reference types="../.svelte-kit/ambient.d.ts" />

declare const self: ServiceWorkerGlobalScope;

import { build, files, prerendered, version } from '$service-worker'
import {cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute, type PrecacheEntry} from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'

cleanupOutdatedCaches();

const precache_list: PrecacheEntry[] = [...build, ...files, ...prerendered].map((s) => ({
	url: s,
	revision: version,
}))

const isProduction = build.length > 0;

// Precache the SPA fallback page (only exists in production builds)
if (isProduction) {
	precache_list.push({ url: '/200.html', revision: version });
}

precacheAndRoute(precache_list, {
	// Ignore all URL search params so /login?returnTo=... matches /login
	ignoreURLParametersMatching: [/./],
})

// For navigation requests that don't match precache, serve the SPA fallback
if (isProduction) {
	const navigationRoute = new NavigationRoute(
		createHandlerBoundToURL('/200.html'),
		{ denylist: [/^\/api\//] }
	);
	registerRoute(navigationRoute);
}

self.addEventListener('install', (event) => {
	event.waitUntil(self.skipWaiting());
});

self.addEventListener('activate', (event) => {
	event.waitUntil(self.clients.claim());
});
