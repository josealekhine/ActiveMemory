//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// generateSummary creates a brief summary for a context file based on its name and content.
func generateSummary(name string, content []byte) string {
	switch name {
	case config.FilenameConstitution:
		return summarizeConstitution(content)
	case config.FilenameTask:
		return summarizeTasks(content)
	case config.FilenameDecision:
		return summarizeDecisions(content)
	case config.FilenameGlossary:
		return summarizeGlossary(content)
	default:
		if len(content) == 0 || isEffectivelyEmpty(content) {
			return "empty"
		}
		return "loaded"
	}
}

func summarizeConstitution(content []byte) string {
	// Count checkbox items (invariants)
	count := bytes.Count(content, []byte("- [ ]")) + bytes.Count(content, []byte("- [x]"))
	if count == 0 {
		return "loaded"
	}
	return fmt.Sprintf("%d invariants", count)
}

func summarizeTasks(content []byte) string {
	// Count active (unchecked) and completed (checked) tasks
	active := bytes.Count(content, []byte("- [ ]"))
	completed := bytes.Count(content, []byte("- [x]"))

	if active == 0 && completed == 0 {
		return "empty"
	}

	var parts []string
	if active > 0 {
		parts = append(parts, fmt.Sprintf("%d active", active))
	}
	if completed > 0 {
		parts = append(parts, fmt.Sprintf("%d completed", completed))
	}
	return strings.Join(parts, ", ")
}

func summarizeDecisions(content []byte) string {
	// Count decision headers (## [date] or ## Decision)
	re := regexp.MustCompile(`(?m)^## `)
	matches := re.FindAll(content, -1)
	count := len(matches)

	if count == 0 {
		return "empty"
	}
	if count == 1 {
		return "1 decision"
	}
	return fmt.Sprintf("%d decisions", count)
}

func summarizeGlossary(content []byte) string {
	// Count definition entries (lines starting with **term** or - **term**)
	re := regexp.MustCompile(`(?m)(?:^|\n)\s*(?:-\s*)?\*\*[^*]+\*\*`)
	matches := re.FindAll(content, -1)
	count := len(matches)

	if count == 0 {
		return "empty"
	}
	if count == 1 {
		return "1 term"
	}
	return fmt.Sprintf("%d terms", count)
}
