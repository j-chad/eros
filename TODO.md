# TODO

> Priority: **P0** = critical path (blocks core features), **P1** = important, **P2** = nice to have

## File Management

Reward nodes can contain images, videos, and files. The backend can store them, but can't serve them to the client, and the admin has no upload UI.

- [ ] Backend #P0
  - [x] implement `FileService` methods in `internal/service/files.go`
  - [x] add client-facing file download/streaming endpoint (`GET /api/files/{fileID}`)
  - [x] file metadata embedded in graph/unlock responses (post-processing in GraphService)
  - [x] replace semantics on upload (re-upload replaces previous file)
  - [x] UNIQUE constraint on `reward_file.node_id`
  - [x] repository additions: `GetFile`, `GetFileByNodeID`, `GetFilesByNodeIDs`, `DeleteFilesByNodeID`
  - [x] bump upload size limit to 50 MB
  - [x] fix storage key regex to allow digits in file extensions
- [ ] Admin #P1
  - [ ] add file-related methods to `admin/src/lib/api.ts` (upload, list, delete)
  - [ ] build file upload component in the reward node edit dialog (`admin/src/routes/graphs/[id]/edit-node-dialog/`)
  - [ ] display and manage uploaded files for reward nodes in the editor
- [ ] Client #P0
  - [ ] image/video reward viewer
  - [ ] other reward type displays (Markdown, favour, calendar, wallet, file)

## Tech Debt and Nice-to-Haves (P2)
- [ ] Backend: registration code validation bypassed — commented out in `internal/service/auth.go`
- [ ] Backend: `CleanupOrphanedFiles` commented out in `internal/service/admin.go`
- [ ] Admin: impersonate button renders but `handleImpersonate` handler is empty (`admin/src/lib/components/Impersonate.svelte`)
- [ ] Client: `src/lib/db/sync.ts` is empty — background sync not implemented
- [ ] Event system (WebSocket/SSE) — replace manual gate polling with real-time approval notifications
- [ ] Location gate: cold-warm-hot continuous proximity feedback (watch position, show distance hint)
- [ ] Code gate: alternate accepted codes — multiple valid answers per code gate node
- [ ] Backend: move rate limiter to Redis/shared store if backend ever runs multiple instances
