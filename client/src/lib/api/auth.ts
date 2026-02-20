import {writable} from "svelte/store";
import {KVKey, KVStore} from "$lib/db/stores/kv";

export const authToken = writable<string | null>(null);

export async function loadToken(): Promise<string| null> {
	const token = await KVStore.get(KVKey.AuthSession);
	authToken.set(token);
	return token;
}

export async function setToken(token: string) {
	console.log('Setting auth token:', token);
	await KVStore.set(KVKey.AuthSession, token);
	authToken.set(token);
}

export async function clearToken() {
	await KVStore.delete(KVKey.AuthSession);
	authToken.set(null);
}

