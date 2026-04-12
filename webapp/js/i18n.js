// i18n.js — lightweight translation lookup. Keys live here instead of in
// separate files so the prototype stays a static file drop with zero build
// step. Parity with the bot's render package: same language codes, same
// spirit, but webapp-specific copy for UI chrome.

const DICT = {
  en: {
    "tab.search":  "Search",
    "tab.map":     "Map",
    "tab.alerts":  "Alerts",
    "tab.saved":   "Saved",
    "tab.info":    "Info",

    "search.title":    "Where to?",
    "search.from":     "From",
    "search.to":       "To",
    "search.date":     "Date",
    "search.today":    "Today",
    "search.tomorrow": "Tomorrow",
    "search.cta":      "Find trains",
    "search.popular":  "Popular routes",
    "search.pick":     "Pick a station",

    "picker.search":   "Search station…",
    "picker.blacklisted": "No train station here",

    "results.direct":      "Direct",
    "results.interchange": "Change in Podgorica",
    "results.none":        "No trains for this route",
    "results.duration":    "duration",
    "results.wait":        "wait",
    "results.price_est":   "Est. price",
    "results.details":     "View all stops",
    "results.reverse":     "Reverse",
    "results.save":        "Save route",

    "map.title":       "Stations map",
    "map.legend.main": "Main stations",
    "map.legend.stop": "Stops & halts",
    "map.legend.sel":  "Selected",
    "map.legend.line": "Route line",
    "map.tap":         "Tap a station to pick it",
    "map.offline":     "Map unavailable offline",
    "map.error":       "Could not load map",

    "alerts.title":         "Alerts",
    "alerts.empty":         "No active alerts.",
    "alerts.new":            "Create alert",
    "alerts.subscribe":      "Notify me",
    "alerts.unsubscribe":    "Cancel",
    "alerts.disruption":     "Service disruption",

    "saved.title": "Saved routes",
    "saved.empty": "Saved routes appear here.",
    "saved.clear": "Clear all",
    "saved.remove":"Remove",

    "info.title":    "About",
    "info.bot":      "Open in Telegram",
    "info.zpcg":     "Official timetable (zpcg.me)",
    "info.site":     "Railways of Montenegro",
    "info.source":   "Source code",
    "info.updated":  "Timetable last updated",
    "info.lang":     "Language",
    "info.theme":    "Theme follows Telegram",
    "info.legal":    "Coordinates and prices in this prototype are approximate.",

    "price.class":    "Class",
    "price.class1":   "1st",
    "price.class2":   "2nd",
    "price.adult":    "Adult",
    "price.child":    "Child (0–6)",
    "price.student":  "Student",
    "price.senior":   "Senior (65+)",
    "price.total":    "Total",
    "price.note":     "Estimate based on €0.06/km — verify at ticket office.",

    "station.type.1": "Station",
    "station.type.2": "Stop",
    "station.type.3": "Crossing",
    "station.type.4": "Main station",
  },
  ru: {
    "tab.search":  "Поиск",
    "tab.map":     "Карта",
    "tab.alerts":  "Оповещения",
    "tab.saved":   "Избранное",
    "tab.info":    "О боте",

    "search.title":    "Куда едем?",
    "search.from":     "Откуда",
    "search.to":       "Куда",
    "search.date":     "Дата",
    "search.today":    "Сегодня",
    "search.tomorrow": "Завтра",
    "search.cta":      "Найти поезда",
    "search.popular":  "Популярные маршруты",
    "search.pick":     "Выберите станцию",

    "picker.search":   "Поиск станции…",
    "picker.blacklisted": "Здесь нет ж/д станции",

    "results.direct":      "Прямой",
    "results.interchange": "Пересадка в Подгорице",
    "results.none":        "Поезда не найдены",
    "results.duration":    "в пути",
    "results.wait":        "ожидание",
    "results.price_est":   "Примерная цена",
    "results.details":     "Все остановки",
    "results.reverse":     "Обратно",
    "results.save":        "Сохранить",

    "map.title":       "Карта станций",
    "map.legend.main": "Главные станции",
    "map.legend.stop": "Остановочные пункты",
    "map.legend.sel":  "Выбрано",
    "map.legend.line": "Линия маршрута",
    "map.tap":         "Нажмите на станцию, чтобы выбрать",
    "map.offline":     "Карта недоступна офлайн",
    "map.error":       "Не удалось загрузить карту",

    "alerts.title":         "Оповещения",
    "alerts.empty":         "Активных оповещений нет.",
    "alerts.new":           "Создать оповещение",
    "alerts.subscribe":     "Подписаться",
    "alerts.unsubscribe":   "Отписаться",
    "alerts.disruption":    "Изменения в расписании",

    "saved.title": "Сохранённые маршруты",
    "saved.empty": "Сохранённые маршруты появятся здесь.",
    "saved.clear": "Очистить",
    "saved.remove":"Удалить",

    "info.title":    "О приложении",
    "info.bot":      "Открыть в Telegram",
    "info.zpcg":     "Официальное расписание (zpcg.me)",
    "info.site":     "Железные дороги Черногории",
    "info.source":   "Исходный код",
    "info.updated":  "Расписание обновлено",
    "info.lang":     "Язык",
    "info.theme":    "Тема по настройкам Telegram",
    "info.legal":    "Координаты и цены в прототипе приблизительные.",

    "price.class":    "Класс",
    "price.class1":   "1",
    "price.class2":   "2",
    "price.adult":    "Взрослый",
    "price.child":    "Дети (0–6)",
    "price.student":  "Студент",
    "price.senior":   "Пенсионер 65+",
    "price.total":    "Итого",
    "price.note":     "Оценка €0.06/км — уточняйте в кассе.",

    "station.type.1": "Станция",
    "station.type.2": "Остановка",
    "station.type.3": "Разъезд",
    "station.type.4": "Главная станция",
  },
  sr: {
    "tab.search":  "Pretraga",
    "tab.map":     "Mapa",
    "tab.alerts":  "Obaveštenja",
    "tab.saved":   "Sačuvano",
    "tab.info":    "Info",

    "search.title":    "Gde putuješ?",
    "search.from":     "Od",
    "search.to":       "Do",
    "search.date":     "Datum",
    "search.today":    "Danas",
    "search.tomorrow": "Sutra",
    "search.cta":      "Pronađi vozove",
    "search.popular":  "Popularne rute",
    "search.pick":     "Izaberi stanicu",

    "picker.search":   "Pretraži stanicu…",
    "picker.blacklisted": "Ovde nema železničke stanice",

    "results.direct":      "Direktno",
    "results.interchange": "Presedanje u Podgorici",
    "results.none":        "Nema vozova za ovu rutu",
    "results.duration":    "trajanje",
    "results.wait":        "čekanje",
    "results.price_est":   "Procena cene",
    "results.details":     "Sve stanice",
    "results.reverse":     "Obratno",
    "results.save":        "Sačuvaj rutu",

    "map.title":       "Mapa stanica",
    "map.legend.main": "Glavne stanice",
    "map.legend.stop": "Stajališta",
    "map.legend.sel":  "Izabrano",
    "map.legend.line": "Linija rute",
    "map.tap":         "Dodirni stanicu za izbor",
    "map.offline":     "Mapa nedostupna offline",
    "map.error":       "Nije moguće učitati mapu",

    "alerts.title":         "Obaveštenja",
    "alerts.empty":         "Nema aktivnih obaveštenja.",
    "alerts.new":            "Kreiraj obaveštenje",
    "alerts.subscribe":      "Prati",
    "alerts.unsubscribe":    "Otkaži",
    "alerts.disruption":     "Prekid saobraćaja",

    "saved.title": "Sačuvane rute",
    "saved.empty": "Sačuvane rute će se pojaviti ovde.",
    "saved.clear": "Obriši sve",
    "saved.remove":"Ukloni",

    "info.title":    "O aplikaciji",
    "info.bot":      "Otvori u Telegramu",
    "info.zpcg":     "Zvaničan red vožnje (zpcg.me)",
    "info.site":     "Železnički prevoz Crne Gore",
    "info.source":   "Izvorni kod",
    "info.updated":  "Red vožnje ažuriran",
    "info.lang":     "Jezik",
    "info.theme":    "Tema prati Telegram",
    "info.legal":    "Koordinate i cene u prototipu su približne.",

    "price.class":    "Klasa",
    "price.class1":   "1.",
    "price.class2":   "2.",
    "price.adult":    "Odrasli",
    "price.child":    "Dete (0–6)",
    "price.student":  "Student",
    "price.senior":   "Penzioner 65+",
    "price.total":    "Ukupno",
    "price.note":     "Procena €0.06/km — proveri na kasi.",

    "station.type.1": "Stanica",
    "station.type.2": "Stajalište",
    "station.type.3": "Ukrsnica",
    "station.type.4": "Glavna stanica",
  },
};

