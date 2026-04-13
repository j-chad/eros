# TODO

> Priority: **P0** = critical path (blocks core features), **P1** = important, **P2** = nice to have

## Graph Play Experience

The core gameplay loop. The client can list graphs, show a countdown to the next one, and display a calendar of unlocked graphs — but cannot navigate or play through them. See `.opencode/plans/gate-unlock-flows/spec.md` for the full spec.

- [ ] Backend: unlock endpoint changes #P0
  - [ ] return `UnlockResult` (unlocked node + newly accessible nodes/edges) instead of 204
  - [ ] case-insensitive code matching (`strings.EqualFold`)
  - [ ] per-node rate limiting middleware (10/min) on `POST /api/nodes/{id}/unlock`
- [ ] Client: API + service layer #P0
  - [ ] add `unlockNode()` to `client/src/lib/api/graph.api.ts` with `UnlockResult` type
  - [ ] add `unlockNode()` to `client/src/lib/services/graph.ts` — call API + merge result into IndexedDB cache
  - [ ] shared reactive online status utility (`client/src/lib/online.svelte.ts`)
- [ ] Client: gate unlock flows #P0
  - [ ] Code gate — input, submit, shake + error on wrong answer, rate-limited state, offline-aware
  - [ ] Location gate — geolocation request, permission denial handling, poor accuracy warning, offline-aware
  - [ ] Manual gate — 30s polling (visibility-aware), approval detection, "Continue" button, offline-aware
- [ ] Client: unlock celebration + transition #P0
  - [ ] standard gate unlock — checkmark + pulse + auto-advance (~1.5s)
  - [ ] reward unlock — elevated entrance (scale, shimmer, glow), longer beat (~2s)
  - [ ] confetti/particle flourish on reward unlock (can slip to P1 if complex)

## File Management

Reward nodes can contain images, videos, and files. The backend can store them, but can't serve them to the client, and the admin has no upload UI.

- [ ] Backend #P0
  - [ ] implement `FileService` methods in `internal/service/files.go` — struct and constructor exist but zero methods
  - [ ] add client-facing file download/streaming endpoint
- [ ] Admin #P1
  - [ ] add file-related methods to `admin/src/lib/api.ts` (upload, list, delete)
  - [ ] build file upload component in the reward node edit dialog (`admin/src/routes/graphs/[id]/edit-node-dialog/`)
  - [ ] display and manage uploaded files for reward nodes in the editor
- [ ] Client #P0
  - [ ] image/video reward viewer
  - [ ] other reward type displays (markdown, favour, calendar, wallet, file)

## Favour System

Backend is fully implemented (both admin and client endpoints). Admin UI is done. Client has nothing.

- [ ] Client: data layer #P1
  - [ ] Create `client/src/lib/api/favour.api.ts`
  - [ ] Create favour DB store in `client/src/lib/db/stores/`
  - [ ] Create favour service in `client/src/lib/services/`
- [ ] Client: UI #P1
  - [ ] Favour browsing page — view available choices and remaining favour count
  - [ ] Favour request flow — select a favour and submit a request
  - [ ] Request history / status page — view pending, approved, and rejected requests


## Tech Debt and Nice-to-Haves (P2)
- [ ] Backend: registration code validation bypassed — commented out in `internal/service/auth.go`
- [ ] Backend: `CleanupOrphanedFiles` commented out in `internal/service/admin.go`
- [ ] Admin: impersonate button renders but `handleImpersonate` handler is empty (`admin/src/lib/components/Impersonate.svelte`)
- [ ] Client: `src/lib/db/sync.ts` is empty — background sync not implemented
- [ ] Event system (WebSocket/SSE) — replace manual gate polling with real-time approval notifications
- [ ] Location gate: cold-warm-hot continuous proximity feedback (watch position, show distance hint)
- [ ] Code gate: alternate accepted codes — multiple valid answers per code gate node
- [ ] Backend: move rate limiter to Redis/shared store if backend ever runs multiple instances
