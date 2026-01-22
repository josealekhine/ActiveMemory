# Tasks — Context CLI Rebuild

This project is a dogfooding exercise: rebuild the `ctx` CLI from scratch using the specs provided.

## In Progress

## Next Up

### Phase 1: Project Scaffolding `#priority:high`
- [ ] Initialize Go module (`go mod init github.com/josealekhine/ctx`)
- [ ] Create directory structure (cmd/ctx, internal/cli, internal/context, internal/templates, etc.)
- [ ] Set up Cobra CLI skeleton in cmd/ctx/main.go
- [ ] Add dependencies (cobra, color, yaml)

### Phase 2: Core Commands `#priority:high`
- [ ] Implement `ctx init` — create .context/ with template files
- [ ] Implement `ctx status` — show context summary
- [ ] Implement `ctx agent` — print AI-ready context packet
- [ ] Implement `ctx load` — output assembled context

### Phase 3: Context Operations `#priority:medium`
- [ ] Implement `ctx add` — add decision/task/learning
- [ ] Implement `ctx complete` — mark task done
- [ ] Implement `ctx drift` — detect stale context
- [ ] Implement `ctx sync` — reconcile with codebase

### Phase 4: Maintenance Commands `#priority:medium`
- [ ] Implement `ctx compact` — archive old items
- [ ] Implement `ctx watch` — watch for update commands
- [ ] Implement `ctx hook` — generate tool config

### Phase 5: Session Management `#priority:medium`
- [ ] Implement `ctx session list` — list saved sessions
- [ ] Implement `ctx session save` — save current session
- [ ] Implement `ctx session load` — load previous session
- [ ] Implement `ctx session parse` — convert transcript to markdown

### Phase 6: Claude Code Integration `#priority:medium`
- [ ] Create CLAUDE.md template with ctx instructions
- [ ] Generate .claude/hooks/auto-save-session.sh
- [ ] Generate .claude/settings.local.json

### Phase 7: Testing & Polish `#priority:low`
- [ ] Add unit tests for context loading
- [ ] Add unit tests for drift detection
- [ ] Add integration tests for CLI commands
- [ ] Create cross-platform build script (hack/build-all.sh)

## Completed (Recent)

## Blocked

## Reference

**Specs to follow** (in `specs/` directory):
- `core-architecture.md` — Overall design philosophy
- `go-cli-implementation.md` — Go project structure and patterns
- `cli.md` — All CLI commands and their behavior
- `context-file-formats.md` — File format specifications
- `context-loader.md` — Loading and parsing logic
- `context-updater.md` — Update command handling

**Autonomy guidelines**:
- Follow specs for behavior, but you have freedom in implementation details
- Use idiomatic Go patterns
- Keep dependencies minimal
- Prioritize simplicity over cleverness
