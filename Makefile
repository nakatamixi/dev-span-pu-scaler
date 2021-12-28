GOOGLE_CLOUD_PROJECT ?=
CURRENT_SHORTHASH := $(shell git rev-parse --short HEAD)
GCR ?= gcr.io
IMAGE_PATH ?= dev-span-pu-scalar
IMAGE_NAME := $(GCR)/$(GOOGLE_CLOUD_PROJECT)/$(IMAGE_PATH)
SERVICE_ACCOUNT ?= dev-span-pu-scaler@$(GOOGLE_CLOUD_PROJECT).iam.gserviceaccount.com
# specify commma separated target instances
INSTANCES ?=
REGION ?= asia-northeast1
.PHONY: build-tools
build-tools:
	go install github.com/google/ko@v0.8.3

.PHONY: ko-publish
ko-publish:
	GOROOT=$(shell go env GOROOT) SOURCE_DATE_EPOCH=$(shell date +%s) KO_DOCKER_REPO=$(IMAGE_NAME) ko publish --bare ./cmd/ --tags=$(CURRENT_SHORTHASH)

.PHONY: deploy
deploy:
	gcloud beta run deploy dev-span-pu-scaler  \
		--image $(IMAGE_NAME):$(CURRENT_SHORTHASH) \
		--platform=managed --region=$(REGION)  \
		--args ^~^-server~-instances=$(INSTANCES) \
		--port=8080 \
		--service-account $(SERVICE_ACCOUNT) \
		--quiet
