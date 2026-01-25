//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"strings"
	"time"
)

// formatTranscriptEntry formats a single transcript entry as Markdown.
//
// Converts a transcript entry into readable Markdown with role headers,
// timestamps, and formatted content blocks. Handles text, thinking (as
// collapsible details), tool_use (with parameters), and tool_result blocks.
//
// Parameters:
//   - entry: Transcript entry to format
//
// Returns:
//   - string: Markdown-formatted representation of the entry
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
