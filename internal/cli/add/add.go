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
//   - --context, -c: Context for decisions/learnings (required)
//   - --rationale, -r: Rationale for decisions (required for decisions)
//   - --consequences: Consequences for decisions (required for decisions)
//   - --lesson, -l: Lesson for learnings (required for learnings)
//   - --application, -a: Application for learnings (required for learnings)
//
// Returns:
//   - *cobra.Command: Configured add command with flags registered
func Cmd() *cobra.Command {
	var (
		priority     string
		section      string
		fromFile     string
		context      string
		rationale    string
		consequences string
		lesson       string
		application  string
	)

	cmd := &cobra.Command{
		Use:   "add <type> [content]",
		Short: "Add a new item to a context file",
		Long: `Add a new decision, task, learning, or convention
to the appropriate context file.

Types:
  decision    Add to DECISIONS.md (requires --context, --rationale, --consequences)
  learning    Add to LEARNINGS.md (requires --context, --lesson, --application)
  task        Add to TASKS.md
  convention  Add to CONVENTIONS.md

Content can be provided as:
  - Command argument: ctx add learning "title here"
  - File: ctx add learning --file /path/to/content.md
  - Stdin: echo "title" | ctx add learning

Examples:
  ctx add decision "Use PostgreSQL" \
    --context "Need a reliable database for production" \
    --rationale "PostgreSQL offers ACID compliance and JSON support" \
    --consequences "Team needs PostgreSQL training"
  ctx add learning "Go embed requires files in same package" \
    --context "Tried to embed files from parent directory" \
    --lesson "go:embed only works with files in same or child directories" \
    --application "Keep embedded files in internal/templates/, not project root"
  ctx add task "Implement user authentication" --priority high`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args, addFlags{
				priority:     priority,
				section:      section,
				fromFile:     fromFile,
				context:      context,
				rationale:    rationale,
				consequences: consequences,
				lesson:       lesson,
				application:  application,
			})
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
	cmd.Flags().StringVarP(
		&context,
		"context", "c", "",
		"Context for decisions: what prompted this decision (required for decisions)",
	)
	cmd.Flags().StringVarP(
		&rationale,
		"rationale", "r", "",
		"Rationale for decisions: why this choice over alternatives (required for decisions)",
	)
	cmd.Flags().StringVar(
		&consequences,
		"consequences", "",
		"Consequences for decisions: what changes as a result (required for decisions)",
	)
	cmd.Flags().StringVarP(
		&lesson,
		"lesson", "l", "",
		"Lesson for learnings: the key insight (required for learnings)",
	)
	cmd.Flags().StringVarP(
		&application,
		"application", "a", "",
		"Application for learnings: how to apply this going forward (required for learnings)",
	)

	return cmd
}

// addFlags holds all flags for the add command.
type addFlags struct {
	priority     string
	section      string
	fromFile     string
	context      string
	rationale    string
	consequences string
	lesson       string
	application  string
}
