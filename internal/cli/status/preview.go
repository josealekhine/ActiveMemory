//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import "strings"

// getContentPreview returns the first n non-empty, meaningful lines
// from content.
//
// Skips empty lines, YAML frontmatter delimiters, and HTML comments.
// Truncates lines longer than 60 characters.
//
// Parameters:
//   - content: The file content to extract preview from
//   - n: Maximum number of lines to return
//
// Returns:
//   - []string: Up to n meaningful lines from the content
func getContentPreview(content string, n int) []string {
	lines := strings.Split(content, "\n")
	var preview []string

	inFrontmatter := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip YAML frontmatter
		if trimmed == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			continue
		}

		// Skip HTML comments
		if strings.HasPrefix(trimmed, "<!--") {
			continue
		}

		// Truncate long lines
		if len(trimmed) > 60 {
			trimmed = trimmed[:57] + "..."
		}

		preview = append(preview, trimmed)
		if len(preview) >= n {
			break
		}
	}

	return preview
}
