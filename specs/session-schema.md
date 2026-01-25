# Session Schema Specification

> Defines the structure of session files produced by AI coding assistants.

## Overview

Session files are **JSONL** (JSON Lines) format, where each line is a 
self-contained JSON object representing a single message in the conversation.

## Message Object Schema

```typescript
interface SessionMessage {
  // Identity
  uuid: string;                    // Unique message ID
  parentUuid: string | null;       // Parent message ID (forms conversation tree)
  sessionId: string;               // Groups messages into a session
  requestId?: string;              // API request ID (for assistant messages)

  // Metadata
  timestamp: string;               // ISO 8601 timestamp
  type: "user" | "assistant";      // Message author
  userType?: "external" | "internal";
  isSidechain?: boolean;           // If message is in a branch

  // Context
  cwd: string;                     // Working directory at message time
  gitBranch?: string;              // Git branch name
  version: string;                 // Client version (e.g., "2.1.14")
  slug: string;                    // Human-friendly session name

  // Content
  message: MessageContent;
}

interface MessageContent {
  id: string;                      // Message ID from API
  type: "message";
  model: string;                   // Model name (e.g., "claude-opus-4-5-20251101")
  role: "user" | "assistant";
  content: ContentBlock[];
  stop_reason?: string | null;
  stop_sequence?: string | null;
  usage?: TokenUsage;
}

type ContentBlock =
  | TextBlock
  | ThinkingBlock
  | ToolUseBlock
  | ToolResultBlock;

interface TextBlock {
  type: "text";
  text: string;
}

interface ThinkingBlock {
  type: "thinking";
  thinking: string;               // The reasoning content
  signature: string;              // Cryptographic signature
}

interface ToolUseBlock {
  type: "tool_use";
  id: string;
  name: string;                   // Tool name (e.g., "bash", "write_file")
  input: Record<string, unknown>;
}

interface ToolResultBlock {
  type: "tool_result";
  tool_use_id: string;
  content: string | ContentBlock[];
  is_error?: boolean;
}

interface TokenUsage {
  input_tokens: number;
  output_tokens: number;
  cache_creation_input_tokens?: number;
  cache_read_input_tokens?: number;
  cache_creation?: {
    ephemeral_5m_input_tokens?: number;
    ephemeral_1h_input_tokens?: number;
  };
  service_tier?: string;
}
```

## Example Messages

### User Message

```json
{
  "uuid": "a1b2c3d4-...",
  "parentUuid": null,
  "sessionId": "session-123",
  "timestamp": "2026-01-21T07:50:00.000Z",
  "type": "user",
  "userType": "external",
  "cwd": "/home/user/project",
  "gitBranch": "main",
  "version": "2.1.14",
  "slug": "async-roaming-allen",
  "message": {
    "id": "msg_user_001",
    "type": "message",
    "role": "user",
    "content": [
      {
        "type": "text",
        "text": "How do I fix this bug in the parser?"
      }
    ]
  }
}
```

### Assistant Message with Thinking

```json
{
  "uuid": "e5f6g7h8-...",
  "parentUuid": "a1b2c3d4-...",
  "sessionId": "session-123",
  "timestamp": "2026-01-21T07:50:30.000Z",
  "type": "assistant",
  "requestId": "req_011CXL...",
  "cwd": "/home/user/project",
  "gitBranch": "main",
  "version": "2.1.14",
  "slug": "async-roaming-allen",
  "message": {
    "id": "msg_01Pq5vq...",
    "type": "message",
    "model": "claude-opus-4-5-20251101",
    "role": "assistant",
    "content": [
      {
        "type": "thinking",
        "thinking": "Let me analyze the parser code. The issue seems to be in the tokenizer where...",
        "signature": "Eu0BCkYI..."
      },
      {
        "type": "text",
        "text": "I see the issue. The parser is not handling escaped quotes correctly..."
      }
    ],
    "stop_reason": "end_turn",
    "usage": {
      "input_tokens": 1500,
      "output_tokens": 300,
      "cache_creation_input_tokens": 0,
      "cache_read_input_tokens": 44061
    }
  }
}
```

### Assistant Message with Tool Use

```json
{
  "uuid": "i9j0k1l2-...",
  "parentUuid": "e5f6g7h8-...",
  "sessionId": "session-123",
  "timestamp": "2026-01-21T07:51:00.000Z",
  "type": "assistant",
  "cwd": "/home/user/project",
  "message": {
    "role": "assistant",
    "content": [
      {
        "type": "tool_use",
        "id": "tool_01ABC",
        "name": "bash",
        "input": {
          "command": "cat parser.go | head -50"
        }
      }
    ]
  }
}
```

## Session Reconstruction

To reconstruct a conversation from a JSONL file:

1. **Read all lines** and parse each as JSON
2. **Group by sessionId** to separate distinct sessions
3. **Sort by timestamp** within each session
4. **Build conversation tree** using parentUuid references
5. **Handle sidechains** (isSidechain: true) as branches

```go
type Session struct {
    ID        string
    Slug      string
    CWD       string
    GitBranch string
    StartTime time.Time
    EndTime   time.Time
    Messages  []SessionMessage
    TokensIn  int
    TokensOut int
}

func ReconstructSession(messages []SessionMessage) *Session {
    sort.Slice(messages, func(i, j int) bool {
        return messages[i].Timestamp.Before(messages[j].Timestamp)
    })

    session := &Session{
        ID:        messages[0].SessionId,
        Slug:      messages[0].Slug,
        CWD:       messages[0].CWD,
        GitBranch: messages[0].GitBranch,
        StartTime: messages[0].Timestamp,
        EndTime:   messages[len(messages)-1].Timestamp,
        Messages:  messages,
    }

    for _, msg := range messages {
        if msg.Message.Usage != nil {
            session.TokensIn += msg.Message.Usage.InputTokens
            session.TokensOut += msg.Message.Usage.OutputTokens
        }
    }

    return session
}
```

## Derived Fields

| Field            | Derivation                                 |
|------------------|--------------------------------------------|
| `project`        | Last component of `cwd` path               |
| `duration`       | Last timestamp - first timestamp           |
| `turn_count`     | Count of user messages                     |
| `total_tokens`   | Sum of input + output tokens               |
| `cache_hit_rate` | cache_read / (cache_read + cache_creation) |
| `has_errors`     | Any tool_result with is_error: true        |

## File Organization

By date:
```
sessions/
├── 2026-01-20/
│   ├── session-abc123.jsonl
│   └── session-def456.jsonl
└── 2026-01-21/
    └── session-ghi789.jsonl
```

By project:
```
sessions/
├── ActiveMemory/
│   └── 2026-01-20-abc123.jsonl
└── SPIKE/
    └── 2026-01-19-ghi789.jsonl
```

## Parsing Considerations

1. **Streaming**: Files can be large; stream line-by-line
2. **Malformed lines**: Skip and log, don't fail entire file
3. **Missing fields**: Use sensible defaults
4. **Timezone**: All timestamps should be treated as UTC
5. **Encoding**: Files are UTF-8

## Validation

```go
func ValidateSession(messages []SessionMessage) error {
    if len(messages) == 0 {
        return errors.New("empty session")
    }

    sessionId := messages[0].SessionId
    for i, msg := range messages {
        if msg.SessionId != sessionId {
            return fmt.Errorf("message %d has different sessionId", i)
        }
        if msg.Timestamp.IsZero() {
            return fmt.Errorf("message %d has no timestamp", i)
        }
    }

    return nil
}
```