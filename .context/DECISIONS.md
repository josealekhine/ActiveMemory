# Decisions

## [2025-01-20] Generic Core with Optional Claude Code Enhancements

**Status**: Accepted

**Context**: `ctx` should work with any AI tool, but Claude Code users could benefit from deeper integration (auto-load, auto-save via hooks).

**Decision**: Keep `ctx` generic as the core tool, but provide optional Claude Code-specific enhancements:
- `ctx hook claude-code` generates Claude-specific configs
- `.claude/hooks/` contains Claude Code hook scripts
- Features work without Claude Code, but are enhanced with it

**Rationale**:
- Maintains tool-agnostic philosophy from core-architecture.md
- Doesn't lock users into Claude Code
- Claude Code users get seamless experience without extra work
- Other AI tools can be supported similarly (`ctx hook cursor`, etc.)

**Consequences**:
- Need to maintain both generic and Claude-specific documentation
- Hook scripts are optional, not required
- Testing must cover both with and without Claude Code

---

## [2025-01-20] Always Generate Claude Hooks in Init (No Flag Needed)

**Status**: Accepted (to be implemented)

**Context**: Setting up Claude Code hooks manually is error-prone. Considered `--claude` flag but realized it's unnecessary.

**Decision**: `ctx init` ALWAYS creates `.claude/hooks/` alongside `.context/`:
```bash
ctx init    # Creates BOTH .context/ AND .claude/hooks/
```

**Rationale**:
- Other AI tools (Cursor, Aider, Copilot) don't know/care about `.claude/`
- No downside to creating hooks that sit unused
- Claude Code users get seamless experience with zero extra steps
- If user later switches to Claude Code, hooks are already there
- Simpler UX - no flags to remember

**Consequences**:
- `ctx init` creates both directories always
- Hook scripts are embedded in binary (like templates)
- Need to detect platform for binary path in hooks
- `.claude/` becomes part of ctx's standard output

---

## [2025-01-20] Two-Tier Context Persistence Model

**Status**: Accepted

**Context**: Need to persist context across sessions. Token budgets limit what can be loaded. But nothing should be truly lost.

**Decision**: Implement two tiers of persistence:

| Tier | Purpose | Location | Token Cost |
|------|---------|----------|------------|
| **Curated** | Quick context reload | `.context/*.md` | Low (budgeted) |
| **Full dump** | Safety net, archaeology | `.context/sessions/*.md` | Zero (not auto-loaded) |

**Rationale**:
- Curated context is token-efficient for daily use
- Full dumps ensure nothing is ever truly lost
- Users can dive into sessions/ when they need deep context
- Separation prevents context bloat

**Consequences**:
- Need both manual and automatic ways to populate both tiers
- Session files grow over time (may need archival strategy)
- `ctx agent` only loads curated tier by default

---

## [2025-01-20] Session Filename Format: YYYY-MM-DD-HHMMSS-topic.md

**Status**: Accepted

**Context**: Multiple sessions per day would overwrite each other. Also, multiple compacts in the same minute could collide.

**Decision**: Use `YYYY-MM-DD-HHMMSS-<topic>.md` format for session files. Two file types:
- **Curated sessions**: `HHMMSS-<topic>.md` - updated throughout session
- **Auto-snapshots**: `HHMMSS-<event>.jsonl` - immutable once created

**Rationale**:
- Human-readable (unlike unix timestamps)
- Naturally sorts chronologically
- Seconds precision prevents collision even with rapid compacts
- Clear distinction between curated notes and raw snapshots

**Consequences**:
- Slightly longer filenames
- Must ensure consistent format in all session-saving code
- Curated files keep getting updated; snapshots are write-once

---

## [2025-01-20] Auto-Save Before Compact

**Status**: Accepted (to be implemented)

**Context**: `ctx compact` archives old tasks. Information could be lost if not captured.

