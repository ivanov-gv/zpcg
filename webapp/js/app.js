// app.js — entry point. Wires up data loading, theme syncing with Telegram,
// the router, and every view the app renders. Views are plain functions
// that return a DocumentFragment to keep the prototype build-tool-free.

import { loadDb } from "./data.js";
import { PathFinder, hhmmToMin, formatDuration } from "./pathfinder.js";
import { Router } from "./router.js";
import { StationsMap } from "./map.js";
import { estimatePrice, formatEur } from "./price.js";
import { t, initLang, setLang, getLang, LANGS, applyStaticTranslations } from "./i18n.js";

// ---- Telegram integration (soft — safe if not inside Telegram) ------------

const tg = window.Telegram?.WebApp;
if (tg) {
  tg.ready();
  tg.expand();
  // Mirror Telegram theme variables onto our own CSS variables so the
  // existing stylesheet stays unaware of Telegram specifics.
  applyTelegramTheme();
  tg.onEvent?.("themeChanged", applyTelegramTheme);
}
function applyTelegramTheme() {
  if (!tg?.themeParams) return;
  const p = tg.themeParams;
  const map = {
    "--tg-bg":           p.bg_color,
    "--tg-text":         p.text_color,
    "--tg-muted":        p.subtitle_text_color || p.hint_color,
    "--tg-hint":         p.hint_color,
    "--tg-link":         p.link_color,
    "--tg-button":       p.button_color,
    "--tg-button-text":  p.button_text_color,
    "--tg-secondary-bg": p.secondary_bg_color,
  };
  const root = document.documentElement;
  for (const [k, v] of Object.entries(map)) if (v) root.style.setProperty(k, v);
  document.body.dataset.tgTheme = tg.colorScheme || "light";
}

// ---- app state ------------------------------------------------------------

// Very small shared state; views read via closures, writes go through set()
// so the current view re-renders. Keeps us out of the framework zoo.
const state = {
  db: null,
  pf: null,
  fromId: null,
  toId: null,
  date: todayLabel(),
  saved: loadSaved(),
  alerts: loadAlerts(),
};

function set(patch) {
  Object.assign(state, patch);
  // Re-render the active route so pickers + swap feed back into the view.
  router.handle();
}

// ---- bootstrap ------------------------------------------------------------

const router = new Router(document.getElementById("view-root"));

init().catch(err => {
  console.error(err);
  document.getElementById("view-root").textContent =
    "Failed to load prototype data: " + err.message;
});

async function init() {
  // Default language: Telegram user → browser → English.
  const preferred =
    tg?.initDataUnsafe?.user?.language_code ||
    navigator.language ||
    "en";
  initLang(preferred);

  state.db = await loadDb();
  state.pf = new PathFinder(state.db);

  // Sensible default = Bar → Podgorica, the #1 real-world query.
  const bar = state.db.stationByLowerName.get("bar");
  const podgorica = state.db.stationByLowerName.get("podgorica");
  if (bar && podgorica) {
    state.fromId = bar.id;
    state.toId   = podgorica.id;
  }

  // Routes ---------------------------------------------------------------
  router.add("/",                    viewSearch);
  router.add("/results",             viewResults);
  router.add("/train/:id",           viewTrainDetails);
  router.add("/train/:id/:fromId/:toId", viewTrainDetails);
  router.add("/station/:id",         viewStation);
  router.add("/map",                 viewMap);
  router.add("/alerts",              viewAlerts);
  router.add("/saved",               viewSaved);
  router.add("/info",                viewInfo);
  router.add("/price/:fromId/:toId/:trainId", viewPrice);

  router.onBeforeMount = updateTabBar;

  document.getElementById("back-btn").addEventListener("click", () => history.back());
  document.getElementById("lang-btn").addEventListener("click", openLangPicker);
  document.querySelectorAll(".tab").forEach(tab => {
    tab.addEventListener("click", () => router.go(tab.dataset.route));
  });

  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('./sw.js').catch(err =>
      console.warn('SW registration failed:', err)
    );
  }

  applyStaticTranslations();
  router.handle();
}

