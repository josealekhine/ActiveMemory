# Learnings

> Ordered reverse-chronologically (newest first) to ensure most relevant items 
> are read first regardless of token budget.

---

## [2026-01-28-191951] Required flags now enforced for learnings

**Context**: Implemented ctx add learning flags to match decision's ADR pattern

**Lesson**: Structured entries with Context/Lesson/Application are more useful than one-liners

**Application**: Always use ctx add learning with all three flags; agents guided via AGENT_PLAYBOOK.md

- **[2026-01-28-072838]** Changelogs document WHAT; blogs explain WHY. Same information, different engagement. Changelogs are for machines (audits, dependency trackers). Blogs are for humans (narrative, context, lessons). When synthesizing session history, output both: changelog for completeness, blog for readability.

- **[2026-01-28-051426]** IDE is already the UI: Discovery, search, and 
editing of .context/ markdown files works better in VS Code/IDE than any 
custom UI we'd build. Full-text search, git integration, extensions - all free. 
Don't reinvent the editor.

- **[2026-01-28-040915]** Subtasks complete does not mean parent task complete

- **[2026-01-28-040251]** AI session JSONL formats are not standardized across tools

---

## [2026-01-27] Always Complete Decision Record Sections

**Context**: Decisions added via `ctx add decision` were left with placeholder 
text like "[Add context here]".

**Lesson**: When recording decisions, always fill in Context 
(what prompted this), Rationale (why this choice over alternatives), and 
Consequences (what changes as a result). Placeholder text is a code smell - 
decisions without rationale lose their value over time.

**Application**: After using `ctx add decision`, immediately edit the file to 
complete all sections. Future: use `--context`, `--rationale`, `--consequences` 
flags when available.

---

## [2026-01-27] Slash Commands Require Matching Permissions

**Context**: Claude Code slash commands using `!` bash syntax require matching 
permissions in settings.local.json.

**Lesson**: When adding new /ctx-* commands, ensure ctx init pre-seeds the 
required `Bash(ctx <subcommand>:*)` permissions. Use additive merging for user 
config - never remove existing permissions.

---

## [2026-01-26] Go json.Marshal Escapes Shell Characters

**Context**: Generated settings.local.json had `2\u003e/dev/null` instead 
of `2>/dev/null`.

**Lesson**: Go's json.Marshal escapes `>`, `<`, and `&` as 
unicode (`\u003e`, `\u003c`, `\u0026`) by default for HTML safety. 
Use `json.Encoder` with `SetEscapeHTML(false)` when generating config files 
that contain shell commands.

---

## [2026-01-26] Claude Code Hook Key Names

**Context**: Hooks weren't working, getting "Invalid key in record" errors.

**Lesson**: Claude Code settings.local.json hook keys are `PreToolUse` and 
`SessionEnd` (not `PreToolUseHooks`/`SessionEndHooks`). The `Hooks` suffix 
causes validation errors.

---

## [2026-01-25] defer os.Chdir Fails errcheck Linter

**Context**: `defer os.Chdir(originalDir)` fails golangci-lint errcheck.

**Lesson**: Use `defer func() { _ = os.Chdir(x) }()` to explicitly ignore the 
error return value.

---

## [2026-01-25] golangci-lint Go Version Mismatch in CI

**Context**: CI was failing with Go version mismatches between golangci-lint 
and the project.

**Lesson**: When golangci-lint is built with an older Go version than the 
project targets, use `install-mode: goinstall` in CI to build the linter from 
source using the project's Go version.

---

## [2026-01-25] CI Tests Need CTX_SKIP_PATH_CHECK

**Context**: CI tests were failing because ctx binary isn't installed on CI runners.

**Lesson**: Tests that call `ctx init` will fail without `CTX_SKIP_PATH_CHECK=1` 
env var, because init checks if ctx is in PATH.

---

## [2026-01-25] AGENTS.md Is Not Auto-Loaded

**Context**: Had both AGENTS.md and CLAUDE.md in project root, causing confusion.

**Lesson**: Only CLAUDE.md is read automatically by Claude Code. Projects 
using ctx should rely on the CLAUDE.md → AGENT_PLAYBOOK.md chain, not AGENTS.md.

---

## [2026-01-25] Hook Regex Can Overfit

**Context**: `.claude/hooks/block-non-path-ctx.sh` was blocking legitimate sed 
commands because the regex `ctx[^ ]*` matched paths containing "ctx" as a 
directory component (e.g., `/home/user/ctx/internal/...`).

**Lesson**: When writing shell hook regexes:
- Test against paths that contain the target string as a substring
- `ctx` as binary vs `ctx` as directory name are different
- Original: `(/home/|/tmp/|/var/)[^ ]*ctx[^ ]* ` — overfits
- Fixed: `(/home/|/tmp/|/var/)[^ ]*/ctx( |$)` — matches binary only

**Application**: Always test hooks with edge cases before deploying.

---

