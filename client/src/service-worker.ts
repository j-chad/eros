/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />
/// <reference types="../.svelte-kit/ambient.d.ts" />

declare const self: ServiceWorkerGlobalScope;

import { build, files, prerendered, version } from '$service-worker'
import {cleanupOutdatedCaches, precacheAndRoute, type PrecacheEntry} from 'workbox-precaching'

cleanupOutdatedCaches();

const precache_list: PrecacheEntry[] = [...build, ...files, ...prerendered].map((s) => ({
	url: s,
	revision: version,
}))

precacheAndRoute(precache_list)

self.addEventListener('install', (event) => {
	self.skipWaiting().catch(() => console.error('Failed to skip waiting during service worker installation'));
});

self.addEventListener('activate', (event) => {
	event.waitUntil((async () => {
		await self.clients.claim().catch(() => console.error('Failed to claim clients during service worker activation'));
	})());
});
