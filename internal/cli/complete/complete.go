//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx complete" command for marking tasks as done.
//
// Tasks can be specified by number, partial text match, or full text.
// The command updates TASKS.md by changing "- [ ]" to "- [x]".
//
// Returns:
//   - *cobra.Command: Configured complete command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete <task-id-or-text>",
		Short: "Mark a task as completed",
		Long: `Mark a task as completed in TASKS.md.

You can specify a task by:
  - Task number (e.g., "ctx complete 3")
  - Partial text match (e.g., "ctx complete auth")
  - Full task text (e.g., "ctx complete 'Implement user authentication'")

The task will be marked with [x] 
and optionally moved to the Completed section.`,
		Args: cobra.ExactArgs(1),
		RunE: runComplete,
	}

	return cmd
}
