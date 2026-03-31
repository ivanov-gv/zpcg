# Seasons Feature — Implementation Plan

## Context

The Montenegro Railways timetable changes twice a year: a summer international train (Bar–Belgrade–Subotica) is added ~June 12 and removed ~September. Currently:
- The exporter fetches the timetable for a **single hardcoded date** — the binary knows only one season
- Summer-only stations (Novi Sad, Indjija, Stara Pazova, Nova Pazova, Subotica) are **permanently blacklisted** with a "coming June 12" message
- The owner must **manually re-run the parser** at midnight on season-change day and redeploy

**Goal:** Parse the timetable once for the whole year (fetching both winter and summer dates), compile all schedules into the binary, and at runtime serve the correct timetable based on today's date.

**Key decisions:**
- **Season boundaries are defined manually in code** (updated once per year), not derived from API validity dates
- **Summer-only stations queried outside summer** get a specific seasonal message (similar to current blacklist behavior, but date-aware)
- **Train validity** is still stored via `ValidFrom`/`ValidTo` from the API and used for runtime filtering

---

## Checklist

### Phase 1: Data Model Changes
- [ ] 1.1 — Add `ValidFrom`, `ValidTo` fields to `TrainInfo` in `internal/model/timetable/train.go`
- [ ] 1.2 — Create season config in a new file `internal/model/timetable/season.go` with manual date boundaries and seasonal station definitions
- [ ] 1.3 — Remove summer stations group from `internal/model/stations/blacklist.go`
- [ ] 1.4 — Add/uncomment summer station aliases in `internal/model/stations/alias.go` (if any exist)

### Phase 2: Exporter Changes
- [ ] 2.1 — Update `cmd/exporter/main.go` to fetch timetable for multiple dates (winter + summer)
- [ ] 2.2 — Update `internal/service/parser/routes/routes.go` deduplication: same train across dates → merge, keep more stops, widen validity range
- [ ] 2.3 — Update `internal/service/parser/parser.go` `MapTimetableToTransferFormat()` to preserve `ValidFrom`/`ValidTo` in `TrainInfo`
- [ ] 2.4 — Regenerate `gen/timetable/timetable.gen.go` with both seasons' data

### Phase 3: Runtime Filtering
- [ ] 3.1 — Update `PathFinder` constructor to accept `TrainIdToTrainInfoMap`
- [ ] 3.2 — Update `PathFinder.FindRoutes()` to accept `time.Time` and filter trains by `ValidFrom`/`ValidTo`
- [ ] 3.3 — Add `PathFinder.IsStationAvailable(stationId, date)` method — checks if any train serving the station is valid on the given date
- [ ] 3.4 — Create date-aware seasonal check in `internal/service/blacklist/` (or `app/`) using season config: if a station is seasonal and current date is outside its season → return seasonal message
- [ ] 3.5 — Update `App.GenerateRouteForStations()` to pass current date to pathfinder and perform seasonal availability check
- [ ] 3.6 — Update `App.NewApp()` to wire new dependencies

### Phase 4: Tests
- [ ] 4.1 — Unit test: `PathFinder` with date filtering — trains outside validity excluded
- [ ] 4.2 — Unit test: `PathFinder` with date filtering — trains within validity included
- [ ] 4.3 — Unit test: `PathFinder` boundary dates — ValidFrom and ValidTo days themselves are inclusive
- [ ] 4.4 — Unit test: `PathFinder` with zero/empty validity — always valid (backward compat)
- [ ] 4.5 — Unit test: `IsStationAvailable()` — summer station available in summer, unavailable in winter
- [ ] 4.6 — App test: query summer station in winter → seasonal message returned
- [ ] 4.7 — App test: query summer station in summer → valid routes returned
- [ ] 4.8 — Uncomment summer station entries in `TestNameClashing` in `internal/app/app_test.go`
- [ ] 4.9 — Update `TestBlackList` — summer stations no longer in permanent blacklist
- [ ] 4.10 — Ensure all existing tests pass: `make test_unit`
- [ ] 4.11 — Integration test: seasonal scenario via test Telegram server (if practical)

---

## Detailed Design

### 1. Data Model Changes

#### 1.1 — `internal/model/timetable/train.go`

Add validity dates to `TrainInfo`:
```go
type TrainInfo struct {
    TrainId      TrainId
    TimetableUrl string
    ValidFrom    time.Time  // new
    ValidTo      time.Time  // new
}
```

This is the core change. `TrainInfo` is stored in `ExportFormat.TrainIdToTrainInfoMap` and compiled into `gen/timetable/timetable.gen.go`. The export uses `%#v` formatting, and the template already imports `"time"`, so the generated code will include `ValidFrom`/`ValidTo` automatically.

#### 1.2 — New file: `internal/model/timetable/season.go`

