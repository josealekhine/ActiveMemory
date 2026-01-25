//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// checkCtxInPath verifies that ctx is available in PATH.
//
// The hooks use "ctx" expecting it to be in PATH, so init should fail
// if the user hasn't installed ctx globally yet.
// Set CTX_SKIP_PATH_CHECK=1 to skip this check (used in tests).
//
// Parameters:
//   - cmd: Cobra command for error output stream
//
// Returns:
//   - error: Non-nil if ctx is not found in PATH
func checkCtxInPath(cmd *cobra.Command) error {
	// Allow skipping for tests
	if os.Getenv("CTX_SKIP_PATH_CHECK") == "1" {
		return nil
	}

	_, err := exec.LookPath("ctx")
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		cmd.PrintErrf("%s ctx is not in your PATH\n\n", red("Error:"))
		cmd.PrintErrln(
			"The hooks created by 'ctx init' require ctx to be in your PATH.",
		)
		cmd.PrintErrln("Without this, Claude Code hooks will fail silently.")
		cmd.PrintErrln()
		cmd.PrintErrf("%s\n", yellow("To fix this:"))
		cmd.PrintErrln("  1. Build:   make build")
		cmd.PrintErrln("  2. Install: sudo make install")
		cmd.PrintErrln()
		cmd.PrintErrln("Or manually:")
		cmd.PrintErrln("  sudo cp ./ctx /usr/local/bin/")
		cmd.PrintErrln()
		cmd.PrintErrln("Then run 'ctx init' again.")

		return fmt.Errorf("ctx not found in PATH")
	}
	return nil
}
