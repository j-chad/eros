import {browser} from '$app/environment';
import {KVKey, KVStore} from "$lib/db/stores/kv";

export async function isAuthenticated(): Promise<boolean> {
	const token = await getToken();
	return !!token;
}

export async function getToken(): Promise<string | null> {
	if (!browser) return null;
	return await KVStore.get(KVKey.AuthSession);
}

export async function setToken(token: string): Promise<void> {
	if (!browser) return;
	await KVStore.set(KVKey.AuthSession, token);
}

export async function clearToken(): Promise<void> {
	if (!browser) return;
	await KVStore.delete(KVKey.AuthSession);
}
