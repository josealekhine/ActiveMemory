//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestCompleteCommand tests the complete command.
func TestCompleteCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-complete-test-*")
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

	// Add a task
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Task to complete"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task command failed: %v", err)
	}

	// Complete the task
	completeCmd := Cmd()
	completeCmd.SetArgs([]string{"Task to complete"})
	if err := completeCmd.Execute(); err != nil {
		t.Fatalf("complete command failed: %v", err)
	}

	// Verify the task was completed
	tasksPath := filepath.Join(tmpDir, ".context", "TASKS.md")
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read TASKS.md: %v", err)
	}

	if !strings.Contains(string(content), "- [x]") {
		t.Errorf("task was not marked as complete")
	}
}
