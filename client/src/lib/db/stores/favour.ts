import { db, promisifyRequest } from '$lib/db/db';
import type { FavourChoice, FavourCount, FavourRequest } from '$lib/types/favour';
import { KVStore, KVKey } from './kv';

export async function getAllFavourChoices(): Promise<FavourChoice[]> {
	const store = await db.getStore('favourChoices', 'readonly');
	return promisifyRequest(store.getAll());
}

export async function putFavourChoices(choices: FavourChoice[]): Promise<void> {
	const tx = await db.transaction('favourChoices', 'readwrite');
	const store = tx.objectStore('favourChoices');
	// Clear stale entries and replace with the latest set.
	await promisifyRequest(store.clear());
	await Promise.all(choices.map((c) => promisifyRequest(store.put(c))));
}

export async function getAllFavourRequests(): Promise<FavourRequest[]> {
	const store = await db.getStore('favourRequests', 'readonly');
	return promisifyRequest(store.getAll());
}

export async function putFavourRequests(requests: FavourRequest[]): Promise<void> {
	const tx = await db.transaction('favourRequests', 'readwrite');
	const store = tx.objectStore('favourRequests');
	await promisifyRequest(store.clear());
	await Promise.all(requests.map((r) => promisifyRequest(store.put(r))));
}

export async function getFavourCount(): Promise<FavourCount | null> {
	return KVStore.get(KVKey.FavourCount);
}

export async function putFavourCount(count: FavourCount): Promise<void> {
	await KVStore.set(KVKey.FavourCount, count);
}
