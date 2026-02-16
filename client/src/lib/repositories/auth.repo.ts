import {request} from "$lib/api/http";

export default {
	login(registration_code: string, device_info: string): Promise<string> {
		const formData = new FormData()
		formData.append('registration_code', registration_code);
		formData.append('device_info', device_info);

		return request('/auth/login', {
			method: 'POST',
			body: formData,
		});
	}
}
