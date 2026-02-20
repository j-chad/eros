import {login as loginRequest, isSessionValid} from "$lib/repositories/auth";
import {setToken, loadToken} from "$lib/api/auth";
import { invalidate } from "$app/navigation";

export const AUTH_DEPENDENCY = 'auth:session';

export async function isAuthenticated(): Promise<boolean> {
	const token = await loadToken();
	if (!token) {
		console.warn('No auth token found, user is not authenticated.');
		return false;
	}

	if (!navigator.onLine) {
		return true;
	}

	if (!await isSessionValid()) {
		console.warn('Auth session is not valid, user is not authenticated.');
		return false;
	}

	return true;
}

export async function login(registration_code: string, device_info: string): Promise<void> {
	const token = await loginRequest(registration_code, device_info);
	await setToken(token);
	await invalidate(AUTH_DEPENDENCY)
}
