import devtoolsJson from 'vite-plugin-devtools-json';
import { sveltekit } from '@sveltejs/kit/vite';
import {defineConfig, loadEnv, type PluginOption} from 'vite';
import tailwindcss from "@tailwindcss/vite";
import fs from 'fs'

export default defineConfig(({mode}) => {
	const env = loadEnv(mode, process.cwd(), '')

	const plugins: PluginOption[] = [
		tailwindcss(),
		sveltekit()
	]

	console.log(mode, env)
	if (mode === 'development') {
		plugins.push(devtoolsJson())
	}

	const proxy = {
		'/api': {
			target: 'http://localhost:8080',
			changeOrigin: true,
		}
	}

	return ({
		server: {
			port: 5174, // Prevent service worker conflicts with the default Vite port (5173)
			allowedHosts: env.ALLOWED_HOSTS ? env.ALLOWED_HOSTS.split(',') : undefined,
			https: env.HTTPS ? {
				key: fs.readFileSync('.cert/key.pem'),
				cert: fs.readFileSync('.cert/cert.pem'),
			} : undefined,
			proxy,
		},
		preview: {
			port: env.HTTPS ? 443 : 80,
			https: env.HTTPS ? {
				key: fs.readFileSync('.cert/key.pem'),
				cert: fs.readFileSync('.cert/cert.pem'),
			} : undefined,
			proxy,
		},
		plugins,
		define: {'process.env.NODE_ENV': JSON.stringify(mode)},
	});
});
