//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import "strings"

// removeEmptySections removes Markdown sections that contain no content.
//
// A section is considered empty if it has a "## " header followed only by
// blank lines until the next section or end of the file.
//
// Parameters:
//   - content: Markdown content to process
//
// Returns:
//   - string: Content with empty sections removed
//   - int: Number of sections removed
func removeEmptySections(content string) (string, int) {
	lines := strings.Split(content, "\n")
	var result []string
	removed := 0

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this is a section header
		if strings.HasPrefix(line, "## ") {
			// Look ahead to see if the section is empty
			sectionStart := i
			i++

			// Skip empty lines
			for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				i++
			}

			// Check if we hit another section or end of the file
			if i >= len(lines) ||
				strings.HasPrefix(lines[i], "## ") ||
				strings.HasPrefix(lines[i], "# ") {
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

// truncateString shortens a string to maxLen, adding "..." if truncated.
//
// Parameters:
//   - s: String to truncate
//   - maxLen: Maximum length including the "..." suffix
//
// Returns:
//   - string: Original string if within limit, otherwise truncated with "..."
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
