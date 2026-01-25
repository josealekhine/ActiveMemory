//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// runComplete executes the complete command logic.
//
// Finds a task in TASKS.md by number or text match and marks it complete
// by changing "- [ ]" to "- [x]".
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments; args[0] is the task number or search text
//
// Returns:
//   - error: Non-nil if the task is not found, multiple matches, or file
//     operations fail
func runComplete(cmd *cobra.Command, args []string) error {
	query := args[0]

	filePath := filepath.Join(config.DirContext, config.FilenameTask)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("TASKS.md not found. Run 'ctx init' first")
	}

	// Read existing content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", err)
	}

	// Parse tasks and find matching one
	lines := strings.Split(string(content), "\n")
	taskPattern := regexp.MustCompile(`^(\s*)-\s*\[\s*]\s*(.+)$`)

	var taskNumber int
	isNumber := false
	if num, err := strconv.Atoi(query); err == nil {
		taskNumber = num
		isNumber = true
	}

	currentTaskNum := 0
	matchedLine := -1
	matchedTask := ""

	for i, line := range lines {
		matches := taskPattern.FindStringSubmatch(line)
		if matches != nil {
			currentTaskNum++
			taskText := matches[2]

			// Match by number
			if isNumber && currentTaskNum == taskNumber {
				matchedLine = i
				matchedTask = taskText
				break
			}

			// Match by text (case-insensitive partial match)
			if !isNumber && strings.Contains(
				strings.ToLower(taskText), strings.ToLower(query),
			) {
				if matchedLine != -1 {
					// Multiple matches - be more specific
					return fmt.Errorf(
						"multiple tasks match %q. Be more specific or use task number",
						query,
					)
				}
				matchedLine = i
				matchedTask = taskText
			}
		}
	}

	if matchedLine == -1 {
		if isNumber {
			return fmt.Errorf(
				"task #%d not found. Use 'ctx status' to see tasks", taskNumber,
			)
		}
		return fmt.Errorf(
			"no task matching %q found. Use 'ctx status' to see tasks", query,
		)
	}

	// Mark the task as complete
	lines[matchedLine] = taskPattern.ReplaceAllString(
		lines[matchedLine], "$1- [x] $2",
	)

	// Write back
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write TASKS.md: %w", err)
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Printf("%s Completed: %s\n", green("✓"), matchedTask)

	return nil
}
