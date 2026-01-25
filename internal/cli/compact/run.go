//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// runCompact executes the compact command logic.
//
// Loads context, optionally saves a pre-compact snapshot, processes
// TASKS.md for completed tasks, and removes empty sections from all
// context files.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - archive: If true, archive old completed tasks to .context/archive/
//   - noAutoSave: If true, skip pre-compact session snapshot
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func runCompact(cmd *cobra.Command, archive, noAutoSave bool) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Auto-save session before compact
	if !noAutoSave {
		if err := preCompactAutoSave(cmd); err != nil {
			cmd.Printf(
				"%s Auto-save failed: %v (continuing anyway)\n", yellow("⚠"), err,
			)
		}
	}

	cmd.Println(cyan("Compact Analysis"))
	cmd.Println(cyan("================"))
	cmd.Println()

	changes := 0

	// Process TASKS.md
	tasksChanges, err := compactTasks(cmd, ctx, archive)
	if err != nil {
		cmd.Printf("%s Error processing TASKS.md: %v\n", yellow("⚠"), err)
	} else {
		changes += tasksChanges
	}

	// Process other files for empty sections
	for _, f := range ctx.Files {
		if f.Name == config.FilenameTask {
			continue
		}
		cleaned, count := removeEmptySections(string(f.Content))
		if count > 0 {
			if err := os.WriteFile(f.Path, []byte(cleaned), 0644); err == nil {
				cmd.Printf("%s Removed %d empty sections from %s\n", green("✓"), count, f.Name)
				changes += count
			}
		}
	}

	if changes == 0 {
		cmd.Printf("%s Nothing to compact — context is already clean\n", green("✓"))
	} else {
		cmd.Printf("\n%s Compacted %d items\n", green("✓"), changes)
	}

	return nil
}
