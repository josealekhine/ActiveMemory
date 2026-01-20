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

## Quick Start

```bash
# Install (download binary for your platform)
# See: https://github.com/josealekhine/ActiveMemory/releases

# Linux/macOS example:
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/amem-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
chmod +x amem-*
sudo mv amem-* /usr/local/bin/amem

# Initialize context directory
amem init

# Check status
amem status

# Load context (what AI sees)
amem load

# Detect stale context
amem drift

# Get AI-ready context packet
amem agent
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
