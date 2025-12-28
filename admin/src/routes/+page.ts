import {api} from "$lib/api";
import {browser} from "$app/environment";
import type {Device} from "$lib/types";

export async function load({fetch}) {
    if (!browser) {
        return {
            registration: null,
            devices: []
        };
    }

    const [registration] = await Promise.all([
        api.registration.get().catch((err) => {
            if (err.status === 404) {
                return null;
            }

            throw err;
        }),
    ]);

    return {registration, devices: [] as Device[]};
}