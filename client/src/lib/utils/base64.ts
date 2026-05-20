export function base64URLDecode(input: string): string {
    const padding = '='.repeat((4 - input.length % 4) % 4);
    const base64 = (input + padding).replace(/-/g, '+').replace(/_/g, '/');

    return atob(base64);
}