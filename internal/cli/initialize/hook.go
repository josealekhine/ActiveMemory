//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/config"
)

// createClaudeHooks creates .claude/hooks/ directory and settings.local.json.
//
// Creates hook scripts (auto-save-session.sh, block-non-path-ctx.sh) and
// merges hooks into existing settings rather than overwriting.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - force: If true, overwrite existing hooks and scripts
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func createClaudeHooks(cmd *cobra.Command, force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get the current working directory for paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create .claude/hooks/ directory
	if err := os.MkdirAll(config.DirClaudeHooks, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", config.DirClaudeHooks, err)
	}

	// Create the auto-save-session.sh script
	scriptPath := filepath.Join(
		config.DirClaudeHooks, config.FileAutoSave,
	)
	if _, err := os.Stat(scriptPath); err == nil && !force {
		cmd.Printf("  %s %s (exists, skipped)\n", yellow("○"), scriptPath)
	} else {
		scriptContent, err := claude.GetAutoSaveScript()
		if err != nil {
			return fmt.Errorf("failed to get auto-save script: %w", err)
		}
		if err := os.WriteFile(scriptPath, scriptContent, 0755); err != nil {
			return fmt.Errorf("failed to write %s: %w", scriptPath, err)
		}
		cmd.Printf("  %s %s\n", green("✓"), scriptPath)
	}

	// Create block-non-path-ctx.sh script
	// (enforces CONSTITUTION.md ctx invocation rules)
	blockScriptPath := filepath.Join(
		config.DirClaudeHooks, config.FileBlockNonPathScript,
	)
	if _, err := os.Stat(blockScriptPath); err == nil && !force {
		cmd.Printf("  %s %s (exists, skipped)\n", yellow("○"), blockScriptPath)
	} else {
		blockScriptContent, err := claude.GetBlockNonPathCtxScript()
		if err != nil {
			return fmt.Errorf("failed to get block-non-path-ctx script: %w", err)
		}
		if err := os.WriteFile(
			blockScriptPath, blockScriptContent, 0755,
		); err != nil {
			return fmt.Errorf("failed to write %s: %w", blockScriptPath, err)
		}
		cmd.Printf("  %s %s\n", green("✓"), blockScriptPath)
	}

	// Handle settings.local.json - merge rather than overwrite
	if err := mergeSettingsHooks(cmd, cwd, force); err != nil {
		return err
	}

	// Create .claude/commands/ directory and ctx skill files
	if err := createClaudeCommands(cmd, force); err != nil {
		return err
	}

	return nil
}

// mergeSettingsHooks creates or merges hooks and permissions into settings.local.json.
//
// Only adds missing hooks and permissions to preserve user customizations.
// Creates the .claude/ directory if needed.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - projectDir: Project root directory for hook paths
//   - force: If true, overwrite existing hooks (permissions are always merged additively)
//
// Returns:
//   - error: Non-nil if JSON parsing or file operations fail
func mergeSettingsHooks(
	cmd *cobra.Command, projectDir string, force bool,
) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if settings.local.json exists
	var settings claude.Settings
	existingContent, err := os.ReadFile(config.FileSettings)
	fileExists := err == nil

	if fileExists {
		if err := json.Unmarshal(existingContent, &settings); err != nil {
			return fmt.Errorf(
				"failed to parse existing %s: %w", config.FileSettings, err,
			)
		}
	}

	// Get our defaults
	defaultHooks := claude.CreateDefaultHooks(projectDir)
	defaultPerms := claude.CreateDefaultPermissions()

	// Check if hooks already exist
	hasPreToolUse := len(settings.Hooks.PreToolUse) > 0
	hasSessionEnd := len(settings.Hooks.SessionEnd) > 0

	// Merge hooks - only add what's missing (or force overwrite)
	hooksModified := false
	if !hasPreToolUse || force {
		settings.Hooks.PreToolUse = defaultHooks.PreToolUse
		hooksModified = true
	}
	if !hasSessionEnd || force {
		settings.Hooks.SessionEnd = defaultHooks.SessionEnd
		hooksModified = true
	}

	// Merge permissions - always additive, never removes existing permissions
	permsModified := mergePermissions(&settings.Permissions, defaultPerms)

	if !hooksModified && !permsModified {
		cmd.Printf(
			"  %s %s (no changes needed)\n", yellow("○"), config.FileSettings,
		)
		return nil
	}

	// Create .claude/ directory if needed
	if err := os.MkdirAll(config.DirClaude, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", config.DirClaude, err)
	}

	// Write settings with pretty formatting (disable HTML escaping to avoid \u003e for >)
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(settings); err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(config.FileSettings, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", config.FileSettings, err)
	}

	// Report what was done
	switch {
	case fileExists && hooksModified && permsModified:
		cmd.Printf("  %s %s (merged hooks and permissions)\n", green("✓"), config.FileSettings)
	case fileExists && hooksModified:
		cmd.Printf("  %s %s (merged hooks)\n", green("✓"), config.FileSettings)
	case fileExists && permsModified:
		cmd.Printf("  %s %s (added ctx permissions)\n", green("✓"), config.FileSettings)
	default:
		cmd.Printf("  %s %s\n", green("✓"), config.FileSettings)
	}

	return nil
}

// mergePermissions adds missing permissions to the allow list.
//
// Only adds permissions that don't already exist. Never removes existing
// permissions to preserve user customizations.
//
// Parameters:
//   - perms: Existing permissions config to modify
//   - defaults: Default permissions to add if missing
//
// Returns:
//   - bool: True if any permissions were added
func mergePermissions(perms *claude.PermissionsConfig, defaults []string) bool {
	// Build a set of existing permissions for fast lookup
	existing := make(map[string]bool)
	for _, p := range perms.Allow {
		existing[p] = true
	}

	// Add missing permissions
	added := false
	for _, p := range defaults {
		if !existing[p] {
			perms.Allow = append(perms.Allow, p)
			added = true
		}
	}

	return added
}
