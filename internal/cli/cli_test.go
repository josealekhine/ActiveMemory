//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestBinaryIntegration is an integration test that builds and runs the actual binary.
//
// This test builds the ctx binary and exercises multiple commands to ensure
// they work correctly end-to-end. It verifies that subcommands execute properly
// (not falling through to root help) and produce expected output.
func TestBinaryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir, err := os.MkdirTemp("", "cli-binary-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Build the binary
	binaryPath := filepath.Join(tmpDir, "ctx-test-binary")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/ctx")
	buildCmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	// Get the project root (go up from internal/cli)
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to get project root: %v", err)
	}
	buildCmd.Dir = projectRoot

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, output)
	}

	// Create a test directory
	testDir := filepath.Join(tmpDir, "test-project")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	// Subtest: ctx init creates expected files
	t.Run("init creates expected files", func(t *testing.T) {
		initCmd := exec.Command(binaryPath, "init")
		initCmd.Dir = testDir
		if output, err := initCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx init failed: %v\n%s", err, output)
		}

		// Check .context directory exists
		ctxDir := filepath.Join(testDir, ".context")
		if _, err := os.Stat(ctxDir); os.IsNotExist(err) {
			t.Fatal(".context directory was not created")
		}

		// Check required files exist
		requiredFiles := []string{
			"CONSTITUTION.md",
			"TASKS.md",
			"DECISIONS.md",
			"LEARNINGS.md",
			"CONVENTIONS.md",
			"ARCHITECTURE.md",
		}
		for _, name := range requiredFiles {
			path := filepath.Join(ctxDir, name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("required file %s was not created", name)
			}
		}
	})

	// Subtest: ctx status returns valid status (not just help text)
	t.Run("status returns valid status", func(t *testing.T) {
		statusCmd := exec.Command(binaryPath, "status")
		statusCmd.Dir = testDir
		output, err := statusCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("ctx status failed: %v\n%s", err, output)
		}

		outputStr := string(output)
		// Verify it's actual status output, not help text
		if strings.Contains(outputStr, "Usage:") || strings.Contains(outputStr, "Available Commands:") {
			t.Error("ctx status returned help text instead of status")
		}
		// Check for expected status output markers
		if !strings.Contains(outputStr, "Context Status") && !strings.Contains(outputStr, "Context Directory") {
			t.Errorf("ctx status did not return expected status output, got:\n%s", outputStr)
		}
	})

	// Subtest: ctx add learning modifies LEARNINGS.md
	t.Run("add learning modifies LEARNINGS.md", func(t *testing.T) {
		addCmd := exec.Command(binaryPath, "add", "learning", "Test learning from integration test")
		addCmd.Dir = testDir
		if output, err := addCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx add learning failed: %v\n%s", err, output)
		}

		// Verify learning was added
		learningsPath := filepath.Join(testDir, ".context", "LEARNINGS.md")
		content, err := os.ReadFile(learningsPath)
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}
		if !strings.Contains(string(content), "Test learning from integration test") {
			t.Error("learning was not added to LEARNINGS.md")
		}
	})

	// Subtest: ctx session save creates session file
	t.Run("session save creates session file", func(t *testing.T) {
		saveCmd := exec.Command(binaryPath, "session", "save")
		saveCmd.Dir = testDir
		if output, err := saveCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx session save failed: %v\n%s", err, output)
		}

		// Check that sessions directory exists and has at least one file
		sessionsDir := filepath.Join(testDir, ".context", "sessions")
		entries, err := os.ReadDir(sessionsDir)
		if err != nil {
			t.Fatalf("failed to read sessions directory: %v", err)
		}
		if len(entries) == 0 {
			t.Error("no session file was created")
		}
	})

	// Subtest: ctx agent returns context packet
	t.Run("agent returns context packet", func(t *testing.T) {
		agentCmd := exec.Command(binaryPath, "agent")
		agentCmd.Dir = testDir
		output, err := agentCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("ctx agent failed: %v\n%s", err, output)
		}

		outputStr := string(output)
		// Verify it's actual agent output, not help text
		if strings.Contains(outputStr, "Usage:") || strings.Contains(outputStr, "Available Commands:") {
			t.Error("ctx agent returned help text instead of context packet")
		}
		// Check for expected context packet markers
		if !strings.Contains(outputStr, "CONSTITUTION") && !strings.Contains(outputStr, "TASKS") {
			t.Errorf("ctx agent did not return expected context packet, got:\n%s", outputStr)
		}
	})

	// Subtest: ctx drift runs without error
	t.Run("drift runs without error", func(t *testing.T) {
		driftCmd := exec.Command(binaryPath, "drift")
		driftCmd.Dir = testDir
		if output, err := driftCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx drift failed: %v\n%s", err, output)
		}
	})

	// Subtest: verify all subcommands execute (not falling through to root help)
	t.Run("subcommands execute without falling through to root help", func(t *testing.T) {
		// Commands that should produce output without "Available Commands:"
		// (which would indicate they fell through to root help)
		subcommands := []struct {
			args     []string
			checkFor string // expected output marker
		}{
			{[]string{"status"}, "Context"},
			{[]string{"agent"}, "Context Packet"},
			{[]string{"drift"}, "Drift"},
			{[]string{"load"}, ""},                 // load outputs context, varies by content
			{[]string{"hook", "cursor"}, "Cursor"}, // hook outputs integration instructions
		}

		for _, tc := range subcommands {
			t.Run(strings.Join(tc.args, "_"), func(t *testing.T) {
				cmd := exec.Command(binaryPath, tc.args...)
				cmd.Dir = testDir
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("ctx %s failed: %v\n%s", strings.Join(tc.args, " "), err, output)
				}

				outputStr := string(output)
				// Critical check: should NOT contain root help indicators
				if strings.Contains(outputStr, "Available Commands:") {
					t.Errorf("ctx %s fell through to root help:\n%s", strings.Join(tc.args, " "), outputStr)
				}
				// If we have an expected marker, check for it
				if tc.checkFor != "" && !strings.Contains(outputStr, tc.checkFor) {
					t.Errorf("ctx %s missing expected output %q:\n%s", strings.Join(tc.args, " "), tc.checkFor, outputStr)
				}
			})
		}
	})
}
