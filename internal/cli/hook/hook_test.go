//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"testing"
)

// TestHookCommand tests the hook command.
func TestHookCommand(t *testing.T) {
	tests := []struct {
		tool     string
		contains string
	}{
		{"claude-code", "Claude Code Integration"},
		{"cursor", "Cursor IDE Integration"},
		{"aider", "Aider Integration"},
		{"copilot", "GitHub Copilot Integration"},
		{"windsurf", "Windsurf Integration"},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			hookCmd := Cmd()
			hookCmd.SetArgs([]string{tt.tool})

			if err := hookCmd.Execute(); err != nil {
				t.Fatalf("hook %s command failed: %v", tt.tool, err)
			}
		})
	}
}

// TestHookCommandUnknownTool tests hook command with unknown tool.
func TestHookCommandUnknownTool(t *testing.T) {
	hookCmd := Cmd()
	hookCmd.SetArgs([]string{"unknown-tool"})

	err := hookCmd.Execute()
	if err == nil {
		t.Error("hook command should fail for unknown tool")
	}
}
