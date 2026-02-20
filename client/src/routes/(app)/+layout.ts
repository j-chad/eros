import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import {get} from "svelte/store";
import {isAuthenticated} from "$lib/services/auth";
import {authToken} from "$lib/api/auth";

export const ssr = false;

export async function load({ url }) {
	if (!browser) return;

	if (!await isAuthenticated()) {
		const returnTo = encodeURIComponent(url.pathname + url.search);
		throw redirect(307, `/login?returnTo=${returnTo}`);
	}

	return {};
}
