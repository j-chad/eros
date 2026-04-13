# Gate Unlock Flows

## Problem

The client PWA can display graphs and navigate between accessible nodes, but every gate interaction is a dead end — buttons are disabled, forms don't submit, and the user can't actually progress through a graph. The backend unlock endpoint (`POST /api/nodes/{id}/unlock`) is fully implemented and validated, but nothing on the client calls it. This means the core gameplay loop — the reason Eros exists — doesn't work.

## Goals

- A user can progress through an entire graph from start to final reward using the client PWA.
- Each gate type (code, location, manual) has a complete, polished interaction flow.
- Unlocking a gate feels rewarding. Unlocking a reward feels genuinely celebratory.
- The unlock endpoint returns enough data for the client to transition without a full re-fetch.
- The backend rate-limits unlock attempts to prevent brute-force guessing.
- The client handles offline, GPS denial, and bad GPS accuracy gracefully.

## Non-Goals

- **Event/push system for manual gate approval.** Polling is the v1 approach. A WebSocket or SSE event system is a future improvement.
- **Cold-warm-hot location proximity.** The location gate is a single-shot check. Continuous proximity feedback is a future enhancement.
- **Alternate code spellings.** Code matching will be case-insensitive, but supporting multiple accepted answers per code gate is out of scope.
- **Reward content rendering.** Reward nodes will display their metadata (title, description, type icon, favour count) as they do today. Rendering actual images, videos, files, and markdown from the file system is tracked separately under File Management.
- **Offline unlock queuing.** Unlocks require server validation. The client will block unlock attempts when offline rather than queue them.

## Overview

This feature wires the existing client gate components to the backend unlock API, adds a response payload to the unlock endpoint (newly accessible nodes and edges), introduces per-node rate limiting, and builds the celebration/transition animations that make unlocking feel like the emotional centrepiece of the app.

## Detailed Design

### Backend: Unlock Response Payload

**Current behaviour:** `POST /api/nodes/{id}/unlock` returns `204 No Content` on success.

**New behaviour:** Returns `200 OK` with a JSON body containing the delta — the newly accessible nodes and edges after the unlock.

```json
{
  "unlocked_node": { /* the just-unlocked node, with unlocked_at now set */ },
  "new_nodes": [ /* nodes that became accessible as a result */ ],
  "new_edges": [ /* edges to/between newly accessible nodes */ ]
}
```

**How the delta is computed:**

1. Snapshot the set of accessible node IDs before the unlock (using the existing `accessible` CTE logic).
2. Perform the unlock (`SET unlocked_at = CURRENT_TIMESTAMP`).
3. Recompute the accessible set.
4. Diff: `new_nodes` = accessible-after minus accessible-before (excluding the unlocked node itself). `new_edges` = edges where the destination is in the new accessible set but wasn't before, plus edges between newly accessible nodes.
5. Return the unlocked node (refreshed, with `unlocked_at` populated) + the diff.

**Branching:** If the unlocked gate fans out to multiple paths, all newly reachable nodes and their connecting edges are included. The client already handles branching display.

**Node data:** Newly accessible nodes include their full type-specific data (gate config, reward metadata). The `ui_position` field is stripped as it already is for client responses.

### Backend: Per-Node Rate Limiting

An in-memory rate limiter on the unlock endpoint, scoped per node ID.

- **Limit:** 10 attempts per node per rolling 60-second window.
- **Implementation:** A simple token bucket or sliding window counter, keyed by node ID. No persistence needed — rate limits reset on server restart, which is acceptable.
- **Response:** `429 Too Many Requests` with a JSON error body: `{ "error": { "code": "RATE_LIMITED", "message": "Too many unlock attempts. Try again shortly." } }`.
- **Scope:** Only the `POST /api/nodes/{id}/unlock` endpoint. Admin unlock endpoints are not rate-limited.

Since this is a single-user app (one couple), there's no need for per-user keying. The node ID alone is sufficient.

### Backend: Case-Insensitive Code Matching

Change `validateCodeGatePayload` to compare using `strings.EqualFold(nodeData.Code, payload)` instead of `nodeData.Code != payload`. The stored code retains its original casing; only comparison is insensitive.

