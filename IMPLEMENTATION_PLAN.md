# Implementation Plan

This file tracks the implementation progress for Active Memory CLI.

## Milestone 1: Project Scaffolding
- [x] Initialize Go module and directory structure
- [x] Create Cobra CLI skeleton in `cmd/ctx/main.go`
- [x] Create embedded templates in `internal/templates/`
- [x] Add all template files to `templates/` directory

## Milestone 2: Core Commands (MVP)
- [x] Implement `ctx init` — Create `.context/` with template files
- [x] Implement `ctx status` — Show context summary with token estimate
- [x] Implement `ctx load` — Output assembled context markdown

## Milestone 3: Context Operations
- [x] Implement `ctx add` — Add decision/task/learning/convention
- [x] Implement `ctx complete` — Mark task as done
- [x] Implement `ctx agent` — Print AI-ready context packet

## Milestone 4: Maintenance Commands
- [x] Implement `ctx drift` — Detect stale paths, broken refs (text output)
- [x] Implement `ctx drift --json` — JSON output for automation
- [x] Implement `ctx sync` — Reconcile context with codebase
- [x] Implement `ctx compact` — Archive completed tasks
- [x] Implement `ctx watch` — Watch for context-update commands

## Milestone 5: Integration
- [x] Implement `ctx hook` — Generate AI tool integration configs
- [x] Add `--help` text for all commands
- [x] Add `--version` flag with build-time version

## Milestone 6: Testing & Release
- [x] Write unit tests for `internal/context/` (loader, parser)
- [x] Write unit tests for `internal/drift/` (detector)
- [x] Write integration tests for CLI commands
- [x] Create `scripts/build-all.sh` for cross-platform builds
- [x] Create `.github/workflows/release.yml` for GitHub Actions
- [x] Create `examples/demo/` with sample `.context/` directory
- [x] Update README.md with installation and usage instructions

## Notes

- Build command: `CGO_ENABLED=0 go build -o ctx ./cmd/ctx`
- CGO is disabled due to gcc cross-compilation issues on ARM64 Linux
