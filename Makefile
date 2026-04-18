include .env

# test

.PHONY: test_unit
test_unit:
	go test -count=1 -race ./internal/...

.PHONY: test_integration
test_integration: # needs tdlib installed. see: https://tdlib.github.io/td/build.html
	go test -count=1 -race ./test/...

.PHONY: test
test: # needs tdlib installed. see: https://tdlib.github.io/td/build.html
	go test -count=1 -race ./...

# deploy

.PHONY: golines_install
golines_install:
	go install github.com/segmentio/golines@latest

.PHONY: parse_timetable
parse_timetable: # needs golines. see target 'golines_install' or use `go install github.com/segmentio/golines@latest`
	go run ./cmd/exporter -file=gen/timetable/timetable.gen.go
	golines ./gen/timetable/timetable.gen.go --no-ignore-generated -w

.PHONY: build
build:
	docker build -t ${DOCKER_IMAGE_TAG} -f deploy/Dockerfile .

.PHONY: push
push:
	docker push ${DOCKER_IMAGE_TAG}

.PHONY: deploy
deploy:
	gcloud run deploy ${GCLOUD_PROJECT} \
    --image=${DOCKER_IMAGE_TAG} \
    --set-env-vars=ENVIRONMENT=${GCLOUD_ENV_VAR_ENVIRONMENT} \
    --set-secrets=TELEGRAM_APITOKEN=${GCLOUD_SECRET_TELEGRAM_APITOKEN} \
    --execution-environment=gen1 \
    --region=${GCLOUD_REGION}\
    --project=${GCLOUD_PROJECT_ID} \
     && gcloud run services update-traffic ${GCLOUD_PROJECT} --to-latest \
     --region ${GCLOUD_REGION}

.PHONY: new_version
new_version: parse_timetable build push deploy

.PHONY: new_version_without_timetable_update
new_version_without_timetable_update: build push deploy

# telegram api

.PHONY: add_webhook
add_webhook:
	curl https://api.telegram.org/bot${TG_TOKEN}/setWebhook?url=${TG_WEBHOOK}

.PHONY: delete_webhook
delete_webhook:
	curl https://api.telegram.org/bot${TG_TOKEN}/setWebhook?url=

.PHONY: get_webhook
get_webhook:
	curl https://api.telegram.org/bot${TG_TOKEN}/getWebhookInfo

# add info

.PHONY: add_info
add_info:
	TELEGRAM_APITOKEN=${TG_TOKEN} go run ./cmd/tg-init

# check

.PHONY: get_default_commands
get_default_commands:
	curl -X GET https://api.telegram.org/bot${TG_TOKEN}/getMyCommands

.PHONY: get_en_commands
get_en_commands:
	curl -X GET \
  -H 'Content-Type: application/json' \
  -d '{"language_code":"en"}' \
  https://api.telegram.org/bot${TG_TOKEN}/getMyCommands

# ci testing with act
# Install act standalone:    https://github.com/nektos/act#installation
# Install act as gh extension: gh extension install nektos/gh-act
# Override the executable:   ACT='gh act' make test-ci-checks
ACT ?= act

.PHONY: test-ci-checks
test-ci-checks: # run checks workflow locally via act; dry_run=true skips test/lint (no container pull needed). copy .github/act/checks.env.example → checks.env first
	$(ACT) workflow_dispatch \
		-W .github/workflows/checks.yml \
		--secret-file .github/act/checks.env \
		--var-file .github/act/checks.env \
		--input dry_run=true

.PHONY: test-ci-ci-image
test-ci-ci-image: # run ci-image workflow locally via act (dry-run: build and push are skipped). copy .github/act/ci-image.env.example → ci-image.env first
	$(ACT) workflow_dispatch \
		-W .github/workflows/ci-image.yml \
		--secret-file .github/act/ci-image.env \
		--var-file .github/act/ci-image.env \
		--input dry_run=true

.PHONY: test-ci-deploy-to-preprod
test-ci-deploy-to-preprod: # run deploy-to-preprod workflow locally via act (dry-run: no registry writes or deployments). copy .github/act/deploy-to-preprod.env.example → deploy-to-preprod.env first
	$(ACT) workflow_dispatch \
		-W .github/workflows/deploy-to-preprod.yml \
		--secret-file .github/act/deploy-to-preprod.env \
		--var-file .github/act/deploy-to-preprod.env \
		--input dry_run=true

.PHONY: test-ci-release
test-ci-release: # run release workflow locally via act (dry-run: no registry writes or deployments). copy .github/act/release.env.example → release.env first
	$(ACT) workflow_dispatch \
		-W .github/workflows/release.yml \
		--secret-file .github/act/release.env \
		--var-file .github/act/release.env \
		--input dry_run=true
