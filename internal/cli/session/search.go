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
)

// findSessionFile finds a session file matching the query.
//
// The query can be:
//   - A numeric index (1 = most recent, 2 = second most recent, etc.)
//   - An exact filename match
//   - A partial match (case-insensitive substring)
//
// Parameters:
//   - query: Index, filename, or partial match string
//
// Returns:
//   - string: Full path to the matched session file
//   - error: Non-nil if no match is found, multiple matches,
//     or index out of range
func findSessionFile(query string) (string, error) {
	// Read directory
	entries, err := os.ReadDir(sessionsDirPath())
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

	// Check if the query is a number (index)
	if idx, err := parseIndex(query); err == nil {
		if idx < 1 || idx > len(sessions) {
			return "", fmt.Errorf("index %d out of range (1-%d)", idx, len(sessions))
		}
		return filepath.Join(sessionsDirPath(), sessions[idx-1]), nil
	}

	// Check for the exact match
	for _, name := range sessions {
		if name == query {
			return filepath.Join(sessionsDirPath(), name), nil
		}
	}

	// Check for a partial match
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

	return filepath.Join(sessionsDirPath(), matches[0]), nil
}
