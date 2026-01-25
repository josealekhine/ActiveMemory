//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import "strings"

// fileNameToTitle converts a context file name to a human-readable title.
//
// Transforms SCREAMING_SNAKE_CASE.md filenames into Title Case strings
// suitable for display (e.g., "TASKS.md" → "Tasks", "AGENT_PLAYBOOK.md" →
// "Agent Playbook").
//
// Parameters:
//   - name: File name to convert (with or without .md extension)
//
// Returns:
//   - string: Title case representation of the file name
func fileNameToTitle(name string) string {
	// Remove .md extension
	name = strings.TrimSuffix(name, ".md")
	// Convert SCREAMING_SNAKE to Title Case
	name = strings.ReplaceAll(name, "_", " ")
	// Title case each word
	words := strings.Fields(name)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, " ")
}
