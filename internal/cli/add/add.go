//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx add" command for appending entries to context files.
//
// Supported types are defined in [config.FileType] (both singular and plural
// forms accepted, e.g., "decision" or "decisions"). Content can be provided
// via command argument, --file flag, or stdin pipe.
//
// Flags:
//   - --priority, -p: Priority level for tasks (high, medium, low)
//   - --section, -s: Target section within the file
//   - --file, -f: Read content from a file instead of argument
//
// Returns:
//   - *cobra.Command: Configured add command with flags registered
func Cmd() *cobra.Command {
	var (
		priority string
		section  string
		fromFile string
	)

	cmd := &cobra.Command{
		Use:   "add <type> [content]",
		Short: "Add a new item to a context file",
		Long: `Add a new decision, task, learning, or convention 
to the appropriate context file.

Types:
  decision    Add to DECISIONS.md
  task        Add to TASKS.md
  learning    Add to LEARNINGS.md
  convention  Add to CONVENTIONS.md

Content can be provided as:
  - Command argument: ctx add learning "text here"
  - File: ctx add learning --file /path/to/content.md
  - Stdin: echo "text" | ctx add learning

Examples:
  ctx add decision "Use PostgreSQL for primary database"
  ctx add task "Implement user authentication" --priority high
  ctx add learning "Vitest mocks must be hoisted"
  ctx add learning --file learning-template.md
  cat notes.md | ctx add decision`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args, priority, section, fromFile)
		},
	}

	cmd.Flags().StringVarP(
		&priority,
		"priority", "p", "",
		"Priority level for tasks (high, medium, low)",
	)
	cmd.Flags().StringVarP(
		&section,
		"section", "s", "",
		"Target section within file",
	)
	cmd.Flags().StringVarP(
		&fromFile,
		"file", "f", "",
		"Read content from file instead of argument",
	)

	return cmd
}
