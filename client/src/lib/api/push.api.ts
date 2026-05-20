import {request} from "$lib/api/http";
import {base64URLDecode} from "$lib/utils/base64";

export async function fetchVapidKey(): Promise<string> {
    const encoded = await request<string>(`/push/vapid-key`);
    return base64URLDecode(encoded);
}