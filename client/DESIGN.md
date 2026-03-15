# Eros Client — Design Language

## Aesthetic

Eros is a private app for two people. The design should feel like it was made specifically for them — warm, intimate, and a little playful. It's not a generic productivity tool; it should have personality.

The visual language is soft and romantic: blush pinks, rounded everything, gentle shadows. At the same time it should feel polished — not a rough side project. Small moments of delight (a pop-in animation, a celebratory unlock) reinforce that this is something special.

**The three words to design to:** warm, playful, premium.

---

## Responsiveness & Layout
The app should work beautifully on both mobile and desktop. The most common use case is mobile, so that should be the primary focus, but the layout should adapt gracefully to wider screens.

## Colour

Built on DaisyUI's `valentine` theme. The palette is blush pinks with deep rose accents. Don't introduce colours outside this set.

- **Cards** sit on `base-100` (white/off-white) with a pink-tinted shadow (`shadow-pink-200/40`)
- **Primary actions** use the DaisyUI `primary` colour (deep rose)
- **Subtle insets** use `base-200` (a light grey-pink)

---

## Shape & Elevation

Everything is heavily rounded — this softness is core to the feel.

- **Cards and large containers:** `rounded-3xl`
- **Buttons, inputs, alerts, chips:** `rounded-2xl`
- **Elevated surfaces** (cards, the brand logo) carry a pink shadow. Plain white surfaces should never float without one.

---

## Typography

System font stack (no custom font loaded). Hierarchy is established through weight and opacity rather than size jumps:

- Headings: `font-bold` or `font-extrabold`
- Supporting text: `text-sm opacity-70` — dimmed, not a different colour
- Monospaced data (codes, IDs): `font-mono`
- Labels: `font-semibold` with a dimmed `opacity-60` hint alongside

---

## Animation

Animation is used intentionally to make the app feel alive — not to show off. The bar is: *does this make the interaction feel better, or is it just noise?*

**General principle:** short durations (150–300ms), ease-out curves. Never block the user waiting for an animation to finish.

**`animate-popIn`** (`app.css`) — the default entrance: a subtle fade + 6px upward drift + very slight scale. Use this whenever content appears due to a state change (step transitions, conditional panels, modals appearing).

**Reward/unlock reveals** should feel genuinely celebratory. This is the emotional centrepiece of the whole app — the moment a partner unlocks something. It warrants more: a longer entrance, a scale-up, possibly a confetti burst or shimmer. The animation should make the person smile. Don't hold back here.

**Loading states** should be smooth and unobtrusive — a DaisyUI spinner inside the button, button disabled. Avoid full-screen loading skeletons unless a page-level fetch takes > ~500ms.

**Transitions between steps/pages** should use `animate-popIn` or a similar fade-up so the new content doesn't just snap in.

---

## UX Patterns

### Loading & Transitions

- Inline loading: spinner inside the triggering button (`loading loading-spinner loading-sm`), button disabled while in-flight.
- Page-level loading: only show a skeleton or spinner if the wait is perceptible. Prefer optimistic UI where safe.
- Step changes within a page: always animate the incoming content with `animate-popIn`.

### Reward / Unlock Reveals

This is the most important moment in the app. Treat it accordingly:

- The reveal should be a distinct, full-attention moment — not just content appearing in a card.
- Use a larger entrance animation (scale up from ~95%, fade in, duration ~350–500ms).
- Consider a celebratory flourish: a confetti burst, a shimmer on the reward content, or a brief heartbeat pulse on the logo.
- The tone should feel like unwrapping a gift — anticipation, then delight.

### Confirmations

For destructive or irreversible actions (deleting a device, resetting progress):

- Use a modal or bottom sheet — don't inline the confirmation in the page.
- The cancel action should be the most visually prominent (ghost or secondary).
- The destructive action should be clearly labelled but not the default focus.
- Keep copy direct and calm — no alarm-bell language, but be honest about what's happening.

### Errors

- Inline errors appear at the top of the relevant card as an `alert alert-error rounded-2xl`.
- Keep error messages short and human. Avoid raw API error strings.
- Where possible, offer a next step ("Check permissions, or enter the code manually.").

---

## Voice & Copy

Short, warm, occasionally cheeky. No emojis — they cheapen the feel. Error messages are honest but not alarming. The app knows its audience is two people who chose to use it — write like that.

---

## Offline-First

The client is a PWA and should work without a network connection. This is a core constraint, not an afterthought.

**Data layer:** all user data is persisted in IndexedDB via the `lib/db/` layer. Pages read from local state first; the server is the source of truth for sync, not for rendering.

**Service worker:** static assets and cached API responses are served from the service worker. Assume the app shell always loads, even offline.

**Optimistic UI:** actions should feel instant. Apply changes locally first and sync in the background. Only block on the network when there's no safe local fallback (e.g. registration).

**Offline indicators:** surface connectivity state when it's relevant to the user — not as a persistent banner, but contextually (e.g. a sync error state when an action couldn't flush). Don't cry wolf on brief drops.

**Sync on reconnect:** pending actions should flush automatically when connectivity is restored. The user shouldn't have to manually retry things that can be retried safely.
