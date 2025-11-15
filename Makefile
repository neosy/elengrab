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

help: ## Список команд
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "Usage: make <commands> \033[36m\033[0m\n" \
	} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { \
		printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	} ' $(MAKEFILE_LIST)

server-run: ## Запуск fastHTTP сервера
	@echo "***** SERVER RUN *****"
	@set -o allexport; \
	. ./cmd/${APP_NAME}/.env; \
	go run ./cmd/${APP_NAME}/main.go

build: update-app-version ## Билд исполняемого файла
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(APP_NAME) ./cmd/$(APP_NAME)/main.go

img-build: update-app-version ## Генерация образа docker контейнера
	docker build -t $(APP_IMG_NAME) .

img-rebuild: update-app-version ## Удаление и генерация образа docker контейнера
	docker rmi -f $(APP_IMG_NAME)
	docker build -t $(APP_IMG_NAME) .

img-rebuild-push: img-rebuild img-push ## Сборка images, обновление в репозитарии и очистка

img-rm: ## Удаление image с тегом latest
	-docker rmi -f $(APP_IMG_NAME)
	
img-push: ## Отправка images в репозитарий с тегом latest
	docker tag $(APP_IMG_NAME) $(APP_IMG_LATEST)
	docker push $(APP_IMG_LATEST)
	
img-push-version: ## Отправка images в репозитарий с тегом актуальной версии
	docker tag $(APP_IMG_NAME) $(APP_IMG):$(VERSION)
	docker push $(APP_IMG):$(VERSION)
	docker rmi $(APP_IMG):$(VERSION)

img-pull: ## Загрузка images из репозитария
	@docker pull $(APP_IMG_LATEST)

docker-run: ## Запуск докера
	docker run -d \
		--name $(APP_NAME) \
		-p $(HTTP_PORT):8080 \
		-v $(APP_NAME)_db:/app_n/sqlite/data \
		-v $(APP_NAME)_downloads:/app_n/downloads \
		$(APP_IMG_NAME_LATEST)

git-push-tag-version: ## Создание тега в git для актуальной версии
	-git tag v$(VERSION)
	git push --tags

git-push-update-tag-version: ## Обновление тега в git для актуальной версии
	git tag -f v$(VERSION)
	git push origin -f v$(VERSION)

update-app-version: ## Update AppVersion in Go
	@sed -i "s|AppVersion = \".*\"|AppVersion = \"${VERSION}\"|" ./infrastructure/constants/app_version.go

version-create: ## Создание файла с номер версии программы
	echo -n $(VERSION_START) > $(VERSION_FILE)
	
version-inc: ## Увеличение номера версии программы и сохранение в файл
	@VERSION_NEW=$$(echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}'); \
	echo -n $$VERSION_NEW > $(VERSION_FILE)
	@$(MAKE) --no-print-directory update-app-version

version-img-list: ## Список версий images
	curl -s $(DOCKER_HTTP_ADRR_TAG_LIST) | jq .

stack-deploy: ## Развертывание контейнеров
	@docker stack deploy -c docker-compose.yml --detach=true $(STACK_NAME)

stack-rm: ## Удаление контейнеров
	@docker stack rm $(STACK_NAME)
