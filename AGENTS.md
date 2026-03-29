# AGENTS.md

## Project Overview

Eros is a full-stack couples app — a graph-based scavenger/treasure hunt system. An admin creates interactive graphs with gate nodes (location, code, manual) and reward nodes (images, videos, favours, files). A client navigates them via a PWA.

**Architecture:** Go backend + two SvelteKit frontends (admin + client) in a monorepo.

**Current roadmap and outstanding work:** See [TODO.md](./TODO.md) for the prioritised task list. Read this file at the start of a session and make sure it is updated as work progresses.

## Build & Run Commands

### Backend (Go)

```bash
# Run the server (from backend/)
go run ./cmd/server

# Build
go build ./cmd/server

# Vet / static analysis
go vet ./...
```

Requires CGo (sqlite3 driver). Go 1.26+. Config loaded from `backend/config.*.json` files; secrets in `config.private.json` (gitignored).

### Admin Frontend (SvelteKit + adapter-static)

```bash
# Install dependencies (from admin/)
npm install

# Dev server
npm run dev

# Build (static site)
npm run build

# Lint
npm run lint

# Format check
npm run format
```

Node 24 (see `.nvmrc`). Uses npm.

### Client Frontend (SvelteKit + adapter-auto)

```bash
# Install dependencies (from client/)
npm install

# Dev server (port 5174)
npm run dev

# Build
npm run build

# Lint
npm run lint

# Format check
npm run format
```

Node 24 (see `.nvmrc`). Uses npm. Dev server runs on port 5174 to avoid service worker conflicts.

### Tests

No test framework is configured. No test files exist in the codebase.

### API Testing

Bruno collections in `bruno/` directory for manual API testing against localhost:8080.

## Tech Stack Details

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.26, stdlib `net/http`, SQLite via `go-sqlite3` (CGo) |
| Admin frontend | SvelteKit (Svelte 5), TypeScript, Vite 7, hand-written scoped CSS |
| Client frontend | SvelteKit (Svelte 5), TypeScript, Vite 7, Tailwind CSS 4, DaisyUI 5 |
| Icons | lucide-svelte (both frontends) |
| Graph editor | @xyflow/svelte (admin only) |

## Code Style — Go Backend

### Imports

Single `import()` block, alphabetically sorted. No consistent blank-line grouping between stdlib and internal packages. Module name is `backend` (not a full URL path). Only external dependency is `github.com/mattn/go-sqlite3`.

### Error Handling

- Repository layer: wrap with `fmt.Errorf("failed to <action>: %w", err)`
- Service layer: sentinel errors as `var Err... = errors.New(...)`, checked with `errors.Is()`
- Handler layer: return `apierror.APIError` via `response.Error(r.Context(), w, apiErr)` then `return`
- Startup: `log.Fatalf`
- Custom `pkg/apierror` type with constructors: `BadRequest()`, `NotFound()`, `Unauthorized()`, `Forbidden()`, `UnknownInternalError(err)` — the last auto-unwraps existing `*APIError`
- Fluent detail builders: `.WithDetail(key, val)`, `.WithDetails(map)`

### Naming

- **Packages:** lowercase single words (`apierror`, `authctx`, `middleware`)
- **Files:** lowercase with underscores (`registration_code.go`, `node_data.go`)
- **Exported types/funcs:** PascalCase. Constructors: `New<Type>()`
- **Receivers:** short, 1-2 chars (`s`, `h`, `n`, `d`, `e`)
- **Variables:** short contextual names (`ctx`, `err`, `req`, `mux`, `conf`)
- **Constants:** PascalCase exported, camelCase unexported. Error code strings use SCREAMING_SNAKE (`"BAD_REQUEST"`)

### Struct Patterns

- Models in `internal/models/` with `json:"snake_case"` tags
- Sensitive fields tagged `json:"-"`
- Optional fields use pointers (`*string`, `*time.Time`) with selective `omitempty`
- Config structs use custom `env:` and `required:` tags
- Validation methods return `*apierror.APIError`

### Handler Patterns

- Signature: `func (h *Handler) Method(w http.ResponseWriter, r *http.Request)`
- Path params via `r.PathValue("id")` (Go 1.22+ routing)
- JSON body: `json.NewDecoder(r.Body).Decode(&target)`
- Responses via `pkg/response`: `response.JSON(w, status, data)`, `response.Error(r.Context(), w, apiErr)`, `response.NoContent(w)`
- Route registration with method-prefixed patterns: `mux.HandleFunc("GET /api/admin/graphs/{id}", h.admin.GetGraph)`

### Architecture

- Clean layers: `handler` → `service` → `repository` (interface) → `sqlite` (impl)
- Manual constructor injection in `main.go`, no DI framework
- Services depend on `repository.Repository` interface; handlers depend on concrete service pointers
- Transactions via `repo.WithTx()` pattern
- Context propagated through all layers; auth context via `authctx` package
- File storage abstracted behind `storage.FileStore` interface (local or S3)
- Structured logging with `log/slog`

### Modern Go Features Used

- Generics (`DecodeInto[T any]`)
- `iter.Seq2` iterators for file listing
- `go:embed` for SQL init scripts and config files
- `slices` package
- `min` builtin

