//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"time"
)

// FormatTask formats a task entry as a markdown checkbox item.
//
// The output includes a timestamp tag for session correlation and an optional
// priority tag. Format: "- [ ] content #priority:level #added:YYYY-MM-DD-HHMM"
//
// Parameters:
//   - content: Task description text
//   - priority: Priority level (high, medium, low); empty string omits the tag
//
// Returns:
//   - string: Formatted task line with trailing newline
func FormatTask(content string, priority string) string {
	// Use YYYY-MM-DD-HHMM timestamp for session correlation
	timestamp := time.Now().Format("2006-01-02-1504")
	var priorityTag string
	if priority != "" {
		priorityTag = fmt.Sprintf(" #priority:%s", priority)
	}
	return fmt.Sprintf("- [ ] %s%s #added:%s\n", content, priorityTag, timestamp)
}

// FormatLearning formats a learning entry as a timestamped markdown list item.
//
// Format: "- **[YYYY-MM-DD-HHMM]** content"
//
// Parameters:
//   - content: Learning description text
//
// Returns:
//   - string: Formatted learning line with trailing newline
func FormatLearning(content string) string {
	timestamp := time.Now().Format("2006-01-02-1504")
	return fmt.Sprintf("- **[%s]** %s\n", timestamp, content)
}

// FormatConvention formats a convention entry as a simple markdown list item.
//
// Format: "- content"
//
// Parameters:
//   - content: Convention description text
//
// Returns:
//   - string: Formatted convention line with trailing newline
func FormatConvention(content string) string {
	return fmt.Sprintf("- %s\n", content)
}

// FormatDecision formats a decision entry as a structured Markdown section.
//
// The output includes a timestamped heading, status, and placeholder sections
// for context, rationale, and consequences of the ADR format.
//
// Parameters:
//   - content: Decision title/summary text
//
// Returns:
//   - string: Formatted decision section with placeholders for details
func FormatDecision(content string) string {
	timestamp := time.Now().Format("2006-01-02-1504")
	return fmt.Sprintf(`## [%s] %s

**Status**: Accepted

**Context**: [Add context here]

**Decision**: %s

**Rationale**: [Add rationale here]

**Consequences**: [Add consequences here]
`, timestamp, content, content)
}
