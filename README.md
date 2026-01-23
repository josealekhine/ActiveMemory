![ctx](assets/ctx-banner.png)

## `ctx` (*Context*)

> **Context is a system, not a prompt.**

A lightweight, file-based system that enables AI coding assistants to persist,
structure, and rehydrate project context across sessions.

`ctx` works with **any AI tool** that can read files; no model or 
vendor lock-in.

**Full documentation: [ctx.ist](https://ctx.ist)**

## The Problem

Most LLM-driven development fails not because models are weak: They fail because
**context is ephemeral**. Every new session starts near zero:

* You re-explain architecture
* The AI repeats past mistakes
* Decisions get rediscovered instead of remembered

## The Solution

`ctx` treats context as infrastructure:

* **Persist**: Tasks, decisions, learnings survive session boundaries
* **Reuse**: Decisions don't get rediscovered; lessons stay learned
* **Align**: Context structure mirrors how engineers actually think
* **Integrate**: Works with any AI tool that can read files

Here's what that looks like in practice:

```text
❯ "Do you remember?"

● Yes. The PreToolUse hook runs ctx agent, and CLAUDE.md tells me to
  check .context/sessions/. I'll have context.

❯ "Summarize all sessions we have had so far?"

● Yes. I can ls .context/sessions/ and read each file:
    - 2025-01-20: The meta-experiment that started it all
    - 2025-01-21: The ctx rename + Claude hooks session
```

That's the whole point: **Temporal continuity across sessions**.

## Installation

Download pre-built binaries from the
[releases page](https://github.com/ActiveMemory/ctx/releases), or build from
source:

```bash
git clone https://github.com/ActiveMemory/ctx.git
cd ctx
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
sudo mv ctx /usr/local/bin/
```

See [installation docs](https://ctx.ist/#installation) for platform-specific
instructions.

## Quick Start

```bash
# Initialize context directory in your project
ctx init

# Check context status
ctx status

# Get an AI-ready context packet
ctx agent --budget 4000

# Add tasks, decisions, learnings
ctx add task "Implement user authentication"
ctx add decision "Use PostgreSQL for primary database"
ctx add learning "Mock functions must be hoisted in Jest"
```

## Documentation

| Guide                                           | Description                            |
|-------------------------------------------------|----------------------------------------|
| [Getting Started](https://ctx.ist)              | Installation, quick start, first steps |
| [CLI Reference](https://ctx.ist/cli-reference/) | All commands and options               |
| [Context Files](https://ctx.ist/context-files/) | File formats and structure             |
| [Integrations](https://ctx.ist/integrations/)   | Claude Code, Cursor, Aider setup       |
| [Ralph Loop](https://ctx.ist/ralph-loop/)       | Autonomous AI development workflows    |

## Design Philosophy

1. **File-based**: No database, no daemon. Just markdown and convention.
2. **Git-native**: Context versions with code, branches with code, merges with
   code.
3. **Human-readable**: Engineers can read, edit, and understand context
   directly.
4. **Token-efficient**: Markdown is cheaper than JSON/XML.
5. **Tool-agnostic**: Works with Claude Code, Cursor, Aider, Copilot, or raw
   CLI.

## Contributing

Contributions welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

All commits must be signed off (`git commit -s`) to certify the
[DCO](CONTRIBUTING_DCO.md).

## License

[Apache 2.0](LICENSE)
