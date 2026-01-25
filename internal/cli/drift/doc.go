//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift implements the "ctx drift" command for detecting stale
// or invalid context.
//
// The drift command checks for broken path references, staleness indicators,
// constitution violations, and missing required files. Results can be
// output as formatted text or JSON.
//
// # File Organization
//
//   - drift.go: Command definition and flag registration
//   - run.go: Main execution logic and context loading
//   - out.go: Output formatting (text and JSON)
//   - types.go: Data structures for JSON output
//   - sanitize.go: Check name formatting utilities
package drift
