/** Returns the path if it's a safe same-origin relative URL, otherwise fallback. */
export function safeReturnTo(value: string | null, fallback = '/'): string {
	if (!value) return fallback;
	if (value.startsWith('/') && !value.startsWith('//')) return value;
	return fallback;
}