function updateTabBar(hash) {
  document.querySelectorAll(".tab").forEach(tab => {
    const matches =
      tab.dataset.route === hash ||
      (tab.dataset.route === "#/" && hash.startsWith("#/results")) ||
      (tab.dataset.route === "#/" && hash.startsWith("#/train")) ||
      (tab.dataset.route === "#/" && hash.startsWith("#/station")) ||
      (tab.dataset.route === "#/" && hash.startsWith("#/price"));
    tab.classList.toggle("active", matches);
  });

  const backBtn = document.getElementById("back-btn");
  const homeish = hash === "#/" || hash === "#/map" || hash === "#/alerts" ||
                  hash === "#/saved" || hash === "#/info";
  backBtn.hidden = homeish;

  const titles = {
    "#/":         t("tab.search"),
    "#/results":  t("results.direct") + " / " + t("results.interchange").replace(/Change.*/, ""),
    "#/map":      t("map.title"),
    "#/alerts":   t("alerts.title"),
    "#/saved":    t("saved.title"),
    "#/info":     t("info.title"),
  };
  const matchKey = Object.keys(titles).find(k => hash === k);
  document.getElementById("page-title").textContent = titles[matchKey] || "Monterails";
}

// ---- views ---------------------------------------------------------------

// Search -------------------------------------------------------------------
function viewSearch() {
  const frag = document.createDocumentFragment();

  const title = el("h1", "section-title");
  title.style.marginTop = "16px";
  title.textContent = t("search.title");
  frag.appendChild(title);

  const card = el("div", "card search-card");
  frag.appendChild(card);

  const fromField = fieldRow(t("search.from"), "🟢",
    state.fromId ? state.db.station(state.fromId).name : t("search.pick"),
    () => openStationPicker(id => set({ fromId: id }))
  );
  const toField   = fieldRow(t("search.to"),   "🔴",
    state.toId ? state.db.station(state.toId).name : t("search.pick"),
    () => openStationPicker(id => set({ toId: id }))
  );

  const fieldWrap = el("div", "field-row");
  const swapBtn = el("button", "swap-btn");
  swapBtn.textContent = "⇅";
  swapBtn.addEventListener("click", () => set({
    fromId: state.toId,
    toId:   state.fromId,
  }));

  const leftCol  = el("div");
  leftCol.style.flex = "1";
  leftCol.appendChild(fromField);
  leftCol.appendChild(toField);

  fieldWrap.appendChild(leftCol);
  fieldWrap.appendChild(swapBtn);
  card.appendChild(fieldWrap);

  // Date pills (today/tomorrow/day-of-week) -------------------------------
  const pills = el("div", "pills");
  const dates = ["today", "tomorrow", "weekday"];
  dates.forEach((k, i) => {
    const b = el("button");
    b.textContent = k === "today" ? t("search.today")
                  : k === "tomorrow" ? t("search.tomorrow")
                  : weekdayLabel();
    if (i === 0) b.style.background = "rgba(192,57,43,0.15)";
    pills.appendChild(b);
  });
  card.appendChild(pills);

  // CTA -----------------------------------------------------------------
  const ctaWrap = el("div", "cta");
  const cta = el("button", "btn-primary");
  cta.textContent = t("search.cta");
  cta.disabled = !(state.fromId && state.toId && state.fromId !== state.toId);
  cta.addEventListener("click", () => router.go("#/results"));
  ctaWrap.appendChild(cta);
  frag.appendChild(ctaWrap);

  // Popular routes ------------------------------------------------------
  const popular = [
    ["Bar", "Podgorica"],
    ["Podgorica", "Nikšić"],
    ["Bar", "Bijelo Polje"],
    ["Podgorica", "Kolašin"],
    ["Bar", "Beograd Centar"],
  ];
  const pTitle = el("div", "section-title");
  pTitle.textContent = t("search.popular");
  frag.appendChild(pTitle);
  popular.forEach(([a, b]) => {
    const aS = state.db.stationByLowerName.get(a.toLowerCase());
    const bS = state.db.stationByLowerName.get(b.toLowerCase());
    if (!aS || !bS) return;
    const row = el("a", "info-item");
    row.href = "javascript:void(0)";
    row.innerHTML = `<span class="info-icon">🚆</span>
      <div><div>${aS.name} → ${bS.name}</div><small>${t("search.popular")}</small></div>
      <span class="chev">›</span>`;
    row.addEventListener("click", () => {
      set({ fromId: aS.id, toId: bS.id });
      router.go("#/results");
    });
    frag.appendChild(row);
  });

  return frag;
}

