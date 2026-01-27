---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Prompting Guide
icon: lucide/message-circle
---

![ctx](images/ctx-banner.png)

## Prompting Guide

Effective prompts for working with AI assistants in ctx-enabled projects.

### Why This Matters

AI assistants don't automatically read context files. The right prompt triggers
the right behavior. This guide documents prompts that reliably produce good results.

---

## Session Start

### "Do you remember?"

**What it does**: Triggers the AI to read AGENT_PLAYBOOK, CONSTITUTION, sessions/,
and other context files before responding.

**When to use**: Start of every important session.

```
Do you remember what we were working on?
```

**Why it works**: The question implies prior context exists. The AI checks files
rather than admitting ignorance.

### "What's the current state?"

**What it does**: Prompts reading of TASKS.md, recent sessions, and status overview.

**When to use**: Resuming work after a break.

**Variants**:

- "Where did we leave off?"
- "What's in progress?"
- "Show me the open tasks"

---

## During Work

### "Why doesn't X work?"

**What it does**: Triggers root cause analysis rather than surface-level fixes.

**When to use**: When something fails unexpectedly.

**Why it works**: Framing as "why" encourages investigation before action.
The AI will trace through code, check configurations, and identify the actual cause.

!!! example "Real example"
    "Why can't I run /ctx-save?" led to discovering missing permissions
    in settings.local.json bootstrapping—a fix that benefited all users.

### "Is this consistent with our decisions?"

**What it does**: Prompts checking DECISIONS.md before implementing.

**When to use**: Before making architectural choices.

**Variants**:

- "Check if we've decided on this before"
- "Does this align with our conventions?"

### "What would break if we..."

**What it does**: Triggers defensive thinking and impact analysis.

**When to use**: Before making significant changes.

```
What would break if we change the Settings struct?
```

### "Before you start, read X"

**What it does**: Ensures specific context is loaded before work begins.

**When to use**: When you know relevant context exists in a specific file.

```
Before you start, read .context/sessions/2026-01-20-auth-discussion.md
```

---

## Reflection & Persistence

### "What did we learn?"

**What it does**: Prompts reflection on the session and often triggers adding
learnings to LEARNINGS.md.

**When to use**: After completing a task or debugging session.

**Why it works**: Explicit reflection prompt. The AI will summarize insights
and often offer to persist them.

### "Add this as a learning/decision"

**What it does**: Explicit persistence request.

**When to use**: When you've discovered something worth remembering.

```
Add this as a learning: "JSON marshal escapes angle brackets by default"
```

### "Save context before we end"

**What it does**: Triggers context persistence before session close.

**When to use**: End of session, or before switching topics.

**Variants**:

- "Let's persist what we did"
- "Update the context files"
- `/ctx-save` (slash command in Claude Code)

---

## Exploration & Research

### "Explore the codebase for X"

**What it does**: Triggers thorough codebase search rather than guessing.

**When to use**: When you need to understand how something works.

**Why it works**: "Explore" signals that investigation is needed, not immediate action.

### "How does X work in this codebase?"

**What it does**: Prompts reading actual code rather than explaining general concepts.

**When to use**: Understanding existing implementation.

```
How does session saving work in this codebase?
```

### "Find all places where X"

**What it does**: Comprehensive search across the codebase.

**When to use**: Before refactoring or understanding impact.

---

## Meta & Process

### "What should we document from this?"

**What it does**: Prompts identifying learnings, decisions, and conventions
worth persisting.

**When to use**: After complex discussions or implementations.

### "Is this the right approach?"

**What it does**: Invites the AI to challenge the current direction.

**When to use**: When you want a sanity check.

**Why it works**: Gives permission to disagree. AIs often default to agreeing;
this prompt signals you want honest assessment.

### "What am I missing?"

**What it does**: Prompts thinking about edge cases, overlooked requirements,
or unconsidered approaches.

**When to use**: Before finalizing a design or implementation.

---

## Anti-Patterns

Prompts that tend to produce poor results:

| Prompt | Problem | Better Alternative |
|--------|---------|-------------------|
| "Fix this" | Too vague, may patch symptoms | "Why is this failing?" |
| "Make it work" | Encourages quick hacks | "What's the right way to solve this?" |
| "Just do it" | Skips planning | "Plan this, then implement" |
| "You should remember" | Confrontational | "Do you remember?" |
| "Obviously..." | Discourages questions | State the requirement directly |

---

## Quick Reference

| Goal | Prompt |
|------|--------|
| Load context | "Do you remember?" |
| Resume work | "What's the current state?" |
| Debug | "Why doesn't X work?" |
| Validate | "Is this consistent with our decisions?" |
| Impact analysis | "What would break if we..." |
| Reflect | "What did we learn?" |
| Persist | "Add this as a learning" |
| Explore | "How does X work in this codebase?" |
| Sanity check | "Is this the right approach?" |
| Completeness | "What am I missing?" |

---

## Contributing

Found a prompt that works well?
[Open an issue](https://github.com/ActiveMemory/ctx/issues) or PR with:

1. The prompt text
2. What behavior it triggers
3. When to use it
4. Why it works (optional but helpful)
