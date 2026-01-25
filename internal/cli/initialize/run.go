//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/templates"
)

// runInit executes the init command logic.
//
// Creates a .context/ directory with template files. Handles existing
// directories, minimal mode, and CLAUDE.md merge operations.
//
// Parameters:
//   - cmd: Cobra command for output and input streams
//   - force: If true, overwrite existing files without prompting
//   - minimal: If true, only create essential files
//   - merge: If true, auto-merge ctx content into existing CLAUDE.md
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func runInit(cmd *cobra.Command, force, minimal, merge bool) error {
	// Check if ctx is in PATH (required for hooks to work)
	if err := checkCtxInPath(cmd); err != nil {
		return err
	}

	contextDir := config.DirContext

	// Check if .context/ already exists
	if _, err := os.Stat(contextDir); err == nil {
		if !force {
			// Prompt for confirmation
			cmd.Printf("%s already exists. Overwrite? [y/N] ", contextDir)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				cmd.Println("Aborted.")
				return nil
			}
		}
	}

	// Create .context/ directory
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", contextDir, err)
	}

	// Get the list of templates to create
	var templatesToCreate []string
	if minimal {
		templatesToCreate = config.RequiredFiles
	} else {
		allTemplates, err := templates.ListTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}
		// Filter out files that go in the project root, not .context/
		for _, t := range allTemplates {
			if t != "IMPLEMENTATION_PLAN.md" && t != "CLAUDE.md" {
				templatesToCreate = append(templatesToCreate, t)
			}
		}
	}

	// Create template files
	green := color.New(color.FgGreen).SprintFunc()
	for _, name := range templatesToCreate {
		targetPath := filepath.Join(contextDir, name)

		// Check if the file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Printf(
				"  %s %s (exists, skipped)\n", color.YellowString("○"), name,
			)
			continue
		}

		content, err := templates.GetTemplate(name)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		cmd.Printf("  %s %s\n", green("✓"), name)
	}

	cmd.Printf("\n%s initialized in %s/\n", green("Context"), contextDir)

	// Create entry templates in .context/templates/
	if err := createEntryTemplates(cmd, contextDir, force); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s Entry templates: %v\n", color.YellowString("⚠"), err)
	}

	// Create IMPLEMENTATION_PLAN.md in the project root (orchestrator directive)
	if err := createImplementationPlan(cmd, force); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf(
			"  %s IMPLEMENTATION_PLAN.md: %v\n", color.YellowString("⚠"), err,
		)
	}

	// Create Claude Code hooks
	cmd.Println("\nSetting up Claude Code integration...")
	if err := createClaudeHooks(cmd, force); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s Claude hooks: %v\n", color.YellowString("⚠"), err)
	}

	// Handle CLAUDE.md creation/merge
	if err := handleClaudeMd(cmd, force, merge); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s CLAUDE.md: %v\n", color.YellowString("⚠"), err)
	}

	cmd.Println("\nNext steps:")
	cmd.Println("  1. Edit .context/TASKS.md to add your current tasks")
	cmd.Println("  2. Run 'ctx status' to see context summary")
	cmd.Println("  3. Run 'ctx agent' to get AI-ready context packet")

	return nil
}