// Results view -------------------------------------------------------------
function viewResults() {
  const frag = document.createDocumentFragment();
  if (!state.fromId || !state.toId) {
    router.go("#/");
    return frag;
  }
  const from = state.db.station(state.fromId);
  const to   = state.db.station(state.toId);

  const header = el("div", "route-header");
  header.innerHTML = `<h2>${from.name} → ${to.name}</h2>
    <span class="muted">${state.date}</span>`;
  frag.appendChild(header);

  const result = state.pf.find(state.fromId, state.toId);

  if (result.kind === "none") {
    const empty = el("div", "empty-state");
    empty.innerHTML = `<div class="big">🚫</div><div>${t("results.none")}</div>`;
    frag.appendChild(empty);

    // Show tips / reverse button
    const actions = el("div", "cta");
    const rev = el("button", "btn-ghost");
    rev.textContent = t("results.reverse");
    rev.addEventListener("click", () => set({
      fromId: state.toId, toId: state.fromId,
    }));
    actions.appendChild(rev);
    frag.appendChild(actions);
    return frag;
  }

  if (result.kind === "direct") {
    const chip = el("div", "chip ok");
    chip.textContent = t("results.direct");
    chip.style.margin = "4px 16px 8px";
    frag.appendChild(chip);
    result.direct.forEach(route => frag.appendChild(renderDirectCard(route, from, to)));
  } else {
    const chip = el("div", "chip warn");
    chip.textContent = t("results.interchange");
    chip.style.margin = "4px 16px 8px";
    frag.appendChild(chip);
    result.interchange.forEach(route => {
      frag.appendChild(renderInterchangeCard(route, from, to));
    });
  }

  // Save / reverse footer
  const footer = el("div", "cta");
  footer.style.display = "flex"; footer.style.gap = "8px";
  const saveBtn = el("button", "btn-ghost");
  saveBtn.textContent = "⭐ " + t("results.save");
  saveBtn.style.flex = "1";
  saveBtn.addEventListener("click", () => {
    addSaved({ fromId: state.fromId, toId: state.toId });
    tg?.showPopup?.({ message: t("results.save") + " ✓" });
    saveBtn.textContent = "✓";
  });
  const revBtn = el("button", "btn-ghost");
  revBtn.textContent = "⇅ " + t("results.reverse");
  revBtn.style.flex = "1";
  revBtn.addEventListener("click", () => set({
    fromId: state.toId, toId: state.fromId,
  }));
  footer.appendChild(saveBtn);
  footer.appendChild(revBtn);
  frag.appendChild(footer);

  return frag;
}

function renderDirectCard(route, from, to) {
  const card = el("div", "train-card direct");
  card.appendChild(renderCardBody(route, from, to, "direct"));
  card.addEventListener("click", e => {
    if (e.target.closest("button")) return;
    router.go(`#/train/${route.trainId}/${from.id}/${to.id}`);
  });
  return card;
}

function renderInterchangeCard(route, from, to) {
  const card = el("div", "train-card interchange");
  const top = el("div", "train-top");
  const ids = el("span", "train-id");
  ids.textContent = `${route.leg1.trainId} → ${route.leg2.trainId}`;
  const times = el("div", "train-times");
  times.innerHTML = `${route.leg1.departure} <span class="dash">→</span> ${route.leg2.arrival}
    <span class="dur">${formatDuration(route.totalMin)}</span>`;
  top.appendChild(ids);
  top.appendChild(times);
  card.appendChild(top);

  const transfer = state.db.station(state.db.transferStationId);
  const inter = el("div", "interchange-block");
  inter.innerHTML = `🔄 <b>${transfer.name}</b> &middot;
    ${route.leg1.arrival} → ${route.leg2.departure}
    (${t("results.wait")} ${formatDuration(route.waitMin)})`;
  card.appendChild(inter);

  const meta = el("div", "train-meta");
  const km = state.db.railKm(route.leg1.trainId, from.id, state.db.transferStationId) || 0;
  const km2 = state.db.railKm(route.leg2.trainId, state.db.transferStationId, to.id) || 0;
  const price = estimatePrice(km + km2);
  meta.innerHTML = `<span>${from.name} → ${to.name}</span>
    <span class="price">${t("results.price_est")} ${formatEur(price?.perPersonEur)}</span>`;
  card.appendChild(meta);
  card.addEventListener("click", e => {
    if (e.target.closest("button")) return;
    router.go(`#/train/${route.leg1.trainId}/${from.id}/${state.db.transferStationId}`);
  });
  return card;
}

