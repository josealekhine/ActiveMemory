# Context Updater Specification

## Overview

The Context Updater handles the write path of the Context system. It parses
structured update commands from AI responses and writes them back to the
appropriate context files while maintaining format consistency.

## Responsibilities

1. **Detection** — Identify update commands in AI output
2. **Parsing** — Extract structured updates from commands
3. **Validation** — Verify updates are well-formed
4. **Writing** — Apply updates to context files
5. **Formatting** — Maintain consistent file structure

## Update Command Format

AI outputs updates using a simple, parseable format:

```markdown
<context-update file="DECISIONS.md" action="add" section="Decisions">
## [2025-01-19] Use Vitest for Testing

**Status**: Accepted

**Context**: Need to choose a test framework for the project.

**Decision**: Use Vitest with native ESM support.

**Rationale**: Faster than Jest, native TypeScript support, compatible with Vite.
</context-update>
```

### Command Attributes

- `file` (required): Target context file name
- `action` (required): `add` | `update` | `remove` | `complete`
- `section` (optional): Target section within file
- `id` (optional): Identifier for update/remove operations

### Actions

| Action     | Behavior                                          |
|------------|---------------------------------------------------|
| `add`      | Append content to section (or file if no section) |
| `update`   | Replace content matching `id`                     |
| `remove`   | Delete content matching `id`                      |
| `complete` | Mark task as done (checkbox toggle)               |

## API

```go
type UpdateCommand struct {
    File    string // Target context file name
    Action  string // "add" | "update" | "remove" | "complete"
    Section string // Target section within file (optional)
    ID      string // Identifier for update/remove operations (optional)
    Content string // Update content
}

type UpdateResult struct {
    Success bool
    File    string
    Action  string
    Error   string // Empty if success
    Diff    string // For preview mode
}

type UpdaterOptions struct {
    ContextDir string // Path to .context/ directory
    DryRun     bool   // Preview without writing
    Backup     bool   // Create .bak before writing
    Validate   bool   // Strict validation mode
}

// Parse update commands from text
func ParseUpdates(text string) ([]UpdateCommand, error)

// Apply a single update
func ApplyUpdate(cmd UpdateCommand, opts *UpdaterOptions) (*UpdateResult, error)

// Apply all updates from text
func ApplyUpdates(text string, opts *UpdaterOptions) ([]UpdateResult, error)

// Validate an update command without applying
func ValidateUpdate(cmd UpdateCommand) (bool, []string)
```

## Detection Rules

Update commands are detected by:

1. Opening tag: `<context-update` with attributes
2. Closing tag: `</context-update>`
3. Content between tags is the update payload

Multiple commands can appear in a single AI response.

## Writing Rules

### Add Action

1. Find target file
2. Find target section (if specified)
3. Append content at end of section
4. Add blank line separator if needed

### Update Action

1. Find target file
2. Find content matching `id` (by header text or first line)
3. Replace entire block (from header to next same-level header)
4. Preserve surrounding content

### Remove Action

1. Find target file
2. Find content matching `id`
3. Remove entire block
4. Clean up extra blank lines

### Complete Action

1. Find target file (typically TASKS.md)
2. Find task matching `id` (by task text)
3. Change `- [ ]` to `- [x]`
4. Optionally move to "Completed" section with date

## Format Preservation

When writing:

1. Preserve file encoding (UTF-8)
2. Preserve line endings (detect and match)
3. Preserve indentation style
4. Maintain consistent blank lines between sections
5. Keep code blocks intact

## Validation

Before writing, validate:

1. Target file exists (or can be created)
2. Action is valid
3. Content is well-formed markdown
4. Section exists (if specified)
5. ID can be found (for update/remove/complete)

## Error Handling

| Error             | Behavior                           |
|-------------------|------------------------------------|
| File not found    | Create if `add`, error otherwise   |
| Section not found | Create if `add`, error otherwise   |
| ID not found      | Error with suggestions             |
| Parse error       | Return partial results with errors |
| Write error       | Restore from backup if available   |

## Conflict Resolution

When multiple updates target the same location:

1. Apply in order received
2. Later updates see earlier changes
3. Log warnings for potential conflicts

## Backup Strategy

When `backup: true`:

1. Copy `FILE.md` to `FILE.md.bak` before first write
2. Keep only one backup (overwrite previous)
3. Backup is same-directory, not in `.context/.backups/`

## Integration Points

### For AI Tools

AI tools should:

1. Include update commands in responses when context changes
2. Use appropriate action for the change type
3. Include enough content for unambiguous matching

### For CLI

The CLI can:

1. Extract and apply updates from AI output logs
2. Show diffs before applying (dry run)
3. Batch updates from multiple sources

## Testing Requirements

- Unit tests for command parsing
- Unit tests for each action type
- Integration tests with real files
- Edge cases: empty files, malformed commands, missing sections
- Conflict tests: multiple updates to same content
- Format preservation tests: round-trip editing
