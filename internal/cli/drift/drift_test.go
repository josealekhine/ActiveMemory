//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestDriftCommand tests the drift command.
func TestDriftCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Run drift - just verify it runs without error
	driftCmd := Cmd()
	driftCmd.SetArgs([]string{})

	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift command failed: %v", err)
	}
}

// TestDriftJSONOutput tests the drift command with JSON output.
func TestDriftJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-json-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test drift with JSON output
	driftCmd := Cmd()
	driftCmd.SetArgs([]string{"--json"})
	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift --json failed: %v", err)
	}
}
