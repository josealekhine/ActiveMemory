//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// sessionsDirPath returns the path to the sessions directory.
func sessionsDirPath() string {
	return filepath.Join(config.DirContext, config.DirSessions)
}

// Cmd returns the session command with subcommands.
//
// Provides commands for managing session snapshots that capture
// the context state at a point in time, including tasks, decisions,
// and learnings.
//
// Returns:
//   - *cobra.Command: The session command with save, list, load,
//     and parse subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage session snapshots",
		Long: `Manage session snapshots in .context/sessions/.

Sessions capture the state of your context at a point in time,
including current tasks, recent decisions, and learnings.

Subcommands:
  save    Save current context state to a session file
  list    List saved sessions with summaries
  load    Load and display a previous session
  parse   Convert .jsonl transcript to readable markdown`,
	}

	cmd.AddCommand(sessionSaveCmd())
	cmd.AddCommand(sessionListCmd())
	cmd.AddCommand(sessionLoadCmd())
	cmd.AddCommand(sessionParseCmd())

	return cmd
}

// sessionSaveCmd returns the session save subcommand.
//
// Returns:
//   - *cobra.Command: Command for saving the current context state to
//     a session file
func sessionSaveCmd() *cobra.Command {
	var sessionType string

	cmd := &cobra.Command{
		Use:   "save [topic]",
		Short: "Save current context state to a session file",
		Long: `Save a snapshot of the current context state to .context/sessions/.

The session file includes:
  - Summary of what was done
  - Current tasks from TASKS.md
  - Recent decisions from DECISIONS.md
  - Recent learnings from LEARNINGS.md

Examples:
  ctx session save "implemented auth"
  ctx session save "refactored API" --type feature
  ctx session save  # prompts for topic interactively`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSessionSave(cmd, args, sessionType)
		},
	}

	cmd.Flags().StringVarP(
		&sessionType,
		"type", "t",
		"session", "Session type (feature, bugfix, refactor, session)",
	)

	return cmd
}

// sessionListCmd returns the session list subcommand.
//
// Returns:
//   - *cobra.Command: Command for listing saved sessions with summaries
func sessionListCmd() *cobra.Command {
	var listLimit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List saved sessions with summaries",
		Long: `List all saved sessions in .context/sessions/.

Shows session date, topic, type, and a brief summary for each session.
Sessions are sorted by date (newest first).

Examples:
  ctx session list
  ctx session list --limit 5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSessionList(cmd, listLimit)
		},
	}

	cmd.Flags().IntVarP(
		&listLimit, "limit", "n", 10, "Maximum number of sessions to display",
	)

	return cmd
}

// sessionLoadCmd returns the session load subcommand.
//
// Returns:
//   - *cobra.Command: Command for loading and displaying a saved session
func sessionLoadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load <file>",
		Short: "Load and display a previous session",
		Long: `Load and display the contents of a saved session.

The file argument can be:
  - A full filename (e.g., 2025-01-21-004900-ctx-rename.md)
  - A partial match (e.g., "ctx-rename" or "2025-01-21")
  - A number from 'ctx session list' output (1 = most recent)

Examples:
  ctx session load 2025-01-21-004900-ctx-rename.md
  ctx session load ctx-rename
  ctx session load 1`,
		Args: cobra.ExactArgs(1),
		RunE: runSessionLoad,
	}

	return cmd
}

// sessionParseCmd returns the session parse subcommand.
//
// Returns:
//   - *cobra.Command: Command for converting JSONL transcripts to Markdown
func sessionParseCmd() *cobra.Command {
	var (
		output  string
		extract bool
	)

	cmd := &cobra.Command{
		Use:   "parse <file.jsonl>",
		Short: "Convert .jsonl transcript to readable markdown",
		Long: `Parse a Claude Code .jsonl transcript file 
and convert it to readable markdown.

The .jsonl files are auto-saved by the SessionEndHooks hook and contain the full
conversation transcript including tool calls and results.

Examples:
  ctx session parse .context/sessions/2026-01-21-072504-session.jsonl
  ctx session parse .context/sessions/2026-01-21-072504-session.jsonl -o conversation.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSessionParse(cmd, args, output, extract)
		},
	}

	cmd.Flags().StringVarP(
		&output,
		"output", "o", "", "Output file (default: stdout)",
	)
	cmd.Flags().BoolVar(
		&extract,
		"extract", false,
		"Extract potential decisions and learnings from transcript",
	)

	return cmd
}
