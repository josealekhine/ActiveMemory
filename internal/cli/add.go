//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	addPriority  string
	addSection   string
	addFromFile  string
)

// fileTypeMap maps short names to actual file names
var fileTypeMap = map[string]string{
	"decision":    "DECISIONS.md",
	"decisions":   "DECISIONS.md",
	"task":        "TASKS.md",
	"tasks":       "TASKS.md",
	"learning":    "LEARNINGS.md",
	"learnings":   "LEARNINGS.md",
	"convention":  "CONVENTIONS.md",
	"conventions": "CONVENTIONS.md",
}

// AddCmd returns the add command.
func AddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <type> [content]",
		Short: "Add a new item to a context file",
		Long: `Add a new decision, task, learning, or convention to the appropriate context file.

Types:
  decision    Add to DECISIONS.md
  task        Add to TASKS.md
  learning    Add to LEARNINGS.md
  convention  Add to CONVENTIONS.md

Content can be provided as:
  - Command argument: ctx add learning "text here"
  - File: ctx add learning --file /path/to/content.md
  - Stdin: echo "text" | ctx add learning

Examples:
  ctx add decision "Use PostgreSQL for primary database"
  ctx add task "Implement user authentication" --priority high
  ctx add learning "Vitest mocks must be hoisted"
  ctx add learning --file learning-template.md
  cat notes.md | ctx add decision`,
		Args: cobra.MinimumNArgs(1),
		RunE: runAdd,
	}

	cmd.Flags().StringVarP(&addPriority, "priority", "p", "", "Priority level for tasks (high, medium, low)")
	cmd.Flags().StringVarP(&addSection, "section", "s", "", "Target section within file")
	cmd.Flags().StringVarP(&addFromFile, "file", "f", "", "Read content from file instead of argument")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	fileType := strings.ToLower(args[0])

	// Determine content source: args, --file, or stdin
	var content string

	if addFromFile != "" {
		// Read from file
		fileContent, err := os.ReadFile(addFromFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", addFromFile, err)
		}
		content = strings.TrimSpace(string(fileContent))
	} else if len(args) > 1 {
		// Content from arguments
		content = strings.Join(args[1:], " ")
	} else {
		// Try reading from stdin (check if it's a pipe)
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// stdin is a pipe, read from it
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = strings.TrimSpace(strings.Join(lines, "\n"))
		}
	}

	if content == "" {
		return fmt.Errorf("no content provided. Use argument, --file, or pipe from stdin")
	}

	fileName, ok := fileTypeMap[fileType]
	if !ok {
		return fmt.Errorf("unknown type %q. Valid types: decision, task, learning, convention", fileType)
	}

	filePath := filepath.Join(contextDirName, fileName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("context file %s not found. Run 'ctx init' first", filePath)
	}

	// Read existing content
	existing, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	// Format the new entry based on type
	var entry string
	switch fileType {
	case "decision", "decisions":
		entry = formatDecision(content)
	case "task", "tasks":
		entry = formatTask(content, addPriority)
	case "learning", "learnings":
		entry = formatLearning(content)
	case "convention", "conventions":
		entry = formatConvention(content)
	}

	// Append to file
	newContent := appendEntry(existing, entry, fileType, addSection)

	if err := os.WriteFile(filePath, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Printf("%s Added to %s\n", green("✓"), fileName)

	return nil
}

func formatDecision(content string) string {
	// Use YYYY-MM-DD-HHMM for precise timestamp correlation with sessions
	timestamp := time.Now().Format("2006-01-02-1504")
	return fmt.Sprintf(`## [%s] %s

**Status**: Accepted

**Context**: [Add context here]

**Decision**: %s

**Rationale**: [Add rationale here]

**Consequences**: [Add consequences here]
`, timestamp, content, content)
}

func formatTask(content string, priority string) string {
	// Use YYYY-MM-DD-HHMM timestamp for session correlation
	timestamp := time.Now().Format("2006-01-02-1504")
	var priorityTag string
	if priority != "" {
		priorityTag = fmt.Sprintf(" #priority:%s", priority)
	}
	return fmt.Sprintf("- [ ] %s%s #added:%s\n", content, priorityTag, timestamp)
}

func formatLearning(content string) string {
	// Use YYYY-MM-DD-HHMM for precise timestamp correlation with sessions
	timestamp := time.Now().Format("2006-01-02-1504")
	return fmt.Sprintf("- **[%s]** %s\n", timestamp, content)
}

func formatConvention(content string) string {
	return fmt.Sprintf("- %s\n", content)
}

func appendEntry(existing []byte, entry string, fileType string, section string) []byte {
	existingStr := string(existing)

	// For tasks, find the appropriate section
	if fileType == "task" || fileType == "tasks" {
		targetSection := section
		if targetSection == "" {
			targetSection = "## Next Up"
		} else if !strings.HasPrefix(targetSection, "##") {
			targetSection = "## " + targetSection
		}

		// Find the section and insert after it
		idx := strings.Index(existingStr, targetSection)
		if idx != -1 {
			// Find the end of the section header line
			lineEnd := strings.Index(existingStr[idx:], "\n")
			if lineEnd != -1 {
				insertPoint := idx + lineEnd + 1
				return []byte(existingStr[:insertPoint] + "\n" + entry + existingStr[insertPoint:])
			}
		}
	}

	// For decisions, insert before the closing comment if present, otherwise append
	if fileType == "decision" || fileType == "decisions" {
		// Just append at the end
		if !strings.HasSuffix(existingStr, "\n") {
			existingStr += "\n"
		}
		return []byte(existingStr + "\n" + entry)
	}

	// Default: append at the end
	if !strings.HasSuffix(existingStr, "\n") {
		existingStr += "\n"
	}
	return []byte(existingStr + "\n" + entry)
}
