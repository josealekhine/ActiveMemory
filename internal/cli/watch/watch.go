//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"github.com/spf13/cobra"
)

var (
	watchLog      string
	watchDryRun   bool
	watchAutoSave bool
)

// Cmd returns the watch command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for context-update commands in AI output",
		Long: `Watch stdin or a log file for <context-update> 
commands and apply them.

This command parses AI output looking for structured update commands:

  <context-update type="task">Implement user auth</context-update>
  <context-update type="decision">Use PostgreSQL</context-update>
  <context-update type="learning">Mock functions must be hoisted</context-update>
  <context-update type="complete">user auth</context-update>

Use --log to watch a specific file instead of stdin.
Use --dry-run to see what would be updated without making changes.
Use --auto-save to periodically save session snapshots (every 5 updates).

Press Ctrl+C to stop watching.`,
		RunE: runWatch,
	}

	cmd.Flags().StringVar(
		&watchLog, "log", "", "Log file to watch (default: stdin)",
	)
	cmd.Flags().BoolVar(
		&watchDryRun, "dry-run", false, "Show updates without applying",
	)
	cmd.Flags().BoolVar(
		&watchAutoSave, "auto-save", false, "Save session snapshots periodically",
	)

	return cmd
}
