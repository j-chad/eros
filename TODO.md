# TODO

> Priority: **P0** = critical path (blocks core features), **P1** = important, **P2** = nice to have

## Graph Play Experience (P0)

The core gameplay loop. The client can list graphs and show a countdown, but cannot navigate or play through them.

- [ ] Backend: register `UnlockNode` route in `router.go` — handler exists in `handler/client/graph.go` but no route is registered
- [ ] Backend: fix `ManualData` missing from `nodeDataDecoders` map in `internal/models/graph.go`
- [ ] Client: data layer (follow existing architecture: types -> API -> DB store -> service)
  - [ ] Define proper `RewardType` enum and full node/gate types in `client/src/lib/types/`
  - [ ] Add graph detail and node unlock methods to `client/src/lib/api/`
  - [ ] Add graph detail DB store for caching nodes, edges, and unlock state in `client/src/lib/db/stores/`
  - [ ] Build graph navigation service with online-first + IndexedDB fallback in `client/src/lib/services/`
- [ ] Client: graph navigation UI
  - [ ] Create graph detail route (e.g. `client/src/routes/(app)/graphs/[id]/`)
  - [ ] Build node map / progression component showing the player's position in the graph
  - [ ] Implement node detail view (shows gate requirements or reward content)
- [ ] Client: gate unlock flows
  - [ ] Location gate — geolocation proximity check + unlock
  - [ ] Code gate — code entry form + validation
  - [ ] Manual gate — "waiting for approval" state display

## File Management (P0 backend, P1 admin UI)

Reward nodes can contain images, videos, and files. The backend can store them, but can't serve them to the client, and the admin has no upload UI.

- [ ] Backend: implement `FileService` methods in `internal/service/files.go` — struct and constructor exist but zero methods
- [ ] Backend: add client-facing file download/streaming endpoint
- [ ] Admin: add file-related methods to `admin/src/lib/api.ts` (upload, list, delete)
- [ ] Admin: build file upload component in the reward node edit dialog (`admin/src/routes/graphs/[id]/edit-node-dialog/`)
- [ ] Admin: display and manage uploaded files for reward nodes in the editor
- [ ] Client: image/video reward viewer
- [ ] Client: other reward type displays (markdown, favour, calendar, wallet, file)

## Favour System (P1)

Backend is fully implemented (both admin and client endpoints). Admin UI is done. Client has nothing.

- [ ] Client: data layer
  - [ ] Create `client/src/lib/api/favour.api.ts`
  - [ ] Create favour DB store in `client/src/lib/db/stores/`
  - [ ] Create favour service in `client/src/lib/services/`
- [ ] Client: UI
  - [ ] Favour browsing page — view available choices and remaining favour count
  - [ ] Favour request flow — select a favour and submit a request
  - [ ] Request history / status page — view pending, approved, and rejected requests

## Tech Debt (P2)

- [ ] Backend: registration code validation bypassed — commented out in `internal/service/auth.go`
- [ ] Backend: `CleanupOrphanedFiles` commented out in `internal/service/admin.go`
- [ ] Admin: impersonate button renders but `handleImpersonate` handler is empty (`admin/src/lib/components/Impersonate.svelte`)
- [ ] Client: `src/lib/db/sync.ts` is empty — background sync not implemented
- [ ] `docker-compose.yml` is empty — no containerization configured
- [ ] No test framework or tests anywhere in the project
