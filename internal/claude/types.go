//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

// HookConfig represents the hooks section of Claude Code's
// settings.local.json.
//
// Hooks are shell commands that Claude Code executes at specific lifecycle
// events. See https://docs.anthropic.com/en/docs/claude-code/hooks for details.
//
// Fields:
//   - PreToolUse: Matchers that run before each tool invocation
//   - SessionEnd: Matchers that run when a session ends
type HookConfig struct {
	PreToolUse  []HookMatcher `json:"PreToolUse,omitempty"`
	SessionEnd  []HookMatcher `json:"SessionEnd,omitempty"`
}

// HookMatcher associates a regex pattern with hooks to execute.
//
// For PreToolUse hooks, the Matcher pattern matches against the tool name
// (e.g., "Bash", "Read"). Use ".*" to match all tools.
//
// Fields:
//   - Matcher: Regex pattern to match; empty string matches all
//   - Hooks: Commands to execute when the pattern matches
type HookMatcher struct {
	Matcher string `json:"matcher,omitempty"`
	Hooks   []Hook `json:"hooks"`
}

// Hook represents a single hook command to execute.
//
// Fields:
//   - Type: Hook type, typically "command"
//   - Command: Shell command or script path to execute
type Hook struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// PermissionsConfig represents the permissions section of Claude Code's
// settings.local.json.
//
// Fields:
//   - Allow: List of tool patterns that are pre-approved (e.g., "Bash(ctx status:*)")
type PermissionsConfig struct {
	Allow []string `json:"allow,omitempty"`
}

// Settings represents the full Claude Code settings.local.json structure.
//
// This is used when reading or writing project-level Claude Code configuration.
//
// Fields:
//   - Hooks: Hook configuration for lifecycle events
//   - Permissions: Tool permission configuration
type Settings struct {
	Hooks       HookConfig        `json:"hooks,omitempty"`
	Permissions PermissionsConfig `json:"permissions,omitempty"`
}
