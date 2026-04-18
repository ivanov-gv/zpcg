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
Findings below are ranked by severity. Each row identifies the file and what
to change.

---

## Findings and proposals

### HIGH

**H1. Local `act` runs of `deploy-to-preprod` / `release` will fail at the GCloud auth step**
- Files: `.github/actions/build-and-push/action.yml`, `.github/actions/deploy/action.yml`
- `google-github-actions/auth@v3` is invoked unconditionally, then later steps
  are gated on `steps.dry_run.outputs.dry_run != 'true'`. Under `act` there is
  no GitHub OIDC token issuer, so WIF token exchange fails and the job errors
  out before the dry-run summary runs.
- README explicitly claims "Authentication … runs normally" under `act`; that
  is only true for static credentials, not WIF.
- **Fix**: gate the `google-github-actions/auth@v3` step (and the Docker GCloud
  configuration that follows) on `steps.dry_run.outputs.dry_run != 'true'`, OR
  additionally skip when `env.ACT == 'true'`. Update the README table to say
  GCloud auth is skipped under `act`.

**H2. `docker manifest inspect` is also ungated and needs GCloud auth to have succeeded**
- File: `.github/actions/build-and-push/action.yml`, `Determine action needed` step
- Under `act` (or any dry run on a fresh runner) `docker manifest inspect
  "$GCLOUD_TAG"` runs after the auth step. If auth was skipped (per H1), this
  lookup against Artifact Registry returns non-zero and the action
  misclassifies the state as `ghcr_exists=false / gcloud_exists=false` →
  `action=build`. Dry-run output is therefore not representative of a real run.
- **Fix**: when dry-run is active, either (a) skip both `manifest inspect`
  probes and set `action=build` explicitly with a note in the summary, or
  (b) keep the probes but don't claim the summary reflects what would really
  happen. Pair with H1.

**H3. `docker buildx imagetools create` in the `retag` path runs without explicit buildx setup**
- File: `.github/actions/build-and-push/action.yml`, `Retag in-registry` step
- On the retag path the `build` action is not invoked, so buildx is never
  initialised via `docker/build-push-action`. It usually works on
  `ubuntu-latest` because buildx ships with Docker 24+, but relying on that
  is fragile — a runner image change will break the whole release path with
  no warning.
- **Fix**: add `docker/setup-buildx-action@v3` as the first step in
  `build-and-push` (before the `Determine action` probe). Cost: ~2–3 s; it
  also enables GHA cache support (see M1).

---

### MEDIUM

**M1. No layer cache for `deploy/Dockerfile` builds**
- File: `.github/actions/build/action.yml`
- Every PR build and every `build-and-push` rebuild is a cold Docker build. On
  PRs, the `checks.build` job currently bears the full cost.
- **Fix**: add `cache-from: type=gha` and `cache-to: type=gha,mode=max` to the
  `docker/build-push-action@v6` call in `build/action.yml`. Requires buildx
  (see H3).

**M2. Missing `concurrency` groups**
- Files: `.github/workflows/checks.yml`, `release.yml`, `deploy-to-preprod.yml`
- `checks` on a PR with rapid pushes runs in parallel against itself, wasting
  minutes. `release` / `deploy-to-preprod` have no serialisation, so two
  simultaneous deploys can race on Cloud Run traffic shifting.
- **Fix**:
  - `checks.yml`: `concurrency: { group: checks-${{ github.ref }}, cancel-in-progress: true }`
  - `release.yml`: `concurrency: { group: release-${{ github.ref }}, cancel-in-progress: false }`
  - `deploy-to-preprod.yml`: `concurrency: { group: deploy-preprod, cancel-in-progress: false }`

**M3. `checks` pulls `ghcr.io/.../ci:latest` — floating, unverified, single point of failure**
- Files: `.github/workflows/checks.yml`, `.github/workflows/ci-image.yml`
- A broken push of `ci:latest` breaks every PR until someone publishes a fix.
  The tag is not smoke-tested before becoming `latest`; there is no rollback.
