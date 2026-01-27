//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"strings"
	"testing"
)

func TestGetAutoSaveScript(t *testing.T) {
	content, err := GetAutoSaveScript()
	if err != nil {
		t.Fatalf("GetAutoSaveScript() unexpected error: %v", err)
	}

	if len(content) == 0 {
		t.Error("GetAutoSaveScript() returned empty content")
	}

	// Check for expected script content
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("GetAutoSaveScript() script missing shebang")
	}
}

func TestGetBlockNonPathCtxScript(t *testing.T) {
	content, err := GetBlockNonPathCtxScript()
	if err != nil {
		t.Fatalf("GetBlockNonPathCtxScript() unexpected error: %v", err)
	}

	if len(content) == 0 {
		t.Error("GetBlockNonPathCtxScript() returned empty content")
	}

	// Check for expected script content
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("GetBlockNonPathCtxScript() script missing shebang")
	}
}

func TestListCommands(t *testing.T) {
	commands, err := ListCommands()
	if err != nil {
		t.Fatalf("ListCommands() unexpected error: %v", err)
	}

	if len(commands) == 0 {
		t.Error("ListCommands() returned empty list")
	}

	// Check that all entries are .md files
	for _, cmd := range commands {
		if !strings.HasSuffix(cmd, ".md") {
			t.Errorf("ListCommands() returned non-.md file: %s", cmd)
		}
	}
}

func TestGetCommand(t *testing.T) {
	// First get the list of commands to test with
	commands, err := ListCommands()
	if err != nil {
		t.Fatalf("ListCommands() failed: %v", err)
	}

	if len(commands) == 0 {
		t.Skip("no commands available to test")
	}

	// Test getting the first command
	content, err := GetCommand(commands[0])
	if err != nil {
		t.Errorf("GetCommand(%q) unexpected error: %v", commands[0], err)
	}
	if len(content) == 0 {
		t.Errorf("GetCommand(%q) returned empty content", commands[0])
	}

	// Test getting nonexistent command
	_, err = GetCommand("nonexistent-command.md")
	if err == nil {
		t.Error("GetCommand(nonexistent) expected error, got nil")
	}
}

func TestCreateDefaultHooks(t *testing.T) {
	tests := []struct {
		name       string
		projectDir string
	}{
		{
			name:       "empty project dir",
			projectDir: "",
		},
		{
			name:       "with project dir",
			projectDir: "/home/user/myproject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hooks := CreateDefaultHooks(tt.projectDir)

			// Check PreToolUse hooks
			if len(hooks.PreToolUse) == 0 {
				t.Error("CreateDefaultHooks() PreToolUse is empty")
			}

			// Check SessionEnd hooks
			if len(hooks.SessionEnd) == 0 {
				t.Error("CreateDefaultHooks() SessionEnd is empty")
			}

			// Check that project dir is used in paths when provided
			if tt.projectDir != "" {
				found := false
				for _, matcher := range hooks.PreToolUse {
					for _, hook := range matcher.Hooks {
						if strings.Contains(hook.Command, tt.projectDir) {
							found = true
							break
						}
					}
				}
				if !found {
					t.Error("CreateDefaultHooks() project dir not found in hook commands")
				}
			}
		})
	}
}

func TestSettingsStructure(t *testing.T) {
	// Test that Settings struct can be instantiated correctly
	settings := Settings{
		Hooks: CreateDefaultHooks(""),
		Permissions: PermissionsConfig{
			Allow: []string{"Bash(ctx status:*)", "Bash(ctx agent:*)"},
		},
	}

	if len(settings.Hooks.PreToolUse) == 0 {
		t.Error("Settings.Hooks.PreToolUse should not be empty")
	}

	if len(settings.Permissions.Allow) == 0 {
		t.Error("Settings.Permissions.Allow should not be empty")
	}
}

func TestCreateDefaultPermissions(t *testing.T) {
	perms := CreateDefaultPermissions()

	if len(perms) == 0 {
		t.Error("CreateDefaultPermissions should return permissions")
	}

	// Check that essential ctx commands are included
	expected := []string{
		"Bash(ctx status:*)",
		"Bash(ctx agent:*)",
		"Bash(ctx add:*)",
		"Bash(ctx session:*)",
	}

	permSet := make(map[string]bool)
	for _, p := range perms {
		permSet[p] = true
	}

	for _, e := range expected {
		if !permSet[e] {
			t.Errorf("Missing expected permission: %s", e)
		}
	}
}
