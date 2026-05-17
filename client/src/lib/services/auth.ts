import { loginRequest } from '$lib/api/auth.api';
import { setToken, loadToken, clearToken } from '$lib/api/auth';
import { invalidate } from '$app/navigation';

export const AUTH_DEPENDENCY = 'auth:session';

/**
 * Optimistic, local-only auth check. If a token exists in IndexedDB, the user
 * is considered authenticated immediately — no network required. Token validity
 * is verified lazily: the HTTP layer handles 401 responses by clearing the
 * token and redirecting to /login.
 */
export async function isAuthenticated(): Promise<boolean> {
	const token = await loadToken();
	if (!token) {
		return false;
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
