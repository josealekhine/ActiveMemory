# Context CLI Makefile
#
# Common targets for Go developers

.PHONY: build test vet fmt lint clean all release dogfood help test-coverage smoke

# Default binary name and output
BINARY := ctx
OUTPUT := $(BINARY)

# Default target
all: build

## build: Build for current platform
build:
	CGO_ENABLED=0 go build -o $(OUTPUT) ./cmd/ctx

## test: Run tests with coverage summary
test:
	@CGO_ENABLED=0 CTX_SKIP_PATH_CHECK=1 go test -cover ./...

## test-v: Run tests with verbose output
test-v:
	CGO_ENABLED=0 go test -v ./...

## test-cover: Run tests with coverage
test-cover:
	CGO_ENABLED=0 go test -cover ./...

## test-coverage: Run tests with coverage and check against target (70%)
test-coverage:
	@echo "Running coverage check (target: 70%)..."
	@echo ""
	@CGO_ENABLED=0 go test -cover ./internal/context ./internal/cli 2>&1 | tee /tmp/ctx-coverage.txt
	@echo ""
	@CONTEXT_COV=$$(grep 'internal/context' /tmp/ctx-coverage.txt | grep -oE '[0-9]+\.[0-9]+%' | sed 's/%//'); \
	CLI_COV=$$(grep 'internal/cli' /tmp/ctx-coverage.txt | grep -oE '[0-9]+\.[0-9]+%' | sed 's/%//'); \
	echo "Coverage summary:"; \
	echo "  internal/context: $${CONTEXT_COV}% (target: 70%)"; \
	echo "  internal/cli: $${CLI_COV}% (target: 70% - aspirational)"; \
	echo ""; \
	if [ $$(echo "$$CONTEXT_COV < 70" | bc -l) -eq 1 ]; then \
		echo "FAIL: internal/context coverage below 70%"; \
		rm -f /tmp/ctx-coverage.txt; \
		exit 1; \
	fi; \
	echo "Coverage check passed (internal/context >= 70%)"; \
	rm -f /tmp/ctx-coverage.txt

## smoke: Build and run basic commands to verify binary works
smoke: build
	@echo "Running smoke tests..."
	@TMPDIR=$$(mktemp -d) && \
	cd $$TMPDIR && \
	echo "  Testing: ctx --help" && \
	$(CURDIR)/$(BINARY) --help > /dev/null && \
	echo "  Testing: ctx init" && \
	CTX_SKIP_PATH_CHECK=1 $(CURDIR)/$(BINARY) init > /dev/null && \
	echo "  Testing: ctx status" && \
	$(CURDIR)/$(BINARY) status > /dev/null && \
	echo "  Testing: ctx agent" && \
	$(CURDIR)/$(BINARY) agent > /dev/null && \
	echo "  Testing: ctx drift" && \
	$(CURDIR)/$(BINARY) drift > /dev/null && \
	echo "  Testing: ctx add task 'smoke test task'" && \
	$(CURDIR)/$(BINARY) add task "smoke test task" > /dev/null && \
	echo "  Testing: ctx session save" && \
	$(CURDIR)/$(BINARY) session save > /dev/null && \
	rm -rf $$TMPDIR && \
	echo "" && \
	echo "Smoke tests passed!"

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
