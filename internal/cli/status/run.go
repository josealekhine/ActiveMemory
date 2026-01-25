//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/spf13/cobra"
)

// runStatus executes the status command logic.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON
//   - verbose: If true, include file content previews
//
// Returns:
//   - error: Non-nil if context loading fails
func runStatus(cmd *cobra.Command, jsonOutput, verbose bool) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	if jsonOutput {
		return outputStatusJSON(cmd, ctx, verbose)
	}

	return outputStatusText(cmd, ctx, verbose)
}