function renderCardBody(route, from, to, kind) {
  const wrap = el("div");
  const top = el("div", "train-top");
  const id  = el("span", "train-id"); id.textContent = `№ ${route.trainId}`;
  const times = el("div", "train-times");
  times.innerHTML = `${route.departure} <span class="dash">→</span> ${route.arrival}
    <span class="dur">${formatDuration(route.durationMin)}</span>`;
  top.appendChild(id);
  top.appendChild(times);
  wrap.appendChild(top);

  const km = state.db.railKm(route.trainId, from.id, to.id);
  const price = estimatePrice(km);

  const meta = el("div", "train-meta");
  meta.innerHTML = `<span>${from.name} → ${to.name}</span>
    <span class="price">${t("results.price_est")} ${formatEur(price?.perPersonEur)}</span>`;
  wrap.appendChild(meta);
  return wrap;
}

// Train details ------------------------------------------------------------
function viewTrainDetails({ params }) {
  const frag = document.createDocumentFragment();
  const train = state.db.train(Number(params.id));
  if (!train) {
    frag.textContent = "Unknown train";
    return frag;
  }

  const fromId = params.fromId ? Number(params.fromId) : null;
  const toId   = params.toId   ? Number(params.toId)   : null;

  const head = el("div", "route-header");
  head.innerHTML = `<h2>№ ${train.id}</h2>
    <span class="muted">${train.stops.length} ${t("results.details")}</span>`;
  frag.appendChild(head);

  const list = el("ol", "stop-list");
  train.stops.forEach(stop => {
    const station = state.db.station(stop.stationId);
    if (!station) return;
    const li = el("li");
    const isHighlight = stop.stationId === fromId || stop.stationId === toId;
    if (isHighlight) li.classList.add("highlight");
    li.innerHTML = `<span>${station.name}</span>
      <span class="t">${stop.arrival}${stop.arrival !== stop.departure ? " → " + stop.departure : ""}</span>`;
    li.addEventListener("click", () => router.go(`#/station/${station.id}`));
    list.appendChild(li);
  });
  frag.appendChild(list);

  if (fromId && toId) {
    const priceCta = el("div", "cta");
    const b = el("button", "btn-ghost");
    b.textContent = "💶 " + t("results.price_est");
    b.addEventListener("click", () =>
      router.go(`#/price/${fromId}/${toId}/${train.id}`));
    priceCta.appendChild(b);
    frag.appendChild(priceCta);
  }

  return frag;
}

// Station details ----------------------------------------------------------
function viewStation({ params }) {
  const frag = document.createDocumentFragment();
  const station = state.db.station(Number(params.id));
  if (!station) return frag;

  const head = el("div", "route-header");
  head.innerHTML = `<h2>${station.name}</h2>
    <span class="muted">${station.nameCyr}</span>`;
  frag.appendChild(head);

  const card = el("div", "card");
  const type = el("div"); type.className = "chip brand";
  type.textContent = t(`station.type.${station.type}`);
  card.appendChild(type);

  const coord = state.db.coordOf(station.id);
  if (coord) {
    const c = el("p"); c.style.margin = "10px 0 0";
    c.style.fontSize = "12px"; c.style.color = "var(--tg-muted)";
    c.textContent = `${coord.lat.toFixed(4)}, ${coord.lng.toFixed(4)}`;
    card.appendChild(c);
  }

  // Trains serving this station
  const trainIds = [...(state.db.stationToTrainIds.get(station.id) || [])].sort((a, b) => a - b);
  const tTitle = el("div", "section-title");
  tTitle.textContent = `${trainIds.length} trains`;
  frag.appendChild(card);
  frag.appendChild(tTitle);

  trainIds.forEach(id => {
    const t = state.db.train(id);
    const stop = t.stops.find(s => s.stationId === station.id);
    const row = el("a", "info-item");
    row.href = "javascript:void(0)";
    row.innerHTML = `<span class="info-icon">🚆</span>
      <div><div>№ ${t.id}</div>
      <small>${stop.arrival}${stop.arrival !== stop.departure ? " → " + stop.departure : ""}</small></div>
      <span class="chev">›</span>`;
    row.addEventListener("click", () => router.go(`#/train/${t.id}`));
    frag.appendChild(row);
  });

  const actions = el("div", "cta");
  actions.style.display = "flex"; actions.style.gap = "8px";
  const useFrom = el("button", "btn-ghost");
  useFrom.textContent = "🟢 " + t("search.from");
  useFrom.style.flex = "1";
  useFrom.addEventListener("click", () => {
    state.fromId = station.id;
    router.go("#/");
  });
  const useTo = el("button", "btn-ghost");
  useTo.textContent = "🔴 " + t("search.to");
  useTo.style.flex = "1";
  useTo.addEventListener("click", () => {
    state.toId = station.id;
    router.go("#/");
  });
  actions.appendChild(useFrom);
  actions.appendChild(useTo);
  frag.appendChild(actions);

  return frag;
}

