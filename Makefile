# Makefile
# fastHTTP server "npulse-watcher"

.DEFAULT_GOAL := help

include .make.env
export

VERSION_FILE := "VERSION"
VERSION_START := "0.1.0"

VERSION := $(shell cat $(VERSION_FILE))
# VERSION_NEW := $(shell echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}')

.DEFAULT_GOAL := help

help: ## List of commands
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "Usage: make <commands> \033[36m\033[0m\n" \
	} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { \
		printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	} ' $(MAKEFILE_LIST)

server-run: ## FastHTTP server startup
	@echo "***** SERVER RUN *****"
	@set -o allexport; \
	. ./cmd/${APP_NAME}/.env; \
	go run ./cmd/${APP_NAME}/main.go

build: update-app-version ## Build executable file
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(APP_NAME) ./cmd/$(APP_NAME)/

build-embedded: update-app-version ## Build executable file with embedded assets
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags embed_assets -v -o $(APP_NAME) ./cmd/$(APP_NAME)/

build-win-embedded: update-app-version ## Build executable file with embedded assets
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags embed_assets -v -o $(APP_NAME).exe ./cmd/$(APP_NAME)/

img-build: update-app-version ## Build Docker image. Arguments: INSTALL_DENO=true
	docker build \
		--build-arg INSTALL_DENO=$(INSTALL_DENO) \
		-t $(APP_IMG_LATEST)$(if $(filter true,$(INSTALL_DENO)),-deno) .

img-rebuild: update-app-version ## Remove and rebuild Docker image. Arguments: INSTALL_DENO=true
	docker rmi -f $(APP_IMG_LATEST)$(if $(filter true,$(INSTALL_DENO)),-deno)
	docker build \
		--build-arg INSTALL_DENO=$(INSTALL_DENO) \
		-t $(APP_IMG_LATEST)$(if $(filter true,$(INSTALL_DENO)),-deno) .

img-rebuild-no-cache: update-app-version ## Remove and rebuild Docker image without cache. Arguments: INSTALL_DENO=true
	docker rmi -f $(APP_IMG_LATEST)$(if $(filter true,$(INSTALL_DENO)),-deno)
	docker build --no-cache \
		--build-arg INSTALL_DENO=$(INSTALL_DENO) \
		-t $(APP_IMG_LATEST)$(if $(filter true,$(INSTALL_DENO)),-deno) .

img-rebuild-push: img-rebuild img-push ## Build images, push to repository, and clean up

img-rm: ## Remove image with the latest tag
	-docker rmi -f $(APP_IMG)
	
img-push: ## Push images to the repository with the latest tag. Arguments: SUFFIX=deno
	docker tag $(APP_IMG_LATEST)$(if $(SUFFIX),-$(SUFFIX)) $(APP_IMG_LATEST)$(if $(SUFFIX),-$(SUFFIX))
	docker push $(APP_IMG_LATEST)$(if $(SUFFIX),-$(SUFFIX))
	
img-push-version: ## Push images to the repository with the current version tag. Arguments: SUFFIX=deno
	docker tag $(APP_IMG_LATEST)$(if $(SUFFIX),-$(SUFFIX)) $(APP_IMG):$(VERSION)$(if $(SUFFIX),-$(SUFFIX))
	docker push $(APP_IMG):$(VERSION)$(if $(SUFFIX),-$(SUFFIX))
	docker rmi $(APP_IMG):$(VERSION)$(if $(SUFFIX),-$(SUFFIX))

img-pull: ## Pull images from the repository
	@docker pull $(APP_IMG_LATEST)

docker-run: ## Run Docker container
	docker run -d \
		--name $(APP_NAME) \
		-p $(HTTP_PORT):8080 \
		-v $(APP_NAME)_db:/app_n/sqlite/data \
		-v $(APP_NAME)_downloads:/app_n/downloads \
		$(APP_IMG_LATEST)

git-push-tag-version: ## Create a Git tag for the current version
	-git tag v$(VERSION)
	git push --tags

git-push-update-tag-version: ## Update the Git tag for the current version
	git tag -f v$(VERSION)
	git push origin -f v$(VERSION)

update-app-version: ## Update AppVersion in Go
	@sed -i "s|AppVersion = \".*\"|AppVersion = \"${VERSION}\"|" ./infrastructure/config/consts.go

version-create: ## Create a file with the application version number
	echo -n $(VERSION_START) > $(VERSION_FILE)
	
version-inc: ## Increment application version and save it to file
# 	@VERSION_NEW=$$(echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}'); \
# 	echo -n $$VERSION_NEW > $(VERSION_FILE)
	@VERSION_NEW=$$(echo $(VERSION) | awk '{match($$0,/[0-9]+$$/); printf "%s%d\n", substr($$0,1,RSTART-1), substr($$0,RSTART,RLENGTH)+1}'); \
	echo -n $$VERSION_NEW > $(VERSION_FILE)
	@$(MAKE) --no-print-directory update-app-version

version-img-list: ## List image versions
	curl -s $(DOCKER_HTTP_ADRR_TAG_LIST) | jq .

stack-deploy: ## Deploy containers
	@docker stack deploy -c docker-compose.yml --detach=true $(STACK_NAME)

stack-rm: ## Remove containers
	@docker stack rm $(STACK_NAME)
