---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

icon: lucide/repeat
---

![ctx](images/ctx-banner.png)

## Ralph + Context

*Like peanut butter and jelly: better together.*

The [Ralph Wiggum technique](https://ghuntley.com/ralph/) is an iterative AI development workflow where 
an agent works autonomously on tasks until completion. Context (`ctx`) and 
Ralph complement each other perfectly:

- **ctx** provides the *memory*: persistent context that survives across sessions
- **Ralph** provides the *loop*: autonomous iteration that runs until done

Together, they enable fully autonomous AI development where the agent remembers 
everything across iterations.

## How It Works

```mermaid
graph TD
    A[Start Loop] --> B[Load PROMPT.md]
    B --> C[AI reads .context/]
    C --> D[AI picks task from TASKS.md]
    D --> E[AI completes task]
    E --> F[AI updates context files]
    F --> G[AI commits changes]
    G --> H{Check signals}
    H -->|SYSTEM_CONVERGED| I[Done - all tasks complete]
    H -->|SYSTEM_BLOCKED| J[Done - needs human input]
    H -->|Continue| B
```

1. Loop reads `PROMPT.md` and invokes AI
2. AI loads context from `.context/`
3. AI picks one task and completes it
4. AI updates context files (mark task done, add learnings)
5. AI commits changes
6. Loop checks for completion signals
7. Repeat until converged or blocked

## Quick Start with Claude Code

Claude Code has a built-in Ralph Loop plugin:

```bash
# Start autonomous loop
/ralph-loop

# Cancel running loop
/cancel-ralph
```

That's it. The loop will:

1. Read your `PROMPT.md` for instructions
2. Pick tasks from `.context/TASKS.md`
3. Work until `SYSTEM_CONVERGED` or `SYSTEM_BLOCKED`

## Manual Loop Setup

For other AI tools, create a `loop.sh`:

```bash
#!/bin/bash
# loop.sh — a minimal Ralph loop

PROMPT_FILE="${1:-PROMPT.md}"
MAX_ITERATIONS="${2:-10}"
OUTPUT_FILE="/tmp/ralph_output.txt"

for i in $(seq 1 $MAX_ITERATIONS); do
  echo "=== Iteration $i ==="

  # Invoke AI with prompt
  cat "$PROMPT_FILE" | claude --print > "$OUTPUT_FILE" 2>&1

  # Display output
  cat "$OUTPUT_FILE"

  # Check for completion signals
  if grep -q "SYSTEM_CONVERGED" "$OUTPUT_FILE"; then
    echo "Loop complete: All tasks done"
    break
  fi

  if grep -q "SYSTEM_BLOCKED" "$OUTPUT_FILE"; then
    echo "Loop blocked: Needs human input"
    break
  fi

  sleep 2
done
```

Make it executable and run:

```bash
chmod +x loop.sh
./loop.sh
```

## The PROMPT.md File

The prompt file instructs the AI on how to work autonomously. Here's a template:

```markdown
# Autonomous Development Prompt

You are working on this project autonomously. Follow these steps:

## 1. Load Context

Read these files in order:
1. `.context/CONSTITUTION.md` — NEVER violate these rules
2. `.context/TASKS.md` — Find work to do
3. `.context/CONVENTIONS.md` — Follow these patterns
4. `.context/DECISIONS.md` — Understand past choices

## 2. Pick One Task

From `.context/TASKS.md`, select ONE task that is:
- Not blocked
- Highest priority available
- Within your capabilities

## 3. Complete the Task

- Write code following conventions
- Run tests if applicable
- Keep changes focused and minimal

## 4. Update Context

After completing work:
- Mark task complete in TASKS.md
- Add any learnings to LEARNINGS.md
- Add any decisions to DECISIONS.md

## 5. Commit Changes

Create a focused commit with clear message.

## 6. Signal Status

End your response with exactly ONE of:

- `SYSTEM_CONVERGED` — All tasks in TASKS.md are complete
- `SYSTEM_BLOCKED` — Cannot proceed, need human input (explain why)
- (no signal) — More work remains, continue to next iteration

## Rules

- ONE task per iteration
- NEVER skip tests
- NEVER violate CONSTITUTION.md
- Commit after each task
```

## Completion Signals

The loop watches for these signals in AI output:

| Signal               | Meaning            | When to Use                              |
|----------------------|--------------------|------------------------------------------|
| `SYSTEM_CONVERGED`   | All tasks complete | No pending tasks in TASKS.md             |
| `SYSTEM_BLOCKED`     | Cannot proceed     | Needs clarification, access, or decision |
| `BOOTSTRAP_COMPLETE` | Initial setup done | Project scaffolding finished             |

### Example Usage

```markdown
I've completed all tasks in TASKS.md:
- [x] Set up project structure
- [x] Implement core API
- [x] Add authentication
- [x] Write tests

No pending tasks remain.

SYSTEM_CONVERGED
```

```markdown
I cannot proceed with the "Deploy to production" task because:
- Missing AWS credentials
- Need confirmation on region selection

Please provide credentials and confirm deployment region.

SYSTEM_BLOCKED
```

## Context Integration

### Why ctx + Ralph Work Well Together

| Without ctx                 | With ctx                             |
|-----------------------------|--------------------------------------|
| Each iteration starts fresh | Each iteration has full history      |
| Decisions get re-made       | Decisions persist in DECISIONS.md    |
| Learnings are lost          | Learnings accumulate in LEARNINGS.md |
| Tasks can be forgotten      | Tasks tracked in TASKS.md            |

### Automatic Context Updates

During the loop, the AI should update context files:

**Mark task complete:**
```bash
ctx complete "implement user auth"
```

Or emit an update command (parsed by `ctx watch`):
```xml
<context-update type="complete">user auth</context-update>
```

**Add learning:**
```bash
ctx add learning "Rate limiting requires Redis connection"
```

Or via update command:
```xml
<context-update type="learning">Rate limiting requires Redis connection</context-update>
```

**Record decision:**
```bash
ctx add decision "Use JWT tokens for API authentication"
```

## Advanced: Watch Mode

Run `ctx watch` alongside the loop to automatically process context updates:

```bash
# Terminal 1: Run the loop
./loop.sh 2>&1 | tee /tmp/loop.log

# Terminal 2: Watch for context updates
ctx watch --log /tmp/loop.log --auto-save
```

The `--auto-save` flag periodically saves session snapshots, creating a 
history of the loop's progress.

## Example Project Setup

```
my-project/
├── .context/
│   ├── CONSTITUTION.md
│   ├── TASKS.md          # Work items for the loop
│   ├── DECISIONS.md
│   ├── LEARNINGS.md
│   ├── CONVENTIONS.md
│   └── sessions/         # Loop iteration history
├── PROMPT.md             # Instructions for the AI
├── loop.sh               # Loop script (if not using Claude Code)
└── src/                  # Your code
```

### Sample TASKS.md for Ralph

```markdown
# Tasks

## Phase 1: Setup

- [x] Initialize project structure
- [x] Set up testing framework

## Phase 2: Core Features

- [ ] Implement user registration `#priority:high`
- [ ] Add email verification `#priority:high`
- [ ] Create password reset flow `#priority:medium`

## Phase 3: Polish

- [ ] Add rate limiting `#priority:medium`
- [ ] Improve error messages `#priority:low`
```

The loop will work through these systematically, marking each complete.

## Troubleshooting

### Loop Runs Forever

**Cause:** AI not emitting completion signals

**Fix:** Ensure PROMPT.md explicitly instructs signaling:
```markdown
End EVERY response with one of:
- SYSTEM_CONVERGED (if all tasks done)
- SYSTEM_BLOCKED (if stuck)
```

### Context Not Persisting

**Cause**: AI not updating context files

**Fix**: Add explicit instructions to PROMPT.md:
```markdown
After completing a task, you MUST:
1. Run: ctx complete "<task>"
2. Add learnings: ctx add learning "..."
```

### Tasks Getting Repeated

**Cause**: Task not marked complete before next iteration

**Fix**: Ensure commit happens after context update:
```markdown
Order of operations:
1. Complete coding work
2. Update context files (ctx complete, ctx add)
3. Commit ALL changes including .context/
4. Then signal status
```

### AI Violating Constitution

**Cause**: Constitution not read first

**Fix**: Make constitution check explicit in PROMPT.md:
```markdown
BEFORE any work:
1. Read .context/CONSTITUTION.md
2. If task would violate ANY rule, emit SYSTEM_BLOCKED
3. Explain which rule prevents the work
```

## Resources

- [Ralph Wiggum Technique](https://ghuntley.com/ralph/) — Original blog post
- [Context CLI](cli-reference.md) — Command reference
- [Integrations](integrations.md) — Tool-specific setup
