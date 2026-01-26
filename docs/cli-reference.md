---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: CLI Reference
icon: lucide/terminal
---

![ctx](images/ctx-banner.png)

## `ctx` CLI

This is a complete reference for all `ctx` commands.

## Global Options

All commands support these flags:

| Flag                   | Description                                      |
|------------------------|--------------------------------------------------|
| `--help`               | Show command help                                |
| `--version`            | Show version                                     |
| `--context-dir <path>` | Override context directory (default: `.context`) |
| `--quiet`              | Suppress non-essential output                    |
| `--no-color`           | Disable colored output                           |

## Commands

### `ctx init`

Initialize a new `.context/` directory with template files.

```bash
ctx init [flags]
```

**Flags**:

| Flag        | Short | Description                                                           |
|-------------|-------|-----------------------------------------------------------------------|
| `--force`   | `-f`  | Overwrite existing context files                                      |
| `--minimal` | `-m`  | Only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md) |
| `--merge`   |       | Auto-merge ctx content into existing CLAUDE.md                        |

**Creates**:

- `.context/` directory with all template files
- `.claude/hooks/` with auto-save script (for Claude Code)
- `.claude/settings.local.json` with hook configuration
- `CLAUDE.md` with bootstrap instructions (or merges into existing)

**Example**:

```bash
# Standard initialization
ctx init

# Minimal setup (just core files)
ctx init --minimal

# Force overwrite existing
ctx init --force

# Merge into existing CLAUDE.md
ctx init --merge
```

---

### `ctx status`

Show the current context summary.

```bash
ctx status [flags]
```

**Flags**:

| Flag        | Short | Description                   |
|-------------|-------|-------------------------------|
| `--json`    |       | Output as JSON                |
| `--verbose` | `-v`  | Include file contents summary |

**Output**:

- Context directory path
- Total files and token estimate
- Status of each file (*loaded, empty, missing*)
- Recent activity (*modification times*)
- Drift warnings if any

**Example**:

```bash
ctx status
ctx status --json
ctx status --verbose
```

---

### `ctx agent`

Print an AI-ready context packet optimized for LLM consumption.

```bash
ctx agent [flags]
```

**Flags**:

| Flag                | Description                  |
|---------------------|------------------------------|
| `--budget <tokens>` | Token budget (default: 8000) |
| `--format md\|json` | Output format (default: md)  |

**Output**:

- Read order for context files
- Constitution rules (never truncated)
- Current tasks
- Key conventions
- Recent decisions

**Example**:

```bash
# Default (8000 tokens, markdown)
ctx agent

# Custom budget
ctx agent --budget 4000

# JSON format
ctx agent --format json
```

**Use case**: Copy-paste into AI chat, pipe to system prompt, or use in hooks.

---

### `ctx load`

Load and display assembled context as AI would see it.

```bash
ctx load [flags]
```

**Flags**:

| Flag                | Description                               |
|---------------------|-------------------------------------------|
| `--budget <tokens>` | Token budget for assembly (default: 8000) |
| `--raw`             | Output raw file contents without assembly |

**Example**:

```bash
ctx load
ctx load --budget 16000
ctx load --raw
```

---

### `ctx add`

Add a new item to a context file.

```bash
ctx add <type> <content> [flags]
```

**Types**:

| Type         | Target File    |
|--------------|----------------|
| `task`       | TASKS.md       |
| `decision`   | DECISIONS.md   |
| `learning`   | LEARNINGS.md   |
| `convention` | CONVENTIONS.md |

**Flags**:

| Flag                 | Description                                 |
|----------------------|---------------------------------------------|
| `--priority <level>` | Priority for tasks: `high`, `medium`, `low` |
| `--section <name>`   | Target section within file                  |
| `--edit`             | Open editor for full entry                  |

**Examples**:

```bash
# Add a task
ctx add task "Implement user authentication"
ctx add task "Fix login bug" --priority high

# Record a decision
ctx add decision "Use PostgreSQL for primary database"

# Note a learning
ctx add learning "Vitest mocks must be hoisted"

# Add to specific section
ctx add learning "Always use --no-gpg-sign" --section "Git"
```

---

### `ctx complete`

Mark a task as completed.

```bash
ctx complete <task-id-or-text>
```

**Arguments**:

- `task-id-or-text`: Task number or partial text match

**Examples**:

```bash
# By text (partial match)
ctx complete "user auth"

# By task number
ctx complete 3
```

---

### `ctx drift`

Detect stale or invalid context.

```bash
ctx drift [flags]
```

**Flags**:

| Flag     | Description                  |
|----------|------------------------------|
| `--json` | Output machine-readable JSON |
| `--fix`  | Auto-fix simple issues       |

**Checks**:

