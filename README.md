# Active Memory

> **Context is a system, not a prompt.**

A lightweight, file-based approach that lets AI coding assistants persist project knowledge across sessions.

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
| `ctx compact` | Archive completed tasks |
| `ctx watch [--log FILE]` | Watch for context-update commands |
| `ctx hook <tool>` | Generate AI tool integration config |

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
└── AGENT_PLAYBOOK.md   # How AI agents should use this system
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

### Claude Code

Add to your project's `CLAUDE.md`:

```markdown
## Active Memory Context

Before starting any task, load the project context:

1. Read .context/CONSTITUTION.md — These rules are INVIOLABLE
2. Read .context/TASKS.md — Current work items
3. Read .context/CONVENTIONS.md — Project patterns
4. Read .context/ARCHITECTURE.md — System overview
5. Read .context/DECISIONS.md — Why things are the way they are

When you make changes:
- Add decisions: <context-update type="decision">Your decision</context-update>
- Add tasks: <context-update type="task">New task</context-update>
- Add learnings: <context-update type="learning">What you learned</context-update>
- Complete tasks: <context-update type="complete">task description</context-update>

Run 'ctx agent' for a quick context summary.
```

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
