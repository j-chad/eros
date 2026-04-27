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

### 12. Offline Auth Guard Bypass

**File:** `client/src/lib/services/auth.ts:14-16`

When offline, any token is accepted — even expired or revoked. Spoofing `navigator.onLine` keeps a stolen token working indefinitely.

**Fix:** Store token expiry alongside the token. Validate expiry even when offline.

---

### 15. Admin API Key in localStorage

**File:** `admin/src/lib/auth.svelte.ts:6,10`

The raw API key persists in localStorage indefinitely. Any XSS reads `localStorage.getItem('admin_api_key')`.

**Fix:** Move to `sessionStorage` at minimum. Implement short-lived session tokens rather than storing the raw API key.

---

## Medium

### 21. No File Type Allowlist on Upload

**File:** `backend/internal/handler/admin/files.go:42-49`

No validation against an allowlist of expected MIME types. An admin could upload HTML files containing JavaScript, leading to stored XSS when served.

**Fix:** Validate MIME type and extension against an allowlist (images, videos, documents). Reject everything else.

---

### 26. Save-Before-Verify Login Pattern (Admin)

**File:** `admin/src/routes/login/+page.svelte:14-29`

The login flow saves the API key to localStorage before verifying it with `api.ping()`. If the ping fails for a non-auth reason, an unverified key persists and the client-side auth gate is bypassed on next load.

**Fix:** Only persist the key after successful verification.

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

