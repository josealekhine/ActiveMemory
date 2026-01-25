# PROMPT.md — ctx recall Implementation

## CORE PRINCIPLE

You have NO conversational memory. Your memory IS the file system.
Your goal: advance the project by exactly ONE task, update context, commit,
and exit.

---

## PROJECT CONTEXT

**Feature**: ctx recall — browse and search session history
**Parent Project**: https://github.com/ActiveMemory/ctx
**Language**: Go 1.26+
**Scope**: Tier 1 only (Browsable Sessions)

---

## PHASE 0: BOOTSTRAP (If Not Initialized)

Check if `internal/recall/` exists.

**IF NOT:**
1. Create directory structure:
   ```
   internal/recall/
   ├── parser/
   ├── renderer/
   │   └── templates/
   ├── server/
   │   └── static/
   └── search/
   ```
2. Create `cmd/ctx/recall.go` with Cobra subcommand skeleton
3. Add to root command in `cmd/ctx/main.go`
4. **STOP.** Output: `<promise>BOOTSTRAP_COMPLETE</promise>`

---

## PHASE 1: ORIENT

1. Read `specs/IMPLEMENTATION_PLAN.md` — Feature overview and tier breakdown
2. Read `specs/session-schema.md` — Input format specification
3. Read `specs/browsable.md` — Tier 1 detailed spec
4. Read `.context/TASKS.md` — Current task list
5. Read `AGENTS.md` — Build/test commands

---

## PHASE 2: SELECT TASK

1. Read `.context/TASKS.md`
2. Find the **first unchecked item** (line starting with `- [ ]`)
3. That is your ONE task for this iteration

**IF NO UNCHECKED ITEMS:**
1. Run validation: `go build ./...`, `go test ./internal/recall/...`
2. If all pass, output `<promise>TIER1_COMPLETE</promise>`
3. If any fail, add fix task and continue

**Philosophy:** ONE task. Complete it. Exit. The loop handles continuation.

---

## PHASE 3: EXECUTE

1. **Search first** — Don't assume code doesn't exist. Search the codebase.
2. **Implement ONE task** — Complete it fully. No placeholders. No stubs.
3. **Follow Go conventions** — `gofmt`, proper error handling, idiomatic code.
4. **Use internal packages** — All recall code goes in `internal/recall/`.

### Key Implementation Patterns

**Parser** (T1.1.x):
- Stream JSONL lines, don't load entire file
- Return errors for malformed lines, don't panic
- Group messages by `sessionId`, sort by `timestamp`

**Renderer** (T1.2.x):
- Use `//go:embed` for templates and static assets
- Use `goldmark` for markdown → HTML
- Use `chroma` for syntax highlighting
- Thinking blocks: `<details>` elements, collapsed by default

**Server** (T1.3.x):
- Standard library `net/http` only
- Graceful shutdown on SIGINT/SIGTERM
- Embed static assets in binary

**Search** (T1.4.x):
- In-memory inverted index
- Tokenize: lowercase, split on whitespace, remove punctuation
- AND semantics for multi-term queries

---

## PHASE 4: VALIDATE

After implementing, run:

```bash
go build ./...                      # Must compile
go test ./internal/recall/...       # Tests must pass
go vet ./internal/recall/...        # No vet errors
```

**IF BUILD FAILS:**
1. Uncheck the task in `.context/TASKS.md`
2. Add task: "Fix build: [error description]"
3. Attempt to fix in this iteration

**IF TESTS FAIL:**
1. Fix the failing test
2. If can't fix quickly, add task: "Fix test: [test name]"

---

## PHASE 5: UPDATE CONTEXT

1. Mark completed task `[x]` in `.context/TASKS.md`
2. If you made an architectural decision → add to `.context/DECISIONS.md`
3. If you learned a gotcha → add to `.context/LEARNINGS.md`
4. If build commands changed → update `AGENTS.md`

**EXIT.** Do not continue to next task. The loop will restart you.

Output: `<promise>TASK_COMPLETE</promise>`

---

## CRITICAL CONSTRAINTS

