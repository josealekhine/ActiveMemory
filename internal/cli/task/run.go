//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/validation"
)

// runTasksSnapshot executes the snapshot subcommand logic.
//
// Creates a point-in-time copy of TASKS.md in the archive directory.
// The snapshot includes a header with the name and timestamp.
//
// Parameters:
//   - cmd: Cobra command (unused, for interface compliance)
//   - args: Optional snapshot name as first argument
//
// Returns:
//   - error: Non-nil if TASKS.md doesn't exist or file operations fail
func runTasksSnapshot(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()
	tasksPath := tasksFilePath()
	archivePath := archiveDirPath()

	// Check if TASKS.md exists
	if _, err := os.Stat(tasksPath); os.IsNotExist(err) {
		return fmt.Errorf("no TASKS.md found")
	}

	// Read TASKS.md
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", err)
	}

	// Ensure the archive directory exists
	if err := os.MkdirAll(archivePath, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Generate snapshot filename
	now := time.Now()
	name := "snapshot"
	if len(args) > 0 {
		name = validation.SanitizeFilename(args[0])
	}
	snapshotFilename := fmt.Sprintf(
		"tasks-%s-%s.md", name, now.Format("2006-01-02-1504"),
	)
	snapshotPath := filepath.Join(archivePath, snapshotFilename)

	// Add snapshot header
	snapshotContent := fmt.Sprintf(
		"# TASKS.md Snapshot — %s\n\nCreated: %s\n\n---\n\n%s",
		name, now.Format(time.RFC3339), string(content),
	)

	// Write snapshot
	if err := os.WriteFile(
		snapshotPath, []byte(snapshotContent), 0644,
	); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	fmt.Printf("%s Snapshot saved to %s\n", green("✓"), snapshotPath)

	return nil
}

// runTaskArchive executes the archive subcommand logic.
//
// Moves completed tasks (marked with [x]) from TASKS.md to a timestamped
// archive file. Pending tasks ([ ]) remain in TASKS.md. If an archive file
// for the current date already exists, completed tasks are appended to it.
//
// Parameters:
//   - cmd: Cobra command (unused, for interface compliance)
//   - dryRun: If true, preview changes without modifying files
//
// Returns:
//   - error: Non-nil if TASKS.md doesn't exist or file operations fail
func runTaskArchive(cmd *cobra.Command, dryRun bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	tasksPath := tasksFilePath()
	archivePath := archiveDirPath()

	// Check if TASKS.md exists
	if _, err := os.Stat(tasksPath); os.IsNotExist(err) {
		return fmt.Errorf("no TASKS.md found")
	}

	// Read TASKS.md
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", err)
	}

	// Parse and separate completed versus pending tasks
	remaining, archived, stats := separateTasks(string(content))

	if stats.completed == 0 {
		fmt.Println("No completed tasks to archive.")
		return nil
	}

	if dryRun {
		fmt.Println(yellow("Dry run - no files modified"))
		fmt.Println()
		fmt.Printf(
			"Would archive %d completed tasks (keeping %d pending)\n",
			stats.completed, stats.pending,
		)
		fmt.Println()
		fmt.Println("Archived content preview:")
		fmt.Println("---")
		fmt.Println(archived)
		fmt.Println("---")
		return nil
	}

	// Ensure the archive directory exists
	if err := os.MkdirAll(archivePath, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Generate archive filename
	now := time.Now()
	archiveFilename := fmt.Sprintf("tasks-%s.md", now.Format("2006-01-02"))
	archiveFilePath := filepath.Join(archivePath, archiveFilename)

	// Check if the archive file already exists for today - append if so
	var archiveContent string
	if existingContent, err := os.ReadFile(archiveFilePath); err == nil {
		archiveContent = string(existingContent) + "\n" + archived
	} else {
		archiveContent = fmt.Sprintf(
			"# Task Archive — %s\n\nArchived from TASKS.md\n\n%s",
			now.Format("2006-01-02"),
			archived,
		)
	}

	// Write the archive file
	if err := os.WriteFile(
		archiveFilePath, []byte(archiveContent), 0644,
	); err != nil {
		return fmt.Errorf("failed to write archive: %w", err)
	}

	// Write updated TASKS.md
	if err := os.WriteFile(
		tasksPath, []byte(remaining), 0644,
	); err != nil {
		return fmt.Errorf("failed to update TASKS.md: %w", err)
	}

	fmt.Printf(
		"%s Archived %d completed tasks to %s\n",
		green("✓"),
		stats.completed,
		archiveFilePath,
	)
	fmt.Printf("  %d pending tasks remain in TASKS.md\n", stats.pending)

	return nil
}
