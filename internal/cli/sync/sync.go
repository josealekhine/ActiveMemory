//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx sync" command for reconciling context with codebase.
//
// The command scans the codebase for changes that should be reflected in
// context files, such as new directories, package manager files, and
// configuration files.
//
// Flags:
//   - --dry-run: Show what would change without modifying files
//
// Returns:
//   - *cobra.Command: Configured sync command with flags registered
func Cmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Reconcile context with codebase",
		Long: `Scan the codebase and reconcile context files with current state.

Actions performed:
  - Scan for new directories that should be in ARCHITECTURE.md
  - Check for package.json/go.mod changes
  - Identify stale references
  - Suggest updates to context files

Use --dry-run to see what would change without modifying files.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSync(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun,
		"dry-run", false, "Show what would change without modifying",
	)

	return cmd
}
