# Seasons Feature — Implementation Plan

## Context

The Montenegro Railways timetable changes twice a year: a summer international train (Bar–Belgrade–Subotica) is added ~June 12 and removed ~September. Currently:
- The exporter fetches the timetable for a **single hardcoded date** — the binary knows only one season
- Summer-only stations (Novi Sad, Indjija, Stara Pazova, Nova Pazova, Subotica) are **permanently blacklisted** with a "coming June 12" message
- The owner must **manually re-run the parser** at midnight on season-change day and redeploy

**Goal:** Parse the timetable once for the whole year (fetching multiple dates), compile all schedules into the binary, and at runtime serve the correct timetable based on today's date.

**Key decisions:**
- **Parsing configuration comes from a config file** (dates, routes, season boundaries) — not hardcoded
- **Train seasonality is derived by comparing fetch sets** — NOT from API's `ValidFrom`/`ValidTo` fields (the API is unreliable except for the route endpoint)
- **Same route URLs are fetched for each date** — summer-only routes simply return empty for winter dates
- **Seasons form a continuous timeline** — no gaps, every date belongs to exactly one season (from -∞ to +∞)
- **Validity is a proper struct** with `IsEmpty()`, setter, getter — impossible to have one date set without the other
- **Summer-only stations queried outside summer** get a specific seasonal message (date-aware, multilingual)
- **Integration tests are mandatory** for season switching scenarios

---

## Core Algorithm

### Parsing

Given a config file with 3 seasons and 3 fetch dates:

```
Season 1 "winter-early":  -∞  to  2026-06-11  (fetch with date 2026-05-25)
Season 2 "summer":        2026-06-12  to  2026-09-15  (fetch with date 2026-06-25)
Season 3 "winter-late":   2026-09-16  to  +∞  (fetch with date 2026-11-25)
```

1. For each season, fetch all routes using that season's `fetch_date`
2. Each fetch returns a set of trains (identified by `TrainNumber`)
3. Classify each train by which seasons returned it:
   - Train in all 3 → **year-round** → Validity is empty (= always valid)
   - Train only in season 2 → **summer-only** → Validity: `[{Jun 12, Sep 15}]`
   - Train in seasons 1+3 but not 2 → **winter-only** → Validity: `[{-∞, Jun 11}, {Sep 16, +∞}]` (two ranges)
   - Train only in season 1 → Validity: `[{-∞, Jun 11}]`
4. Merge all trains into one timetable, each `TrainInfo` carrying its validity ranges

### Runtime

For a user query:
1. Resolve station names
2. Check permanent blacklist (Budva, Tivat, etc.)
3. Check seasonal availability — is either station seasonal and currently unavailable?
4. `PathFinder.FindRoutes(origin, dest, today)` — filters trains where `today` falls within at least one validity range
5. Render timetable

---

## Checklist

### Phase 1: Data Model Changes
- [ ] 1.1 — Create `Validity` struct in `internal/model/timetable/validity.go` with `IsEmpty() bool`, `Set(from, to)`, `Get() []ValidityRange`, and `Contains(date) bool`
- [ ] 1.2 — Add `Validity` field to `TrainInfo` in `internal/model/timetable/train.go`
- [ ] 1.3 — Create season config types in `internal/model/timetable/season.go` (Season, SeasonConfig)
- [ ] 1.4 — Move summer stations seasonal message from `internal/model/stations/blacklist.go` to season config; remove summer stations from permanent blacklist
- [ ] 1.5 — Uncomment/add summer station aliases in `internal/model/stations/alias.go`

### Phase 2: Parser Config File
- [ ] 2.1 — Define config file format (YAML or JSON) with seasons, fetch dates, route templates, summer station names
- [ ] 2.2 — Create config parser in `internal/service/parser/config.go`
- [ ] 2.3 — Create the config file (e.g. `cmd/exporter/config.yaml`)
- [ ] 2.4 — Update `cmd/exporter/main.go` to read config file instead of hardcoded dates/routes

