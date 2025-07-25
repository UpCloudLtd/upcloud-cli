# UpCloud CLI Makefile

GO       = go
PYTHON   = python3
PIP      = pip3
CLI      = upctl
MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
LATEST_RELEASE ?= $(shell git describe --tags --match=v* --abbrev=0 | grep -Eo "[0-9]+\.[0-9]+\.[0-9]+")
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
			'{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			$(PKGS))

BIN_DIR              = $(CURDIR)/bin
CLI_BIN              = $(CLI)
BIN_LINUX            = $(CLI_BIN)-$(VERSION)-linux-amd64
BIN_DARWIN           = $(CLI_BIN)-$(VERSION)-darwin-amd64
BIN_WINDOWS          = $(CLI_BIN)-$(VERSION)-windows-amd64.exe
BIN_FREEBSD          = $(CLI_BIN)-$(VERSION)-freebsd-amd64


V = 0
Q = $(if $(filter 1,$V),,@)

TOOLS_DIR:=$(CURDIR)/.ci/bin

export GO111MODULE=on

.PHONY: build
build: fmt | $(BIN_DIR) ; $(info building executable for the current target…) @ ## Build program binary for current os/arch
	$Q $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(CLI_BIN) cmd/$(CLI)/main.go

.PHONY: clean-md-docs
clean-md-docs:
	rm -f docs/changelog.md
	rm -rf docs/commands_reference/
	rm -rf docs/examples/

.PHONY: md-docs
md-docs: clean-md-docs ## Generate documentation (markdown)
	$(GO) run ./.ci/docs/
	cp CHANGELOG.md docs/changelog.md
	mkdir -p docs/examples/

.PHONY: clean-docs
clean-docs:
	rm -f mkdocs.yaml
	rm -rf site/

.PHONY: install-docs-tools
install-docs-tools:
	$(PIP) install -r requirements.txt
	cd .ci/tools && GOBIN=$(TOOLS_DIR) go install github.com/UpCloudLtd/mdtest

.PHONY: docs
docs: clean-docs md-docs install-docs-tools ## Generate documentation (mkdocs site)
	$(TOOLS_DIR)/mdtest normalise examples/ -o docs/examples/ -t filename=title --quote-values always
	$(PYTHON) .ci/docs/generate_dynamic_nav.py
	echo "latest_release: $(LATEST_RELEASE)" > vars.yaml
	mkdocs build

.PHONY: build-all
build-all: build-linux build-darwin build-windows build-freebsd ## Build all targets

.PHONY: build-linux
build-linux: ; $(info building executable for Linux x86_64…) @ ## Build program binary for linux x86_64
	$Q GOOS=linux GOARCH=amd64 $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(BIN_LINUX) cmd/$(CLI)/main.go

.PHONY: build-freebsd
build-freebsd: ; $(info building executable for FreeBSD x86_64…) @ ## Build program binary for freebsd x86_64
	$Q GOOS=freebsd GOARCH=amd64 $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(BIN_FREEBSD) cmd/$(CLI)/main.go

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

.PHONY: release-notes
release-notes: CHANGELOG_HEADER = ^\#\# \[
release-notes: CHANGELOG_VERSION = $(subst v,,$(VERSION))
release-notes:
	@awk \
		'/${CHANGELOG_HEADER}${CHANGELOG_VERSION}/ { flag = 1; next } \
		/${CHANGELOG_HEADER}/ { if ( flag ) { exit; } } \
		flag { if ( n ) { print prev; } n++; prev = $$0 }' \
		CHANGELOG.md
