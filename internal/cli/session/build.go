//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// buildSessionContent creates the Markdown content for a session file.
//
// Assembles a session document with metadata, current tasks, recent decisions,
// and learnings. Uses CTX_SESSION_START environment variable for session
// correlation if available.
//
// Parameters:
//   - topic: Session topic used as the document title
//   - sessionType: Type of session (e.g., "manual", "auto-save")
//   - timestamp: Time used for end_time and fallback start_time
//
// Returns:
//   - string: Complete Markdown content for the session file
//   - error: Currently always nil (reserved for future validation)
func buildSessionContent(
	topic, sessionType string, timestamp time.Time,
) (string, error) {
	var sb strings.Builder

	// Header with timestamp fields for session correlation
	sb.WriteString(fmt.Sprintf("# Session: %s\n\n", topic))
	sb.WriteString(fmt.Sprintf("**Date**: %s\n", timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s\n", timestamp.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("**Type**: %s\n", sessionType))

	// Session correlation timestamps
	// (YYYY-MM-DD-HHMM format matches ctx add timestamps)
	// start_time: When session began
	// (use CTX_SESSION_START env var if available, else save time)
	startTime := timestamp
	if envStart := os.Getenv("CTX_SESSION_START"); envStart != "" {
		if parsed, err := time.Parse("2006-01-02-1504", envStart); err == nil {
			startTime = parsed
		}
	}
	sb.WriteString(
		fmt.Sprintf("**start_time**: %s\n", startTime.Format("2006-01-02-1504")),
	)
	sb.WriteString(
		fmt.Sprintf("**end_time**: %s\n", timestamp.Format("2006-01-02-1504")),
	)
	sb.WriteString("\n---\n\n")

	// Summary section (placeholder for the user to fill in)
	sb.WriteString("## Summary\n\n")
	sb.WriteString("[Describe what was accomplished in this session]\n\n")
	sb.WriteString("---\n\n")

	// Current Tasks
	sb.WriteString("## Current Tasks\n\n")
	tasks, err := readContextSection(
		"TASKS.md", "## In Progress", "## Next Up",
	)
	if err == nil && tasks != "" {
		sb.WriteString("### In Progress\n\n")
		sb.WriteString(tasks)
		sb.WriteString("\n")
	}
	nextTasks, err := readContextSection(
		"TASKS.md", "## Next Up", "## Completed",
	)
	if err == nil && nextTasks != "" {
		sb.WriteString("### Next Up\n\n")
		sb.WriteString(nextTasks)
		sb.WriteString("\n")
	}
	sb.WriteString("---\n\n")

	// Recent Decisions
	sb.WriteString("## Recent Decisions\n\n")
	decisions, err := readRecentDecisions()
	if err == nil && decisions != "" {
		sb.WriteString(decisions)
	} else {
		sb.WriteString("[No recent decisions found]\n")
	}
	sb.WriteString("\n---\n\n")

	// Recent Learnings
	sb.WriteString("## Recent Learnings\n\n")
	learnings, err := readRecentLearnings()
	if err == nil && learnings != "" {
		sb.WriteString(learnings)
	} else {
		sb.WriteString("[No recent learnings found]\n")
	}
	sb.WriteString("\n---\n\n")

	// Tasks for Next Session
	sb.WriteString("## Tasks for Next Session\n\n")
	sb.WriteString("[List tasks to continue in the next session]\n\n")

	return sb.String(), nil
}
