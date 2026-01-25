//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx loop" command for generating Ralph loop scripts.
//
// The command generates a shell script that runs an AI assistant in a loop
// until a completion signal is detected, enabling iterative development
// where the AI builds on previous work.
//
// Flags:
//   - --prompt, -p: Prompt file to use (default "PROMPT.md")
//   - --tool, -t: AI tool - claude, aider, or generic (default "claude")
//   - --max-iterations, -n: Maximum iterations, 0 for unlimited (default 0)
//   - --completion, -c: Completion signal to detect
//     (default "SYSTEM_CONVERGED")
//   - --output, -o: Output script filename (default "loop.sh")
//
// Returns:
//   - *cobra.Command: Configured loop command with flags registered
func Cmd() *cobra.Command {
	var (
		promptFile    string
		tool          string
		maxIterations int
		completionMsg string
		outputFile    string
	)

	cmd := &cobra.Command{
		Use:   "loop",
		Short: "Generate a Ralph loop script",
		Long: `Generate a ready-to-use shell script for running a Ralph loop.

A Ralph loop continuously runs an AI assistant with the same prompt until
a completion signal is detected. This enables iterative development where
the AI can build on its previous work.

Examples:
  ctx loop                           # Generate loop.sh for Claude
  ctx loop --tool aider              # Generate for Aider
  ctx loop --prompt TASKS.md         # Use custom prompt file
  ctx loop --max-iterations 10       # Limit to 10 iterations
  ctx loop -o my-loop.sh             # Output to custom file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoop(
				cmd, promptFile, tool, maxIterations, completionMsg, outputFile,
			)
		},
	}

	cmd.Flags().StringVarP(&promptFile,
		"prompt", "p", "PROMPT.md", "Prompt file to use",
	)
	cmd.Flags().StringVarP(
		&tool, "tool", "t", "claude", "AI tool: claude, aider, or generic",
	)
	cmd.Flags().IntVarP(
		&maxIterations,
		"max-iterations", "n",
		0, "Maximum iterations (0 = unlimited)",
	)
	cmd.Flags().StringVarP(
		&completionMsg,
		"completion", "c", "SYSTEM_CONVERGED", "Completion signal to detect",
	)
	cmd.Flags().StringVarP(
		&outputFile,
		"output", "o",
		"loop.sh", "Output script filename",
	)

	return cmd
}
