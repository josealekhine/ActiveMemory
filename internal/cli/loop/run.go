//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runLoop executes the loop command logic.
//
// Validates the tool selection, generates the loop script, and writes it
// to the output file. Prints usage instructions after generation.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - promptFile: Path to the prompt file for the AI
//   - tool: AI tool to use (claude, aider, or generic)
//   - maxIterations: Maximum loop iterations (0 for unlimited)
//   - completionMsg: Signal string that indicates loop completion
//   - outputFile: Path for the generated script
//
// Returns:
//   - error: Non-nil if the tool is invalid or file write fails
func runLoop(
	cmd *cobra.Command,
	promptFile, tool string,
	maxIterations int,
	completionMsg, outputFile string,
) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Validate tool
	validTools := map[string]bool{"claude": true, "aider": true, "generic": true}
	if !validTools[tool] {
		return fmt.Errorf(
			"invalid tool %q: must be claude, aider, or generic", tool,
		)
	}

	// Generate the script
	script := generateLoopScript(promptFile, tool, maxIterations, completionMsg)

	// Write to the file
	if err := os.WriteFile(outputFile, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputFile, err)
	}

	cmd.Printf("%s Generated %s\n", green("✓"), outputFile)
	cmd.Println()
	cmd.Println("To start the loop:")
	cmd.Printf("  ./%s\n", outputFile)
	cmd.Println()
	cmd.Printf("Tool: %s\n", tool)
	cmd.Printf("Prompt: %s\n", promptFile)
	if maxIterations > 0 {
		cmd.Printf("Max iterations: %d\n", maxIterations)
	} else {
		cmd.Println("Max iterations: unlimited")
	}
	cmd.Printf("Completion signal: %s\n", completionMsg)

	return nil
}
