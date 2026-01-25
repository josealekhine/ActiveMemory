//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestAgentCommand tests the agent command.
func TestAgentCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-agent-test-*")
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

	// Run agent - just verify it runs without error
	agentCmd := Cmd()
	agentCmd.SetArgs([]string{})

	if err := agentCmd.Execute(); err != nil {
		t.Fatalf("agent command failed: %v", err)
	}
}

// TestAgentJSONOutput tests the agent command with JSON output.
func TestAgentJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-agent-json-test-*")
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

	// Test agent with JSON output
	agentCmd := Cmd()
	agentCmd.SetArgs([]string{"--format", "json"})
	if err := agentCmd.Execute(); err != nil {
		t.Fatalf("agent --format json failed: %v", err)
	}
}
