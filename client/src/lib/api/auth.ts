import {derived, writable} from "svelte/store";
import {KVKey, KVStore} from "$lib/db/stores/kv";

export const authToken = writable<string | null>(null);

export const isAuthenticated = derived(authToken, ($authToken) => !!$authToken);

export async function initAuth(){
	const token = await KVStore.get(KVKey.AuthSession);
	authToken.set(token);
}

export async function setToken(token: string) {
	await KVStore.set(KVKey.AuthSession, token);
	authToken.set(token);
}

export async function clearToken() {
	await KVStore.delete(KVKey.AuthSession);
	authToken.set(null);
}

