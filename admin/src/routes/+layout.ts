import {auth} from "$lib/auth.svelte";
import {redirect} from "@sveltejs/kit";
import {browser} from "$app/environment";

export async function load({ url }) {
    const isAuthenticated = auth.isAuthenticated;
    if (browser && !isAuthenticated && url.pathname !== '/login') {
        throw redirect(303, '/login');
    }

    return {};
}