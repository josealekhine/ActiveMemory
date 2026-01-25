//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"regexp"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// extractConstitutionRules extracts checkbox items from CONSTITUTION.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: List of constitution rules; nil if the file is not found
func extractConstitutionRules(ctx *context.Context) []string {
	for _, f := range ctx.Files {
		if f.Name == config.FilenameConstitution {
			return extractCheckboxItems(string(f.Content))
		}
	}
	return nil
}

// extractActiveTasks extracts unchecked task items from TASKS.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: List of active tasks with "- [ ]" prefix; nil if
//     the file is not found
func extractActiveTasks(ctx *context.Context) []string {
	for _, f := range ctx.Files {
		if f.Name == config.FilenameTask {
			return extractUncheckedTasks(string(f.Content))
		}
	}
	return nil
}

// extractConventions extracts bullet items from CONVENTIONS.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: Up to 5 convention items; nil if the file is not found
func extractConventions(ctx *context.Context) []string {
	for _, f := range ctx.Files {
		if f.Name == config.FilenameConvention {
			return extractBulletItems(string(f.Content), 5)
		}
	}
	return nil
}

// extractRecentDecisions extracts the most recent decision titles from
// DECISIONS.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//   - limit: Maximum number of decisions to return
//
// Returns:
//   - []string: Decision titles (most recent last); nil if the file
//     is not found
func extractRecentDecisions(
	ctx *context.Context, limit int,
) []string {
	for _, f := range ctx.Files {
		if f.Name == config.FilenameDecision {
			return extractDecisionTitles(string(f.Content), limit)
		}
	}
	return nil
}

// extractCheckboxItems extracts text from Markdown checkbox items.
//
// Matches both checked "- [x]" and unchecked "- [ ]" items.
//
// Parameters:
//   - content: Markdown content to parse
//
// Returns:
//   - []string: Text content of each checkbox item
func extractCheckboxItems(content string) []string {
	re := regexp.MustCompile(`(?m)^-\s*\[[ x]]\s*(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		items = append(items, strings.TrimSpace(m[1]))
	}
	return items
}

// extractUncheckedTasks extracts unchecked Markdown checkbox items.
//
// Only matches "- [ ]" items (not checked). Returns items with the
// "- [ ]" prefix preserved for display.
//
// Parameters:
//   - content: Markdown content to parse
//
// Returns:
//   - []string: Unchecked task items with "- [ ]" prefix
func extractUncheckedTasks(content string) []string {
	re := regexp.MustCompile(`(?m)^-\s*\[\s*]\s*(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		items = append(items, "- [ ] "+strings.TrimSpace(m[1]))
	}
	return items
}

// extractBulletItems extracts Markdown bullet items up to a limit.
//
// Skips empty items and lines starting with "#" (headers).
//
// Parameters:
//   - content: Markdown content to parse
//   - limit: Maximum number of items to return
//
// Returns:
//   - []string: Bullet item text without the "- " prefix
func extractBulletItems(content string, limit int) []string {
	re := regexp.MustCompile(`(?m)^-\s+(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, limit)
	for i, m := range matches {
		if i >= limit {
			break
		}
		text := strings.TrimSpace(m[1])
		// Skip empty or header-only items
		if text != "" && !strings.HasPrefix(text, "#") {
			items = append(items, text)
		}
	}
	return items
}

// extractDecisionTitles extracts decision titles from Markdown headings.
//
// Matches headings in the format "## [YYYY-MM-DD] Title" and returns
// the most recent decisions (those appearing last in the file).
//
// Parameters:
//   - content: Markdown content to parse
//   - limit: Maximum number of decision titles to return
//
// Returns:
//   - []string: Decision titles without a timestamp prefix
func extractDecisionTitles(content string, limit int) []string {
	re := regexp.MustCompile(`(?m)^##\s+\[[\d-]+]\s*(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, limit)
	// Get the most recent (last) decisions
	start := len(matches) - limit
	if start < 0 {
		start = 0
	}
	for i := start; i < len(matches); i++ {
		items = append(items, strings.TrimSpace(matches[i][1]))
	}
	return items
}
