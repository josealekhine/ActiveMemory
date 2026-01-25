# CLI Specification

## Overview

The Context CLI provides commands for initializing, managing, and inspecting 
project context. It's the primary human interface to the system.

## Implementation

- **Language**: Go (minimal dependencies, single binary)
- **Distribution**: GitHub Releases (pre-built binaries for Linux, macOS, Windows)
- **Repository**: https://github.com/ActiveMemory/ctx

## Installation

```bash
# Download the latest release for your platform
# From: https://github.com/ActiveMemory/ctx/releases

# Linux (amd64)
curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-linux-amd64
chmod +x ctx-linux-amd64
sudo mv ctx-linux-amd64 /usr/local/bin/ctx

# Linux (arm64)
curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-linux-arm64
chmod +x ctx-linux-arm64
sudo mv ctx-linux-arm64 /usr/local/bin/ctx

# macOS (Apple Silicon)
curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-darwin-arm64
chmod +x ctx-darwin-arm64
sudo mv ctx-darwin-arm64 /usr/local/bin/ctx

# macOS (Intel)
curl -LO https://github.com/ActiveMemory/ctx/releases/latest/download/ctx-darwin-amd64
chmod +x ctx-darwin-amd64
sudo mv ctx-darwin-amd64 /usr/local/bin/ctx

# Windows (PowerShell)
# Download ctx-windows-amd64.exe from releases page
# Add to PATH or move to a directory in PATH

# Verify installation
ctx --version
```

## Usage

```bash
ctx <command> [options]
```

## Commands

### `ctx init`

Initialize a new `.context/` directory with template files.

```bash
ctx init [--force] [--minimal]
```

**Behavior**:
1. Check if `.context/` already exists
2. If exists and no `--force`, prompt for confirmation
3. Create `.context/` directory
4. Create template files from templates (or minimal set)
5. Add `.context/` section to `.gitignore` comments (not ignoring)

**Options**:
- `--force`: Overwrite existing context files
- `--minimal`: Only create essential files (TASKS.md, DECISIONS.md)

**Templates**: See `templates/` directory for file templates

---

### `ctx status`

Show current context summary.

```bash
ctx status [--json] [--verbose]
```

**Output**:
```
Context Status
====================

📁 Context Directory: .context/
📊 Total Files: 9
📝 Token Estimate: 3,850 tokens

Files:
  ✓ CONSTITUTION.md (5 invariants)
  ✓ TASKS.md (3 active, 2 completed)
  ✓ DECISIONS.md (5 decisions)
  ✓ CONVENTIONS.md (loaded)
  ✓ ARCHITECTURE.md (loaded)
  ✓ GLOSSARY.md (12 terms)
  ✓ DRIFT.md (loaded)
  ✓ AGENT_PLAYBOOK.md (loaded)
  ✗ LEARNINGS.md (empty)

Recent Activity:
  - TASKS.md modified 2 hours ago
  - DECISIONS.md modified 1 day ago

Drift Status: ⚠️ 2 warnings (run 'ctx drift' for details)
```

**Options**:
- `--json`: Output as JSON
- `--verbose`: Include file contents summary

---

### `ctx agent`

Print a concise "agent context packet" for AI consumption.

```bash
ctx agent [--budget <tokens>] [--format md|json]
```

**Behavior**:
1. Assemble context in optimal read order (per AGENT_PLAYBOOK.md)
2. Respect token budget
3. Output ready-to-paste context for AI tools

**Output**:
```markdown
# Context Context Packet
Generated: 2025-01-19T10:30:00Z | Budget: 8000 tokens | Used: 3850

## Read These Files (in order)
1. .context/CONSTITUTION.md
2. .context/TASKS.md
3. .context/CONVENTIONS.md
4. .context/ARCHITECTURE.md
5. .context/DECISIONS.md

## Constitution (NEVER VIOLATE)
- Never commit secrets
- All code must pass tests
- No TODO comments in main

## Current Tasks
- [ ] Implement context loader #priority:high
- [ ] Add drift detection #priority:medium

## Key Conventions
[Summarized from CONVENTIONS.md]

## Recent Decisions
[Last 3 from DECISIONS.md]
```

**Options**:
- `--budget <tokens>`: Token budget (default: 8000)
- `--format md|json`: Output format (default: md)

**Use Case**: Copy-paste into AI chat, or pipe to system prompt.

---

### `ctx load`

Load and display assembled context (what AI sees).

```bash
ctx load [--budget <tokens>] [--raw]
```

**Output**: The assembled markdown context as it would be provided to an AI.

**Options**:
- `--budget <tokens>`: Token budget for assembly (default: 8000)
- `--raw`: Output raw file contents without assembly

---

### `ctx sync`

Reconcile context with current codebase state.

