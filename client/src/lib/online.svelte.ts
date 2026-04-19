let onlineState = $state(typeof navigator !== 'undefined' ? navigator.onLine : true);

if (typeof window !== 'undefined') {
	window.addEventListener('online', () => {
		onlineState = true;
	});

	window.addEventListener('offline', () => {
		onlineState = false;
	});
}

export function useOnlineStatus(): boolean {
	return onlineState;
}