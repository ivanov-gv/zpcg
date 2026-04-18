# CI/CD

Three composable actions cover the full delivery lifecycle; four workflows wire them together.

```
[checks]           push / PR gate   → build check + go test -race + golangci-lint
[ci-image]         manual, rare     → build & publish the CI base image
[deploy-to-preprod] manual          → build-and-push → deploy to preprod
[release]          v*.*.* tag       → build-and-push → deploy preprod → deploy prod
```

All Google Cloud authentication uses **Workload Identity Federation** — no service account
key files are generated or stored anywhere.

---

## Actions

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
| Yes | No | Pull from GHCR → tag → push to GCloud only |
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

| Name | Required | Description |
|---|---|---|
| `github_token` | Yes | `github.token` from the caller (requires `packages: write`). Passed explicitly because composite actions cannot read the caller's token directly. |
| `gcloud_project_id` | Yes | GCP project ID |
| `gcloud_region` | Yes | Artifact Registry region (e.g. `europe-west1`) |
| `gcloud_identity_provider` | Yes | Workload Identity Provider resource name |
| `gcloud_service_account` | Yes | Service account email |
| `dry-run` | No | If `true`, print the execution plan without running any mutating steps. Default: `false`. |

**Usage**

```yaml
# Standard deployment
- id: push
  uses: ./.github/actions/build-and-push
  with:
    github_token: ${{ github.token }}
    gcloud_project_id: ${{ vars.GCLOUD_PROJECT_ID }}
    gcloud_region: ${{ vars.GCLOUD_REGION }}
    gcloud_identity_provider: ${{ vars.GCLOUD_IDENTITY_PROVIDER }}
    gcloud_service_account: ${{ vars.GCLOUD_SERVICE_ACCOUNT }}

# Validate pipeline config without pushing anything
- uses: ./.github/actions/build-and-push
  with:
    github_token: ${{ github.token }}
    gcloud_project_id: ${{ vars.GCLOUD_PROJECT_ID }}
    gcloud_region: ${{ vars.GCLOUD_REGION }}
    gcloud_identity_provider: ${{ vars.GCLOUD_IDENTITY_PROVIDER }}
    gcloud_service_account: ${{ vars.GCLOUD_SERVICE_ACCOUNT }}
    dry-run: 'true'
```

**Dry-run mode** prints a full execution plan to the job summary — resolved image paths,
registry existence checks, and which steps would run — without touching any registry.
Useful for validating variable configuration before a real deployment, or for local testing
with [`act`](https://github.com/nektos/act):

```bash
act workflow_dispatch -W .github/workflows/deploy-to-preprod.yml \
  --var GCLOUD_PROJECT_ID=my-project \
  --var GCLOUD_REGION=europe-west1 \
  --var GCLOUD_IDENTITY_PROVIDER=projects/... \
  --var GCLOUD_SERVICE_ACCOUNT=sa@my-project.iam.gserviceaccount.com
```

---

### [`deploy`](actions/deploy)

Authenticates to Google Cloud via Workload Identity Federation and deploys a container
image to Cloud Run. Sets `ENVIRONMENT` as an env var and mounts the Telegram API token
from Secret Manager by version reference.

**Inputs**

| Name | Required | Description |
|---|---|---|
| `gcloud_image` | Yes | Full Artifact Registry image ref to deploy — use `gcloud_tag` output from `build-and-push` |
| `gcloud_project_id` | Yes | GCP project ID |
| `gcloud_region` | Yes | Cloud Run region |
| `gcloud_run_service` | Yes | Cloud Run service name |
| `env_var_environment` | Yes | Value injected as the `ENVIRONMENT` env var into the service |
| `gcloud_secret_version_telegram_token` | Yes | Secret Manager version ref for the Telegram API token |
| `gcloud_identity_provider` | Yes | Workload Identity Provider resource name |
| `gcloud_service_account` | Yes | Service account email |

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
    gcloud_identity_provider: ${{ vars.GCLOUD_IDENTITY_PROVIDER }}
    gcloud_service_account: ${{ vars.GCLOUD_SERVICE_ACCOUNT }}
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
| `ci:tdlib-<commit>` | Immutable, content-addressable reference |
| `ci:latest` | Floating tag used by `checks` jobs |

**When to run**: after modifying `deploy/ci.dockerfile` or bumping `TDLIB_COMMIT` in the
workflow file. This is intentionally a manual step — the CI image changes rarely and
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

## Environment variables and secrets

All configuration lives in GitHub Environment variables (`vars.*`). No credentials are
committed to the repository.

| Variable | Environments | Description |
|---|---|---|
| `GCLOUD_PROJECT_ID` | preprod, prod | GCP project ID |
| `GCLOUD_REGION` | preprod, prod | GCP region for Cloud Run and Artifact Registry |
| `GCLOUD_IDENTITY_PROVIDER` | preprod, prod | Workload Identity Provider resource name |
| `GCLOUD_SERVICE_ACCOUNT` | preprod, prod | Service account email |
| `GCLOUD_RUN_SERVICE` | preprod, prod | Cloud Run service name |
| `ENVIRONMENT` | preprod, prod | Value injected as `ENVIRONMENT` env var into the service |
| `GCLOUD_SECRET_TELEGRAM_APITOKEN_VERSION` | preprod, prod | Secret Manager version ref for the Telegram API token |
