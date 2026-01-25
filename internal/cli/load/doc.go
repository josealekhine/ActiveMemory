//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package load provides the command for outputting assembled context.
//
// The load command assembles context files from .context/ and outputs them
// in the recommended read order, suitable for providing to an AI assistant.
// This is the primary mechanism for giving AI tools access to project context.
//
// # Assembly Order
//
// Context files are assembled in priority order:
//
//  1. CONSTITUTION.md - Hard rules and constraints
//  2. TASKS.md - Current work items
//  3. CONVENTIONS.md - Code patterns and standards
//  4. ARCHITECTURE.md - System design
//  5. DECISIONS.md - Architectural decisions with rationale
//  6. LEARNINGS.md - Gotchas, tips, lessons learned
//  7. GLOSSARY.md - Domain terminology
//  8. DRIFT.md - Drift detection results
//  9. AGENT_PLAYBOOK.md - AI-specific instructions
//
// # Token Budget
//
// The --budget flag limits output to approximately the specified token count.
// This is useful for AI assistants with context window limitations. Files are
// prioritized by importance, with lower-priority files truncated or omitted
// when budget constraints are reached.
//
// # Raw Output
//
// The --raw flag outputs file contents directly without assembly headers or
// priority-based ordering. This is useful for debugging or when exact file
// contents are needed.
//
// # File Organization
//
//   - load.go: Command definition and flag handling
//   - run.go: Main load execution logic
//   - convert.go: Format conversion utilities
//   - sort.go: Priority-based file sorting
//   - out.go: Output formatting
package load
