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

    const [registration, devices] = await Promise.all([
        api.registration.get().catch((err) => {
            if (err.status === 404) {
                return null;
            }

            throw err;
        }),
        api.devices.list()
    ]);

    return {registration, devices};
}