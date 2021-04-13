# UpCloud CLI Makefile

GO       = go
CLI      = upctl
DOC_GEN  = doc-gen
MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
			'{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			$(PKGS))

BIN_DIR         = $(CURDIR)/bin
CLI_BIN         = $(CLI)
DOC_GEN_BIN     = $(DOC_GEN)
BIN_LINUX       = $(CLI_BIN)-$(VERSION)-linux-amd64
BIN_DARWIN      = $(CLI_BIN)-$(VERSION)-darwin-amd64
BIN_WINDOWS     = $(CLI_BIN)-$(VERSION)-windows-amd64.exe


V = 0
Q = $(if $(filter 1,$V),,@)

export GO111MODULE=on

.PHONY: build
build: fmt | $(BIN_DIR) ; $(info building executable for the current target…) @ ## Build program binary for current os/arch
	$Q $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(CLI_BIN) cmd/$(CLI)/main.go

doc: $(BIN_DIR) ; $(info building documentation generator ) @ ## Build documentation (markdown)
	$Q $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(DOC_GEN_BIN) cmd/$(DOC_GEN)/main.go

.PHONY: build-all
build-all: build-linux build-darwin build-windows ## Build all targets

.PHONY: build-linux
build-linux: $(BIN_DIR) ; $(info building executable for Linux x86_64…) @ ## Build program binary for linux x86_64
	$Q GOOS=linux GOARCH=amd64 $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(BIN_LINUX) cmd/$(CLI)/main.go

.PHONY: build-darwin
build-darwin: $(BIN_DIR) ; $(info building executable for Darwin x86_64…) @ ## Build program binary for darwin x86_64
	$Q GOOS=darwin GOARCH=amd64 $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(BIN_DARWIN) cmd/$(CLI)/main.go

.PHONY: build-windows
build-windows: $(BIN_DIR) ; $(info building executable for Windows x86_64…) @ ## Build program binary for windows x86_64
	$Q GOOS=windows GOARCH=amd64 $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(BIN_WINDOWS) cmd/$(CLI)/main.go


# Tests

.PHONY: test
test: fmt; $(info running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test $(TESTPKGS)

.PHONY: fmt
fmt: ; $(info running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)

# Misc

$(BIN_DIR):
	@mkdir -p $@

.PHONY: clean
clean: ; $(info cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN_DIR)

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