```bash
ctx sync [--dry-run]
```

**Behavior**:
1. Scan codebase for structural changes
2. Compare with ARCHITECTURE.md
3. Check DEPENDENCIES.md against package.json
4. Identify stale or outdated context
5. Prompt for updates (or show diff in dry-run)

**Options**:
- `--dry-run`: Show what would change without modifying

---

### `ctx compact`

Consolidate and clean up context files.

```bash
ctx compact [--archive]
```

**Behavior**:
1. Move completed tasks older than 7 days to archive
2. Deduplicate learnings with similar content
3. Remove empty sections
4. Consolidate related decisions (prompt for confirmation)

**Options**:
- `--archive`: Create `.context/archive/` for old content

---

### `ctx drift`

Detect stale or invalid context (drift detection).

```bash
ctx drift [--json] [--fix]
```

**Behavior**:
1. Check all path references in ARCHITECTURE.md, CONVENTIONS.md exist
2. Check task references in TASKS.md are valid
3. Check CONSTITUTION.md rules aren't violated (heuristic)
4. Check for staleness indicators (old files, many completed tasks)
5. Output human-readable report

**Output Example**:
```
Drift Detection Report
======================

⚠️  WARNINGS (3)

  Path References:
  - ARCHITECTURE.md:42 references 'src/utils/helpers.ts' (not found)
  - CONVENTIONS.md:18 references 'scripts/lint.sh' (not found)

  Staleness:
  - TASKS.md has 15 completed items (consider archiving)

❌ VIOLATIONS (1)

  Constitution:
  - Found potential secret pattern in 'config/dev.env' (rule: no secrets)

✅ PASSED (4)
  - GLOSSARY.md terms are consistent
  - DEPENDENCIES.md packages exist in package.json
  - DECISIONS.md references are valid
  - LEARNINGS.md is under size limit
```

**Options**:
- `--json`: Output machine-readable JSON report
- `--fix`: Auto-fix simple issues (remove dead refs, archive old tasks)

**JSON Output**:
```json
{
  "timestamp": "2025-01-19T10:30:00Z",
  "status": "warning",
  "warnings": [
    {"file": "ARCHITECTURE.md", "line": 42, "type": "dead_path", "path": "src/utils/helpers.ts"}
  ],
  "violations": [
    {"file": "config/dev.env", "type": "potential_secret", "rule": "no_secrets"}
  ],
  "passed": ["glossary_consistency", "dependency_check", "decision_refs", "learnings_size"]
}
```

---

### `ctx add`

Add a new item to a context file.

```bash
ctx add <file> <content>
ctx add decision "Use PostgreSQL for primary database"
ctx add task "Implement user authentication" --priority high
ctx add learning "Vitest mocks must be hoisted"
```

**Arguments**:
- `file`: Target file (decision, task, learning, convention)
- `content`: Item content (quoted string)

**Options**:
- `--priority`: Priority level for tasks (high, medium, low)
- `--section`: Target section within file
- `--edit`: Open editor for full item entry

---

### `ctx complete`

Mark a task as completed.

```bash
ctx complete <task-id-or-text>
ctx complete "Implement user authentication"
ctx complete 3  # By task number
```

---

### `ctx watch`

Watch for AI output and auto-apply context updates.

```bash
ctx watch [--log <file>] [--dry-run]
```

**Behavior**:
1. Watch specified log file (or stdin)
2. Parse for `<context-update>` commands
3. Apply updates to context files
4. Show confirmation for each update

**Options**:
- `--log <file>`: Log file to watch (default: stdin)
- `--dry-run`: Show updates without applying

---

### `ctx hook`

Generate AI tool integration hooks.

```bash
ctx hook <tool>
ctx hook claude-code
ctx hook cursor
ctx hook aider
```

**Output**: Instructions and configuration for integrating with specified AI tool.

---

## Global Options

All commands support:

- `--help`: Show command help
- `--version`: Show version
- `--context-dir <path>`: Override context directory path
- `--quiet`: Suppress non-essential output
- `--no-color`: Disable colored output

## Configuration

Optional `.contextrc` or `context.config.js` at project root:

```javascript
module.exports = {
  contextDir: '.context',      // Context directory name
  tokenBudget: 8000,           // Default token budget
  priorityOrder: [             // File loading priority
    'TASKS.md',
    'DECISIONS.md',
    'CONVENTIONS.md',
  ],
  autoArchive: true,           // Auto-archive old items
  archiveAfterDays: 7,         // Days before archiving
};
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Context not found |
| 3 | Invalid arguments |
| 4 | File operation error |

## Testing Requirements

- Unit tests for argument parsing
- Integration tests for each command
- E2E tests with real file operations
- Edge cases: missing context, corrupted files, permission errors
