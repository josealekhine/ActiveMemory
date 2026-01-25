---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

icon: lucide/files
---

![ctx](images/ctx-banner.png)

## Context Files

Each context file in `.context/` serves a specific purpose. Files are designed 
to be human-readable, AI-parseable, and token-efficient.

## File Overview

| File              | Purpose                                | Priority    |
|-------------------|----------------------------------------|-------------|
| CONSTITUTION.md   | Hard rules that must NEVER be violated | 1 (highest) |
| TASKS.md          | Current and planned work               | 2           |
| CONVENTIONS.md    | Project patterns and standards         | 3           |
| ARCHITECTURE.md   | System overview and components         | 4           |
| DECISIONS.md      | Architectural decisions with rationale | 5           |
| LEARNINGS.md      | Lessons learned, gotchas, tips         | 6           |
| GLOSSARY.md       | Domain terms and abbreviations         | 7           |
| DRIFT.md          | Staleness signals and update triggers  | 8           |
| AGENT_PLAYBOOK.md | Instructions for AI agents             | 9 (lowest)  |

## Read Order Rationale

The priority order follows a logical progression for AI agents:

1. **CONSTITUTION** — Inviolable rules first. The agent must know what it
   *cannot* do before attempting anything.

2. **TASKS** — Current work items. What the agent should focus on.

3. **CONVENTIONS** — How to write code. Patterns and standards to follow
   when implementing tasks.

4. **ARCHITECTURE** — System structure. Understanding of components and
   boundaries before making changes.

5. **DECISIONS** — Historical context. Why things are the way they are,
   to avoid re-debating settled decisions.

6. **LEARNINGS** — Gotchas and tips. Lessons from past work that inform
   current implementation.

7. **GLOSSARY** — Reference material. Domain terms and abbreviations for
   lookup as needed.

8. **DRIFT** — Staleness indicators. Lower priority since it's primarily
   for maintenance workflows, not active development.

9. **AGENT_PLAYBOOK** — Meta instructions last. How to use this context
   system itself. Loaded last because the agent should understand the
   *content* (rules, tasks, patterns) before the *operating manual*.

---

## CONSTITUTION.md

**Purpose:** Define hard invariants—rules that must **NEVER** be violated, 
regardless of the task.

AI agents read this first and should refuse tasks that violate these rules.

### Structure

```markdown
# Constitution

These rules are INVIOLABLE. If a task requires violating these, the task 
is wrong.

## Security Invariants

- [ ] Never commit secrets, tokens, API keys, or credentials
- [ ] Never store customer/user data in context files
- [ ] Never disable security linters without documented exception

## Quality Invariants

- [ ] All code must pass tests before commit
- [ ] No `any` types in TypeScript without documented reason
- [ ] No TODO comments in main branch (move to TASKS.md)

## Process Invariants

- [ ] All architectural changes require a decision record
- [ ] Breaking changes require version bump
- [ ] Generated files are never committed
```

### Guidelines

* Keep rules minimal and absolute
* Each rule should be enforceable (can verify compliance)
* Use checkbox format for clarity
* Never compromise on these rules

---

## TASKS.md

**Purpose:** Track current work, planned work, and blockers.

### Structure

```markdown
# Tasks

## In Progress

- [ ] Task description `#priority:high` `#area:core`
  - Current status: What's happening now
  - Blockers: Any impediments

## Next Up

- [ ] Task description `#priority:medium`
- [ ] Task description `#priority:low`

## Completed (Recent)

- [x] Task description — Completed YYYY-MM-DD

## Blocked

- [ ] Task description
  - Blocked by: Reason
  - Unblock strategy: How to resolve
```

### Tags

Use inline backtick-wrapped tags for metadata:

| Tag            | Values                         | Purpose                   |
|----------------|--------------------------------|---------------------------|
| `#priority`    | `high`, `medium`, `low`        | Task urgency              |
| `#area`        | `core`, `cli`, `docs`, `tests` | Codebase area             |
| `#estimate`    | `1h`, `4h`, `1d`               | Time estimate (optional)  |
| `#in-progress` | (none)                         | Currently being worked on |

### Status Markers

| Marker | Meaning                  |
|--------|--------------------------|
| `[ ]`  | Pending                  |
| `[x]`  | Completed                |
| `[-]`  | Skipped (include reason) |

### Guidelines

- Never delete tasks—mark as completed or skipped
- Keep completed tasks for 7 days, then archive
- One task should be `#in-progress` at a time

---

## DECISIONS.md

**Purpose:** Record architectural decisions with rationale so they don't
get re-debated.

### Structure

```markdown
# Decisions

## [YYYY-MM-DD] Decision Title

**Status**: Accepted | Superseded | Deprecated

**Context**: What situation prompted this decision?

**Decision**: What was decided?

**Rationale**: Why was this the right choice?

**Consequences**: What are the implications?

**Alternatives Considered**:
- Alternative A: Why rejected
- Alternative B: Why rejected
```

### Example

```markdown
## [2025-01-15] Use TypeScript Strict Mode

**Status**: Accepted

**Context**: Starting new project, need to choose type checking level.

**Decision**: Enable TypeScript strict mode with all strict flags.

**Rationale**: Catches more bugs at compile time. Team has experience
with strict mode. Upfront cost pays off in reduced runtime errors.

**Consequences**: More verbose type annotations required. Some
third-party libraries need type assertions.

**Alternatives Considered**:
- Basic TypeScript: Rejected because it misses null checks
- JavaScript with JSDoc: Rejected because tooling support is weaker
```

### Status Values

