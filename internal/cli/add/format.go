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
// priority tag. Format: "- [ ] content #priority:level #added:YYYY-MM-DD-HHMMSS"
//
// Parameters:
//   - content: Task description text
//   - priority: Priority level (high, medium, low); empty string omits the tag
//
// Returns:
//   - string: Formatted task line with trailing newline
func FormatTask(content string, priority string) string {
	// Use YYYY-MM-DD-HHMMSS timestamp for session correlation
	timestamp := time.Now().Format("2006-01-02-150405")
	var priorityTag string
	if priority != "" {
		priorityTag = fmt.Sprintf(" #priority:%s", priority)
	}
	return fmt.Sprintf("- [ ] %s%s #added:%s\n", content, priorityTag, timestamp)
}

// FormatLearning formats a learning entry as a structured markdown section.
//
// The output includes a timestamped heading and complete sections for context,
// lesson, and application.
//
// Parameters:
//   - title: Learning title/summary text
//   - context: What prompted this learning
//   - lesson: The key insight
//   - application: How to apply this going forward
//
// Returns:
//   - string: Formatted learning section with all fields
func FormatLearning(title, context, lesson, application string) string {
	timestamp := time.Now().Format("2006-01-02-150405")
	return fmt.Sprintf(`## [%s] %s

**Context**: %s

**Lesson**: %s

**Application**: %s
`, timestamp, title, context, lesson, application)
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
// The output includes a timestamped heading, status, and complete ADR sections
// for context, rationale, and consequences.
//
// Parameters:
//   - title: Decision title/summary text
//   - context: What prompted this decision
//   - rationale: Why this choice over alternatives
//   - consequences: What changes as a result
//
// Returns:
//   - string: Formatted decision section with all ADR fields
func FormatDecision(title, context, rationale, consequences string) string {
	timestamp := time.Now().Format("2006-01-02-150405")
	return fmt.Sprintf(`## [%s] %s

**Status**: Accepted

**Context**: %s

**Decision**: %s

**Rationale**: %s

**Consequences**: %s
`, timestamp, title, context, title, rationale, consequences)
}
