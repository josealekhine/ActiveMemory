---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

icon: lucide/git-compare
---

![ctx](images/ctx-banner.png)

## How `ctx` Differs from Similar Tools

There are many tools in the AI ecosystem that touch *parts* of the context
problem:

* Some manage prompts.  
* Some retrieve data.  
* Some provide runtime context objects.  
* Some offer enterprise platforms.

`ctx` focuses on a different layer entirely.

This page explains where `ctx` fits, and where it **intentionally** does not.

---

## The Core Distinction

Most tools treat context as **input**.

`ctx` treats context as **infrastructure**.

That single difference explains nearly all of `ctx`'s design choices.

| Question                 | Most tools                | ctx              |
|--------------------------|---------------------------|------------------|
| Where does context live? | In prompts or APIs        | In files         |
| How long does it last?   | One request / one session | Across time      |
| Who can read it?         | The model                 | Humans and tools |
| How is it updated?       | Implicitly                | Explicitly       |
| Is it inspectable?       | Rarely                    | Always           |

---

## Prompt Management Tools

Examples include:

* prompt templates
* reusable system prompts
* prompt libraries
* prompt versioning tools

These tools help you *start* a session.

They do not help you *continue* one.

Prompt tools:
* inject text at session start
* are ephemeral by design
* do not evolve with the project

`ctx`:
* persists knowledge over time
* accumulates decisions and learnings
* makes the context part of the repository itself

Prompt tooling and `ctx` are complementary; not competing. 
Yet, they operate at different layers.

---

## Retrieval-Augmented Generation (RAG)

RAG systems typically:

* index documents
* embed text
* retrieve chunks dynamically at runtime

They are excellent for:
* large knowledge bases
* static documentation
* reference material

RAG answers questions like:

> “What information might be relevant right now?”

`ctx` answers a different question:

> “What have we already decided, learned, or committed to?”

Here are some key differences:

| RAG                   | ctx                   |
|-----------------------|-----------------------|
| Statistical relevance | Intentional relevance |
| Embedding-based       | File-based            |
| Opaque retrieval      | Explicit structure    |
| Runtime query         | Persistent memory     |

`ctx` does not replace RAG.
Instead, it defines a persistent context layer that RAG can optionally augment.

> RAG belongs to the **data plane**; ctx defines the **context control plane**.

It focuses on **project memory**, not knowledge search.

---

## Agent Frameworks

Agent frameworks often provide:

* task loops
* tool orchestration
* planner/executor patterns
* autonomous iteration

These systems are powerful, but they typically assume that:

* memory is external
* context is injected
* state is transient

Agent frameworks answer:

> “How should the agent act?”

`ctx` answers:

> “What should the agent remember?”

Without persistent context, agents tend to:

* rediscover decisions
* repeat mistakes
* lose architectural intent

This is why `ctx` pairs well with loop-based workflows such as the
[Ralph Wiggum technique](https://ghuntley.com/ralph/):

* **Ralph** provides iteration
* **ctx** provides continuity

Together, loops become cumulative instead of forgetful.

---

## SDK-Level Context Objects

Some SDKs expose "context" objects that exist:

* inside a process
* during a request
* for the lifetime of a call chain

These are extremely useful and completely different.

SDK context objects:

* are in-memory
* disappear when the process ends
* are not shared across sessions

`ctx`:

* survives process restarts
* survives new chats
* survives new days

They share a name, not a purpose.

---

## Enterprise Context Platforms

Enterprise platforms often provide:

* centralized context services
* dashboards
* access control
* organizational knowledge layers

These tools are designed for:

* teams
* governance
* compliance
* managed environments

`ctx` is intentionally:

* local-first
* file-based
* dependency-free
* CLI-driven
* developer-controlled

It does not require:

* a server
* a database
* an account
* a SaaS backend

`ctx` optimizes for *individual and small-team workflows* where context should
live next to code; **not** behind a service boundary.

---

## When `ctx` Is a Good Fit

`ctx` works best when:

* you want AI work to compound over time
* architectural decisions matter
* context must be inspectable
* humans and AI must share the same source of truth
* Git history should include *why*, not just *what*

---

## When `ctx` Is Not the Right Tool

`ctx` is probably not what you want if:

* you only need one-off prompts
* you rely exclusively on RAG
* you want autonomous agents without a human-readable state
* you require centralized enterprise control
* you want black-box memory systems

These are valid goals; just different ones.

---

## Summary

Many tools help AI *think*.

`ctx` helps AI *remember*.

**NOT** by storing thoughts,
**BUT** by preserving intent.
