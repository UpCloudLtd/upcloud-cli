ALL_FILES=./...

default: build

build:
	@echo 'building a new upctl client'
	@cd cmd && go build -o ../upctl

test:
	go test $(ALL_FILES)

fmt:
	go fmt $(ALL_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	go vet $(ALL_FILES)

all: build test fmt vet

.PHONY: build test fmt fmtcheck vet
