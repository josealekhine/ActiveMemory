# Context:Recall Implementation Plan

This is the next phase of `ctx`.
* Check ./specs/released/v0.1.0/ for the alread-implemented features.

> Transform ephemeral AI session data into persistent, searchable, 
> actionable context.

## Overview

`ctx recall` enables ctx to learn from its own history. It processes session 
files from AI coding assistants (Claude Code, etc.) and transforms them into 
retrievable knowledge.

**Design principle**: `ctx` provides forward-looking structure 
(tasks, decisions, learnings). Recall provides backward-looking retrieval 
(what happened, what worked, what failed).

## Tiers

| Tier | Name                            | Value                                        | Effort   |
|------|---------------------------------|----------------------------------------------|----------|
| 1    | [Browsable](specs/browsable.md) | Human visibility into session history        | 1-2 days |
| 2    | [Import](specs/import.md)       | Backfill ctx from historical sessions        | 3-5 days |
| 3    | [Thinking](specs/thinking.md)   | Extract reasoning patterns and anti-patterns | 2-3 days |
| 4    | [RAG](specs/rag.md)             | Query history at session start               | 3-5 days |

---

## Input Format

Session files are JSONL (one JSON object per line). See 
[session-schema.md](specs/session-schema.md) for full specification.

```json
{
  "sessionId": "af7cba21-...",
  "timestamp": "2026-01-21T07:51:32.271Z",
  "cwd": "/home/user/project",
  "gitBranch": "main",
  "message": {
    "content": [
      {"type": "thinking", "thinking": "..."},
      {"type": "text", "text": "..."}
    ],
    "usage": {"input_tokens": 8, "output_tokens": 1}
  }
}
```

---

## CLI Integration

```bash
# Tier 1: Browse sessions
ctx recall serve ./sessions

# Tier 2: Import to ctx format
ctx recall import ./sessions --project myapp

# Tier 3: Search reasoning patterns
ctx recall reason "how to debug connection refused"

# Tier 4: Auto-context injection
ctx recall --auto --budget 4000
```

---

## Tier 1: Browsable Archive

**Goal**: Human-readable HTML rendering of session history.

### Deliverables
- [ ] Session parser (Go)
- [ ] Markdown/HTML renderer
- [ ] Local HTTP server with index
- [ ] Search/filter by date, project, branch

### Architecture

```
Session Files → Parser → Rendered HTML → HTTP Server
   (JSONL)                                localhost:8080
```

### Key Features
- Collapsible thinking blocks
- Syntax-highlighted code blocks
- Metadata sidebar (tokens, branch, cwd)
- Full-text search

**Details**: [browsable.md](specs/browsable.md)

---

## Tier 2: ctx Import Pipeline

**Goal**: Extract decisions, learnings, and tasks into ctx format.

### Deliverables
- [ ] Pattern-based extraction (deterministic)
- [ ] LLM refinement pass (optional)
- [ ] Deduplication across sessions
- [ ] ctx-compatible markdown output

### Extraction Targets

| Type     | Indicators                            | Example               |
|----------|---------------------------------------|-----------------------|
| Decision | "Let's use X", "Going with Y because" | Architecture choices  |
| Learning | "TIL", "The issue was", "Fixed by"    | Debugging discoveries |
| Task     | "TODO", "Next steps", "Still need to" | Incomplete work       |

### Output

```
.context/
├── decisions/2026-01-21-use-error-wrapping.md
├── learnings/2026-01-21-regex-lookbehind-go.md
└── tasks/2026-01-21-edge-case-tests.md
```

**Details**: [import.md](specs/import.md)

---

## Tier 3: Thinking Mining

**Goal**: Extract and index reasoning patterns from thinking blocks.

### Deliverables
- [ ] Thinking block parser
- [ ] Pattern classifier (decomposition, hypothesis, pivot, error analysis)
- [ ] Outcome detector (success/failure/abandoned)
- [ ] Searchable reasoning corpus

### Categories

| Category       | Pattern             | Example                         |
|----------------|---------------------|---------------------------------|
| Decomposition  | Breaking into steps | "Let me break this down: 1)..." |
| Hypothesis     | Forming theories    | "Could be A, B, or C..."        |
| Pivot          | Changing approach   | "Actually, wait..."             |
| Error Analysis | Debugging           | "The stack trace shows..."      |

### Anti-Pattern Detection

Identify reasoning approaches that consistently fail → warn against repeating.

**Details**: [thinking.md](specs/thinking.md)

---

## Tier 4: Session-Aware RAG

**Goal**: Query session history at the start of new sessions.

### Deliverables
- [ ] Chunking pipeline (turns, thinking, code, tools)
- [ ] Embedding + vector storage
- [ ] Query engine with metadata filtering
- [ ] ctx agent integration (PreToolUse hook)

### Architecture

```
Sessions → Chunker → Embedder → Vector Store
                                     ↓
New Session → Query: "What have I done here?" → Relevant Context
```

### Automatic Context

```bash
# In .claude/hooks/pre_tool_use.yaml
hooks:
  - name: session-context
    trigger: session_start
    command: ctx recall --auto --budget 4000
```

**Details**: [rag.md](specs/rag.md)

---

## Implementation Order

```
Week 1: Tier 1 (Browsable)
Week 2: Tier 2 (Import)
Week 3: Tier 3 (Thinking) + Tier 4 Start
Week 4: Tier 4 (RAG) Complete
```

### Dependencies

```
Tier 1 ──→ Tier 2 ──→ Tier 3
  │          │          │
  └──────────┴──────────┴──→ Tier 4 (all feed into RAG)
```

---

## File Structure

```
cmd/ctx/
├── recall.go              # ctx recall subcommand
├── recall_serve.go        # Tier 1: HTTP server
├── recall_import.go       # Tier 2: Import pipeline
├── recall_reason.go       # Tier 3: Reasoning search
└── recall_auto.go         # Tier 4: Auto-context

internal/
├── recall/
│   ├── parser/            # Session parsing
│   ├── renderer/          # HTML generation
│   ├── extractor/         # Decision/learning extraction
│   ├── thinking/          # Reasoning classification
│   ├── embedder/          # Vector embeddings
│   └── store/             # SQLite + FAISS
```

---

## Success Metrics

| Tier | Metric               | Target |
|------|----------------------|--------|
| 1    | Parse 1000 sessions  | < 5s   |
| 2    | Extraction precision | > 80%  |
| 3    | Category accuracy    | > 75%  |
| 4    | Retrieval relevance  | > 75%  |

---

## Risk Mitigation

| Risk                | Mitigation               |
|---------------------|--------------------------|
| Large session files | Stream processing        |
| LLM extraction cost | Optional, batched        |
| Embedding cost      | Local model (all-MiniLM) |
| Dedup errors        | Conservative defaults    |

---

## Next Steps

TODO: Decide. Maybe focus on Tier 1 only first.
