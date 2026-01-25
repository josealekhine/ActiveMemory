//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/context"
)

// runAgent executes the agent command logic.
//
// Loads context from .context/ and outputs a context packet in the
// specified format (Markdown or JSON).
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - budget: Token budget to include in the output
//   - format: Output format, "json" for JSON, or any other value for Markdown
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func runAgent(cmd *cobra.Command, budget int, format string) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	if format == "json" {
		return outputAgentJSON(cmd, ctx, budget)
	}

	return outputAgentMarkdown(cmd, ctx, budget)
}