### Phase 3: Exporter Multi-Season Parsing
- [ ] 3.1 — Update `parser.ParseTimetable()` to accept season configs, fetch routes per season
- [ ] 3.2 — Implement train set comparison: for each train, determine which seasons returned it
- [ ] 3.3 — Assign validity ranges to each train based on its season presence
- [ ] 3.4 — Handle deduplication: same train across seasons → keep version with more stops
- [ ] 3.5 — Update `MapTimetableToTransferFormat()` to preserve `Validity` in `TrainInfo`
- [ ] 3.6 — Regenerate `gen/timetable/timetable.gen.go` with all seasons' data

### Phase 4: Runtime Filtering
- [ ] 4.1 — Update `PathFinder` constructor to accept `TrainIdToTrainInfoMap`
- [ ] 4.2 — Update `PathFinder.FindRoutes()` to accept `time.Time` and filter trains by `Validity.Contains(date)`
- [ ] 4.3 — Add `PathFinder.IsStationAvailable(stationId, date)` method
- [ ] 4.4 — Add seasonal availability check in `App.GenerateRouteForStations()` using season config
- [ ] 4.5 — Update `App.NewApp()` to wire new dependencies
- [ ] 4.6 — Make `App` date-testable: allow injecting a fixed date for testing

### Phase 5: Tests (mandatory)
- [ ] 5.1 — Unit test: `Validity` struct — `IsEmpty()`, `Contains()`, `Set()`, edge cases
- [ ] 5.2 — Unit test: `PathFinder` date filtering — summer train excluded in winter
- [ ] 5.3 — Unit test: `PathFinder` date filtering — summer train included in summer
- [ ] 5.4 — Unit test: `PathFinder` boundary dates — season start/end days are inclusive
- [ ] 5.5 — Unit test: `PathFinder` empty validity — always valid (year-round trains)
- [ ] 5.6 — Unit test: `PathFinder` non-contiguous validity — winter train valid in seasons 1+3
- [ ] 5.7 — Unit test: `IsStationAvailable()` — summer station available/unavailable by date
- [ ] 5.8 — Unit test: config file parsing — valid config, missing fields, malformed dates
- [ ] 5.9 — Unit test: train set comparison — correct season classification
- [ ] 5.10 — App test: query summer station in winter → seasonal message
- [ ] 5.11 — App test: query summer station in summer → valid routes
- [ ] 5.12 — Uncomment summer station entries in `TestNameClashing`
- [ ] 5.13 — Update `TestBlackList` — summer stations no longer in permanent blacklist
- [ ] 5.14 — Ensure all existing tests pass: `make test_unit`
- [ ] 5.15 — Integration test: winter scenario — send "Novi Sad, Bar" → seasonal message via test Telegram server
- [ ] 5.16 — Integration test: summer scenario — send "Novi Sad, Bar" → valid timetable via test Telegram server
- [ ] 5.17 — Integration test: year-round route unaffected — send "Nikšić, Bar" → valid timetable in any season
- [ ] 5.18 — Integration test: season boundary — test on exact switch date

---

## Detailed Design

### 1. Validity Struct

**New file: `internal/model/timetable/validity.go`**

```go
package timetable

import "time"

// ValidityRange represents a date range [From, To] inclusive.
// Zero From means -infinity (no start constraint).
// Zero To means +infinity (no end constraint).
type ValidityRange struct {
    from time.Time
    to   time.Time
}

// Validity represents when a train is valid to run.
// Empty ranges means the train is always valid (year-round).
type Validity struct {
    ranges []ValidityRange
}

// NewValidity creates a Validity from one or more ranges.
// Both from and to must be zero or both must be non-zero within each range.
func NewValidity(ranges ...ValidityRange) Validity { ... }

// NewValidityRange creates a range. Pass time.Time{} for -∞ or +∞.
func NewValidityRange(from, to time.Time) ValidityRange { ... }

// IsEmpty returns true if the train is year-round (no restrictions).
func (v Validity) IsEmpty() bool { return len(v.ranges) == 0 }

// Contains returns true if the given date falls within any validity range.
func (v Validity) Contains(date time.Time) bool {
    if v.IsEmpty() {
        return true // year-round
    }
    d := truncateToDate(date)
    for _, r := range v.ranges {
        fromOk := r.from.IsZero() || !d.Before(r.from)
        toOk := r.to.IsZero() || !d.After(r.to)
        if fromOk && toOk {
            return true
        }
    }
    return false
}

// Ranges returns the validity ranges (read-only access).
func (v Validity) Ranges() []ValidityRange { return v.ranges }
```

