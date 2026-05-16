# AGENTS.md

Workflow-specific instructions for AI coding agents.

## Workflows are lint-checked in CI

`.github/workflows/` is dry-run-linted on every PR by the `workflows-lint` job in
`pr-checks.yml`, which runs `make test-all-workflows ACT="act -n"` and iterates
over per-workflow `test-<name>` targets in the root `Makefile`.

**When you add or rename a workflow file, you MUST:**

1. Add or rename the matching `test-<workflow-name>` target in `Makefile`
   (mirror the existing `test-pr-checks`, `test-ci`, `test-cd-pre-release`, etc.).
2. List the target as a prerequisite of `test-all-workflows`.

Nothing else catches a missing test target — a new workflow file alone is invisible
to the lint job until it's wired in.

See [`README.md`](README.md) for the broader workflow conventions
(test-run model, env-block discipline, fan-in patterns, etc.).
