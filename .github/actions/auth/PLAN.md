### [`auth`](actions/auth)

Composite action for centralizing authentication. `build-and-push` uses it before registry
manifest checks so dry-run can verify private GHCR and Artifact Registry state with
read-only credentials. `ci-image` and `deploy` still keep their narrower inline auth.

The action has a deliberately small contract: determine the effective dry-run mode for its
own authentication policy, authenticate for that mode, and fail if authentication cannot be
established. It should not return a matrix of partial-auth booleans. Either the action
authenticated correctly for the current mode, or the action fails.

Unlike the existing `dry-run` action, this action intentionally reads the `ACT` environment
variable directly. Authentication policy depends on whether the run is local: local `act`
runs cannot use GitHub OIDC / Workload Identity Federation, so they must use credentials
provided through local environment files or preconfigured local tools.

**Inputs**

| Name      | Required | Description                                                                                                                   |
|-----------|----------|-------------------------------------------------------------------------------------------------------------------------------|
| `dry_run` | No       | Requested dry-run value (`true` or `false`). Default: `false`. If `ACT=true`, the effective value is always forced to `true`. |

Use the existing `dry_run` input name to stay consistent with the current workflows. In prose,
this is the dry-run flag.

**Outputs**

No outputs.

The effective `dry_run` output already belongs to `.github/actions/dry-run`. Workflows that
need a reusable `dry_run` value for later `if:` guards should keep calling `dry-run`
separately and pass `steps.dry_run.outputs.dry_run` into `auth`.

There are no outputs for dry-run mode, access mode, GHCR login, GCloud login, Docker
configuration, or whether authentication happened. The action must fail when required
authentication was not performed successfully.

**Credential sources**

Credentials are not action inputs. They are provided by the workflow environment, repository
settings, GitHub context, or local `act` env files.

| Mode                                             | Credential source                                                                                                                                       |
|--------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| Local dry-run (`ACT=true`)                       | Local env files or locally configured tools. Use read-only GHCR credentials and local GCloud credentials because WIF/OIDC is not available under `act`. |
| Remote dry-run (`dry_run=true`, `ACT` not set)   | Read-only credentials from repository or environment settings. GHCR must be read-only; GCloud must use a read-only Artifact Registry identity.          |
| Remote real run (`dry_run=false`, `ACT` not set) | Normal workflow credentials: `github.token`/GHCR package permissions and the existing WIF service account used for publish/deploy.                      |
