# Core Architecture Specification

## Overview

Context is a file-based context engineering system that persists project 
knowledge across AI coding sessions. It treats context as infrastructure, 
not ephemeral prompt content.

## Design Philosophy

### First Principles

1. **Memory IS the filesystem** — AI agents have no conversational memory; 
   everything important must be written to files
2. **Context is a system, not a prompt** — Context should be structured, 
   versioned, and maintained like code
3. **Decisions compound** — Every architectural choice, pattern selection, 
   and lesson learned should persist
4. **Memory matches workflow** — Context structure mirrors how engineers 
   actually think about projects
5. **Tool-agnostic by design** — Works with any AI tool that can read files 
   (Claude Code, Cursor, Aider, etc.)
6. **Git-native** — Context files are text, diffable, and committable

### Why File-Based?

- **No dependencies** — No database, no daemon, no runtime
- **Version controlled** — Context evolves with code, branches with code, merges with code
- **Human-readable** — Engineers can read, edit, and understand context directly
- **Token-efficient** — Markdown is more token-efficient than JSON/XML
- **Portable** — Works anywhere files work

## Core Components

### 1. Context Directory (`.context/`)

A dedicated directory at project root containing all context files:

```
.context/
├── CONSTITUTION.md     # Hard invariants — rules that must NEVER be violated
├── TASKS.md            # Current and planned work
├── DECISIONS.md        # Architectural decisions (or decisions/ dir for many)
├── LEARNINGS.md        # Lessons learned, gotchas, tips
├── CONVENTIONS.md      # Project patterns and standards
├── ARCHITECTURE.md     # System overview and component relationships
├── DEPENDENCIES.md     # Key dependencies and why they were chosen
├── GLOSSARY.md         # Domain terms, abbreviations, project vocabulary
├── DRIFT.md            # Staleness signals — when to update what
├── AGENT_PLAYBOOK.md   # How AI agents should read/apply/update context
└── sessions/           # Optional: session-specific context
    └── 2025-01-19.md   # Daily session notes (auto-generated)
```

For projects with many architectural decisions, use a `decisions/` directory:

```
.context/decisions/
├── 0001-use-typescript-strict.md
├── 0002-choose-vitest.md
└── 0003-postgresql-primary-db.md
```

### 2. Context Loader

A module that:
- Discovers and reads context files
- Parses markdown into structured data
- Assembles context for AI consumption
- Handles missing files gracefully

### 3. Context Updater

A module that:
- Watches for AI outputs suggesting context updates
- Parses structured update commands from AI responses
- Writes updates back to appropriate context files
- Maintains file format consistency

### 4. CLI Interface

A single Go binary (`ctx`) with commands for human operators:
- `ctx init` — Initialize `.context/` with templates
- `ctx status` — Show current context summary
- `ctx sync` — Reconcile context with codebase state
- `ctx compact` — Consolidate and deduplicate context
- `ctx drift` — Detect stale paths, broken refs, constitution violations
- `ctx agent` — Print AI-ready context packet

**Implementation**: Go with minimal dependencies. Single binary, cross-platform.

**Distribution**: GitHub Releases at https://github.com/ActiveMemory/ctx

## Constraints

1. **No binary files** — All context must be text-based markdown
2. **No external services** — Everything runs locally, offline-capable
3. **No magic** — Explicit is better than implicit; all context loading visible
4. **No lock-in** — If you delete the CLI, the files remain useful
5. **Total context budget** — System must track and respect token limits

## Non-Goals

These are explicitly OUT OF SCOPE:

1. **Not a SaaS** — No server, hosted database, or cloud sync
2. **Not proprietary** — Must work without any vendor integrations
3. **Not a secrets manager** — Never store tokens, credentials, API keys, or customer data
4. **Not a replacement for git** — Context files are git-tracked, not a parallel VCS
5. **Not an IDE plugin** — CLI and files first; IDE integrations are optional extras

## Success Criteria

- AI sessions can cold-start with full project context in <30 seconds
- Decisions made in session N are available in session N+1
- Context can be manually edited without breaking the system
- System works with zero configuration for simple projects
- Context overhead is <5% of available tokens
