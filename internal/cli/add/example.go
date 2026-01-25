//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import "github.com/ActiveMemory/ctx/internal/config"

// getExamplesForType returns example usage strings for a given entry type.
//
// The examples are displayed in error messages when content is missing,
// helping users understand the correct command syntax.
//
// Parameters:
//   - fileType: Entry type (e.g., "decision", "task", "learning", "convention")
//
// Returns:
//   - string: Formatted example commands; returns a generic example for
//     unrecognized types
func getExamplesForType(fileType string) string {
	switch fileType {
	case config.UpdateTypeDecision, config.UpdateTypeDecisions:
		return `  ctx add decision "Use PostgreSQL for primary database"
  ctx add decision "Adopt Go 1.22 for range-over-func support"`
	case config.UpdateTypeTask, config.UpdateTypeTasks:
		return `  ctx add task "Implement user authentication"
  ctx add task "Fix login bug" --priority high`
	case config.UpdateTypeLearning, config.UpdateTypeLearnings:
		return `  ctx add learning "Vitest mocks must be hoisted above imports"
  ctx add learning "Go embed requires files in same package directory"`
	case config.UpdateTypeConvention, config.UpdateTypeConventions:
		return `  ctx add convention "Use camelCase for function names"
  ctx add convention "All API responses use JSON"`
	default:
		return `  ctx add <type> "your content here"`
	}
}
