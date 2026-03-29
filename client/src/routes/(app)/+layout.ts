import { redirect } from '@sveltejs/kit';
import { AUTH_DEPENDENCY, isAuthenticated } from '$lib/services/auth';

export async function load({ url, depends }) {
	depends(AUTH_DEPENDENCY)

	if (!await isAuthenticated()) {
		const returnTo = encodeURIComponent(url.pathname + url.search);
		throw redirect(307, `/login?returnTo=${returnTo}`);
	}

	return {};
}
