// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
interface ApiError {
	message: string;
	base: string;
	endpoint: string;
	method: string;
	body: unknown;
}

declare global {
	namespace App {
		interface Error {
			message: string;
			base?: string;
			endpoint?: string;
			method?: string;
			body?: unknown;
		}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
