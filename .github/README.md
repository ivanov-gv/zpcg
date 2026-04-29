# CI/CD

Five composable actions cover the full delivery lifecycle; four workflows wire them together.

```
[checks]            push / PR gate   → build check + go test -race + golangci-lint
[ci-image]          manual, rare     → build & publish the CI base image
[deploy-to-preprod] manual           → build-and-push → deploy to preprod
[release]           v*.*.* tag       → build-and-push → deploy preprod → deploy prod
```

GitHub-hosted Google Cloud authentication uses **Workload Identity Federation**. Local
`act` runs can optionally use a short-lived access token or service account JSON for
read-only dry-run checks.

---

## Actions

### [`dry-run`](actions/dry-run)

Single source of truth for dry-run mode. This is the **only** place in the repository that
reads the `act`-internal `ACT` environment variable.

Computes the effective `dry_run` value and exposes it as a step output. All composite actions
and workflows gate their write-only steps on `steps.dry_run.outputs.dry_run` — no direct `ACT`
references anywhere.

| Condition | Effective `dry_run` |
|---|---|
| `inputs.dry_run == 'true'` | `true` |
| Running locally via `act` (`ACT=true`) | `true` (enforced, regardless of input) |
| `inputs.dry_run == 'false'` or empty (e.g. tag push) | `false` |

When dry-run is active, a `> [!WARNING]` notice is written to the job summary explaining that
write operations are skipped; authentication may still run when read-only steps need private
registry or Google Cloud access.

**Inputs**: `dry_run` (optional, default `'false'`) — pass the `workflow_dispatch` input directly.
**Outputs**: `dry_run`, `local_run` — effective values after act detection.

---

### [`auth`](actions/auth)

Centralized workflow authentication. Workflows call this action directly after `dry-run`;
composite actions do not authenticate themselves.

It can:

| Capability                    | GitHub-hosted CI                                                | Local `act`                                                                                                 |
|-------------------------------|-----------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------|
| GitHub CLI/API                | Exports `GH_TOKEN` / `GITHUB_TOKEN` from `secrets.GITHUB_TOKEN` | Same, from `.github/act/*.env`                                                                              |
| GHCR Docker auth              | `docker/login-action` with `secrets.GITHUB_TOKEN`               | Same, usually with a PAT in `GITHUB_TOKEN`                                                                  |
| Google Cloud auth             | `google-github-actions/auth` with WIF                           | service account JSON via `GCLOUD_CREDENTIALS_JSON`, or access token for Docker-only Artifact Registry reads |
| Artifact Registry Docker auth | Docker login with the generated access token                    | Docker login with `GCLOUD_ACCESS_TOKEN` or `_json_key` credentials                                          |

**Inputs** include `dry_run`, `local_run`, `github_token`, `gcloud_project_id`,
`gcloud_region`, `gcloud_identity_provider`, `gcloud_service_account`,
`gcloud_credentials_json`, and `gcloud_access_token`.

---

### [`build`](actions/build)

Builds `deploy/Dockerfile` into a **local** Docker image (`load: true`, no push).
Callers decide what to do with the result — push, retag, or simply verify the image compiles.

**Why separate from push?** The `checks` workflow only needs to confirm the image builds;
no registry write is required. The `build-and-push` action calls this internally when
a local image is needed, then handles pushing itself.

**Inputs**

| Name | Required | Description |
|---|---|---|
| `ghcr_tag` | No | Tag applied to the built image (e.g. `ghcr.io/org/repo:abc1234`). Omit for a build-only check. |

**Outputs**: `imageid`, `digest`

**Prerequisites**: `actions/checkout`

**Usage**

```yaml
# Build check only — verify the Dockerfile compiles, no registry write
- uses: ./.github/actions/build

# Build and tag for a subsequent push (used internally by build-and-push)
- uses: ./.github/actions/build
  with:
    ghcr_tag: ghcr.io/org/repo:${{ github.sha }}
```

---

### [`build-and-push`](actions/build-and-push)

