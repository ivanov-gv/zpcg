### [`auth`](actions/auth)

Composite action for centralizing authentication. `build-and-push` uses it twice — once
in a read-only-permissioned job for dry-run manifest checks, once in a write-permissioned
job for the real publish. `ci-image` and `deploy` still keep their narrower inline auth.

The action has a deliberately small contract: authenticate to GHCR and Google Cloud, or
fail. It does not know whether the surrounding job is a dry-run. Read vs. write access
is the caller's responsibility, expressed through the workflow's `permissions:` block
(for `GITHUB_TOKEN`/GHCR) and through which Workload Identity Federation service account
the job is bound to (for Google Cloud). Either authentication succeeds, or the action
fails.

The action's only branch is on whether it is running under `act`. Local `act` runs
cannot use GitHub OIDC / Workload Identity Federation, so they must use credentials
provided through local environment files or preconfigured local tools. The action reads
the `ACT` environment variable directly to detect this.

**Inputs**

None.

**Outputs**

None. The action must fail when authentication is not performed successfully.

The existing `.github/actions/dry-run` action is unaffected. Workflows that need a
reusable `dry_run` value for `if:` guards keep calling it separately. `auth` and
`dry-run` are independent.

**Credential sources**

Credentials are not action inputs. They are provided by the workflow environment,
repository settings, GitHub context, or local `act` env files.

| Mode                 | Credential source                                                                                                                                                                       |
|----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Local (`ACT=true`)   | Local env files or locally configured tools (`gcloud`, `docker login`). The caller is responsible for supplying credentials matching the surrounding job's intent (read-only or write). |
| Remote (`ACT` unset) | `GITHUB_TOKEN` for GHCR and Workload Identity Federation for Google Cloud. The job's `permissions:` block and the bound WIF service account determine read-only vs. write access.      |
