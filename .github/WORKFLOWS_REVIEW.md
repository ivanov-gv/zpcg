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
| L1 | LOW | **DO** | Standardize all composite-action inputs on `dry_run` (underscore) |
| L2 | LOW | **N/A (rechecked)** | See recheck note below — not a real issue |
| L5 | LOW | **DO** (repo var) | Move `TDLIB_COMMIT` to a repository variable |
| L3, L4, L6, L7, L8, L9 | LOW | **DO NOT TOUCH** | Out of scope |

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

### L1 — standardize on `dry_run` (underscore)
- Files: `.github/actions/build-and-push/action.yml`,
  `.github/actions/deploy/action.yml`,
  `.github/workflows/release.yml`,
  `.github/workflows/deploy-to-preprod.yml`,
  `.github/README.md`
- Rename the composite-action input `dry-run` → `dry_run` in `build-and-push`
  and `deploy`. Update every caller (`with: dry-run:` → `with: dry_run:`).
  Update README tables and usage examples. The input on workflows was already
  `dry_run`; the directory name for the composite stays `dry-run` (action
  *name*, not input name).

### L5 — `TDLIB_COMMIT` via repository variable
- Files: `.github/workflows/checks.yml`, `.github/workflows/ci-image.yml`,
  `Makefile`, `.github/act/checks.env.example`,
  `.github/act/ci-image.env.example`
- Replace the hardcoded `TDLIB_COMMIT: 971684a` in both workflows with
  `TDLIB_COMMIT: ${{ vars.TDLIB_COMMIT }}`. Add `--var-file` to
  `test-ci-checks` and `test-ci-ci-image` Makefile targets so `act` picks up
  the variable from the env-example files. Add a `TDLIB_COMMIT=971684a` line
  to both env-example files, mirroring the pattern already used for
  deploy-to-preprod / release.

---

## Recheck — L2 (`packages: write` on `deploy-to-preprod.yml`)

User interpretation confirmed: `build-and-push` probes both registries and
picks the cheapest path (build / pull / retag), so re-dispatching
`deploy-to-preprod` for an already-published commit does not rebuild. On a
`workflow_dispatch` re-run where the image is in both registries, the
retag path is reached *and* `docker/metadata-action` produces an empty tag
set (there is no `type=ref,event=tag` match outside a tag push), so the
`Retag in-registry` step is skipped entirely by its
`steps.meta.outputs.tags != ''` guard. Net effect: **zero GHCR writes on a
repeat preprod dispatch.**

`packages: write` is still requested upfront because GitHub Actions does not
support conditional job permissions; dropping it would mean splitting
`deploy-to-preprod` into two jobs (build-and-push with packages:write vs a
pure deploy job with id-token only). That's extra pipeline complexity for a
permission that goes unused in a cold-path scenario — not a real issue.

**Verdict**: L2 is not a real issue. No change to make.

---

## Out of scope for this iteration

- **H1 / H2**: superseded by the local-run design discussion. The preferred
  path is host ADC + impersonation of a read-only GCloud SA, with the WIF
  auth step gated on `env.ACT != 'true'`. Treated as a design decision to
  implement separately, not as a bug in the current workflow code.
- **M4, M5, M7, M8**: deferred.
- **L3, L4, L6, L7, L8, L9**: explicitly not touched.

---

## Pre-merge action required (PR description warning)

Before merging this branch, add a **repository variable** in GitHub so the
workflows can resolve `vars.TDLIB_COMMIT`:

> ⚠️ **Before merging**: add a repository variable `TDLIB_COMMIT` with value
> `971684a` under **Settings → Secrets and variables → Actions →
> Variables → New repository variable**. Without it, `checks.yml` pulls an
> invalid `ci:tdlib-` tag and `ci-image.yml` builds with an empty TDLib
> commit arg.
>
> To bump the TDLib version later, update this repository variable and run
> the `CI Image` workflow — no code change required.

---

## Files touched in this iteration

- `.github/actions/build-and-push/action.yml` — H3, M6, L1
- `.github/actions/build/action.yml` — M1 (+ buildx setup so it runs standalone)
- `.github/actions/deploy/action.yml` — L1
- `.github/workflows/checks.yml` — M2, M3, L5
- `.github/workflows/ci-image.yml` — L5
- `.github/workflows/release.yml` — M2, L1
- `.github/workflows/deploy-to-preprod.yml` — M2, L1
- `.github/README.md` — L1
- `Makefile` — L5 (act `--var-file` for checks / ci-image targets)
- `.github/act/checks.env.example`, `.github/act/ci-image.env.example` — L5

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
