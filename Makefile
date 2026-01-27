# Context CLI Makefile
#
# Common targets for Go developers

.PHONY: build test vet fmt lint clean all release build-all dogfood help test-coverage smoke site site-serve site-setup

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

## release: Full release process (build, tag, push)
release:
	./hack/release.sh

## build-all: Build binaries for all platforms (no tag)
build-all:
	./hack/build-all.sh $$(cat VERSION | tr -d '[:space:]')

## release-notes: Generate release notes (use Claude Code slash command)
release-notes:
	@echo "To generate release notes, run in Claude Code:"
	@echo ""
	@echo "  /release-notes"
	@echo ""
	@echo "This will analyze commits since the last tag and write to dist/RELEASE_NOTES.md"

## dogfood: Start dogfooding in a target folder
dogfood:
	@test -n "$(TARGET)" || (echo "Usage: make dogfood TARGET=~/WORKSPACE/ctx-dogfood" && exit 1)
	./hack/start-dogfood.sh $(TARGET)

## install: Install to /usr/local/bin (run as: make build && sudo make install)
install:
	@test -f $(BINARY) || (echo "Binary not found. Run 'make build' first, then 'sudo make install'" && exit 1)
	cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "Installed ctx to /usr/local/bin/ctx"

## site-setup: Create venv and install zensical
site-setup:
	python3 -m venv .venv
	.venv/bin/pip install --upgrade pip
	.venv/bin/pip install zensical

## site: Build documentation site
site:
	.venv/bin/zensical build

## site-serve: Serve documentation site locally
site-serve:
	.venv/bin/zensical serve

## help: Show this help
help:
	@echo "Context CLI - Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
