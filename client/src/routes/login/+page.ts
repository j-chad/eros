import { isAuthenticated } from '$lib/services/auth';
import { safeReturnTo } from '$lib/utils/url';
import { redirect } from '@sveltejs/kit';

export async function load({ url }) {
	if (await isAuthenticated()) {
		throw redirect(307, safeReturnTo(url.searchParams.get('returnTo')));
	}

	return {};
}
