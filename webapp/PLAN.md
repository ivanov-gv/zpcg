# PLAN: Ship the Monterails Mini App

## What we're building

A **Telegram Mini App** for the Montenegro Railways timetable bot (`@Monterails_bot`).
It is a pure-JS single-page web app that runs inside Telegram's built-in browser.
Users pick two stations, see train times, browse a map, estimate ticket prices — all
without leaving Telegram.

The app **must work offline** (trains go through mountain tunnels with no signal).
It deploys to **GitHub Pages** automatically when code is pushed.

## What already exists

Everything lives in `webapp/`:

```
webapp/
  index.html              ← single HTML shell
  css/app.css             ← styles (Telegram theme-aware)
  js/
    app.js                ← main entry: 9 views, state, helpers
    router.js             ← hash-based page router
    data.js               ← loads JSON, builds indexes, search
    pathfinder.js         ← finds direct + 1-transfer routes
    price.js              ← rough €/km ticket estimate
    map.js                ← Leaflet map with station markers
    i18n.js               ← EN / RU / SR translations
  data/
    stations.json         ← 75 stations (generated from Go timetable)
    trains.json           ← 30 trains with stop times
    meta.json             ← transfer station id, station types
    coordinates.json      ← lat/lng per station (approximate)
```

The prototype works in a browser today. What's missing: one bug fix, error
handling, offline support (Service Worker), and a deploy pipeline.

## Rules

- **Only touch files inside `webapp/` and `.github/workflows/`.**
- **Do NOT touch any Go code** (`cmd/`, `internal/`, `gen/`, etc.).
- No build tools, no npm, no bundler. Everything stays as plain files.
- Use ES modules (`type="module"`). Modern browsers only (Telegram WebView is Chromium).

---

## TODO

### Task 1 — Fix overnight train bug

**What's wrong:**
`webapp/js/pathfinder.js`, the `#directRoutes` method (lines 53–75).

Train 432 goes from Bijelo Polje (departs 23:20) to Belgrade (arrives 07:08 the next day).
The current code on line 64 does:
```js
if (hhmmToMin(fromStop.departure) >= hhmmToMin(toStop.arrival)) continue;
```
This compares 1400 >= 428 → true → skips the route. But the route IS valid — it just
crosses midnight.

**Why it happens:**
The stops array in `trains.json` is sorted by departure time as a string.
For overnight trains, "00:10" sorts before "23:20". So the array order is:
`[..stations after midnight.., ..stations before midnight..]`.
For normal daytime trains, array order = travel order. For overnight trains, it wraps.

**How to fix:**
Replace the `#directRoutes` method. Instead of comparing times only, also use the
array index position of each stop:

```
fromIdx = index of departure station in train.stops
toIdx   = index of arrival station in train.stops

if fromIdx < toIdx → normal (daytime) train
  valid when depMin < arrMin

if fromIdx > toIdx → overnight train (array wrapped at midnight)
  valid when depMin > arrMin (departs in evening, arrives after midnight)
```

For duration: `arrMin - depMin`. If negative, add `24 * 60`.

**Also fix these related spots:**
1. `#mergeInterchange` (same file, lines 94–95): `totalMin` and `waitMin` can go
   negative for overnight combos. Add `if (x < 0) x += 24 * 60`.
2. `webapp/js/app.js`, function `computeMapRouteIds` (lines 589–612): the
   `[lo, hi]` slice doesn't work for overnight trains where `fromIdx > toIdx`.
   Fix: when `fromIdx > toIdx`, return `stops[fromIdx..end] + stops[0..toIdx+1]`.

**How to check it's done:**
1. Start a local HTTP server: `cd webapp && python3 -m http.server 8080` (or any server)
2. Open the app in a browser
3. Search **"Podgorica → Beograd Centar"** → train 432 MUST appear in results
4. Search **"Bar → Podgorica"** → trains 6100–6105 MUST still appear (not broken)
5. Search **"Bar → Beograd Centar"** → train 432 MUST appear
6. Click train 432 → stop list should show all stations in travel order
7. No JavaScript errors in the browser console

---

### Task 2 — Add Leaflet error handling

**What's wrong:**
`webapp/js/map.js`, the `mount()` method (line 19).

