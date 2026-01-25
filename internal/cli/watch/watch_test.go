//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config"
)

// TestApplyUpdate tests the applyUpdate function routing.
func TestApplyUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-apply-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	tests := []struct {
		name        string
		update      ContextUpdate
		checkFile   string
		checkFor    string
		expectError bool
	}{
		{
			name:      "task update",
			update:    ContextUpdate{Type: config.UpdateTypeTask, Content: "Test task from watch"},
			checkFile: config.FilenameTask,
			checkFor:  "Test task from watch",
		},
		{
			name:      "decision update",
			update:    ContextUpdate{Type: config.UpdateTypeDecision, Content: "Test decision from watch"},
			checkFile: config.FilenameDecision,
			checkFor:  "Test decision from watch",
		},
		{
			name:      "learning update",
			update:    ContextUpdate{Type: config.UpdateTypeLearning, Content: "Test learning from watch"},
			checkFile: config.FilenameLearning,
			checkFor:  "Test learning from watch",
		},
		{
			name:      "convention update",
			update:    ContextUpdate{Type: config.UpdateTypeConvention, Content: "Test convention from watch"},
			checkFile: config.FilenameConvention,
			checkFor:  "Test convention from watch",
		},
		{
			name:        "unknown type",
			update:      ContextUpdate{Type: "invalid", Content: "Should fail"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := applyUpdate(tt.update)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("applyUpdate failed: %v", err)
			}

			// Verify content was added
			filePath := filepath.Join(config.DirContext, tt.checkFile)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tt.checkFile, err)
			}
			if !strings.Contains(string(content), tt.checkFor) {
				t.Errorf("expected %s to contain %q", tt.checkFile, tt.checkFor)
			}
		})
	}
}

// TestApplyCompleteUpdate tests the complete update type.
func TestApplyCompleteUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-complete-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add a task to complete
	tasksPath := filepath.Join(config.DirContext, config.FilenameTask)
	tasksContent := `# Tasks

## Next Up

- [ ] Implement authentication
- [ ] Write tests
`
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0644); err != nil {
		t.Fatalf("failed to write tasks: %v", err)
	}

	// Complete the task
	update := ContextUpdate{Type: config.UpdateTypeComplete, Content: "authentication"}
	if err := applyUpdate(update); err != nil {
		t.Fatalf("applyUpdate failed: %v", err)
	}

	// Verify task was marked complete
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read tasks: %v", err)
	}
	if !strings.Contains(string(content), "- [x] Implement authentication") {
		t.Error("task was not marked complete")
	}
	if !strings.Contains(string(content), "- [ ] Write tests") {
		t.Error("other task should remain unchecked")
	}
}

// TestBuildWatchSession tests session snapshot generation.
func TestBuildWatchSession(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 30, 0, 0, time.UTC)
	updates := []ContextUpdate{
		{Type: config.UpdateTypeTask, Content: "Task 1"},
		{Type: config.UpdateTypeTask, Content: "Task 2"},
		{Type: config.UpdateTypeDecision, Content: "Decision 1"},
		{Type: config.UpdateTypeLearning, Content: "Learning 1"},
	}

	result := buildWatchSession(timestamp, updates)

	// Check metadata
	if !strings.Contains(result, "**Date**: 2026-01-15") {
		t.Error("missing date in session")
	}
	if !strings.Contains(result, "**Time**: 14:30:00") {
		t.Error("missing time in session")
	}
	if !strings.Contains(result, "watch-auto-save") {
		t.Error("missing session type")
	}

	// Check updates by type
	if !strings.Contains(result, "### Tasks") {
		t.Error("missing Tasks section")
	}
	if !strings.Contains(result, "- Task 1") {
		t.Error("missing Task 1")
	}
	if !strings.Contains(result, "- Task 2") {
		t.Error("missing Task 2")
	}
	if !strings.Contains(result, "### Decisions") {
		t.Error("missing Decisions section")
	}
	if !strings.Contains(result, "### Learnings") {
		t.Error("missing Learnings section")
	}
}

// TestProcessStream tests stream processing applies updates.
func TestProcessStream(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-stream-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Ensure dry-run is off
	watchDryRun = false
	watchAutoSave = false

	input := `Some AI output text
<context-update type="task">Stream test task</context-update>
More output
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var output bytes.Buffer
	cmd.SetOut(&output)

	err = processStream(cmd, reader)
	if err != nil {
		t.Fatalf("processStream failed: %v", err)
	}

	// Verify task was written
	tasksPath := filepath.Join(config.DirContext, config.FilenameTask)
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read tasks: %v", err)
	}
	if !strings.Contains(string(content), "Stream test task") {
		t.Error("task should have been added to file")
	}
}

// TestRunCompleteSilentNoMatch tests complete with no matching task.
func TestRunCompleteSilentNoMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-nomatch-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Try to complete a non-existent task
	err = runCompleteSilent([]string{"nonexistent task query"})
	if err == nil {
		t.Error("expected error for non-matching task")
	}
	if !strings.Contains(err.Error(), "no task matching") {
		t.Errorf("unexpected error: %v", err)
	}
}