Builds (if needed) and publishes the image to both **GHCR** and **Google Cloud Artifact Registry**.
Before doing any work, it inspects both registries for the commit-SHA tag and takes the
cheapest available path:

| Image in GHCR | Image in GCloud | Action taken |
|---|---|---|
| No | — | Build locally → tag → push to both |
| Yes | No | Pull from GHCR → tag → push to both (commit-SHA already in GHCR; extra tags added there too) |
| Yes | Yes | **In-registry retag** (`buildx imagetools create`) — zero data transfer |

This means re-running a pipeline for the same commit is essentially free, and the `release`
workflow never rebuilds what `deploy-to-preprod` already pushed.

**Tagging**

Every published image carries a **commit-SHA tag** (`repo:abc1234`) for traceability.
On a `v*.*.*` tag push, `docker/metadata-action` additionally appends a **semantic version
tag** to both registries. For `workflow_dispatch`, only the commit-SHA tag is used.

**Outputs**: `ghcr_tag`, `gcloud_tag` — full image refs with the commit-SHA tag, ready to
pass directly to the `deploy` action.

**Inputs**

| Name                | Required | Description                                                                                                                                   |
|---------------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| `gcloud_project_id` | Yes      | GCP project ID                                                                                                                                |
| `gcloud_region`     | Yes      | Artifact Registry region (e.g. `europe-west1`)                                                                                                |
| `dry_run`           | No       | If `true`, skip remote-write steps (push, retag). Build and pull run normally. Prints an execution plan to the job summary. Default: `false`. |

**Usage**

```yaml
# Standard deployment
- uses: ./.github/actions/auth
  with:
    dry_run: ${{ steps.dry_run.outputs.dry_run }}
    local_run: ${{ steps.dry_run.outputs.local_run }}
    authenticate_ghcr: 'true'
    authenticate_gcloud: 'true'
    configure_docker: 'true'
    github_token: ${{ secrets.GITHUB_TOKEN }}
    gcloud_project_id: ${{ vars.GCLOUD_PROJECT_ID }}
    gcloud_region: ${{ vars.GCLOUD_REGION }}
    gcloud_identity_provider: ${{ vars.GCLOUD_IDENTITY_PROVIDER }}
    gcloud_service_account: ${{ vars.GCLOUD_SERVICE_ACCOUNT }}

- id: push
  uses: ./.github/actions/build-and-push
  with:
    gcloud_project_id: ${{ vars.GCLOUD_PROJECT_ID }}
    gcloud_region: ${{ vars.GCLOUD_REGION }}

# Validate pipeline config without pushing anything
- uses: ./.github/actions/build-and-push
  with:
    gcloud_project_id: ${{ vars.GCLOUD_PROJECT_ID }}
    gcloud_region: ${{ vars.GCLOUD_REGION }}
    dry_run: 'true'
```

**Dry-run mode** prints a full execution plan to the job summary — resolved image paths,
registry existence checks, and which steps would run. Local build and image pull execute
normally; registry push and in-registry retag are skipped.

---

### [`deploy`](actions/deploy)

Deploys a container image to Cloud Run. Sets `ENVIRONMENT` as an env var and mounts the
Telegram API token from Secret Manager by version reference. The caller workflow must run
`auth` first when deployment is not a dry-run.

**Inputs**

| Name                                   | Required | Description                                                                                       |
|----------------------------------------|----------|---------------------------------------------------------------------------------------------------|
| `gcloud_image`                         | Yes      | Full Artifact Registry image ref to deploy — use `gcloud_tag` output from `build-and-push`        |
| `gcloud_project_id`                    | Yes      | GCP project ID                                                                                    |
| `gcloud_region`                        | Yes      | Cloud Run region                                                                                  |
| `gcloud_run_service`                   | Yes      | Cloud Run service name                                                                            |
| `env_var_environment`                  | Yes      | Value injected as the `ENVIRONMENT` env var into the service                                      |
| `gcloud_secret_version_telegram_token` | Yes      | Secret Manager version ref for the Telegram API token                                             |
| `dry_run`                              | No       | If `true`, skip Cloud Run deployment. Prints a plan summary to the job summary. Default: `false`. |

