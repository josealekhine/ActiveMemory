//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agent implements the "ctx agent" command for generating
// AI-ready context packets.
//
// The agent command reads context files from .context/ and produces
// a concise, token-budgeted output optimized for AI consumption.
// Output can be in Markdown (default) or JSON format.
//
// # File Organization
//
//   - agent.go: Command definition and flag registration
//   - run.go: Main execution logic and context loading
//   - extract.go: Functions for extracting content from context files
//   - sort.go: Priority sorting for tasks and decisions
//   - out.go: Output formatting (Markdown and JSON)
//   - types.go: Data structures for context packets
package agent
