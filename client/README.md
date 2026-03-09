# Eros — Client App

The partner-facing Progressive Web App (PWA). Users register their device via QR code, then follow their reward journey, unlock nodes, and manage their favour wallet.

## Development

```bash
npm install
npm run dev        # start dev server
npm run dev -- --open  # start and open in browser
```

## Building

```bash
npm run build      # production build
npm run preview    # preview the production build locally
```

> The app uses `@sveltejs/adapter-auto`. For deployment you may need to install a specific [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.

## Other commands

```bash
npm run check      # TypeScript type-check
npm run lint       # ESLint + Prettier check
npm run format     # auto-format
```

## Environment

The app expects the backend API to be reachable. Configure the API base URL via the `PUBLIC_API_BASE` environment variable (or `.env` file) if it differs from the default.
