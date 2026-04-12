# TODO

> Priority: **P0** = critical path (blocks core features), **P1** = important, **P2** = nice to have

## Graph Play Experience

The core gameplay loop. The client can list graphs, show a countdown to the next one, and display a calendar of unlocked graphs — but cannot navigate or play through them.

- [ ] Client: gate unlock flows #P0
  - [ ] Location gate — geolocation proximity check + unlock
  - [ ] Code gate — code entry form + validation
  - [ ] Manual gate — "waiting for approval" state display

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
