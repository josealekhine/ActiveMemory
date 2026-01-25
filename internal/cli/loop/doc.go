//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package loop provides the command for generating Ralph loop scripts.
//
// A Ralph loop is an iterative development technique where an AI assistant
// runs repeatedly with the same prompt until a completion signal is detected.
// This enables autonomous development where the AI builds on its previous work
// across multiple iterations.
//
// # How It Works
//
// The generated script:
//
//  1. Reads the prompt file (default: PROMPT.md)
//  2. Runs the AI tool with the prompt
//  3. Checks output for a completion signal
//  4. Repeats until signal is detected or max iterations reached
//
// # Supported Tools
//
// The loop command generates scripts for different AI tools:
//
//   - claude: Claude Code CLI (default)
//   - aider: Aider AI pair programming tool
//   - generic: Template for custom tools
//
// # Completion Signal
//
// The completion signal (default: "SYSTEM_CONVERGED") indicates the AI has
// finished its work. The AI should output this signal when it determines
// that the task is complete. The loop script watches for this signal and
// exits when detected.
//
// # File Organization
//
//   - loop.go: Command definition and flag handling
//   - run.go: Main loop script generation logic
//   - script.go: Shell script templates for each tool
package loop
