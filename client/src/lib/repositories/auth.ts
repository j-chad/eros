import {request} from "$lib/api/http";


export function login(registration_code: string, device_info: string): Promise<string> {
	return request('/device', {
		method: 'POST',
		body: JSON.stringify({
			registration_code,
			device_info,
		}),
	});
}
