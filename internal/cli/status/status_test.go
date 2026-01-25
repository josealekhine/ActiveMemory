//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestStatusCommand tests the status command.
func TestStatusCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-status-test-*")
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

	// Then status - just verify it runs without error
	statusCmd := Cmd()
	statusCmd.SetArgs([]string{})

	if err := statusCmd.Execute(); err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

// TestStatusJSONOutput tests the status command with JSON output.
func TestStatusJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-status-json-test-*")
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

	// Test status with JSON output
	statusCmd := Cmd()
	statusCmd.SetArgs([]string{"--json"})
	if err := statusCmd.Execute(); err != nil {
		t.Fatalf("status --json failed: %v", err)
	}
}
