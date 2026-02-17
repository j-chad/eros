import { login as loginRequest } from "$lib/repositories/auth";
import {setToken} from "$lib/api/auth";

export async function login(registration_code: string, device_info: string): Promise<void> {
	const token = await loginRequest(registration_code, device_info);
	return setToken(token)
}
