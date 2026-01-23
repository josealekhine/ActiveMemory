//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	sessionsDirName = ".context/sessions"
)

// SessionCmd returns the session command with subcommands.
func SessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage session snapshots",
		Long: `Manage session snapshots in .context/sessions/.

Sessions capture the state of your context at a point in time,
including current tasks, recent decisions, and learnings.

Subcommands:
  save    Save current context state to a session file
  list    List saved sessions with summaries
  load    Load and display a previous session
  parse   Convert .jsonl transcript to readable markdown`,
	}

	cmd.AddCommand(sessionSaveCmd())
	cmd.AddCommand(sessionListCmd())
	cmd.AddCommand(sessionLoadCmd())
	cmd.AddCommand(sessionParseCmd())

	return cmd
}

var sessionType string

// sessionSaveCmd returns the session save subcommand.
func sessionSaveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save [topic]",
		Short: "Save current context state to a session file",
		Long: `Save a snapshot of the current context state to .context/sessions/.

The session file includes:
  - Summary of what was done
  - Current tasks from TASKS.md
  - Recent decisions from DECISIONS.md
  - Recent learnings from LEARNINGS.md

Examples:
  ctx session save "implemented auth"
  ctx session save "refactored API" --type feature
  ctx session save  # prompts for topic interactively`,
		Args: cobra.MaximumNArgs(1),
		RunE: runSessionSave,
	}

	cmd.Flags().StringVarP(&sessionType, "type", "t", "session", "Session type (feature, bugfix, refactor, session)")

	return cmd
}

func runSessionSave(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Get topic from args or use default
	topic := "manual-save"
	if len(args) > 0 {
		topic = args[0]
	}

	// Sanitize topic for filename
	topic = sanitizeFilename(topic)

	// Ensure sessions directory exists
	if err := os.MkdirAll(sessionsDirName, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02-150405"), topic)
	filePath := filepath.Join(sessionsDirName, filename)

	// Build session content
	content, err := buildSessionContent(topic, sessionType, now)
	if err != nil {
		return fmt.Errorf("failed to build session content: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	cmd.Printf("%s Session saved to %s\n", green("✓"), filePath)
	return nil
}

var listLimit int

// sessionListCmd returns the session list subcommand.
func sessionListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List saved sessions with summaries",
		Long: `List all saved sessions in .context/sessions/.

Shows session date, topic, type, and a brief summary for each session.
Sessions are sorted by date (newest first).

Examples:
  ctx session list
  ctx session list --limit 5`,
		RunE: runSessionList,
	}

	cmd.Flags().IntVarP(&listLimit, "limit", "n", 10, "Maximum number of sessions to display")

	return cmd
}

func runSessionList(cmd *cobra.Command, args []string) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if sessions directory exists
	if _, err := os.Stat(sessionsDirName); os.IsNotExist(err) {
		cmd.Println("No sessions found. Use 'ctx session save' to create one.")
		return nil
	}

	// Read directory
	entries, err := os.ReadDir(sessionsDirName)
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

		info, err := parseSessionFile(filepath.Join(sessionsDirName, name))
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

	// Sort by date (newest first) - filenames are date-prefixed so reverse sort works
	for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	}

	// Limit output
	if listLimit > 0 && len(sessions) > listLimit {
		sessions = sessions[:listLimit]
	}

	// Display
	cmd.Printf("Sessions in %s:\n\n", sessionsDirName)
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

// sessionInfo holds parsed information about a session file.
type sessionInfo struct {
	Filename string
	Topic    string
	Date     string
	Type     string
	Summary  string
}

