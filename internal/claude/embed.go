// Package claude provides Claude Code integration templates and utilities.
package claude

import (
	"embed"
	"fmt"
)

//go:embed auto-save-session.sh
var FS embed.FS

// GetAutoSaveScript returns the auto-save session script.
func GetAutoSaveScript() ([]byte, error) {
	content, err := FS.ReadFile("auto-save-session.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to read auto-save-session.sh: %w", err)
	}
	return content, nil
}

// SettingsHooks represents the hooks section of settings.local.json
type SettingsHooks struct {
	PreToolUse []HookMatcher `json:"PreToolUse,omitempty"`
	SessionEnd []HookMatcher `json:"SessionEnd,omitempty"`
}

// HookMatcher represents a hook matcher with optional pattern
type HookMatcher struct {
	Matcher string `json:"matcher,omitempty"`
	Hooks   []Hook `json:"hooks"`
}

// Hook represents a single hook command
type Hook struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// Settings represents the full settings.local.json structure
type Settings struct {
	Hooks       SettingsHooks          `json:"hooks,omitempty"`
	Permissions map[string]interface{} `json:"permissions,omitempty"`
}

// CreateDefaultHooks returns the default ctx hooks configuration.
// Hooks use "ctx" expecting it to be in PATH.
func CreateDefaultHooks(projectDir string) SettingsHooks {
	hooksDir := ".claude/hooks"
	if projectDir != "" {
		hooksDir = fmt.Sprintf("%s/.claude/hooks", projectDir)
	}

	return SettingsHooks{
		PreToolUse: []HookMatcher{
			{
				Matcher: ".*",
				Hooks: []Hook{
					{
						Type:    "command",
						Command: "ctx agent --budget 4000 2>/dev/null || true",
					},
				},
			},
		},
		SessionEnd: []HookMatcher{
			{
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/auto-save-session.sh", hooksDir),
					},
				},
			},
		},
	}
}