Define season boundaries and seasonal stations:
```go
package timetable

import "time"

type Season struct {
    Start time.Time // inclusive
    End   time.Time // inclusive
}

// SummerSeason defines when the summer timetable is active.
// Update these dates annually when ZPCG publishes the new schedule.
var SummerSeason = Season{
    Start: time.Date(2026, time.June, 12, 0, 0, 0, 0, time.UTC),
    End:   time.Date(2026, time.September, 15, 0, 0, 0, 0, time.UTC),
}

// SummerOnlyStationNames lists station names that are only reachable
// during summer season. Used for seasonal unavailability messages.
var SummerOnlyStationNames = []string{
    "Novi Sad", "Indjija", "Stara Pazova", "Nova Pazova", "Subotica",
    // Novi Beograd intentionally excluded to not clash with "Beograd Centar"
}
```

The seasonal message text (multilingual) can live here too, or stay near the render/blacklist packages. It's the same message currently in `blacklist.go` for summer stations, just repurposed for date-aware checks.

#### 1.3 — `internal/model/stations/blacklist.go`

Remove the first entry (the `// summer season stations` group with Novi Sad, Indjija, etc.). All other blacklisted stations (Budva, Tivat, Kotor, etc.) remain — they genuinely have no train stations.

#### 1.4 — `internal/model/stations/alias.go`

Check if summer station aliases are commented out. If so, uncomment them so fuzzy name matching works for Novi Sad, Stara Pazova, etc.

---

### 2. Exporter Changes

#### 2.1 — `cmd/exporter/main.go`

Fetch for two dates:
```go
var Dates = []string{
    "2025-12-14",  // winter schedule
    "2026-07-01",  // summer schedule
}
```

Build route URLs for each date and combine them all into a single `parser.ParseTimetable(allUrls...)` call.

The route set for summer should include the Bar–Subotica route that doesn't exist in winter. The existing routes (Bar–Zemun, Podgorica–Bar, etc.) should be fetched for both dates to capture any schedule differences.

#### 2.2 — `internal/service/parser/routes/routes.go`

Current deduplication keeps the train entry with more stops. Update to handle the same `TrainNumber` appearing across different dates:
- If same train, same stops → keep one, but widen validity range: `ValidFrom = min(both)`, `ValidTo = max(both)`
- If same train, different stops → keep the one with more stops (existing behavior), inherit wider validity
- Different trains → no conflict

#### 2.3 — `internal/service/parser/parser.go`

In `MapTimetableToTransferFormat()`, update the `trainIdToTrainInfoMap` construction:
```go
trainIdToTrainInfoMap[trainId] = timetable.TrainInfo{
    TrainId:      trainId,
    TimetableUrl: route.TimetableUrl,
    ValidFrom:    route.ValidFrom,  // new
    ValidTo:      route.ValidTo,    // new
}
```

No other changes needed in this function — station maps, name maps, etc. will naturally include summer stations since those trains and stops are now in the parsed data.

#### 2.4 — Regenerate timetable

Run `make parse_timetable` (or `go run ./cmd/exporter -file gen/timetable/timetable.gen.go`). The generated file will now contain ~10-20 additional trains (summer schedule) with their `ValidFrom`/`ValidTo` dates. All summer stations appear as regular positive-ID stations.

---

### 3. Runtime Filtering

#### 3.1 — `internal/service/pathfinder/pathfinder.go` — Constructor

```go
func NewPathFinder(
    stationIdToTrainIdSetMap map[timetable.StationId]timetable.TrainIdSet,
    trainIdToStationsMap     map[timetable.TrainId]timetable.StationIdToStationMap,
    trainIdToTrainInfoMap    map[timetable.TrainId]timetable.TrainInfo,  // new
    transferStation          timetable.StationId,
) *PathFinder
```

Store `trainIdToTrainInfoMap` in the struct.

#### 3.2 — `pathfinder.go` — Date-filtered route finding

Update `FindRoutes` signature:
```go
func (p *PathFinder) FindRoutes(aStation, bStation timetable.StationId, currentDate time.Time) (routes []timetable.Path, isDirectRoute bool)
```

Pass `currentDate` through to `findDirectPaths`. In the loop over candidate trains, add:
```go
if !p.isTrainValidOnDate(trainId, currentDate) {
    continue
}
```

Helper method:
```go
func (p *PathFinder) isTrainValidOnDate(trainId timetable.TrainId, date time.Time) bool {
    info, ok := p.trainIdToTrainInfoMap[trainId]
    if !ok {
        return true // no info → always valid
    }
    if info.ValidFrom.IsZero() && info.ValidTo.IsZero() {
        return true // no constraints → always valid
    }
    // Compare dates only (ignore time component)
    d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
    return !d.Before(info.ValidFrom) && !d.After(info.ValidTo)
}
```

#### 3.3 — `pathfinder.go` — Station availability

```go
func (p *PathFinder) IsStationAvailable(stationId timetable.StationId, currentDate time.Time) bool {
    trainIds, ok := p.stationIdToTrainIdSetMap[stationId]
    if !ok {
        return false
    }
    for trainId := range trainIds {
        if p.isTrainValidOnDate(trainId, currentDate) {
            return true
        }
    }
    return false
}
```

#### 3.4 — Seasonal unavailability messages

Two options for where to put this logic:

**In `app/app.go`** (recommended for simplicity): After resolving station names and before the permanent-blacklist check, check if the station is seasonal and currently unavailable:

