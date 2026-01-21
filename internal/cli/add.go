package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	addPriority string
	addSection  string
)

// fileTypeMap maps short names to actual file names
var fileTypeMap = map[string]string{
	"decision":   "DECISIONS.md",
	"decisions":  "DECISIONS.md",
	"task":       "TASKS.md",
	"tasks":      "TASKS.md",
	"learning":   "LEARNINGS.md",
	"learnings":  "LEARNINGS.md",
	"convention": "CONVENTIONS.md",
	"conventions": "CONVENTIONS.md",
}

// AddCmd returns the add command.
func AddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <type> <content>",
		Short: "Add a new item to a context file",
		Long: `Add a new decision, task, learning, or convention to the appropriate context file.

Types:
  decision    Add to DECISIONS.md
  task        Add to TASKS.md
  learning    Add to LEARNINGS.md
  convention  Add to CONVENTIONS.md

Examples:
  ctx add decision "Use PostgreSQL for primary database"
  ctx add task "Implement user authentication" --priority high
  ctx add learning "Vitest mocks must be hoisted"
  ctx add convention "All API routes must be versioned"`,
		Args: cobra.MinimumNArgs(2),
		RunE: runAdd,
	}

	cmd.Flags().StringVarP(&addPriority, "priority", "p", "", "Priority level for tasks (high, medium, low)")
	cmd.Flags().StringVarP(&addSection, "section", "s", "", "Target section within file")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	fileType := strings.ToLower(args[0])
	content := strings.Join(args[1:], " ")

	fileName, ok := fileTypeMap[fileType]
	if !ok {
		return fmt.Errorf("unknown type %q. Valid types: decision, task, learning, convention", fileType)
	}

	filePath := filepath.Join(contextDirName, fileName)

	// Check if file exists
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
	fmt.Printf("%s Added to %s\n", green("✓"), fileName)

	return nil
}

func formatDecision(content string) string {
	date := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`## [%s] %s

**Status**: Accepted

**Context**: [Add context here]

**Decision**: %s

**Rationale**: [Add rationale here]

**Consequences**: [Add consequences here]
`, date, content, content)
}

func formatTask(content string, priority string) string {
	var priorityTag string
	if priority != "" {
		priorityTag = fmt.Sprintf(" #priority:%s", priority)
	}
	return fmt.Sprintf("- [ ] %s%s\n", content, priorityTag)
}

func formatLearning(content string) string {
	date := time.Now().Format("2006-01-02")
	return fmt.Sprintf("- **[%s]** %s\n", date, content)
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
