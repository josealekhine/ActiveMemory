//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"bufio"
	"regexp"
	"strings"
)

// separateTasks parses TASKS.md and separates completed from pending tasks.
//
// The function scans TASKS.md line by line, identifying task items by their
// checkbox markers ([x] for completed, [ ] for pending). It preserves phase
// headers (### Phase ...) in the archived content for traceability.
//
// Subtasks (indented task items) follow their parent task:
//   - Subtasks of completed tasks are archived with the parent
//   - Subtasks of pending tasks remain with the parent
//
// Parameters:
//   - content: Full content of TASKS.md as a string
//
// Returns:
//   - remaining: Content with only pending tasks (to write back to TASKS.md)
//   - archived: Content with completed tasks and their phase headers
//   - stats: Counts of completed and pending tasks processed
func separateTasks(content string) (string, string, taskStats) {
	var remaining strings.Builder
	var archived strings.Builder
	var stats taskStats

	// Track the current phase header
	var currentPhase string
	var phaseHasArchivedTasks bool
	var phaseArchiveBuffer strings.Builder

	completedPattern := regexp.MustCompile(`^\s*-\s*\[x\]`)
	pendingPattern := regexp.MustCompile(`^\s*-\s*\[\s*\]`)
	phasePattern := regexp.MustCompile(`^###\s+Phase`)
	subTaskPattern := regexp.MustCompile(`^\s{2,}-\s*\[`)

	scanner := bufio.NewScanner(strings.NewReader(content))
	var inCompletedTask bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check for phase headers
		if phasePattern.MatchString(line) {
			// Flush previous phase's archived tasks
			if phaseHasArchivedTasks {
				archived.WriteString(currentPhase + "\n")
				archived.WriteString(phaseArchiveBuffer.String())
				archived.WriteString("\n")
			}

			currentPhase = line
			phaseHasArchivedTasks = false
			phaseArchiveBuffer.Reset()
			remaining.WriteString(line + "\n")
			inCompletedTask = false
			continue
		}

		// Check for completed tasks
		if completedPattern.MatchString(line) {
			stats.completed++
			phaseHasArchivedTasks = true
			phaseArchiveBuffer.WriteString(line + "\n")
			inCompletedTask = true
			continue
		}

		// Check for pending tasks
		if pendingPattern.MatchString(line) {
			stats.pending++
			remaining.WriteString(line + "\n")
			inCompletedTask = false
			continue
		}

		// Handle subtasks (indented task items)
		if subTaskPattern.MatchString(line) {
			if inCompletedTask {
				// Subtask of a completed task - archive it
				phaseArchiveBuffer.WriteString(line + "\n")
			} else {
				// Subtask of a pending task - keep it
				remaining.WriteString(line + "\n")
			}
			continue
		}

		// Non-task lines go to the remaining
		remaining.WriteString(line + "\n")
		inCompletedTask = false
	}

	// Flush final phase's archived tasks
	if phaseHasArchivedTasks {
		archived.WriteString(currentPhase + "\n")
		archived.WriteString(phaseArchiveBuffer.String())
	}

	return remaining.String(), archived.String(), stats
}