// Price breakdown ----------------------------------------------------------
function viewPrice({ params }) {
  const frag = document.createDocumentFragment();
  const fromId = Number(params.fromId);
  const toId   = Number(params.toId);
  const trainId = Number(params.trainId);
  const from = state.db.station(fromId);
  const to   = state.db.station(toId);

  const head = el("div", "route-header");
  head.innerHTML = `<h2>${from.name} → ${to.name}</h2>
    <span class="muted">№ ${trainId}</span>`;
  frag.appendChild(head);

  const km = state.db.railKm(trainId, fromId, toId) ||
             state.db.haversineKm(fromId, toId);

  // Class + passenger state is local to this screen via closures. The whole
  // card is rebuilt on every change so the chip state stays in sync with
  // the breakdown — keeps the view truly stateless.
  let pClass = 2;
  let pType = "adult";
  let pCount = 1;

  const card = el("div", "card");
  frag.appendChild(card);

  const note = el("p");
  note.style.cssText = "font-size:11px;color:var(--tg-muted);padding:0 16px";
  note.textContent = t("price.note");
  frag.appendChild(note);

  function redraw() {
    card.replaceChildren();
    card.appendChild(selectorRow(t("price.class"),
      [[2, t("price.class2")], [1, t("price.class1")]],
      () => pClass, v => { pClass = v; redraw(); }));
    card.appendChild(selectorRow("👤",
      [["adult", t("price.adult")], ["child", t("price.child")],
       ["student", t("price.student")], ["senior", t("price.senior")]],
      () => pType, v => { pType = v; redraw(); }));
    card.appendChild(selectorRow("×",
      [[1, "1"], [2, "2"], [3, "3"], [4, "4"]],
      () => pCount, v => { pCount = v; redraw(); }));

    const breakdown = el("div");
    breakdown.style.marginTop = "12px";
    card.appendChild(breakdown);

    const est = estimatePrice(km, {
      ticketClass: pClass, passenger: pType, count: pCount,
    });
    if (!est) { breakdown.textContent = "—"; return; }
    breakdown.appendChild(priceRow(`≈ ${est.km} km × €${(0.06).toFixed(2)}`, formatEur(est.baseEur)));
    if (est.classMult !== 1)
      breakdown.appendChild(priceRow(t("price.class1") + " ×" + est.classMult, ""));
    if (est.discount !== 1)
      breakdown.appendChild(priceRow(t("price." + pType) + " ×" + est.discount, ""));
    if (pCount > 1)
      breakdown.appendChild(priceRow("× " + pCount, ""));
    breakdown.appendChild(priceRow(t("price.total"), formatEur(est.totalEur)));
  }
  redraw();
  return frag;
}

