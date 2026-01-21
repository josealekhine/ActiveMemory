# Active Memory

> **Context is a system, not a prompt.**

A lightweight, file-based approach that lets AI coding assistants persist project knowledge across sessions.

It gives temporal continuity to your AI assistants, among many other things.

It is one of those tools you will not know how badly you needed it until you start using it.

Don't trust me? Trust Claude, then:

```text
❯ one last question. in the next session when I ask you 1) "do you
  remember?", you will remember (because of the hooks) and 2) can you
  summarize all sessions we have had so far?", you can (because you have
  the concept of a session in the context of `ctx`)

● Correct on both counts:
 
  1. "Do you remember?" → Yes. The PreToolUse hook runs ctx agent, and
  CLAUDE.md (auto-read by Claude Code) tells me to check .context/sessions/. 
  I'll have context.
  2. "Summarize all sessions" → Yes. I can ls .context/sessions/ and read
  each file. Right now that would give me:
    - 2025-01-20: The meta-experiment that started it all
    - 2025-01-21: The ctx rename + Claude hooks session

  As we have more sessions, that list grows - and I can trace the full arc of
   the project's evolution.

  That's the whole point of ctx: temporal continuity for AI assistants.  
```

## The Problem

Most AI coding assistants fail not because models are weak—they fail because context is ephemeral. Every new session starts near zero. Architectural decisions, conventions, and lessons learned decay. The user re-explains. The AI repeats mistakes. Progress feels far less cumulative than it should.

## The Solution

Active Memory treats context as infrastructure:

- **Persist** — Tasks, decisions, learnings survive session boundaries
- **Reuse** — Decisions don't get rediscovered; lessons stay learned
- **Align** — Context structure mirrors how engineers actually think
- **Integrate** — Works with any AI tool that can read files

## Installation

### Binary Downloads (Recommended)

