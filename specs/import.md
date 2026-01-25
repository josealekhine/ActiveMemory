# Tier 2: Import Pipeline

> Extract decisions, learnings, and tasks from sessions into ctx format.

## Goals

- **Backfill**: Populate `.context/` from historical sessions
- **Structure**: Transform conversation into ctx artifacts
- **Deduplicate**: Avoid redundant entries across sessions

## CLI

```bash
ctx recall import ./sessions
ctx recall import ./sessions --project myapp --since 2026-01-01
ctx recall import ./sessions --dry-run --llm-refine
```

## Architecture

```
Sessions → Deterministic → LLM Refine → Dedup → ctx Format
           Extractor       (optional)          Writer
```

## Extraction Targets

### Decisions

Choices about implementation with reasoning.

**Indicators**:
- "Let's use X instead of Y"
- "Going with X because..."
- "The approach will be..."

**Output**: `.context/decisions/YYYY-MM-DD-slug.md`

```markdown
---
date: 2026-01-21
source: session://async-roaming-allen
tags: [parser, error-handling]
---

# Use Error Wrapping Instead of Custom Types

Use Go's `%w` error wrapping instead of custom error types.
Standard library support, better composition, easier testing.
```

### Learnings

Discoveries from debugging or research.

**Indicators**:
- "TIL:", "Learned that..."
- "The issue was...", "Fixed by..."
- "Turns out..."

**Output**: `.context/learnings/YYYY-MM-DD-slug.md`

```markdown
---
date: 2026-01-21
source: session://async-roaming-allen
tags: [go, regex]
---

# Regex Negative Lookbehind Not Supported in Go

Go's regexp uses RE2 which lacks negative lookbehind.
Use `[^\\]"` instead of `(?<!\\)"` for unescaped quotes.
```

### Tasks

Work items not completed in session.

**Indicators**:
- "TODO:", "FIXME:"
- "Next steps:", "Still need to..."

**Output**: `.context/tasks/YYYY-MM-DD-slug.md`

```markdown
---
date: 2026-01-21
source: session://async-roaming-allen
status: open
---

# Add Edge Case Tests for Escaped Quotes

Cover: empty strings, nested escaping, unicode mixed with quotes.
```

## Extraction Pipeline

### Phase 1: Deterministic (High Precision)

```go
var patterns = []Pattern{
    {`(?i)let's use (\w+) instead of (\w+)`, "decision"},
    {`(?i)going with (.+?) because`, "decision"},
    {`(?i)TIL[:\s]+(.+?)(?:\.|$)`, "learning"},
    {`(?i)the (?:issue|problem) was (.+?)(?:\.|$)`, "learning"},
    {`(?i)TODO[:\s]+(.+?)(?:\n|$)`, "task"},
}

func Extract(session *Session) []Extraction {
    var results []Extraction
    for _, msg := range session.Messages {
        for _, block := range msg.Content {
            // Extract from text AND thinking blocks
            text := block.Text()
            for _, p := range patterns {
                if matches := p.Regex.FindAllStringSubmatch(text, -1); matches != nil {
                    results = append(results, buildExtraction(p, matches, msg))
                }
            }
        }
    }
    return results
}
```

### Phase 2: LLM Refinement (Optional)

```go
const prompt = `Analyze this session excerpt. Extract:
1. DECISIONS: Implementation/architecture choices
2. LEARNINGS: Debugging discoveries
3. TASKS: Incomplete work items

Return JSON array with: category, title, content, confidence (0-1), tags`

func Refine(session *Session, existing []Extraction) ([]Extraction, error) {
    // Batch into 8k token chunks
    // Call Claude API with structured output
    // Filter by confidence threshold (default: 0.7)
}
```

### Phase 3: Deduplication

```go
func Deduplicate(extractions []Extraction, threshold float64) []Extraction {
    // 1. Embed all extractions
    // 2. Find pairs with similarity > threshold
    // 3. Cluster duplicates
    // 4. Select canonical (highest confidence, most recent)
    // 5. Track evolution (same topic, different content over time)
}
```

## Output Structure

```
.context/
├── decisions/
│   ├── 2026-01-20-use-error-wrapping.md
│   └── index.json
├── learnings/
│   ├── 2026-01-20-regex-lookbehind-go.md
│   └── index.json
├── tasks/
│   └── 2026-01-21-edge-case-tests.md
└── import-log.json
```

## Tasks

| Phase | Task               | Hours |
|-------|--------------------|-------|
| 2.1   | Pattern extraction | 8     |
| 2.2   | LLM refinement     | 8     |
| 2.3   | ctx format writer  | 8     |
| 2.4   | Deduplication      | 8     |
| 2.5   | CLI integration    | 8     |

## Success Metrics

| Metric               | Target           |
|----------------------|------------------|
| Extraction precision | > 80%            |
| Recall (with LLM)    | > 70%            |
| Duplicate detection  | > 90%            |
| Import 100 sessions  | < 5 min (no LLM) |
