# Tasks — Context CLI

### Phase 1: Project Scaffolding `#priority:high` `#area:setup`
- [x] Initialize Go module (`go mod init github.com/ActiveMemory/ctx`)
- [x] Create directory structure (cmd/ctx, internal/cli, internal/context, internal/templates)
- [x] Set up Cobra CLI skeleton in cmd/ctx/main.go
- [x] Add dependencies (cobra, color, yaml)

### Phase 2: Core Commands `#priority:high` `#area:cli`
- [x] Implement `ctx init` — create .context/ with template files
- [x] Implement `ctx status` — show context summary
- [x] Implement `ctx agent` — print AI-ready context packet
- [x] Implement `ctx load` — output assembled context

### Phase 3: Context Operations `#priority:high` `#area:cli`
- [x] Implement `ctx add` — add decision/task/learning
- [x] Implement `ctx complete` — mark task done
- [x] Implement `ctx drift` — detect stale context
- [x] Implement `ctx sync` — reconcile with codebase

### Phase 4: Maintenance Commands `#priority:medium` `#area:cli`
- [x] Implement `ctx compact` — archive old items
- [x] Implement `ctx watch` — watch for update commands
- [x] Implement `ctx watch --auto-save` mode
- [x] Implement `ctx hook` — generate tool config

### Phase 5: Session Management `#priority:medium` `#area:cli`
- [x] Implement `ctx session save` — manually dump context to sessions/
- [x] Implement `ctx session list` — list saved sessions with summaries
- [x] Implement `ctx session load <file>` — load/summarize a previous session
- [x] Implement `ctx session parse` — convert .jsonl transcript to readable markdown
- [x] Add `--extract` flag to session parse — extract decisions/learnings from transcript

### Phase 6: Claude Code Integration `#priority:high` `#area:integration`
- [x] Create `.context/sessions/` directory structure
- [x] Create CLAUDE.md for native Claude Code bootstrapping
- [x] Set up PreToolUse hook for auto-load
- [x] Set up SessionEnd hook for auto-save
- [x] Enhance `ctx init` to create Claude hooks (embedded scripts, settings.local.json)
- [x] Handle CLAUDE.md creation/merge in `ctx init` (backup, markers, --merge flag)
- [x] Add PATH check to `ctx init` — verify ctx is in PATH before creating hooks
- [x] Document session persistence in AGENT_PLAYBOOK.md

### Phase 7: Testing & Verification `#priority:high` `#area:quality`
- [x] Add headers to all files
- [x] Add integration tests — invoke actual binary, verify output
  - [x] `ctx init` creates expected files
  - [x] `ctx status` returns valid status (not just help text)
  - [x] `ctx add learning "test"` modifies LEARNINGS.md
  - [x] `ctx session save` creates session file
  - [x] `ctx agent` returns context packet
- [x] Set unit test coverage target (70% for internal/cli, internal/context)
- [x] Add coverage reporting to `make test`
- [x] Add smoke test to CI/Makefile: build binary, run basic commands
- [ ] Verify built binary executes subcommands (not silently falling through to root help)

### Phase 8: Task Archival & Snapshots `#priority:medium` `#area:cli`
- [ ] Implement `ctx tasks archive` — move completed tasks to timestamped archive file
- [ ] Implement `ctx tasks snapshot` — create point-in-time snapshot of TASKS.md
- [ ] Archive location: `.context/archive/tasks-YYYY-MM-DD.md`
- [ ] Keep Phase structure in archives for traceability
- [ ] Update CONSTITUTION.md: archival is allowed, deletion is not

### Phase 9: Claude Slash Commands (Skills) `#priority:medium` `#area:cli`
- [ ] Research how existing skills are registered (check ralph-loop pattern)
- [ ] Create `/ctx-save` skill — calls `ctx session save`
- [ ] Create `/ctx-status` skill — calls `ctx status`
- [ ] Create `/ctx-add-learning` skill — calls `ctx add learning`
- [ ] Create `/ctx-add-decision` skill — calls `ctx add decision`
- [ ] Create `/ctx-add-task` skill — calls `ctx add task`
- [ ] Create `/ctx-agent` skill — calls `ctx agent` (manual context load)
- [ ] Create `/ctx-archive` skill — calls `ctx tasks archive`
- [ ] Create `/ctx-loop` skill — calls `ctx loop` (generate Ralph loop script)
- [ ] Update `ctx init` to create skill definitions in `.claude/commands/`

