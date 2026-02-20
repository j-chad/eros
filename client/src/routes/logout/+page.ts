import {redirect} from "@sveltejs/kit";
import {logout} from "$lib/services/auth";
import {browser} from "$app/environment";

export async function load() {
	if (!browser) return

	await logout()
	return redirect(307, '/');
}
