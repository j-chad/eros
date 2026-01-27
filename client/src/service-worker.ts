import { build, files, prerendered, version } from '$service-worker'
import { cleanupOutdatedCaches, precacheAndRoute, type PrecacheEntry } from 'workbox-precaching'

const precache_list: PrecacheEntry[] = [...build, ...files, ...prerendered].map((s) => ({
	url: s,
	revision: version,
}))

precacheAndRoute(precache_list)
