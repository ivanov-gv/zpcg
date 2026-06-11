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


TEST_RUN_IMAGE_REPO ?= ttl.sh
TEST_RUN_IMAGE_PATH = $(TEST_RUN_IMAGE_REPO)/ivanov-gv/zpcg
TEST_RUN_TAG ?= test-run-tag

.PHONY: push-test-run-image
push-test-run-image:
	docker pull hello-world
	docker tag hello-world $(TEST_RUN_IMAGE_PATH):$(TEST_RUN_TAG)
	docker push $(TEST_RUN_IMAGE_PATH):$(TEST_RUN_TAG)

.PHONY: push-test-run-image-v0.0.0
push-test-run-image-v0.0.0:
	docker pull hello-world
	docker tag hello-world $(TEST_RUN_IMAGE_PATH):v0.0.0
	docker push $(TEST_RUN_IMAGE_PATH):v0.0.0

.PHONY: test-all-workflows
test-all-workflows: test-ci-pr-checks test-ci-pr-checks-image-build test-ci test-cd-pre-release test-cd-release

# Override the executable:   ACT=act make test-ci-checks
ACT ?= gh act

GCLOUD_SCOPE := https://www.googleapis.com/auth/cloud-platform.read-only
# lazy evaluation for gcloud token with read-only scope
GCLOUD_TOKEN = $(eval GCLOUD_TOKEN := $(shell gcloud auth print-access-token --scopes='$(GCLOUD_SCOPE)'))$(GCLOUD_TOKEN)

.PHONY: test-ci-pr-checks
test-ci-pr-checks:
	$(ACT) pull_request \
		-W .github/workflows/ci-pr-checks.yml \
		--secret-file .github/act/secret.env \
		--var-file .github/act/var.env

.PHONY: test-ci-pr-checks-image-build
test-ci-pr-checks-image-build:
	$(ACT) workflow_dispatch \
		-W .github/workflows/ci-pr-checks-image-build.yml \
		--secret-file .github/act/secret.env \
		--var-file .github/act/var.env

.PHONY: test-ci
test-ci:
	$(ACT) push \
		-W .github/workflows/ci.yml \
		-e .github/act/event-push-tag.json \
		--secret-file .github/act/secret.env \
		--var-file .github/act/var.env

.PHONY: test-cd-pre-release
test-cd-pre-release: push-test-run-image-v0.0.0
	$(ACT) release \
		-W .github/workflows/cd-pre-release.yml \
		-e .github/act/event-release-prerelease.json \
		--secret-file .github/act/secret.env \
        --var-file .github/act/var.env \
        --env CLOUDSDK_AUTH_ACCESS_TOKEN=$(GCLOUD_TOKEN)

.PHONY: test-cd-release
test-cd-release: push-test-run-image-v0.0.0
	$(ACT) release \
		-W .github/workflows/cd-release.yml \
		-e .github/act/event-release-release.json \
		--secret-file .github/act/secret.env \
        --var-file .github/act/var.env \
        --env CLOUDSDK_AUTH_ACCESS_TOKEN=$(GCLOUD_TOKEN)
