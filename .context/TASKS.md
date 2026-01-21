# Tasks

## In Progress

## Next Up

### Project Rename `#priority:high` `#area:branding`
- [ ] Rename project from "Active Memory" to "Context"
  - Update README.md title and references
  - Update Go module path (github.com/josealekhine/ActiveMemory → TBD)
  - Update all import paths in Go files
  - Update CLAUDE.md references
  - Update any branding/taglines
  - Keep `ctx` as binary name (short for context)
  - Handle GitHub repo rename

### Auto-Save Enhancements `#priority:medium` `#area:cli`
- [x] Implement `ctx watch --auto-save` mode — 2026-01-21

### Documentation `#priority:medium` `#area:docs`
- [x] Document Claude Code integration in README — 2026-01-21
- [x] Add "Dogfooding Guide" — how to use ctx on ctx itself — 2026-01-21
- [x] Document session auto-save setup for new users — 2026-01-21

## Completed (Recent)

- [x] Extract key decisions/learnings from transcript (--extract flag on session parse) — 2026-01-21
- [x] Add PreCompact behavior — auto-save before `ctx compact` runs — 2026-01-21
- [x] Implement `ctx session parse` — convert .jsonl transcript to readable markdown — 2026-01-21
- [x] Implement `ctx session load <file>` — load/summarize a previous session — 2026-01-21
- [x] Implement `ctx session list` — list saved sessions with summaries — 2026-01-21
- [x] Implement `ctx session save` — manually dump context to sessions/ — 2026-01-21
- [x] Handle CLAUDE.md creation/merge in `ctx init` (template, backup, markers, --merge flag, idempotency) — 2026-01-21
- [x] Enhance `ctx init` to create Claude hooks (embedded scripts, settings.local.json, platform detection) — 2025-01-21
- [x] Set up PreToolUse hook for auto-load — 2025-01-20
- [x] Set up SessionEnd hook for auto-save — 2025-01-20
- [x] Create `.context/sessions/` directory structure — 2025-01-20
- [x] Create CLAUDE.md for native Claude Code bootstrapping — 2025-01-20
- [x] Document session persistence in AGENT_PLAYBOOK.md — 2025-01-20
- [x] Decide: always create .claude/ hooks (no --claude flag needed) — 2025-01-20

## Blocked
