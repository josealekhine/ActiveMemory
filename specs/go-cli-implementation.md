# Go CLI Implementation Specification

## Overview

The `amem` CLI is implemented in Go for single-binary distribution with zero runtime dependencies.

## Project Structure

```
active-memory/
├── cmd/
│   └── amem/
│       └── main.go           # Entry point
├── internal/
│   ├── context/
│   │   ├── loader.go         # Context loading and assembly
│   │   ├── parser.go         # Markdown parsing
│   │   └── token.go          # Token estimation
│   ├── files/
│   │   ├── constitution.go   # CONSTITUTION.md handler
│   │   ├── tasks.go          # TASKS.md handler
│   │   ├── decisions.go      # DECISIONS.md handler
│   │   ├── learnings.go      # LEARNINGS.md handler
│   │   ├── conventions.go    # CONVENTIONS.md handler
│   │   ├── architecture.go   # ARCHITECTURE.md handler
│   │   ├── glossary.go       # GLOSSARY.md handler
│   │   └── drift.go          # DRIFT.md handler
│   ├── drift/
│   │   ├── detector.go       # Drift detection logic
│   │   ├── paths.go          # Path existence checks
│   │   └── rules.go          # Constitution rule checks
│   ├── cli/
│   │   ├── init.go           # amem init
│   │   ├── status.go         # amem status
│   │   ├── load.go           # amem load
│   │   ├── sync.go           # amem sync
│   │   ├── compact.go        # amem compact
│   │   ├── drift.go          # amem drift
│   │   ├── agent.go          # amem agent
│   │   ├── add.go            # amem add
│   │   ├── complete.go       # amem complete
│   │   ├── watch.go          # amem watch
│   │   └── hook.go           # amem hook
│   └── templates/
│       └── embed.go          # Embedded template files
├── templates/
│   ├── CONSTITUTION.md
│   ├── TASKS.md
│   ├── DECISIONS.md
│   ├── LEARNINGS.md
│   ├── CONVENTIONS.md
│   ├── ARCHITECTURE.md
│   ├── GLOSSARY.md
│   ├── DRIFT.md
│   └── AGENT_PLAYBOOK.md
├── scripts/
│   └── build-all.sh          # Cross-platform build script
├── examples/
│   └── demo/                 # Example project with .context/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Dependencies

Minimal dependencies (standard library preferred):

```go
// go.mod
module github.com/josealekhine/ActiveMemory

go 1.22

require (
    github.com/spf13/cobra v1.8.0      // CLI framework
    github.com/fatih/color v1.16.0     // Colored output
    gopkg.in/yaml.v3 v3.0.1            // YAML parsing (for tasks.yaml option)
)
```

## Build Process

### Local Development

```bash
# Build for current platform
go build -o amem ./cmd/amem

# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Lint
golangci-lint run
```

### Release Build

```bash
# scripts/build-all.sh
#!/bin/bash
set -e

VERSION=${1:-"dev"}
LDFLAGS="-s -w -X main.Version=$VERSION"

platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output="dist/amem-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        output+=".exe"
    fi
    
    echo "Building $output..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "$output" ./cmd/amem
done

echo "Done. Binaries in dist/"
```

### GitHub Actions Release

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Build binaries
        run: ./scripts/build-all.sh ${{ github.ref_name }}
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
```

## Embedded Templates

Use Go's `embed` package for template files:

```go
// internal/templates/embed.go
package templates

import "embed"

//go:embed *.md
var FS embed.FS

func GetTemplate(name string) ([]byte, error) {
    return FS.ReadFile(name)
}
```

## CLI Framework

Using Cobra for command structure:

```go
// cmd/amem/main.go
package main

import (
    "os"
    "github.com/spf13/cobra"
    "github.com/josealekhine/ActiveMemory/internal/cli"
)

var Version = "dev"

func main() {
    root := &cobra.Command{
        Use:     "amem",
        Short:   "Active Memory - persistent context for AI coding assistants",
        Version: Version,
    }
    
    root.AddCommand(
        cli.InitCmd(),
        cli.StatusCmd(),
        cli.LoadCmd(),
        cli.SyncCmd(),
        cli.CompactCmd(),
        cli.DriftCmd(),
        cli.AgentCmd(),
        cli.AddCmd(),
        cli.CompleteCmd(),
        cli.WatchCmd(),
        cli.HookCmd(),
    )
    
    if err := root.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Context not found |
| 3 | Invalid arguments |
| 4 | File operation error |
| 5 | Drift detected (for CI usage) |

## Testing

```go
// internal/context/loader_test.go
func TestLoadContext(t *testing.T) {
    // Create temp directory with .context/
    dir := t.TempDir()
    contextDir := filepath.Join(dir, ".context")
    os.MkdirAll(contextDir, 0755)
    
    // Write test files
    os.WriteFile(
        filepath.Join(contextDir, "TASKS.md"),
        []byte("# Tasks\n\n- [ ] Test task"),
        0644,
    )
    
    // Test loading
    ctx, err := LoadContext(dir)
    require.NoError(t, err)
    assert.Len(t, ctx.Files, 1)
}
```

## Performance Targets

- `amem status`: < 50ms
- `amem load`: < 100ms for typical project
- `amem drift`: < 500ms (filesystem scanning)
- `amem agent`: < 100ms
- Binary size: < 10MB

## Versioning

Semantic versioning with git tags:

- `v0.1.0` — Initial release
- `v0.2.0` — New features
- `v1.0.0` — Stable API

Version embedded at build time via `-ldflags`.
