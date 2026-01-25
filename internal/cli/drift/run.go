//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
)

// runDrift executes the drift command logic.
//
// Loads context, runs drift detection, and outputs results in the
// specified format. When `fix` is true, attempts to auto-fix supported
// issue types (staleness, missing_file).
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON; otherwise output as text
//   - fix: If true, attempt to auto-fix supported issues
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func runDrift(cmd *cobra.Command, jsonOutput, fix bool) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	report := drift.Detect(ctx)

	// Apply fixes if requested
	if fix && (len(report.Warnings) > 0 || len(report.Violations) > 0) {
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		cmd.Println("Applying fixes...")
		cmd.Println()

		result := applyFixes(cmd, ctx, report)

		cmd.Println()
		if result.fixed > 0 {
			cmd.Printf("%s Fixed %d issue(s)\n", green("✓"), result.fixed)
		}
		if result.skipped > 0 {
			cmd.Printf("%s Skipped %d issue(s) (cannot auto-fix)\n",
				yellow("○"), result.skipped)
		}
		for _, errMsg := range result.errors {
			cmd.Printf("%s Error: %s\n", yellow("⚠"), errMsg)
		}

		// Re-run detection to show updated status
		if result.fixed > 0 {
			cmd.Println()
			cmd.Println("Re-checking after fixes...")
			ctx, _ = context.Load("")
			report = drift.Detect(ctx)
		}
	}

	if jsonOutput {
		return outputDriftJSON(cmd, report)
	}

	return outputDriftText(cmd, report)
}
