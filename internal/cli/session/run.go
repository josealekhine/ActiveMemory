//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/validation"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runSessionLoad loads and displays a saved session file.
//
// Finds a session file matching the query (by filename, partial match, or index)
// and displays its contents. The query can be a full filename, a substring match,
// or a numeric index from the session list (1 = most recent).
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments where args[0] is the search query
//
// Returns:
//   - error: Non-nil if the sessions directory doesn't exist,
//     the file is not found, or read fails
func runSessionLoad(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Check if the sessions directory exists
	if _, err := os.Stat(sessionsDirPath()); os.IsNotExist(err) {
		return fmt.Errorf(
			"no sessions directory found. Run 'ctx session save' first",
		)
	}

	// Find the matching session file
	filePath, err := findSessionFile(query)
	if err != nil {
		return err
	}

	// Read and display
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read session file: %w", err)
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	cmd.Printf("%s Loading: %s\n\n", cyan("●"), filepath.Base(filePath))
	cmd.Println(string(content))

	return nil
}

// runSessionParse parses a JSONL transcript file and outputs formatted content.
//
// Converts a Claude Code JSONL transcript to readable Markdown. Can optionally
// extract potential decisions and learnings from the conversation using pattern
// matching. Output goes to stdout or a specified file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments where args[0] is the input JSONL file path
//   - output: Output file path (empty string for stdout)
//   - extract: If true, extract decisions/learnings instead of full transcript
//
// Returns:
//   - error: Non-nil if the file not found, parse fails, or write fails
func runSessionParse(
	cmd *cobra.Command, args []string, output string, extract bool,
) error {
	inputPath := args[0]

	// Check if the file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputPath)
	}

	green := color.New(color.FgGreen).SprintFunc()

	if extract {
		// Extract decisions and learnings
		decisions, learnings, err := extractInsights(inputPath)
		if err != nil {
			return fmt.Errorf("failed to extract insights: %w", err)
		}

		// Display extracted insights
		cmd.Println("# Extracted Insights")
		cmd.Println()
		cmd.Printf("**Source**: %s\n\n", filepath.Base(inputPath))

		cmd.Println("## Potential Decisions")
		cmd.Println()
		if len(decisions) == 0 {
			cmd.Println("No decisions detected.")
			cmd.Println()
		} else {
			for _, d := range decisions {
				cmd.Printf("- %s\n", d)
			}
			cmd.Println()
		}

		cmd.Println("## Potential Learnings")
		cmd.Println()
		if len(learnings) == 0 {
			cmd.Println("No learnings detected.")
			cmd.Println()
		} else {
			for _, l := range learnings {
				cmd.Printf("- %s\n", l)
			}
			cmd.Println()
		}

		cmd.Printf(
			"\n*Found %d potential decisions and %d potential learnings*\n",
			len(decisions), len(learnings),
		)
		return nil
	}

	// Parse the jsonl file
	content, err := parseJsonlTranscript(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse transcript: %w", err)
	}

	// Output
	if output != "" {
		if err := os.WriteFile(output, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		cmd.Printf("%s Parsed transcript saved to %s\n", green("✓"), output)
	} else {
		cmd.Println(content)
	}

	return nil
}

// runSessionSave saves the current context state to a session file.
//
// Creates a Markdown file in .context/sessions/ containing the current state
// of tasks, decisions, and learnings. The filename includes a timestamp and
// sanitized topic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments where args[0] is the optional topic
//   - sessionType: Type of session (feature, bugfix, refactor, session)
//
// Returns:
//   - error: Non-nil if directory creation, content building, or file write fails
func runSessionSave(
	cmd *cobra.Command, args []string, sessionType string,
) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Get topic from args or use default
	topic := "manual-save"
	if len(args) > 0 {
		topic = args[0]
	}

	// Sanitize the topic for filename
	topic = validation.SanitizeFilename(topic)

	// Ensure sessions directory exists
	if err := os.MkdirAll(sessionsDirPath(), 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02-150405"), topic)
	filePath := filepath.Join(sessionsDirPath(), filename)

	// Build session content
	content, err := buildSessionContent(topic, sessionType, now)
	if err != nil {
		return fmt.Errorf("failed to build session content: %w", err)
	}

	// Write the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	cmd.Printf("%s Session saved to %s\n", green("✓"), filePath)
	return nil
}

// runSessionList lists saved sessions with summaries.
//
// Reads all session files from .context/sessions/, parses their metadata,
// and displays them sorted by date (newest first). Output includes topic,
// date, type, summary, and filename for each session.
//
// Parameters:
//   - cmd: Cobra command for output
//   - limit: Maximum number of sessions to display (0 for unlimited)
//
// Returns:
//   - error: Non-nil if reading sessions directory fails
func runSessionList(cmd *cobra.Command, limit int) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if the `sessions` directory exists
	if _, err := os.Stat(sessionsDirPath()); os.IsNotExist(err) {
		cmd.Println("No sessions found. Use 'ctx session save' to create one.")
		return nil
	}

	// Read directory
	entries, err := os.ReadDir(sessionsDirPath())
	if err != nil {
		return fmt.Errorf("failed to read sessions directory: %w", err)
	}

	// Filter and collect session files
	var sessions []sessionInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Only show .md files (not .jsonl transcripts)
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		// Skip summary files that accompany jsonl files
		if strings.HasSuffix(name, "-summary.md") {
			continue
		}

		info, err := parseSessionFile(filepath.Join(sessionsDirPath(), name))
		if err != nil {
			// Skip files that can't be parsed
			continue
		}
		info.Filename = name
		sessions = append(sessions, info)
	}

	if len(sessions) == 0 {
		cmd.Println("No sessions found. Use 'ctx session save' to create one.")
		return nil
	}

	// Sort by date (newest first) - filenames are date-prefixed
	// so the reverse sort works
	for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	}

	// Limit output
	if limit > 0 && len(sessions) > limit {
		sessions = sessions[:limit]
	}

	// Display
	cmd.Printf("Sessions in %s:\n\n", sessionsDirPath())
	for _, s := range sessions {
		cmd.Printf("%s %s\n", cyan("●"), s.Topic)
		cmd.Printf("  %s %s | %s %s\n",
			gray("Date:"), s.Date,
			gray("Type:"), s.Type)
		if s.Summary != "" {
			cmd.Printf("  %s %s\n", gray("Summary:"), truncate(s.Summary, 60))
		}
		cmd.Printf("  %s %s\n", yellow("File:"), s.Filename)
		cmd.Println()
	}

	cmd.Printf("Total: %d session(s)\n", len(sessions))
	return nil
}
