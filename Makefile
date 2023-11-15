GO_VERSION ?= 1.18.8
GOOS ?= linux
GOARCH ?= amd64
GOPATH ?= $(shell go env GOPATH)
NODE_VERSION ?= 16.11.1


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

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server


.PHONY: build-web-deps
build-web-deps: ## 安装web依赖包
	cd web && pnpm install #--registry https://registry.npm.taobao.org

##@ Build

.PHONY: web
web: ## 编译构建web
	cd web && pnpm install && pnpm run build

.PHONY: backend
backend: ## 编译构建api
	@echo "Running ${@}"
	./scripts/release.sh build main.go

.PHONY: install
install: web backend ## 编译构建本项目,先构建web，再构建api
