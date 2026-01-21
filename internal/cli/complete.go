package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// CompleteCmd returns the complete command.
func CompleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete <task-id-or-text>",
		Short: "Mark a task as completed",
		Long: `Mark a task as completed in TASKS.md.

You can specify a task by:
  - Task number (e.g., "ctx complete 3")
  - Partial text match (e.g., "ctx complete auth")
  - Full task text (e.g., "ctx complete 'Implement user authentication'")

The task will be marked with [x] and optionally moved to the Completed section.`,
		Args: cobra.ExactArgs(1),
		RunE: runComplete,
	}

	return cmd
}

func runComplete(cmd *cobra.Command, args []string) error {
	query := args[0]

	filePath := filepath.Join(contextDirName, "TASKS.md")

	// Check if file exists
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
	taskPattern := regexp.MustCompile(`^(\s*)-\s*\[\s*\]\s*(.+)$`)

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
			if !isNumber && strings.Contains(strings.ToLower(taskText), strings.ToLower(query)) {
				if matchedLine != -1 {
					// Multiple matches - be more specific
					return fmt.Errorf("multiple tasks match %q. Be more specific or use task number", query)
				}
				matchedLine = i
				matchedTask = taskText
			}
		}
	}

	if matchedLine == -1 {
		if isNumber {
			return fmt.Errorf("task #%d not found. Use 'ctx status' to see tasks", taskNumber)
		}
		return fmt.Errorf("no task matching %q found. Use 'ctx status' to see tasks", query)
	}

	// Mark the task as complete
	lines[matchedLine] = taskPattern.ReplaceAllString(lines[matchedLine], "$1- [x] $2")

	// Write back
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write TASKS.md: %w", err)
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s Completed: %s\n", green("✓"), matchedTask)

	return nil
}
