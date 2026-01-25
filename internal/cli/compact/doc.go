//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compact implements the "ctx compact" command for cleaning up
// and consolidating context files.
//
// The compact command performs maintenance on .context/ files including
// moving completed tasks to a dedicated section, optionally archiving
// old content, and removing empty sections.
//
// # File Organization
//
//   - compact.go: Command definition and flag registration
//   - run.go: Main execution logic and orchestration
//   - process.go: File processing and section manipulation
//   - task.go: Task extraction and completion detection
//   - sanitize.go: Content cleaning and normalization
package compact
