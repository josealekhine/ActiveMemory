//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import "strings"

// cleanInsight cleans and truncates an extracted insight.
//
// Removes leading/trailing whitespace and punctuation, and truncates long
// insights to 150 characters at a word boundary when possible.
//
// Parameters:
//   - s: Raw insight string to clean
//
// Returns:
//   - string: Cleaned and potentially truncated insight
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

// truncate shortens a string to maxLen characters, adding "..." if truncated.
//
// Parameters:
//   - s: String to truncate
//   - maxLen: Maximum length including the ellipsis
//
// Returns:
//   - string: Original string if within maxLen, otherwise truncated with "..."
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
