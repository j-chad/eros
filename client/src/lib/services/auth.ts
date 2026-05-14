import { loginRequest, isSessionValid } from '$lib/api/auth.api';
import { setToken, loadToken, clearToken } from '$lib/api/auth';
import { ServerUnreachableError } from '$lib/api/http';
import { isOnline } from '$lib/online.svelte';
import { invalidate } from '$app/navigation';

export const AUTH_DEPENDENCY = 'auth:session';

export async function isAuthenticated(): Promise<boolean> {
	const token = await loadToken();
	if (!token) {
		console.warn('No auth token found, user is not authenticated.');
		return false;
	}

	if (!isOnline()) {
		return true; // trust cached token when server is unreachable
	}

	try {
		if (!(await isSessionValid())) {
			console.warn('Auth session is not valid, user is not authenticated.');
			return false;
		}
	} catch (e) {
		if (e instanceof ServerUnreachableError) {
			// Server went down between the isOnline() check and the ping.
			// Trust the cached token rather than logging the user out.
			return true;
		}
		throw e;
	}

	return true;
}

export async function login(registration_code: string, device_info: string): Promise<void> {
	const token = await loginRequest(registration_code, device_info);
	await setToken(token);
	await invalidate(AUTH_DEPENDENCY);
}

export async function logout(): Promise<void> {
	await clearToken();
}
