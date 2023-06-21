GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
GOPRIVATE="github.com/Allen-Career-Institute/*"
GOVERSION := $(shell go version | cut -d " " -f 3 | cut -c 3-)
GOROOT:=$(shell go env GOROOT)
PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name '*.proto'")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name '*.proto'")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name '*.proto')
	API_PROTO_FILES=$(shell find api -name '*.proto')
endif

.PHONY: init
init:
	@echo ''
	@echo 'Init:'
	@echo ' make [init]'
	@echo ''
	@echo 'Installing dependencies:'
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2




lint:
	golangci-lint run

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out
	go test -v ./... -covermode=count -coverprofile=sonar_coverage.out

.PHONY: config
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
 	       --validate_out=paths=source_relative,lang=go:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: error
# generate error proto
error:
	protoc --proto_path=./api \
             --proto_path=./third_party \
             --go_out=paths=source_relative:./api \
             --go-errors_out=paths=source_relative:./api \
             $(API_PROTO_FILES)

.PHONY: build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: test
test:
	go test -short -coverprofile=bin/cov.out `go list ./... | grep -v vendor/`
	go tool cover -func=bin/cov.out

.PHONY: generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: all
all:
	make api;
	make config;
	make generate;
	make error;
	make lint;
	make build;


help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