If the Leaflet CDN (`unpkg.com`) is unreachable (user is offline and the Service Worker
hasn't cached it yet), the global `L` is undefined. Calling `L.map()` throws a
ReferenceError and the entire Map tab crashes with a blank screen.

**How to fix:**
1. At the start of `mount()`, check `typeof L === 'undefined'`. If true, render a
   friendly fallback message in the container element ("Map unavailable offline") and return.
2. Wrap the rest of `mount()` in a `try { ... } catch` that shows "Could not load map".
3. Add `if (!this.map) return;` at the top of `showRoute()`, `highlightStations()`, and
   `flyToStation()` so they don't crash if mount failed.
4. In `webapp/js/i18n.js`, add two new keys to all 3 language dictionaries:
   - `"map.offline"` — EN: "Map unavailable offline", RU: "Карта недоступна офлайн",
     SR: "Mapa nedostupna offline"
   - `"map.error"` — EN: "Could not load map", RU: "Не удалось загрузить карту",
     SR: "Nije moguće učitati mapu"

**How to check it's done:**
1. Open the app normally → Map tab works with Leaflet, stations visible ✓
2. In Chrome DevTools → Network tab → right-click `unpkg.com` → "Block request domain"
3. Hard-reload the page (Ctrl+Shift+R)
4. Open the Map tab → should show a friendly message, NOT a crash or blank screen ✓
5. No JavaScript errors in the console ✓
6. Unblock `unpkg.com`, reload → map works again ✓

---

### Task 3 — Create the Service Worker

**What to create:**
New file: `webapp/sw.js` (in the same folder as `index.html`).

**What it must do:**

1. **On install** — open a cache named `monterails-v1` and add every local file:
   ```
   ./
   ./index.html
   ./css/app.css
   ./js/app.js
   ./js/router.js
   ./js/data.js
   ./js/pathfinder.js
   ./js/price.js
   ./js/map.js
   ./js/i18n.js
   ./data/stations.json
   ./data/trains.json
   ./data/meta.json
   ./data/coordinates.json
   ```
   Also try to cache these CDN files (but don't fail if they're unreachable):
   ```
   https://telegram.org/js/telegram-web-app.js
   https://unpkg.com/leaflet@1.9.4/dist/leaflet.css
   https://unpkg.com/leaflet@1.9.4/dist/leaflet.js
   ```
   Call `self.skipWaiting()` so the new SW activates immediately.

2. **On activate** — delete any old caches whose name is not `monterails-v1`.
   Call `self.clients.claim()` so the SW takes control of the page right away.

3. **On fetch** — cache-first strategy:
   - Check the cache. If found, return it.
   - If not in cache, fetch from network.
   - If the network response is OK and is a GET request, put a copy in the cache
     (this caches map tiles on the fly).
   - Return the response.

**How to update the app later:**
Change `monterails-v1` to `monterails-v2` (or any new string). On next visit, the
browser detects the new SW, installs a new cache, and deletes the old one.

**How to check it's done:**
1. Open the app in Chrome
2. Go to DevTools → Application → Service Workers
3. The SW should show status **"activated and is running"** ✓
4. Go to Application → Cache Storage → `monterails-v1` should exist with all files listed ✓
5. In DevTools → Network tab → check the **"Offline"** checkbox
6. Reload the page → the app MUST still load and work ✓
7. Search "Bar → Podgorica" while offline → results MUST appear ✓
8. Map tab while offline → should show the fallback message from Task 2 (no tiles) ✓
9. Uncheck "Offline", reload, open Map tab → map tiles load from network ✓

---

### Task 4 — Register the Service Worker

**What to change:**
`webapp/js/app.js`, inside the `init()` function, right before the `router.handle()` call
(currently line 111).

**What to add:**
```js
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('./sw.js').catch(err =>
    console.warn('SW registration failed:', err)
  );
}
```

This is non-blocking. The app doesn't wait for the SW to install before rendering.

**How to check it's done:**
1. Reload the page
2. DevTools → Application → Service Workers → SW is listed and active ✓
3. No errors in the console about SW registration ✓
4. The app works exactly as before (registration doesn't break anything) ✓

---

### Task 5 — Add GitHub Pages deployment workflow

**What to create:**
New file: `.github/workflows/deploy-pages.yml`

**Contents:**
```yaml
name: Deploy Mini App to GitHub Pages

on:
  push:
    branches: [main]
    paths: ['webapp/**']
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/configure-pages@v5
      - uses: actions/upload-pages-artifact@v3
        with:
          path: webapp
      - id: deployment
        uses: actions/deploy-pages@v4
```

**What this does:**
- Triggers when files under `webapp/` are pushed to `main`
- Also has a manual trigger button (`workflow_dispatch`)
- Uploads only the `webapp/` folder (not the whole repo) to GitHub Pages
- Uses the standard GitHub-provided actions (no custom scripts)

**How to check it's done:**
1. The file exists at `.github/workflows/deploy-pages.yml` ✓
2. The YAML is valid (no syntax errors) ✓
3. After merging to main, the workflow appears in the repo's Actions tab ✓
4. Pushing a change to `webapp/index.html` triggers the workflow ✓
5. Pushing a change to a Go file (e.g. `internal/`) does NOT trigger it ✓

---

## Execution order

```
Task 1  →  Task 2  →  Task 3  →  Task 4  →  Task 5
  │            │          │          │          │
  │            │          │          └── needs sw.js from Task 3
  │            │          └── do after 1+2 so cached JS has fixes
  │            └── independent
  └── most critical, do first

Task 5 can be done in parallel with any of the others.
```

## Final end-to-end check

After all 5 tasks are done:

1. Serve the app: `cd webapp && python3 -m http.server 8080`
2. Open `http://localhost:8080/` in Chrome
3. Search "Bar → Podgorica" → results appear ✓
4. Search "Podgorica → Beograd Centar" → train 432 appears ✓ (Task 1)
5. Open Map tab → Leaflet renders with station markers ✓
6. DevTools → Application → Service Workers → "activated" ✓ (Tasks 3+4)
7. DevTools → Network → check "Offline" → reload → app still works ✓ (Task 3)
8. Map tab while offline → friendly fallback, no crash ✓ (Task 2)
9. `.github/workflows/deploy-pages.yml` exists and is valid YAML ✓ (Task 5)
10. No JavaScript errors in the console at any point ✓
