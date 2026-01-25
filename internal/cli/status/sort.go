//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"sort"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// sortFilesByPriority sorts files in-place by the recommended read order.
//
// Uses config.FilePriority to determine ordering (CONSTITUTION first,
// then TASKS, CONVENTIONS, etc.).
//
// Parameters:
//   - files: Slice of files to sort (modified in place)
func sortFilesByPriority(files []context.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return config.FilePriority(
			files[i].Name,
		) < config.FilePriority(files[j].Name)
	})
}

// getRecentFiles returns the n most recently modified files.
//
// Parameters:
//   - files: Source files to select from
//   - n: Maximum number of files to return
//
// Returns:
//   - []context.FileInfo: Up to n files sorted by modification time
//     (newest first)
func getRecentFiles(files []context.FileInfo, n int) []context.FileInfo {
	sorted := make([]context.FileInfo, len(files))
	copy(sorted, files)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ModTime.After(sorted[j].ModTime)
	})
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	return sorted
}