The struct is immutable after construction — `ranges` is unexported. `NewValidityRange` validates that both dates are either zero or non-zero (prevents one-sided ranges). However: -∞ or +∞ is represented by a zero `time.Time`, which IS valid for the first/last season.

**Clarification on -∞ / +∞ semantics:** A `ValidityRange{from: time.Time{}, to: someDate}` means "valid from the beginning of time until someDate". A `ValidityRange{from: someDate, to: time.Time{}}` means "valid from someDate forever". Both `from` and `to` being zero simultaneously is not allowed (use `IsEmpty()` for year-round instead).

### 2. Config File

**New file: `cmd/exporter/config.yaml`**

```yaml
seasons:
  - name: winter-early
    end: "2026-06-11"
    fetch_date: "2026-05-25"
  - name: summer
    start: "2026-06-12"
    end: "2026-09-15"
    fetch_date: "2026-06-25"
  - name: winter-late
    start: "2026-09-16"
    fetch_date: "2026-11-25"

routes:
  - start: Bar
    finish: Subotica
  - start: Subotica
    finish: Bar
  - start: Bar
    finish: Zemun
  - start: Zemun
    finish: Bar
  - start: Podgorica
    finish: Bar
  - start: Bar
    finish: Podgorica
  - start: Podgorica
    finish: "Nikšić"
  - start: "Nikšić"
    finish: Podgorica
  - start: Bar
    finish: Bijelo Polje
  - start: Bijelo Polje
    finish: Bar

summer_only_stations:
  - Novi Sad
  - Indjija
  - Stara Pazova
  - Nova Pazova
  - Subotica
```

**Rules:**
- First season has no `start` (= -∞), last season has no `end` (= +∞)
- Each season's `end` + 1 day = next season's `start` (continuous, no gaps)
- Each season has a `fetch_date` that falls within its date range
- Config parser validates continuity and completeness

**New file: `internal/service/parser/config.go`**
- Parses YAML config
- Validates season continuity (no gaps, no overlaps)
- Validates fetch dates fall within their season range
- Returns typed `ParserConfig` struct

### 3. Exporter Changes

**`cmd/exporter/main.go`**
- Read config file path from `-config` flag (default: `config.yaml`)
- Parse config
- For each season: build route URLs as `fmt.Sprintf("/routes?start=%s&finish=%s&date=%s", route.Start, route.Finish, season.FetchDate)`
- Call updated `parser.ParseTimetableMultiSeason()`

**`internal/service/parser/parser.go`**

New function:
```go
func ParseTimetableMultiSeason(config ParserConfig) (timetable.ExportFormat, error) {
    // 1. Parse stations (once)
    // 2. For each season: fetch routes → map[TrainId]DetailedTimetable
    // 3. Build presence map: trainId → set of season indices
    // 4. For each train:
    //    - Present in all seasons → Validity{} (empty = year-round)
    //    - Otherwise → Validity with ranges from the seasons where it appears
    //    - Pick version with most stops for the train data
    // 5. Merge all trains → single map
    // 6. MapTimetableToTransferFormat (existing, with Validity)
    // 7. Add blacklist + aliases (existing)
}
```

**Train set comparison example with 3 seasons:**

| Train | S1 (winter) | S2 (summer) | S3 (winter) | Validity |
|-------|:-----------:|:-----------:|:-----------:|----------|
| 6100  | ✓ | ✓ | ✓ | empty (year-round) |
| 9001  | ✗ | ✓ | ✗ | `[{Jun 12, Sep 15}]` |
| 432   | ✓ | ✗ | ✓ | `[{-∞, Jun 11}, {Sep 16, +∞}]` |

### 4. Runtime Changes

**`internal/service/pathfinder/pathfinder.go`**

Add `trainIdToTrainInfoMap` to struct. Update signatures:
```go
func NewPathFinder(
    stationIdToTrainIdSetMap map[timetable.StationId]timetable.TrainIdSet,
    trainIdToStationsMap     map[timetable.TrainId]timetable.StationIdToStationMap,
    trainIdToTrainInfoMap    map[timetable.TrainId]timetable.TrainInfo,  // new
    transferStation          timetable.StationId,
) *PathFinder

func (p *PathFinder) FindRoutes(aStation, bStation timetable.StationId, currentDate time.Time) ([]timetable.Path, bool)

func (p *PathFinder) IsStationAvailable(stationId timetable.StationId, currentDate time.Time) bool
```

