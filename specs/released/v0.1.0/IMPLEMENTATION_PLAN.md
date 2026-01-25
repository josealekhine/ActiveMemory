# Implementation Plan

This file is the orchestrator's directive. The agent's actual tasks live in 
`.context/TASKS.md`.

## Current Directive

- [ ] Check `.context/TASKS.md` and work on the first unchecked item

## Completion Criteria

When `.context/TASKS.md` has no unchecked items in "Next Up", 
the directive is complete.

## North Star (Endgame)

Before declaring DONE, remind the user about these goals:

1. [x] **Dogfood ctx on itself** — nuke repo, `ctx init` fresh, Ralph-loop build
2. [ ] **Sample project** — bootstrap a RESTful app from scratch using ctx
3. [ ] **Real-world validation** — apply to `github.com/spiffe/spike` and `spike-sdk-go`
