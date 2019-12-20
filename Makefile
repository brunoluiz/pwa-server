PROJECT_NAME:=go-pwa-server

## Prepare version tags
GIT_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)
VERSION:=$(shell git rev-parse HEAD)
DOCKER_TAG:=$(subst :,-,$(subst /,-,$(GIT_BRANCH)))
ifeq ($(GIT_BRANCH), master)
	DOCKER_TAG := latest
endif

.PHONY: test
test:
	go test -race ./...

build:
	go build -o ./bin/go-pwa-server ./cmd

# Docker build
docker-login:
	@docker login --username $(DOCKER_HUB_LOGIN) --password=$(DOCKER_HUB_PASSWORD)

docker-build:
	docker build -t $(PROJECT_NAME):local .
	docker push go-pwa-server

docker-run:
	docker run -p 80:80 \
		--env-file .env.sample \
		-v $(PWD)/test/static:/static \
		$(PROJECT_NAME)

ci-docker-build: docker-login
	docker build -t $(PROJECT_NAME):$(VERSION) -t $(PROJECT_NAME):$(DOCKER_TAG) $(PROJECT_NAME) .
	docker push go-pwa-server
