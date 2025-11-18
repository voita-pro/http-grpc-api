ENV_FILE = .env
PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0
PLATFORMS=linux/amd64

ifneq ($(wildcard $(ENV_FILE)),)
	include $(ENV_FILE)
endif

install-deps:
	GOBIN=$(PROJECT_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(PROJECT_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(PROJECT_BIN) go install github.com/bufbuild/buf/cmd/buf@latest
	GOBIN=$(PROJECT_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest


migration-status:
	goose status -v

migration-add:
	goose create $(name) sql

migration-up:
	goose up -v

migration-down:
	goose  down -v