### Phase 9b: Ralph Loop Integration `#priority:medium` `#area:cli`
- [ ] Implement `ctx loop` command — generate a ready-to-use loop.sh script
  - [ ] Detect AI tool in use (claude, aider, etc.) and generate appropriate invocation
  - [ ] Include configurable max iterations, prompt file path
  - [ ] Include completion signal detection (SYSTEM_CONVERGED, SYSTEM_BLOCKED)
  - [ ] Make script executable by default
- [ ] Add `ctx loop --prompt PROMPT.md` — specify custom prompt file
- [ ] Add `ctx loop --tool claude|aider|generic` — target specific AI CLI
- [ ] Document in README that `/ralph-loop` exists for Claude Code users

### Phase 10: Project Rename `#priority:medium` `#area:branding`
- [x] Rename project from "Active Memory" to "Context"
  - [x] Update README.md title and references
  - [x] Update Go module path (github.com/ActiveMemory/ctx)
  - [x] Update all import paths in Go files
  - [x] Update CLAUDE.md references
  - [x] Keep `ctx` as binary name (short for context)
- [ ] Handle GitHub repo rename (manual step)

### Phase 11: Documentation `#priority:low` `#area:docs`
- [x] Document Claude Code integration in README
- [x] Add "Dogfooding Guide" — how to use ctx on ctx itself
- [x] Document session auto-save setup for new users
- [x] Create actual documentation site in `docs/` folder
  - [x] Getting started guide
  - [x] CLI command reference
  - [x] Context file format reference
  - [x] Integration guides (Claude Code, Cursor, Aider, etc.)
  - [x] Ralph Loop pairing guide
- [ ] Set up Cloudflare Pages to serve docs at ctx.ist
- [ ] Review docs/ and README.md for accuracy and completeness `#human-in-the-loop`
  - Verify CLI examples work as documented
  - Check for inconsistencies between README.md and docs/
  - Requires human confirmation before marking complete
- [ ] Simplify `docs/index.md` to avoid README.md duplication `#blocked-by:ctx.ist-live`
  - Keep minimal intro + installation
  - Link to full docs at ctx.ist for details
  - Reduces drift between README.md and docs/

## Blocked

## Reference

**Specs** (in `specs/` directory):
- `core-architecture.md` — Overall design philosophy
- `go-cli-implementation.md` — Go project structure and patterns
- `cli.md` — All CLI commands and their behavior
- `context-file-formats.md` — File format specifications
- `context-loader.md` — Loading and parsing logic
- `context-updater.md` — Update command handling

**Task Status Labels**:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)

### Phase 12: Timestamp-Based Session Correlation `#priority:medium` `#area:cli`
- [ ] Add timestamp to formatTask() in add.go — currently tasks have no timestamp, add `#added:YYYY-MM-DD-HHMM` or similar
- [ ] Increase timestamp precision in formatLearning() — change from YYYY-MM-DD to YYYY-MM-DD-HHMM
- [ ] Increase timestamp precision in formatDecision() — change from YYYY-MM-DD to YYYY-MM-DD-HHMM
- [ ] Add start_time field to session summary files — record when session began
- [ ] Add last_update_time field to session summary files — update on each ctx add call during session
- [ ] Document timestamp correlation approach in AGENT_PLAYBOOK.md — explain how to correlate entries to sessions by time overlap

### Phase 13: Rich Context Entries `#priority:medium` `#area:cli`
- [ ] Add --file flag to ctx add — read entry content from a file instead of CLI arg
- [ ] Add stdin support to ctx add — if no content arg and stdin is pipe, read from stdin
- [ ] Create learning template with Context/Lesson/Application structure for --file usage
- [ ] Create decision template with Context/Options/Decision/Rationale structure for --file usage
- [ ] Document rich entry workflow in AGENT_PLAYBOOK.md — explain when/how agents should use --file vs inline