### Client: API Layer

Add to `client/src/lib/api/graph.api.ts`:

```typescript
export interface UnlockResult {
  unlocked_node: AnyNode;
  new_nodes: AnyNode[];
  new_edges: Edge[];
}

export async function unlockNode(nodeId: string, payload: string): Promise<UnlockResult> {
  return request<UnlockResult>(`/nodes/${nodeId}/unlock`, {
    method: 'POST',
    body: payload,
    headers: { 'Content-Type': 'text/plain' },
  });
}
```

The `Content-Type` is `text/plain` because the backend reads the raw body as the payload string (not JSON-encoded).

### Client: Service Layer

Add to `client/src/lib/services/graph.ts`:

```typescript
export async function unlockNode(nodeId: string, payload: string): Promise<UnlockResult> {
  const result = await unlockNodeAPI(nodeId, payload);
  // Merge result into cached graph detail in IndexedDB
  // (update the unlocked node's unlocked_at, append new_nodes and new_edges)
  return result;
}
```

The service merges the unlock delta into the IndexedDB-cached `GraphDetail` so that subsequent offline reads reflect the new state.

### Client: Code Gate Flow

**Component:** `client/src/lib/ui/nodes/CodeGate.svelte`

**States:**
1. **Input** (default) — Text input + "Unlock" button. Input is enabled, button is enabled.
2. **Submitting** — Button shows spinner, input disabled. Triggered on form submit.
3. **Incorrect** — Input shakes briefly (CSS animation), inline error text appears: "That's not quite right — try again." Input clears and re-focuses. Returns to Input state.
4. **Rate limited** — Inline error: "Too many attempts. Try again shortly." Input and button disabled briefly (5s client-side cooldown, then re-enable).
5. **Offline** — Input enabled (user can still type), button disabled with muted appearance. Hint text below: "You're offline. Connect to try unlocking."
6. **Success** — Transitions to celebration (see Celebration section).

**Behaviour:**
- On submit, call `unlockNode(node.id, inputValue)`.
- On 403 response → Incorrect state.
- On 429 response → Rate limited state.
- On network error / offline → show offline hint.
- On 200 → pass `UnlockResult` to parent for celebration + transition.

**The component accepts an `onUnlock` callback prop** that receives the `UnlockResult`. The page-level component (`+page.svelte`) handles the celebration and graph state update.

### Client: Location Gate Flow

**Component:** `client/src/lib/ui/nodes/LocationGate.svelte`

**States:**
1. **Ready** (default) — "Check location" button enabled. Description shows required proximity.
2. **Requesting permission** — Button shows spinner, text changes to "Getting your location..." The browser's permission prompt appears.
3. **Permission denied** — Button disabled. Message: "Location access is needed to continue. Check your browser settings and try again." A "Try again" link/button that re-requests.
4. **Acquiring fix** — Spinner with "Getting a precise fix..." Waiting for the Geolocation API response.
5. **Poor accuracy warning** — If `position.coords.accuracy` exceeds 2x the gate's radius: "Your location accuracy is low ({accuracy}m). Try moving to an open area." Button remains available to submit anyway ("Check anyway") or retry ("Try again").
6. **Submitting** — Sending coordinates to backend. Button shows spinner.
7. **Not in range** — "You're not at the right spot yet. Get closer and try again." Button re-enables.
8. **Rate limited** — Same pattern as code gate.
9. **Offline** — Button disabled. "You're offline. Connect to check your location."
10. **Success** — Celebration + transition.

**Behaviour:**
- On button tap, call `navigator.geolocation.getCurrentPosition()`.
- On success, format as `"lat,lng"` string and call `unlockNode(node.id, payload)`.
- If the position error is `PERMISSION_DENIED` → Permission denied state.
- If the position error is `POSITION_UNAVAILABLE` or `TIMEOUT` → show error with retry.
- On 403 → Not in range state.
- On 200 → celebration.

**Geolocation options:** `{ enableHighAccuracy: true, timeout: 15000, maximumAge: 0 }` — always request a fresh, high-accuracy fix.