**Usage**

```yaml
- uses: ./.github/actions/deploy
  with:
    gcloud_image: ${{ needs.build-and-push.outputs.gcloud_tag }}
    gcloud_project_id: ${{ vars.GCLOUD_PROJECT_ID }}
    gcloud_region: ${{ vars.GCLOUD_REGION }}
    gcloud_run_service: ${{ vars.GCLOUD_RUN_SERVICE }}
    env_var_environment: ${{ vars.ENVIRONMENT }}
    gcloud_secret_version_telegram_token: ${{ vars.GCLOUD_SECRET_TELEGRAM_APITOKEN_VERSION }}
```

---

## Workflows

### [`checks`](workflows/checks.yml)

**Triggers**: push to `main`, pull requests, `workflow_dispatch`

Three parallel jobs guard every change:

| Job | Command | Notes |
|---|---|---|
| `build` | `actions/build` (Dockerfile) | Catches build errors before they reach a registry |
| `test` | `go test -count=1 -race ./...` | `-race` detects data races at test time; `-count=1` disables result caching |
| `lint` | `golangci-lint v2.11.4` | Pinned version for reproducible results |

`test` and `lint` run inside the pre-built [CI image](#ci-image), so TDLib (a 15+ minute
compile) is already present — no rebuild on every run.

**Permissions**: each job declares only what it needs (`contents: read`, `packages: read`
for the CI image pull), following the principle of least privilege.

**Dependency**: the CI image must be current. Rebuild it whenever `deploy/ci.dockerfile`
or `TDLIB_COMMIT` changes.

---

### [`ci-image`](workflows/ci-image.yml)

**Trigger**: `workflow_dispatch` (manual, infrequent)

Builds `deploy/ci.dockerfile` with a pinned TDLib commit and pushes to GHCR:

| Tag | Purpose |
|---|---|
| `ci:tdlib-<commit>` | Immutable, content-addressable reference used by `checks` jobs |
| `ci:latest` | Floating convenience tag |

**When to run**: after modifying `deploy/ci.dockerfile` or bumping the `TDLIB_COMMIT`
repository variable (GitHub Settings → Variables). This is intentionally a manual step — the CI image changes rarely and
building it on every push would waste significant compute.

---

### [`deploy-to-preprod`](workflows/deploy-to-preprod.yml)

**Trigger**: `workflow_dispatch` (manual)

Accepts an optional `description` input to annotate what is being deployed, then runs
`build-and-push` followed by `deploy` in a single job.

Because `build-and-push` skips the build when the commit image already exists in GHCR,
re-deploying the same commit to preprod is fast — it pulls or retags rather than rebuilds.

**Environment**: `preprod` — configure branch protection rules or required reviewers in
GitHub Settings → Environments to gate access.

**Required variables** (set on the `preprod` environment):
`GCLOUD_PROJECT_ID`, `GCLOUD_REGION`, `GCLOUD_IDENTITY_PROVIDER`, `GCLOUD_SERVICE_ACCOUNT`,
`GCLOUD_RUN_SERVICE`, `ENVIRONMENT`, `GCLOUD_SECRET_TELEGRAM_APITOKEN_VERSION`

---

### [`release`](workflows/release.yml)

**Trigger**: push of a `v*.*.*` tag

Three jobs run sequentially with environment gates:

```
build-and-push  ──►  deploy-to-preprod  ──►  deploy-to-prod
                      environment: preprod     environment: prod
```

The image is **built once** in `build-and-push` and its `gcloud_tag` is propagated as a
job output. Both deploy jobs consume `needs.build-and-push.outputs.gcloud_tag` — they never
rebuild or re-push.

`deploy-to-prod` declares `needs: [deploy-to-preprod, build-and-push]`: the first enforces
deployment order; the second gives it direct access to the `build-and-push` job outputs.

The `prod` environment can be configured with a required manual approval in GitHub Settings
→ Environments, enforcing a human gate between preprod and production.

**Permissions by job**:

