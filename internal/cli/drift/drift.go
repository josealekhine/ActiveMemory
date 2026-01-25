//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx drift" command for detecting stale context.
//
// The command checks for broken path references, staleness indicators,
// constitution violations, and missing required files.
//
// Flags:
//   - --json: Output results as JSON for machine parsing
//   - --fix: Auto-fix supported issues (staleness, missing_file)
//
// Returns:
//   - *cobra.Command: Configured drift command with flags registered
func Cmd() *cobra.Command {
	var (
		jsonOutput bool
		fix        bool
	)

	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect stale or invalid context",
		Long: `Run drift detection to find stale paths,
broken references, and constitution violations.

Checks performed:
  - Path references in ARCHITECTURE.md and CONVENTIONS.md exist
  - Staleness indicators (many completed tasks)
  - Constitution rule violations (potential secrets)
  - Required files are present

Use --json for machine-readable output.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDrift(cmd, jsonOutput, fix)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&fix,
		"fix", false, "Auto-fix supported issues (staleness, missing files)",
	)

	return cmd
}
