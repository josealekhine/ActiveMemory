//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

// Output represents the JSON output format for the status command.
//
// Fields:
//   - ContextDir: Path to the .context/ directory
//   - TotalFiles: Number of context files found
//   - TotalTokens: Estimated total token count across all files
//   - TotalSize: Total size in bytes across all files
//   - Files: Individual file status entries
type Output struct {
	ContextDir  string       `json:"context_dir"`
	TotalFiles  int          `json:"total_files"`
	TotalTokens int          `json:"total_tokens"`
	TotalSize   int64        `json:"total_size"`
	Files       []FileStatus `json:"files"`
}

// FileStatus represents a single file's status in JSON output.
//
// Fields:
//   - Name: Filename (e.g., "TASKS.md")
//   - Tokens: Estimated token count for this file
//   - Size: File size in bytes
//   - IsEmpty: True if the file has no meaningful content
//   - Summary: Brief description of file contents
//   - ModTime: Last modification time (RFC3339 format)
//   - Preview: Content preview lines (only with --verbose)
type FileStatus struct {
	Name    string   `json:"name"`
	Tokens  int      `json:"tokens"`
	Size    int64    `json:"size"`
	IsEmpty bool     `json:"is_empty"`
	Summary string   `json:"summary"`
	ModTime string   `json:"mod_time"`
	Preview []string `json:"preview,omitempty"`
}
