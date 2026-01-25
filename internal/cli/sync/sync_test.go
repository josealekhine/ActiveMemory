//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestSyncCommand tests the sync command.
func TestSyncCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-sync-test-*")
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

	// Test sync command
	syncCmd := Cmd()
	syncCmd.SetArgs([]string{})
	if err := syncCmd.Execute(); err != nil {
		t.Fatalf("sync command failed: %v", err)
	}
}