// Supported codes the bot also supports. Unsupported → fall back to English.
export const LANGS = [
  { code: 'en', label: 'EN', name: 'English'    },
  { code: 'ru', label: 'RU', name: 'Русский'    },
  { code: 'sr', label: 'SR', name: 'Srpski'     },
];

let current = 'en';

export function setLang(code) {
  if (DICT[code]) current = code;
  else current = 'en';
  try { localStorage.setItem('monterails.lang', current); } catch {}
  document.documentElement.lang = current;
  applyStaticTranslations();
}

export function getLang() { return current; }

export function t(key, fallback = key) {
  return (DICT[current] && DICT[current][key])
      || (DICT.en && DICT.en[key])
      || fallback;
}

export function initLang(preferredCode) {
  let code = null;
  try { code = localStorage.getItem('monterails.lang'); } catch {}
  if (!code && preferredCode) code = preferredCode.slice(0, 2).toLowerCase();
  if (!DICT[code]) code = 'en';
  setLang(code);
}

// Replace [data-i18n] and [data-i18n-ph] attributes in static markup.
export function applyStaticTranslations(root = document) {
  root.querySelectorAll('[data-i18n]').forEach(el => {
    el.textContent = t(el.dataset.i18n, el.textContent);
  });
  root.querySelectorAll('[data-i18n-ph]').forEach(el => {
    el.setAttribute('placeholder', t(el.dataset.i18nPh, el.placeholder));
  });
}
