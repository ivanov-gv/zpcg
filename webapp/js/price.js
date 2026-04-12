// price.js — rough ticket-price estimator.
//
// The real Montenegro timetable does NOT publish a machine-readable price
// table. ZPCG sells tickets only at ticket offices or on board. The numbers
// here are an illustrative model based on publicly quoted fares from
// traveller reports and the 2024 Bar↔Belgrade price (~€21). Use the
// estimator ONLY as a "roughly how much" hint — the UI is clear about it.

const BASE_EUR_PER_KM = 0.06;        // ~6 cents/km (observed on local lines)
const MIN_FARE_EUR    = 1.20;        // commuter minimum
const CLASS1_MULT     = 1.5;
const DISCOUNT = {
  adult:   1.00,
  child:   0.50,   // children 6–14
  student: 0.80,   // student card
  senior:  0.70,   // 65+
};

export function estimatePrice(km, { ticketClass = 2, passenger = 'adult', count = 1 } = {}) {
  if (km == null || km <= 0) return null;

  const base = Math.max(MIN_FARE_EUR, km * BASE_EUR_PER_KM);
  const classMult = ticketClass === 1 ? CLASS1_MULT : 1;
  const discount  = DISCOUNT[passenger] ?? 1;

  const perPerson = base * classMult * discount;
  const total = perPerson * count;
  return {
    perPersonEur: round2(perPerson),
    totalEur:     round2(total),
    baseEur:      round2(base),
    km:           Math.round(km),
    classMult,
    discount,
  };
}

function round2(n) { return Math.round(n * 100) / 100; }

export function formatEur(n) {
  if (n == null) return "—";
  return `€${n.toFixed(2)}`;
}
