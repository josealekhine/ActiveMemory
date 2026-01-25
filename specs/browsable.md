# Tier 1: Browsable Sessions

> A local web interface for browsing and searching session history.

## Goals

- **Visibility**: See all sessions in human-readable format
- **Navigation**: Find sessions by date, project, or content
- **Insight**: Understand token costs and conversation flow
- **Foundation**: Validate parsing before higher tiers

## CLI

```bash
ctx recall serve ./sessions
ctx recall serve ./sessions --port 8080 --open
```

## Architecture

```
Session Files → Parser → Rendered HTML → HTTP Server
   (JSONL)                 + Index       localhost:8080
```

## Routes

| Route                  | Description                  |
|------------------------|------------------------------|
| `GET /`                | Index page with session list |
| `GET /session/:id`     | Session detail page          |
| `GET /api/sessions`    | JSON session list            |
| `GET /api/session/:id` | JSON session detail          |

## Index Page

```
┌─────────────────────────────────────────────────────────────┐
│  ctx recall                                     [Search 🔍] │
├─────────────────────────────────────────────────────────────┤
│  [All Projects ▼] [Last 7 days ▼] [All branches ▼]         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐   │
│  │ async-roaming-allen                    Jan 21, 2026 │   │
│  │ ActiveMemory • main • 15 turns • 45K tokens         │   │
│  │ "How do I fix this bug in the parser..."            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ brave-sailing-mercury                  Jan 20, 2026 │   │
│  │ SPIKE • feature/auth • 8 turns • 22K tokens         │   │
│  │ "Implement JWT validation..."                       │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  142 sessions • 1.2M tokens • 23 projects                   │
└─────────────────────────────────────────────────────────────┘
```

## Session Detail Page

```
┌────────────────────────────────┬────────────────────────────┐
│                                │  Metadata                  │
│  ┌──────────────────────────┐  │  Date: Jan 21, 2026        │
│  │ 👤 User         07:50:00 │  │  Duration: 12m 34s         │
│  │ How do I fix this bug?   │  │  Project: ActiveMemory     │
│  └──────────────────────────┘  │  Branch: main              │
│                                │                            │
│  ┌──────────────────────────┐  │  Tokens                    │
│  │ 🤖 Assistant    07:50:30 │  │  In: 44,061 Out: 1,169     │
│  │                          │  │                            │
│  │ [▶ Thinking]             │  │  Tools: bash (3)           │
│  │                          │  │                            │
│  │ I see the issue...       │  │  [Export MD] [Export JSON] │
│  └──────────────────────────┘  │                            │
└────────────────────────────────┴────────────────────────────┘
```

## Data Structures

```go
type SessionIndex struct {
    Sessions  []SessionSummary
    ByProject map[string][]string  // project → sessionIds
    ByDate    map[string][]string  // YYYY-MM-DD → sessionIds
}

type SessionSummary struct {
    ID           string
    Slug         string
    Project      string
    Branch       string
    StartTime    time.Time
    TurnCount    int
    TokensIn     int
    TokensOut    int
    FirstMessage string  // Truncated preview
}
```

## Rendering

- Markdown → HTML via goldmark + GFM
- Code highlighting via chroma
- Thinking blocks: collapsed by default, click to expand
- Dark mode CSS

## Search

In-memory inverted index:
```go
type SearchIndex struct {
    Terms   map[string][]string  // term → sessionIds
    Content map[string]string    // sessionId → full text
}
```

## Tasks

| Phase | Task                      | Hours |
|-------|---------------------------|-------|
| 1.1   | Session parser + grouping | 4     |
| 1.2   | HTML renderer + templates | 4     |
| 1.3   | HTTP server + routes      | 4     |
| 1.4   | Search + filters          | 4     |

## Success Metrics

| Metric              | Target  |
|---------------------|---------|
| Parse 1000 sessions | < 5s    |
| Render page         | < 100ms |
| Search response     | < 200ms |