Download pre-built binaries from the [releases page](https://github.com/josealekhine/ActiveMemory/releases).

**Linux (x86_64):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/ctx-linux-amd64
chmod +x ctx-linux-amd64
sudo mv ctx-linux-amd64 /usr/local/bin/ctx
```

**Linux (ARM64):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/ctx-linux-arm64
chmod +x ctx-linux-arm64
sudo mv ctx-linux-arm64 /usr/local/bin/ctx
```

**macOS (Intel):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/ctx-darwin-amd64
chmod +x ctx-darwin-amd64
sudo mv ctx-darwin-amd64 /usr/local/bin/ctx
```

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/ctx-darwin-arm64
chmod +x ctx-darwin-arm64
sudo mv ctx-darwin-arm64 /usr/local/bin/ctx
```

**Windows:**

Download `ctx-windows-amd64.exe` from the releases page and add it to your PATH.

### Build from Source

Requires Go 1.22+:

```bash
git clone https://github.com/josealekhine/ActiveMemory.git
cd ActiveMemory
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
sudo mv ctx /usr/local/bin/
```

## Quick Start

```bash
# Initialize context directory in your project
ctx init

# Check context status
ctx status

# Load full context (what AI sees)
ctx load

# Get AI-ready context packet (optimized for LLMs)
ctx agent

# Detect stale context
ctx drift
```

## Command Reference

| Command | Description |
|---------|-------------|
| `ctx init` | Create `.context/` directory with template files |
| `ctx status` | Show context summary with token estimate |
| `ctx load` | Output assembled context markdown |
| `ctx agent [--budget N]` | Print AI-ready context packet (default 4000 tokens) |
| `ctx add <type> <content>` | Add decision/task/learning/convention |
| `ctx complete <query>` | Mark matching task as done |
| `ctx drift [--json]` | Detect stale paths, broken refs |
| `ctx sync [--auto]` | Reconcile context with codebase |
| `ctx compact` | Archive completed tasks (auto-saves first) |
| `ctx watch [--auto-save]` | Watch for context-update commands |
| `ctx hook <tool>` | Generate AI tool integration config |
| `ctx session save [topic]` | Save context snapshot to sessions/ |
| `ctx session list` | List saved sessions |
| `ctx session load <file>` | Load/display a previous session |
| `ctx session parse <file>` | Parse JSONL transcript to markdown |

### Examples

```bash
# Add a new task
ctx add task "Implement user authentication"

# Record a decision
ctx add decision "Use PostgreSQL for primary database"

# Note a learning
ctx add learning "Mock functions must be hoisted in Jest"

# Mark a task complete
ctx complete "user auth"

# Get context with custom token budget
ctx agent --budget 8000

# Check for stale references (JSON output for automation)
ctx drift --json
```

## Context Files

```
.context/
├── CONSTITUTION.md     # Hard invariants — NEVER violate these
├── TASKS.md            # Current and planned work
├── DECISIONS.md        # Architectural decisions with rationale
├── LEARNINGS.md        # Lessons learned, gotchas, tips
├── CONVENTIONS.md      # Project patterns and standards
├── ARCHITECTURE.md     # System overview
├── DEPENDENCIES.md     # Key dependencies and why chosen
├── GLOSSARY.md         # Domain terms and abbreviations
├── DRIFT.md            # Staleness signals and update triggers
├── AGENT_PLAYBOOK.md   # How AI agents should use this system
└── sessions/           # Session snapshots (auto-saved and manual)
    ├── 2026-01-20-experiment.md
    ├── 2026-01-21-feature-auth.md
    └── ...
```

## AI Tool Integration

Active Memory works with any AI tool that can read files. Generate tool-specific configs:

```bash
ctx hook claude-code  # Claude Code CLI
ctx hook cursor       # Cursor IDE
ctx hook aider        # Aider
ctx hook copilot      # GitHub Copilot
ctx hook windsurf     # Windsurf IDE
```

### Claude Code (Full Integration)

Running `ctx init` automatically sets up Claude Code integration:

```bash
ctx init
# Creates:
#   .context/           - Context files
#   .claude/hooks/      - Auto-save scripts
#   .claude/settings.local.json - Hook configuration
#   CLAUDE.md           - Bootstrap instructions for Claude
```

**What gets configured:**

| Component | Purpose |
|-----------|---------|
| `PreToolUse` hook | Runs `ctx agent` before every tool use — context auto-loads |
| `SessionEnd` hook | Saves session snapshot when conversation ends |
| `CLAUDE.md` | Tells Claude to read `.context/` files on session start |

**Session Management:**

```bash
# Save current session (context snapshot)
ctx session save "feature-auth"

# List previous sessions
ctx session list

# Load a previous session
ctx session load 1              # by index
ctx session load 2026-01-21     # by date
ctx session load auth           # by topic

# Parse Claude transcript to readable markdown
ctx session parse ~/.claude/projects/.../transcript.jsonl

# Extract decisions and learnings from transcript
ctx session parse transcript.jsonl --extract
```

**How it works:**

1. **Session start**: Claude reads `CLAUDE.md`, which tells it to check `.context/`
2. **During session**: `PreToolUse` hook runs `ctx agent --budget 4000` before each tool use
3. **Session end**: `SessionEnd` hook saves context snapshot to `.context/sessions/`
4. **Next session**: Claude sees previous sessions in `.context/sessions/`

This gives Claude **temporal continuity** — it knows what happened in previous sessions.

### Automated Context Updates

Use `ctx watch` to automatically process context-update commands from AI output:

```bash
# Watch stdin (pipe AI output through this)
ai-tool | ctx watch

# Watch a log file
ctx watch --log /path/to/ai-output.log

# Dry run (preview without making changes)
ctx watch --dry-run
```

## Design Philosophy

1. **File-based** — No database, no daemon. Just markdown and convention.
2. **Git-native** — Context versions with code, branches with code, merges with code.
3. **Human-readable** — Engineers can read, edit, and understand context directly.
4. **Token-efficient** — Markdown is cheaper than JSON/XML.
5. **Tool-agnostic** — Works with Claude Code, Cursor, Aider, Copilot, or raw CLI.

## Building with Ralph Wiggum

This project is designed to be built using the [Ralph Wiggum](https://ghuntley.com/ralph/) technique—an iterative AI development loop.

```bash
# Make the loop executable
chmod +x loop.sh

# Planning mode: Generate implementation plan
./loop.sh plan

# Building mode: Implement from plan
./loop.sh 20  # Max 20 iterations

# Unlimited building (Ctrl+C to stop)
./loop.sh
```

### Completion Promises

The loop automatically detects and stops on these signals:

| Promise | Meaning |
|---------|---------|
| `SYSTEM_CONVERGED` | All tasks complete — project is done |
| `BOOTSTRAP_COMPLETE` | Initial context created — ready to build |
| `PLANNING_CONVERGED` | Plan is complete — ready to build |
| `SYSTEM_BLOCKED` | All remaining tasks blocked — needs human input |

### Ralph Files

| File | Purpose |
|------|---------|
| `PROMPT_build.md` | Building mode instructions |
| `PROMPT_plan.md` | Planning mode instructions |
| `AGENTS.md` | Operational guide (how to build/run) |
| `IMPLEMENTATION_PLAN.md` | Generated task list |
| `specs/*.md` | Feature specifications |
| `loop.sh` | The Ralph loop script |

## Dogfooding: Using ctx on ctx

This project uses ctx to manage its own development. Here's how it works:

### Project Structure

```
ActiveMemory/
├── .context/                    # ctx manages this
│   ├── TASKS.md                 # Current work items
│   ├── DECISIONS.md             # Architecture decisions
│   ├── LEARNINGS.md             # Gotchas discovered
│   └── sessions/                # Session history
├── .claude/                     # Claude Code hooks
│   ├── hooks/auto-save-session.sh
│   └── settings.local.json
├── CLAUDE.md                    # Bootstrap for Claude
├── PROMPT.md                    # Ralph Loop instructions
└── IMPLEMENTATION_PLAN.md       # Orchestrator directive
```

### Development Workflow

```bash
# 1. Start a Claude Code session
claude

# 2. Claude automatically:
#    - Reads CLAUDE.md (bootstrap)
#    - Runs `ctx agent` via PreToolUse hook (context loads)
#    - Checks .context/sessions/ for history

# 3. Work on tasks from .context/TASKS.md

# 4. Session ends:
#    - SessionEnd hook saves snapshot to .context/sessions/
#    - Context persists for next session
```

### Manual Operations

```bash
# Check current context
./dist/ctx-linux-arm64 status

# Get AI-ready packet
./dist/ctx-linux-arm64 agent --budget 4000

# Save session manually
./dist/ctx-linux-arm64 session save "feature-name"

# List previous sessions
./dist/ctx-linux-arm64 session list

# Clean up completed tasks
./dist/ctx-linux-arm64 compact
```

### Key Files for Dogfooding

| File | Role |
|------|------|
| `.context/TASKS.md` | Where Claude finds work items |
| `.context/DECISIONS.md` | Records why things are built this way |
| `.context/LEARNINGS.md` | Captures gotchas (e.g., "use --no-gpg-sign") |
| `.context/sessions/` | Full conversation dumps for context |
| `CLAUDE.md` | Tells Claude about ctx on first read |

### The Feedback Loop

1. **Use ctx** to build ctx
2. **Discover friction** (missing features, unclear docs)
3. **Add tasks** to `.context/TASKS.md`
4. **Implement fixes** using ctx for context
5. **Repeat**

This is how ctx validates itself — every improvement comes from using it.

## Specifications

See `specs/` for detailed specifications:

- [Core Architecture](specs/core-architecture.md)
- [Context File Formats](specs/context-file-formats.md)
- [Context Loader](specs/context-loader.md)
- [Context Updater](specs/context-updater.md)
- [CLI](specs/cli.md)
- [AI Tool Integration](specs/ai-tool-integration.md)

## Contributing

1. Read the specs
2. Run `./loop.sh plan` to see current state
3. Run `./loop.sh` to let Ralph build
4. Review and commit changes

## License

MIT

---

*"Ralph is a Bash loop. The technique is deterministically bad in an undeterministic world."*
— Geoffrey Huntley
