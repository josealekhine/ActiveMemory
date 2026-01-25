//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx hook" command for generating AI tool integrations.
//
// The command outputs configuration snippets and instructions for integrating
// Context with various AI coding tools like Claude Code, Cursor, Aider, etc.
//
// Returns:
//   - *cobra.Command: Configured hook command that accepts a tool name argument
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook <tool>",
		Short: "Generate AI tool integration configs",
		Long: `Generate configuration and instructions 
for integrating Context with AI tools.

Supported tools:
  claude-code  - Anthropic's Claude Code CLI
  cursor       - Cursor IDE
  aider        - Aider AI coding assistant
  copilot      - GitHub Copilot
  windsurf     - Windsurf IDE

Example:
  ctx hook claude-code`,
		Args: cobra.ExactArgs(1),
		RunE: runHook,
	}

	return cmd
}
