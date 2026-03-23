import { redirect } from "@sveltejs/kit";

export const prerender = true;
export const ssr = true;

export async function load() {
	return redirect(307, '/app');
}
