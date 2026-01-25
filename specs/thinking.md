# Tier 3: Thinking Mining

> Extract and index reasoning patterns from thinking blocks.

## Goals

- **Capture Reasoning**: Preserve the "how" not just the "what"
- **Learn from Failure**: Identify anti-patterns and dead ends
- **Build a Corpus**: Searchable database of problem-solving approaches

## CLI

```bash
ctx recall reason "how to debug connection refused"
ctx recall reason --category error_analysis "database"
ctx recall reason --outcome failure "caching"  # Anti-patterns
```

## Reasoning Categories

### Decomposition
Breaking problems into steps.
```
"Let me break this down:
1. First, understand the input format
2. Then, identify transformation rules..."
```

### Hypothesis
Forming and testing theories.
```
"The error could be:
- A: Missing dependency
- B: Wrong version
Let me check A first..."
```

### Pivot
Recognizing dead ends.
```
"Actually, this won't work because the API
doesn't support batch operations. Let me try..."
```

### Error Analysis
Debugging from stack traces.
```
"The stack trace shows nil pointer at line 45.
So config isn't loading the database section..."
```

## Outcome Detection

| Outcome   | Indicators                                       |
|-----------|--------------------------------------------------|
| Success   | Tool succeeds, "that worked", continued progress |
| Failure   | `is_error: true`, "didn't work", pivot follows   |
| Abandoned | Topic changes, no resolution                     |

```go
func DetectOutcome(thinking *ThinkingBlock, session *Session, idx int) Outcome {
    // Look at next 5 messages
    for i := idx + 1; i < min(idx+5, len(session.Messages)); i++ {
        msg := session.Messages[i]
        
        // Check for tool errors
        if hasToolError(msg) {
            return OutcomeFailure
        }
        
        // Check for success language
        if containsSuccess(msg.Text) {
            return OutcomeSuccess
        }
        
        // Check for pivot in next thinking
        if nextThinking := getThinking(msg); nextThinking != nil {
            if ClassifyThinking(nextThinking) == CategoryPivot {
                return OutcomeFailure
            }
        }
    }
    return OutcomeUnknown
}
```

## Data Model

```go
type ReasoningUnit struct {
    ID           string
    SessionID    string
    Timestamp    time.Time
    ThinkingText string
    Category     string    // decomposition, hypothesis, pivot, error_analysis
    Outcome      string    // success, failure, abandoned, unknown
    Tags         []string
    Embedding    []float32
}
```

## Storage

TODO: can we stay with a file-system based approach; what benefits 
do we get from a relational database?

```sql
CREATE TABLE reasoning_units (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    timestamp DATETIME,
    thinking_text TEXT,
    category TEXT,
    outcome TEXT
);

CREATE VIRTUAL TABLE reasoning_fts USING fts5(thinking_text);
```

## Anti-Pattern Detection

Reasoning that consistently fails:

```go
type AntiPattern struct {
    Pattern       string
    FailureRate   float64  // > 60%
    Occurrences   int      // >= 3
    BetterApproach *ReasoningUnit
}

func DetectAntiPatterns(units []ReasoningUnit) []AntiPattern {
    // Cluster similar reasoning by embedding
    // Calculate failure rate per cluster
    // Flag clusters with >60% failure rate and 3+ occurrences
    // Find successful alternatives for the same problem type
}
```

## Query Interface

```
$ ctx recall reason "debug connection refused"

Found 3 relevant patterns:

1. [SUCCESS] Error Analysis - Jan 20 (async-roaming-allen)
   "Check if port is in use, then firewall rules..."
   → Port was already bound

2. [FAILURE] Hypothesis - Jan 18 (brave-sailing-mercury)
   "Maybe the config has wrong hostname..."
   → Wrong hypothesis, actual issue was firewall

3. [SUCCESS] Decomposition - Jan 15 (calm-dancing-neptune)
   "Break down: 1) Check DNS, 2) Check routing..."
   → DNS caching issue
```

## Tasks

| Phase | Task                  | Hours |
|-------|-----------------------|-------|
| 3.1   | Thinking block parser | 4     |
| 3.2   | Category classifier   | 8     |
| 3.3   | Outcome detector      | 4     |
| 3.4   | Index + search        | 8     |

## Success Metrics

| Metric            | Target |
|-------------------|--------|
| Category accuracy | > 75%  |
| Outcome detection | > 80%  |
| Search relevance  | > 70%  |