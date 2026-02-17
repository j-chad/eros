import {writable} from "svelte/store";
import {getToken, setToken as setTokenDB, clearToken as clearTokenDB} from "$lib/db/services/auth";

export const authToken = writable<string | null>(null);

export async function initAuth(){
	const token = await getToken();
	authToken.set(token);
}

export async function setToken(token: string) {
	await setTokenDB(token)
	authToken.set(token);
}

export async function clearToken() {
	await clearTokenDB();
	authToken.set(null);
}

