include config/deploy/.env

webhook_info:
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/getWebhookInfo"

webhook_delete:
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/deleteWebhook"

webhook_create:
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/setWebhook" --header 'content-type: application/json' --data "{\"url\": \"$(SERVERLESS_APIGW_URL)\"}"

service_account_key_create:
	yc iam key create --service-account-name $(SERVICE_ACCOUNT_YDB_NAME) --output key.json

create_iam_token:
	yc iam create-token

build:
	docker build --platform=linux/amd64 --pull --rm -t cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME) .

push: webhook_create
	docker push cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME)

run_local:
	export $$(cat config/dev/.env | xargs -d '\n'); go run cmd/local/main.go