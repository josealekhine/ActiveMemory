---
title: "Building ctx Using ctx: A Meta-Experiment in AI-Assisted Development"
date: 2026-01-27
author: Jose Alekhinne
---

*Jose Alekhinne / 2026-01-27*

# Building ctx Using ctx: A Meta-Experiment in AI-Assisted Development

> What happens when you build a tool designed to give AI memory, using that very 
same tool to remember what you are building? 

This is the story of `ctx`, how it evolved from a hasty "*YOLO mode*" experiment 
to a disciplined system for **persistent AI context**, and what I have 
learned along the way.

!!! info "Context is a Record"
    **Context** *is a* **persistent record**.
    
    By "*context*", I don’t mean model memory or stored thoughts: 
    
    I mean **the durable record of decisions, learnings, and intent** 
    that normally *evaporates* between sessions.

## AI Amnesia

Every developer who works with AI code generators knows the frustration: 
you have a deep, productive session where the AI understands your codebase, 
your conventions, your decisions. And then you close the terminal. 

Tomorrow it's a blank slate. The AI has forgotten everything.

That is "*reset amnesia*", and it's not just annoying: it's expensive. 

Every session starts with re-explaining context, re-reading files, 
re-discovering decisions that were already made.

> "I don't want to lose this discussion... 
> I am a brain-dead developer YOLO'ing my way out"

That's exactly what I said to Claude when I first started working on `ctx`.

## The Genesis

The project started as "*Active Memory*" (`amem`): a CLI tool to persist AI 
context across sessions. 

The core idea was simple: create a `.context/` directory with structured 
Markdown files for decisions, learnings, tasks, and conventions. 
The AI reads these at session start and writes to them before the session ends.

**The first commit** was just scaffolding. But within hours, the 
**Ralph Loop**—an iterative AI development workflow—had produced a working CLI:

```
feat(cli): implement amem init command
feat(cli): implement amem status command
feat(cli): implement amem add command
feat(cli): implement amem agent command
...
```

Fourteen core commands shipped in rapid succession. 

I was YOLO'ing like there was no tomorrow:
auto-accept every change, let the AI run free, ship features fast.

## The Meta-Experiment: Using `amem` to Build `amem`

Here's where it gets interesting: On January 20th, I asked: 

> **"Can I use `amem` to help you remember this context when I restart?"**

The answer was yes—but with a gap: 

