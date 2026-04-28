# TODO

> Priority: **P0** = critical path (blocks core features), **P1** = important, **P2** = nice to have

## Tech Debt and Nice-to-Haves (P2)
- [ ] Backend: `CleanupOrphanedFiles` commented out in `internal/service/admin.go`
- [ ] Admin: impersonate button renders but `handleImpersonate` handler is empty (`admin/src/lib/components/Impersonate.svelte`)
- [ ] Client: `src/lib/db/sync.ts` is empty — background sync not implemented
- [ ] Event system (WebSocket/SSE) — replace manual gate polling with real-time approval notifications
- [ ] Location gate: cold-warm-hot continuous proximity feedback (watch position, show distance hint)
