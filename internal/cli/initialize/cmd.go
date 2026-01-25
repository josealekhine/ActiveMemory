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

	"github.com/ActiveMemory/ctx/internal/claude"
)

// createClaudeCommands creates .claude/commands/ with ctx skill files.
//
// Copies embedded command files to the .claude/commands/ directory for
// use as Claude Code slash commands.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - force: If true, overwrite existing command files
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func createClaudeCommands(cmd *cobra.Command, force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	commandsDir := ".claude/commands"
	if err := os.MkdirAll(commandsDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", commandsDir, err)
	}

	// Get the list of embedded command files
	commands, err := claude.ListCommands()
	if err != nil {
		return fmt.Errorf("failed to list commands: %w", err)
	}

	for _, cmdName := range commands {
		cmdPath := filepath.Join(commandsDir, cmdName)
		if _, err := os.Stat(cmdPath); err == nil && !force {
			cmd.Printf("  %s %s (exists, skipped)\n", yellow("○"), cmdPath)
			continue
		}

		content, err := claude.GetCommand(cmdName)
		if err != nil {
			return fmt.Errorf("failed to get command %s: %w", cmdName, err)
		}

		if err := os.WriteFile(cmdPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", cmdPath, err)
		}
		cmd.Printf("  %s %s\n", green("✓"), cmdPath)
	}

	return nil
}
