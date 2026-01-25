//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInitCommand tests the init command creates the .context directory.
func TestInitCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save and restore working directory
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Run the init command
	cmd := Cmd()
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Check that .context directory was created
	ctxDir := filepath.Join(tmpDir, ".context")
	info, err := os.Stat(ctxDir)
	if err != nil {
		t.Fatalf(".context directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal(".context should be a directory")
	}

	// Check that required files exist
	requiredFiles := []string{
		"CONSTITUTION.md",
		"TASKS.md",
		"DECISIONS.md",
		"CONVENTIONS.md",
		"ARCHITECTURE.md",
	}

	for _, name := range requiredFiles {
		path := filepath.Join(ctxDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("required file %s was not created", name)
		}
	}
}

// TestInitWithExistingClaudeMdWithCtxMarker tests init when CLAUDE.md
// already exists with ctx marker.
func TestInitWithExistingClaudeMdWithCtxMarker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-existing-claude-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create existing CLAUDE.md with ctx marker already present
	existingContent := `# My Project

This is my existing CLAUDE.md content.

<!-- ctx:context -->
Old ctx content here
<!-- ctx:end -->

## Custom Section

Some custom content here.
`
	if err := os.WriteFile("CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	// Run init
	initCmd := Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Check that CLAUDE.md was updated
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	// Should still contain ctx marker (updated)
	if !strings.Contains(string(content), "ctx:context") {
		t.Error("CLAUDE.md missing ctx:context marker")
	}

	// Should preserve custom section
	if !strings.Contains(string(content), "Custom Section") {
		t.Error("CLAUDE.md lost custom section")
	}
}
