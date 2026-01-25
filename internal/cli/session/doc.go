//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package session provides commands for managing session snapshots.
//
// Sessions capture the state of context at a point in time, including current
// tasks, recent decisions, and learnings. The package supports saving sessions
// to Markdown files, listing saved sessions, loading previous sessions, and
// parsing Claude Code JSONL transcripts.
//
// # Commands
//
// The session command provides four subcommands:
//
//   - save: Save current context state to a session file in .context/sessions/
//   - list: List saved sessions with summaries, sorted by date
//   - load: Load and display a previous session by name, partial match, or index
//   - parse: Convert Claude Code JSONL transcripts to readable Markdown
//
// # Session Files
//
// Session files are stored in .context/sessions/ with timestamped filenames
// (YYYY-MM-DD-HHMMSS-topic.md). Each file contains:
//
//   - Metadata (date, time, type, start_time, end_time)
//   - Summary section for describing what was accomplished
//   - Current tasks from TASKS.md
//   - Recent decisions from DECISIONS.md
//   - Recent learnings from LEARNINGS.md
//   - Tasks for next session
//
// # Transcript Parsing
//
// The parse subcommand converts Claude Code JSONL transcript files into
// readable Markdown. It supports:
//
//   - Full transcript conversion with message headers and timestamps
//   - Extraction mode (--extract) for identifying potential decisions and
//     learnings
//   - Output to file (-o) or stdout
package session
