//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

// transcriptEntry represents a single entry in the JSONL transcript.
//
// Used when parsing Claude Code transcript files to extract conversation
// history, tool calls, and assistant responses.
//
// Fields:
//   - Type: Entry type ("user", "assistant", or other message types)
//   - Message: The message content with role and content data
//   - Timestamp: ISO 8601 timestamp of when the entry was created
//   - UUID: Unique identifier for the entry
type transcriptEntry struct {
	Type      string        `json:"type"`
	Message   transcriptMsg `json:"message"`
	Timestamp string        `json:"timestamp"`
	UUID      string        `json:"uuid"`
}

// transcriptMsg represents the message content within a transcript entry.
//
// Fields:
//   - Role: Message role ("user" or "assistant")
//   - Content: Message content, either a string or array of content blocks
//     (text, thinking, tool_use, tool_result)
type transcriptMsg struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or []interface{}
}

// sessionInfo holds parsed information about a session file.
//
// Fields:
//   - Filename: Name of the session file
//   - Topic: Session topic extracted from the header
//   - Date: Session date in YYYY-MM-DD format
//   - Type: Session type (feature, bugfix, refactor, session)
//   - Summary: Brief description of what was accomplished
type sessionInfo struct {
	Filename string
	Topic    string
	Date     string
	Type     string
	Summary  string
}