## [2026-01-25] Autonomous Mode Creates Technical Debt

**Context**: Compared commits from autonomous "YOLO mode" (auto-accept, 
agent-driven) vs human-guided refactoring sessions.

**Lesson**: YOLO mode is effective for feature velocity but accumulates technical debt:

| YOLO Pattern                           | Human-Guided Fix                      |
|----------------------------------------|---------------------------------------|
| `"TASKS.md"` scattered in 10 files     | `config.FilenameTask` constant        |
| `dir + "/" + file`                     | `filepath.Join(dir, file)`            |
| `{"task": "TASKS.md"}`                 | `{UpdateTypeTask: FilenameTask}`      |
| Monolithic `cli_test.go` (1500+ lines) | Colocated `package/package_test.go`   |
| `package initcmd` in `init/` folder    | `package initialize` in `initialize/` |

**Application**:
1. Schedule periodic consolidation sessions (not just feature sprints)
2. When same literal appears 3+ times, extract to constant
3. Constants should reference constants (self-referential maps)
4. Tests belong next to implementations, not in monoliths

---

## [2026-01-23] ctx agent vs Manual File Reading Trade-offs

**Context**: User asked "Do you remember?" and agent used parallel file reads 
instead of `ctx agent`. Compared outputs to understand the delta.

**Lesson**: `ctx agent` is optimized for task execution:
- Filters to pending tasks only
- Surfaces constitution rules inline
- Provides prioritized read order
- Token-budget aware

Manual file reading is better for exploratory/memory questions:
- Session history access
- Timestamps ("modified 8 min ago")
- Completed task context
- Parallel reads for speed

**Application**: No need to mandate one approach. Agents naturally pick appropriately:
- "Do you remember?" → parallel file reads (need history)
- "What should I work on?" → `ctx agent` (need tasks)

---

## [2026-01-23] Claude Code Skills Format

**Context**: Needed to understand how to create custom slash commands.

**Lesson**: Claude Code skills are markdown files in `.claude/commands/` with 
YAML frontmatter (`description`, `argument-hint`, `allowed-tools`). Body is 
the prompt. Use code blocks with `!` prefix for shell execution. `$ARGUMENTS` 
passes command args.

---

## [2026-01-23] Infer Intent on "Do You Remember?" Questions

**Context**: User asked "Do you remember?" at session start. Agent asked for 
clarification instead of proactively checking context files.

**Lesson**: In a ctx-enabled project, "do you remember?" has an obvious 
meaning: check the `.context/` files and report what you know from previous 
sessions. Don't ask for clarification - just do it.

**Application**: When user asks memory-related questions ("do you remember?", 
"what were we working on?", "where did we leave off?"), immediately:
1. Read `.context/TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`
2. Check `.context/sessions/` for recent session files
3. Summarize what you find

Don't ask "would you like me to check the context files?" - that's the 
obvious intent.

---

## [2026-01-23] Always Use ctx from PATH

**Context**: Agent used `./dist/ctx-linux-arm64` and `go run ./cmd/ctx` 
instead of just `ctx`, even though the binary was installed to PATH.

**Lesson**: When working on a ctx-enabled project, always use `ctx` directly:
```bash
ctx status        # correct
ctx agent         # correct
./dist/ctx        # avoid hardcoded paths
go run ./cmd/ctx  # avoid unless developing ctx itself
```

**Application**: Check `which ctx` if unsure. The binary is installed during 
setup (`sudo make install` or `sudo cp ./ctx /usr/local/bin/`).

---

## [2026-01-21] Exit Criteria Must Include Verification

**Context**: Dogfooding experiment had another Claude rebuild `ctx` from specs. 
All tasks were marked complete, Ralph Loop exited successfully. But the built 
binary didn't work — commands just printed help text instead of executing.

**Lesson**: "All tasks checked off" ≠ "Implementation works." This applies to 
US too, not just the dogfooding clone. Our own verification is based on manual 
testing, not automated proof. Blind spots exist in both projects.

Exit criteria must include:
- **Integration tests**: Binary executes commands correctly (not just unit tests)
- **Coverage targets**: Quantifiable proof that code paths are tested
- **Smoke tests**: Basic "does it run" verification in CI

**Application**:
1. Add integration test suite that invokes the actual binary
2. Set coverage targets (e.g., 70% for core packages)
3. Add verification tasks to TASKS.md — we have the same blind spot
4. Being proud of our achievement doesn't prove its validity

---

## [2026-01-21] Orchestrator vs Agent Tasks Must Be Separate

**Context**: Ralph Loop checked `IMPLEMENTATION_PLAN.md`, found all tasks 
done, exited — ignoring `.context/TASKS.md`.

**Lesson**: Separate concerns:
- **`IMPLEMENTATION_PLAN.md`** = Orchestrator directive ("check your tasks")
- **`.context/TASKS.md`** = Agent's mind (actual task list)

The orchestrator shouldn't maintain a parallel ledger. It just says 
"check your mind."