| Status     | Meaning                                 |
|------------|-----------------------------------------|
| Accepted   | Current, active decision                |
| Superseded | Replaced by newer decision (link to it) |
| Deprecated | No longer relevant                      |

---

## LEARNINGS.md

**Purpose:** Capture lessons learned, gotchas, and tips that shouldn't
be forgotten.

### Structure

```markdown
# Learnings

## Category Name

### Learning Title

**Discovered**: YYYY-MM-DD

**Context**: When/how was this learned?

**Lesson**: What's the takeaway?

**Application**: How should this inform future work?
```

### Example

```markdown
## Testing

### Vitest Mocks Must Be Hoisted

**Discovered**: 2025-01-15

**Context**: Tests were failing intermittently when mocking fs module.

**Lesson**: Vitest requires `vi.mock()` calls to be hoisted to the
top of the file. Dynamic mocks need `vi.doMock()` instead.

**Application**: Always use `vi.mock()` at file top. Use `vi.doMock()`
only when mock needs runtime values.
```

### Categories

Organize learnings by topic:

- Testing
- Build & Deploy
- Performance
- Security
- Third-Party Libraries
- Git & Workflow

---

## CONVENTIONS.md

**Purpose**: Document project patterns, naming conventions, and standards.

### Structure

```markdown
# Conventions

## Naming

- **Files**: kebab-case for all source files
- **Components**: PascalCase for React components
- **Functions**: camelCase, verb-first (getUser, parseConfig)
- **Constants**: SCREAMING_SNAKE_CASE

## Patterns

### Pattern Name

**When to use**: Situation description

**Implementation**:
// in triple backticks
// Example code

**Why**: Rationale for this pattern
```

### Guidelines

* Include concrete examples
* Explain the "why" not just the "what"
* Keep patterns minimal—only document what's non-obvious

---

## ARCHITECTURE.md

**Purpose**: Provide system overview and component relationships.

### Structure

```markdown
# Architecture

## Overview

Brief description of what the system does and how it's organized.

## Components

### Component Name

**Responsibility**: What this component does

**Dependencies**: What it depends on

**Dependents**: What depends on it

**Key Files**:
- path/to/file.ts — Description

## Data Flow

Description or diagram of how data moves through the system.

## Boundaries

What's in scope vs out of scope for this codebase.
```

### Guidelines

* Keep diagrams simple (Mermaid works well)
* Focus on boundaries and interfaces
* Update when major structural changes occur

---

## GLOSSARY.md

**Purpose**: Define domain terms, abbreviations, and project vocabulary.

### Structure

```markdown
# Glossary

## Domain Terms

### Term Name

**Definition**: What it means in this project's context

**Not to be confused with**: Similar terms that mean different things

**Example**: How it's used

## Abbreviations

| Abbrev | Expansion                     | Context                |
|--------|-------------------------------|------------------------|
| ADR    | Architectural Decision Record | Decision documentation |
| SUT    | System Under Test             | Testing                |
```

### Guidelines

* Define project-specific meanings
* Clarify potentially ambiguous terms
* Include abbreviations used in code or docs

---

## DRIFT.md

**Purpose**: Define signals that the context is stale and needs updating.

Used by `ctx drift` command to detect staleness.

### Structure

```markdown
# Drift Detection

## Automatic Checks

These are checked by `ctx drift`:

### Path References

- [ ] All paths in ARCHITECTURE.md exist in filesystem
- [ ] All paths in CONVENTIONS.md exist
- [ ] All file references in DECISIONS.md are valid

### Task References

- [ ] All issues referenced in TASKS.md exist
- [ ] No tasks older than 30 days without update

### Constitution Violations

- [ ] No secrets patterns detected in committed files

## Manual Review Triggers

Update context when:

- [ ] New team member joins (review CONVENTIONS.md)
- [ ] Major dependency upgraded (review ARCHITECTURE.md)
- [ ] Sprint/milestone completed (archive old tasks)

## Staleness Indicators

| File            | Stale If       | Action                 |
|-----------------|----------------|------------------------|
| ARCHITECTURE.md | >30 days old   | Review component list  |
| TASKS.md        | >50% completed | Archive and refresh    |
| LEARNINGS.md    | >20 items      | Consolidate or archive |
```

---

## AGENT_PLAYBOOK.md

**Purpose**: Explicit instructions for how AI agents should read, apply, 
and update context.

### Key Sections

**Read Order**: Priority order for loading context files

**When to Update**: Events that trigger context updates

**How to Avoid Hallucinating Memory:** Critical rules:

1. Never assume—if not in files, you don't know it
2. Never invent history—don't claim "we discussed" without evidence
3. Verify before referencing—search files before citing
4. When uncertain, say so
5. Trust files over intuition

**Context Update Commands:** Format for automated updates via `ctx watch`:

```xml
<context-update type="learning">Key takeaway from today's work</context-update>
<context-update type="decision">Use Redis for caching</context-update>
<context-update type="complete">user auth</context-update>
```

See [Integrations](integrations.md#context-update-commands) for full documentation.

---

## Parsing Rules

All context files follow these conventions:

1. **Headers define structure** — `#` for title, `##` for sections, `###` for 
   items
2. **Bold keys for fields** — `**Key**:` followed by value
3. **Code blocks are literal** — Never parse code block content as structure
4. **Lists are ordered** — Items appear in priority/chronological order
5. **Tags are inline** — Backtick-wrapped tags like `#priority:high`

## Token Efficiency

Keep context files concise:

- Use abbreviations in tags, not prose
- Omit obvious words ("The", "This")
- Prefer bullet points over paragraphs
- Keep examples minimal but illustrative
- Archive old completed items periodically
