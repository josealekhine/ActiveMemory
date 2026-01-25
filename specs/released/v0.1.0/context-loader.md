# Context Loader Specification

## Overview

The Context Loader discovers, reads, parses, and assembles context files into a 
unified context payload suitable for AI consumption. It's the read path of the 
Context system.

## Responsibilities

1. **Discovery** — Find the `.context/` directory and enumerate files
2. **Reading** — Load file contents with proper encoding
3. **Parsing** — Convert markdown to structured data
4. **Assembly** — Combine parsed files into unified context
5. **Budgeting** — Respect token limits and prioritize content

## API

### Core Types

```go
type LoaderOptions struct {
    RootDir       string   // Project root (default: cwd)
    ContextDir    string   // Context directory name (default: ".context")
    TokenBudget   int      // Max tokens to emit (default: 8000)
    PriorityOrder []string // File loading priority
}

type LoadedContext struct {
    Files      []ContextFile // Parsed context files
    Summary    string        // Assembled context for AI
    TokenCount int           // Estimated tokens used
    Truncated  bool          // Whether content was truncated
    Missing    []string      // Expected but missing files
}

type ContextFile struct {
    Name       string        // File name without extension
    Path       string        // Full path
    Content    string        // Raw markdown content
    Parsed     ParsedContent // Structured data
    TokenCount int           // Estimated tokens for this file
}

// Main entry point
func LoadContext(opts *LoaderOptions) (*LoadedContext, error)

// Individual file loading
func LoadContextFile(path string) (*ContextFile, error)

// Check if context directory exists
func HasContext(rootDir string) (bool, error)
```

## Discovery Rules

1. Start from `rootDir` (default: current working directory)
2. Look for `.context/` directory
3. If not found, walk up parent directories (max 5 levels)
4. Stop at filesystem root or `.git` boundary
5. Return `missing` list for expected-but-absent files

## Priority Order

Default loading priority (highest first):

1. `CONSTITUTION.md` — Hard rules; must never be violated
2. `TASKS.md` — Most immediately actionable
3. `DECISIONS.md` — Prevents re-debating
4. `CONVENTIONS.md` — Ensures consistency
5. `ARCHITECTURE.md` — Provides orientation
6. `GLOSSARY.md` — Correct terminology
7. `LEARNINGS.md` — Avoids repeated mistakes
8. `DEPENDENCIES.md` — Reference material
9. `DRIFT.md` — Maintenance guidance
10. `AGENT_PLAYBOOK.md` — Agent instructions (often already known)

When token budget is tight, lower-priority files get truncated or omitted.
CONSTITUTION.md is NEVER truncated.

## Parsing

### Markdown to Structure

The parser extracts structure from markdown using these rules:

```go
type ParsedContent struct {
    Title    string            // From first H1
    Sections []Section         // From H2 headers
    Metadata map[string]string // From frontmatter (optional)
}

type Section struct {
    Title   string  // Section header text
    Items   []Item  // Parsed items within section
    Content string  // Raw content if not itemized
}

type Item struct {
    Title    string            // Item header (H3) or first line
    Fields   map[string]string // **Key**: Value pairs
    Tags     []string          // Inline tags like #priority:high
    Status   string            // "todo" | "done" | "blocked" (from checkbox)
    Children []string          // Sub-bullets
}
```

### Error Handling

- **Malformed markdown**: Parse what's possible, log warnings
- **Missing files**: Include in `missing` array, continue loading others
- **Encoding issues**: Assume UTF-8, replace invalid sequences
- **Empty files**: Include with empty parsed content

## Token Estimation

Use simple heuristic: 1 token ≈ 4 characters for English text.

More accurate estimation can be added later with tiktoken or similar.

## Assembly

The `summary` field contains the assembled context as a single markdown string:

```markdown
# Project Context

## Current Tasks
[Summarized from TASKS.md]

## Key Decisions
[Recent/relevant from DECISIONS.md]

## Conventions
[From CONVENTIONS.md]

## Architecture Overview
[From ARCHITECTURE.md]

## Recent Learnings
[From LEARNINGS.md]
```

### Assembly Rules

1. Each section has a header indicating source
2. Content is included verbatim when under budget
3. When over budget, summarize or truncate with `[truncated]` marker
4. Always include at least the headers from each file
5. Preserve code blocks exactly

## Caching

Optional caching layer:

- Cache parsed files by path + mtime
- Invalidate on file change
- Cache assembled summary by combined hash
- No cache by default (simplicity over speed)

## Testing Requirements

- Unit tests for each parsing rule
- Integration tests with sample `.context/` directories
- Edge cases: empty files, malformed markdown, missing directory
- Token budget tests: verify truncation behavior
- Performance: loading should complete in <100ms for typical projects
