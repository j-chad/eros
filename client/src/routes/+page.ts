import {error, redirect} from "@sveltejs/kit";

export async function load() {
	return redirect(307, '/app');
}
