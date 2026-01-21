// Package cli implements the CLI commands for ctx.
package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/josealekhine/ActiveMemory/internal/claude"
	"github.com/josealekhine/ActiveMemory/internal/templates"
	"github.com/spf13/cobra"
)

const (
	contextDirName      = ".context"
	claudeDirName       = ".claude"
	claudeHooksDirName  = ".claude/hooks"
	settingsFileName    = ".claude/settings.local.json"
	autoSaveScriptName  = "auto-save-session.sh"
)

var (
	initForce   bool
	initMinimal bool
)

// minimalTemplates are the essential files created with --minimal flag
var minimalTemplates = []string{
	"TASKS.md",
	"DECISIONS.md",
	"CONSTITUTION.md",
}

// InitCmd returns the init command.
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new .context/ directory with template files",
		Long: `Initialize a new .context/ directory with template files for
maintaining persistent context for AI coding assistants.

The following files are created:
  - CONSTITUTION.md  — Hard invariants that must never be violated
  - TASKS.md         — Current and planned work
  - DECISIONS.md     — Architectural decisions with rationale
  - LEARNINGS.md     — Lessons learned, gotchas, tips
  - CONVENTIONS.md   — Project patterns and standards
  - ARCHITECTURE.md  — System overview
  - GLOSSARY.md      — Domain terms and abbreviations
  - DRIFT.md         — Staleness signals and update triggers
  - AGENT_PLAYBOOK.md — How AI agents should use this system

Use --minimal to only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md).`,
		RunE: runInit,
	}

	cmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing context files")
	cmd.Flags().BoolVarP(&initMinimal, "minimal", "m", false, "Only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md)")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	contextDir := contextDirName

	// Check if .context/ already exists
	if _, err := os.Stat(contextDir); err == nil {
		if !initForce {
			// Prompt for confirmation
			fmt.Printf("%s already exists. Overwrite? [y/N] ", contextDir)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}
	}

	// Create .context/ directory
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", contextDir, err)
	}

	// Get list of templates to create
	var templatesToCreate []string
	if initMinimal {
		templatesToCreate = minimalTemplates
	} else {
		var err error
		templatesToCreate, err = templates.ListTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}
	}

	// Create template files
	green := color.New(color.FgGreen).SprintFunc()
	for _, name := range templatesToCreate {
		targetPath := filepath.Join(contextDir, name)

		// Check if file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !initForce {
			fmt.Printf("  %s %s (exists, skipped)\n", color.YellowString("○"), name)
			continue
		}

		content, err := templates.GetTemplate(name)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		fmt.Printf("  %s %s\n", green("✓"), name)
	}

	fmt.Printf("\n%s initialized in %s/\n", green("Active Memory"), contextDir)

	// Create Claude Code hooks
	fmt.Println("\nSetting up Claude Code integration...")
	if err := createClaudeHooks(initForce); err != nil {
		// Non-fatal: warn but continue
		fmt.Printf("  %s Claude hooks: %v\n", color.YellowString("⚠"), err)
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit .context/TASKS.md to add your current tasks")
	fmt.Println("  2. Run 'ctx status' to see context summary")
	fmt.Println("  3. Run 'ctx agent' to get AI-ready context packet")

	return nil
}

// createClaudeHooks creates .claude/hooks/ directory and settings.local.json
// It merges hooks into existing settings rather than overwriting.
func createClaudeHooks(force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get current working directory for paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create .claude/hooks/ directory
	if err := os.MkdirAll(claudeHooksDirName, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", claudeHooksDirName, err)
	}

	// Create auto-save-session.sh script
	scriptPath := filepath.Join(claudeHooksDirName, autoSaveScriptName)
	if _, err := os.Stat(scriptPath); err == nil && !force {
		fmt.Printf("  %s %s (exists, skipped)\n", yellow("○"), scriptPath)
	} else {
		scriptContent, err := claude.GetAutoSaveScript(cwd)
		if err != nil {
			return fmt.Errorf("failed to get auto-save script: %w", err)
		}
		if err := os.WriteFile(scriptPath, scriptContent, 0755); err != nil {
			return fmt.Errorf("failed to write %s: %w", scriptPath, err)
		}
		fmt.Printf("  %s %s\n", green("✓"), scriptPath)
	}

	// Handle settings.local.json - merge rather than overwrite
	if err := mergeSettingsHooks(cwd, force); err != nil {
		return err
	}

	return nil
}

// mergeSettingsHooks creates or merges hooks into settings.local.json
func mergeSettingsHooks(projectDir string, force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if settings.local.json exists
	var settings claude.Settings
	existingContent, err := os.ReadFile(settingsFileName)
	fileExists := err == nil

	if fileExists {
		if err := json.Unmarshal(existingContent, &settings); err != nil {
			return fmt.Errorf("failed to parse existing %s: %w", settingsFileName, err)
		}
	}

	// Get our default hooks
	defaultHooks := claude.CreateDefaultHooks(projectDir)

	// Check if hooks already exist
	hasPreToolUse := len(settings.Hooks.PreToolUse) > 0
	hasSessionEnd := len(settings.Hooks.SessionEnd) > 0

	if fileExists && hasPreToolUse && hasSessionEnd && !force {
		fmt.Printf("  %s %s (hooks exist, skipped)\n", yellow("○"), settingsFileName)
		return nil
	}

	// Merge hooks - only add what's missing
	modified := false
	if !hasPreToolUse || force {
		settings.Hooks.PreToolUse = defaultHooks.PreToolUse
		modified = true
	}
	if !hasSessionEnd || force {
		settings.Hooks.SessionEnd = defaultHooks.SessionEnd
		modified = true
	}

	if !modified {
		fmt.Printf("  %s %s (no changes needed)\n", yellow("○"), settingsFileName)
		return nil
	}

	// Create .claude/ directory if needed
	if err := os.MkdirAll(claudeDirName, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", claudeDirName, err)
	}

	// Write settings with pretty formatting
	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsFileName, output, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", settingsFileName, err)
	}

	if fileExists {
		fmt.Printf("  %s %s (merged hooks)\n", green("✓"), settingsFileName)
	} else {
		fmt.Printf("  %s %s\n", green("✓"), settingsFileName)
	}

	return nil
}
