//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// runAdd executes the add command logic.
//
// It reads content from the specified source (argument, file, or stdin),
// validates the entry type, formats the entry, and appends it to the
// appropriate context file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments; args[0] is the entry type, args[1:] is content
//   - flags: All flag values from the command
//
// Returns:
//   - error: Non-nil if content is missing, type is invalid, required flags
//     are missing, or file operations fail
func runAdd(cmd *cobra.Command, args []string, flags addFlags) error {
	fType := strings.ToLower(args[0])

	// Validate required flags for decisions
	if fType == config.UpdateTypeDecision || fType == config.UpdateTypeDecisions {
		var missing []string
		if flags.context == "" {
			missing = append(missing, "--context")
		}
		if flags.rationale == "" {
			missing = append(missing, "--rationale")
		}
		if flags.consequences == "" {
			missing = append(missing, "--consequences")
		}
		if len(missing) > 0 {
			return fmt.Errorf(`decisions require complete ADR format

Missing required flags: %s

Usage:
  ctx add decision "Decision title" \
    --context "What prompted this decision" \
    --rationale "Why this choice over alternatives" \
    --consequences "What changes as a result"

Example:
  ctx add decision "Use PostgreSQL for primary database" \
    --context "Need a reliable database for production workloads" \
    --rationale "PostgreSQL offers ACID compliance, JSON support, and team familiarity" \
    --consequences "Team needs PostgreSQL training; must set up replication"`,
				strings.Join(missing, ", "))
		}
	}

	// Validate required flags for learnings
	if fType == config.UpdateTypeLearning || fType == config.UpdateTypeLearnings {
		var missing []string
		if flags.context == "" {
			missing = append(missing, "--context")
		}
		if flags.lesson == "" {
			missing = append(missing, "--lesson")
		}
		if flags.application == "" {
			missing = append(missing, "--application")
		}
		if len(missing) > 0 {
			return fmt.Errorf(`learnings require complete format

Missing required flags: %s

Usage:
  ctx add learning "Learning title" \
    --context "What prompted this learning" \
    --lesson "The key insight" \
    --application "How to apply this going forward"

Example:
  ctx add learning "Go embed requires files in same package" \
    --context "Tried to embed files from parent directory, got compile error" \
    --lesson "go:embed only works with files in same or child directories" \
    --application "Keep embedded files in internal/templates/, not project root"`,
				strings.Join(missing, ", "))
		}
	}

	// Determine the content source: args, --file, or stdin
	var content string

	if flags.fromFile != "" {
		// Read from the file
		fileContent, err := os.ReadFile(flags.fromFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", flags.fromFile, err)
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
		examples := getExamplesForType(fType)
		return fmt.Errorf(`no content provided

Usage:
  ctx add %s "your content here"
  ctx add %s --file /path/to/content.md
  echo "content" | ctx add %s

Examples:
%s`, fType, fType, fType, examples)
	}

	fName, ok := config.FileType[fType]
	if !ok {
		return fmt.Errorf(
			"unknown type %q. Valid types: decision, task, learning, convention",
			fType,
		)
	}

	filePath := filepath.Join(config.DirContext, fName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf(
			"context file %s not found. Run 'ctx init' first", filePath,
		)
	}

	// Read existing content
	existing, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	// Format the new entry based on type
	var entry string
	switch fType {
	case config.UpdateTypeDecision, config.UpdateTypeDecisions:
		entry = FormatDecision(content, flags.context, flags.rationale, flags.consequences)
	case config.UpdateTypeTask, config.UpdateTypeTasks:
		entry = FormatTask(content, flags.priority)
	case config.UpdateTypeLearning, config.UpdateTypeLearnings:
		entry = FormatLearning(content, flags.context, flags.lesson, flags.application)
	case config.UpdateTypeConvention, config.UpdateTypeConventions:
		entry = FormatConvention(content)
	}

	// Append to file
	newContent := AppendEntry(existing, entry, fType, flags.section)

	if err := os.WriteFile(filePath, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Printf("%s Added to %s\n", green("✓"), fName)

	return nil
}