| Job | `contents` | `packages` | `id-token` |
|---|---|---|---|
| `build-and-push` | read | **write** | write |
| `deploy-to-preprod` | read | — | write |
| `deploy-to-prod` | read | — | write |

Deploy jobs do not push images, so they do not request `packages: write`.

---

## Local testing with `act`

All workflows can be run locally using [`act`](https://github.com/nektos/act). Makefile
targets handle the flags; `act` automatically enforces dry-run mode by setting `ACT=true`,
which the [`dry-run`](#dry-run) action detects.

```bash
make test-ci-checks            # run checks workflow locally (build + test + lint); set GITHUB_TOKEN in checks.env to pull the CI container
make test-ci-ci-image          # build CI image locally without pushing (WARNING: ~15 min TDLib compile)
make test-ci-deploy-to-preprod # run deploy-to-preprod locally (no registry writes or deployment)
make test-ci-release           # run release workflow locally (no registry writes or deployments)
```

Override the `act` binary if needed (e.g. when using standalone `act` instead of the `gh` extension):

```bash
make test-ci-deploy-to-preprod ACT=act
```

**What dry-run mode means per step type:**

| Step type                                            | Dry-run behaviour                                                          |
|------------------------------------------------------|----------------------------------------------------------------------------|
| Authentication (GitHub, GHCR, Google Cloud)          | Runs when configured so read-only private registry checks can authenticate |
| Docker login for Artifact Registry                   | Runs when `GCLOUD_ACCESS_TOKEN` or `GCLOUD_CREDENTIALS_JSON` is available  |
| Registry existence check (`docker manifest inspect`) | Runs normally (failures silently ignored)                                  |
| Image build (local, `deploy/Dockerfile`)             | Runs normally                                                              |
| Image pull (from GHCR)                               | Runs normally                                                              |
| Registry push / in-registry retag                    | **Skipped**                                                                |
| Cloud Run deployment                                 | **Skipped**                                                                |

**Setup**: each workflow has a corresponding `.env.example` template under `.github/act/`.
Copy and fill in real values before running:

```bash
cp .github/act/deploy-to-preprod.env.example .github/act/deploy-to-preprod.env
# edit .github/act/deploy-to-preprod.env with real variable values
make test-ci-deploy-to-preprod
```

The `*.env` files are gitignored — only the `*.env.example` templates are committed.

---

## Environment variables and secrets

All configuration lives in GitHub Environment variables (`vars.*`). No credentials are
committed to the repository.

| Variable | Environments | Description |
|---|---|---|
| `TDLIB_COMMIT` | repository | TDLib commit SHA used to build and pull the CI container image |
| `GCLOUD_PROJECT_ID` | preprod, prod | GCP project ID |
| `GCLOUD_REGION` | preprod, prod | GCP region for Cloud Run and Artifact Registry |
| `GCLOUD_IDENTITY_PROVIDER` | preprod, prod | Workload Identity Provider resource name |
| `GCLOUD_SERVICE_ACCOUNT` | preprod, prod | Service account email |
| `GCLOUD_RUN_SERVICE` | preprod, prod | Cloud Run service name |
| `ENVIRONMENT` | preprod, prod | Value injected as `ENVIRONMENT` env var into the service |
| `GCLOUD_SECRET_TELEGRAM_APITOKEN_VERSION` | preprod, prod | Secret Manager version ref for the Telegram API token |

Optional local/GitHub secrets:

| Secret                    | Scope                                 | Description                                                                                                  |
|---------------------------|---------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `GITHUB_TOKEN`            | local `act`                           | PAT used for GitHub API/GHCR auth; use `read:packages` for dry-run reads or `write:packages` when publishing |
| `GCLOUD_ACCESS_TOKEN`     | local `act`                           | Short-lived token from `gcloud auth print-access-token`; used only for Artifact Registry Docker auth         |
| `GCLOUD_CREDENTIALS_JSON` | local `act` or fallback GitHub secret | Minified service account key JSON for Google auth; prefer WIF in GitHub-hosted CI                            |
