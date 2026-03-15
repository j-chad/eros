import {error, redirect} from "@sveltejs/kit";

export async function load() {
	// throw error(401, {
	// 	message: 'Unauthorized: No authentication token found',
	// 	body: null,
	// });
	return redirect(307, '/app');
}
