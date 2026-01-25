# PROMPT.md — Context Go CLI Implementation

## CORE PRINCIPLE

You have NO conversational memory. Your memory IS the file system.
Your goal: advance the project by exactly ONE task, update context, commit, 
and exit.

---

## PROJECT CONTEXT

**Project**: Context — persistent context for AI coding assistants
**Repository**: https://github.com/ActiveMemory/ctx
**Language**: Go 1.25+
**Distribution**: Single binary via GitHub Releases

---

## PHASE 0: BOOTSTRAP (If Project Not Initialized)

Check if `go.mod` exists.

**IF NOT:**
1. Run `go mod init github.com/ActiveMemory/ctx`
2. Create directory structure:
   ```
   cmd/ctx/main.go
   internal/cli/
   internal/context/
   internal/files/
   internal/drift/
   internal/templates/
   templates/
   hack/
   examples/
   ```
3. Install dependencies: `go get github.com/spf13/cobra@latest github.com/fatih/color@latest gopkg.in/yaml.v3@latest`
4. Create minimal `cmd/ctx/main.go` with Cobra skeleton
5. Create `IMPLEMENTATION_PLAN.md` with task list from Phase breakdown below
6. **STOP.** Output: `<promise>BOOTSTRAP_COMPLETE</promise>`

---

## PHASE 1: ORIENT

1. Read `specs/core-architecture.md` — Overall design philosophy
2. Read `specs/go-cli-implementation.md` — Go project structure and patterns
3. Read `specs/cli.md` — All CLI commands and their behavior
4. Read `specs/context-file-formats.md` — File format specifications
5. Read `specs/context-loader.md` — Loading and parsing logic
6. Read `specs/context-updater.md` — Update command handling
7. Read `IMPLEMENTATION_PLAN.md` — Current task list
8. Read `AGENTS.md` — Build/test commands

---

## PHASE 2: SELECT TASK

1. Read `IMPLEMENTATION_PLAN.md` for the current directive
2. Follow the directive — typically: "Check `.context/TASKS.md`"
3. Read `.context/TASKS.md` and pick the **first unchecked item** in "Next Up"

**IF NO UNCHECKED ITEMS in `.context/TASKS.md`:**
1. Check `IMPLEMENTATION_PLAN.md` for North Star goals
2. Remind user about Endgame goals before exiting
3. Output `<promise>DONE</promise>`

**Philosophy:** Tasks live in the agent's mind (`.context/TASKS.md`). 
The orchestrator (`IMPLEMENTATION_PLAN.md`) provides the meta-directive, 
not the task list

---

## PHASE 3: EXECUTE

1. **Search first** — Don't assume code doesn't exist. Search the codebase.
2. **Implement ONE task** — Complete it fully. No placeholders. No stubs.
3. **Follow Go conventions** — `gofmt`, proper error handling, idiomatic code.
4. **Use internal packages** — Put reusable code in `internal/`, not `cmd/`.

---

## PHASE 4: VALIDATE

After implementing, run:

```bash
go build ./...          # Must compile
go test ./...           # Tests must pass
go vet ./...            # No vet errors
```

**IF BUILD FAILS:**
1. Uncheck the task
2. Add task: "Fix build: [error description]"
3. Attempt to fix in this iteration

**IF TESTS FAIL:**
1. Fix the failing test
2. If can't fix quickly, add task: "Fix test: [test name]"

---

## PHASE 5: UPDATE CONTEXT

1. Mark completed task `[x]` in `.context/TASKS.md`
2. Move task to "Completed (Recent)" section with date
3. If you made an architectural decision → run `ctx add decision "description"`
4. If you learned a gotcha → run `ctx add learning "description"`
5. If build commands changed → update `AGENTS.md`

---

## PHASE 6: COMMIT & EXIT

```bash
git add -A
git commit -m "feat(cli): implement [command/feature]"  # or fix/docs/test/chore
git push origin main
```

**EXIT.** Do not continue to next task. The loop will restart you.

---

## CRITICAL CONSTRAINTS

### ONE TASK ONLY
Complete ONE task, then stop. The loop handles continuation.

### NO CHAT
Never ask questions. If blocked:
1. Move task to "Blocked" section in `.context/TASKS.md` with reason
2. Move to next task in "Next Up"
3. If ALL tasks blocked: `<promise>SYSTEM_BLOCKED</promise>`

### MEMORY IS THE FILESYSTEM
You will not remember this conversation. Write everything important to files.

