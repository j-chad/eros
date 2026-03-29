import {redirect} from "@sveltejs/kit";
import {logout} from "$lib/services/auth";

export async function load() {
	await logout()
	throw redirect(307, '/');
}
