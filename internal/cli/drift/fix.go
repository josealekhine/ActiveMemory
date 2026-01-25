//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/templates"
)

// fixResult tracks fixes applied during drift fix.
type fixResult struct {
	fixed   int
	skipped int
	errors  []string
}

// applyFixes attempts to auto-fix issues in the drift report.
//
// Currently, supports fixing:
//   - staleness: Archives completed tasks from TASKS.md
//   - missing_file: Creates missing required files from templates
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context
//   - report: Drift report containing issues to fix
//
// Returns:
//   - *fixResult: Summary of fixes applied
func applyFixes(
	cmd *cobra.Command, ctx *context.Context, report *drift.Report,
) *fixResult {
	result := &fixResult{}
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Process warnings (staleness, missing_file, dead_path)
	for _, issue := range report.Warnings {
		switch issue.Type {
		case "staleness":
			if err := fixStaleness(cmd, ctx); err != nil {
				result.errors = append(result.errors,
					fmt.Sprintf("staleness: %v", err))
			} else {
				cmd.Printf("%s Fixed staleness in %s (archived completed tasks)\n",
					green("✓"), issue.File)
				result.fixed++
			}

		case "missing_file":
			if err := fixMissingFile(cmd, issue.File); err != nil {
				result.errors = append(result.errors,
					fmt.Sprintf("missing %s: %v", issue.File, err))
			} else {
				cmd.Printf("%s Created missing file: %s\n", green("✓"), issue.File)
				result.fixed++
			}

		case "dead_path":
			cmd.Printf("%s Cannot auto-fix dead path in %s:%d (%s)\n",
				yellow("○"), issue.File, issue.Line, issue.Path)
			result.skipped++
		}
	}

	// Process violations (potential_secret) - never auto-fix
	for _, issue := range report.Violations {
		if issue.Type == "potential_secret" {
			cmd.Printf("%s Cannot auto-fix potential secret: %s\n",
				yellow("○"), issue.File)
			result.skipped++
		}
	}

	return result
}

// fixStaleness archives completed tasks from TASKS.md.
//
// Moves completed tasks to .context/archive/tasks-YYYY-MM-DD.md and removes
// them from the Completed section in TASKS.md.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context containing the files
//
// Returns:
//   - error: Non-nil if file operations fail
func fixStaleness(cmd *cobra.Command, ctx *context.Context) error {
	var tasksFile *context.FileInfo
	for i := range ctx.Files {
		if ctx.Files[i].Name == config.FilenameTask {
			tasksFile = &ctx.Files[i]
			break
		}
	}

	if tasksFile == nil {
		return fmt.Errorf("TASKS.md not found")
	}

	content := string(tasksFile.Content)
	lines := strings.Split(content, "\n")

	// Find completed tasks in the Completed section
	completedPattern := regexp.MustCompile(`^-\s*\[x]\s*(.+)$`)
	var completedTasks []string
	var newLines []string
	inCompletedSection := false

	for _, line := range lines {
		// Track if we're in the Completed section
		if strings.HasPrefix(line, "## Completed") {
			inCompletedSection = true
			newLines = append(newLines, line)
			continue
		}
		if strings.HasPrefix(line, "## ") && inCompletedSection {
			inCompletedSection = false
		}

		// Collect completed tasks from the Completed section for archiving
		if inCompletedSection && completedPattern.MatchString(line) {
			matches := completedPattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				completedTasks = append(completedTasks, matches[1])
				continue // Remove from file
			}
		}

		newLines = append(newLines, line)
	}

	if len(completedTasks) == 0 {
		return fmt.Errorf("no completed tasks to archive")
	}

	// Create an archive directory
	archiveDir := filepath.Join(config.DirContext, "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Write to the archive file
	archiveFile := filepath.Join(
		archiveDir,
		fmt.Sprintf("tasks-%s.md", time.Now().Format("2006-01-02")),
	)

	archiveContent := fmt.Sprintf(
		"# Archived Tasks - %s\n\n", time.Now().Format("2006-01-02"),
	)
	for _, task := range completedTasks {
		archiveContent += fmt.Sprintf("- [x] %s\n", task)
	}

	// Append to existing archive file if it exists
	if existing, err := os.ReadFile(archiveFile); err == nil {
		archiveContent = string(existing) + "\n" + archiveContent
	}

	if err := os.WriteFile(
		archiveFile, []byte(archiveContent), 0644,
	); err != nil {
		return fmt.Errorf("failed to write archive: %w", err)
	}

	// Write updated TASKS.md
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(
		tasksFile.Path, []byte(newContent), 0644,
	); err != nil {
		return fmt.Errorf("failed to update TASKS.md: %w", err)
	}

	cmd.Printf("  Archived %d completed tasks to %s\n",
		len(completedTasks), archiveFile)

	return nil
}

// fixMissingFile creates a missing required context file from template.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - filename: Name of the file to create (e.g., "CONSTITUTION.md")
//
// Returns:
//   - error: Non-nil if the template not is found or file write fails
func fixMissingFile(cmd *cobra.Command, filename string) error {
	content, err := templates.GetTemplate(filename)
	if err != nil {
		return fmt.Errorf("no template available for %s: %w", filename, err)
	}

	targetPath := filepath.Join(config.DirContext, filename)

	// Ensure .context/ directory exists
	if err := os.MkdirAll(config.DirContext, 0755); err != nil {
		return fmt.Errorf("failed to create .context/: %w", err)
	}

	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetPath, err)
	}

	return nil
}
