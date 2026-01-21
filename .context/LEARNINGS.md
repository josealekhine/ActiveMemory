# Learnings

## ctx vs Ralph Loop

### They Are Separate Systems
**Discovered**: 2025-01-20

**Context**: User asked "How do I use the ctx binary to recreate this project?"

**Lesson**: `ctx` and Ralph Loop are two distinct systems:
- `ctx init` creates `.context/` for context management (decisions, learnings, tasks)
- Ralph Loop uses PROMPT.md, IMPLEMENTATION_PLAN.md, specs/ for iterative AI development
- `ctx` does NOT create Ralph Loop infrastructure

**Application**: To bootstrap a new project with both:
1. Run `ctx init` to create `.context/`
2. Manually copy/adapt PROMPT.md, AGENTS.md, specs/ from a reference project
3. Create IMPLEMENTATION_PLAN.md with your tasks
4. Run `/ralph-loop` to start iterating

---

## Claude Code Integration

### Binary Path Must Be Absolute
**Discovered**: 2025-01-20

**Context**: Set up PreToolUse hook but referenced wrong binary path.

**Lesson**: The built binary is in `dist/ctx-linux-arm64`, not at project root. Use full path in hooks.

**Application**: Always verify binary location with `ls dist/` before configuring hooks. For this project: `/home/parallels/WORKSPACE/ActiveMemory/dist/ctx-linux-arm64`

### `.context/` Is NOT a Claude Code Primitive
**Discovered**: 2025-01-20

**Context**: User asked if Claude Code natively understands `.context/`.

**Lesson**: Claude Code only natively reads:
- `CLAUDE.md` (auto-loaded at session start)
- `.claude/settings.json` (hooks and permissions)

The `.context/` directory is an ctx convention. Claude won't know about it unless:
1. A hook runs `ctx agent` to inject context
2. CLAUDE.md explicitly instructs reading `.context/`

**Application**: Always create CLAUDE.md as the bootstrap entry point.

### Session Filename Must Include Time
**Discovered**: 2025-01-20

**Context**: Using just date (`2025-01-20-topic.md`) would overwrite multiple sessions per day.

**Lesson**: Use `YYYY-MM-DD-HHMM-<topic>.md` format to prevent overwrites.

**Application**: Always include hour+minute in session filenames.

### SessionEnd Hook Catches Ctrl+C
**Discovered**: 2025-01-20

**Context**: Needed to auto-save context even when user force-quits with Ctrl+C.

**Lesson**: Claude Code's `SessionEnd` hook fires on ALL exits including Ctrl+C. It provides:
- `transcript_path` - full session transcript (.jsonl)
- `reason` - why session ended (exit, clear, logout, etc.)
- `session_id` - unique session identifier

**Application**: Use SessionEnd hook to auto-save transcripts to `.context/sessions/`. See `.claude/hooks/auto-save-session.sh`.

---

## Context Persistence Patterns

### Two Tiers of Persistence
**Discovered**: 2025-01-20

**Context**: User wanted to ensure nothing is lost when session ends.

**Lesson**: Two levels serve different needs:

| Tier | Content | Purpose | Location |
|------|---------|---------|----------|
| Curated | Key learnings, decisions, tasks | Quick reload, token-efficient | `.context/*.md` |
| Full dump | Entire conversation | Safety net, deep dive | `.context/sessions/*.md` |

**Application**: Before session ends, save BOTH tiers.

### Auto-Load Works, Auto-Save Was Missing
**Discovered**: 2025-01-20

**Context**: Explored how to persist context across Claude Code sessions.

**Lesson**: Initial state was asymmetric:
- **Auto-load**: Works via `PreToolUse` hook running `ctx agent`
- **Auto-save**: Did NOT exist

**Solution implemented**: `SessionEnd` hook that copies transcript to `.context/sessions/`

---

## Init Command Design

### Always Backup Before Modifying User Files
**Discovered**: 2025-01-20

**Context**: `ctx init` needs to create/modify CLAUDE.md, but user may have existing customizations.

**Lesson**: When modifying user files (especially config files like CLAUDE.md):
1. **Always backup first** — `file.bak` before any modification
2. **Check for existing content** — use marker comments for idempotency
3. **Offer merge, don't overwrite** — respect user's customizations
4. **Provide escape hatch** — `--merge` flag for automation, manual merge for control

**Application**: Any `ctx` command that modifies user files should follow this pattern.

---

## Build & Platform

### CGO Must Be Disabled for ARM64 Linux
**Discovered**: During project build

**Context**: `go test` failed with `gcc: error: unrecognized command-line option '-m64'`

**Lesson**: On ARM64 Linux, CGO causes cross-compilation issues. Always use `CGO_ENABLED=0`.

**Application**:
```bash
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
CGO_ENABLED=0 go test ./...
```
