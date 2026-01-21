package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/josealekhine/ActiveMemory/internal/context"
	"github.com/spf13/cobra"
)

var (
	compactArchive bool
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

Use --archive to create .context/archive/ for old content.`,
		RunE: runCompact,
	}

	cmd.Flags().BoolVar(&compactArchive, "archive", false, "Create .context/archive/ for old content")

	return cmd
}

func runCompact(cmd *cobra.Command, args []string) error {
	ctx, err := context.Load("")
	if err != nil {
		if _, ok := err.(*context.NotFoundError); ok {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(cyan("Compact Analysis"))
	fmt.Println(cyan("================"))
	fmt.Println()

	changes := 0

	// Process TASKS.md
	tasksChanges, err := compactTasks(ctx, compactArchive)
	if err != nil {
		fmt.Printf("%s Error processing TASKS.md: %v\n", yellow("⚠"), err)
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
				fmt.Printf("%s Removed %d empty sections from %s\n", green("✓"), count, f.Name)
				changes += count
			}
		}
	}

	if changes == 0 {
		fmt.Printf("%s Nothing to compact — context is already clean\n", green("✓"))
	} else {
		fmt.Printf("\n%s Compacted %d items\n", green("✓"), changes)
	}

	return nil
}

func compactTasks(ctx *context.Context, archive bool) (int, error) {
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

	completedPattern := regexp.MustCompile(`^-\s*\[x\]\s*(.+)$`)

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

		// If completed task outside Completed section, collect it
		if !inCompletedSection && completedPattern.MatchString(line) {
			matches := completedPattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				completedTasks = append(completedTasks, matches[1])
				fmt.Printf("%s Moving completed task: %s\n", green("✓"), truncateString(matches[1], 50))
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
				for insertIdx < len(newLines) && newLines[insertIdx] != "" && !strings.HasPrefix(newLines[insertIdx], "## ") {
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
				fmt.Printf("%s Archived %d tasks to %s\n", green("✓"), len(completedTasks), archiveFile)
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
			// Look ahead to see if section is empty
			sectionStart := i
			i++

			// Skip empty lines
			for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				i++
			}

			// Check if we hit another section or end of file
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