Auto-load worked (*via Claude Code's `PreToolUse` hook*), but auto-save was 
missing. If the user quit with Ctrl+C, everything since the last manual save 
was lost.

That session became the first real test of the system. 

Here is the first session file we recorded:

```markdown
## Key Discussion Points

### 1. amem vs Ralph Loop - They're Separate Systems

**User's question**: "How do I use the binary to recreate this project?"

**Answer discovered**: amem is for context management, Ralph Loop is for 
development workflow. They're complementary but separate.

### 2. Two Tiers of Context Persistence

| Tier      | What                        | Why                           | Where                  |
|-----------|-----------------------------|-------------------------------|------------------------|
| Curated   | Learnings, decisions, tasks | Quick reload, token-efficient | .context/*.md          |
| Full dump | Entire conversation         | Safety net, nothing lost      | .context/sessions/*.md |
```

This session file—written by the AI to preserve its own context—became the 
template for how `ctx` handles session persistence.

## The Rename

By January 21st, I realized "*Active Memory*" was too generic, and (*arguably*)
too marketing-smelly. 

Besides, the binary was already called `ctx` (*short for Context*), 
the directory was `.context/`, and the slash commands would be `/ctx-*`. 

So it followed that the project should be renamed to `ctx` to make things 
make sense.

The rename touched **100+ files** but was clean—a find-and-replace with Go's 
type system catching any misses.

The git history tells the story:

```
0e8f6bb feat: rename amem to ctx and add Claude Code integration
87dcfa1 README.
4f0e195 feat: separate orchestrator directive from agent tasks
```

## YOLO Mode: Fast, But Dangerous

The Ralph Loop made feature development incredibly fast. 

But **it created technical debt that I didn't notice until later**.

A comparison session on January 25th revealed the patterns:

| YOLO Pattern                           | What We Found                                   |
|----------------------------------------|-------------------------------------------------|
| `"TASKS.md"` scattered in 10 files     | Same string literal everywhere, no constants    |
| `dir + "/" + file`                     | Should be `filepath.Join()`                     |
| Monolithic `cli_test.go` (1500+ lines) | Tests disconnected from implementations         |
| `package initcmd` in `init/` folder    | Go's "init" is reserved—subtle naming collision |

Here is another analysis made by `ctx`:

```text
● Based on my analysis, here are the key differences:

YOLO Mode (Pre-040ce99)
- Feature-first: Added slash commands, tests, templates rapidly
- Scattered magic strings: "TASKS.md", "decision", ".context" spread across files
- Quick file creation: New files without organizational patterns
- Working but inconsistent: Code functioned but lacked systematic structure

Human-Guided Mode (Post-040ce99)
- Consolidation focus: Centralized constants in config package
- Consistent naming: Dir, File, Filename, UpdateType prefixes
- Self-referential constants: FileType map uses constants as keys, not literals
- Proper path construction: filepath.Join() instead of +"/"+
- Colocated tests: Tests next to implementations
- Canonical naming: Package name = folder name
```

The fix required a human-guided refactoring session.

We introduced `internal/config/config.go` with semantic prefixes:

```go
const (
    DirContext     = ".context"
    DirArchive     = "archive"
    DirSessions    = "sessions"
    FilenameTask   = "TASKS.md"
    UpdateTypeTask = "task"
)
```

What I begrudgingly learned was: 
**YOLO mode is effective for velocity but accumulates debt**. 

So I took a mental note to schedule **periodic consolidation sessions** 
from that point onward.

## The Dogfooding Test That Failed

On January 21st, I ran an experiment: have another Claude instance rebuild 
`ctx` from scratch using only the specs and `PROMPT.md`. 

The Ralph Loop ran, all tasks got checked off, the loop exited successfully.

**But the binary was broken**!

Commands just printed help text instead of executing. 

All tasks were marked "**complete**" but the implementation didn't work.

Here's what `ctx` discovered:

```markdown
## Key Findings

### Dogfooding Binary Is Broken
- Commands don't execute — they just print root help text
- All tasks were marked complete but binary doesn't work
- Lesson: "tasks checked off" ≠ "implementation works"
```

This was humbling—to say the least.

I realized, I had the same blind spot in my own codebase:
no integration tests that actually invoked the binary. 

So I added:

- Integration tests for all commands
- Coverage targets (60-80% per package)
- Smoke tests in CI
- A **constitution** rule: "**All code must pass tests before commit**"

## The Constitution versus Conventions

As lessons accumulated, there was the temptation to add everything to 
`CONSTITUTION.md` as "inviolable rules". 

But I resisted.

The constitution should contain only truly inviolable invariants:
- Security (*no secrets, no customer data*)
- Quality (*tests must pass*)
- Process (*decisions need records*)
- `ctx` invocation (*always use `PATH`, never fallback*)

Everything else—coding style, file organization, naming 
conventions—should go in to `CONVENTIONS.md`. 

Here's how `ctx` explained why the distinction was important: 

> "Overly strict constitution creates friction and gets ignored. 
> Conventions can be bent; constitution cannot."
> — Decision record, 2026-01-25

## Hooks: Harder Than They Look

Claude Code hooks seemed simple: run a script before/after certain events. 

But I hit multiple gotchas:

**1. Key names matter**

```text
// WRONG - "Invalid key in record" error
"PreToolUseHooks": [...]

// RIGHT
"PreToolUse": [...]
```

**2. Blocking requires specific output**

```bash
# WRONG - just exits, doesn't block
exit 1

# RIGHT - JSON output + exit 0
echo '{"decision": "block", "reason": "Use ctx from PATH"}'
exit 0
```

**3. Go's JSON escaping**

`json.Marshal` escapes `>`, `<`, `&` as unicode (`\u003e`) by default. 

When generating shell commands in JSON:

```go
encoder := json.NewEncoder(file)
encoder.SetEscapeHTML(false) // Prevent 2>/dev/null → 2\u003e/dev/null
```

**4. Regex overfitting**

Our hook to block non-PATH `ctx` invocations initially matched too broadly:

```bash
# WRONG - matches /home/user/ctx/internal/file.go (ctx as directory)
(/home/|/tmp/|/var/)[^ ]*ctx[^ ]*

# RIGHT - matches ctx as binary only
(/home/|/tmp/|/var/)[^ ]*/ctx( |$)
```

## The Session Files

By the time of this writing this project's `ctx` sessions (`.context/sessions/`) 
contains 40+ files from this project's development.

They are not part of the source code due to security, privacy, and size concerns.

However, they are invaluable for the project's progress.

Each **session file** is a timestamped Markdown with:

- Summary of what has been accomplished
- Key decisions made
- Learnings discovered
- Tasks for the next session
- Technical context (*platform, versions*)

These files are **not auto-loaded** (*that would bust the token budget*). 

They are what I see as the "*archaeological record*" of `ctx`:
When the AI needs deeper information about why something was done, it
digs into the sessions.

Auto-generated session files use a naming convention:

```
2026-01-23-115432-session-prompt_input_exit-summary.md
2026-01-25-220244-manual-save.md
2026-01-27-052107-session-other-summary.md
```

Also, the `SessionEnd` hook captures transcripts automatically. 
Even `Ctrl+C `is caught.

## The Decision Log: 18 Architectural Decisions

`ctx` helps record every significant architectural choice in 
`.context/DECISIONS.md`. 

Here are some highlights:

**Reverse-chronological order (2026-01-27)**

```markdown
**Context**: With chronological order, oldest items consume tokens first, and
newest (most relevant) items risk being truncated.

**Decision**: Use reverse-chronological order (newest first) for DECISIONS.md
and LEARNINGS.md.
```

**PATH over hardcoded paths (2026-01-21)**

```markdown
**Context**: Original implementation hardcoded absolute paths in hooks.
This breaks when sharing configs with other developers.

**Decision**: Hooks use `ctx` from PATH. `ctx init` checks PATH before proceeding.
```

**Generic core with Claude enhancements (2026-01-20)**

```markdown
**Context**: ctx should work with any AI tool, but Claude Code users could
benefit from deeper integration.

**Decision**: Keep ctx generic as the core tool, but provide optional
Claude Code-specific enhancements.
```

## The Learning Log: 24 Gotchas and Insights

The `.context/LEARNINGS.md` file captures gotchas that would otherwise be 
forgotten. Each has Context, Lesson, and Application sections:

**CGO on ARM64**

```markdown
**Context**: `go test` failed with `gcc: error: unrecognized command-line option '-m64'`
**Lesson**: On ARM64 Linux, CGO causes cross-compilation issues. Always use `CGO_ENABLED=0`.
```

**Claude Code skills format**

```markdown
**Lesson**: Claude Code skills are Markdown files in .claude/commands/ with `YAML`
frontmatter (*description, argument-hint, allowed-tools*). Body is the prompt.
```

**"Do you remember?" handling**

```markdown
**Lesson**: In a `ctx`-enabled project, "*do you remember?*" has an obvious meaning:
check the `.context/` files. Don't ask for clarification—just do it.
```

## Task Archives: The Completed Work

Completed tasks are archived to `.context/archive/` with timestamps. 

The archive from January 23rd shows 13 phases of work:

- Phase 1: Project Scaffolding (Go module, Cobra CLI)
- Phase 2-4: Core Commands 
  (init, status, agent, add, complete, drift, sync, compact, watch, hook)
- Phase 5: Session Management (save, list, load, parse, --extract)
- Phase 6: Claude Code Integration (hooks, settings, CLAUDE.md handling)
- Phase 7: Testing & Verification
- Phase 8: Task Archival
- Phase 9: Slash Commands
- Phase 9b: Ralph Loop Integration
- Phase 10: Project Rename
- Phase 11: Documentation
- Phase 12: Timestamp Correlation
- Phase 13: Rich Context Entries

That's an impressive 173 commits across 8 days of development.

## What I Learned About AI-Assisted Development

**1. Memory changes everything**

When the AI remembers decisions, it doesn't repeat mistakes. When it knows 
your conventions, it follows them. 

`ctx` makes the AI a better collaborator because it's not starting from zero.

**2. Two-tier persistence works**

Curated context (`DECISIONS.md`, `LEARNINGS.md`, `TASKS.md`) is for 
**quick reload**. 

Full session dumps are for **archaeology**. 

It's a futile effort to try to fit everything in the token budget.

**Persist more, load less**.

**3. YOLO mode has its place**

For rapid prototyping, letting the AI run free is effective. 

But I had to schedule **consolidation sessions**. 

**Technical debt accumulates silently**.

**4. The constitution should be small**

Only truly inviolable rules go in `CONSTITUTION.md`. 
Everything else is a convention. 

If you put too much in the constitution, it will get ignored.

**5. Verification is non-negotiable**

"All tasks complete" means nothing if you haven't run the tests. 

Integration tests that invoke the actual binary caught bugs that 
the unit tests missed.

**6. Session files are underrated**

The ability to grep through 40 session files and find exactly when and why a 
decision was made helped me a lot. 

It's not about loading them into context: It is about having them when you 
need them.

## The Future: Recall System

The next phase of `ctx` is the **Recall System**:

- **Parser**: Parse session capture markdowns, enrich with JSONL data
- **Renderer**: Goldmark + Chroma for syntax highlighting, dark mode UI
- **Server**: Local HTTP server for browsing sessions
- **Search**: Inverted index for searching across sessions
- **CLI**: `ctx recall serve <path>` to start the server

The goal is to make the archaeological record browsable—not just `grep`-able.

Because not everyone always lives in the terminal—me included.

## Conclusion

Building `ctx` using ctx was a meta-experiment in AI-assisted development. 

I learned that **memory isn't just convenient—it's transformative**:

* An AI that remembers your decisions doesn't repeat mistakes.
* An AI that knows your conventions doesn't need them re-explained.

If you are reading this, chances are that you already have heard about `ctx`.

* `ctx` is open source at 
[github.com/ActiveMemory/ctx](https://github.com/ActiveMemory/ctx),
* and the documentation lives at [ctx.ist](https://ctx.ist).

If you're a mere mortal tired of reset amnesia, give `ctx` a try. 

And when you do, check `.context/sessions/` sometime. 

The archaeological record might surprise you.

---

*This blog post was written with the help of `ctx` with full access to the 
`ctx` session files, decision log, learning log, task archives, and 
git history of `ctx—The meta continues.*