- **Fix**: two options, pick one.
  - (preferred) pin `checks.yml`'s container to `ci:tdlib-${{ env.TDLIB_COMMIT }}`
    and keep `TDLIB_COMMIT` as the single source of truth for both files.
    Removes the floating tag entirely.
  - Alternatively, add a smoke step in `ci-image.yml` that runs `go version`
    and `pkg-config --exists tdlib` inside the freshly built image before
    tagging `latest`.

**M4. Makefile + `.env.example` doc claim: "dry_run=true skips test/lint in act"**
- Files: `Makefile` (comment on `test-ci-checks`), `.github/act/checks.env.example`
- `checks.yml` has no `if:` on the `test` / `lint` jobs, so they always run
  (and always try to pull the CI container) even with `dry_run=true` under
  `act`. The `dry_run` input on `checks.yml` is completely unwired.
- **Fix**: either (a) wire the input — add `if: inputs.dry_run != 'true'` to
  the `test` and `lint` jobs — or (b) delete the `dry_run` input and update
  the Makefile comment / env-example accordingly. (a) is cheap and makes local
  runs actually cheap.

**M5. Third-party actions pinned by major version, not SHA**
- Files: all composite actions, all workflows
- `actions/checkout@v4`, `docker/login-action@v3`, `docker/build-push-action@v6`,
  `docker/metadata-action@v5`, `google-github-actions/auth@v3`,
  `google-github-actions/deploy-cloudrun@v3`, `golangci/golangci-lint-action@v7`.
  Major tags are mutable — a compromised release would execute with
  `packages: write` + `id-token: write` inside this repo.
- **Fix**: SHA-pin every third-party action and add Dependabot config
  (`.github/dependabot.yml`, `package-ecosystem: github-actions`) to automate
  the bumps. Keep major-tag comments next to the SHAs for readability.

**M6. Short-SHA tag width is non-deterministic**
- File: `.github/actions/build-and-push/action.yml`, `Determine action needed`
- `git rev-parse --short HEAD` uses a repo-global default length that git may
  expand automatically when the corpus is ambiguous. The same commit could be
  tagged `abc1234` today and `abc12345` later, causing a spurious cache miss
  and forcing a rebuild.
- **Fix**: either `git rev-parse --short=7 HEAD`, or derive from
  `"${GITHUB_SHA:0:7}"` — consistent across runs.

**M7. Nested `dry-run` composite double-prints the warning banner**
- Files: all workflows, `build-and-push/action.yml`, `deploy/action.yml`
- Workflow calls `dry-run`; then `build-and-push` internally calls `dry-run`
  again; same for `deploy`. On a dispatched dry run that chains build-and-push
  → deploy, the `> [!WARNING]` banner lands in the summary three times.
- **Fix**: inside `build-and-push` and `deploy`, inline the ACT/dry_run
  evaluation (two lines of bash) and drop the nested `uses:` call. Or, skip
  the warning step when called as a sub-action (gate on presence of an input
  like `nested: true`).

**M8. `release` workflow — workflow_dispatch UX is risky**
- File: `.github/workflows/release.yml`
- `on.workflow_dispatch.inputs.dry_run` default is `true`, but the *tag push*
  trigger passes no input, so the composite resolves to `false` and deploys
  to prod. Fine — but a maintainer who uses `workflow_dispatch` to re-run a
  tag's pipeline and forgets to flip `dry_run=false` will get a surprise
  no-op. The safer-default is good; just make the summary state loudly
  whether real writes occurred (it already does for the dry branch; add a
  matching "real run" header to the non-dry branch for symmetry).

---

### LOW

**L1. Inconsistent input naming: `dry-run` vs `dry_run`**
- `dry-run` composite takes `dry_run` (underscore); `build-and-push` and
  `deploy` expose `dry-run` (hyphen). Pick one; hyphen is the GH convention
  elsewhere but `dry_run` matches the workflow input. Cosmetic.

**L2. `deploy-to-preprod.yml` requests `packages: write` for every dispatch**
- It is required only on the build-local / pull-then-push paths, not on the
  retag path. Simpler to keep; mentioning for completeness.