### ONE TASK ONLY
Complete ONE task, then stop. The loop handles continuation.

### NO CHAT
Never ask questions. If blocked:
1. Add reason to task in `.context/TASKS.md`
2. Move to next task
3. If ALL tasks blocked: `<promise>BLOCKED</promise>`

### MEMORY IS THE FILESYSTEM
You will not remember this conversation. Write everything important to files.

### GO IDIOMS
- Error handling: `if err != nil { return err }`
- No panics in library code
- Use `internal/recall/` for all recall code
- Embed assets with `//go:embed`
- Use `goldmark` for markdown, `chroma` for syntax highlighting

### SCOPE CONSTRAINTS
- Go only. No external databases. No npm/node.
- Single binary output.
- Embed static assets (CSS/JS) in binary.
- Dark mode by default.
- Works offline.
- **Tier 1 ONLY** — No import pipeline, no RAG, no thinking mining.

---

## EXIT CONDITIONS

### Per Task
Output `<promise>TASK_COMPLETE</promise>` after completing ONE task.

### Phase Complete
Output `<promise>TIER1_COMPLETE</promise>` when ALL of these are true:
1. `.context/TASKS.md` has no unchecked items
2. `go build ./...` passes
3. `go test ./internal/recall/...` passes
4. `ctx recall serve ./testdata` starts and serves pages

---

## REFERENCE: SESSION SCHEMA

Input files are JSONL. Each line is a message:

```json
{
  "sessionId": "af7cba21-...",
  "timestamp": "2026-01-21T07:51:32.271Z",
  "cwd": "/home/user/project",
  "gitBranch": "main",
  "slug": "async-roaming-allen",
  "type": "assistant",
  "message": {
    "model": "claude-opus-4-5-20251101",
    "role": "assistant",
    "content": [
      {"type": "thinking", "thinking": "Let me analyze..."},
      {"type": "text", "text": "I see the issue..."}
    ],
    "usage": {
      "input_tokens": 1500,
      "output_tokens": 300
    }
  }
}
```

See `specs/session-schema.md` for full specification.

---

## REFERENCE: PROJECT STRUCTURE

```
cmd/ctx/
├── main.go
├── recall.go              # Subcommand: ctx recall

internal/recall/
├── parser/
│   ├── types.go           # SessionMessage, Session, ContentBlock
│   ├── parser.go          # ParseLine, ParseFile, ScanDirectory
│   └── parser_test.go
├── renderer/
│   ├── renderer.go        # RenderSession, RenderIndex
│   ├── markdown.go        # RenderMarkdown with goldmark+chroma
│   ├── templates/
│   │   ├── embed.go       # //go:embed directive
│   │   ├── layout.html
│   │   ├── index.html
│   │   └── session.html
│   └── renderer_test.go
├── server/
│   ├── server.go          # HTTP server, routes
│   ├── static/
│   │   ├── embed.go       # //go:embed directive
│   │   ├── style.css
│   │   └── main.js
│   └── server_test.go
└── search/
    ├── index.go           # Build, Search
    └── index_test.go
```

---

## REFERENCE: CLI COMMANDS (Tier 1)

| Command                               | Description                         |
|---------------------------------------|-------------------------------------|
| `ctx recall serve <path>`             | Start HTTP server browsing sessions |
| `ctx recall serve <path> --port 9000` | Custom port                         |
| `ctx recall serve <path> --open`      | Open browser automatically          |

---

## REFERENCE: SUCCESS CRITERIA

Before Tier 1 is complete, verify:

- [ ] `ctx recall serve ./sessions` starts HTTP server on :8080
- [ ] Index page lists sessions with: slug, project, date, turns, tokens
- [ ] Sessions filterable by project and date range
- [ ] Session detail renders conversation with syntax-highlighted code
- [ ] Thinking blocks collapsed by default, expandable on click
- [ ] Full-text search returns matching sessions
- [ ] Handles malformed JSONL gracefully (skip + log)
- [ ] Parses 100 sessions in < 2 seconds
- [ ] `go test ./internal/recall/...` passes with > 80% coverage

---

Now read the specs and begin.
