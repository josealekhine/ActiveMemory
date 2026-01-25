//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task implements the "ctx tasks" command for managing task archival
// and snapshots.
//
// The task package provides subcommands to:
//   - archive: Move completed tasks to timestamped archive files
//   - snapshot: Create point-in-time copies of TASKS.md
//
// Archive files preserve phase structure for traceability, while snapshots
// copy the entire file as-is without modification.
package task

import (
	"github.com/spf13/cobra"
)

// Cmd returns the tasks command with subcommands.
//
// The tasks command provides utilities for managing the task lifecycle:
//   - archive: Move completed tasks out of TASKS.md
//   - snapshot: Create point-in-time backup without modification
//
// Returns:
//   - *cobra.Command: Configured tasks command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Manage task archival and snapshots",
		Long: `Manage task archival and snapshots.

Tasks can be archived to move completed items out of TASKS.md while
preserving them for historical reference. Snapshots create point-in-time
copies without modifying the original.

Subcommands:
  archive   Move completed tasks to timestamped archive file
  snapshot  Create point-in-time snapshot of TASKS.md`,
	}

	cmd.AddCommand(archiveCmd())
	cmd.AddCommand(snapshotCmd())

	return cmd
}

// archiveCmd returns the tasks archive subcommand.
//
// The archive command moves completed tasks (marked with [x]) from TASKS.md
// to a timestamped archive file in .context/archive/. Pending tasks ([ ])
// remain in TASKS.md.
//
// Flags:
//   - --dry-run: Preview changes without modifying files
//
// Returns:
//   - *cobra.Command: Configured archive subcommand
func archiveCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Move completed tasks to timestamped archive file",
		Long: `Move completed tasks from TASKS.md to an archive file.

Archive files are stored in .context/archive/ with timestamped names:
  .context/archive/tasks-YYYY-MM-DD.md

The archive preserves Phase structure for traceability. Completed tasks
(marked with [x]) are moved; pending tasks ([ ]) remain in TASKS.md.

Use --dry-run to preview changes without modifying files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTaskArchive(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun,
		"dry-run",
		false,
		"Preview changes without modifying files",
	)

	return cmd
}

// snapshotCmd returns the tasks snapshot subcommand.
//
// The snapshot command creates a point-in-time copy of TASKS.md without
// modifying the original. Snapshots are stored in .context/archive/ with
// timestamped names.
//
// Arguments:
//   - [name]: Optional name for the snapshot (defaults to "snapshot")
//
// Returns:
//   - *cobra.Command: Configured snapshot subcommand
func snapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot [name]",
		Short: "Create point-in-time snapshot of TASKS.md",
		Long: `Create a point-in-time snapshot of TASKS.md without modifying the original.

Snapshots are stored in .context/archive/ with timestamped names:
  .context/archive/tasks-snapshot-YYYY-MM-DD-HHMM.md

Unlike archive, snapshot copies the entire file as-is.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runTasksSnapshot,
	}

	return cmd
}