### Client: Manual Gate Flow

**Component:** `client/src/lib/ui/nodes/ManualGate.svelte`

**States:**
1. **Waiting** — "Waiting for approval. Check back soon." Refresh button visible. Auto-polls every 30 seconds while the page is open.
2. **Approved** — `node.data.unlocked_at` is set. Badge changes to "Approved". "Continue" button appears.
3. **Submitting** — "Continue" button shows spinner. Calling unlock endpoint.
4. **Offline (waiting)** — Polling paused. "You're offline. Approval checks will resume when you reconnect." Refresh button disabled.
5. **Success** — Celebration + transition.

**Behaviour:**
- On mount, start a 30-second `setInterval` that re-fetches the graph via the service layer. If `node.data.unlocked_at` becomes non-null, transition to Approved state and stop polling.
- Polling only runs while the page is visible (`document.visibilityState === 'visible'`). Pause on hide, resume on show.
- Polling stops when approval is detected.
- "Refresh" button triggers an immediate re-fetch (resets the 30s timer).
- "Continue" calls `unlockNode(node.id, "")` (empty payload — the backend only checks that the admin has approved).
- On 403 (admin revoked approval between render and tap) → revert to Waiting state, resume polling.
- On 200 → celebration.

**Offline awareness:** check `navigator.onLine` before polling and before the Continue action. Listen to `online`/`offline` events to toggle state.

### Celebration and Transition

This is the emotional core of the app. Two tiers of celebration:

#### Standard Gate Unlock (gate → next gate)

1. **Unlock confirmed:** The gate UI transforms — icon swaps to a checkmark, colour shifts to `success`, a brief pulse animation plays. Text: "Unlocked!" Duration: ~1.5s.
2. **Auto-advance:** After the 1.5s beat, the current card content fades out (`animate-popIn` in reverse or a simple fade-down), and the next node's content fades in with `animate-popIn`. The graph state is already updated in memory from the unlock response.

#### Reward Unlock (gate → reward)

1. **Unlock confirmed:** Same checkmark + pulse as above, but held slightly longer (~2s) to build anticipation.
2. **Reveal transition:** The card content fades out, then the reward appears with an elevated entrance: scale up from ~90%, longer fade-in (~400ms), and a pink shimmer/glow effect on the card. The reward icon gets a brief heartbeat pulse.
3. **Optional flourish:** A subtle confetti burst or particle effect layered above the card. Should be tasteful — a few pink/rose particles drifting down, not a party cannon. If implementation complexity is high, ship without confetti in v0.1 and add it in v0.2.

#### Implementation Approach

The `+page.svelte` orchestrates the transition:

1. Gate component calls `onUnlock(result)`.
2. Page sets a `celebrationState` (`'gate-success'` | `'reward-reveal'` | null).
3. During celebration, the page merges `result.new_nodes` and `result.new_edges` into the reactive `nodes` and `edges` arrays, and updates the unlocked node's `unlocked_at`.
4. After the celebration delay (`setTimeout`), `celebrationState` resets to null and the derived node resolution logic (`activeNodes`, `currentNode`, etc.) picks up the new state automatically, rendering the next node.

The celebration UI is rendered as an overlay or replacement within the existing card, controlled by `celebrationState`.

### Offline Awareness (Cross-Cutting)

All three gate components share the same offline pattern:

- Track `navigator.onLine` as reactive state (listen to `online`/`offline` events on `window`).
- When offline: disable the primary action button, show hint text below it.
- When connectivity returns: re-enable automatically, no user action needed.
- Consider extracting a shared `useOnlineStatus()` utility in `$lib/` that returns a reactive boolean, so all three components don't duplicate the listener logic.

## Technical Approach

### Backend Changes

