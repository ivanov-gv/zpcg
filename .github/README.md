# `.github/` conventions

This document describes the conventions used by the workflows in
`.github/workflows/` and the composite actions in `.github/actions/`.
Most of them look like duplication or footguns at first glance — they aren't.
Read this before "simplifying" anything here.

# TL;DR

Every workflow:

1. Must have 3 ways to run with exactly the same steps and env variables definitions:
    1. Locally with `nektos/act`
    2. On GitHub Actions using `workflow_dispatch` with `inputs.test_run`
    3. Default production run
2. Must define `env:` sections in job's scope that map `secrets.*` and `vars.*` into named env vars
3. Must use only variables from the closest `env:` block or `github.*` context in `steps:` section
4. Must only use variables from `env:` blocks via `${{ env.* }}`
5. Must not use `test_run` for branching in `steps:` sections. Steps must not know whether they're
   running in test-run mode or not.
6. Must use test_run mode if running locally with `nektos/act`.

Workflows are strictly divided into CI and CD parts.

CI workflows:

1. Must not use a pre/production environment in any way.
2. Must use only GHCR.io artifactory, which is considered a part of the development environment.
3. Must not create tags or releases.

CD workflows:

1. Must have `cd-` prefix in the name.
2. Must be triggered by a release or workflow_dispatch events.
3. Must not be triggered by a push event of any kind, including tags.
4. Must use a pre/production environment in default mode.
5. Must not use a pre/production environment in test-run mode, only test-run environment.
6. Must only pull Docker images from GHCR.io, never push to GHCR.io.
7. Must only pull Docker images from GHCR.io with `ghcr.io/${{ github.repository }}` path prefix and tags, related to
   the release it was triggered by – either the release tag or commit SHA.
8. Must never use GHCR.io directly to push to pre/production environments. Must always create a copy of the image
   in the environment it deploys to.
9. Must never delete any image from pre/production environments.

# Workflows

<details>
<summary>An example:</summary>

```yaml

name: CI - push to ghcr.io
on:
  push:
    tags:
      - 'v*.*.*'
  workflow_dispatch:
    inputs:
      test_run:
        description: Test run
        type: boolean
        default: true

jobs:
  switch:
    name: Test-run switch
    runs-on: ubuntu-latest
    outputs:
      test_run: ${{ steps.test_run.outputs.test_run }}
    steps:
      - uses: actions/checkout@v6
      - id: test_run
        uses: ./.github/actions/test_run
        with:
          test_run: ${{ inputs.test_run }}

  build-and-push:
    name: Build and push ZPCG image
    needs: switch
    if: needs.switch.outputs.test_run != 'true'
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
      id-token: write
    environment: dev
    env: &build-and-push-env
      IMAGE_TAG_SHA: ${{ vars.ENV_REGISTRY }}/${{ github.repository }}:${{ github.sha }}
      IMAGE_TAG_TAG: ${{ vars.ENV_REGISTRY }}/${{ github.repository }}:${{ github.ref_name }}
    steps: &build-and-push-steps
      - uses: actions/checkout@v6

      - name: Login to GHCR
        uses: docker/login-action@v4
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ github.token }}

      - name: Build image
        uses: docker/build-push-action@v7
        with:
          context: .
          file: deploy/Dockerfile
          load: 'true'
          tags: |
            ${{ env.IMAGE_TAG_SHA }}
            ${{ env.IMAGE_TAG_TAG }}

      - name: Push
        run: |
          docker push ${{ env.IMAGE_TAG_SHA }}
          docker push ${{ env.IMAGE_TAG_TAG }}

  build-and-push-test-run:
    name: Build CI image (test-run)
    needs: switch
    if: needs.switch.outputs.test_run == 'true'
    runs-on: ubuntu-latest
    environment: test-run
    permissions:
      contents: read
      packages: read # <-- test-run mode
      id-token: write
    env: *build-and-push-env
    steps: *build-and-push-steps


```

</details>

## CI vs CD

#TODO: write this up

## `env:` blocks

```yaml

env:
  GLOBAL_ENV_VAR: value

jobs:
  - job1:
      env: # maps outside secrets and vars to job-scoped env vars
        - DOESNT_WORK: ${{ env.GLOBAL_ENV_VAR }} # <-- substitution from global env: DOESN'T WORK in ACT
        - JOB1_VAR2: ${{ secrets.GLOBAL_SECRET }}
        - JOB1_VAR2: ${{ vars.GLOBAL_VAR }}
      steps:
        - run: |
            LOCAL_VAR=${{ env.JOB1_VAR }} # explicitly shows that env.JOB1_VAR is defined in job1's env: section
            LOCAL_VAR2=LOCAL_VAR          # explicitly shows that LOCAL_VAR is a local variable
```

Each job declares an `env:` block that maps `secrets.*` and `vars.*` into named env vars. Steps then reference only
`${{ env.* }}` in `with:` and `run:` blocks — never `secrets.*` or `vars.*` directly. The benefit is that each step is a
**self-contained block** — to understand
or edit a step you don't have to scroll up and cross-reference secret/var definitions.

## Test-run model

Workflow diagram with a test-run switch:

```mermaid

graph LR
    A[test-run switch] --> B[build]
    B --> C[push]
    C --> D[retag]
    D --> E[deploy]
    A --> B1[test-run: build]
    B1 --> C1[test-run: retag]
    C1 --> D1[test-run: deploy]
    D1 --> E1[test-run: deploy]
```

Implementation:

```yaml
jobs:
  - switch:
      outputs:
        test_run: ${{ steps.test_run_switch.outputs.test_run }} # true or false
      steps:
        - # decide - default or test-run
        - # set outputs.test_run to true or false 
  - build:
      name: Deploy to preprod
      needs: [ switch ]
      if: ${{ needs.switch.outputs.test_run != 'true' }} # depends on switch
      env: &build-env       # env anchor
      # ...
      steps: &build-steps   # steps anchor
      # ...

  - deploy-test-run:
      name: (test-run) Deploy to preprod
      needs: [ switch ]
      if: ${{ needs.switch.outputs.test_run == 'true' }}
      environment: test-run # test environment
      permissions:
        registry: read      # protects from accidental writes
      env: *build-env       # same envs
      steps: *build-steps   # same steps
```

Production and test-run jobs are deployed as **parallel pairs** that share a YAML-anchor step list
(e.g. `*build-and-push-steps`, `*retag-steps`, `*deploy-steps`). The pair is selected by
`needs.switch.outputs.test_run`. Both jobs run the **same steps**.

Safety comes from `environment: test-run` resolving secrets and vars to non-production targets —
for example `vars.ENV_REGISTRY` resolves to a public no-auth sink like `ttl.sh` in the test-run
environment, and WIF credentials resolve to a test-run GCP project. Reduced `permissions:`
(e.g. `packages: read` on GHCR) provide a second layer.

This means a step that "actually runs" in the test-run job is fine, as long as its destination
comes from environment-scoped secrets or vars. Do **not** add step-level `if:` skips to make
test-run steps "no-op" — that defeats the purpose of exercising the same code path.

The one exception is `golang-tdlib-image-build.yml`'s Push step, see below.

### Outputs don't work

#TODO: write this up

Fan-in jobs preserve outputs across the test-run pair

`cd-pre-release.yml` has a `retag-out` job (lines ~94-102) whose only role is to coalesce
`outputs.image_tag` from the two mutually-exclusive `retag` and `retag-test-run` jobs into a
single output that the `deploy` job can consume. It looks like a no-op but it is **load-bearing**:
without it, the `deploy` job would need conditional logic to pick the right upstream, and the
test-run path would break.
