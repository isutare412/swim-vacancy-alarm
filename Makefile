ifneq (,$(wildcard ./.env.local))
    include .env.local
	export $(shell sed 's/=.*//' ./.env.local)
endif

# images
TAG ?= latest
IMAGE = redshoore/swim-vacancy-alarm:$(TAG)

# docker hub
DOCKER_USER ?= <docker_hub_username>
DOCKER_PASSWORD ?= <docker_hub_secret>

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: build
build: ## Build docker image.
	docker build -f Dockerfile -t $(IMAGE) .

.PHONY: push
push: ## Push docker image.
	echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USER) --password-stdin
	docker push $(IMAGE)
