package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/josealekhine/ActiveMemory/internal/context"
	"github.com/spf13/cobra"
)

var (
	watchLog    string
	watchDryRun bool
)

// WatchCmd returns the watch command.
func WatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for context-update commands in AI output",
		Long: `Watch stdin or a log file for <context-update> commands and apply them.

This command parses AI output looking for structured update commands:

  <context-update type="task">Implement user auth</context-update>
  <context-update type="decision">Use PostgreSQL</context-update>
  <context-update type="learning">Mock functions must be hoisted</context-update>
  <context-update type="complete">user auth</context-update>

Use --log to watch a specific file instead of stdin.
Use --dry-run to see what would be updated without making changes.

Press Ctrl+C to stop watching.`,
		RunE: runWatch,
	}

	cmd.Flags().StringVar(&watchLog, "log", "", "Log file to watch (default: stdin)")
	cmd.Flags().BoolVar(&watchDryRun, "dry-run", false, "Show updates without applying")

	return cmd
}

// ContextUpdate represents a parsed context update command.
type ContextUpdate struct {
	Type    string
	Content string
}

func runWatch(cmd *cobra.Command, args []string) error {
	// Check if context exists
	if !context.Exists("") {
		return fmt.Errorf("no .context/ directory found. Run 'amem init' first")
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Println(cyan("Watching for context updates..."))
	if watchDryRun {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Println(yellow("DRY RUN — No changes will be made"))
	}
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	var reader io.Reader
	if watchLog != "" {
		file, err := os.Open(watchLog)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	return processStream(reader)
}

func processStream(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	// Use a larger buffer for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	// Pattern to match context-update tags
	updatePattern := regexp.MustCompile(`<context-update\s+type="([^"]+)"[^>]*>([^<]+)</context-update>`)

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	for scanner.Scan() {
		line := scanner.Text()

		// Check for context-update commands
		matches := updatePattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				update := ContextUpdate{
					Type:    strings.ToLower(match[1]),
					Content: strings.TrimSpace(match[2]),
				}

				if watchDryRun {
					fmt.Printf("%s Would apply: [%s] %s\n", yellow("○"), update.Type, update.Content)
				} else {
					err := applyUpdate(update)
					if err != nil {
						fmt.Printf("%s Failed to apply [%s]: %v\n", color.RedString("✗"), update.Type, err)
					} else {
						fmt.Printf("%s Applied: [%s] %s\n", green("✓"), update.Type, update.Content)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

func applyUpdate(update ContextUpdate) error {
	switch update.Type {
	case "task":
		return applyTaskUpdate(update.Content)
	case "decision":
		return applyDecisionUpdate(update.Content)
	case "learning":
		return applyLearningUpdate(update.Content)
	case "convention":
		return applyConventionUpdate(update.Content)
	case "complete":
		return applyCompleteUpdate(update.Content)
	default:
		return fmt.Errorf("unknown update type: %s", update.Type)
	}
}

func applyTaskUpdate(content string) error {
	// Reuse the add command logic
	args := []string{"task", content}
	return runAdd(nil, args)
}

func applyDecisionUpdate(content string) error {
	args := []string{"decision", content}
	// Suppress output from add command during watch
	return runAddSilent(args)
}

func applyLearningUpdate(content string) error {
	args := []string{"learning", content}
	return runAddSilent(args)
}

func applyConventionUpdate(content string) error {
	args := []string{"convention", content}
	return runAddSilent(args)
}

func applyCompleteUpdate(content string) error {
	args := []string{content}
	return runCompleteSilent(args)
}

// runAddSilent runs the add command without output
func runAddSilent(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("insufficient arguments")
	}

	fileType := strings.ToLower(args[0])
	content := strings.Join(args[1:], " ")

	fileName, ok := fileTypeMap[fileType]
	if !ok {
		return fmt.Errorf("unknown type %q", fileType)
	}

	filePath := contextDirName + "/" + fileName

	existing, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var entry string
	switch fileType {
	case "decision", "decisions":
		entry = formatDecision(content)
	case "task", "tasks":
		entry = formatTask(content, "")
	case "learning", "learnings":
		entry = formatLearning(content)
	case "convention", "conventions":
		entry = formatConvention(content)
	}

	newContent := appendEntry(existing, entry, fileType, "")
	return os.WriteFile(filePath, newContent, 0644)
}

// runCompleteSilent runs the complete command without output
func runCompleteSilent(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no task specified")
	}

	query := args[0]
	filePath := contextDirName + "/TASKS.md"

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	taskPattern := regexp.MustCompile(`^(\s*)-\s*\[\s*\]\s*(.+)$`)

	matchedLine := -1
	for i, line := range lines {
		matches := taskPattern.FindStringSubmatch(line)
		if matches != nil {
			taskText := matches[2]
			if strings.Contains(strings.ToLower(taskText), strings.ToLower(query)) {
				matchedLine = i
				break
			}
		}
	}

	if matchedLine == -1 {
		return fmt.Errorf("no task matching %q found", query)
	}

	lines[matchedLine] = taskPattern.ReplaceAllString(lines[matchedLine], "$1- [x] $2")
	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}
