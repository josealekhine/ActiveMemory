//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/config"
)

// applyTaskUpdate appends a task entry to TASKS.md.
//
// Parameters:
//   - content: Task description text
//
// Returns:
//   - error: Non-nil if the file operation fails
func applyTaskUpdate(content string) error {
	args := []string{config.UpdateTypeTask, content}
	return runAddSilent(args)
}

// applyDecisionUpdate appends a decision entry to DECISIONS.md.
//
// Parameters:
//   - content: Decision description text
//
// Returns:
//   - error: Non-nil if the file operation fails
func applyDecisionUpdate(content string) error {
	args := []string{config.UpdateTypeDecision, content}
	return runAddSilent(args)
}

// applyLearningUpdate appends a learning entry to LEARNINGS.md.
//
// Parameters:
//   - content: Learning description text
//
// Returns:
//   - error: Non-nil if the file operation fails
func applyLearningUpdate(content string) error {
	args := []string{config.UpdateTypeLearning, content}
	return runAddSilent(args)
}

// applyConventionUpdate appends a convention entry to CONVENTIONS.md.
//
// Parameters:
//   - content: Convention description text
//
// Returns:
//   - error: Non-nil if the file operation fails
func applyConventionUpdate(content string) error {
	args := []string{config.UpdateTypeConvention, content}
	return runAddSilent(args)
}

// applyCompleteUpdate marks a matching task as complete in TASKS.md.
//
// Parameters:
//   - content: Search query to match against task descriptions
//
// Returns:
//   - error: Non-nil if no matching task is found or file operation fails
func applyCompleteUpdate(content string) error {
	args := []string{content}
	return runCompleteSilent(args)
}

// applyUpdate routes a context update to the appropriate handler.
//
// Dispatches based on update type to add entries to context files
// or mark tasks complete.
//
// Parameters:
//   - update: ContextUpdate containing type and content
//
// Returns:
//   - error: Non-nil if type is unknown or the handler fails
func applyUpdate(update ContextUpdate) error {
	switch update.Type {
	case config.UpdateTypeTask:
		return applyTaskUpdate(update.Content)
	case config.UpdateTypeDecision:
		return applyDecisionUpdate(update.Content)
	case config.UpdateTypeLearning:
		return applyLearningUpdate(update.Content)
	case config.UpdateTypeConvention:
		return applyConventionUpdate(update.Content)
	case config.UpdateTypeComplete:
		return applyCompleteUpdate(update.Content)
	default:
		return fmt.Errorf("unknown update type: %s", update.Type)
	}
}
