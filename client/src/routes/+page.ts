import {redirect} from "@sveltejs/kit";
import {isAuthenticated} from "$lib/services/auth";
import {browser} from "$app/environment";

export async function load() {
	return redirect(307, '/test');
}
