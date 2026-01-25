//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/templates"
)

// createEntryTemplates creates .context/templates/ with entry templates for
// rich entries.
//
// These templates help users format detailed learnings and decisions using
// the --file flag with ctx add commands.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - contextDir: Path to the .context/ directory
//   - force: If true, overwrite existing templates
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func createEntryTemplates(
	cmd *cobra.Command, contextDir string, force bool,
) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	templatesDir := filepath.Join(contextDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", templatesDir, err)
	}

	// Get list of entry templates
	entryTemplates, err := templates.ListEntryTemplates()
	if err != nil {
		return fmt.Errorf("failed to list entry templates: %w", err)
	}

	for _, name := range entryTemplates {
		targetPath := filepath.Join(templatesDir, name)

		// Check if the file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Printf("  %s templates/%s (exists, skipped)\n", yellow("○"), name)
			continue
		}

		content, err := templates.GetEntryTemplate(name)
		if err != nil {
			return fmt.Errorf("failed to read entry template %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		cmd.Printf("  %s templates/%s\n", green("✓"), name)
	}

	return nil
}
