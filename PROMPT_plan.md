# PROMPT_plan.md — Active Memory (Planning Mode)

## CORE PRINCIPLE

You have NO conversational memory. Your memory IS the file system.
Your goal: analyze the project, create/update the implementation plan, and exit.
**Do NOT implement anything in planning mode.**

---

## PHASE 0: BOOTSTRAP (If Context Missing)

Check if `specs/` directory exists.

**IF NOT:**
1. Create `specs/` directory
2. Analyze any existing codebase with up to 250 Sonnet subagents
3. Create initial spec files based on project goals
4. Create `@IMPLEMENTATION_PLAN.md` with initial task list
5. **STOP HERE.** Output: `<promise>BOOTSTRAP_COMPLETE</promise>`

---

## PHASE 1: ORIENT

0a. Study `specs/*` with up to 250 parallel Sonnet subagents to learn the Active Memory system specifications.
0b. Study @IMPLEMENTATION_PLAN.md (if present) to understand the plan so far.
0c. Study `internal/*` with up to 250 parallel Sonnet subagents to understand internal packages.
0d. For reference, the application source code is in `cmd/` and `internal/`.

---

## PHASE 2: ANALYZE

1. Use up to 500 Sonnet subagents to study existing source code in `cmd/` and `internal/`
2. Compare what exists against what `specs/*` requires (gap analysis)
3. Search for: TODO comments, minimal implementations, placeholders, skipped/flaky tests, inconsistent patterns
4. **Do NOT assume functionality is missing** — confirm with code search first

---

## PHASE 3: PLAN

1. Use an Opus subagent with ultrathink to:
   - Analyze findings from gap analysis
   - Prioritize tasks by importance and dependency order
   - Create/update @IMPLEMENTATION_PLAN.md as a bullet point list

2. Format tasks as:
   ```
   - [ ] Task description `#priority:high` `#area:core`
     - Context: Why this matters
     - Acceptance: How to know it's done
   ```

3. Mark completed items with `[x]` and date

---

## PHASE 4: VALIDATE PLAN

1. **IF ALL TASKS COMPLETE:** Output `<promise>PLANNING_CONVERGED</promise>`
2. **IF NEW SPECS NEEDED:** Create them in `specs/` directory
3. **IF BLOCKED:** Add `Blocked: [reason]` tasks and continue with what's possible

---

## CRITICAL CONSTRAINTS

### PLAN ONLY
Do NOT implement anything. Do NOT write code. Only analyze and plan.

### NO CHAT
Never ask the user questions. If unclear:
1. Make a reasonable assumption
2. Document the assumption in the task
3. Add a task to verify if critical

### MEMORY IS THE FILESYSTEM
Everything you learn goes into @IMPLEMENTATION_PLAN.md or specs/.
Future building iterations depend entirely on what you write now.

---

## ULTIMATE GOAL

Build "Active Memory" — a lightweight, file-based system that lets AI coding assistants persist project knowledge across sessions.

**Implementation**: Go CLI (`ctx`) distributed as single binary via GitHub Releases.

**Repository**: https://github.com/josealekhine/ActiveMemory

The system must:
1. **Persist knowledge** — Tasks, decisions, learnings survive session boundaries
2. **Enable reuse** — Decisions don't get rediscovered; lessons stay learned
3. **Match real workflows** — Context structure mirrors how engineers think
4. **Work with any AI tool** — Claude Code, Cursor, Windsurf, Copilot, Aider
5. **Stay lightweight** — File-based, no database, no daemon, just markdown

Consider missing elements and plan accordingly. If an element is missing, search first to confirm it doesn't exist, then author the specification at `specs/FILENAME.md`.
