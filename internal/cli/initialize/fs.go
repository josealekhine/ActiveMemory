//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/templates"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// updateCtxSection replaces the existing ctx section between markers with
// new content.
//
// Locates the ctx markers in the existing content and replaces that section
// with the corresponding section from the template. Creates a timestamped
// backup before modifying.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - existing: Current file content containing ctx markers
//   - newTemplate: Template content with updated ctx section
//
// Returns:
//   - error: Non-nil if the markers are not found or file operations fail
func updateCtxSection(
	cmd *cobra.Command, existing string, newTemplate []byte,
) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Find the start marker
	startIdx := strings.Index(existing, config.CtxMarkerStart)
	if startIdx == -1 {
		return fmt.Errorf("ctx start marker not found")
	}

	// Find the end marker
	endIdx := strings.Index(existing, config.CtxMarkerEnd)
	if endIdx == -1 {
		// No end marker - append from start marker to end
		endIdx = len(existing)
	} else {
		endIdx += len(config.CtxMarkerEnd)
	}

	// Extract the ctx content from the template (between markers)
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, config.CtxMarkerStart)
	templateEnd := strings.Index(templateStr, config.CtxMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return fmt.Errorf("template missing ctx markers")
	}
	ctxContent := templateStr[templateStart : templateEnd+len(config.CtxMarkerEnd)]

	// Build new content: before ctx + new ctx content + after ctx
	newContent := existing[:startIdx] + ctxContent + existing[endIdx:]

	// Back up before updating
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FileClaudeMd, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	cmd.Printf("  %s %s (backup)\n", green("✓"), backupName)

	if err := os.WriteFile(
		config.FileClaudeMd, []byte(newContent), 0644,
	); err != nil {
		return fmt.Errorf("failed to update %s: %w", config.FileClaudeMd, err)
	}
	cmd.Printf(
		"  %s %s (updated ctx section)\n", green("✓"), config.FileClaudeMd,
	)

	return nil
}

// createImplementationPlan creates IMPLEMENTATION_PLAN.md in the project root.
//
// This is the orchestrator directive that points to .context/TASKS.md,
// used by AI agents for task management.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - force: If true, overwrite existing file
//
// Returns:
//   - error: Non-nil if template read or file write fails
func createImplementationPlan(cmd *cobra.Command, force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	const planFileName = "IMPLEMENTATION_PLAN.md"

	// Check if file exists
	if _, err := os.Stat(planFileName); err == nil && !force {
		cmd.Printf("  %s %s (exists, skipped)\n", yellow("○"), planFileName)
		return nil
	}

	// Get template content
	content, err := templates.GetTemplate(planFileName)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	if err := os.WriteFile(planFileName, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	cmd.Printf("  %s %s (orchestrator directive)\n", green("✓"), planFileName)
	return nil
}
