//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx init" command for initializing a .context/ directory.
//
// The command creates template files for maintaining persistent context
// for AI coding assistants. Files include constitution rules, tasks,
// decisions, learnings, conventions, and architecture documentation.
//
// Flags:
//   - --force, -f: Overwrite existing context files without prompting
//   - --minimal, -m: Only create essential files
//     (TASKS, DECISIONS, CONSTITUTION)
//   - --merge: Auto-merge ctx content into existing CLAUDE.md
//
// Returns:
//   - *cobra.Command: Configured init command with flags registered
func Cmd() *cobra.Command {
	var (
		force   bool
		minimal bool
		merge   bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new .context/ directory with template files",
		Long: `Initialize a new .context/ directory with template files for
maintaining persistent context for AI coding assistants.

The following files are created:
  - CONSTITUTION.md  — Hard invariants that must never be violated
  - TASKS.md         — Current and planned work
  - DECISIONS.md     — Architectural decisions with rationale
  - LEARNINGS.md     — Lessons learned, gotchas, tips
  - CONVENTIONS.md   — Project patterns and standards
  - ARCHITECTURE.md  — System overview
  - GLOSSARY.md      — Domain terms and abbreviations
  - DRIFT.md         — Staleness signals and update triggers
  - AGENT_PLAYBOOK.md — How AI agents should use this system

Use --minimal to only create essential files 
(TASKS.md, DECISIONS.md, CONSTITUTION.md).

Examples:
  ctx init           # Full initialization with all templates
  ctx init --minimal # Only essential files (TASKS, DECISIONS, CONSTITUTION)
  ctx init --force   # Overwrite existing files without prompting
  ctx init --merge   # Auto-merge ctx content into existing CLAUDE.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cmd, force, minimal, merge)
		},
	}

	cmd.Flags().BoolVarP(
		&force,
		"force", "f", false, "Overwrite existing context files",
	)
	cmd.Flags().BoolVarP(
		&minimal,
		"minimal", "m", false,
		"Only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md)",
	)
	cmd.Flags().BoolVar(
		&merge, "merge", false,
		"Auto-merge ctx content into existing CLAUDE.md without prompting",
	)

	return cmd
}
