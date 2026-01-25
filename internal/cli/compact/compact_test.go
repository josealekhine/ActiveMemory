//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/complete"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestTruncateString tests the truncateString helper function.
func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a longer string", 10, "this is..."},
		{"", 10, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestRemoveEmptySections tests the removeEmptySections helper function.
func TestRemoveEmptySections(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		removed  int
	}{
		{
			name:     "no empty sections",
			input:    "# Title\n\n## Section\n\nContent here\n",
			expected: "# Title\n\n## Section\n\nContent here\n",
			removed:  0,
		},
		{
			name:     "single empty section",
			input:    "# Title\n\n## Empty\n\n## HasContent\n\nSome content\n",
			expected: "# Title\n\n## HasContent\n\nSome content\n",
			removed:  1,
		},
		{
			name:     "multiple empty sections",
			input:    "# Title\n\n## Empty1\n\n## Empty2\n\n## HasContent\n\nContent\n",
			expected: "# Title\n\n## HasContent\n\nContent\n",
			removed:  2,
		},
		{
			name:     "empty section at end",
			input:    "# Title\n\n## Content\n\nText\n\n## EmptyAtEnd\n",
			expected: "# Title\n\n## Content\n\nText\n",
			removed:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, count := removeEmptySections(tt.input)
			if count != tt.removed {
				t.Errorf("removeEmptySections() removed %d sections, want %d", count, tt.removed)
			}
			if result != tt.expected {
				t.Errorf("removeEmptySections() result:\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

// TestCompactCommand tests the compact command.
func TestCompactCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-compact-test-*")
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

	// Run compact with --no-auto-save to skip pre-compact session
	compactCmd := Cmd()
	compactCmd.SetArgs([]string{"--no-auto-save"})
	if err := compactCmd.Execute(); err != nil {
		t.Fatalf("compact failed: %v", err)
	}
}

// TestCompactWithTasks tests the compact command with actual completed tasks.
func TestCompactWithTasks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-compact-tasks-test-*")
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

	// Add and complete a task
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Task to complete"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task failed: %v", err)
	}

	completeCmd := complete.Cmd()
	completeCmd.SetArgs([]string{"Task to complete"})
	if err := completeCmd.Execute(); err != nil {
		t.Fatalf("complete task failed: %v", err)
	}

	// Run compact without auto-save
	compactCmd := Cmd()
	compactCmd.SetArgs([]string{"--no-auto-save"})
	if err := compactCmd.Execute(); err != nil {
		t.Fatalf("compact failed: %v", err)
	}
}