// parseSessionFile extracts metadata from a session file.
func parseSessionFile(path string) (sessionInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return sessionInfo{}, err
	}

	contentStr := string(content)
	info := sessionInfo{}

	// Extract topic from first line (# Session: topic)
	if strings.HasPrefix(contentStr, "# Session:") {
		lineEnd := strings.Index(contentStr, "\n")
		if lineEnd != -1 {
			info.Topic = strings.TrimSpace(contentStr[11:lineEnd])
		}
	} else if strings.HasPrefix(contentStr, "# ") {
		// Alternative format: # Topic
		lineEnd := strings.Index(contentStr, "\n")
		if lineEnd != -1 {
			info.Topic = strings.TrimSpace(contentStr[2:lineEnd])
		}
	}

	// Extract date
	if idx := strings.Index(contentStr, "**Date**:"); idx != -1 {
		lineEnd := strings.Index(contentStr[idx:], "\n")
		if lineEnd != -1 {
			info.Date = strings.TrimSpace(contentStr[idx+9 : idx+lineEnd])
		}
	}

	// Extract type
	if idx := strings.Index(contentStr, "**Type**:"); idx != -1 {
		lineEnd := strings.Index(contentStr[idx:], "\n")
		if lineEnd != -1 {
			info.Type = strings.TrimSpace(contentStr[idx+9 : idx+lineEnd])
		}
	}

	// Extract summary (first non-empty line after ## Summary)
	if idx := strings.Index(contentStr, "## Summary"); idx != -1 {
		afterSummary := contentStr[idx+10:]
		lines := strings.Split(afterSummary, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "[") {
				info.Summary = line
				break
			}
		}
	}

	return info, nil
}

// truncate shortens a string to maxLen characters, adding "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// sessionLoadCmd returns the session load subcommand.
func sessionLoadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load <file>",
		Short: "Load and display a previous session",
		Long: `Load and display the contents of a saved session.

The file argument can be:
  - A full filename (e.g., 2025-01-21-004900-ctx-rename.md)
  - A partial match (e.g., "ctx-rename" or "2025-01-21")
  - A number from 'ctx session list' output (1 = most recent)

Examples:
  ctx session load 2025-01-21-004900-ctx-rename.md
  ctx session load ctx-rename
  ctx session load 1`,
		Args: cobra.ExactArgs(1),
		RunE: runSessionLoad,
	}

	return cmd
}

func runSessionLoad(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Check if sessions directory exists
	if _, err := os.Stat(sessionsDirName); os.IsNotExist(err) {
		return fmt.Errorf("no sessions directory found. Run 'ctx session save' first")
	}

	// Find matching session file
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

// findSessionFile finds a session file matching the query.
// Query can be a full filename, partial match, or numeric index.
func findSessionFile(query string) (string, error) {
	// Read directory
	entries, err := os.ReadDir(sessionsDirName)
	if err != nil {
		return "", fmt.Errorf("failed to read sessions directory: %w", err)
	}

	// Collect .md files (excluding -summary.md)
	var sessions []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		if strings.HasSuffix(name, "-summary.md") {
			continue
		}
		sessions = append(sessions, name)
	}

	if len(sessions) == 0 {
		return "", fmt.Errorf("no sessions found")
	}

	// Reverse sort (newest first) for numeric indexing
	for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	}

	// Check if query is a number (index)
	if idx, err := parseIndex(query); err == nil {
		if idx < 1 || idx > len(sessions) {
			return "", fmt.Errorf("index %d out of range (1-%d)", idx, len(sessions))
		}
		return filepath.Join(sessionsDirName, sessions[idx-1]), nil
	}

	// Check for exact match
	for _, name := range sessions {
		if name == query {
			return filepath.Join(sessionsDirName, name), nil
		}
	}

	// Check for partial match
	query = strings.ToLower(query)
	var matches []string
	for _, name := range sessions {
		if strings.Contains(strings.ToLower(name), query) {
			matches = append(matches, name)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no session found matching %q", query)
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("multiple sessions match %q: %v", query, matches)
	}

	return filepath.Join(sessionsDirName, matches[0]), nil
}

// parseIndex attempts to parse a string as a positive integer index.
func parseIndex(s string) (int, error) {
	var idx int
	_, err := fmt.Sscanf(s, "%d", &idx)
	if err != nil {
		return 0, err
	}
	if idx < 1 {
		return 0, fmt.Errorf("index must be positive")
	}
	return idx, nil
}

var (
	parseOutput  string
	parseExtract bool
)

// sessionParseCmd returns the session parse subcommand.
func sessionParseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse <file.jsonl>",
		Short: "Convert .jsonl transcript to readable markdown",
		Long: `Parse a Claude Code .jsonl transcript file and convert it to readable markdown.

The .jsonl files are auto-saved by the SessionEnd hook and contain the full
conversation transcript including tool calls and results.

Examples:
  ctx session parse .context/sessions/2026-01-21-072504-session.jsonl
  ctx session parse .context/sessions/2026-01-21-072504-session.jsonl -o conversation.md`,
		Args: cobra.ExactArgs(1),
		RunE: runSessionParse,
	}

	cmd.Flags().StringVarP(&parseOutput, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().BoolVar(&parseExtract, "extract", false, "Extract potential decisions and learnings from transcript")

	return cmd
}

