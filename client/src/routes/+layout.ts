import {initAuth} from "$lib/api/auth";

export async function load() {
	await initAuth();
	console.log('Auth initialized');
}
