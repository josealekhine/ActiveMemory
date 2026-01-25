---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

icon: lucide/rocket
---

![ctx](images/ctx-banner.png)

## `ctx`

`ctx` (*Context*) is a file-based system that enables AI coding assistants to
persist project knowledge across sessions. Instead of re-explaining your
codebase every time, context files let AI tools remember decisions,
conventions, and learnings:

* A session is interactive.
* `ctx` enables **cognitive continuity**.
* **Cognitive continuity** enables durable, *symbiotic-like* human–AI workflows.

## Community

**Open source is better together**.

> ⭐️ **If the idea behind `ctx` resonates, a star helps it reach engineers  
> who run into context drift every day.**
> 
> → https://github.com/ActiveMemory/ctx

Join the community to ask questions, share feedback, and connect with
other users:

- [:fontawesome-brands-discord: **Discord**](https://discord.gg/kampus):
  Chat, get help, share your workflows
- [:fontawesome-brands-github: **GitHub**](https://github.com/ActiveMemory/ctx):
  Star the repo, report issues, contribute

## Why?

Most AI-driven development fails not because models are weak—they fail because 
**context is ephemeral**. Every new session starts near zero:

* You re-explain architecture
* The AI repeats past mistakes
* Decisions get rediscovered instead of remembered

Context solves this by treating context as infrastructure: 
files that version with your code and persist across sessions.

## Installation

### Binary Downloads (Recommended)

Download pre-built binaries from the 
[releases page](https://github.com/ActiveMemory/ctx/releases).

=== "Linux (x86_64)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-linux-amd64
    chmod +x ctx-linux-amd64
    sudo mv ctx-linux-amd64 /usr/local/bin/ctx
    ```

=== "Linux (ARM64)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-linux-arm64
    chmod +x ctx-linux-arm64
    sudo mv ctx-linux-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Apple Silicon)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-darwin-arm64
    chmod +x ctx-darwin-arm64
    sudo mv ctx-darwin-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Intel)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-darwin-amd64
    chmod +x ctx-darwin-amd64
    sudo mv ctx-darwin-amd64 /usr/local/bin/ctx
    ```

=== "Windows"

    Download `ctx-windows-amd64.exe` from the releases page and add it to your `PATH`.

### Build from Source

Requires [Go 1.26+](https://go.dev/):

```bash
git clone https://github.com/ActiveMemory/ctx.git
cd ctx
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
sudo mv ctx /usr/local/bin/
```

Verify installation:

```bash
ctx --version
```

## Quick Start

### 1. Initialize Context

```bash
cd your-project
ctx init
```

This creates a `.context/` directory with template files and configures 
AI tool hooks (*for Claude Code*).

### 2. Check Status

```bash
ctx status
```

Shows context summary: files present, token estimate, and recent activity.

### 3. Start Using with AI

With Claude Code, context loads automatically via hooks. For other tools,
paste the output of:

```bash
ctx agent --budget 8000
```

### 4. Verify It Works

Ask your AI: **"Do you remember?"**

It should cite specific context: current tasks, recent decisions, 
or previous session topics.

## What Gets Created

```
.context/
├── CONSTITUTION.md     # Hard rules — NEVER violate these
├── TASKS.md            # Current and planned work
├── CONVENTIONS.md      # Project patterns and standards
├── ARCHITECTURE.md     # System overview
├── DECISIONS.md        # Architectural decisions with rationale
├── LEARNINGS.md        # Lessons learned, gotchas, tips
├── GLOSSARY.md         # Domain terms and abbreviations
├── DRIFT.md            # Staleness signals
├── AGENT_PLAYBOOK.md   # How AI agents should use this
└── sessions/           # Session snapshots
```

See [Context Files](context-files.md) for detailed documentation of each file.

## Common Workflows

### Add a Task

```bash
ctx add task "Implement user authentication"
```

### Record a Decision

```bash
ctx add decision "Use PostgreSQL for primary database"
```

### Note a Learning

```bash
ctx add learning "Mock functions must be hoisted in Jest"
```

### Mark Task Complete

```bash
ctx complete "user auth"
```

### Check for Stale Context

```bash
ctx drift
```

## Next Steps

- [CLI Reference](cli-reference.md) — All commands and options
- [Context Files](context-files.md) — File formats and structure
- [Ralph Loop Integration](ralph-loop.md) — Autonomous AI development workflows
- [Integrations](integrations.md) — Setup for Claude Code, Cursor, Aider