// Map view -----------------------------------------------------------------
let currentMap = null;
function viewMap() {
  // Leaflet needs a real, sized container. We return a wrapper that
  // absolutely-positions the Leaflet div inside the view-root.
  const wrap = el("div"); wrap.id = "map-wrap";
  wrap.style.cssText = "position:relative;height:calc(100vh - var(--header-h) - var(--tab-h));";

  const map = el("div"); map.id = "leaflet-map";
  wrap.appendChild(map);

  const legend = el("div", "map-legend");
  legend.innerHTML = `
    <div><span class="dot" style="background:#c0392b"></span>${t("map.legend.main")}</div>
    <div><span class="dot" style="background:#fff;border:1px solid #707579"></span>${t("map.legend.stop")}</div>
    <div><span class="dot" style="background:#d4a017"></span>${t("map.legend.sel")}</div>
    <div><span class="dot" style="background:#c0392b;height:2px;border-radius:0"></span>${t("map.legend.line")}</div>
    <div style="margin-top:6px;color:var(--tg-muted);font-size:10px">${t("map.tap")}</div>`;
  wrap.appendChild(legend);

  if (currentMap) { currentMap.destroy(); currentMap = null; }
  // Defer Leaflet init until after the element is in the DOM.
  setTimeout(() => {
    currentMap = new StationsMap(state.db);
    const highlight = [state.fromId, state.toId].filter(Boolean);
    currentMap.mount(map, {
      highlight,
      routeIds: computeMapRouteIds(),
      onStationClick: s => {
        if (!state.fromId) { set({ fromId: s.id }); return; }
        if (!state.toId || state.fromId === s.id) { set({ toId: s.id }); return; }
        // If both picked, clicking replaces "to".
        set({ toId: s.id });
      },
    });
  }, 0);

  return wrap;
}

function computeMapRouteIds() {
  if (!state.fromId || !state.toId) return [];
  const r = state.pf.find(state.fromId, state.toId);

  // Extract station IDs along a segment, handling overnight trains where
  // fromIdx > toIdx (the stops array wraps at midnight).
  const seg = (train, aId, bId) => {
    const ai = train.stops.findIndex(s => s.stationId === aId);
    const bi = train.stops.findIndex(s => s.stationId === bId);
    if (ai < 0 || bi < 0) return [];
    if (ai <= bi) {
      return train.stops.slice(ai, bi + 1).map(s => s.stationId);
    }
    // Overnight: wrap around the end of the array.
    return [
      ...train.stops.slice(ai).map(s => s.stationId),
      ...train.stops.slice(0, bi + 1).map(s => s.stationId),
    ];
  };

  if (r.kind === "direct" && r.direct[0]) {
    return seg(state.db.train(r.direct[0].trainId), state.fromId, state.toId);
  }
  if (r.kind === "interchange" && r.interchange[0]) {
    const transfer = state.db.transferStationId;
    const first  = state.db.train(r.interchange[0].leg1.trainId);
    const second = state.db.train(r.interchange[0].leg2.trainId);
    return [...seg(first, state.fromId, transfer), ...seg(second, transfer, state.toId).slice(1)];
  }
  return [];
}

// Alerts view --------------------------------------------------------------
function viewAlerts() {
  const frag = document.createDocumentFragment();

  // Mocked service-disruption cards (what ZPCG announces on Facebook/Viber).
  const disruption = el("div", "card alert-card");
  disruption.innerHTML = `
    <div class="chip warn" style="margin-bottom:8px">⚠️ ${t("alerts.disruption")}</div>
    <div class="alert-title">Podgorica — Nikšić</div>
    <div class="alert-time">2026-04-12 · 09:00–17:00</div>
    <div class="alert-route" style="margin-top:4px">Track work Spuž–Danilovgrad.<br>Trains 7102–7106 replaced by bus.</div>`;
  frag.appendChild(disruption);

  const tTitle = el("div", "section-title");
  tTitle.textContent = "Your subscriptions";
  frag.appendChild(tTitle);

  if (state.alerts.length === 0) {
    const empty = el("div", "empty-state");
    empty.innerHTML = `<div class="big">🔕</div><div>${t("alerts.empty")}</div>`;
    frag.appendChild(empty);
  } else {
    state.alerts.forEach((a, i) => {
      const row = el("div", "card alert-card");
      row.innerHTML = `
        <div class="alert-title">${a.label}</div>
        <div class="alert-time">${a.days} · ${a.time}</div>
        <div class="alert-actions">
          <button class="btn-ghost">${t("alerts.unsubscribe")}</button>
        </div>`;
      row.querySelector("button").addEventListener("click", () => {
        state.alerts.splice(i, 1); saveAlerts(); set({});
      });
      frag.appendChild(row);
    });
  }

  const cta = el("div", "cta");
  const add = el("button", "btn-primary");
  add.textContent = "+ " + t("alerts.new");
  add.addEventListener("click", () => {
    if (!state.fromId || !state.toId) {
      tg?.showPopup?.({ message: "Pick a route first" }) ||
        alert("Pick a route first");
      return;
    }
    const from = state.db.station(state.fromId);
    const to   = state.db.station(state.toId);
    state.alerts.push({
      label: `${from.name} → ${to.name}`,
      days:  "Mon–Fri",
      time:  "08:00 ± 30min",
    });
    saveAlerts();
    set({});
  });
  cta.appendChild(add);
  frag.appendChild(cta);

  return frag;
}

