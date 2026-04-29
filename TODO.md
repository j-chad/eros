# TODO

> Priority: **P0** = critical path (blocks core features), **P1** = important, **P2** = nice to have

## Bugs (P0)
- [ ] Client: date type inconsistency — `GraphSummary` uses `Date` objects but IndexedDB round-trips them to strings, so `Calendar.svelte` will throw offline when calling `.getTime()` on a string

## Important (P1)
- [ ] Client: `ManualGate.svelte` polling uses `window.location.reload()` — should update reactive state instead (the pattern exists in `handleUnlock`)
- [ ] Client: markdown rewards render as plain text — `RewardNode.svelte` wraps payload in `whitespace-pre-wrap` but has no Markdown parser
- [ ] Client: service worker has no update notification — users silently stay on old cached version after deploys
- [ ] Client: no IndexedDB availability fallback — private browsing or quota exhaustion silently breaks the entire service layer

## Tech Debt and Nice-to-Haves (P2)
- [ ] Backend: uncomment and wire up `CleanupOrphanedFiles` in `internal/service/admin.go:217-236` — the supporting infrastructure (`FileStore.List()`, `FileStore.DeleteMany()`) is implemented on both local and S3, the orchestration method just needs enabling and a call site (e.g., on startup or periodic schedule)
- [ ] Backend: implement log collector / tracing — `config.CollectorConfig` (`Enabled`, `MaxSpans`, `TraceHeader`) is parsed but nothing reads it; wire up middleware to extract the trace header and propagate a per-request logger
- [ ] Backend: wire up request-scoped logging — `logging.NewContext()` is defined but never called, so `FromContext()` always falls through to `slog.Default()`; add middleware that injects a logger with request-scoped fields (request ID, method, path) via `NewContext()`
- [ ] Admin: impersonate button renders but `handleImpersonate` handler is empty (`admin/src/lib/components/Impersonate.svelte`)
- [ ] Client: `src/lib/db/sync.ts` is empty — background sync not implemented
- [ ] Event system (WebSocket/SSE) — replace manual gate polling with real-time approval notifications
- [ ] Location gate: cold-warm-hot continuous proximity feedback (watch position, show distance hint)
- [ ] Client: no graph completion screen — finishing an adventure just shows the last reward node again, no celebratory end state (despite `GraphStatus` having a `'completed'` variant)
- [ ] Client: favour fulfilment has no real-time update — user must manually revisit the page to see status changes
- [ ] Client: `Input` base component (`lib/ui/base/Input.svelte`) is built but unused — login page and `CodeGate` hand-roll inputs
- [ ] Client: duplicated unlock logic across 4 gate components — extract shared helper to `lib/ui/nodes/useGateUnlock.svelte.ts`
- [ ] Client: dynamic imports in event handlers (`CodeGate`, `ManualGate`, `favours/+page`) where static imports would suffice
- [ ] Client: empty `src/lib/db/sync.ts` and unused `ApiError` interface in `app.d.ts` — dead code to clean up
