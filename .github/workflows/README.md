# Workflows

Callable workflows:

- checks.yml - test, lint, build checks for PRs and 'main' branch
- stage.yml - builds and deploys to stage
- release.yml - builds and deploys to stage, production, and creates a release

Reusable workflows:

- build.yml - builds and optionally pushed docker images
- deploy.yml - deploys a docker image to GCloud

# build.yml

Builds a docker image and tags it for 2 registries: GitHub's ghcr.io and GCloud's docker.pkg.dev.

- ghcr.io - a backup registry
- pkg.dev - the primary one

Build and push logic:

1. Check if `ghcr.io` has the image built from the requested commit
2. If exists - OK, return
3. If not - build, push to BOTH ghcr.io AND pkg.dev

Result: no unnecessary rebuilds, 2 up-to-date registries 
