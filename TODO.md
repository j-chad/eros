# TODO

> Priority: **P0** = critical path (blocks core features), **P1** = important, **P2** = nice to have

## Backend

- [ ] **Missing Routes & Fixes (P0)** — these gaps block the client graph play experience
  - [ ] Register `UnlockNode` route in `router.go` — handler exists in `handler/client/graph.go` but no route is registered
  - [ ] Fix `ManualData` missing from `nodeDataDecoders` map in `internal/models/graph.go`
- [ ] **File Serving (P0)** — the client PWA needs to fetch reward images/videos, but there's no way to serve them
  - [ ] Implement `FileService` methods in `internal/service/files.go` — struct and constructor exist but zero methods
  - [ ] Add client-facing file download/streaming endpoint
- [ ] **Tech Debt (P2)**
  - [ ] Registration code validation bypassed — commented out in `internal/service/auth.go` (lines 62-69)
  - [ ] `CleanupOrphanedFiles` commented out in `internal/service/admin.go` (lines 179-198)

## Admin

- [ ] **File Upload UI (P1)** — backend endpoints exist (`PUT /api/admin/nodes/{node_id}/files`, `GET /api/admin/nodes/{node_id}/files`) but the admin frontend has no way to use them
  - [ ] Add file-related methods to `src/lib/api.ts` (upload, list, delete)
  - [ ] Build file upload component in the reward node edit dialog (`src/routes/graphs/[id]/edit-node-dialog/`)
  - [ ] Display and manage uploaded files for reward nodes in the editor
- [ ] **Tech Debt (P2)**
  - [ ] Impersonate button renders but `handleImpersonate` handler is empty (`src/lib/components/Impersonate.svelte`)

## Client

- [ ] **Graph Play Experience (P0)** — the core gameplay loop; the client can list graphs and show a countdown, but cannot navigate or play through a graph
  - [ ] Data layer (follow existing architecture: types -> API -> DB store -> service)
    - [ ] Define proper `RewardType` enum and full node/gate types in `src/lib/types/` (currently `reward_type` is an untyped `string`)
    - [ ] Add graph detail and node unlock methods to the API layer (`src/lib/api/`)
    - [ ] Add graph detail DB store for caching nodes, edges, and unlock state (`src/lib/db/stores/`)
    - [ ] Build graph navigation service with online-first + IndexedDB fallback (`src/lib/services/`)
  - [ ] UI — Graph navigation
    - [ ] Create graph detail route (e.g. `src/routes/(app)/graphs/[id]/`)
    - [ ] Build node map / progression component showing the player's position in the graph
    - [ ] Implement node detail view (shows gate requirements or reward content)
  - [ ] UI — Gate unlock flows
    - [ ] Location gate — geolocation proximity check + unlock
    - [ ] Code gate — code entry form + validation
    - [ ] Manual gate — "waiting for approval" state display
  - [ ] UI — Reward display
    - [ ] Image/video reward viewer (depends on backend file serving)
    - [ ] Other reward types: markdown, favour, calendar, wallet, file
- [ ] **Favour System (P1)** — backend is fully implemented (both admin and client endpoints), admin UI is done, client has nothing
  - [ ] Data layer
    - [ ] Create `src/lib/api/favour.api.ts`
    - [ ] Create favour DB store in `src/lib/db/stores/`
    - [ ] Create favour service in `src/lib/services/`
  - [ ] UI
    - [ ] Favour browsing page — view available choices and remaining favour count
    - [ ] Favour request flow — select a favour and submit a request
    - [ ] Request history / status page — view pending, approved, and rejected requests
- [ ] **Tech Debt (P2)**
  - [ ] `src/lib/db/sync.ts` is empty — background sync not implemented

## General

- [ ] `docker-compose.yml` is empty — no containerization configured (P2)
- [ ] No test framework or tests anywhere in the project (P2)
