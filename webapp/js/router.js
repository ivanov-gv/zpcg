// router.js — the smallest hash router I could make that still supports
// route params. Each route handler receives a { params } object and must
// return a DocumentFragment or HTMLElement to be mounted into #view-root.

export class Router {
  constructor(rootEl) {
    this.root = rootEl;
    this.routes = []; // [{ pattern: RegExp, paramNames: string[], handler: fn }]
    this.onBeforeMount = null;
    this.currentHash = null;
    window.addEventListener("hashchange", () => this.handle());
  }

  add(path, handler) {
    const paramNames = [];
    const pattern = new RegExp("^" + path.replace(/:([\w]+)/g, (_, name) => {
      paramNames.push(name);
      return "([^/?#]+)";
    }) + "$");
    this.routes.push({ pattern, paramNames, handler });
    return this;
  }

  go(hash) {
    if (window.location.hash === hash) this.handle();
    else window.location.hash = hash;
  }

  async handle() {
    const hash = window.location.hash || "#/";
    this.currentHash = hash;
    const path = hash.slice(1);
    for (const { pattern, paramNames, handler } of this.routes) {
      const m = path.match(pattern);
      if (!m) continue;
      const params = {};
      paramNames.forEach((n, i) => { params[n] = decodeURIComponent(m[i + 1]); });
      try {
        if (this.onBeforeMount) this.onBeforeMount(hash, params);
        const node = await handler({ params });
        this.root.replaceChildren();
        if (node) this.root.appendChild(node);
        this.root.scrollTo?.(0, 0);
      } catch (err) {
        console.error("route error", path, err);
        this.root.textContent = "Error: " + err.message;
      }
      return;
    }
    // unknown route → home
    window.location.hash = "#/";
  }
}