**Application**: For new projects, `IMPLEMENTATION_PLAN.md` has ONE task: 
"Check `.context/TASKS.md`"

---

## [2026-01-21] One Templates Directory, Not Two

**Context**: Confusion arose about `templates/` (root) vs 
`internal/templates/` (embedded).

**Lesson**: Only `internal/templates/` matters — it's where Go embeds files 
into the binary. A root `templates/` directory is spec baggage that serves 
no purpose.

**The actual flow:**
```
internal/templates/  ──[ctx init]──>  .context/
     (baked into binary)              (agent's working copy)
```

**Application**: Don't create duplicate template directories. One source of truth.

---

## [2026-01-21] Hooks Should Use PATH, Not Hardcoded Paths

**Context**: Original hooks used hardcoded absolute paths like 
`/home/user/project/dist/ctx-linux-arm64`. This caused issues when dogfooding 
or sharing configs.

**Lesson**: Hooks should assume `ctx` is in the user's PATH:
- More portable across machines/users
- Standard Unix practice
- `ctx init` now checks if `ctx` is in PATH before proceeding
- Hooks use `ctx agent` instead of `/full/path/to/ctx-linux-arm64 agent`

**Application**:
1. Users must install ctx to PATH: `sudo make install` or `sudo cp ./ctx /usr/local/bin/`
2. `ctx init` will fail with clear instructions if ctx is not in PATH
3. Tests can skip this check with `CTX_SKIP_PATH_CHECK=1`

**Supersedes**: Previous learning "Binary Path Must Be Absolute" (2026-01-20)

---

## [2026-01-20] ctx and Ralph Loop Are Separate Systems

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

## [2026-01-20] .context/ Is NOT a Claude Code Primitive

**Context**: User asked if Claude Code natively understands `.context/`.

**Lesson**: Claude Code only natively reads:
- `CLAUDE.md` (auto-loaded at session start)
- `.claude/settings.json` (hooks and permissions)

The `.context/` directory is a ctx convention. Claude won't know about it unless:
1. A hook runs `ctx agent` to inject context
2. CLAUDE.md explicitly instructs reading `.context/`

**Application**: Always create CLAUDE.md as the bootstrap entry point.

---

## [2026-01-20] SessionEnd Hook Catches Ctrl+C

**Context**: Needed to auto-save context even when user force-quits with Ctrl+C.

**Lesson**: Claude Code's `SessionEnd` hook fires on ALL exits including Ctrl+C. It provides:
- `transcript_path` - full session transcript (.jsonl)
- `reason` - why session ended (exit, clear, logout, etc.)
- `session_id` - unique session identifier

**Application**: Use SessionEnd hook to auto-save transcripts to 
`.context/sessions/`. See `.claude/hooks/auto-save-session.sh`.

---

## [2026-01-20] Session Filename Must Include Time

**Context**: Using just date (`2026-01-20-topic.md`) would overwrite multiple sessions per day.

**Lesson**: Use `YYYY-MM-DD-HHMM-<topic>.md` format to prevent overwrites.

**Application**: Always include hour+minute in session filenames.

---

## [2026-01-20] Two Tiers of Persistence

**Context**: User wanted to ensure nothing is lost when session ends.

**Lesson**: Two levels serve different needs:

| Tier      | Content                         | Purpose                       | Location                 |
|-----------|---------------------------------|-------------------------------|--------------------------|
| Curated   | Key learnings, decisions, tasks | Quick reload, token-efficient | `.context/*.md`          |
| Full dump | Entire conversation             | Safety net, deep dive         | `.context/sessions/*.md` |

**Application**: Before session ends, save BOTH tiers.

---

## [2026-01-20] Auto-Load Works, Auto-Save Was Missing

**Context**: Explored how to persist context across Claude Code sessions.

**Lesson**: Initial state was asymmetric:
- **Auto-load**: Works via `PreToolUse` hook running `ctx agent`
- **Auto-save**: Did NOT exist

**Solution implemented**: `SessionEnd` hook that copies transcript to `.context/sessions/`

---

## [2026-01-20] Always Backup Before Modifying User Files

**Context**: `ctx init` needs to create/modify CLAUDE.md, but user may have existing customizations.

**Lesson**: When modifying user files (especially config files like CLAUDE.md):
1. **Always backup first** — `file.bak` before any modification
2. **Check for existing content** — use marker comments for idempotency
3. **Offer merge, don't overwrite** — respect user's customizations
4. **Provide escape hatch** — `--merge` flag for automation, manual merge for control

**Application**: Any `ctx` command that modifies user files should follow this pattern.

---

## [undated] CGO Must Be Disabled for ARM64 Linux

**Context**: `go test` failed with `gcc: error: unrecognized command-line option '-m64'`

**Lesson**: On ARM64 Linux, CGO causes cross-compilation issues. Always use `CGO_ENABLED=0`.

**Application**:
```bash
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
CGO_ENABLED=0 go test ./...
```