```go
// After name resolution, before permanent blacklist check:
if seasonMsg, isSeasonal := a.checkSeasonalAvailability(originStationId, destinationStationId, languageTag); isSeasonal {
    return seasonMsg, nil
}
```

The `checkSeasonalAvailability` method uses the season config to determine if a station is seasonal, and the pathfinder's `IsStationAvailable` to check if it has trains today. If seasonal and no trains → return the multilingual seasonal message.

The multilingual seasonal message text (currently in `blacklist.go` summer entry) moves to the season config or a companion file.

#### 3.5 — `internal/app/app.go` — Wire it together

In `GenerateRouteForStations`:
```go
// 1. Resolve names (existing)
// 2. Check permanent blacklist (existing — now without summer stations)
// 3. NEW: Check seasonal availability
// 4. Find routes with date (updated call)
routes, isDirect := a.finder.FindRoutes(originStationId, destinationStationId, a.dateService.CurrentDateAsTime())
```

#### 3.6 — `internal/app/app.go` — `NewApp()` constructor

Pass `TrainIdToTrainInfoMap` to `NewPathFinder`:
```go
finder := pathfinder.NewPathFinder(
    _timetable.StationIdToTrainIdSet,
    _timetable.TrainIdToStationMap,
    _timetable.TrainIdToTrainInfoMap,  // new
    _timetable.TransferStationId,
)
```

---

### 4. Test Plan

#### 4.1–4.4 — PathFinder unit tests (`internal/service/pathfinder/pathfinder_test.go`)

Create a small synthetic timetable with:
- Train A: `ValidFrom=2026-01-01, ValidTo=2026-12-31` (year-round)
- Train B: `ValidFrom=2026-06-12, ValidTo=2026-09-15` (summer only)
- Train C: zero validity dates (always valid)

Test cases:
- Winter date → only trains A and C returned
- Summer date → trains A, B, and C returned
- Boundary: June 12 → train B included
- Boundary: September 15 → train B included
- September 16 → train B excluded

#### 4.5 — `IsStationAvailable` tests

Using the same synthetic timetable:
- Station only served by summer train → available in summer, unavailable in winter
- Station served by year-round train → always available

#### 4.6–4.7 — App seasonal tests (`internal/app/app_test.go`)

These need the app's date to be controllable. Options:
- Add `NewAppWithDate(date time.Time)` or `App.SetDate(date)` for testing
- Or make `DateService` injectable with a fixed-date variant

Test:
- Create app with winter date → query "Bar, Novi Sad" → seasonal message
- Create app with summer date → query "Bar, Novi Sad" → valid routes

#### 4.8 — Uncomment summer stations in `TestNameClashing`

The commented-out entries for Novi Beograd, Stara Pazova, Nova Pazova, Novi Sad become active. These test that fuzzy name matching works for summer stations.

#### 4.9 — Update `TestBlackList`

Summer stations are no longer in `blacklist.BlackListedStations`, so the test loop naturally skips them. Verify the test still passes with the reduced blacklist.

#### 4.10 — Full test pass: `make test_unit`

#### 4.11 — Integration test (stretch goal)

Add a test case in `test/integration/telegram_test.go` that sends a message like "Novi Sad, Bar" and validates the response is either a seasonal message or valid timetable depending on the configured test date.

---

## Files to Modify (Summary)

| File | Change |
|------|--------|
| `internal/model/timetable/train.go` | Add `ValidFrom`, `ValidTo` to `TrainInfo` |
| `internal/model/timetable/season.go` | **New file** — season config, seasonal station list, seasonal messages |
| `internal/model/stations/blacklist.go` | Remove summer stations group |
| `internal/model/stations/alias.go` | Uncomment/add summer station aliases |
| `cmd/exporter/main.go` | Multiple dates, expanded route URLs |
| `internal/service/parser/routes/routes.go` | Merge validity ranges for same train across dates |
| `internal/service/parser/parser.go` | Preserve `ValidFrom`/`ValidTo` in `TrainInfo` |
| `internal/service/pathfinder/pathfinder.go` | Add date filtering, train info map, `IsStationAvailable` |
| `internal/app/app.go` | Pass date to pathfinder, seasonal availability check, wire new deps |
| `internal/service/pathfinder/pathfinder_test.go` | Date-filtered pathfinding tests |
| `internal/app/app_test.go` | Seasonal behavior tests, uncomment summer stations |
| `gen/timetable/timetable.gen.go` | Regenerated with both seasons |

## Implementation Order

1. **Phase 1** first — safe, no behavioral changes until the exporter and runtime are updated
2. **Phase 2** — exporter changes can be tested in isolation by running the exporter and inspecting output
3. **Phase 3** — this is where behavior changes; do runtime filtering + app wiring together
4. **Phase 4** — tests throughout, but especially after Phase 3

## Verification

1. `make test_unit` — all tests pass
2. Query winter-only route (Nikšić → Bar) → works year-round, no regression
3. Query summer station (Bar → Novi Sad) with winter date → seasonal message
4. Query summer station (Bar → Novi Sad) with summer date → valid routes
5. Permanent blacklist stations (Budva, Tivat, etc.) → still show "no train station" messages
