//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"os"
	"testing"
)

// TestLoopCommand tests the loop command.
func TestLoopCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-loop-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a PROMPT.md file
	if err := os.WriteFile("PROMPT.md", []byte("# Test Prompt\n"), 0644); err != nil {
		t.Fatalf("failed to create PROMPT.md: %v", err)
	}

	// Test loop command
	loopCmd := Cmd()
	loopCmd.SetArgs([]string{"--tool", "generic"})
	if err := loopCmd.Execute(); err != nil {
		t.Fatalf("loop command failed: %v", err)
	}

	// Verify loop.sh was created
	if _, err := os.Stat("loop.sh"); os.IsNotExist(err) {
		t.Error("loop.sh was not created")
	}
}