// Saved view ---------------------------------------------------------------
function viewSaved() {
  const frag = document.createDocumentFragment();

  if (state.saved.length === 0) {
    const empty = el("div", "empty-state");
    empty.innerHTML = `<div class="big">⭐</div><div>${t("saved.empty")}</div>`;
    frag.appendChild(empty);
    return frag;
  }
  state.saved.forEach((s, i) => {
    const from = state.db.station(s.fromId);
    const to   = state.db.station(s.toId);
    if (!from || !to) return;
    const row = el("a", "info-item");
    row.href = "javascript:void(0)";
    row.innerHTML = `<span class="info-icon">⭐</span>
      <div><div>${from.name} → ${to.name}</div></div>
      <span class="chev">›</span>`;
    row.addEventListener("click", () => {
      set({ fromId: s.fromId, toId: s.toId });
      router.go("#/results");
    });
    frag.appendChild(row);
  });

  const clr = el("div", "cta");
  const b = el("button", "btn-ghost");
  b.textContent = t("saved.clear");
  b.addEventListener("click", () => { state.saved = []; saveSaved(); set({}); });
  clr.appendChild(b);
  frag.appendChild(clr);
  return frag;
}

// Info view ----------------------------------------------------------------
function viewInfo() {
  const frag = document.createDocumentFragment();

  const head = el("div", "section-title");
  head.textContent = t("info.title");
  frag.appendChild(head);

  const rows = [
    ["🤖", t("info.bot"),    null, "https://t.me/Monterails_bot"],
    ["🌐", t("info.zpcg"),   null, "https://zpcg.me/search"],
    ["🏛️", t("info.site"),   null, "https://www.zcg-prevoz.me/"],
    ["📦", t("info.source"), null, "https://github.com/ivanov-gv/zpcg"],
  ];
  rows.forEach(([icon, label, sub, href]) => {
    const a = el("a", "info-item");
    a.href = href; a.target = "_blank"; a.rel = "noopener";
    a.innerHTML = `<span class="info-icon">${icon}</span>
      <div><div>${label}</div>${sub ? `<small>${sub}</small>` : ""}</div>
      <span class="chev">↗</span>`;
    a.addEventListener("click", e => {
      // If inside Telegram, route external links through the host WebApp API.
      if (tg?.openLink) { e.preventDefault(); tg.openLink(href); }
    });
    frag.appendChild(a);
  });

  const meta = el("div", "card");
  meta.style.marginTop = "12px";
  const updated = new Date(state.db.meta.generatedAt);
  meta.innerHTML = `
    <div class="price-row muted">${t("info.updated")}</div>
    <div class="price-row"><span>UTC</span><span>${updated.toISOString().slice(0, 16).replace("T", " ")}</span></div>
    <div class="price-row muted">${t("info.theme")}</div>
    <div class="price-row muted">${t("info.legal")}</div>`;
  frag.appendChild(meta);

  // Language picker (inline)
  const lpTitle = el("div", "section-title");
  lpTitle.textContent = t("info.lang");
  frag.appendChild(lpTitle);
  LANGS.forEach(l => {
    const row = el("a", "info-item");
    row.href = "javascript:void(0)";
    const active = l.code === getLang();
    row.innerHTML = `<span class="info-icon">${active ? "✓" : "🌐"}</span>
      <div><div>${l.name}</div><small>${l.label}</small></div>`;
    row.addEventListener("click", () => { setLang(l.code); set({}); });
    frag.appendChild(row);
  });

  return frag;
}

