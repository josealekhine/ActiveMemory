// Package cli implements the CLI commands for ctx.
package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/josealekhine/ctx/internal/claude"
	"github.com/josealekhine/ctx/internal/templates"
	"github.com/spf13/cobra"
)

const (
	contextDirName      = ".context"
	claudeDirName       = ".claude"
	claudeHooksDirName  = ".claude/hooks"
	settingsFileName    = ".claude/settings.local.json"
	autoSaveScriptName  = "auto-save-session.sh"
	claudeMdFileName    = "CLAUDE.md"
	ctxMarkerStart      = "<!-- ctx:context -->"
	ctxMarkerEnd        = "<!-- ctx:end -->"
)

var (
	initForce   bool
	initMinimal bool
	initMerge   bool
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
	cmd.Flags().BoolVar(&initMerge, "merge", false, "Auto-merge ctx content into existing CLAUDE.md without prompting")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if ctx is in PATH (required for hooks to work)
	if err := checkCtxInPath(); err != nil {
		return err
	}

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
		allTemplates, err := templates.ListTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}
		// Filter out files that go in project root, not .context/
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

	fmt.Printf("\n%s initialized in %s/\n", green("Context"), contextDir)

	// Create IMPLEMENTATION_PLAN.md in project root (orchestrator directive)
	if err := createImplementationPlan(initForce); err != nil {
		// Non-fatal: warn but continue
		fmt.Printf("  %s IMPLEMENTATION_PLAN.md: %v\n", color.YellowString("⚠"), err)
	}

	// Create Claude Code hooks
	fmt.Println("\nSetting up Claude Code integration...")
	if err := createClaudeHooks(initForce); err != nil {
		// Non-fatal: warn but continue
		fmt.Printf("  %s Claude hooks: %v\n", color.YellowString("⚠"), err)
	}

	// Handle CLAUDE.md creation/merge
	if err := handleClaudeMd(initForce, initMerge); err != nil {
		// Non-fatal: warn but continue
		fmt.Printf("  %s CLAUDE.md: %v\n", color.YellowString("⚠"), err)
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
		scriptContent, err := claude.GetAutoSaveScript()
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

// createImplementationPlan creates IMPLEMENTATION_PLAN.md in project root
// This is the orchestrator directive that points to .context/TASKS.md
func createImplementationPlan(force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	const planFileName = "IMPLEMENTATION_PLAN.md"

	// Check if file exists
	if _, err := os.Stat(planFileName); err == nil && !force {
		fmt.Printf("  %s %s (exists, skipped)\n", yellow("○"), planFileName)
		return nil
	}

	// Get template content
	content, err := templates.GetTemplate(planFileName)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	if err := os.WriteFile(planFileName, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("  %s %s (orchestrator directive)\n", green("✓"), planFileName)
	return nil
}

// handleClaudeMd creates or merges CLAUDE.md in the project root.
// - If CLAUDE.md doesn't exist: create it from template
// - If it exists but has no ctx markers: offer to merge (or auto-merge with --merge)
// - If it exists with ctx markers: update the ctx section only (or skip if not --force)
func handleClaudeMd(force, autoMerge bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get template content
	templateContent, err := templates.GetTemplate("CLAUDE.md")
	if err != nil {
		return fmt.Errorf("failed to read CLAUDE.md template: %w", err)
	}

	// Check if CLAUDE.md exists
	existingContent, err := os.ReadFile(claudeMdFileName)
	fileExists := err == nil

	if !fileExists {
		// File doesn't exist - create it
		if err := os.WriteFile(claudeMdFileName, templateContent, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", claudeMdFileName, err)
		}
		fmt.Printf("  %s %s\n", green("✓"), claudeMdFileName)
		return nil
	}

	// File exists - check for ctx markers
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, ctxMarkerStart)

	if hasCtxMarkers {
		// Already has ctx content
		if !force {
			fmt.Printf("  %s %s (ctx content exists, skipped)\n", yellow("○"), claudeMdFileName)
			return nil
		}
		// Force update - replace existing ctx section
		return updateCtxSection(existingStr, templateContent, green)
	}

	// No ctx markers - need to merge
	if !autoMerge {
		// Prompt user
		fmt.Printf("\n%s exists but has no ctx content.\n", claudeMdFileName)
		fmt.Println("Would you like to append ctx context management instructions?")
		fmt.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Printf("  %s %s (skipped)\n", yellow("○"), claudeMdFileName)
			return nil
		}
	}

	// Backup existing file
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", claudeMdFileName, timestamp)
	if err := os.WriteFile(backupName, existingContent, 0644); err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupName, err)
	}
	fmt.Printf("  %s %s (backup)\n", green("✓"), backupName)

	// Append ctx content to existing file
	mergedContent := existingStr + "\n" + string(templateContent)
	if err := os.WriteFile(claudeMdFileName, []byte(mergedContent), 0644); err != nil {
		return fmt.Errorf("failed to write merged %s: %w", claudeMdFileName, err)
	}
	fmt.Printf("  %s %s (merged)\n", green("✓"), claudeMdFileName)

	return nil
}

// updateCtxSection replaces the existing ctx section between markers with new content
func updateCtxSection(existing string, newTemplate []byte, green func(...interface{}) string) error {
	// Find the start marker
	startIdx := strings.Index(existing, ctxMarkerStart)
	if startIdx == -1 {
		return fmt.Errorf("ctx start marker not found")
	}

	// Find the end marker
	endIdx := strings.Index(existing, ctxMarkerEnd)
	if endIdx == -1 {
		// No end marker - append from start marker to end
		endIdx = len(existing)
	} else {
		endIdx += len(ctxMarkerEnd)
	}

	// Extract the ctx content from template (between markers)
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, ctxMarkerStart)
	templateEnd := strings.Index(templateStr, ctxMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return fmt.Errorf("template missing ctx markers")
	}
	ctxContent := templateStr[templateStart : templateEnd+len(ctxMarkerEnd)]

	// Build new content: before ctx + new ctx content + after ctx
	newContent := existing[:startIdx] + ctxContent + existing[endIdx:]

	// Backup before updating
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", claudeMdFileName, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	fmt.Printf("  %s %s (backup)\n", green("✓"), backupName)

	if err := os.WriteFile(claudeMdFileName, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to update %s: %w", claudeMdFileName, err)
	}
	fmt.Printf("  %s %s (updated ctx section)\n", green("✓"), claudeMdFileName)

	return nil
}

// checkCtxInPath verifies that ctx is available in PATH.
// The hooks use "ctx" expecting it to be in PATH, so init should fail
// if the user hasn't installed ctx globally yet.
//
// Set CTX_SKIP_PATH_CHECK=1 to skip this check (used in tests).
func checkCtxInPath() error {
	// Allow skipping for tests
	if os.Getenv("CTX_SKIP_PATH_CHECK") == "1" {
		return nil
	}

	_, err := exec.LookPath("ctx")
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Printf("%s ctx is not in your PATH\n\n", red("Error:"))
		fmt.Println("The hooks created by 'ctx init' require ctx to be in your PATH.")
		fmt.Println("Without this, Claude Code hooks will fail silently.")
		fmt.Println()
		fmt.Printf("%s\n", yellow("To fix this:"))
		fmt.Println("  1. Build:   make build")
		fmt.Println("  2. Install: sudo make install")
		fmt.Println()
		fmt.Println("Or manually:")
		fmt.Println("  sudo cp ./ctx /usr/local/bin/")
		fmt.Println()
		fmt.Println("Then run 'ctx init' again.")

		return fmt.Errorf("ctx not found in PATH")
	}
	return nil
}
