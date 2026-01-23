//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude provides Claude Code integration templates and utilities.
package claude

import (
	"embed"
	"fmt"
)

//go:embed tpl/auto-save-session.sh tpl/block-non-path-ctx.sh
var FS embed.FS

// GetAutoSaveScript returns the auto-save session script.
func GetAutoSaveScript() ([]byte, error) {
	content, err := FS.ReadFile("tpl/auto-save-session.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to read auto-save-session.sh: %w", err)
	}
	return content, nil
}

// GetBlockNonPathCtxScript returns the script that blocks non-PATH ctx invocations.
func GetBlockNonPathCtxScript() ([]byte, error) {
	content, err := FS.ReadFile("tpl/block-non-path-ctx.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to read block-non-path-ctx.sh: %w", err)
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
				// Block non-PATH ctx invocations (./ctx, ./dist/ctx, go run ./cmd/ctx)
				Matcher: "Bash",
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/block-non-path-ctx.sh", hooksDir),
					},
				},
			},
			{
				// Auto-load context on every tool use
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
