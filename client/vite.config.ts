import devtoolsJson from 'vite-plugin-devtools-json';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
	server: {
		port: 5174 // Prevent service worker conflicts with the default Vite port (5173)
	},
	plugins: [tailwindcss(), sveltekit(), devtoolsJson()],
	define: { 'process.env.NODE_ENV': '"production"' }
});
