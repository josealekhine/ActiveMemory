//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

// ContextUpdate represents a parsed context update command.
//
// Extracted from <context-update> XML tags in the input stream.
//
// Fields:
//   - Type: Update type (task, decision, learning, convention, complete)
//   - Content: The entry text or search query for complete
type ContextUpdate struct {
	Type    string
	Content string
}
