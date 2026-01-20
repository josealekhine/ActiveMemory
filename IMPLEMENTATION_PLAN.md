# Implementation Plan

This file tracks the implementation progress for Active Memory CLI.

## Milestone 1: Project Scaffolding
- [x] Initialize Go module and directory structure
- [x] Create Cobra CLI skeleton in `cmd/amem/main.go`
- [ ] Create embedded templates in `internal/templates/`
- [ ] Add all template files to `templates/` directory

## Milestone 2: Core Commands (MVP)
- [ ] Implement `amem init` — Create `.context/` with template files
- [ ] Implement `amem status` — Show context summary with token estimate
- [ ] Implement `amem load` — Output assembled context markdown

## Milestone 3: Context Operations
- [ ] Implement `amem add` — Add decision/task/learning/convention
- [ ] Implement `amem complete` — Mark task as done
- [ ] Implement `amem agent` — Print AI-ready context packet

## Milestone 4: Maintenance Commands
- [ ] Implement `amem drift` — Detect stale paths, broken refs (text output)
- [ ] Implement `amem drift --json` — JSON output for automation
- [ ] Implement `amem sync` — Reconcile context with codebase
- [ ] Implement `amem compact` — Archive completed tasks
- [ ] Implement `amem watch` — Watch for context-update commands

## Milestone 5: Integration
- [ ] Implement `amem hook` — Generate AI tool integration configs
- [ ] Add `--help` text for all commands
- [ ] Add `--version` flag with build-time version

## Milestone 6: Testing & Release
- [ ] Write unit tests for `internal/context/` (loader, parser)
- [ ] Write unit tests for `internal/drift/` (detector)
- [ ] Write integration tests for CLI commands
- [ ] Create `scripts/build-all.sh` for cross-platform builds
- [ ] Create `.github/workflows/release.yml` for GitHub Actions
- [ ] Create `examples/demo/` with sample `.context/` directory
- [ ] Update README.md with installation and usage instructions

## Notes

- Build command: `CGO_ENABLED=0 go build -o amem ./cmd/amem`
- CGO is disabled due to gcc cross-compilation issues on ARM64 Linux
