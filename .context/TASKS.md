# Tasks — Context CLI

# Tasks

### Phase 1: Parser

- [ ] T1.1.1: Define data structures in `internal/recall/parser/types.go`
  - `SessionMessage`, `Session`, `ContentBlock` interface
  - `TextBlock`, `ThinkingBlock`, `ToolUseBlock`, `ToolResultBlock`
  - `TokenUsage` struct
  - Must match `specs/recall/v0.1.0/01-session-schema.md`

- [ ] T1.1.2: Implement line parser in `internal/recall/parser/parser.go`
  - `ParseLine(line []byte) (*SessionMessage, error)`
  - Handle malformed JSON (return error, don't panic)
  - Handle missing optional fields with defaults

- [ ] T1.1.3: Implement session grouper
  - `ParseFile(path string) (map[string]*Session, error)`
  - Stream lines, don't load entire file into memory
  - Sort messages by timestamp within each session
  - Calculate: TurnCount, TotalTokensIn, TotalTokensOut, Duration

- [ ] T1.1.4: Implement directory scanner
  - `ScanDirectory(path string) ([]*Session, error)`
  - Recurse into subdirectories
  - Skip non-JSONL files
  - Aggregate sessions across files

### Phase 2: Renderer

- [ ] T1.2.1: Set up template system in `internal/recall/renderer/`
  - Use `//go:embed` for templates
  - Create `templates/layout.html` with dark mode CSS
  - Create `templates/index.html`
  - Create `templates/session.html`

- [ ] T1.2.2: Implement markdown renderer
  - Add goldmark + chroma dependencies
  - `RenderMarkdown(text string) template.HTML`
  - Enable GFM extensions
  - Syntax highlighting with monokai theme

- [ ] T1.2.3: Implement session renderer
  - `RenderSession(session *Session) (*RenderedSession, error)`
  - Wrap thinking blocks in `<details>` (collapsed by default)
  - Format tool calls with syntax-highlighted input/output

- [ ] T1.2.4: Implement index renderer
  - `RenderIndex(sessions []*Session, filters Filters) (*RenderedIndex, error)`
  - Sort by date (newest first)
  - Include preview (first 100 chars of first user message)
  - Show aggregate stats

### Phase 3: Server

- [ ] T1.3.1: Set up HTTP server in `internal/recall/server/`
  - Standard library `net/http`
  - Graceful shutdown on SIGINT/SIGTERM
  - Embed static assets with `//go:embed`

- [ ] T1.3.2: Implement index route
  - `GET /` — render index page
  - `GET /?project=X` — filter by project
  - `GET /?after=DATE&before=DATE` — filter by date
  - `GET /?q=QUERY` — search sessions

- [ ] T1.3.3: Implement session detail route
  - `GET /session/:id` — render session detail
  - 404 if not found
  - Include back link to index

- [ ] T1.3.4: Implement API routes
  - `GET /api/sessions` — JSON session list
  - `GET /api/session/:id` — JSON session detail

### Phase 4: Search

- [ ] T1.4.1: Implement search index in `internal/recall/search/`
  - Inverted index: term → sessionIds
  - `Build(sessions []*Session)`
  - Tokenize: lowercase, split whitespace, remove punctuation

- [ ] T1.4.2: Implement search query
  - `Search(query string, limit int) []string`
  - AND semantics (all terms must match)
  - Sort by term frequency score

### Phase 5: CLI

- [ ] T1.5.1: Add recall subcommand
  - Create `cmd/ctx/recall.go`
  - `ctx recall serve <path>` — start server
  - `--port` flag (default: 8080)
  - `--open` flag to open browser

- [ ] T1.5.2: Add help and validation
  - Validate path exists and has JSONL files
  - Print URL when server starts
  - Print stats (sessions loaded, time taken)

## Backlog

- [x] Rename vars in the config package.
- [x] Why is agent runbook lowest in reading priority order?
  - follow-up: is it enforced?
- [x] Create a list of what CLI options (if any) are not implemented yet.
- [ ] Verify all Markdown files by "actually reading them"; take notes for
  follow-up actions.
- [x] All go code should have godoc and testing.
- [ ] GitHub CI linter is giving errors that need fixing.
- [x] Manual code review. take notes.
- [x] Add tests per file.
- [ ] validate everything in the docs with a skeptical eye.
- [x] consider the case where `ctx` is not called from within AI prompt:
  - does the command still make sense?
  - does it create the expected output?
- [ ] Cut the first release.
  - Versioning strategy.
  - Always have a `latest` tag pointing to the latest release.
  - Or, maybe just use the `latest` tag at all times?
- [ ] have a proper email for security vulnerability reports.
- [ ] compare versions of recent change and the last AI-assisted version and
      ask AI what we have learned about this.
- [ ] CREATE SHORTS and VODS
- [ ] Trace the entire git history and sessions, create an extensive document
      of what we did and how it progressed, and then create a blog post about it.

