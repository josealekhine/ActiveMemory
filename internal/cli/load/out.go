//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/context"
)

// outputRaw outputs context files without assembly or headers.
//
// Files are output in read order (see [config.FileReadOrder]), separated
// by blank lines. Content is printed as-is without modification.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing files to output
//
// Returns:
//   - error: Always nil (included for interface consistency)
func outputRaw(cmd *cobra.Command, ctx *context.Context) error {
	// Sort files by read order
	files := sortByReadOrder(ctx.Files)

	for i, f := range files {
		if i > 0 {
			cmd.Println()
		}
		cmd.Print(string(f.Content))
	}
	return nil
}

// outputAssembled outputs context as formatted Markdown with token budgeting.
//
// Assembles context files into a single Markdown document with headers,
// respecting the token budget. Files are included in read order until the
// budget is exhausted. Truncated files are noted in the output.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing files to assemble
//   - budget: Maximum token count for the output
//
// Returns:
//   - error: Always nil (included for interface consistency)
func outputAssembled(
	cmd *cobra.Command, ctx *context.Context, budget int,
) error {
	var sb strings.Builder

	// Header
	sb.WriteString("# Context\n\n")
	sb.WriteString(
		fmt.Sprintf(
			"Token Budget: %d | Available: %d\n\n",
			budget, ctx.TotalTokens,
		),
	)
	sb.WriteString("---\n\n")

	// Sort files by read order
	files := sortByReadOrder(ctx.Files)

	tokensUsed := context.EstimateTokensString(sb.String())

	for _, f := range files {
		// Skip empty files
		if f.IsEmpty {
			continue
		}

		// Check if we have the budget for this file
		fileTokens := f.Tokens
		if tokensUsed+fileTokens > budget {
			// Add a truncation notice
			sb.WriteString(
				fmt.Sprintf("\n---\n\n*[Truncated: %s and remaining files "+
					"excluded due to token budget]*\n", f.Name),
			)
			break
		}

		// Add the file section
		sb.WriteString(fmt.Sprintf("## %s\n\n", fileNameToTitle(f.Name)))
		sb.Write(f.Content)
		if !strings.HasSuffix(string(f.Content), "\n") {
			sb.WriteString("\n")
		}
		sb.WriteString("\n---\n\n")

		tokensUsed += fileTokens
	}

	cmd.Print(sb.String())
	return nil
}
