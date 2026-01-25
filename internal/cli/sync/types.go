//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

// Action represents a suggested sync action.
//
// Fields:
//   - Type: Category of action (e.g., "DEPS", "CONFIG", "NEW_DIR")
//   - File: Context file that should be updated
//   - Description: What was detected that needs attention
//   - Suggestion: Recommended action to take
type Action struct {
	Type        string
	File        string
	Description string
	Suggestion  string
}
