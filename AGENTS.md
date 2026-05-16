# AGENTS.md

Project-level instructions for AI coding agents working in this repository.

## Workflows are lint-checked in CI

`.github/workflows/` is dry-run-linted on every PR by the `workflows-lint` job in
`.github/workflows/pr-checks.yml`. That job runs `make test-all-workflows ACT="act -n"`,
which iterates over per-workflow `test-<name>` targets in `Makefile`.

**When you add or rename a workflow file, you MUST:**

1. Add or rename the matching `test-<workflow-name>` target in `Makefile`
   (mirror the existing `test-pr-checks`, `test-ci`, `test-cd-pre-release`, etc.).
2. List the target as a prerequisite of `test-all-workflows`.

Nothing else catches a missing test target — a new workflow file alone is invisible
to the lint job until it's wired in.

See [`.github/README.md`](.github/README.md) for the broader workflow conventions
(test-run model, env-block discipline, fan-in patterns, etc.).
