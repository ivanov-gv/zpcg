# Monterails Telegram Mini App — Design & Research

> **Status**: prototype / research.
> This document accompanies the working prototype in `webapp/`.

---

## Table of contents

1. [What this is](#what-this-is)
2. [Problem statement](#problem-statement)
3. [Prototype walkthrough](#prototype-walkthrough)
4. [Feature inventory](#feature-inventory)
5. [Architecture options](#architecture-options)
6. [Technology choices: HTMX vs. SPA vs. hybrid](#technology-choices)
7. [Integration with the existing bot](#integration-with-the-existing-bot)
8. [Data pipeline](#data-pipeline)
9. [Hosting & deployment](#hosting--deployment)
10. [Telegram Mini App specifics](#telegram-mini-app-specifics)
11. [Risks & open questions](#risks--open-questions)
12. [Roadmap ideas](#roadmap-ideas)
13. [Decision log](#decision-log)

---

## What this is

A **Telegram Mini App** (formerly "WebApp") for
[@Monterails_bot](https://t.me/Monterails_bot) — the Montenegro Railways
timetable bot described in the main README. The mini app is opened from
inside Telegram and gives users a richer interface than plain chat messages
can provide: interactive station search, a visual map, price estimates,
saved routes, and alert subscriptions.

The prototype in `webapp/` is a **fully static, zero-build, pure-vanilla-JS
single-page app** that loads JSON data exported from the same
`gen/timetable/timetable.gen.go` the bot itself compiles. It runs on GitHub
Pages (or any static host) and needs no server beyond what the bot already
has.

---

## Problem statement

The README's "Ways to improve" section identifies two key UX gaps:

1. **"I prefer clicking buttons over raw keyboard input."** — Common user
   feedback. The text-based `Podgorica, Niksic` input is powerful but
   alien to casual users.
2. **Extra information (prices, route maps, ticket offices) can't fit in
   one small Telegram message.** — A rich web view is needed.

A Mini App addresses both: it provides tap-friendly station pickers, a
visual map, and room for price breakdowns, train details, and alerts — all
without leaving Telegram.

---

## Prototype walkthrough

### Screen 1 — Search (home)

```
┌──────────────────────────────┐
│ ‹  Monterails               EN│
├──────────────────────────────┤
│                               │
│  Where to?                    │
│  ┌────────────────────┐  ⇅   │
│  │🟢 From              │      │
│  │   Bar               │      │
│  ├────────────────────┤      │
│  │🔴 To                │      │
│  │   Podgorica         │      │
│  └────────────────────┘      │
│                               │
│  [Today] [Tomorrow] [Wed]     │
│                               │
│  ┌─────────────────────────┐ │
│  │     Find trains          │ │
│  └─────────────────────────┘ │
│                               │
│  POPULAR ROUTES               │
│  🚆 Bar → Podgorica          │
│  🚆 Podgorica → Nikšić       │
│  🚆 Bar → Bijelo Polje       │
│  🚆 Podgorica → Kolašin      │
│  🚆 Bar → Beograd Centar     │
│                               │
├───────┬───────┬───────┬─────┤
│ 🔍    │ 🗺️    │ 🔔    │ ⭐  │ ℹ️
│Search │ Map   │Alerts │Saved│Info
└───────┴───────┴───────┴─────┘
```

**What it does:**
- Tap "From" / "To" → opens a **bottom sheet** station picker with fuzzy
  search (Latin, Cyrillic, diacritic-insensitive).
- Swap button (⇅) reverses from↔to.
- "Find trains" navigates to results.
- Popular routes are 1-tap shortcuts.

### Screen 2 — Station picker (bottom sheet)

```
┌──────────────────────────────┐
│ ┌──────────────────────┐  ✕ │
│ │ Search station…       │    │
│ └──────────────────────┘    │
│                               │
│  Bar                     main │
│  Бар                          │
│                               │
│  Podgorica               main │
│  Подгорица                    │
│                               │
│  Bijelo Polje            main │
│  Бијело Поље                  │
│                               │
│  Nikšić                  main │
│  Никшић                       │
│  ...                          │
└──────────────────────────────┘
```

**Details:**
- Shows station name (Latin), Cyrillic below, station type on the right.
- Blacklisted stations (Budva, Kotor, etc.) are dimmed; tapping them shows
  the "no train station here" warning — same copy as the bot.
- Fuzzy: typing "pdgrca" still finds Podgorica.

### Screen 3 — Results

```
┌──────────────────────────────┐
│ ‹  Search results         EN │
├──────────────────────────────┤
│  Bar → Podgorica    2026-04-12│
│                               │
│  ✓ Direct                     │
│                               │
│  ┌──────────────────────────┐│
│  │ № 6100                    ││
│  │ 05:13  →  06:16   1h 03m ││
│  │ Bar → Podgorica   ≈€3.36 ││
│  └──────────────────────────┘│
│  ┌──────────────────────────┐│
│  │ № 6102                    ││
│  │ 09:40  →  10:49   1h 09m ││
│  │ Bar → Podgorica   ≈€3.36 ││
│  └──────────────────────────┘│
│  ... more cards ...           │
│                               │
│  [⭐ Save route] [⇅ Reverse] │
│                               │
└──────────────────────────────┘
```

**For interchange routes** (e.g. Bar → Nikšić):

```
  ┌──────────────────────────┐
  │ 6100 → 7100              │
  │ 05:13  →  09:07   3h 54m │
  │ 🔄 Podgorica              │
  │   06:16 → 08:00 (wait 1h) │
  │ Bar → Nikšić     ≈€6.72  │
  └──────────────────────────┘
```

### Screen 4 — Train details

```
┌──────────────────────────────┐
│ ‹  № 6100                 EN │
├──────────────────────────────┤
│  № 6100       26 stops       │
│                               │
│  ● Bar              05:13    │
│  │                            │
│  ○ Šušanj           05:15    │
│  │                            │
│  ○ Sutomore         05:27    │
│  │                            │
│  ○ Crmnica          05:35    │
│  │                            │
│  ○ Virpazar     05:40→05:45  │
│  │                            │
│  ... (all stops listed) ...   │
│  │                            │
│  ● Podgorica    06:16→06:21  │
│  │                            │
│  ○ Zlatica          06:26    │
│  │  ...                       │
│  ● Bijelo Polje     08:41    │
│                               │
│  [💶 Est. price]              │
└──────────────────────────────┘
```

Departure and arrival stops are highlighted with filled dots (●).

### Screen 5 — Price estimator

```
┌──────────────────────────────┐
│ ‹  Bar → Podgorica        EN │
├──────────────────────────────┤
│  Class    [2nd] [1st]        │
│  👤       [Adult] [Child]    │
│           [Student] [Senior] │
│  ×        [1] [2] [3] [4]   │
│                               │
│  ≈ 56 km × €0.06 .... €3.36 │
│  Total ............... €3.36 │
│                               │
│  Estimate €0.06/km — verify  │
│  at ticket office.            │
└──────────────────────────────┘
```

### Screen 6 — Map

```
┌──────────────────────────────┐
│ ‹  Stations map            EN │
├──────────────────────────────┤
│                               │
│     🟡 Nikšić     ┌ Legend ─┐│
│       \           │● Main   ││
│        \          │○ Stop   ││
│   ●Podgorica      │🟡 Sel   ││
│         |         │— Route  ││
│    ○ Virpazar     └─────────┘│
│         |                     │
│    ● Bar          ▭ Leaflet  │
│                               │
│  Tap a station to pick it     │
│                               │
├───────┬───────┬───────┬─────┤
│ 🔍    │ 🗺️    │ 🔔    │ ⭐  │ ℹ️
└───────┴───────┴───────┴─────┘
```

Uses **Leaflet + OpenStreetMap**. Route polyline drawn for the selected
from→to pair. Tapping a marker picks it as from/to.

### Screen 7 — Alerts (mocked)

Shows disruption cards (from ZPCG announcements) and user alert
subscriptions stored in localStorage. "Create alert" adds a new one.
In production this would push notifications through the bot.

### Screen 8 — Saved routes

LocalStorage-backed list of from→to pairs. Tapping one loads it into
search and navigates to results.

### Screen 9 — Info / About

Links to the bot, zpcg.me, zcg-prevoz.me, GitHub. Shows timetable
generation timestamp, language picker.

---

## Feature inventory

Everything below is categorised as **built** (in the prototype),
**mockable** (UI exists but no real backend), or **idea** (not built yet).

| # | Feature | Status | Value | Effort |
|---|---------|--------|-------|--------|
| 1 | Station picker with fuzzy search | Built | High | Done |
| 2 | Direct + 1-transfer pathfinding | Built | High | Done |
| 3 | Train details (all stops, times) | Built | High | Done |
| 4 | Leaflet stations map | Built | High | Done |
| 5 | Route polyline on map | Built | Medium | Done |
| 6 | Per-station detail page | Built | Medium | Done |
| 7 | Price estimator (€/km model) | Built | Medium | Done |
| 8 | Class / passenger / count selectors | Built | Medium | Done |
| 9 | Saved routes (localStorage) | Built | Medium | Done |
| 10 | Popular-route shortcuts | Built | Medium | Done |
| 11 | Language picker (EN / RU / SR) | Built | High | Done |
| 12 | Telegram theme sync (light/dark) | Built | Medium | Done |
| 13 | Alert subscriptions UI | Mockable | High | Low |
| 14 | Service disruption cards | Mockable | High | Low |
| 15 | Push notifications via bot | Idea | High | Medium |
| 16 | Real price table (from ZPCG) | Idea | High | Medium |
| 17 | Ticket office hours per station | Idea | Medium | Medium |
| 18 | "Next train" countdown timer | Idea | High | Low |
| 19 | Calendar / date filter | Idea | Medium | Low |
| 20 | Share route (Telegram deep link) | Idea | Medium | Low |
| 21 | Nearest station (geolocation) | Idea | Medium | Low |
| 22 | Accessibility (a11y / screen reader) | Idea | Medium | Medium |
| 23 | Offline mode (service worker cache) | Idea | Medium | Medium |
| 24 | Platform-specific fare zones | Idea | Low | High |
| 25 | Live train position (if API exists) | Idea | High | High |
| 26 | Multi-day journey planner | Idea | Low | High |
| 27 | Integration with Google/Apple Wallet | Idea | Low | High |

### Detailed notes on selected features

#### 15 — Push notifications via bot

Telegram bots can send messages to users who have interacted with them.
The mini app can POST alert subscriptions to a bot endpoint, which stores
them in-memory (fits the stateless constraint if alerts are simple) or in a
lightweight KV store. A cron-like trigger checks timetable times and sends
messages N minutes before departure.

**Constraint**: the bot is fully stateless today — no database. Adding
alerts requires *some* persistence (even a JSON file on GCS or a free-tier
Firestore). This is the biggest architectural decision for alerts.

#### 16 — Real price table

ZPCG does not publish a machine-readable price list. Options:
- Scrape the ticket office price board (photos exist on Google Maps reviews).
- Maintain a hand-curated CSV in the repo, updated annually.
- Use the €0.06/km model with a calibration table for known routes.

#### 18 — "Next train" countdown

Given the current time and the selected route, show "next train in 47 min"
with a live countdown. Extremely useful for commuters. Implementation is
trivial — compare `Date.now()` with timetable departure times.

#### 21 — Nearest station (geolocation)

`navigator.geolocation` + haversine against the coordinates table. Would
auto-fill the "From" field. Telegram's `requestLocation` API is another
option for more reliable location inside the Telegram client.

#### 25 — Live train position

ZPCG does not expose a real-time vehicle-tracking API. The only path
would be crowd-sourced reports or GTFS-RT feed if ZPCG ever publishes one.
Marked as high-effort/high-value but blocked on data availability.

---

## Architecture options

### Option A — Fully static (current prototype)

```
┌──────────┐   CDN / GitHub Pages   ┌──────────┐
│ Telegram │ ─────────────────────→ │ index.html│
│ WebView  │                        │ + JS + JSON│
└──────────┘                        └──────────┘
```

- **Pros**: zero ops cost, instant deploy via GitHub Pages, no server
  changes needed, same data as the bot (generated JSON).
- **Cons**: no live data, no push alerts, no server-side fuzzy matching.
- **Best for**: launching quickly and validating the UX with real users.

### Option B — HTMX + Go backend

```
┌──────────┐        ┌─────────────────────┐
│ Telegram │  HTTPS │ Go server (Cloud Run)│
│ WebView  │ ──────→│  /api/search         │
│  HTMX    │ ←──────│  /api/train/:id      │
└──────────┘  HTML  │  /api/map (static)   │
              frags │  serves webapp/      │
                    └─────────────────────┘
```

- HTMX sends `hx-get="/api/search?from=bar&to=podgorica"` and the server
  returns an HTML fragment — no JSON ↔ DOM dance.
- **Pros**: server-side rendering, real fuzzy matching (approximate-match
  lib already exists in Go), one deployment unit, HTML fragments are tiny.
- **Cons**: adds latency (Cloud Run cold start ~1s), slightly higher ops
  cost, couples the mini app to the bot's deployment lifecycle.
- **Best for**: if the bot already needs a persistent HTTP endpoint (e.g.
  for alerts or live disruptions).

#### HTMX example

```html
<!-- station search triggers a server fragment swap -->
<input type="search" name="q"
       hx-get="/api/stations"
       hx-trigger="input changed delay:200ms"
       hx-target="#station-list" />
<ul id="station-list">
  <!-- server returns <li> fragments with station names -->
</ul>
```

```go
// Go handler returns HTML fragments
func stationsHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    matches := fuzzySearch(q)
    for _, s := range matches {
        fmt.Fprintf(w, `<li hx-get="/api/pick/%d">%s</li>`, s.Id, s.Name)
    }
}
```

This is clean and simple. HTMX shines here because:
- Fragments are tiny (~200 bytes per station `<li>`)
- No JavaScript state management needed
- Server uses the same Go fuzzy-matching logic the bot already has
- Dark theme is CSS-only — server doesn't care

### Option C — Hybrid (static shell + HTMX fragments for live data)

```
┌──────────┐   GitHub Pages    ┌───────────┐
│ Telegram │ ────────────────→ │ static    │
│ WebView  │                   │ shell     │
│          │   hx-get          │           │
│          │ ────────────────→ │ Go server │ (only for /api/alerts etc.)
└──────────┘                   └───────────┘
```

- Static HTML shell + JSON for the timetable (works offline).
- HTMX calls to the bot server only for live features: alerts, disruptions,
  "next train" status.
- **Pros**: best of both — free static hosting for the core, small server
  surface for live features only.
- **Cons**: two deploy targets, CORS setup needed.

### Key tradeoff: offline capability

The static prototype (Option A) loads all JSON once and then works **fully
offline**. This matters for real users: the Podgorica–Kolašin corridor runs
through mountains with poor or no mobile signal, and tunnels cut connection
entirely. A passenger who opened the app at the station can still browse
timetables mid-journey.

A pure HTMX app (Option B) **cannot work offline at all** — every tap
fetches an HTML fragment from the server. No server = blank screen.

The hybrid (Option C) preserves offline timetable browsing while adding
HTMX only for features that inherently need a server (alerts, disruptions,
live status). This is the key reason to prefer C over B.

### Recommendation

**Start with Option A** (fully static, current prototype) to validate UX
and gather user feedback. **Migrate to Option C** (hybrid) when alerts or
live disruptions become a priority. Option B is viable only if offline
support is not a concern — which, given Montenegro's mountain rail
corridors, it is.

---

## Technology choices

### Why vanilla JS (not React, Vue, Svelte)?

1. **Zero build step** — drop files on GitHub Pages and it works.
2. **Tiny bundle** — the whole prototype is ~55 KB of JS (minified by
   esbuild). React alone is 40+ KB.
3. **Telegram WebView constraint** — the in-app browser is Chromium-based
   but memory-constrained on older Android devices. Less JS = faster load.
4. **Prototype scope** — a framework adds value when there are many
   interconnected stateful components. This app has ~8 screens with
   minimal shared state (from, to, language).

### Why HTMX is interesting for production

If the app goes to Option B/C architecture, HTMX is a strong fit because:
- The Go server already exists (the bot is a Go HTTP handler).
- HTMX needs no build step either — `<script src="htmx.min.js">`.
- Server-side rendering means the app works even if JavaScript is slow to
  load (progressive enhancement).
- Station search is the highest-interaction feature and HTMX's
  `hx-trigger="input changed delay:200ms"` handles it elegantly.

### Why Leaflet for the map

- Open source, lightweight (~40 KB), uses OpenStreetMap tiles (free).
- Alternatives: Mapbox GL JS (heavier, needs API key), Google Maps
  (cost per load), no map at all (link to Google Maps like the bot does).
- For this prototype, Leaflet is the clear winner on simplicity and cost.

---

## Integration with the existing bot

### Launching the mini app

Telegram bots can attach a **Menu Button** that opens a URL as a WebApp.
Configuration via [@BotFather](https://t.me/BotFather):

```
/mybots → @Monterails_bot → Bot Settings → Menu Button → Configure
→ URL: https://ivanov-gv.github.io/zpcg/webapp/
→ Title: "Timetable"
```

Alternatively, the bot can send an inline keyboard with a `web_app` button
on `/start` or any message:

```go
keyboard := [][]gotgbot.InlineKeyboardButton{{
    {Text: "Open timetable", WebApp: &gotgbot.WebAppInfo{
        Url: "https://ivanov-gv.github.io/zpcg/webapp/",
    }},
}}
```

### Passing context bot → app

When the mini app opens, `Telegram.WebApp.initDataUnsafe` contains:
- `user.language_code` — used to auto-select language.
- `user.id` — could be used for alert subscriptions.
- `start_param` — a bot can pass `?startapp=bar-podgorica` to deep-link
  directly to a route.

### Passing data app → bot

The mini app can call `Telegram.WebApp.sendData(json)` to pass a message
back to the bot (appears as a special message from the user). Use cases:
- "Search this route in the bot" → sends `{"from":"Bar","to":"Podgorica"}`.
- "Subscribe to alert" → sends `{"alert":"Bar-Podgorica","time":"08:00"}`.

The bot's webhook handler parses this message and responds normally.

---

## Data pipeline

```
 zpcg.me/search
      │
      ▼
 cmd/exporter          parses HTML, emits Go code
      │
      ▼
 gen/timetable/timetable.gen.go    compiled into bot binary
      │
      ▼
 cmd/webapp-export     reads the gen package, emits JSON
      │
      ▼
 webapp/data/
    stations.json    (75 stations)
    trains.json      (30 trains, ordered stops)
    meta.json        (transfer station, types, timestamp)
    coordinates.json (hand-curated lat/lng — MUST be verified)
```

**Keeping data fresh**: after every `make parse_timetable`, re-run
`go run ./cmd/webapp-export -out=webapp/data` and commit the updated JSON.
A GitHub Action can automate this.

### Coordinates

Station coordinates are **approximate** in this prototype. The main
stations (Bar, Podgorica, Nikšić, Bijelo Polje, Belgrade, etc.) are
accurate. Minor halts between them are interpolated along the corridor.

**To verify/fix**: query OpenStreetMap's Overpass API for
`railway=station` + `railway=halt` in Montenegro and Serbia, then match
by name. This could be a one-time script in `cmd/`.

---

## Hosting & deployment

### GitHub Pages (recommended for prototype)

1. Add `webapp/` to the repo (already done).
2. In repo Settings → Pages → Source: "Deploy from a branch" → `main` →
   `/webapp` folder (or use a `gh-pages` branch).
3. The app is live at `https://ivanov-gv.github.io/zpcg/webapp/`.
4. Cost: free.

### Cloud Run (for Option B/C)

The bot already runs on Cloud Run. Serving the mini app from the same
container adds zero marginal cost — just serve `webapp/` as static files
alongside the webhook handler. Cloud Run's cold-start is <1s which is
acceptable for an HTMX fragment server.

### Cloudflare Pages / Vercel / Netlify

All offer free-tier static hosting with HTTPS and global CDN. Any of them
work as alternatives to GitHub Pages.

---

## Telegram Mini App specifics

### Theme integration

The Telegram WebApp SDK injects CSS variables for the user's theme:
`--tg-theme-bg-color`, `--tg-theme-text-color`, etc. The prototype mirrors
these onto its own CSS variables so the app automatically matches the
user's Telegram theme (light or dark).

### Viewport & safe areas

- `Telegram.WebApp.expand()` requests full-height mode.
- `viewport-fit=cover` + `env(safe-area-inset-bottom)` handles notch/pill
  areas on iPhones.
- The tab bar accounts for the safe area.

### Back button

The SDK provides `Telegram.WebApp.BackButton`. The prototype wires the
header "‹" button to `history.back()` for now; in production, the native
back button should be used for a more integrated feel.

### Main button

`Telegram.WebApp.MainButton` is a full-width button at the bottom of the
WebView, styled by the OS. It could replace the "Find trains" button for a
more native feel.

### Haptic feedback

`Telegram.WebApp.HapticFeedback.impactOccurred("light")` can be called
on button taps for a native feel on supported devices.

### Permissions

| Feature | Telegram API | Notes |
|---------|-------------|-------|
| Language | `initDataUnsafe.user.language_code` | Auto, no permission |
| Location | `requestLocation()` | Prompts user |
| Send data | `sendData()` | One-shot, closes app |
| Open link | `openLink()` | External browser |

---

## Risks & open questions

| # | Risk / question | Mitigation |
|---|----------------|------------|
| 1 | **Coordinates are approximate.** Minor halts may be 1-5 km off on the map. | Verify against OSM Overpass API. One-time script. |
| 2 | **Price model is a guess.** No official machine-readable pricing. | Clearly label as "estimate". Maintain a calibration CSV. |
| 3 | **Timetable freshness.** If someone forgets to re-run the exporter after a parse, the app shows stale data. | Add a CI step: parse → export → commit JSON. |
| 4 | **No push notifications without persistence.** The bot is stateless by design. | Evaluate free-tier Firestore or a JSON file on GCS for alert storage. |
| 5 | **Cross-midnight trains.** Train 432 (Bijelo Polje → Belgrade) departs at 23:20 and arrives next day. The prototype's time comparison may fail for these. | Add a `+24h` offset flag in the exporter for cross-midnight trains. |
| 6 | **Blacklisted stations are hard-coded in JS.** The bot's list lives in Go. | Export the blacklist as part of `meta.json` so there's one source of truth. |
| 7 | **Telegram WebView on old Android.** Some Android 7/8 devices have an old WebView that may not support ES modules. | esbuild can bundle to a single IIFE for compat. |
| 8 | **No HTTPS on localhost.** Telegram requires HTTPS for Mini Apps. | Use ngrok or Telegram's test environment for local dev. |
| 9 | **Seasonal timetable changes.** Summer schedule adds Subotica route, blacklist changes. | Already handled by the exporter — re-run after seasonal parse. |
| 10 | **Accessibility.** Prototype has no ARIA roles beyond `role="dialog"` on the picker. | Add aria-labels, focus trapping, keyboard navigation in production. |

---

## Roadmap ideas

### Phase 1 — Ship the static prototype (2-3 days)

- [x] Build the webapp (this PR)
- [ ] Verify coordinates against OSM
- [ ] Wire up BotFather menu button
- [ ] Deploy to GitHub Pages
- [ ] Announce in bot `/start` message: "Try the new timetable app!"
- [ ] Gather user feedback for 2 weeks

### Phase 2 — Live features (1-2 weeks)

- [ ] "Next train" countdown on results screen
- [ ] Share route via Telegram deep link (`tg://resolve?domain=Monterails_bot&startapp=bar-podgorica`)
- [ ] Nearest station via geolocation
- [ ] Export blacklist to meta.json
- [ ] CI: auto-regenerate JSON on timetable parse
- [ ] Cross-midnight train handling fix

### Phase 3 — HTMX server + alerts (2-4 weeks)

- [ ] Add lightweight `/api/search` and `/api/alerts` endpoints to the bot
- [ ] Migrate station search to HTMX (server-side fuzzy matching)
- [ ] Alert subscriptions with Firestore persistence
- [ ] Push notifications: bot sends message N min before departure
- [ ] Service disruption feed (if ZPCG publishes one)

### Phase 4 — Polish & expand (ongoing)

- [ ] Real price table (hand-curated or scraped)
- [ ] Ticket office hours per station
- [ ] Offline mode with service worker
- [ ] Additional languages (UK, BE, DE, SK, TR, HR)
- [ ] Accessibility audit
- [ ] Analytics (privacy-respecting — no cookies, just bot-side event counts)

---

## Decision log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-04-12 | Vanilla JS, no framework | Zero build step; tiny bundle (~55 KB); sufficient for 8 screens |
| 2026-04-12 | Static JSON, not live API | Bot is stateless; JSON is generated from the same compiled timetable |
| 2026-04-12 | Leaflet + OSM, not Google Maps | Free, open source, lightweight, no API key |
| 2026-04-12 | €0.06/km price model | No official price data available; clearly labeled as estimate |
| 2026-04-12 | Hash-based router, not History API | Works on `file://` for local dev; no server config needed |
| 2026-04-12 | EN/RU/SR languages only (for now) | Highest-traffic languages; matches bot's render package |
| 2026-04-12 | Approximate coordinates | Main stations real, halts interpolated; good enough for prototype |
| 2026-04-12 | localStorage for saved/alerts | No persistence needed for prototype; production needs server |

---

## Appendix: all Montenegro railway stations

Extracted from the bot's compiled timetable. 75 stations across 4 types:

### Main stations (type 4)
Bar, Podgorica, Danilovgrad, Nikšić, Bijelo Polje, Beograd Centar,
Novi Sad, Subotica

### Stations (type 1)
Sutomore, Golubovci, Kolašin, Mojkovac, Trebešica, Spuž, Kosjerić,
Lajkovac, Lazarevac, Rakovica, Zemun, Virpazar, Novi Beograd, Čačak,
Kraljevo, Kragujevac, Lapovo, Velika Plana, Branešci, Brodarevo, Vrbnica,
Zmajevo, Inđija, Stara Pazova, Nova Pazova, Lovćenac

### Stops (type 2)
Šušanj, Crmnica, Vranjina, Morača, Aerodrom, Zlatica, Kruševački Potok,
Selište, Mateševo, Padež, Oblutak, Štitarička Rijeka, Žari, Slijepač Most,
Ravna Rijeka, Lješnica, Pričelje, Ljutotuk, Slap, Bare Šumanovića,
Dabovići, Stubica, Prijepolje, Priboj, Vrbas, Beška, Bačka Topola

### Crossings (type 3)
Zeta, Bioče, Kos, Trebaljevo, Mijatovo Kolo, Kruševo, Ostrog, Šobajići,
Požega, Bratonožići, Lutovo

### Rail branches
- **Bar branch**: Bar — Šušanj — Sutomore — Crmnica — Virpazar — Vranjina — Zeta — Morača — Aerodrom — Golubovci — **Podgorica**
- **Nikšić branch**: **Podgorica** — Pričelje — Spuž — Ljutotuk — Danilovgrad — Slap — Bare Šumanovića — Šobajići — Ostrog — Dabovići — Stubica — **Nikšić**
- **Belgrade branch**: **Podgorica** — Zlatica — Bioče — Bratonožići — ... — Kolašin — ... — Mojkovac — ... — **Bijelo Polje** — ... — Prijepolje — Priboj — Užice — ... — **Beograd Centar** — ... — **Novi Sad** — ... — **Subotica**

---

*Generated 2026-04-12 alongside the prototype PR.*
