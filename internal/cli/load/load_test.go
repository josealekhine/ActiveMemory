//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestLoadCommand tests the load command.
func TestLoadCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-load-test-*")
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

	// Run load - just verify it runs without error
	loadCmd := Cmd()
	loadCmd.SetArgs([]string{})

	if err := loadCmd.Execute(); err != nil {
		t.Fatalf("load command failed: %v", err)
	}
}

// TestLoadRawOutput tests the load command with raw output.
func TestLoadRawOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-load-raw-test-*")
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

	// Test load with raw output
	loadCmd := Cmd()
	loadCmd.SetArgs([]string{"--raw"})
	if err := loadCmd.Execute(); err != nil {
		t.Fatalf("load --raw failed: %v", err)
	}
}
