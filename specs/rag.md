# Tier 4: Session-Aware RAG

> Query session history at the start of new sessions for relevant context.

TODO: make some research on RAG before diving into this phase.

## Goals

- **Automatic Context**: Surface relevant history without manual lookup
- **Project Awareness**: Know what's been done in this codebase
- **Mistake Prevention**: Recall what failed before
- **Continuity**: Pick up where previous sessions left off

## CLI

```bash
# Manual search
ctx recall "how did I handle authentication?"

# Auto-context for current project
ctx recall --auto --budget 4000

# Index management
ctx recall --index ./sessions
ctx recall --sync
```

## Architecture

```
Sessions → Chunker → Embedder → Vector Store
                                     ↓
New Session → Query → Relevant Context → Inject
```

## Chunking Strategy

### Chunk Types

| Type      | Content                    | Use Case                     |
|-----------|----------------------------|------------------------------|
| Turn Pair | User Q + Assistant A       | General recall               |
| Thinking  | Reasoning blocks           | How problems were approached |
| Code      | Code + surrounding context | Implementation patterns      |
| Tool      | Tool calls + results       | What commands were run       |

```go
type Chunk struct {
    Type      string    // turn_pair, thinking, code, tool
    SessionID string
    Timestamp time.Time
    Content   string
    Project   string
    Branch    string
    Embedding []float32
}

func ChunkSession(session *Session) []Chunk {
    var chunks []Chunk
    
    // Turn pairs
    chunks = append(chunks, extractTurnPairs(session)...)
    
    // Thinking blocks (reuse from Tier 3)
    chunks = append(chunks, extractThinking(session)...)
    
    // Code blocks with context
    chunks = append(chunks, extractCode(session)...)
    
    // Significant tool calls
    chunks = append(chunks, extractTools(session)...)
    
    return chunks
}
```

## Embedding

**Model**: `all-MiniLM-L6-v2` (384 dims, local, free)

```go
func (c *Chunk) EmbeddingText() string {
    switch c.Type {
    case "turn_pair":
        return fmt.Sprintf("Q: %s\nA: %s", c.Question, c.Answer)
    case "thinking":
        return fmt.Sprintf("Reasoning [%s]: %s", c.Category, c.Content)
    case "code":
        return fmt.Sprintf("Code (%s): %s", c.Language, c.Content)
    case "tool":
        return fmt.Sprintf("Tool %s: %s → %s", c.ToolName, c.Input, c.Output)
    }
}
```

## Storage

TODO: can we live with file-based storage, instead of a DB?

SQLite + FAISS hybrid:

```sql
CREATE TABLE chunks (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    type TEXT,
    timestamp DATETIME,
    content TEXT,
    project TEXT,
    branch TEXT
);

CREATE TABLE embeddings (
    chunk_id TEXT PRIMARY KEY,
    embedding BLOB
);
```

## Query Engine

```go
type QueryEngine struct {
    Store    *VectorStore
    Embedder Embedder
}

func (qe *QueryEngine) Search(query string, filters Filters, limit int) []Chunk {
    // 1. Embed query
    embedding := qe.Embedder.Embed(query)
    
    // 2. Vector search (retrieve 3x for filtering)
    candidates := qe.Store.Search(embedding, limit*3)
    
    // 3. Apply filters (project, date, type)
    filtered := applyFilters(candidates, filters)
    
    // 4. Deduplicate similar chunks
    deduped := deduplicate(filtered, 0.9)
    
    // 5. Return top N
    return deduped[:min(limit, len(deduped))]
}

func (qe *QueryEngine) AutoContext(project, branch, task string, budget int) []Chunk {
    var chunks []Chunk
    
    // Recent from this project
    chunks = append(chunks, qe.getRecent(project, branch, 10)...)
    
    // Relevant to current task
    if task != "" {
        chunks = append(chunks, qe.Search(task, Filters{Project: project}, 10)...)
    }
    
    // Past failures in this project
    chunks = append(chunks, qe.getFailures(project, 5)...)
    
    // Fit to token budget
    return fitToBudget(deduplicate(chunks, 0.85), budget)
}
```

## ctx Integration

### PreToolUse Hook

```yaml
# .claude/hooks/pre_tool_use.yaml
hooks:
  - name: recall-context
    trigger: session_start
    command: ctx recall --auto --budget 4000 --format context
```

### Output Format

```markdown
## Session History

### Recent (ActiveMemory/main)
- Jan 20: Implemented error wrapping
- Jan 19: Fixed regex lookbehind issue

### Related
1. **Authentication** (Jan 15, SPIKE)
   "Used JWT with SPIFFE SVIDs..."

### Failures to Avoid
- ❌ Negative lookbehind in Go regex (not supported)
- ❌ DNS caching without TTL
```

## Token Budget

```go
func fitToBudget(chunks []Chunk, maxTokens int) []Chunk {
    // Sort by relevance
    sort.Slice(chunks, func(i, j int) bool {
        return chunks[i].Score > chunks[j].Score
    })
    
    var result []Chunk
    used := 0
    for _, chunk := range chunks {
        tokens := countTokens(chunk.Content)
        if used + tokens > maxTokens {
            break
        }
        result = append(result, chunk)
        used += tokens
    }
    return result
}
```

## Incremental Sync

```go
func (w *Watcher) Sync() error {
    lastProcessed := w.Store.LastProcessedTime()
    newSessions := findSessionsAfter(lastProcessed)
    
    for _, session := range newSessions {
        chunks := ChunkSession(session)
        embeddings := w.Embedder.BatchEmbed(chunks)
        w.Store.Add(chunks, embeddings)
    }
    
    return nil
}
```

## Tasks

| Phase | Task                  | Hours |
|-------|-----------------------|-------|
| 4.1   | Chunking pipeline     | 8     |
| 4.2   | Embedding integration | 8     |
| 4.3   | Vector storage        | 8     |
| 4.4   | Query engine          | 8     |
| 4.5   | ctx integration       | 8     |

## Storage Estimates

| Sessions | Chunks | Embeddings | Total   |
|----------|--------|------------|---------|
| 100      | 5K     | 7 MB       | ~25 MB  |
| 1,000    | 50K    | 73 MB      | ~225 MB |
| 10,000   | 500K   | 730 MB     | ~2.3 GB |

## Success Metrics

| Metric                | Target   |
|-----------------------|----------|
| Retrieval relevance   | > 75%    |
| Query latency         | < 1s     |
| Index 1000 sessions   | < 30 min |
| Token budget accuracy | ±10%     |