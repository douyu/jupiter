PROJECT_NAME := "jupiter"
PKG := "github.com/douyu/jupiter"
PKG_LIST := $(shell go list ${PKG}/... | grep /pkg/)
GO_FILES := $(shell find . -name '*.go' | grep /pkg/ | grep -v _test.go)

.DEFAULT_GOAL := default
.PHONY: all test lint fmt fmtcheck cmt errcheck license

GOCMT := $(shell command -v gocmt 2 > /dev/null)
REVIVE := $(shell command -v revive 2 > /dev/null)
ERRCHECK := $(shell command -v errcheck 2 > /dev/null)
all: fmt errcheck lint build

fmt: ## Format the files
	@gofmt -l -w $(GO_FILES)

fmtcheck: ## Check and format the files
	@gofmt -l -s $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed for:"; gofmt -l -s $(GO_FILES); fi

########################################################
lint:  ## lint check
	@hash revive 2>&- || go get -u github.com/mgechev/revive
	@revive -formatter stylish pkg/...

########################################################
cmt: ## auto comment exported Function
	@hash gocmt 2>&- || go get -u github.com/Gnouc/gocmt
	@gocmt -d pkg -i

########################################################
errcheck: ## check error
	@hash errcheck 2>&- || go get -u github.com/kisielk/errcheck
	@errcheck pkg/...

########################################################
test: ## Run unittests
	@go test -short ${PKG_LIST}

########################################################
race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

########################################################
msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

########################################################
dep: ## Get the dependencies
	@go get -v -d ./...

########################################################
version: ## Print git revision info
	@echo $(expr substr $(git rev-parse HEAD) 1 8)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

default: help

demo: ## Build jupiter demo and Run it
	@APP_NAME=dev go run example/all/cmd/demo/main.go --config=example/all/config/config.toml --watch

demo.build: ## Build jupiter Demo
	@JUPITER_MODE=dev go build -ldflags  "-X github.com/douyu/jupiter/initialize.AppName=hello" example/all/cmd/demo/main.go

license: ## Add license header for all code files
	@find . -name \*.go -exec sh -c "if ! grep -q 'LICENSE' '{}'; then mv '{}' tmp && cp doc/LICENSEHEADER.txt '{}' && cat tmp >> '{}' && rm tmp; fi" \;
