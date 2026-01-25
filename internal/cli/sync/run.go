//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/context"
)

// runSync executes the sync command logic.
//
// Loads context, detects discrepancies between codebase and documentation,
// and displays suggested actions. In dry-run mode, only shows what would
// be suggested without prompting for changes.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - dryRun: If true, only show suggestions without prompting for changes
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func runSync(cmd *cobra.Command, dryRun bool) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	actions := detectSyncActions(ctx)

	if len(actions) == 0 {
		green := color.New(color.FgGreen).SprintFunc()
		cmd.Printf("%s Context is in sync with codebase\n", green("✓"))
		return nil
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	cmd.Println(cyan("Sync Analysis"))
	cmd.Println(cyan("============="))
	cmd.Println()

	if dryRun {
		cmd.Println(yellow("DRY RUN — No changes will be made"))
		cmd.Println()
	}

	for i, action := range actions {
		cmd.Printf("%d. [%s] %s\n", i+1, action.Type, action.Description)
		if action.Suggestion != "" {
			cmd.Printf("   Suggestion: %s\n", action.Suggestion)
		}
		cmd.Println()
	}

	if dryRun {
		cmd.Printf(
			"Found %d items to sync. Run without --dry-run to apply suggestions.\n",
			len(actions),
		)
	} else {
		cmd.Printf(
			"Found %d items. Review and update context files manually.\n",
			len(actions),
		)
	}

	return nil
}
