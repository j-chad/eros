/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />
/// <reference types="../.svelte-kit/ambient.d.ts" />

import type {PushMessage} from "./lib/types/push";

declare const self: ServiceWorkerGlobalScope;

import { build, files, prerendered, version } from '$service-worker';

const CACHE_NAME = `cache-${version}`;

const PRECACHE_ASSETS = [...build, ...files, ...prerendered];
const APP_SHELL = '/__app_shell';

self.addEventListener('install', (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then(async (cache) => {
			// Cache each asset individually so a single failure doesn't block installation
			const results = await Promise.allSettled(PRECACHE_ASSETS.map((url) => cache.add(url)));

			const failed = results
				.map((r, i) => (r.status === 'rejected' ? PRECACHE_ASSETS[i] : null))
				.filter(Boolean);

			if (failed.length > 0) {
				console.warn('[SW] Failed to precache:', failed);
			}

			// Cache the SPA shell for offline navigation fallback.
			// We fetch / which returns the processed 200.html (with absolute paths and
			// modulepreload hints injected by the server), then store it under a synthetic
			// key so it doesn't collide with normal route caching.
			const shell = await fetch('/');
			if (shell.ok) {
				await cache.put(APP_SHELL, shell);
			}

			await self.skipWaiting();
		})
	);
});

self.addEventListener('activate', (event) => {
	event.waitUntil(
		caches
			.keys()
			.then((keys) =>
				Promise.all(keys.filter((key) => key !== CACHE_NAME).map((key) => caches.delete(key)))
			)
			.then(() => self.clients.claim())
	);
});

self.addEventListener('fetch', (event) => {
	const { request } = event;
	const url = new URL(request.url);

	// Only handle same-origin GET requests
	if (request.method !== 'GET') return;
	if (url.origin !== self.location.origin) return;

	// Never cache API calls -- the app uses IndexedDB for offline data
	if (url.pathname.startsWith('/api/')) return;

	event.respondWith(
		caches.match(request).then((cached) => {
			if (cached) return cached;

		// Navigation requests that miss the cache get the SPA fallback,
		// so the client-side router can handle the URL.
		// Exception: /_app/ paths are static build assets (JS, source maps, etc.)
		// that should be fetched from the network, not replaced with the SPA shell.
		if (request.mode === 'navigate' && !url.pathname.startsWith('/_app/')) {
			return caches.match(APP_SHELL).then((fallback) => fallback || fetch(request));
		}

			return fetch(request);
		})
	);
});

self.addEventListener('push', async (event) => {
	const {title, ...options} = event.data.json() as PushMessage
	await self.registration.showNotification(title, options)
})

self.addEventListener('notificationclick', (event) => {
	event.notification.close()

	if (event.notification.data?.url) {
		event.waitUntil(
			self.clients.matchAll({ type: 'window' }).then((clients) => {
				for (const client of clients) {
					if ('focus' in client) {
						client.focus();
						client.navigate(event.notification.data.url);
						return;
					}
				}

				self.clients.openWindow?.(event.notification.data.url);
			})
		);
	}
})