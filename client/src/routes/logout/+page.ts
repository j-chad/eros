import {redirect} from "@sveltejs/kit";
import {logout} from "$lib/services/auth";
import {browser} from "$app/environment";

export async function load() {
	if (!browser) return

	await logout()
	throw redirect(307, '/');
}
