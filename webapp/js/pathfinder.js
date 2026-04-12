// pathfinder.js — JS port of the Go path-finding algorithm
// (internal/service/pathfinder/pathfinder.go).
//
// Assumptions carried over from the Go implementation:
//   1. Each train visits every station in its route at most once.
//   2. Podgorica is the only transfer station; for any pair of stations,
//      either a direct route exists OR a route with exactly one transfer
//      in Podgorica does.
//   3. Direct routes are always preferred.
//
// This port keeps the same public contract: given (fromId, toId) return
// either a list of direct routes, or a list of { leg1, leg2 } interchange
// routes through Podgorica.

export class PathFinder {
  constructor(db) { this.db = db; }

  /**
   * @param {number} fromId
   * @param {number} toId
   * @returns {{ kind: 'direct'|'interchange'|'none',
   *             direct?: DirectRoute[],
   *             interchange?: InterchangeRoute[] }}
   */
  find(fromId, toId) {
    if (fromId === toId) return { kind: 'none' };

    const direct = this.#directRoutes(fromId, toId);
    if (direct.length) {
      direct.sort((a, b) => hhmmToMin(a.departure) - hhmmToMin(b.departure));
      return { kind: 'direct', direct };
    }

    const transferId = this.db.transferStationId;
    if (fromId === transferId || toId === transferId) {
      // no direct *and* one of the endpoints is the transfer → nothing to do
      return { kind: 'none' };
    }

    const legA = this.#directRoutes(fromId, transferId);
    const legB = this.#directRoutes(transferId, toId);
    if (!legA.length || !legB.length) return { kind: 'none' };

    const merged = this.#mergeInterchange(legA, legB);
    if (!merged.length) return { kind: 'none' };
    return { kind: 'interchange', interchange: merged };
  }

  /**
   * Finds every train that stops at both stations, where the first stop
   * precedes the second (i.e. the train actually goes A → B, not B → A).
   */
  #directRoutes(fromId, toId) {
    const trainsA = this.db.stationToTrainIds.get(fromId) || new Set();
    const trainsB = this.db.stationToTrainIds.get(toId)   || new Set();
    const routes = [];
    for (const trainId of trainsA) {
      if (!trainsB.has(trainId)) continue;
      const train = this.db.train(trainId);
      const fromStop = train.stops.find(s => s.stationId === fromId);
      const toStop   = train.stops.find(s => s.stationId === toId);
      if (!fromStop || !toStop) continue;
      // Sort order of train.stops is departure-ascending — compare directly.
      if (hhmmToMin(fromStop.departure) >= hhmmToMin(toStop.arrival)) continue;
      routes.push({
        trainId,
        fromStationId: fromId,
        toStationId:   toId,
        departure:     fromStop.departure,
        arrival:       toStop.arrival,
        durationMin:   hhmmToMin(toStop.arrival) - hhmmToMin(fromStop.departure),
      });
    }
    return routes;
  }

  /**
   * Pairs each leg-A train with the soonest-arriving leg-B train that
   * departs the transfer station after leg-A arrives (with at least 1
   * minute of slack — we keep it loose so tight connections still show).
   */
  #mergeInterchange(legA, legB) {
    legA.sort((a, b) => hhmmToMin(a.arrival) - hhmmToMin(b.arrival));
    legB.sort((a, b) => hhmmToMin(a.departure) - hhmmToMin(b.departure));

    const out = [];
    for (const a of legA) {
      const arrivalMin = hhmmToMin(a.arrival);
      const candidate = legB.find(b => hhmmToMin(b.departure) >= arrivalMin + 1);
      if (!candidate) continue;
      out.push({
        leg1: a,
        leg2: candidate,
        totalMin: hhmmToMin(candidate.arrival) - hhmmToMin(a.departure),
        waitMin:  hhmmToMin(candidate.departure) - arrivalMin,
      });
    }
    // Sort final routes by when the user actually leaves.
    out.sort((x, y) => hhmmToMin(x.leg1.departure) - hhmmToMin(y.leg1.departure));
    return out;
  }
}

export function hhmmToMin(s) {
  if (!s) return 0;
  const [h, m] = s.split(":").map(Number);
  return h * 60 + m;
}

export function formatDuration(min) {
  if (min < 0) min += 24 * 60;
  const h = Math.floor(min / 60), m = min % 60;
  if (h === 0) return `${m}m`;
  return m === 0 ? `${h}h` : `${h}h ${m}m`;
}
