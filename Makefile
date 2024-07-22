.DEFAULT_GOAL := test
SHELL := /bin/bash

# environment variables
BIN_DIR := $(CURDIR)/bin
GOBIN := $(BIN_DIR)
PATH := $(abspath $(BIN_DIR)):$(PATH)
export

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

PROTOC := $(BIN_DIR)/protoc
PROTOC_VERSION := 3.20.0
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-$(UNAME_OS)-$(UNAME_ARCH).zip
ifeq "$(UNAME_OS)" "Darwin"
	PROTOC_ZIP=protoc-$(PROTOC_VERSION)-osx-$(UNAME_ARCH).zip
endif
$(PROTOC): | $(BIN_DIR)
	@curl -sSOL \
		"https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)"
	@unzip -j -o $(PROTOC_ZIP) -d $(BIN_DIR) bin/protoc
	@unzip -o $(PROTOC_ZIP) -d $(BIN_DIR) "include/*"
	@rm -f $(PROTOC_ZIP)

PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
$(PROTOC_GEN_GO): | $(BIN_DIR)
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26

PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc
$(PROTOC_GEN_GO_GRPC): | $(BIN_DIR)
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

.PHONY: test
test: ## run all tests
	go test ./... -race

.PHONY: build
build: build/parser-server build/parser ## build binaries

.PHONY: build/%
build/%:
	@CGO_ENABLED=0 go build -v -o $(BIN_DIR)/$* ./cmd/$*

export SERVER_PORT ?= 8888

PROTO_DIR := $(CURDIR)/proto
GEN_PB_DIR := $(CURDIR)/gen/pb
PROTOC_OPT := -I$(PROTO_DIR)
PROTOC_GO_OPT := --plugin=${BIN_DIR}/protoc-gen-go --go_out=$(GEN_PB_DIR) --go_opt=paths=source_relative
PROTOC_GO_GRPC_OPT := --go-grpc_out=require_unimplemented_servers=false:$(GEN_PB_DIR) --go-grpc_opt=paths=source_relative
.PHONY: generate
generate: $(PROTOC) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) ## generate Go code from proto files
	@rm -rf $(GEN_PB_DIR)
	@mkdir -p $(GEN_PB_DIR)
	@find $(PROTO_DIR) -name '*.proto' | xargs -P8 protoc $(PROTOC_OPT) $(PROTOC_GO_OPT) $(PROTOC_GO_GRPC_OPT)

.PHONY: clean
clean: ## remove binaries, tools and generated files
	@rm -rf $(BIN_DIR) $(GEN_PB_DIR)

.PHONY: help
help: ## print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'