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
	defer func() { _ = os.Chdir(origDir) }()

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
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test adding a decision with required flags
	t.Run("add decision", func(t *testing.T) {
		addCmd := Cmd()
		addCmd.SetArgs([]string{
			"decision", "Use PostgreSQL for database",
			"--context", "Need a reliable database",
			"--rationale", "PostgreSQL is well-supported",
			"--consequences", "Team needs training",
		})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add decision failed: %v", err)
		}

		content, err := os.ReadFile(".context/DECISIONS.md")
		if err != nil {
			t.Fatalf("failed to read DECISIONS.md: %v", err)
		}
		contentStr := string(content)
		if !strings.Contains(contentStr, "Use PostgreSQL for database") {
			t.Error("decision title was not added to DECISIONS.md")
		}
		if !strings.Contains(contentStr, "Need a reliable database") {
			t.Error("decision context was not added to DECISIONS.md")
		}
		if !strings.Contains(contentStr, "PostgreSQL is well-supported") {
			t.Error("decision rationale was not added to DECISIONS.md")
		}
		if !strings.Contains(contentStr, "Team needs training") {
			t.Error("decision consequences was not added to DECISIONS.md")
		}
	})

	// Test that decision without required flags fails
	t.Run("add decision without flags fails", func(t *testing.T) {
		addCmd := Cmd()
		addCmd.SetArgs([]string{"decision", "Incomplete decision"})
		err := addCmd.Execute()
		if err == nil {
			t.Fatal("expected error when adding decision without required flags")
		}
		if !strings.Contains(err.Error(), "--context") {
			t.Errorf("error should mention missing --context flag: %v", err)
		}
	})

	// Test adding a learning with required flags
	t.Run("add learning", func(t *testing.T) {
		addCmd := Cmd()
		addCmd.SetArgs([]string{
			"learning", "Always check for nil before dereferencing",
			"--context", "Got a nil pointer panic in production",
			"--lesson", "Always validate pointers before use",
			"--application", "Add nil checks in all pointer-receiving functions",
		})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add learning failed: %v", err)
		}

		content, err := os.ReadFile(".context/LEARNINGS.md")
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}
		contentStr := string(content)
		if !strings.Contains(contentStr, "Always check for nil before dereferencing") {
			t.Error("learning title was not added to LEARNINGS.md")
		}
		if !strings.Contains(contentStr, "Got a nil pointer panic in production") {
			t.Error("learning context was not added to LEARNINGS.md")
		}
		if !strings.Contains(contentStr, "Always validate pointers before use") {
			t.Error("learning lesson was not added to LEARNINGS.md")
		}
		if !strings.Contains(contentStr, "Add nil checks in all pointer-receiving functions") {
			t.Error("learning application was not added to LEARNINGS.md")
		}
	})

	// Test that learning without required flags fails
	t.Run("add learning without flags fails", func(t *testing.T) {
		addCmd := Cmd()
		addCmd.SetArgs([]string{"learning", "Incomplete learning"})
		err := addCmd.Execute()
		if err == nil {
			t.Fatal("expected error when adding learning without required flags")
		}
		if !strings.Contains(err.Error(), "--context") {
			t.Errorf("error should mention missing --context flag: %v", err)
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

// TestPrependOrder tests that decisions and learnings are prepended (newest first).
func TestPrependOrder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-prepend-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	t.Run("decisions are prepended", func(t *testing.T) {
		// Add first decision
		addCmd := Cmd()
		addCmd.SetArgs([]string{
			"decision", "First decision",
			"--context", "First context",
			"--rationale", "First rationale",
			"--consequences", "First consequences",
		})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add first decision failed: %v", err)
		}

		// Add second decision
		addCmd = Cmd()
		addCmd.SetArgs([]string{
			"decision", "Second decision",
			"--context", "Second context",
			"--rationale", "Second rationale",
			"--consequences", "Second consequences",
		})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add second decision failed: %v", err)
		}

		content, err := os.ReadFile(".context/DECISIONS.md")
		if err != nil {
			t.Fatalf("failed to read DECISIONS.md: %v", err)
		}

		contentStr := string(content)
		firstIdx := strings.Index(contentStr, "First decision")
		secondIdx := strings.Index(contentStr, "Second decision")

		if firstIdx == -1 || secondIdx == -1 {
			t.Fatal("decisions not found in file")
		}
		if secondIdx >= firstIdx {
			t.Errorf("second decision should appear before first (prepended), but first at %d, second at %d", firstIdx, secondIdx)
		}
	})

	t.Run("learnings are prepended", func(t *testing.T) {
		// Add first learning
		addCmd := Cmd()
		addCmd.SetArgs([]string{
			"learning", "First learning",
			"--context", "First context",
			"--lesson", "First lesson",
			"--application", "First application",
		})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add first learning failed: %v", err)
		}

		// Add second learning
		addCmd = Cmd()
		addCmd.SetArgs([]string{
			"learning", "Second learning",
			"--context", "Second context",
			"--lesson", "Second lesson",
			"--application", "Second application",
		})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add second learning failed: %v", err)
		}

		content, err := os.ReadFile(".context/LEARNINGS.md")
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}

		contentStr := string(content)
		firstIdx := strings.Index(contentStr, "First learning")
		secondIdx := strings.Index(contentStr, "Second learning")

		if firstIdx == -1 || secondIdx == -1 {
			t.Fatal("learnings not found in file")
		}
		if secondIdx >= firstIdx {
			t.Errorf("second learning should appear before first (prepended), but first at %d, second at %d", firstIdx, secondIdx)
		}
	})
}

