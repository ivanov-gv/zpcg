const CACHE_VERSION = 'monterails-v1';

const LOCAL_ASSETS = [
  './',
  './index.html',
  './css/app.css',
  './js/app.js',
  './js/router.js',
  './js/data.js',
  './js/pathfinder.js',
  './js/price.js',
  './js/map.js',
  './js/i18n.js',
  './data/stations.json',
  './data/trains.json',
  './data/meta.json',
  './data/coordinates.json',
];

const CDN_ASSETS = [
  'https://telegram.org/js/telegram-web-app.js',
  'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css',
  'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js',
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_VERSION).then(async (cache) => {
      await cache.addAll(LOCAL_ASSETS);
      for (const url of CDN_ASSETS) {
        try { await cache.add(url); } catch { /* CDN unreachable — skip */ }
      }
    })
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((names) =>
      Promise.all(
        names
          .filter((name) => name !== CACHE_VERSION)
          .map((name) => caches.delete(name))
      )
    )
  );
  self.clients.claim();
});

self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;
  event.respondWith(
    caches.match(event.request).then((cached) => {
      if (cached) return cached;
      return fetch(event.request).then((response) => {
        if (response.ok) {
          const clone = response.clone();
          caches.open(CACHE_VERSION).then((cache) => cache.put(event.request, clone));
        }
        return response;
      });
    })
  );
});
