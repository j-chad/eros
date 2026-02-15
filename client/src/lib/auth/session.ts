import {browser} from '$app/environment';
import {KVKey, KVStore} from "$lib/db/stores/kv";

export async function isAuthenticated(): Promise<boolean> {
	const token = await readCachedToken();
	return !!token;
}

export async function readCachedToken(): Promise<string | null> {
	if (!browser) return null;
	const cached = await KVStore.get(KVKey.AuthSession);
	if (!cached) return null;

	const now = Math.floor(Date.now() / 1000);
	if (cached.expiresAt && cached.expiresAt <= now) {
		await clearCachedToken()
		return null;
	}

	return cached.token;
}

export async function cacheToken(token: string, expiresAt: number): Promise<void> {
	if (!browser) return;
	await KVStore.set(KVKey.AuthSession, { token, expiresAt });
}

export async function clearCachedToken(): Promise<void> {
	if (!browser) return;
	await KVStore.delete(KVKey.AuthSession);
}
