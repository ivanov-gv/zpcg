# Review — `ci/test-composite-actions`

## Context

Branch `ci/test-composite-actions` replaces the old `build.yml` / `deploy.yml` /
`stage.yml` pair with a cleaner composite-action architecture: four composite
actions (`dry-run`, `build`, `build-and-push`, `deploy`) glued into four
workflows (`checks`, `ci-image`, `deploy-to-preprod`, `release`). The design is
time-efficient: `checks` runs `test`/`lint` inside a pre-built CI image so the
15-minute TDLib compile is amortised; `build-and-push` does a
`manifest-inspect` lookup against both registries and picks the cheapest of
build-local / pull-from-GHCR / in-registry retag. Local runs via `act` are
wired with `.env.example` templates, a Makefile target per workflow, and a
single `dry-run` composite that centralises `ACT` detection.

Overall: the design is sound and noticeably better than what was on `main`.

---

## Action plan (agreed)

Legend: **DO** = apply in this iteration; **SKIP** = intentionally deferred;
**N/A** = resolved by design discussion, no code change needed.

| ID | Severity | Decision | Summary |
|---|---|---|---|
| H1 | HIGH | N/A | Resolved via planned host-ADC + impersonation approach for local runs |
| H2 | HIGH | N/A | Same — probe becomes accurate (option 1) or is skipped (option 3) |
| H3 | HIGH | **DO** | Add `docker/setup-buildx-action@v3` early in `build-and-push` |
| M1 | MEDIUM | **DO** | Add GHA layer cache to `build` composite |
| M2 | MEDIUM | **DO** | Add `concurrency` groups to `checks`, `release`, `deploy-to-preprod` |
| M3 | MEDIUM | **DO** (pin version) | Pin `checks.yml` container to `ci:tdlib-${{ env.TDLIB_COMMIT }}` |
| M4 | MEDIUM | SKIP | Misleading `dry_run` docs on `checks.yml` — deferred |
| M5 | MEDIUM | SKIP | SHA-pinning + Dependabot — deferred |
| M6 | MEDIUM | **DO** | Change `git rev-parse --short HEAD` to `git rev-parse --short=7 HEAD` |
| M7 | MEDIUM | SKIP | Nested `dry-run` banner duplication — deferred |
| M8 | MEDIUM | SKIP | `release` dispatch UX polish — deferred |
| L1–L9 | LOW | **DO NOT TOUCH** | Out of scope for this iteration |

---

## Changes to apply

### H3 — buildx setup in `build-and-push`
- File: `.github/actions/build-and-push/action.yml`
- Add `docker/setup-buildx-action@v3` as the first step (before `Login to GHCR`),
  so the `retag` path's `docker buildx imagetools create` works regardless of
  runner-image defaults.

### M1 — layer cache for `deploy/Dockerfile` builds
- File: `.github/actions/build/action.yml`
- Add `cache-from: type=gha` and `cache-to: type=gha,mode=max` to the
  `docker/build-push-action@v6` call. Relies on buildx being configured
  (covered by H3 when invoked from `build-and-push`; add
  `docker/setup-buildx-action@v3` at the top of `build` composite too so the
  standalone `checks.build` job gets the cache as well).

### M2 — concurrency groups
- File: `.github/workflows/checks.yml`
  ```yaml
  concurrency:
    group: checks-${{ github.ref }}
    cancel-in-progress: true
  ```
- File: `.github/workflows/release.yml`
  ```yaml
  concurrency:
    group: release-${{ github.ref }}
    cancel-in-progress: false
  ```
- File: `.github/workflows/deploy-to-preprod.yml`
  ```yaml
  concurrency:
    group: deploy-preprod
    cancel-in-progress: false
  ```

### M3 — pin checks to the concrete CI image tag
- File: `.github/workflows/checks.yml`
- Replace `ghcr.io/${{ github.repository }}/ci:latest` with
  `ghcr.io/${{ github.repository }}/ci:tdlib-${{ env.TDLIB_COMMIT }}` on both
  `test` and `lint` jobs. `TDLIB_COMMIT` is already declared at workflow
  top-level, so no new env is needed. The floating `ci:latest` tag remains
  published by `ci-image.yml` but is no longer on the critical path.

### M6 — deterministic short SHA
- File: `.github/actions/build-and-push/action.yml`, `Determine action needed`
  step
- Change `COMMIT_SHORT_SHA="$(git rev-parse --short HEAD)"` to
  `COMMIT_SHORT_SHA="$(git rev-parse --short=7 HEAD)"`.

---

## Out of scope for this iteration

- **H1 / H2**: superseded by the local-run design discussion. The preferred
  path is host ADC + impersonation of a read-only GCloud SA, with the WIF
  auth step gated on `env.ACT != 'true'`. Treated as a design decision to
  implement separately, not as a bug in the current workflow code.
- **M4, M5, M7, M8**: deferred.
- **All LOW (L1–L9)**: explicitly not touched.

---

## Files touched in this iteration

- `.github/actions/build-and-push/action.yml` — H3, M6
- `.github/actions/build/action.yml` — M1 (+ buildx setup so it runs standalone)
- `.github/workflows/checks.yml` — M2, M3
- `.github/workflows/release.yml` — M2
- `.github/workflows/deploy-to-preprod.yml` — M2

No other workflow, action, doc, Makefile, or env-example file is modified.

---

## Verification

- `checks` workflow on a PR: confirm `build` job reports a cache hit on the
  second run (M1), that pushing a new commit cancels the in-flight run (M2),
  and that the container image resolves to `ci:tdlib-<TDLIB_COMMIT>` (M3).
- `release` / `deploy-to-preprod`: confirm a second dispatch queues behind the
  first rather than running in parallel (M2).
- `build-and-push`: on a commit-SHA collision-free state, verify
  `steps.check.outputs.ghcr_tag` ends in a 7-char SHA (M6) and that
  `buildx imagetools create` in the retag path no longer relies on implicit
  buildx (H3).
