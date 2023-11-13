GO_VERSION ?= 1.18.8
GOOS ?= linux
GOARCH ?= amd64
GOPATH ?= $(shell go env GOPATH)
NODE_VERSION ?= 16.11.1
COMPOSE_PROJECT_NAME := ${TAG}-$(shell git rev-parse --abbrev-ref HEAD)
BRANCH_NAME ?= $(shell git rev-parse --abbrev-ref HEAD | sed "s!/!-!g")
GIT_TAG := $(shell git rev-parse --short HEAD)
ifeq (${BRANCH_NAME},main)
	TAG    := ${GIT_TAG}-go${GO_VERSION}
	TRACKED_BRANCH := true
	LATEST_TAG := latest
else
	TAG := ${GIT_TAG}-${BRANCH_NAME}-go${GO_VERSION}
	ifneq (,$(findstring release-,$(BRANCH_NAME)))
		TRACKED_BRANCH := true
		LATEST_TAG := ${BRANCH_NAME}-latest
	endif
endif
CUSTOMTAG ?=


FILEEXT :=
ifeq (${GOOS},windows)
FILEEXT := .exe
endif


.DEFAULT_GOAL := help
.PHONY: help
help:
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n"\
	} \
	/^[a-zA-Z_-]+:.*?##/ { \
		printf "  \033[36m%-17s\033[0m %s\n", $$1, $$2 \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	} ' $(MAKEFILE_LIST)

##@ Dependencies

.PHONY: init
init:	## 初始化安装开发工具
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/cosmtrek/air@latest
	go install golang.org/x/tools/cmd/stringer

.PHONY: bootstrap
bootstrap:
	cd ./deploy/docker-compose && docker compose up -d && cd ../../
	go run ./cmd/migration
	nunu run ./cmd/server

.PHONY: mock
mock:
	mockgen -source=internal/service/user.go -destination test/mocks/service/user.go
	mockgen -source=internal/repository/user.go -destination test/mocks/repository/user.go
	mockgen -source=internal/repository/repository.go -destination test/mocks/repository/repository.go

.PHONY: test
test:
	go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./coverage.out ./test/server/...
	go tool cover -html=./coverage.out -o coverage.html

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

.PHONY: docker
docker:
	docker build -f deploy/build/Dockerfile --build-arg APP_RELATIVE_PATH=./cmd/job -t 1.1.1.1:5000/demo-job:v1 .
	docker run --rm -i 1.1.1.1:5000/demo-job:v1

.PHONY: build-web-deps
build-web-deps: ## 安装web依赖包
	cd web && npm ci --registry https://registry.npm.taobao.org

##@ Build

.PHONY: web
web: ## 编译构建web
	cd web && pnpm install --registry https://registry.npm.taobao.org && pnpm run build

.PHONY: backend
backend: ## 编译构建api
	@echo "Running ${@}"
	./scripts/release.sh build main.go

.PHONY: install
install: web backend ## 编译构建本项目,先构建web，再构建api