// ---- picker sheet --------------------------------------------------------

function openStationPicker(onPick) {
  const tpl = document.getElementById("tpl-picker").content.cloneNode(true);
  const backdrop = tpl.querySelector(".sheet-backdrop");
  const sheet    = tpl.querySelector(".sheet");
  const input    = tpl.querySelector(".sheet-input");
  const list     = tpl.querySelector(".station-list");
  const close    = tpl.querySelector(".sheet-close");

  document.body.appendChild(tpl);
  const wrapper = [backdrop, sheet];
  input.focus();

  function refresh() {
    const matches = state.db.search(input.value, { includeAll: false });
    list.innerHTML = "";
    matches.slice(0, 80).forEach(({ station }) => {
      const li = el("li");
      const bl = state.db.isBlacklisted(station.name);
      if (bl) li.classList.add("blacklisted");
      li.innerHTML = `
        <div>
          <div class="st-name">${station.name}</div>
          <small class="st-cyr">${station.nameCyr}</small>
        </div>
        <span class="st-type">${t(`station.type.${station.type}`)}</span>`;
      li.addEventListener("click", () => {
        if (bl) { alert(t("picker.blacklisted")); return; }
        onPick(station.id);
        dismiss();
      });
      list.appendChild(li);
    });
  }
  refresh();
  input.addEventListener("input", refresh);
  close.addEventListener("click", dismiss);
  backdrop.addEventListener("click", dismiss);

  applyStaticTranslations(sheet);

  function dismiss() {
    wrapper.forEach(el => el.remove());
  }
}

// ---- language picker (tiny — pops from the header) ----------------------

function openLangPicker() {
  const labels = LANGS.map(l => l.label).join(" / ");
  const chosen = prompt("Language? " + labels, getLang().toUpperCase());
  if (!chosen) return;
  const code = chosen.toLowerCase().slice(0, 2);
  if (LANGS.some(l => l.code === code)) {
    setLang(code);
    document.getElementById("lang-btn").textContent = code.toUpperCase();
    set({});
  }
}

// ---- helpers -------------------------------------------------------------

function el(tag, cls) {
  const e = document.createElement(tag);
  if (cls) e.className = cls;
  return e;
}

function fieldRow(label, icon, value, onClick) {
  const f = el("div", "field");
  if (value === t("search.pick")) f.classList.add("empty");
  f.innerHTML = `<span class="field-icon">${icon}</span>
    <div><span class="field-label">${label}</span>
    <span class="field-value">${value}</span></div>`;
  f.addEventListener("click", onClick);
  return f;
}

function priceRow(label, value) {
  const r = el("div", "price-row");
  r.innerHTML = `<span>${label}</span><span>${value}</span>`;
  return r;
}

function selectorRow(label, items, get, setter) {
  const row = el("div");
  row.style.cssText = "display:flex;align-items:center;gap:6px;margin-top:8px;flex-wrap:wrap";
  const l = el("span"); l.style.cssText = "font-size:12px;color:var(--tg-muted);width:50px";
  l.textContent = label;
  row.appendChild(l);
  items.forEach(([val, text]) => {
    const b = el("button", "chip");
    b.textContent = text;
    if (get() === val) b.classList.add("brand");
    b.addEventListener("click", () => setter(val));
    row.appendChild(b);
  });
  return row;
}

function todayLabel() {
  const d = new Date();
  return d.toISOString().slice(0, 10);
}

function weekdayLabel() {
  return new Date().toLocaleDateString(getLang(), { weekday: "short" });
}

function loadSaved() {
  try { return JSON.parse(localStorage.getItem("monterails.saved") || "[]"); }
  catch { return []; }
}
function saveSaved() {
  try { localStorage.setItem("monterails.saved", JSON.stringify(state.saved)); } catch {}
}
function addSaved(entry) {
  const exists = state.saved.some(s => s.fromId === entry.fromId && s.toId === entry.toId);
  if (!exists) state.saved.push(entry);
  saveSaved();
}

function loadAlerts() {
  try { return JSON.parse(localStorage.getItem("monterails.alerts") || "[]"); }
  catch { return []; }
}
function saveAlerts() {
  try { localStorage.setItem("monterails.alerts", JSON.stringify(state.alerts)); } catch {}
}
