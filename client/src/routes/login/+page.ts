import {browser} from "$app/environment";
import {isAuthenticated} from "$lib/services/auth";
import {redirect} from "@sveltejs/kit";

export const prerender = true

export async function load() {
	if (!browser) return;

	if (await isAuthenticated()) {
		const returnURL = new URLSearchParams(window.location.search).get('returnTo') ?? '/';
		throw redirect(307, returnURL);
	}

	return {};
}