func runSessionParse(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Check file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputPath)
	}

	green := color.New(color.FgGreen).SprintFunc()

	if parseExtract {
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

		cmd.Printf("\n*Found %d potential decisions and %d potential learnings*\n", len(decisions), len(learnings))
		return nil
	}

	// Parse the jsonl file
	content, err := parseJsonlTranscript(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse transcript: %w", err)
	}

	// Output
	if parseOutput != "" {
		if err := os.WriteFile(parseOutput, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		cmd.Printf("%s Parsed transcript saved to %s\n", green("✓"), parseOutput)
	} else {
		cmd.Println(content)
	}

	return nil
}

// transcriptEntry represents a single entry in the jsonl transcript.
type transcriptEntry struct {
	Type      string        `json:"type"`
	Message   transcriptMsg `json:"message"`
	Timestamp string        `json:"timestamp"`
	UUID      string        `json:"uuid"`
}

// transcriptMsg represents the message content.
type transcriptMsg struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or []interface{}
}

// contentBlock represents a block in assistant content array.
type contentBlock struct {
	Type     string      `json:"type"`
	Text     string      `json:"text,omitempty"`
	Thinking string      `json:"thinking,omitempty"`
	Name     string      `json:"name,omitempty"`  // tool name
	Input    interface{} `json:"input,omitempty"` // tool input
}

