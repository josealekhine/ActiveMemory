//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap initializes the ctx CLI application.
//
// It provides functions to create the root command and register all
// subcommands. The typical usage pattern is:
//
//	cmd := bootstrap.Initialize(bootstrap.RootCmd())
//	if err := cmd.Execute(); err != nil {
//	    // handle error
//	}
package bootstrap

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/agent"
	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/cli/complete"
	"github.com/ActiveMemory/ctx/internal/cli/drift"
	"github.com/ActiveMemory/ctx/internal/cli/hook"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/cli/load"
	"github.com/ActiveMemory/ctx/internal/cli/loop"
	"github.com/ActiveMemory/ctx/internal/cli/session"
	"github.com/ActiveMemory/ctx/internal/cli/status"
	"github.com/ActiveMemory/ctx/internal/cli/sync"
	"github.com/ActiveMemory/ctx/internal/cli/task"
	"github.com/ActiveMemory/ctx/internal/cli/watch"
)

// version is set at build time via ldflags
const version = "dev"

// RootCmd creates and returns the root cobra command for the ctx CLI.
//
// The root command provides the entry point for all ctx subcommands and
// displays help information when invoked without arguments.
//
// Returns:
//   - *cobra.Command: The configured root command with usage and version info
func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ctx",
		Short: "Context - persistent context for AI coding assistants",
		Long: `Context (ctx) maintains persistent context files that help
  AI coding assistants understand your project's architecture, conventions,
  decisions, and current tasks.
  
  Use 'ctx init' to create a .context/ directory in your project,
  then use 'ctx status', 'ctx load', and 'ctx agent' to work with context.`,
		Version: version,
	}
}

// Initialize registers all ctx subcommands with the root command.
//
// This function attaches all available subcommands (init, status, load, add,
// complete, agent, drift, sync, compact, watch, hook, session, tasks, loop)
// to the provided root command.
//
// Parameters:
//   - cmd: The root cobra command to attach subcommands to
//
// Returns:
//   - *cobra.Command: The same command with all subcommands registered
func Initialize(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(initialize.Cmd())
	cmd.AddCommand(status.Cmd())
	cmd.AddCommand(load.Cmd())
	cmd.AddCommand(add.Cmd())
	cmd.AddCommand(complete.Cmd())
	cmd.AddCommand(agent.Cmd())
	cmd.AddCommand(drift.Cmd())
	cmd.AddCommand(sync.Cmd())
	cmd.AddCommand(compact.Cmd())
	cmd.AddCommand(watch.Cmd())
	cmd.AddCommand(hook.Cmd())
	cmd.AddCommand(session.Cmd())
	cmd.AddCommand(task.Cmd())
	cmd.AddCommand(loop.Cmd())

	return cmd
}
