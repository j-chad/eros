import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import {authToken, isAuthenticated} from '$lib/api/auth';
import {get} from "svelte/store";

export const ssr = false;

export async function load({ url }) {
	if (!browser) return;

	const authed = get(isAuthenticated);
	const token = get(authToken)
	console.log('isAuthenticated', authed, 'token', token);

	if (!get(isAuthenticated)) {
		const returnTo = encodeURIComponent(url.pathname + url.search);
		throw redirect(307, `/login?returnTo=${returnTo}`);
	}

	return {};
}