### GO IDIOMS
- Error handling: `if err != nil { return err }`
- No panics in library code
- Use `internal/` for non-exported packages
- Embed templates with `//go:embed`

---

## TEMPLATE FILES TO CREATE

When implementing `ctx init`, embed these templates:

### templates/CONSTITUTION.md
```markdown
# Constitution

These rules are INVIOLABLE. If a task requires violating these, 
the task is wrong.

## Security Invariants

- [ ] Never commit secrets, tokens, API keys, or credentials
- [ ] Never store customer/user data in context files

## Quality Invariants

- [ ] All code must pass tests before commit
- [ ] No TODO comments in main branch (move to TASKS.md)

## Process Invariants

- [ ] All architectural changes require a decision record
```

### templates/TASKS.md
```markdown
# Tasks

## In Progress

## Next Up

## Completed (Recent)

## Blocked
```

### templates/DECISIONS.md
```markdown
# Decisions

<!-- Use this format for each decision:

## [YYYY-MM-DD] Decision Title

**Status**: Accepted | Superseded | Deprecated

**Context**: What situation prompted this decision?

**Decision**: What was decided?

**Rationale**: Why was this the right choice?

**Consequences**: What are the implications?
-->
```

### templates/LEARNINGS.md
```markdown
# Learnings

<!-- Add gotchas, tips, and lessons learned here -->
```

### templates/CONVENTIONS.md
```markdown
# Conventions

## Naming

## Patterns

## Testing
```

### templates/ARCHITECTURE.md
```markdown
# Architecture

## Overview

## Components

## Data Flow
```

### templates/GLOSSARY.md
```markdown
# Glossary

## Domain Terms

## Abbreviations
```

### templates/DRIFT.md
```markdown
# Drift Detection

## Automatic Checks

## Manual Review Triggers

## Staleness Indicators
```

### templates/AGENT_PLAYBOOK.md
```markdown
# Agent Playbook

## Read Order

1. CONSTITUTION.md
2. TASKS.md
3. CONVENTIONS.md
4. ARCHITECTURE.md
5. DECISIONS.md
6. LEARNINGS.md
7. GLOSSARY.md

## When to Update Memory

## How to Avoid Hallucinating Memory

Never assume. If you don't see it in files, you don't know it.
```

---

## EXIT CONDITIONS

Output `<promise>DONE</promise>` ONLY when ALL of these are true:

1. `.context/TASKS.md` has no unchecked items in "Next Up"
2. `go build ./...` passes
3. `go test ./...` passes
4. You have reminded the user about the North Star goals in `IMPLEMENTATION_PLAN.md`

---

## REFERENCE: CLI COMMANDS

| Command                    | Description                       |
|----------------------------|-----------------------------------|
| `ctx init`                 | Create `.context/` with templates |
| `ctx status`               | Show context summary              |
| `ctx load`                 | Output assembled context          |
| `ctx agent`                | Print AI-ready context packet     |
| `ctx add <type> "content"` | Add decision/task/learning        |
| `ctx complete <task>`      | Mark task done                    |
| `ctx drift`                | Detect stale context              |
| `ctx sync`                 | Reconcile with codebase           |
| `ctx compact`              | Archive old items                 |
| `ctx watch`                | Watch for update commands         |
| `ctx hook <tool>`          | Generate tool config              |

---

## REFERENCE: PROJECT STRUCTURE

```
active-memory/
├── cmd/
│   └── ctx/
│       └── main.go           # Entry point, Cobra root command
├── internal/
│   ├── cli/                  # Command implementations
│   │   ├── init.go
│   │   ├── status.go
│   │   ├── load.go
│   │   ├── agent.go
│   │   ├── add.go
│   │   ├── complete.go
│   │   ├── drift.go
│   │   ├── sync.go
│   │   ├── compact.go
│   │   ├── watch.go
│   │   └── hook.go
│   ├── context/              # Core context logic
│   │   ├── loader.go
│   │   ├── parser.go
│   │   └── token.go
│   ├── files/                # File type handlers
│   │   └── handlers.go
│   ├── drift/                # Drift detection
│   │   ├── detector.go
│   │   └── rules.go
│   └── templates/            # Embedded templates
│       └── embed.go
├── templates/                # Template source files
│   ├── CONSTITUTION.md
│   ├── TASKS.md
│   └── ... (all template files)
├── hack/
│   └── build-all.sh
├── examples/
│   └── demo/
├── specs/                    # Specifications (read-only reference)
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

Now read the specs and begin.
