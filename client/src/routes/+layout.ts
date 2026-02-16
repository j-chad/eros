import {initAuth, setToken} from "$lib/api/auth";

export async function load() {
	await initAuth();
}