// parseJsonlTranscript parses a .jsonl file and returns formatted markdown.
func parseJsonlTranscript(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	sb.WriteString("# Conversation Transcript\n\n")
	sb.WriteString(fmt.Sprintf("**Source**: %s\n\n", filepath.Base(path)))
	sb.WriteString("---\n\n")

	scanner := bufio.NewScanner(file)
	// Increase buffer size for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024) // 10MB max line size

	messageCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry transcriptEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Skip unparseable lines
			continue
		}

		// Skip non-message entries
		if entry.Type != "user" && entry.Type != "assistant" {
			continue
		}

		messageCount++
		formatted := formatTranscriptEntry(entry)
		if formatted != "" {
			sb.WriteString(formatted)
			sb.WriteString("\n---\n\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	sb.WriteString(fmt.Sprintf("*Total messages: %d*\n", messageCount))

	return sb.String(), nil
}

// formatTranscriptEntry formats a single transcript entry as markdown.
func formatTranscriptEntry(entry transcriptEntry) string {
	var sb strings.Builder

	// Header with role and timestamp
	role := strings.ToUpper(entry.Message.Role[:1]) + entry.Message.Role[1:]
	sb.WriteString(fmt.Sprintf("## %s\n\n", role))

	if entry.Timestamp != "" {
		// Parse and format timestamp
		if t, err := time.Parse(time.RFC3339, entry.Timestamp); err == nil {
			sb.WriteString(fmt.Sprintf("*%s*\n\n", t.Format("2006-01-02 15:04:05")))
		}
	}

	// Handle content
	switch content := entry.Message.Content.(type) {
	case string:
		sb.WriteString(content)
		sb.WriteString("\n")
	case []interface{}:
		for _, block := range content {
			blockMap, ok := block.(map[string]interface{})
			if !ok {
				continue
			}

			blockType, _ := blockMap["type"].(string)
			switch blockType {
			case "text":
				if text, ok := blockMap["text"].(string); ok {
					sb.WriteString(text)
					sb.WriteString("\n")
				}
			case "thinking":
				if thinking, ok := blockMap["thinking"].(string); ok {
					sb.WriteString("<details>\n<summary>💭 Thinking</summary>\n\n")
					sb.WriteString(thinking)
					sb.WriteString("\n</details>\n\n")
				}
			case "tool_use":
				name, _ := blockMap["name"].(string)
				sb.WriteString(fmt.Sprintf("**🔧 Tool: %s**\n", name))
				if input, ok := blockMap["input"].(map[string]interface{}); ok {
					// Show key parameters
					for k, v := range input {
						vStr := fmt.Sprintf("%v", v)
						if len(vStr) > 100 {
							vStr = vStr[:100] + "..."
						}
						sb.WriteString(fmt.Sprintf("- %s: `%s`\n", k, vStr))
					}
				}
				sb.WriteString("\n")
			case "tool_result":
				sb.WriteString("**📋 Tool Result**\n")
				if result, ok := blockMap["content"].(string); ok {
					if len(result) > 500 {
						result = result[:500] + "...(truncated)"
					}
					sb.WriteString("```\n")
					sb.WriteString(result)
					sb.WriteString("\n```\n\n")
				}
			}
		}
	}

	return sb.String()
}

// sanitizeFilename converts a topic string to a safe filename component.
func sanitizeFilename(s string) string {
	// Replace spaces and special chars with hyphens
	re := regexp.MustCompile(`[^a-zA-Z0-9-]+`)
	s = re.ReplaceAllString(s, "-")
	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")
	// Convert to lowercase
	s = strings.ToLower(s)
	// Limit length
	if len(s) > 50 {
		s = s[:50]
	}
	if s == "" {
		s = "session"
	}
	return s
}

// buildSessionContent creates the markdown content for a session file.
func buildSessionContent(topic, sessionType string, timestamp time.Time) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# Session: %s\n\n", topic))
	sb.WriteString(fmt.Sprintf("**Date**: %s\n", timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s\n", timestamp.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("**Type**: %s\n", sessionType))
	sb.WriteString("\n---\n\n")

	// Summary section (placeholder for user to fill in)
	sb.WriteString("## Summary\n\n")
	sb.WriteString("[Describe what was accomplished in this session]\n\n")
	sb.WriteString("---\n\n")

	// Current Tasks
	sb.WriteString("## Current Tasks\n\n")
	tasks, err := readContextSection("TASKS.md", "## In Progress", "## Next Up")
	if err == nil && tasks != "" {
		sb.WriteString("### In Progress\n\n")
		sb.WriteString(tasks)
		sb.WriteString("\n")
	}
	nextTasks, err := readContextSection("TASKS.md", "## Next Up", "## Completed")
	if err == nil && nextTasks != "" {
		sb.WriteString("### Next Up\n\n")
		sb.WriteString(nextTasks)
		sb.WriteString("\n")
	}
	sb.WriteString("---\n\n")

	// Recent Decisions
	sb.WriteString("## Recent Decisions\n\n")
	decisions, err := readRecentDecisions()
	if err == nil && decisions != "" {
		sb.WriteString(decisions)
	} else {
		sb.WriteString("[No recent decisions found]\n")
	}
	sb.WriteString("\n---\n\n")

	// Recent Learnings
	sb.WriteString("## Recent Learnings\n\n")
	learnings, err := readRecentLearnings()
	if err == nil && learnings != "" {
		sb.WriteString(learnings)
	} else {
		sb.WriteString("[No recent learnings found]\n")
	}
	sb.WriteString("\n---\n\n")

	// Tasks for Next Session
	sb.WriteString("## Tasks for Next Session\n\n")
	sb.WriteString("[List tasks to continue in the next session]\n\n")

	return sb.String(), nil
}

// readContextSection reads a section from a context file between two headers.
func readContextSection(filename, startHeader, endHeader string) (string, error) {
	filePath := filepath.Join(contextDirName, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Find start
	startIdx := strings.Index(contentStr, startHeader)
	if startIdx == -1 {
		return "", fmt.Errorf("section not found: %s", startHeader)
	}
	startIdx += len(startHeader)

	// Find end
	endIdx := len(contentStr)
	if endHeader != "" {
		idx := strings.Index(contentStr[startIdx:], endHeader)
		if idx != -1 {
			endIdx = startIdx + idx
		}
	}

	section := strings.TrimSpace(contentStr[startIdx:endIdx])
	return section, nil
}

// readRecentDecisions extracts the most recent decisions from DECISIONS.md.
func readRecentDecisions() (string, error) {
	filePath := filepath.Join(contextDirName, "DECISIONS.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Find decision headers (## [YYYY-MM-DD] Title)
	re := regexp.MustCompile(`(?m)^## \[\d{4}-\d{2}-\d{2}\].*$`)
	matches := re.FindAllStringIndex(contentStr, -1)

	if len(matches) == 0 {
		return "", nil
	}

	// Get the last 3 decisions (most recent)
	limit := 3
	if len(matches) < limit {
		limit = len(matches)
	}

	var decisions []string
	for i := len(matches) - limit; i < len(matches); i++ {
		start := matches[i][0]
		end := len(contentStr)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		decision := strings.TrimSpace(contentStr[start:end])
		// Only include the header for brevity
		headerEnd := strings.Index(decision, "\n")
		if headerEnd != -1 {
			decisions = append(decisions, "- "+decision[:headerEnd])
		}
	}

	return strings.Join(decisions, "\n"), nil
}

// readRecentLearnings extracts the most recent learnings from LEARNINGS.md.
func readRecentLearnings() (string, error) {
	filePath := filepath.Join(contextDirName, "LEARNINGS.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Find learning entries (- **[YYYY-MM-DD]** text)
	re := regexp.MustCompile(`(?m)^- \*\*\[\d{4}-\d{2}-\d{2}\]\*\*.*$`)
	matches := re.FindAllString(contentStr, -1)

	if len(matches) == 0 {
		return "", nil
	}

	// Get the last 5 learnings (most recent)
	limit := 5
	if len(matches) < limit {
		limit = len(matches)
	}

	return strings.Join(matches[len(matches)-limit:], "\n"), nil
}

// extractInsights parses a jsonl transcript and extracts potential decisions and learnings.
func extractInsights(path string) ([]string, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var decisions []string
	var learnings []string

	// Patterns for detecting decisions
	decisionPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)decided to\s+(.{20,100})`),
		regexp.MustCompile(`(?i)decision:\s*(.{20,100})`),
		regexp.MustCompile(`(?i)we('ll| will) use\s+(.{10,80})`),
		regexp.MustCompile(`(?i)going with\s+(.{10,80})`),
		regexp.MustCompile(`(?i)chose\s+(.{10,80})\s+(over|instead)`),
	}

	// Patterns for detecting learnings
	learningPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)learned that\s+(.{20,100})`),
		regexp.MustCompile(`(?i)gotcha:\s*(.{20,100})`),
		regexp.MustCompile(`(?i)lesson:\s*(.{20,100})`),
		regexp.MustCompile(`(?i)TIL:?\s*(.{20,100})`),
		regexp.MustCompile(`(?i)turns out\s+(.{20,100})`),
		regexp.MustCompile(`(?i)important to (note|remember):\s*(.{20,100})`),
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	seen := make(map[string]bool) // Deduplicate

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry transcriptEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		// Only look at assistant messages
		if entry.Type != "assistant" {
			continue
		}

		// Extract text content
		texts := extractTextContent(entry)

		for _, text := range texts {
			// Check for decisions
			for _, pattern := range decisionPatterns {
				matches := pattern.FindAllStringSubmatch(text, -1)
				for _, match := range matches {
					if len(match) > 1 {
						insight := cleanInsight(match[1])
						if insight != "" && !seen[insight] {
							seen[insight] = true
							decisions = append(decisions, insight)
						}
					}
				}
			}

			// Check for learnings
			for _, pattern := range learningPatterns {
				matches := pattern.FindAllStringSubmatch(text, -1)
				for _, match := range matches {
					if len(match) > 1 {
						insight := cleanInsight(match[len(match)-1])
						if insight != "" && !seen[insight] {
							seen[insight] = true
							learnings = append(learnings, insight)
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return decisions, learnings, nil
}

// extractTextContent extracts all text content from a transcript entry.
func extractTextContent(entry transcriptEntry) []string {
	var texts []string

	switch content := entry.Message.Content.(type) {
	case string:
		texts = append(texts, content)
	case []interface{}:
		for _, block := range content {
			blockMap, ok := block.(map[string]interface{})
			if !ok {
				continue
			}
			if text, ok := blockMap["text"].(string); ok {
				texts = append(texts, text)
			}
			if thinking, ok := blockMap["thinking"].(string); ok {
				texts = append(texts, thinking)
			}
		}
	}

	return texts
}

// cleanInsight cleans and truncates an extracted insight.
func cleanInsight(s string) string {
	s = strings.TrimSpace(s)
	// Remove trailing punctuation fragments
	s = strings.TrimRight(s, ".,;:!?")
	// Truncate if too long
	if len(s) > 150 {
		// Try to cut at word boundary
		idx := strings.LastIndex(s[:150], " ")
		if idx > 100 {
			s = s[:idx] + "..."
		} else {
			s = s[:147] + "..."
		}
	}
	return s
}
