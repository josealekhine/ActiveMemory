//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// parseIndex attempts to parse a string as a positive integer index.
//
// Parameters:
//   - s: String to parse as an integer
//
// Returns:
//   - int: Parsed positive integer
//   - error: Non-nil if parsing fails or index is not positive
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

// parseJsonlTranscript parses a .jsonl file and returns formatted Markdown.
//
// Reads a Claude Code JSONL transcript and converts it to readable Markdown
// with message headers, timestamps, and formatted content blocks.
//
// Parameters:
//   - path: Path to the JSONL transcript file
//
// Returns:
//   - string: Markdown-formatted transcript
//   - error: Non-nil if the file cannot be opened or read
func parseJsonlTranscript(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close file: %w", err)
		}
	}(file)

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

// parseSessionFile extracts metadata from a session file.
//
// Parses a session Markdown file to extract topic, date, type, and summary
// information. Handles both "# Session: topic" and "# topic" header formats.
//
// Parameters:
//   - path: Path to the session file
//
// Returns:
//   - sessionInfo: Parsed session metadata
//   - error: Non-nil if file cannot be read
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
			if line != "" && !strings.HasPrefix(line, "#") &&
				!strings.HasPrefix(line, "---") &&
				!strings.HasPrefix(line, "[") {
				info.Summary = line
				break
			}
		}
	}

	return info, nil
}
