//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	compactArchive    bool
	compactNoAutoSave bool
)

// CompactCmd returns the compact command.
func CompactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compact",
		Short: "Archive completed tasks and clean up context",
		Long: `Consolidate and clean up context files.

Actions performed:
  - Move completed tasks to "Completed (Recent)" section
  - Archive old completed tasks (with --archive)
  - Remove empty sections from context files
  - Report on potential duplicates

Use --archive to create .context/archive/ for old content.

Examples:
  ctx compact                  # Clean up context, move completed tasks
  ctx compact --archive        # Also archive old tasks to .context/archive/
  ctx compact --no-auto-save   # Skip pre-compact session snapshot`,
		RunE: runCompact,
	}

	cmd.Flags().BoolVar(
		&compactArchive,
		"archive",
		false,
		"Create .context/archive/ for old content",
	)
	cmd.Flags().BoolVar(
		&compactNoAutoSave,
		"no-auto-save",
		false,
		"Skip auto-saving session before compact",
	)

	return cmd
}

func runCompact(cmd *cobra.Command, _ []string) error {
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
	if !compactNoAutoSave {
		if err := preCompactAutoSave(cmd); err != nil {
			cmd.Printf("%s Auto-save failed: %v (continuing anyway)\n", yellow("⚠"), err)
		}
	}

	cmd.Println(cyan("Compact Analysis"))
	cmd.Println(cyan("================"))
	cmd.Println()

	changes := 0

	// Process TASKS.md
	tasksChanges, err := compactTasks(cmd, ctx, compactArchive)
	if err != nil {
		cmd.Printf("%s Error processing TASKS.md: %v\n", yellow("⚠"), err)
	} else {
		changes += tasksChanges
	}

	// Process other files for empty sections
	for _, f := range ctx.Files {
		if f.Name == "TASKS.md" {
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

func compactTasks(
		cmd *cobra.Command, ctx *context.Context, archive bool,
) (int, error) {
	var tasksFile *context.FileInfo
	for i := range ctx.Files {
		if ctx.Files[i].Name == "TASKS.md" {
			tasksFile = &ctx.Files[i]
			break
		}
	}

	if tasksFile == nil {
		return 0, nil
	}

	content := string(tasksFile.Content)
	lines := strings.Split(content, "\n")

	completedPattern := regexp.MustCompile(`^-\s*\[x]\s*(.+)$`)

	var completedTasks []string
	var newLines []string
	inCompletedSection := false
	changes := 0

	green := color.New(color.FgGreen).SprintFunc()

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

		// If completed task outside the Completed section, collect it
		if !inCompletedSection && completedPattern.MatchString(line) {
			matches := completedPattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				completedTasks = append(completedTasks, matches[1])
				cmd.Printf(
					"%s Moving completed task: %s\n", green("✓"),
					truncateString(matches[1], 50),
				)
				changes++
				continue // Don't add to newLines
			}
		}

		newLines = append(newLines, line)
	}

	// If we have completed tasks to move, add them to the Completed section
	if len(completedTasks) > 0 {
		// Find the Completed section and add tasks there
		for i, line := range newLines {
			if strings.HasPrefix(line, "## Completed") {
				// Find the next line that's either empty or another section
				insertIdx := i + 1
				for insertIdx < len(newLines) && newLines[insertIdx] != "" &&
						!strings.HasPrefix(newLines[insertIdx], "## ") {
					insertIdx++
				}

				// Insert completed tasks
				var tasksToInsert []string
				for _, task := range completedTasks {
					tasksToInsert = append(tasksToInsert, fmt.Sprintf("- [x] %s", task))
				}

				// Insert at the right position
				newContent := append(newLines[:insertIdx], append(tasksToInsert, newLines[insertIdx:]...)...)
				newLines = newContent
				break
			}
		}
	}

	// Archive old content if requested
	if archive && len(completedTasks) > 0 {
		archiveDir := filepath.Join(contextDirName, "archive")
		if err := os.MkdirAll(archiveDir, 0755); err == nil {
			archiveFile := filepath.Join(archiveDir, fmt.Sprintf("tasks-%s.md", time.Now().Format("2006-01-02")))
			archiveContent := fmt.Sprintf("# Archived Tasks - %s\n\n", time.Now().Format("2006-01-02"))
			for _, task := range completedTasks {
				archiveContent += fmt.Sprintf("- [x] %s\n", task)
			}
			if err := os.WriteFile(archiveFile, []byte(archiveContent), 0644); err == nil {
				cmd.Printf("%s Archived %d tasks to %s\n", green("✓"), len(completedTasks), archiveFile)
			}
		}
	}

	// Write back
	newContent := strings.Join(newLines, "\n")
	if newContent != content {
		if err := os.WriteFile(tasksFile.Path, []byte(newContent), 0644); err != nil {
			return 0, err
		}
	}

	return changes, nil
}

func removeEmptySections(content string) (string, int) {
	lines := strings.Split(content, "\n")
	var result []string
	removed := 0

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this is a section header
		if strings.HasPrefix(line, "## ") {
			// Look ahead to see if the section is empty
			sectionStart := i
			i++

			// Skip empty lines
			for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				i++
			}

			// Check if we hit another section or end of the file
			if i >= len(lines) || strings.HasPrefix(lines[i], "## ") || strings.HasPrefix(lines[i], "# ") {
				// Section is empty, skip it
				removed++
				continue
			}

			// Section has content, keep it
			result = append(result, lines[sectionStart:i]...)
			continue
		}

		result = append(result, line)
		i++
	}

	return strings.Join(result, "\n"), removed
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// preCompactAutoSave saves a session snapshot before compacting.
func preCompactAutoSave(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Ensure sessions directory exists
	sessionsDir := filepath.Join(contextDirName, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf("%s-pre-compact.md", now.Format("2006-01-02-150405"))
	filePath := filepath.Join(sessionsDir, filename)

	// Build minimal session content
	content := buildPreCompactSession(now)

	// Write the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	cmd.Printf("%s Auto-saved pre-compact snapshot to %s\n\n", green("✓"), filePath)
	return nil
}

// buildPreCompactSession creates a minimal session snapshot before compact.
func buildPreCompactSession(timestamp time.Time) string {
	var sb strings.Builder

	sb.WriteString("# Pre-Compact Snapshot\n\n")
	sb.WriteString(fmt.Sprintf("**Date**: %s\n", timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s\n", timestamp.Format("15:04:05")))
	sb.WriteString("**Type**: pre-compact\n\n")
	sb.WriteString("---\n\n")

	sb.WriteString("## Purpose\n\n")
	sb.WriteString("This snapshot was automatically created before running `ctx compact`.\n")
	sb.WriteString("It preserves the state of context files before any cleanup operations.\n\n")
	sb.WriteString("---\n\n")

	// Read and include current TASKS.md content
	tasksPath := filepath.Join(contextDirName, "TASKS.md")
	if tasksContent, err := os.ReadFile(tasksPath); err == nil {
		sb.WriteString("## Tasks (Before Compact)\n\n")
		sb.WriteString("```markdown\n")
		sb.WriteString(string(tasksContent))
		sb.WriteString("\n```\n\n")
	}

	return sb.String()
}