In `findDirectPaths`, after building candidate set, filter:
```go
if !p.trainIdToTrainInfoMap[trainId].Validity.Contains(currentDate) {
    continue
}
```

**`internal/app/app.go`**

In `GenerateRouteForStations`, after name resolution and permanent blacklist check:
```go
// Check seasonal availability
if msg, unavailable := a.checkSeasonalAvailability(originStationId, destinationStationId, languageTag); unavailable {
    return msg, nil
}
// Find routes with date filtering
routes, isDirect := a.finder.FindRoutes(originStationId, destinationStationId, a.dateService.CurrentDateAsTime())
```

Date injection for testing — add options pattern or `NewAppWithDate()`:
```go
func NewApp(opts ...AppOption) (*App, error)

type AppOption func(*App)

func WithFixedDate(date time.Time) AppOption {
    return func(a *App) { a.dateService = newFixedDateService(date) }
}
```

### 5. Integration Tests

**`test/integration/telegram_test.go`**

All scenarios use the real HTTP server with mocked Telegram client, injecting a fixed date:

**5.15 — Winter scenario:**
- App date: 2026-01-15
- Send: `{"message": {"text": "Novi Sad, Bar", ...}}`
- Mock expects `sendMessage` call
- Assert response text contains seasonal unavailability message (multilingual)
- Assert response does NOT contain train departure times

**5.16 — Summer scenario:**
- App date: 2026-07-15
- Send: `{"message": {"text": "Novi Sad, Bar", ...}}`
- Mock expects `sendMessage` call
- Assert response text contains train IDs and times
- Assert response does NOT contain seasonal message

**5.17 — Year-round route unaffected:**
- App date: 2026-01-15 AND 2026-07-15
- Send: `{"message": {"text": "Nikšić, Bar", ...}}`
- Assert both return valid timetable with trains

**5.18 — Season boundary:**
- App date: 2026-06-11 → "Novi Sad, Bar" → seasonal message
- App date: 2026-06-12 → "Novi Sad, Bar" → valid timetable

---

## Files to Modify / Create

| File | Change |
|------|--------|
| `internal/model/timetable/validity.go` | **New** — `Validity`, `ValidityRange` structs |
| `internal/model/timetable/validity_test.go` | **New** — Validity unit tests |
| `internal/model/timetable/train.go` | Add `Validity` field to `TrainInfo` |
| `internal/model/timetable/season.go` | **New** — season config types |
| `internal/model/stations/blacklist.go` | Remove summer stations group |
| `internal/model/stations/alias.go` | Uncomment summer station aliases |
| `cmd/exporter/config.yaml` | **New** — parser config file |
| `cmd/exporter/main.go` | Read config file, multi-date fetching |
| `internal/service/parser/config.go` | **New** — config file parser + validation |
| `internal/service/parser/config_test.go` | **New** — config parsing tests |
| `internal/service/parser/parser.go` | Multi-season parsing, train classification |
| `internal/service/pathfinder/pathfinder.go` | Date filtering, `IsStationAvailable` |
| `internal/service/pathfinder/pathfinder_test.go` | Date-filtered tests |
| `internal/app/app.go` | Date wiring, seasonal check, date injection |
| `internal/app/app_test.go` | Seasonal behavior tests |
| `test/integration/telegram_test.go` | Season switching integration tests |
| `gen/timetable/timetable.gen.go` | Regenerated with all seasons |

## Implementation Order

1. **Phase 1** — Data model (Validity struct, season config types, blacklist cleanup)
2. **Phase 2** — Config file format + parser
3. **Phase 3** — Exporter multi-season parsing
4. **Phase 4** — Runtime filtering in PathFinder + App
5. **Phase 5** — Tests throughout, integration tests last

## Verification

1. `make test_unit` — all tests pass
2. `make test_integration` — season switching works end-to-end
3. Winter route (Nikšić → Bar) → works in any season
4. Summer station (Bar → Novi Sad) in winter → seasonal message
5. Summer station (Bar → Novi Sad) in summer → valid routes
6. Season boundary (Jun 11 vs Jun 12) → correct behavior
7. Permanent blacklist (Budva, Tivat) → unchanged
