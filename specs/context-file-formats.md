# Context File Formats Specification

## Overview

Each context file serves a specific purpose and follows a consistent markdown structure. Files are designed to be human-readable, AI-parseable, and token-efficient.

## File Specifications

### DECISIONS.md

**Purpose**: Record architectural decisions with rationale so they don't get re-debated.

**Structure**:
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

**Example**:
```markdown
## [2025-01-15] Use TypeScript Strict Mode

**Status**: Accepted

**Context**: Starting new project, need to choose type checking level.

**Decision**: Enable TypeScript strict mode with all strict flags.

**Rationale**: Catches more bugs at compile time. Team has experience with strict mode. Upfront cost pays off in reduced runtime errors.

**Consequences**: More verbose type annotations required. Some third-party libraries need type assertions.

**Alternatives Considered**:
- Basic TypeScript: Rejected because it misses null checks
- JavaScript with JSDoc: Rejected because tooling support is weaker
```

---

### TASKS.md

**Purpose**: Track current work, planned work, and blockers.

**Structure**:
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
- [x] Task description — Completed YYYY-MM-DD

## Blocked
- [ ] Task description
  - Blocked by: Reason
  - Unblock strategy: How to resolve
```

**Tags**: Use inline tags for metadata:
- `#priority:high|medium|low`
- `#area:core|cli|docs|tests`
- `#estimate:Xh` (optional)

---

### LEARNINGS.md

**Purpose**: Capture lessons learned, gotchas, and tips that shouldn't be forgotten.

**Structure**:
```markdown
# Learnings

## Category Name

### Learning Title

**Discovered**: YYYY-MM-DD

**Context**: When/how was this learned?

**Lesson**: What's the takeaway?

**Application**: How should this inform future work?
```

**Example**:
```markdown
## Testing

### Vitest Mocks Must Be Hoisted

**Discovered**: 2025-01-15

**Context**: Tests were failing intermittently when mocking fs module.

**Lesson**: Vitest requires `vi.mock()` calls to be hoisted to the top of the file. Dynamic mocks need `vi.doMock()` instead.

**Application**: Always use `vi.mock()` at file top. Use `vi.doMock()` only when mock needs runtime values.
```

---

### CONVENTIONS.md

**Purpose**: Document project patterns, naming conventions, and standards.

**Structure**:
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
```code
// Example code
```

**Why**: Rationale for this pattern
```

---

### ARCHITECTURE.md

**Purpose**: Provide system overview and component relationships.

**Structure**:
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
- `path/to/file.ts` — Description

## Data Flow

Description or diagram of how data moves through the system.

## Boundaries

What's in scope vs out of scope for this codebase.
```

---

### DEPENDENCIES.md

**Purpose**: Document key dependencies and why they were chosen.

**Structure**:
```markdown
# Dependencies

## Runtime

### dependency-name

**Version**: ^X.Y.Z

**Purpose**: What we use it for

**Why this one**: Why chosen over alternatives

**Concerns**: Any known issues or limitations

## Development

### dev-dependency-name

**Version**: ^X.Y.Z

**Purpose**: What we use it for
```

---

### CONSTITUTION.md

**Purpose**: Define hard invariants — rules that must NEVER be violated, regardless of task.

**Structure**:
```markdown
# Constitution

These rules are INVIOLABLE. If a task requires violating these, the task is wrong.

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

**Usage**: AI agents MUST read this first and refuse tasks that violate invariants.

---

### GLOSSARY.md

**Purpose**: Define domain terms, abbreviations, and project vocabulary to ensure consistent language.

**Structure**:
```markdown
# Glossary

## Domain Terms

### Term Name
**Definition**: What it means in this project's context
**Not to be confused with**: Similar terms that mean different things
**Example**: How it's used

## Abbreviations

| Abbrev | Expansion | Context |
|--------|-----------|---------|
| ADR | Architectural Decision Record | Decision documentation |
| SUT | System Under Test | Testing |

## Project-Specific Terms

