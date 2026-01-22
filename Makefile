# Context CLI Makefile
#
# Common targets for Go developers

.PHONY: build test vet fmt lint clean all release dogfood help

# Default binary name and output
BINARY := ctx
OUTPUT := $(BINARY)

# Default target
all: build

## build: Build for current platform
build:
	CGO_ENABLED=0 go build -o $(OUTPUT) ./cmd/ctx

## test: Run tests
test:
	CGO_ENABLED=0 CTX_SKIP_PATH_CHECK=1 go test ./...

## test-v: Run tests with verbose output
test-v:
	CGO_ENABLED=0 go test -v ./...

## test-cover: Run tests with coverage
test-cover:
	CGO_ENABLED=0 go test -cover ./...

## vet: Run go vet
vet:
	go vet ./...

## fmt: Format code
fmt:
	go fmt ./...

## lint: Run golangci-lint (requires golangci-lint installed)
lint:
	golangci-lint run

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/

## release: Build for all platforms
release:
	./hack/build-all.sh

## release-version: Build for all platforms with version
release-version:
	@test -n "$(VERSION)" || (echo "Usage: make release-version VERSION=1.0.0" && exit 1)
	./hack/build-all.sh $(VERSION)

## dogfood: Start dogfooding in a target folder
dogfood:
	@test -n "$(TARGET)" || (echo "Usage: make dogfood TARGET=~/WORKSPACE/ctx-dogfood" && exit 1)
	./hack/start-dogfood.sh $(TARGET)

## install: Install to /usr/local/bin
install: build
	sudo cp $(BINARY) /usr/local/bin/$(BINARY)

## help: Show this help
help:
	@echo "Context CLI - Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
