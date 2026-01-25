//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx load" command for outputting assembled context.
//
// The command loads context files from .context/ and outputs them in the
// recommended read order, suitable for providing to an AI assistant.
//
// Flags:
//   - --budget: Token budget for assembly (default 8000)
//   - --raw: Output raw file contents without headers or assembly
//
// Returns:
//   - *cobra.Command: Configured load command with flags registered
func Cmd() *cobra.Command {
	var (
		budget int
		raw    bool
	)

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Output assembled context markdown",
		Long: `Load and display the assembled context 
as it would be provided to an AI.

The context files are assembled in the recommended read order:
  1. CONSTITUTION.md
  2. TASKS.md
  3. CONVENTIONS.md
  4. ARCHITECTURE.md
  5. DECISIONS.md
  6. LEARNINGS.md
  7. GLOSSARY.md
  8. DRIFT.md
  9. AGENT_PLAYBOOK.md

Use --raw to output raw file contents without headers or assembly.
Use --budget to limit output to a specific token count.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoad(cmd, budget, raw)
		},
	}

	cmd.Flags().IntVar(
		&budget, "budget", 8000, "Token budget for assembly",
	)
	cmd.Flags().BoolVar(
		&raw, "raw", false, "Output raw file contents without assembly",
	)

	return cmd
}
