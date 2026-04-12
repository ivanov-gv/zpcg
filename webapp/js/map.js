// map.js — Leaflet initialization and rendering. Kept in a separate module
// so the search-heavy screens don't pull Leaflet internals into their
// closures. The map view is lazy-mounted by the router.

import { t } from "./i18n.js";

export class StationsMap {
  constructor(db) {
    this.db = db;
    this.map = null;
    this.markers = new Map();     // stationId → Leaflet marker
    this.routeLayer = null;
    this.onStationClick = null;
  }

  mount(containerEl, { onStationClick, highlight = [], routeIds = [] } = {}) {
    this.onStationClick = onStationClick;

    if (typeof L === 'undefined') {
      containerEl.innerHTML = `<div style="display:flex;align-items:center;justify-content:center;height:100%;color:var(--hint-color,#707579);padding:2rem;text-align:center">${escapeHtml(t("map.offline"))}</div>`;
      return;
    }

    try {
      this.map = L.map(containerEl, {
        center: [42.6, 19.3],
        zoom: 8,
        zoomControl: true,
        attributionControl: true,
      });
      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: '&copy; OpenStreetMap',
        maxZoom: 18,
      }).addTo(this.map);

      for (const s of this.db.stations) {
        const coord = this.db.coordOf(s.id);
        if (!coord) continue;
        const isMain = s.type === 4;
        const marker = L.circleMarker([coord.lat, coord.lng], {
          radius: isMain ? 7 : 4,
          color: isMain ? "#c0392b" : "#707579",
          weight: isMain ? 2 : 1,
          fillColor: isMain ? "#c0392b" : "#ffffff",
          fillOpacity: isMain ? 0.9 : 0.9,
        });
        const typeName = t(`station.type.${s.type}`);
        marker.bindPopup(
          `<b>${escapeHtml(s.name)}</b><br>` +
          `<span style="color:#707579">${escapeHtml(s.nameCyr)}</span><br>` +
          `<small>${escapeHtml(typeName)}</small>`
        );
        marker.on("click", () => { this.onStationClick?.(s); });
        marker.addTo(this.map);
        this.markers.set(s.id, marker);
      }

      if (routeIds.length >= 2) this.showRoute(routeIds);
      if (highlight.length) this.highlightStations(highlight);
    } catch (e) {
      console.warn("Map init failed:", e);
      containerEl.innerHTML = `<div style="display:flex;align-items:center;justify-content:center;height:100%;color:var(--hint-color,#707579);padding:2rem;text-align:center">${escapeHtml(t("map.error"))}</div>`;
    }
  }

  showRoute(stationIds) {
    if (!this.map) return;
    this.clearRoute();
    const latlngs = stationIds
      .map(id => this.db.coordOf(id))
      .filter(Boolean)
      .map(c => [c.lat, c.lng]);
    if (latlngs.length < 2) return;
    this.routeLayer = L.polyline(latlngs, {
      color: "#c0392b",
      weight: 4,
      opacity: 0.75,
      dashArray: "1 8",
      lineCap: "round",
    }).addTo(this.map);
    this.map.fitBounds(this.routeLayer.getBounds(), { padding: [40, 40] });
  }

  clearRoute() {
    if (this.routeLayer) {
      this.routeLayer.remove();
      this.routeLayer = null;
    }
  }

  highlightStations(ids) {
    if (!this.map) return;
    for (const id of ids) {
      const marker = this.markers.get(id);
      if (!marker) continue;
      marker.setStyle({
        radius: 9,
        fillColor: "#d4a017",
        color: "#c0392b",
        weight: 3,
      });
      marker.bringToFront();
    }
  }

  flyToStation(id) {
    if (!this.map) return;
    const coord = this.db.coordOf(id);
    if (!coord) return;
    this.map.flyTo([coord.lat, coord.lng], 11, { duration: 0.6 });
    this.markers.get(id)?.openPopup();
  }

  destroy() {
    if (this.map) {
      this.map.remove();
      this.map = null;
    }
    this.markers.clear();
  }
}

function escapeHtml(s) {
  return String(s)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}