// TestAppendEntry tests the AppendEntry function directly.
func TestAppendEntry(t *testing.T) {
	t.Run("decision prepends after header", func(t *testing.T) {
		// Use timestamp format "## [" to match what FormatDecision produces
		existing := []byte("# Decisions\n\n## [2026-01-01] Old Decision\n\nContent\n")
		entry := "## [2026-01-02] New Decision\n\nNew content\n"

		result := AppendEntry(existing, entry, "decision", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New Decision")
		oldIdx := strings.Index(resultStr, "Old Decision")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatalf("decisions not found in result: %s", resultStr)
		}
		if newIdx >= oldIdx {
			t.Errorf("new decision should appear before old, but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("learning prepends after separator", func(t *testing.T) {
		// Use timestamp format "- **[" to match what FormatLearning produces
		existing := []byte("# Learnings\n\n<!-- comment -->\n\n- **[2026-01-01]** Old Learning\n")
		entry := "- **[2026-01-02]** New Learning\n"

		result := AppendEntry(existing, entry, "learning", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New Learning")
		oldIdx := strings.Index(resultStr, "Old Learning")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatalf("learnings not found in result: %s", resultStr)
		}
		if newIdx >= oldIdx {
			t.Errorf("new learning should appear before old, but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("convention appends at end", func(t *testing.T) {
		existing := []byte("# Conventions\n\n- Old convention\n")
		entry := "- New convention\n"

		result := AppendEntry(existing, entry, "convention", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New convention")
		oldIdx := strings.Index(resultStr, "Old convention")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatal("conventions not found in result")
		}
		if newIdx <= oldIdx {
			t.Errorf("new convention should appear after old (appended), but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("decision on empty file", func(t *testing.T) {
		existing := []byte("# Decisions\n\n<!-- Add decisions here -->\n")
		entry := "## [2026-01-01] First Decision\n\nContent\n"

		result := AppendEntry(existing, entry, "decision", "")

		if !strings.Contains(string(result), "First Decision") {
			t.Errorf("decision not found in result: %s", result)
		}
	})

	t.Run("learning on empty file", func(t *testing.T) {
		existing := []byte("# Learnings\n\n<!-- Add gotchas here -->\n")
		entry := "- **[2026-01-01]** First Learning\n"

		result := AppendEntry(existing, entry, "learning", "")

		if !strings.Contains(string(result), "First Learning") {
			t.Errorf("learning not found in result: %s", result)
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
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Create a file with content (title)
	contentFile := filepath.Join(tmpDir, "learning-content.md")
	if err := os.WriteFile(contentFile, []byte("Content from file test"), 0644); err != nil {
		t.Fatalf("failed to create content file: %v", err)
	}

	// Test adding from file (still needs flags for learning)
	addCmd := Cmd()
	addCmd.SetArgs([]string{
		"learning", "--file", contentFile,
		"--context", "Testing file input",
		"--lesson", "File input works",
		"--application", "Use --file for long content",
	})
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
