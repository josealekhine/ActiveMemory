# Agent Playbook

## Read Order

1. CONSTITUTION.md — Hard rules, NEVER violate
2. TASKS.md — What to work on next
3. CONVENTIONS.md — How to write code
4. ARCHITECTURE.md — Where things go
5. DECISIONS.md — Why things are the way they are
6. LEARNINGS.md — Gotchas to avoid
7. GLOSSARY.md — Correct terminology

## Session History

**IMPORTANT**: Check `.context/sessions/` for full conversation dumps from previous sessions.

If you're confused about context or need a deep dive into past discussions:
```
ls .context/sessions/
```

**Curated session files** are named `YYYY-MM-DD-HHMMSS-<topic>.md` (e.g., `2025-01-20-164600-feature-discussion.md`). These are updated throughout the session.

**Auto-snapshot files** are named `YYYY-MM-DD-HHMMSS-<event>.jsonl` (e.g., `2025-01-20-170830-pre-compact.jsonl`). These are immutable once created.

**Auto-save triggers** (for Claude Code users):
- **SessionEnd hook** → auto-saves transcript on exit, including Ctrl+C
- **PreCompact** → saves before `ctx compact` archives old tasks
- **Manual** → `ctx session save`

See `.claude/hooks/auto-save-session.sh` for the implementation.

## Session File Structure (Suggested)

Adapt this structure based on session type. Not all sections are needed for every session.

### Core Sections (Always Include)
```markdown
# Session: <Topic>

**Date**: YYYY-MM-DD
**Topic**: Brief description
**Type**: feature | bugfix | architecture | exploration | planning

---

## Summary
What was discussed/accomplished (2-3 sentences)

## Key Decisions
Bullet points of decisions made (if any)

## Tasks for Next Session
What to pick up next
```

### Context-Dependent Sections

| Session Type | Additional Sections |
|--------------|---------------------|
| **Feature discussion** | Requirements, Design options, Implementation plan |
| **Bug investigation** | Symptoms, Root cause, Fix applied, Prevention |
| **Architecture decision** | Context, Options considered, Trade-offs, Decision |
| **Exploration/Research** | Questions, Findings, Open questions |
| **Planning** | Goals, Milestones, Dependencies, Risks |
| **Quick fix** | Problem, Solution, Files changed (minimal format) |

### When to Go Minimal

For quick sessions (<15 min), just capture:
```markdown
# Session: Quick Fix - <Topic>
**Date**: YYYY-MM-DD
**Summary**: One sentence
**Files changed**: List
```

### When to Go Deep

For complex sessions (architecture, debugging), include:
- User quotes that capture key insights
- Technical context (platform, versions, constraints)
- Links to related sessions or decisions
- Code snippets or error messages if relevant

---

## When to Update Memory

| Event | Action |
|-------|--------|
| Made architectural decision | Add to DECISIONS.md |
| Discovered gotcha/bug | Add to LEARNINGS.md |
| Established new pattern | Add to CONVENTIONS.md |
| Completed task | Mark [x] in TASKS.md |
| Had important discussion | Save to sessions/ |

## Before Session Ends

**CRITICAL**: Before the user ends the session, offer to save context:
1. Curated summary → LEARNINGS.md, DECISIONS.md, TASKS.md
2. Full conversation dump → `.context/sessions/YYYY-MM-DD-<topic>.md`

## How to Avoid Hallucinating Memory

Never assume. If you don't see it in files, you don't know it.

- Don't claim "we discussed X" without file evidence
- Don't invent history - check sessions/ for actual discussions
- If uncertain, say "I don't see this documented"
- Trust files over intuition
