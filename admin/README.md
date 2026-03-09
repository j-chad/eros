# Eros — Admin Panel

The admin web interface for managing device registration, designing reward graphs, and handling favour requests.

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

## Authentication

All admin API requests are authenticated with the `ADMIN_API_KEY` you configure in the backend. Enter this key in the admin login screen; it is stored in `localStorage` for the session.

## Environment

Configure the backend API URL via the `PUBLIC_API_BASE` environment variable (or `.env` file) if it differs from the default.
