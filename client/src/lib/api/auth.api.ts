import { rawRequest, request, ServerUnreachableError } from '$lib/api/http';

export async function isSessionValid(): Promise<boolean> {
	try {
		const response = await rawRequest('/ping', { method: 'GET' });
		return response.ok;
	} catch (e) {
		// Let ServerUnreachableError propagate so the auth service can
		// distinguish "server down" from "token invalid".
		if (e instanceof ServerUnreachableError) {
			throw e;
		}
		return false;
	}
}

export async function loginRequest(registration_code: string, device_info: string): Promise<string> {
	const response = await request<{ token: string }>(
		'/device',
		{
			method: 'POST',
			body: JSON.stringify({ registration_code, device_info }),
		},
		false,
	);
	return response.token;
}
