# Security Audit Report

**Date:** 2026-04-04
**Scope:** Full monorepo — Go backend, admin frontend, client frontend
**Findings:** 6 Critical, 9 High, 14 Medium, 9 Low

---

## Critical

### 6. Auth Token Stored in Plaintext IndexedDB

**Files:** `client/src/lib/api/auth.ts:12-14`, `client/src/lib/db/stores/kv.ts:23-27`

The bearer token sits in IndexedDB unencrypted. Any XSS — from a dependency, extension, or future bug — reads it trivially. Combined with no CSP (#14), this is a full account takeover vector.

**Fix:** Move to httpOnly secure cookies set by the backend. If architecturally impossible, implement short token expiry with rotation and add CSP to reduce XSS surface.

---

## High

### 8. Device Tokens Stored in Plaintext in Database

**Files:** `backend/internal/repository/sqlite/device.go`, `backend/internal/repository/sqlite/sqlite_init.sql:24`

Tokens are stored and compared as plaintext. Database compromise exposes all active tokens.

**Fix:** Store `SHA-256(token)` in the database. Hash the incoming token before comparison.

---

### 10. CORS Wildcard With Credentials

**Files:** `backend/internal/config/config.develop.json:5`, `backend/internal/handler/middleware/cors.go:15,39`

Dev config allows `*` origin with `Access-Control-Allow-Credentials: true`. Any website can make authenticated cross-origin requests when running in development mode.

**Fix:** Never allow wildcard origins when credentials are enabled. Use explicit allowed origins in all environments.

---

### 11. Race Condition in Favour Request (Transaction Bypass)

**File:** `backend/internal/service/favour.go:50`

Uses `s.repo.CreateFavourRequest` instead of `txRepo.CreateFavourRequest` — the insert bypasses the transaction entirely. Concurrent requests can overspend favours.

**Fix:** Change line 50 to `txRepo.CreateFavourRequest(ctx, request)`.

---

### 12. Offline Auth Guard Bypass

**File:** `client/src/lib/services/auth.ts:14-16`

When offline, any token is accepted — even expired or revoked. Spoofing `navigator.onLine` keeps a stolen token working indefinitely.

**Fix:** Store token expiry alongside the token. Validate expiry even when offline.

---

### 13. Sensitive Game Data Cached Unencrypted

**Files:** `client/src/lib/db/stores/graph.ts`, `client/src/lib/db/schema.ts`

Gate codes, GPS coordinates, reward payloads — the entire treasure hunt — stored in plaintext IndexedDB. Anyone with DevTools can read all secrets.

**Fix:** The backend should strip sensitive fields from locked nodes in API responses. Never send data the client should not see.

---

### 14. No Content-Security-Policy (Both Frontends)

**Files:** `admin/src/app.html`, `client/src/app.html`

No CSP headers or meta tags anywhere. Zero defence-in-depth against XSS. Combined with token storage in IndexedDB/localStorage, one XSS equals full compromise.

**Fix:** Add CSP meta tags or configure CSP headers on the serving infrastructure:

```html
<meta http-equiv="Content-Security-Policy"
      content="default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline';
             connect-src 'self' https://your-api-domain.com; img-src 'self' data: https:;
             frame-ancestors 'none';">
```

---

### 15. Admin API Key in localStorage

**File:** `admin/src/lib/auth.svelte.ts:6,10`

The raw API key persists in localStorage indefinitely. Any XSS reads `localStorage.getItem('admin_api_key')`.

**Fix:** Move to `sessionStorage` at minimum. Implement short-lived session tokens rather than storing the raw API key.

---

## Medium

### 16. No Rate Limiting on Registration Endpoint

**File:** `backend/internal/handler/router.go:35`

`POST /api/device` is unauthenticated with no rate limiting. Even with code validation restored, brute force against registration codes is feasible at high request rates.

**Fix:** Add IP-based rate limiting middleware with exponential backoff.

---

### 17. No Rate Limiting on Node Unlock Endpoint

**File:** `backend/internal/handler/router.go:69`

`POST /api/nodes/{id}/unlock` has no rate limiting. Code gate codes can be brute-forced without restriction.

**Fix:** Add per-device rate limiting on the unlock endpoint.

---

### 18. Trusted `Content-Length` for Buffer Allocation

**File:** `backend/internal/handler/admin/device.go:44`

`make([]byte, r.ContentLength)` allocates based on the client-supplied header. A malicious `Content-Length: 2147483647` with a 1-byte body allocates ~2GB of memory (denial of service).

**Fix:** Use `io.ReadAll(io.LimitReader(r.Body, maxSize))` instead of trusting `Content-Length`.

---

### 19. Internal Error Details Leaked to Admin Responses

**File:** `backend/pkg/response/response.go:42-44`

When the context indicates admin, raw internal error messages (Go stack details, SQL errors, file paths) are included in JSON responses.

**Fix:** Log internal errors server-side. Return generic error messages in API responses.

---

### 20. Raw Error in `http.Error()`

**File:** `backend/internal/handler/admin/graphs.go:16`

`ListStartNodes` uses `http.Error()` with the raw Go error instead of the structured `response.Error()` pattern. Leaks internal details in plaintext.

**Fix:** Replace with `response.Error(r.Context(), w, apierror.UnknownInternalError(err))`.

---

### 21. No File Type Allowlist on Upload

**File:** `backend/internal/handler/admin/files.go:42-49`

No validation against an allowlist of expected MIME types. An admin could upload HTML files containing JavaScript, leading to stored XSS when served.

**Fix:** Validate MIME type and extension against an allowlist (images, videos, documents). Reject everything else.

---

### 22. No `MaxBytesReader` on JSON Request Bodies

**Files:** Multiple handler files using `json.NewDecoder(r.Body).Decode()`

JSON request bodies are decoded without size limits. Large payloads consume excessive memory.

**Fix:** Wrap `r.Body` with `http.MaxBytesReader(w, r.Body, 1<<20)` (1MB) before decoding.

---

### 23. No TLS Configuration

**File:** `backend/cmd/server/main.go:58`

Server listens on plain HTTP (port 80 in production). Auth tokens and admin API keys transmitted in cleartext.

**Fix:** Use `ListenAndServeTLS()` or deploy behind a TLS-terminating reverse proxy and document this requirement.

---

### 24. Error Page Leaks API Internals (Client)

**File:** `client/src/routes/+error.svelte:49-50`

The error page renders the full `page.error` object as JSON, including API base URL, endpoint paths, and potentially internal error strings from the backend.

**Fix:** Strip `base`, `endpoint`, and `body` fields from errors shown to users. Log them to the console only.

---

### 25. Error Page Leaks API Details (Admin)

**File:** `admin/src/lib/api.ts:39-45`

`App.Error` includes `base`, `endpoint`, `method`, and `body` (the full backend error response). Internal implementation details are exposed.

**Fix:** Sanitise error objects before storing in SvelteKit error state.

---

### 26. Save-Before-Verify Login Pattern (Admin)

**File:** `admin/src/routes/login/+page.svelte:14-29`

The login flow saves the API key to localStorage before verifying it with `api.ping()`. If the ping fails for a non-auth reason, an unverified key persists and the client-side auth gate is bypassed on next load.

**Fix:** Only persist the key after successful verification.

---

### 27. No Input Validation on Admin Forms

**Files:** Multiple admin dialog components (`EditCodeDialog.svelte`, `EditLocationDialog.svelte`, `EditRewardDialog.svelte`, `FavourChoiceTable.svelte`)

Client-side forms accept arbitrary input with minimal validation. No length limits, no bounds checking on lat/long, no URL scheme validation.

**Fix:** Add client-side validation as defence-in-depth: lat -90 to 90, long -180 to 180, HTTPS URLs only, maxlength on strings.

---

### 28. User-Supplied URLs Rendered as Media Sources

**File:** `admin/src/routes/graphs/[id]/edit-node-dialog/EditRewardDialog.svelte:225,259`

`<img src={payloadForm.url}>` and `<video src={payloadForm.url}>` render arbitrary URLs. Leaks admin IP to URL owner, enables data URIs.

**Fix:** Validate HTTPS scheme before rendering previews.

---

### 29. Service Worker No Response Integrity Validation

**File:** `client/src/service-worker.ts:66-78`

Cached responses are served without integrity checks. Cache poisoning (via compromised CDN or DNS hijacking during initial population) persists until the next service worker version change.

**Fix:** Add subresource integrity checks for critical assets. Consider network-first strategy for the app shell.

---

## Low

### 30. `APIError.Err` Field JSON-Serialisable

**File:** `backend/pkg/apierror/apierror.go:14`

The `Err` field is tagged `json:"err,omitempty"` instead of `json:"-"`. If any code path serialises an `APIError` directly, internal errors leak.

**Fix:** Change tag to `json:"-"`.

---

### 31. `context.Background()` in Auth Validation

**File:** `backend/internal/service/auth.go:44,54`

Database queries use `context.Background()` instead of the request context. Cannot be cancelled when clients disconnect.

**Fix:** Accept and propagate the request context.

---

### 32. Nil Body Close on S3 GET

**File:** `backend/internal/repository/storage/s3.go:91`

`defer req.Body.Close()` on a GET request where `req.Body` is `http.NoBody`. Potential nil pointer panic.

**Fix:** Remove `defer req.Body.Close()` from the GET method. Only response bodies need closing.

---

### 33. UUID Panic on Crypto Failure

**File:** `backend/internal/crypto/uuid.go:16`

`panic(err)` if `crypto/rand.Reader` fails. Crashes the entire server.

**Fix:** Return an error instead of panicking.

---

### 34. Memory Leak in DateDisplay Component

**File:** `admin/src/lib/components/DateDisplay.svelte:55-57`

`setInterval` runs at module level without cleanup. Each component instance creates an interval that runs forever.

**Fix:** Wrap in `onMount` with cleanup: `return () => clearInterval(id)`.

---

### 35. Console Statements Leak Sensitive Data

**Files:** Multiple files in both frontends (graph objects, form data, auth state, database operations)

`console.log` and `console.error` expose internal data to anyone with DevTools open.

**Fix:** Remove or gate behind `dev` checks. Use a logging utility that is silent in production.

---

### 36. No Clickjacking Protection

**Files:** Neither frontend configures `X-Frame-Options` or `frame-ancestors`

The app can be embedded in an iframe on an attacker's site.

**Fix:** Add `frame-ancestors 'none'` to CSP (see #14), or configure `X-Frame-Options: DENY` on serving infrastructure.

---

### 37. `skipWaiting()` Called Unconditionally

**File:** `client/src/service-worker.ts:39`

Forces new service worker to activate immediately. A compromised update takes effect instantly rather than waiting for all tabs to close.

**Fix:** Prompt the user to reload rather than force-activating.

---

### 38. Credentials Header Sent When No Origin Matched

**File:** `backend/internal/handler/middleware/cors.go:15`

`Access-Control-Allow-Credentials: true` is set unconditionally, even when no origin is matched (production config has empty origins).

**Fix:** Only set the credentials header when an origin is actually matched.

---

## Positive Findings

The following were done correctly:

- **SQL injection:** All queries use parameterised `?` placeholders. Zero string concatenation in SQL.
- **XSS via templates:** Zero `{@html}` usage across both frontends. Svelte escaping is active everywhere.
- **Code execution:** Zero `eval()`, `Function()`, `innerHTML`, or `document.write` calls.
- **Path traversal:** Local storage uses UUID regex allowlist in `safePath()`.
- **Admin auth timing:** `crypto/subtle.ConstantTimeCompare` used for admin token validation.
- **Token generation:** `crypto/rand` with sufficient entropy.
- **S3 keys:** Server-generated UUIDs, not user input.
- **Open redirects (admin):** `goto()` uses hardcoded paths only.
- **Config layering:** Environment variables override file-based config as highest priority.

---

## Priority Actions

1. **Uncomment registration code validation** (`backend/internal/service/auth.go:62-69`)
2. **Fix `go:embed` glob** to exclude `config.private.json`
3. **Replace `admin123`** with a strong key and rotate S3 credentials
4. **Fix open redirect** in client login with `returnTo` validation
5. **Add Content-Security-Policy** to both frontends
6. **Fix missing `return`** in `CreateGraph` handler
7. **Fix transaction bypass** in `RequestFavour`
8. **Use constant-time comparison** for code gates
9. **Add rate limiting** to registration and unlock endpoints
10. **Add `MaxBytesReader`** to all JSON-decoding handlers