## Code Style — SvelteKit Frontends

### Svelte 5 Runes

Both frontends use Svelte 5 runes consistently:
- `$props()` for component props (admin uses `let` destructuring, client `$lib/ui/` uses `const`)
- `$state()` for mutable local state; `$state.raw()` for high-frequency non-deep-reactive state
- `$derived()` for computed values; `$derived(() => ...)` closure form for complex derivations
- `$effect()` for side effects
- `$bindable()` for two-way-bindable props
- `{#snippet}` / `{@render}` for composition (replaces slots)
- `.svelte.ts` extension for files with rune syntax outside components

### TypeScript

- `<script lang="ts">` on all components
- `strict: true` in tsconfig
- `interface` for object shapes, `type` for unions/aliases/utility types
- `enum` with string values for discriminated unions (`NodeType`, `RewardType`)
- `import type` for type-only imports; inline `type` keyword for mixed imports
- `$lib` alias for all `src/lib/` imports; relative imports for route-local components

### Formatting (Prettier)

Config: tabs, single quotes, no trailing commas, 100 char print width, `prettier-plugin-svelte`. Note: formatting is not consistently applied — expect mixed quote styles and semicolon usage.

### ESLint

Flat config: `@eslint/js` recommended + `typescript-eslint` recommended + `eslint-plugin-svelte` recommended + `eslint-config-prettier`. `no-undef: 'off'`.

### Naming Conventions

- **Component files:** PascalCase (`Button.svelte`, `GraphRenderer.svelte`)
- **TS modules:** camelCase or dot-separated (`auth.ts`, `graph.types.ts`, `auth.svelte.ts`)
- **Route dirs:** kebab-case or `[param]`; route groups in parens `(app)/`
- **State variables:** camelCase; booleans prefixed `is`/`has`/`can`/`show`
- **Event handlers:** `handle` prefix (`handleSave`, `handleDeleteDevice`)
- **Callback props:** `on` prefix (`onSave`, `onClose`); native DOM: lowercase (`onclick`)
- **Constants:** UPPER_SNAKE (`API_BASE`, `DB_NAME`)

### API Interaction

- Admin: centralized `api` object in `lib/api.ts` with namespaced methods (`api.graph.list()`)
- Client: layered architecture — see below
- Both use generic `request<T>()` with `PUBLIC_SERVER_URL` env var (no `/api` suffix — all endpoint strings include `/api/...`)
- Admin auth: `Authorization: Admin <key>`, Client auth: `Authorization: Bearer <token>`
- Data loading in `+page.ts` with browser guard (`if (!browser) return defaults`)

### Client Data Layer

The client uses a strict layered architecture for offline-first data access:

| Layer | Path | Responsibility |
|---|---|---|
| **Types** | `src/lib/types/` | TypeScript interfaces and enums only. No logic. |
| **API** | `src/lib/api/*.api.ts` | Raw HTTP calls via `request()`/`rawRequest()`. One file per resource (e.g. `graph.api.ts`, `auth.api.ts`). No storage, no business logic. |
| **DB stores** | `src/lib/db/stores/` | IndexedDB read/write for a single object store. One file per store (e.g. `graph.ts`, `kv.ts`). No network calls. |
| **Services** | `src/lib/services/` | Orchestration: if online → fetch from API and write to IndexedDB; if offline → read from IndexedDB. This is what pages call. |

`src/lib/db/db.ts` — typed `Database` wrapper around the raw IndexedDB API.
`src/lib/db/schema.ts` — `DBSchema` interface and `createObjectStores()` for upgrades. `StoredGraph` and similar stored types live here (dates as ISO strings, matching what the API returns).
`src/lib/db/index.ts` — barrel re-export for the db layer.
`src/lib/api/auth.ts` — `authToken` Svelte store + `loadToken`/`setToken`/`clearToken` backed by `KVStore`. Not an `.api.ts` file because it manages client-side token state, not just HTTP.

**Auth token hydration:** `api/http.ts` reads `authToken` from the store, falling back to `loadToken()` from IndexedDB if the store is empty. Pages inside `(app)/` can assume the token is available without any extra setup — `rawRequest` handles it transparently.

### Styling

- **Admin:** hand-written scoped `<style>` blocks, hardcoded hex colors, `:global()` selectors, flexbox/grid
- **Client:** Tailwind CSS v4 utility classes + DaisyUI v5 component classes (valentine theme), minimal scoped CSS. See `client/DESIGN.md` for the full design language reference — colours, typography, spacing, radius, shadows, components, motion, and conventions.
- **No emojis in UI copy.** They cheapen the feel.

### Error Handling

- API errors throw SvelteKit `error()` with structured `App.Error` body
- Event handlers: `try/catch` with `console.error` + `alert()` (admin) or `$state` error variables (client)
- Load functions: `.catch()` for expected failures (e.g., 404 → null)

### State Management

- Admin: class-based singleton with `$state` rune in `.svelte.ts` file
- Client: mix of Svelte 5 `$state` and Svelte 4 `writable` stores backed by IndexedDB
- Page data copied into `$state` for mutation: `let items = $state(data.items)`
