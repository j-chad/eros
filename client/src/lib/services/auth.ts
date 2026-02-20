import { login as loginRequest } from "$lib/repositories/auth";
import {setToken} from "$lib/api/auth";
import { invalidate } from "$app/navigation";

export const AUTH_DEPENDENCY = 'auth:session';

export async function login(registration_code: string, device_info: string): Promise<void> {
	const token = await loginRequest(registration_code, device_info);
	await setToken(token);
	await invalidate(AUTH_DEPENDENCY)
}
