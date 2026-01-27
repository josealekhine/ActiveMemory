//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/templates"
)

// GetAutoSaveScript returns the auto-save session script content.
//
// The script automatically saves Claude Code session transcripts when a
// session ends. It is installed to .claude/hooks/ during ctx init --claude.
//
// Returns:
//   - []byte: Raw bytes of the auto-save-session.sh script
//   - error: Non-nil if the embedded file cannot be read
func GetAutoSaveScript() ([]byte, error) {
	content, err := templates.GetClaudeHook("auto-save-session.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to read auto-save-session.sh: %w", err)
	}
	return content, nil
}

// GetBlockNonPathCtxScript returns the script that blocks non-PATH ctx
// invocations.
//
// The script prevents Claude from running ctx via relative paths (./ctx,
// ./dist/ctx) or "go run", ensuring only the installed PATH version is used.
// It is installed to .claude/hooks/ during ctx init --claude.
//
// Returns:
//   - []byte: Raw bytes of the block-non-path-ctx.sh script
//   - error: Non-nil if the embedded file cannot be read
func GetBlockNonPathCtxScript() ([]byte, error) {
	content, err := templates.GetClaudeHook("block-non-path-ctx.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to read block-non-path-ctx.sh: %w", err)
	}
	return content, nil
}

// ListCommands returns the list of embedded command file names.
//
// These are Claude Code slash command definitions (e.g., "ctx-status.md",
// "ctx-reflect.md") from internal/templates/claude/commands/. They can be
// installed to .claude/commands/ via "ctx init".
//
// Returns:
//   - []string: Filenames of available command definitions
//   - error: Non-nil if the commands directory cannot be read
func ListCommands() ([]string, error) {
	names, err := templates.ListClaudeCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to list commands: %w", err)
	}
	return names, nil
}

// GetCommand returns the content of a command file by name.
//
// Parameters:
//   - name: Filename as returned by [ListCommands] (e.g., "ctx-status.md")
//
// Returns:
//   - []byte: Raw bytes of the command definition file
//   - error: Non-nil if the command file does not exist or cannot be read
func GetCommand(name string) ([]byte, error) {
	content, err := templates.GetClaudeCommand(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read command %s: %w", name, err)
	}
	return content, nil
}

// CreateDefaultPermissions returns the default permissions for ctx commands.
//
// These permissions allow Claude Code to run ctx CLI commands without
// prompting for approval. All ctx subcommands are pre-approved.
//
// Returns:
//   - []string: List of permission patterns for ctx commands
func CreateDefaultPermissions() []string {
	return []string{
		"Bash(ctx status:*)",
		"Bash(ctx agent:*)",
		"Bash(ctx add:*)",
		"Bash(ctx session:*)",
		"Bash(ctx tasks:*)",
		"Bash(ctx loop:*)",
	}
}

// CreateDefaultHooks returns the default ctx hooks configuration for
// Claude Code.
//
// The returned hooks configure PreToolUse to block non-PATH ctx
// invocations and auto-load context on every tool use, and SessionEnd
// to run auto-save-session.sh for persisting session transcripts.
//
// Parameters:
//   - projectDir: Project root directory for absolute hook paths; if empty,
//     paths are relative (e.g., ".claude/hooks/")
//
// Returns:
//   - HookConfig: Configured hooks for PreToolUse and SessionEnd events
func CreateDefaultHooks(projectDir string) HookConfig {
	hooksDir := ".claude/hooks"
	if projectDir != "" {
		hooksDir = fmt.Sprintf("%s/.claude/hooks", projectDir)
	}

	return HookConfig{
		PreToolUse: []HookMatcher{
			{
				// Block non-PATH ctx invocations and require approval for git push
				Matcher: "Bash",
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/block-non-path-ctx.sh", hooksDir),
					},
					{
						Type:    "command",
						Command: `if echo "$CLAUDE_TOOL_INPUT" | grep -qE 'git\s+push'; then echo '{"decision": "block", "reason": "git push not allowed - ask user first"}'; exit 0; fi`,
					},
				},
			},
			{
				// Autoload context on every tool use
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
