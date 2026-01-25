//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestAddCommand tests the add command.
func TestAddCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-test-*")
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

	// Test adding a task
	addCmd := Cmd()
	addCmd.SetArgs([]string{"task", "Test task for integration"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task command failed: %v", err)
	}

	// Verify the task was added
	tasksPath := filepath.Join(tmpDir, ".context", "TASKS.md")
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read TASKS.md: %v", err)
	}

	if !strings.Contains(string(content), "Test task for integration") {
		t.Errorf("task was not added to TASKS.md")
	}
}

// TestAddDecisionAndLearning tests adding decisions and learnings.
func TestAddDecisionAndLearning(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-dl-test-*")
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

	// Test adding a decision
	t.Run("add decision", func(t *testing.T) {
		addCmd := Cmd()
		addCmd.SetArgs([]string{"decision", "Use PostgreSQL for database"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add decision failed: %v", err)
		}

		content, err := os.ReadFile(".context/DECISIONS.md")
		if err != nil {
			t.Fatalf("failed to read DECISIONS.md: %v", err)
		}
		if !strings.Contains(string(content), "Use PostgreSQL for database") {
			t.Error("decision was not added to DECISIONS.md")
		}
	})

	// Test adding a learning
	t.Run("add learning", func(t *testing.T) {
		addCmd := Cmd()
		addCmd.SetArgs([]string{"learning", "Always check for nil before dereferencing"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add learning failed: %v", err)
		}

		content, err := os.ReadFile(".context/LEARNINGS.md")
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}
		if !strings.Contains(string(content), "Always check for nil before dereferencing") {
			t.Error("learning was not added to LEARNINGS.md")
		}
	})

	// Test adding a convention
	t.Run("add convention", func(t *testing.T) {
		addCmd := Cmd()
		addCmd.SetArgs([]string{"convention", "Use camelCase for variable names"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add convention failed: %v", err)
		}

		content, err := os.ReadFile(".context/CONVENTIONS.md")
		if err != nil {
			t.Fatalf("failed to read CONVENTIONS.md: %v", err)
		}
		if !strings.Contains(string(content), "Use camelCase for variable names") {
			t.Error("convention was not added to CONVENTIONS.md")
		}
	})
}

// TestAddFromFile tests adding content from a file.
func TestAddFromFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-file-test-*")
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

	// Create a file with content
	contentFile := filepath.Join(tmpDir, "learning-content.md")
	if err := os.WriteFile(contentFile, []byte("Content from file test"), 0644); err != nil {
		t.Fatalf("failed to create content file: %v", err)
	}

	// Test adding from file
	addCmd := Cmd()
	addCmd.SetArgs([]string{"learning", "--file", contentFile})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add from file failed: %v", err)
	}

	content, err := os.ReadFile(".context/LEARNINGS.md")
	if err != nil {
		t.Fatalf("failed to read LEARNINGS.md: %v", err)
	}
	if !strings.Contains(string(content), "Content from file test") {
		t.Error("content from file was not added to LEARNINGS.md")
	}
}
