//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// runWatch executes the watch command logic.
//
// Sets up a reader from either a log file (--log) or stdin, then
// processes the stream for context update commands. Displays status
// messages and respects the --dry-run flag.
//
// Parameters:
//   - cmd: Cobra command for output
//   - _: Unused positional arguments
//
// Returns:
//   - error: Non-nil if context directory is missing, log file cannot
//     be opened, or stream processing fails
func runWatch(cmd *cobra.Command, _ []string) error {
	if !context.Exists("") {
		return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	cmd.Println(cyan("Watching for context updates..."))
	if watchDryRun {
		yellow := color.New(color.FgYellow).SprintFunc()
		cmd.Println(yellow("DRY RUN — No changes will be made"))
	}
	cmd.Println("Press Ctrl+C to stop")
	cmd.Println()

	var reader io.Reader
	if watchLog != "" {
		file, err := os.Open(watchLog)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("failed to close log file: %v\n", err)
			}
		}(file)
		reader = file
	} else {
		reader = os.Stdin
	}

	return processStream(cmd, reader)
}

// runAddSilent appends an entry to a context file without output.
//
// Used by the watch command to silently apply updates detected in
// the input stream. Formats the entry based on type and appends it
// to the appropriate context file.
//
// Parameters:
//   - args: Slice where args[0] is the entry type (task, decision, etc.)
//     and args[1:] is the content
//
// Returns:
//   - error: Non-nil if args has fewer than 2 elements, type is unknown,
//     or file operations fail
func runAddSilent(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("insufficient arguments")
	}

	fileType := strings.ToLower(args[0])
	content := strings.Join(args[1:], " ")

	fileName, ok := config.FileType[fileType]
	if !ok {
		return fmt.Errorf("unknown type %q", fileType)
	}

	filePath := filepath.Join(config.DirContext, fileName)

	existing, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var entry string
	switch fileType {
	case config.UpdateTypeDecision, config.UpdateTypeDecisions:
		entry = add.FormatDecision(content)
	case config.UpdateTypeTask, config.UpdateTypeTasks:
		entry = add.FormatTask(content, "")
	case config.UpdateTypeLearning, config.UpdateTypeLearnings:
		entry = add.FormatLearning(content)
	case config.UpdateTypeConvention, config.UpdateTypeConventions:
		entry = add.FormatConvention(content)
	}

	newContent := add.AppendEntry(existing, entry, fileType, "")
	return os.WriteFile(filePath, newContent, 0644)
}

// runCompleteSilent marks a task as complete without output.
//
// Used by the watch command to silently complete tasks detected in
// the input stream. Searches for an unchecked task matching the query
// and marks it as done by changing [ ] to [x].
//
// Parameters:
//   - args: Slice where args[0] is the search query to match against
//     task descriptions (case-insensitive substring match)
//
// Returns:
//   - error: Non-nil if args is empty, no matching task is found,
//     or file operations fail
func runCompleteSilent(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no task specified")
	}

	query := args[0]
	filePath := filepath.Join(config.DirContext, config.FilenameTask)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	taskPattern := regexp.MustCompile(`^(\s*)-\s*\[\s*]\s*(.+)$`)

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

	lines[matchedLine] = taskPattern.ReplaceAllString(
		lines[matchedLine], "$1- [x] $2",
	)
	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}
