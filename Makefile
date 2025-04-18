include .env

.PHONY: parse_timetable
parse_timetable:
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