**L3. `github.repository` not lowercased when composing GHCR refs in `checks.yml`**
- Works for `ivanov-gv/zpcg`, but a fork into an uppercase account would
  silently break. The `build-and-push` action already handles this
  (`tr '[:upper:]' '[:lower:]'`); mirror in `checks.yml`.

**L4. `checks.yml` does not filter by paths**
- A docs-only PR still runs Docker build + container pulls for test/lint.
  Add `paths-ignore: ['**.md', 'docs/**']` or explicit `paths:` to trim.

**L5. `ci-image.yml` + `checks.yml` duplicate `TDLIB_COMMIT: 971684a`**
- Drift risk. Move to a single file (e.g. `.github/tdlib-commit`) read by
  both, or reference via `repository_variables`.

**L6. `release.yml`'s `build-and-push` job outputs only `gcloud_tag`**
- Also expose `ghcr_tag` for symmetry and to enable future rollback tooling.

**L7. `deploy/Dockerfile` builds the app with `CGO_ENABLED=0`**
- Out of scope for CI review but worth flagging: the app depends on TDLib via
  cgo; the production image appears to ship a build that would fail at
  runtime import. Verify this is intentional (separate binary paths?) before
  the next release.

**L8. Removed `docker/setup-qemu-action`**
- The old pipeline set up QEMU for cross-arch builds; the new one doesn't.
  Fine for amd64-only Cloud Run, but explicitly confirm no arm64 target is
  needed.

**L9. Summary noise**
- The `build` composite writes a summary unconditionally, including on PR
  builds where it's low-value. Harmless; trim if summary real estate matters.

---

## Proposed implementation order (when greenlit)

1. H1 + H2 + H3 as one PR — correctness of dry-run/local flow and retag path.
2. M1 + M2 — biggest real CI-time wins.
3. M3 + M4 — stabilise the floating `ci:latest` dependency and fix misleading
   docs.
4. M5 — supply-chain pinning + Dependabot.
5. M6 + M7 — polish and determinism.
6. L-series — opportunistic cleanup.

## Files that will be touched

- `.github/actions/build-and-push/action.yml` — H1, H2, H3, M6, M7
- `.github/actions/build/action.yml` — M1
- `.github/actions/deploy/action.yml` — H1, M7
- `.github/workflows/checks.yml` — M2, M3, M4, L3, L4
- `.github/workflows/ci-image.yml` — M3 (smoke test), L5
- `.github/workflows/deploy-to-preprod.yml` — M2
- `.github/workflows/release.yml` — M2, M8, L6
- `.github/dependabot.yml` (new) — M5
- `Makefile`, `.github/act/checks.env.example` — M4
- `.github/README.md` — keep in sync with H1/M3/M4/M7 changes

## Verification

- Local: `make test-ci-checks`, `make test-ci-deploy-to-preprod`,
  `make test-ci-release`, `make test-ci-ci-image`. Each should complete
  without hitting the WIF auth step error (H1).
- On a throwaway PR: confirm `checks` build job is faster with GHA cache
  (M1), that concurrent pushes cancel the older run (M2), that the
  `ci:tdlib-<commit>` pin is used (M3).
- On a `vX.Y.Z-rc1` tag: dry-dispatch the `release` workflow with
  `dry_run=true` and confirm the summary lists `action=retag` after the first
  real run.
- On the real release tag: verify `build-and-push` runs once, both deploy
  jobs consume `needs.build-and-push.outputs.gcloud_tag`, and `ghcr_tag`
  (L6) if added.

## Notes on what is already good

- Four-composite shape is the right abstraction.
- WIF-only GCloud auth — no key files anywhere.
- Per-job least-privilege permissions, correctly split between
  `build-and-push` (packages: write) and deploys (id-token only).
- Build/pull/retag cost model; image built once per commit.
- Single source of truth for ACT detection (`dry-run` composite).
- `.env.example` templates gitignored with real files.
- README is accurate on architecture; only the dry-run/act claims need
  correcting after H1/M4.
