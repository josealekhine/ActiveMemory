//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

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
)

// compactTasks moves completed tasks to the "Completed" section in TASKS.md.
//
// Scans TASKS.md for checked items ("- [x]") outside the Completed section,
// moves them into the Completed section, and optionally archives them to
// .context/archive/.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context containing the files
//   - archive: If true, write completed tasks to a dated archive file
//
// Returns:
//   - int: Number of tasks moved
//   - error: Non-nil if file write fails
func compactTasks(
	cmd *cobra.Command, ctx *context.Context, archive bool,
) (int, error) {
	var tasksFile *context.FileInfo
	for i := range ctx.Files {
		if ctx.Files[i].Name == config.FilenameTask {
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
				newContent := append(newLines[:insertIdx],
					append(tasksToInsert, newLines[insertIdx:]...,
					)...,
				)
				newLines = newContent
				break
			}
		}
	}

	// Archive old content if requested
	if archive && len(completedTasks) > 0 {
		archiveDir := filepath.Join(config.DirContext, "archive")
		if err := os.MkdirAll(archiveDir, 0755); err == nil {
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
			if err := os.WriteFile(
				archiveFile, []byte(archiveContent), 0644,
			); err == nil {
				cmd.Printf(
					"%s Archived %d tasks to %s\n", green("✓"),
					len(completedTasks), archiveFile,
				)
			}
		}
	}

	// Write back
	newContent := strings.Join(newLines, "\n")
	if newContent != content {
		if err := os.WriteFile(
			tasksFile.Path, []byte(newContent), 0644,
		); err != nil {
			return 0, err
		}
	}

	return changes, nil
}
