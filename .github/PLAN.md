# Plan: Wire auth action everywhere — one source of auth

## Context

`auth` composite action is the single source of all authentication. No inline auth anywhere — not in composite actions, not in workflows.

**Two job variants per logical operation, gated on `inputs.dry_run`:**
- **Default job** (`if: inputs.dry_run != true`): write permissions, no dry-run steps, only default execution
- **Dry-run job** (`if: inputs.dry_run == true`): read-only permissions, dry-run plan/summary steps, no default-mode work

`inputs.dry_run` for tag push is empty → `!= true` → default jobs run. For `workflow_dispatch` with `dry_run: true` → dry-run jobs run.

## Behavior changes (intentional)

- **ACT no longer forces dry-run.** Currently the `dry-run` action overrides `dry_run` to `true` when running under `act`. With job-level gating, running `act` with `dry_run=false` now executes the default job. `auth` still handles ACT for credential source (local creds vs WIF). Local users must explicitly pass `dry_run=true` for local dry-run.
- **`dry-run` action stays in dry-run jobs only**, as an explicit marker step. Removed from default jobs and from inside composite actions.

## Composite action changes

### `.github/actions/build-and-push/action.yml`
- **Remove inputs**: `github_token`, `gcloud_identity_provider`, `gcloud_service_account`
- **Keep inputs**: `gcloud_region`, `gcloud_project_id`, `dry_run`
- **Remove steps**: "Login to GHCR", "Authenticate to Google Cloud", "Configure Docker for Artifact Registry", and the internal `./.github/actions/dry-run` call
- Replace `steps.dry_run.outputs.dry_run` with `inputs.dry_run` throughout
- Update dry-run summary table: drop rows for the three removed inputs

### `.github/actions/deploy/action.yml`
- **Remove inputs**: `gcloud_identity_provider`, `gcloud_service_account`
- **Remove steps**: "Authenticate to Google Cloud", and the internal `./.github/actions/dry-run` call
- Replace `steps.dry_run.outputs.dry_run` with `inputs.dry_run`

## Workflow: `.github/workflows/release.yml`

Replace the current 3 jobs with **6 jobs** (3 logical operations × 2 modes):

### `build-and-push-dry-run` (new)
```yaml
if: inputs.dry_run == true
permissions: { contents: read, packages: read, id-token: write }
env: { GCLOUD_PROJECT_ID, GCLOUD_IDENTITY_PROVIDER, GCLOUD_SERVICE_ACCOUNT, GCLOUD_REGION (all from vars.*) }
outputs:
  gcloud_tag: ${{ steps.build.outputs.gcloud_tag }}
steps:
  - uses: actions/checkout@v4
  - uses: ./.github/actions/dry-run        # marker step in dry-run mode
    with: { dry_run: true }
  - uses: ./.github/actions/auth
  - id: build
    uses: ./.github/actions/build-and-push
    with: { gcloud_project_id, gcloud_region, dry_run: 'true' }
```

### `build-and-push` (replaces existing)
```yaml
if: inputs.dry_run != true
permissions: { contents: read, packages: write, id-token: write }
env: { GCLOUD_PROJECT_ID, GCLOUD_IDENTITY_PROVIDER, GCLOUD_SERVICE_ACCOUNT, GCLOUD_REGION }
outputs:
  gcloud_tag: ${{ steps.push.outputs.gcloud_tag }}
steps:
  - uses: actions/checkout@v4
  # NO dry-run step here — default job
  - uses: ./.github/actions/auth
  - id: push
    uses: ./.github/actions/build-and-push
    with: { gcloud_project_id, gcloud_region, dry_run: 'false' }
```

### `deploy-to-preprod-dry-run` (new)
```yaml
if: inputs.dry_run == true
needs: build-and-push-dry-run
environment: preprod
permissions: { contents: read, id-token: write }
env: { GCLOUD_* }
steps:
  - uses: actions/checkout@v4
  - uses: ./.github/actions/dry-run
    with: { dry_run: true }
  - uses: ./.github/actions/auth
  - uses: ./.github/actions/deploy
    with:
      gcloud_image: ${{ needs.build-and-push-dry-run.outputs.gcloud_tag }}
      dry_run: 'true'
      # ...other inputs (env_var_environment, gcloud_project_id, gcloud_region, gcloud_run_service, gcloud_secret_version_telegram_token)
      # NO gcloud_identity_provider, NO gcloud_service_account
```

### `deploy-to-preprod` (updated)
```yaml
if: inputs.dry_run != true
needs: build-and-push
environment: preprod
permissions: { contents: read, id-token: write }
env: { GCLOUD_* }
steps:
  - uses: actions/checkout@v4
  - uses: ./.github/actions/auth
  - uses: ./.github/actions/deploy
    with:
      gcloud_image: ${{ needs.build-and-push.outputs.gcloud_tag }}
      dry_run: 'false'
      # ...other inputs (no auth-related ones)
```

### `deploy-to-prod-dry-run` (new)
```yaml
if: inputs.dry_run == true
needs: [deploy-to-preprod-dry-run, build-and-push-dry-run]
environment: prod
# same shape as deploy-to-preprod-dry-run
```

### `deploy-to-prod` (updated)
```yaml
if: inputs.dry_run != true
needs: [deploy-to-preprod, build-and-push]
environment: prod
# same shape as deploy-to-preprod
```

## Workflow: `.github/workflows/ci-image.yml`

Replace the single `publish` job with **two variants**:

### `publish-dry-run`
```yaml
if: inputs.dry_run == true
permissions: { contents: read, packages: read, id-token: write }
env: { TDLIB_COMMIT, GCLOUD_* }   # GCLOUD_* needed because auth runs them
steps:
  - uses: actions/checkout@v4
  - uses: ./.github/actions/dry-run
    with: { dry_run: true }
  - uses: ./.github/actions/auth
  - uses: docker/build-push-action@v6
    with: { ..., push: false, load: true }
  - dry-run summary step
```

### `publish`
```yaml
if: inputs.dry_run != true
permissions: { contents: read, packages: write, id-token: write }
env: { TDLIB_COMMIT, GCLOUD_* }
steps:
  - uses: actions/checkout@v4
  - uses: ./.github/actions/auth
  - uses: docker/build-push-action@v6
    with: { ..., push: true }
```

> `ci-image` only pushes to GHCR but `auth` always runs `gcloud auth configure-docker`, so the job needs `id-token: write` and the `GCLOUD_*` env vars even though it has no GCloud Artifact Registry interaction. Trade-off accepted for single-source-of-auth.

## Notes / known trade-offs

- **Same WIF service account in both dry-run and default jobs**. True least-privilege would need a separate read-only SA for dry-run jobs; deferred — out of scope.
- **6 jobs in release.yml** (vs 3 today) means more job-startup overhead and some YAML duplication of env+checkout+auth blocks. GitHub Actions doesn't have great YAML anchor support; reusable workflows would add more abstraction than it's worth here. Accepted.
- **`dry-run` composite action** stays alive but is now only called from dry-run jobs as an explicit "this is dry-run" marker. Could be simplified or removed in a follow-up.

## Verification

1. **Local act with `dry_run=true`**: dry-run jobs run, auth uses local credentials, plan summaries printed, no push/deploy.
2. **Local act with `dry_run=false`** (new behavior): default jobs run, auth uses local credentials, actual push/deploy attempted with local creds. User must opt in to dry-run explicitly for local safety.
3. **GitHub `workflow_dispatch` with `dry_run=true`**: dry-run jobs run (`packages: read`, WIF), default jobs skipped.
4. **GitHub tag push**: default jobs run (`packages: write`), full pipeline executes.
