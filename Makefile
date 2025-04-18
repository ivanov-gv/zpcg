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
