//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestSeparateTasks tests the separateTasks helper function.
func TestSeparateTasks(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedCompleted int
		expectedPending   int
	}{
		{
			name:              "mixed tasks",
			input:             "# Tasks\n\n### Phase 1\n- [x] Done task\n- [ ] Pending task\n",
			expectedCompleted: 1,
			expectedPending:   1,
		},
		{
			name:              "all completed",
			input:             "# Tasks\n\n- [x] Task 1\n- [x] Task 2\n",
			expectedCompleted: 2,
			expectedPending:   0,
		},
		{
			name:              "all pending",
			input:             "# Tasks\n\n- [ ] Task 1\n- [ ] Task 2\n",
			expectedCompleted: 0,
			expectedPending:   2,
		},
		{
			name:              "no tasks",
			input:             "# Tasks\n\nNo tasks here.\n",
			expectedCompleted: 0,
			expectedPending:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, stats := separateTasks(tt.input)
			if stats.completed != tt.expectedCompleted {
				t.Errorf("separateTasks() completed = %d, want %d", stats.completed, tt.expectedCompleted)
			}
			if stats.pending != tt.expectedPending {
				t.Errorf("separateTasks() pending = %d, want %d", stats.pending, tt.expectedPending)
			}
		})
	}
}

// TestTasksCommands tests the tasks subcommands.
func TestTasksCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-tasks-test-*")
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

	// Add some tasks
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Test task 1"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task failed: %v", err)
	}

	// Test tasks snapshot
	t.Run("tasks snapshot", func(t *testing.T) {
		tasksCmd := Cmd()
		tasksCmd.SetArgs([]string{"snapshot", "test-snapshot"})
		if err := tasksCmd.Execute(); err != nil {
			t.Fatalf("tasks snapshot failed: %v", err)
		}

		// Verify snapshot was created
		entries, err := os.ReadDir(".context/archive")
		if err != nil {
			t.Fatalf("failed to read archive dir: %v", err)
		}
		found := false
		for _, e := range entries {
			if strings.Contains(e.Name(), "test-snapshot") {
				found = true
				break
			}
		}
		if !found {
			t.Error("snapshot file was not created")
		}
	})

	// Test tasks archive (dry-run)
	t.Run("tasks archive dry-run", func(t *testing.T) {
		tasksCmd := Cmd()
		tasksCmd.SetArgs([]string{"archive", "--dry-run"})
		if err := tasksCmd.Execute(); err != nil {
			t.Fatalf("tasks archive failed: %v", err)
		}
	})
}
