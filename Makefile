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

default: help

# Lint the go files
golint:
	golangci-lint run -v

# Lint markdown files
lintmd:
	markdownlint -c .github/markdown_lint_config.json website/docs README.md pkg

# Run e2e test
e2e-test:
	cd test/e2e \
		&& go mod tidy \
		&& ginkgo -r -race -cover -covermode=atomic -coverprofile=coverage.txt --randomize-suites --trace -coverpkg=github.com/douyu/jupiter/... .\
		&& cd -

# Get the coverage of e2e test
covsh-e2e:
	gocovsh --profile test/e2e/coverage.txt

# Run unit test
unit-test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# Get the coverage of unit test
covsh-unit:
	gocovsh --profile coverage.txt

# install tools
init:
	go install github.com/bufbuild/buf/cmd/buf@v1.13.1
	go install github.com/srikrsna/protoc-gen-gotag@v0.6.2
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.4
	go install ./cmd/protoc-gen-go-echo

# update buf mod
update:
	cd api && buf mod update

.PHONY: generate
# generate code
generate:
	buf generate
	cd proto && buf generate --template buf.gen.tag.yaml

.PHONY: lint
# lint
lint:
	buf lint

# breaking
breaking:
	buf breaking --against https://github.com/douyu/proto/.git#branch=main,ref=HEAD~1,subdir=api

# test
test-proto2echo:
	go test -v -cover ./proto/...

# validate openapi docs
validate:
	swagger validate proto/helloworld/v1/helloworld.swagger.json

# serve openapi docs
serve:
	swagger serve proto/helloworld/v1/helloworld.swagger.json

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help