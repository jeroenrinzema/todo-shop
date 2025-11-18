SHELL := /bin/bash
.DEFAULT_GOAL := generate

BIN			= $(CURDIR)/bin
BUILD_DIR	= $(CURDIR)/build

GOPATH			= $(HOME)/go
GOBIN			= $(GOPATH)/bin
GO				?= GOGC=off $(shell which go)
NODE			?= $(shell which node)
PNPM			?= $(shell which pnpm)
PKGS			= $(or $(PKG),$(shell env $(GO) list ./...))
VERSION			?= $(shell git describe --tags --always --match=v*)
SHORT_COMMIT	?= $(shell git rev-parse --short HEAD)

PATH := $(GOBIN):$(BIN):$(PATH)

# Printing
V ?= 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

$(BUILD_DIR):
	@mkdir -p $@

# Tools
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(@F)…)
	$Q GOBIN=$(BIN) $(GO) install $(shell $(GO) list tool | grep $(@F))

OAPI_CODEGEN = $(BIN)/oapi-codegen
TOOLCHAIN = $(OAPI_CODEGEN)

# Targets
.PHONY: lint
lint: | $(EMBEDDED) $(GOLANGCI_LINT) ; $(info $(M) running linters…) @ ## Run the project linters
	$Q $(GOLANGCI_LINT) run --max-issues-per-linter 10 --timeout 5m

.PHONY: test
test: | $(EMBEDDED) ; $(info $(M) running tests) @ ## Run all tests
	$Q $(GO) test $(PKGS) -timeout 30s -race -count 1

.PHONY: fmt
fmt: | $(EMBEDDED) ; $(info $(M) running go fmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)

.PHONY: generate
generate: | $(EMBEDDED) $(TOOLCHAIN) ; $(info $(M) running go generate…) @ ## Run gogenerate on all source files
	$Q $(GO) generate $(PKGS)
	$Q $(MAKE) fmt

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf $(BUILD)
	@rm -rf $(BUILD_DIR)
	@find . -name '*_mock_test.go' -exec rm -r {} \;
	@find . -name '*_string.go' -exec rm -r {} \;
	@find . -name '*_gen.go' -exec rm -r {} \;
	@find . -name '*.pb.go' -exec rm -r {} \;

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