### Ralph Loop
A continuous AI iteration loop where the agent works on tasks until completion.
See: [Ralph Wiggum technique](https://ghuntley.com/ralph/)
```

---

### DRIFT.md

**Purpose**: Define signals that context is stale and needs updating. Used by `ctx drift` command.

**Structure**:
```markdown
# Drift Detection

## Automatic Checks

These are checked by `ctx drift`:

### Path References
- [ ] All paths in ARCHITECTURE.md exist in filesystem
- [ ] All paths in CONVENTIONS.md exist
- [ ] All file references in DECISIONS.md are valid

### Task References
- [ ] All issues referenced in TASKS.md exist (if issue tracker linked)
- [ ] No tasks older than 30 days without update

### Constitution Violations
- [ ] No secrets patterns detected in committed files
- [ ] No generated file patterns in git (if "no generated files" rule)

## Manual Review Triggers

Update context when:

- [ ] New team member joins (review CONVENTIONS.md)
- [ ] Major dependency upgraded (review DEPENDENCIES.md)
- [ ] Architecture diagram doesn't match code (weekly check)
- [ ] Sprint/milestone completed (archive old tasks)

## Staleness Indicators

| File | Stale If | Action |
|------|----------|--------|
| ARCHITECTURE.md | >30 days old | Review component list |
| TASKS.md | >50% completed | Archive and refresh |
| LEARNINGS.md | >20 items | Consolidate or archive |
| DECISIONS.md | References old patterns | Mark superseded |
```

---

### AGENT_PLAYBOOK.md

**Purpose**: Explicit instructions for how AI agents should read, apply, and update Active Memory.

**Structure**:
```markdown
# Agent Playbook

## Read Order

Load context files in this order (most critical first):

1. **CONSTITUTION.md** — Hard rules; refuse tasks that violate these
2. **TASKS.md** — What to work on
3. **CONVENTIONS.md** — How to write code
4. **ARCHITECTURE.md** — Where things go
5. **DECISIONS.md** — Why things are the way they are
6. **LEARNINGS.md** — Gotchas to avoid
7. **GLOSSARY.md** — Correct terminology
8. **DEPENDENCIES.md** — Reference material

## When to Update Memory

| Event | Update |
|-------|--------|
| Made architectural decision | Add to DECISIONS.md |
| Discovered gotcha/bug | Add to LEARNINGS.md |
| Established new pattern | Add to CONVENTIONS.md |
| Completed task | Mark [x] in TASKS.md |
| Added dependency | Add to DEPENDENCIES.md |
| Introduced new term | Add to GLOSSARY.md |

## How to Propose a Decision (ADR)

1. Check DECISIONS.md for existing related decisions
2. If new decision needed, add with format:
   - Date, Status, Context, Decision, Rationale, Consequences
3. Reference the decision in commit message

## How to Avoid Hallucinating Memory

**CRITICAL**: You have NO memory between sessions. All knowledge comes from files.

1. **Never assume** — If you don't see it in files, you don't know it
2. **Never invent history** — Don't claim "we discussed" or "we decided" without file evidence
3. **Verify before referencing** — Search files before citing past decisions
4. **When uncertain, say so** — "I don't see this documented" is valid
5. **Trust files over intuition** — If CONVENTIONS.md says X, do X even if it feels wrong

## Context Update Commands

Emit these to update context (processed by `ctx watch`):

\`\`\`xml
<context-update file="LEARNINGS.md" action="add" section="Testing">
### Vitest Mocks Must Be Hoisted
**Discovered**: 2025-01-19
**Lesson**: vi.mock() calls must be at file top...
</context-update>
\`\`\`
```

---

## Parsing Rules

1. **Headers define structure** — `#` for file title, `##` for sections, `###` for items
2. **Bold keys for fields** — `**Key**:` followed by value
3. **Code blocks are literal** — Never parse code block content as structure
4. **Lists are ordered** — Items appear in priority/chronological order
5. **Tags are inline** — Backtick-wrapped tags for metadata

## Token Efficiency

- Use abbreviations in tags, not prose
- Omit obvious words ("The", "This", etc.)
- Prefer bullet points over paragraphs for facts
- Keep examples minimal but illustrative
- Archive old completed items periodically