**Decision**: `ctx compact` should auto-save a session dump before archiving:
1. Save current state to `.context/sessions/YYYY-MM-DD-HHMM-pre-compact.md`
2. Then perform the compaction

**Rationale**:
- Safety net before destructive-ish operation
- User can always recover pre-compact state
- No extra user action required

**Consequences**:
- Compact command becomes slightly slower
- Sessions directory grows with each compact
- May want `--no-save` flag for automation

---

## [2025-01-20] Handle CLAUDE.md Creation/Merge in ctx init

**Status**: Accepted (to be implemented)

**Context**: Both `claude init` and `ctx init` want to create/modify CLAUDE.md. Users of ctx will likely want ctx's context-aware version, but may already have a CLAUDE.md from `claude init`.

**Decision**: `ctx init` handles CLAUDE.md intelligently:
- **No CLAUDE.md exists** → Create it with ctx's context-loading template
- **CLAUDE.md exists** → Don't overwrite. Instead:
  1. **Backup first** → Copy to `CLAUDE.md.<unix_timestamp>.bak` (e.g., `CLAUDE.md.1737399000.bak`)
  2. Check if it already has ctx content (idempotent check via marker comment)
  3. If not, output the snippet to append and offer to merge
  4. `ctx init --merge` flag to auto-append without prompting

**Rationale**:
- Timestamped backups preserve history across multiple runs
- Unix timestamp is fine for backups (rarely read by humans, easy to sort)
- Respects user's existing CLAUDE.md customizations
- Doesn't silently overwrite important config
- Idempotency prevents duplicate content on re-runs

**Consequences**:
- Need to detect existing ctx content (marker comment like `<!-- ctx:context -->`)
- Backup files accumulate: `CLAUDE.md.<timestamp>.bak` (may want cleanup command later)
- Init output must clearly show what was created vs what needs manual merge
- Should work gracefully even if user runs `ctx init` multiple times

---

## [2025-01-20] Use SessionEnd Hook for Auto-Save

**Status**: Accepted (implemented)

**Context**: Need to save context even when user exits abruptly (Ctrl+C).

**Decision**: Use Claude Code's `SessionEnd` hook to auto-save transcript:
- Hook fires on all exits including Ctrl+C
- Copies `transcript_path` to `.context/sessions/`
- Creates both .jsonl (raw) and .md (summary) files

**Rationale**:
- Catches all exit scenarios
- Transcript contains full conversation
- No user action required
- Graceful degradation (just doesn't save if hook fails)

**Consequences**:
- Only works with Claude Code (other tools need different approach)
- Requires jq for JSON parsing in hook script
- Session files are .jsonl format (need tooling to read)

---

## [2025-01-21] Separate Orchestrator Directive from Agent Tasks

**Status**: Accepted

**Context**: Two task systems existed: `IMPLEMENTATION_PLAN.md` (Ralph Loop orchestrator) and `.context/TASKS.md` (ctx's own context). Ralph would find IMPLEMENTATION_PLAN.md complete and exit, ignoring .context/TASKS.md.

**Decision**: Clean separation of concerns:
- **`.context/TASKS.md`** = Agent's mind. Tasks the agent decided need doing.
- **`IMPLEMENTATION_PLAN.md`** = Orchestrator's directive. A single meta-task: "Check your tasks."

The orchestrator doesn't maintain a parallel ledger — it just tells the agent to check its own mind.

**Rationale**:
- Agent autonomy: the agent owns its task list
- Single source of truth for tasks
- Orchestrator is minimal, not a micromanager
- Fresh `ctx init` deployments can have one directive: "Check .context/TASKS.md"
- Prevents task list drift between two files

**Consequences**:
- `PROMPT.md` now references `.context/TASKS.md` for task selection
- `IMPLEMENTATION_PLAN.md` becomes a thin directive layer
- Historical milestones are archived, not active tasks
- North Star goals live in IMPLEMENTATION_PLAN.md (meta-level, not tasks)
