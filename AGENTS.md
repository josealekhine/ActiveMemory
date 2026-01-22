# AGENTS.md — Context Operational Guide

## Core Principle

**You have NO conversational memory. Your memory IS the file system.**

Everything important must be written to files. Future iterations depend entirely on what you write now.

## Context Read Order

Load context in this order (most critical first):

1. `.context/CONSTITUTION.md` — NEVER violate these rules
2. `.context/TASKS.md` — What to work on
3. `.context/CONVENTIONS.md` — How to write code
4. `.context/ARCHITECTURE.md` — Where things go
5. `.context/GLOSSARY.md` — Correct terminology
6. `.context/DECISIONS.md` — Why things are the way they are
7. `.context/LEARNINGS.md` — Gotchas to avoid

## Build & Run

```bash
# Build the CLI (CGO disabled for cross-platform compatibility)
CGO_ENABLED=0 go build -o ctx ./cmd/ctx

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Lint (requires golangci-lint)
golangci-lint run

# Format code
go fmt ./...

# Build for all platforms (release)
./hack/build-all.sh
```

## Validation

Run these after implementing to get immediate feedback:

- **Tests**: `go test ./...` — All context operations must be tested
- **Lint**: `golangci-lint run` — Static analysis
- **Format**: `go fmt ./...` — Code formatting
- **Build**: `CGO_ENABLED=0 go build ./...` — Compilation check
- **Integration**: `go test -tags=integration ./...` — Test with real files

## Key Directories

```
cmd/
└── ctx/               # CLI entry point (main.go)

internal/
├── context/            # Core context engine (loading, parsing, merging)
├── files/              # Context file handlers (DECISIONS.md, TASKS.md, etc.)
├── drift/              # Drift detection logic
└── cli/                # CLI command implementations

pkg/
└── ctx/               # Public API (if any)

specs/                  # Feature specifications
examples/               # Example .context/ directories
hack/
└── build-all.sh        # Cross-platform build script
```

## Codebase Patterns

- **Context files are markdown** — Human-readable, token-efficient, git-friendly
- **Handlers are stateless** — Each file type has a parser and serializer, no side effects
- **CLI is thin** — Business logic lives in core/, CLI just wires things together
- **Tests simulate cold starts** — Every test begins with no loaded context

## Operational Notes

- Context files go in `.context/` at project root (configurable)
- All context files must be valid markdown with consistent structure
- File naming: `SCREAMING_SNAKE.md` for top-level, kebab-case for subdirectories
- Never store secrets or credentials in context files
