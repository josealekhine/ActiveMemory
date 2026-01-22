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

### Hooks Should Use PATH, Not Hardcoded Paths
**Discovered**: 2026-01-21

**Context**: Original hooks used hardcoded absolute paths like `/home/user/project/dist/ctx-linux-arm64`. This caused issues when dogfooding or sharing configs.

**Lesson**: Hooks should assume `ctx` is in the user's PATH:
- More portable across machines/users
- Standard Unix practice
- `ctx init` now checks if `ctx` is in PATH before proceeding
- Hooks use `ctx agent` instead of `/full/path/to/ctx-linux-arm64 agent`

**Application**:
1. Users must install ctx to PATH: `sudo make install` or `sudo cp ./ctx /usr/local/bin/`
2. `ctx init` will fail with clear instructions if ctx is not in PATH
3. Tests can skip this check with `CTX_SKIP_PATH_CHECK=1`

**Supersedes**: Previous learning "Binary Path Must Be Absolute" (2025-01-20)

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

---

## Project Structure

### One Templates Directory, Not Two
**Discovered**: 2025-01-21

**Context**: Confusion arose about `templates/` (root) vs `internal/templates/` (embedded).

**Lesson**: Only `internal/templates/` matters — it's where Go embeds files into the binary. A root `templates/` directory is spec baggage that serves no purpose.

**The actual flow:**
```
internal/templates/  ──[ctx init]──>  .context/
     (baked into binary)              (agent's working copy)
```

**Application**: Don't create duplicate template directories. One source of truth.

### Orchestrator vs Agent Tasks
**Discovered**: 2025-01-21

**Context**: Ralph Loop checked `IMPLEMENTATION_PLAN.md`, found all tasks done, exited — ignoring `.context/TASKS.md`.

**Lesson**: Separate concerns:
- **`IMPLEMENTATION_PLAN.md`** = Orchestrator directive ("check your tasks")
- **`.context/TASKS.md`** = Agent's mind (actual task list)

The orchestrator shouldn't maintain a parallel ledger. It just says "check your mind."

**Application**: For new projects, `IMPLEMENTATION_PLAN.md` has ONE task: "Check `.context/TASKS.md`"

---

## Ralph Loop & Dogfooding

### Exit Criteria Must Include Verification, Not Just Task Completion
**Discovered**: 2026-01-21

**Context**: Dogfooding experiment had another Claude rebuild `ctx` from specs. All tasks were marked complete, Ralph Loop exited successfully. But the built binary didn't work — commands just printed help text instead of executing.

**Lesson**: "All tasks checked off" ≠ "Implementation works." This applies to US too, not just the dogfooding clone. Our own verification is based on manual testing, not automated proof. Blind spots exist in both projects.

Exit criteria must include:
- **Integration tests**: Binary executes commands correctly (not just unit tests)
- **Coverage targets**: Quantifiable proof that code paths are tested
- **Smoke tests**: Basic "does it run" verification in CI

**Application**:
1. Add integration test suite that invokes the actual binary
2. Set coverage targets (e.g., 70% for core packages)
3. Add verification tasks to TASKS.md — we have the same blind spot
4. Being proud of our achievement doesn't prove its validity
