# Tasks

## In Progress

## Next Up

### Enhance `ctx init` to Create Claude Hooks `#priority:high` `#area:cli`
- [ ] Embed hook scripts in binary (like templates)
- [ ] Create `.claude/hooks/auto-save-session.sh` during init
- [ ] Create `.claude/settings.local.json` with PreToolUse and SessionEnd hooks
- [ ] Detect platform to set correct binary path in hooks
- [ ] Update `ctx init` output to mention Claude Code integration

### Handle CLAUDE.md Creation/Merge `#priority:high` `#area:cli`
- [ ] Create CLAUDE.md if it doesn't exist
- [ ] If CLAUDE.md exists, backup to CLAUDE.md.<unix_timestamp>.bak before any modification
- [ ] Detect existing ctx content via marker comment (`<!-- ctx:context -->`)
- [ ] If no ctx content, offer to merge (output snippet + prompt)
- [ ] Add `--merge` flag to auto-append without prompting
- [ ] Ensure idempotency — running init twice doesn't duplicate content

### Session Management Commands `#priority:high` `#area:cli`
- [ ] Implement `ctx session save` — manually dump context to sessions/
- [ ] Implement `ctx session list` — list saved sessions with summaries
- [ ] Implement `ctx session load <file>` — load/summarize a previous session
- [ ] Implement `ctx session parse` — convert .jsonl transcript to readable markdown

### Auto-Save Enhancements `#priority:medium` `#area:cli`
- [ ] Add PreCompact behavior — auto-save before `ctx compact` runs
- [ ] Extract key decisions/learnings from transcript automatically
- [ ] Consider `ctx watch --auto-save` mode

### Documentation `#priority:medium` `#area:docs`
- [ ] Document Claude Code integration in README
- [ ] Add "Dogfooding Guide" — how to use ctx on ctx itself
- [ ] Document session auto-save setup for new users

## Completed (Recent)

- [x] Set up PreToolUse hook for auto-load — 2025-01-20
- [x] Set up SessionEnd hook for auto-save — 2025-01-20
- [x] Create `.context/sessions/` directory structure — 2025-01-20
- [x] Create CLAUDE.md for native Claude Code bootstrapping — 2025-01-20
- [x] Document session persistence in AGENT_PLAYBOOK.md — 2025-01-20
- [x] Decide: always create .claude/ hooks (no --claude flag needed) — 2025-01-20

## Blocked