- Path references in ARCHITECTURE.md and CONVENTIONS.md exist
- Task references are valid
- Constitution rules aren't violated (*heuristic*)
- Staleness indicators (*old files, many completed tasks*)

**Example**:

```bash
ctx drift
ctx drift --json
ctx drift --fix
```

**Exit codes**:

| Code | Meaning           |
|------|-------------------|
| 0    | All checks passed |
| 1    | Warnings found    |
| 2    | Violations found  |

---

### `ctx sync`

Reconcile context with the current codebase state.

```bash
ctx sync [flags]
```

**Flags**:

| Flag        | Description                              |
|-------------|------------------------------------------|
| `--dry-run` | Show what would change without modifying |

**What it does:**

* Scans codebase for structural changes
* Compares with ARCHITECTURE.md
* Suggests documenting dependencies if package files exist
* Identifies stale or outdated context

**Example**:

```bash
ctx sync
ctx sync --dry-run
```

---

### `ctx compact`

Consolidate and clean up context files.

* Moves completed tasks older than 7 days to the archive
* Deduplicates the "*learning*"s with similar content
* Removes empty sections

```bash
ctx compact [flags]
```

**Flags**:

| Flag        | Description                                |
|-------------|--------------------------------------------|
| `--archive` | Create `.context/archive/` for old content |

**Example**:

```bash
ctx compact
ctx compact --archive
```

---

### `ctx watch`

Watch for AI output and auto-apply context updates.

Parses `<context-update>` XML commands from AI output and applies
them to context files.

```bash
ctx watch [flags]
```

**Flags**:

| Flag           | Description                         |
|----------------|-------------------------------------|
| `--log <file>` | Log file to watch (default: stdin)  |
| `--dry-run`    | Preview updates without applying    |
| `--auto-save`  | Periodically save session snapshots |

**Example**:

```bash
# Watch stdin
ai-tool | ctx watch

# Watch a log file
ctx watch --log /path/to/ai-output.log

# Preview without applying
ctx watch --dry-run
```

---

### `ctx hook`

Generate AI tool integration configuration.

```bash
ctx hook <tool>
```

**Supported tools**:

| Tool          | Description     |
|---------------|-----------------|
| `claude-code` | Claude Code CLI |
| `cursor`      | Cursor IDE      |
| `aider`       | Aider CLI       |
| `copilot`     | GitHub Copilot  |
| `windsurf`    | Windsurf IDE    |

**Example**:

```bash
ctx hook claude-code
ctx hook cursor
ctx hook aider
```

---

### `ctx session`

Manage session snapshots.

#### ctx session save

Save the current context snapshot.

```bash
ctx session save [topic] [flags]
```

**Flags**:

| Flag            | Description                                              |
|-----------------|----------------------------------------------------------|
| `--type <type>` | Session type: `feature`, `bugfix`, `refactor`, `session` |

**Example**:

```bash
ctx session save
ctx session save "feature-auth"
ctx session save "bugfix" --type bugfix
```

#### `ctx session list`

List saved sessions.

```bash
ctx session list
```

**Output**: Table of sessions with index, date, topic, and type.

#### `ctx session load`

Load and display a previous session.

```bash
ctx session load <index|date|topic>
```

**Arguments**:

* `index`: Numeric index from `session list`
* `date`: Date pattern (e.g., `2026-01-21`)
* `topic`: Topic keyword match

**Example:**

```bash
ctx session load 1           # by index
ctx session load 2026-01-21  # by date
ctx session load auth        # by topic
```

#### `ctx session parse`

Parse JSONL transcript to readable markdown.

```bash
ctx session parse <file> [flags]
```

**Flags:**

| Flag        | Description                                     |
|-------------|-------------------------------------------------|
| `--extract` | Extract decisions and learnings from transcript |

**Example**:

```bash
ctx session parse ~/.claude/projects/.../transcript.jsonl
ctx session parse transcript.jsonl --extract
```

---

## Exit Codes

| Code | Meaning              |
|------|----------------------|
| 0    | Success              |
| 1    | General error        |
| 2    | Context not found    |
| 3    | Invalid arguments    |
| 4    | File operation error |

## Environment Variables

| Variable   | Description                             |
|------------|-----------------------------------------|
| `CTX_DIR`  | Override default context directory path |
| `NO_COLOR` | Disable colored output when set         |

## Configuration File

Optional `.contextrc` or `context.config.js` at project root:

```javascript
module.exports = {
  contextDir: '.context', // Context directory name
  tokenBudget: 8000,      // Default token budget
  priorityOrder: [        // File loading priority
    'TASKS.md',
    'DECISIONS.md',
    'CONVENTIONS.md',
  ],
  autoArchive: true,   // Auto-archive old items
  archiveAfterDays: 7, // Days before archiving
};
```