| File | Change |
|------|--------|
| `internal/service/graph.go` | `UnlockNode` returns `UnlockResult` struct. Calls repo to compute accessible diff. `validateCodeGatePayload` uses `strings.EqualFold`. |
| `internal/models/graph.go` | Add `UnlockResult` struct (`UnlockedNode`, `NewNodes`, `NewEdges`). |
| `internal/handler/client/graph.go` | `UnlockNode` handler returns `200` with JSON body instead of `204`. |
| `internal/repository/repository.go` | Add method to `Repository` interface for computing accessible nodes/edges for a graph (or extend `UnlockNode` to return the diff). |
| `internal/repository/sqlite/graph.go` | Implement the accessible-set diff query. |
| New: `internal/middleware/ratelimit.go` | In-memory per-key rate limiter middleware. Applied to the unlock route only. |

### Client Changes

| File | Change |
|------|--------|
| `src/lib/api/graph.api.ts` | Add `unlockNode()` function and `UnlockResult` type. |
| `src/lib/services/graph.ts` | Add `unlockNode()` that calls API and merges result into IndexedDB cache. |
| `src/lib/ui/nodes/CodeGate.svelte` | Full interactive flow: input, submit, error states, offline. Calls `onUnlock` on success. |
| `src/lib/ui/nodes/LocationGate.svelte` | Geolocation flow: permission, accuracy, submit, error states, offline. Calls `onUnlock` on success. |
| `src/lib/ui/nodes/ManualGate.svelte` | Polling, approval detection, continue button, offline. Calls `onUnlock` on success. |
| `src/routes/(app)/graphs/[id]/+page.svelte` | Orchestrate celebration state, merge unlock results into reactive graph data, auto-advance after delay. |
| New: `src/lib/online.svelte.ts` | Shared reactive online status utility using `$state` + event listeners. |
| New: `src/lib/ui/nodes/UnlockCelebration.svelte` | Celebration overlay component (standard + reward variants). |

### Patterns

- Follows the existing layered architecture: API → Service → Component.
- Rate limiter is a middleware wrapping the single unlock route, not a global concern.
- Gate components use `onUnlock` callback prop (existing `on`-prefix convention).
- Celebration is a presentational component controlled by the page, not embedded in gate components.

## Open Questions

1. **Rate limiter storage:** A simple `map[string][]time.Time` with a mutex works for a single-instance server. If the backend ever runs multiple instances, this would need to move to Redis or similar. Is single-instance fine for now? (Almost certainly yes.)
2. **Confetti/particles:** What's the right library or approach for the reward celebration flourish? Pure CSS animation? Canvas-based particles? A small library like `canvas-confetti`? Needs a spike during implementation.
3. **Accessible diff query performance:** Computing the diff requires two passes of the accessible CTE (before and after unlock). For small graphs this is trivial. Could be optimised to a single pass if needed, but premature for now.
4. **Manual gate poll interval:** 30 seconds is a starting point. In practice, how long does admin approval take? If it's typically fast (within a minute), 30s is fine. If it's hours, polling is wasteful and the event system becomes more urgent.

## Milestones

### v0.1 — Core Unlock Loop

The minimum to make graphs playable end-to-end.

- Backend: case-insensitive code matching.
- Backend: unlock endpoint returns `UnlockResult` with newly accessible nodes/edges.
- Backend: per-node rate limiting (10/min) on unlock endpoint.
- Client: `unlockNode()` in API and service layers.
- Client: Code gate — interactive input, submit, incorrect/rate-limited/offline states.
- Client: Location gate — geolocation request, submit, error states, offline.
- Client: Manual gate — polling, approval detection, continue button, offline.
- Client: Standard celebration (checkmark + pulse + auto-advance).
- Client: Shared online status utility.

### v0.2 — Polish and Reward Celebration

The upgrade from "functional" to "delightful."

- Client: Reward-specific celebration with elevated animation (scale, shimmer, glow).
- Client: Confetti/particle flourish on reward unlock.
- Client: Poor GPS accuracy warning on location gate.
- Client: IndexedDB cache merge in service layer (so offline reads reflect unlocks).
- Backend: tests for rate limiter, unlock result computation, case-insensitive matching.

### Future (Nice-to-Haves)

- Event system (WebSocket/SSE) to replace manual gate polling.
- Cold-warm-hot continuous proximity feedback for location gates.
- Alternate accepted codes per code gate (multiple valid answers).
