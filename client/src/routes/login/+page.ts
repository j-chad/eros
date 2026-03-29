import { isAuthenticated } from '$lib/services/auth';
import { redirect } from '@sveltejs/kit';

export async function load({ url }) {
	if (await isAuthenticated()) {
		const returnTo = url.searchParams.get('returnTo') ?? '/';
		throw redirect(307, returnTo);
	}

	return {};
}
