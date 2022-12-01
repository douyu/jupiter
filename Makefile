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

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

default: help

lint: ## Lint the go files
	golangci-lint run -v

lintmd: ## Lint markdown files
	markdownlint -c .github/markdown_lint_config.json website/docs README.md pkg

e2e-test: ## Run e2e test
	cd test/e2e \
		&& go mod tidy \
		&& ginkgo -r -race -cover -covermode=atomic -coverprofile=coverage.txt --randomize-suites --trace -coverpkg=github.com/douyu/jupiter/... .\
		&& cd -

covsh-e2e: ## Get the coverage of e2e test
	gocovsh --profile test/e2e/coverage.txt

unit-test: ## Run unit test
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

covsh-unit: ## Get the coverage of unit test
	gocovsh --profile coverage.txt
