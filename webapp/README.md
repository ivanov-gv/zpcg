# Monterails Mini App · prototype

A **Telegram Mini App** prototype for the [@Monterails_bot](https://t.me/Monterails_bot)
Montenegro Railways timetable bot. This directory contains a pure-static web
app (HTML + vanilla JS + Leaflet) that can be hosted on GitHub Pages and
opened via the bot using Telegram's WebApp feature.

The full design rationale, feature roadmap, architectural trade-offs and
open questions live in [`docs/miniapp/DESIGN.md`](../docs/miniapp/DESIGN.md).

---

## What's in here

| Path | What it is |
| --- | --- |
| `index.html`              | Single-page shell, Telegram SDK + Leaflet + stylesheet |
| `css/app.css`             | Mobile-first styles; reads Telegram theme CSS variables |
| `js/app.js`               | Entry point, views, state, helpers |
| `js/router.js`            | ~40-line hash router |
| `js/data.js`              | Loads `data/*.json` and indexes it |
| `js/pathfinder.js`        | JS port of `internal/service/pathfinder` (direct + 1-transfer) |
| `js/price.js`             | Rough €/km ticket estimator (not authoritative) |
| `js/map.js`               | Leaflet layer + route polyline |
| `js/i18n.js`              | EN / RU / SR strings; picks Telegram user language |
| `data/stations.json`      | 75 stations, generated from `gen/timetable/timetable.gen.go` |
| `data/trains.json`        | 30 trains with ordered stops (HH:MM) |
| `data/meta.json`          | `transferStationId`, `generatedAt`, station types |
| `data/coordinates.json`   | Hand-curated approximate lat/lng per station |

Everything under `data/` (except `coordinates.json`) is rebuilt by:

```bash
go run ./cmd/webapp-export -out=webapp/data
```

Re-run this after every `make parse_timetable` so the mini app stays in
sync with the bot.

## Run it locally

```bash
# from repo root
cd webapp
python3 -m http.server 8080
# then open http://localhost:8080/
```

or

```bash
npx http-server webapp -p 8080
```

`file://` won't work — the browser blocks `fetch()` for local files, so you
do need a real HTTP server.

## Run it inside Telegram

1. Deploy the `webapp/` folder to GitHub Pages (or any HTTPS host).
2. In [@BotFather](https://t.me/BotFather) → your bot → *Bot Settings* →
   *Menu Button* → *Configure menu button* → paste the Pages URL.
3. In the bot, tapping the new button opens the mini app with
   `Telegram.WebApp.initDataUnsafe.user.language_code` pre-set — the app
   picks up the language automatically.

No server changes needed for a first slice: the app talks only to its own
static JSON. Once the bot gets an HTTP endpoint for live status, the app
can fetch it directly (see `DESIGN.md § Architecture options`).

## What works today

- [x] Route search with "From / To" pickers (latin + cyrillic + diacritic-loose match)
- [x] Direct and 1-transfer path-finding, same algorithm as the Go bot
- [x] Train details with highlighted departure/arrival stops
- [x] Per-station view with the list of trains serving it
- [x] Stations map (Leaflet + OSM) with route polyline overlay
- [x] Rough price estimator with class / passenger type / count
- [x] Popular-routes shortcuts
- [x] Alerts and saved-routes screens (LocalStorage only — no push yet)
- [x] EN / RU / SR UI strings, language auto-detected from Telegram
- [x] Telegram theme variables piped into CSS (light + dark)
- [x] Back button wired to `history.back()`, tab bar, hash router

## What is mocked (on purpose)

- **Station coordinates** — only main stations are real; minor halts are
  approximated along the corridor. See note at the top of
  `data/coordinates.json`.
- **Price estimator** — €0.06/km flat; UI marks it as an estimate.
- **Alerts** — stored in LocalStorage, not pushed through the bot yet.
- **Live disruptions** — hard-coded example card on the Alerts screen.

## Browser support

Modern evergreen browsers only (ES modules, `fetch`, template literals).
Telegram's in-app browser is Chromium-based and covers this.
