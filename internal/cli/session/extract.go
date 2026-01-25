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
	"regexp"
)

// extractInsights parses a JSONL transcript and extracts potential decisions
// and learnings.
//
// Scans assistant messages for patterns indicating decisions
// (e.g., "decided to", "we'll use", "chose X over Y") and learnings
// (e.g., "learned that", "gotcha", "TIL"). Results are deduplicated.
//
// Parameters:
//   - path: Path to the JSONL transcript file
//
// Returns:
//   - []string: Extracted decision insights
//   - []string: Extracted learning insights
//   - error: Non-nil if file cannot be opened or read
func extractInsights(path string) ([]string, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_ = fmt.Errorf("error closing file: %v", err)
		}
	}(file)

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
//
// Handles both string content and array content (with text and thinking blocks).
//
// Parameters:
//   - entry: Transcript entry to extract text from
//
// Returns:
//   - []string: All text content found in the entry
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
