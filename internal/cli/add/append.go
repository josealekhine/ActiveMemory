//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import "strings"

// AppendEntry inserts a formatted entry into existing file content.
//
// For task entries, the function locates the target section header and inserts
// the entry immediately after it. For all other entry types, the entry is
// appended to the end of the file with appropriate newline handling.
//
// Parameters:
//   - existing: Current file content as bytes
//   - entry: Pre-formatted entry text to insert
//   - fileType: Entry type (e.g., "task", "decision", "learning", "convention")
//   - section: Target section header for tasks; defaults to "## Next Up" if
//     empty; a "## " prefix is added automatically if missing
//
// Returns:
//   - []byte: Modified file content with the entry inserted
func AppendEntry(
	existing []byte, entry string, fileType string, section string,
) []byte {
	existingStr := string(existing)

	// For tasks, find the appropriate section
	if fileType == "task" || fileType == "tasks" {
		targetSection := section
		if targetSection == "" {
			targetSection = "## Next Up"
		} else if !strings.HasPrefix(targetSection, "##") {
			targetSection = "## " + targetSection
		}

		// Find the section and insert after it
		idx := strings.Index(existingStr, targetSection)
		if idx != -1 {
			// Find the end of the section header line
			lineEnd := strings.Index(existingStr[idx:], "\n")
			if lineEnd != -1 {
				insertPoint := idx + lineEnd + 1
				return []byte(existingStr[:insertPoint] + "\n" +
					entry + existingStr[insertPoint:])
			}
		}
	}

	// For decisions, insert before the closing comment if present,
	// otherwise append
	if fileType == "decision" || fileType == "decisions" {
		// Just append at the end
		if !strings.HasSuffix(existingStr, "\n") {
			existingStr += "\n"
		}
		return []byte(existingStr + "\n" + entry)
	}

	// Default: append at the end
	if !strings.HasSuffix(existingStr, "\n") {
		existingStr += "\n"
	}
	return []byte(existingStr + "\n" + entry)
}
