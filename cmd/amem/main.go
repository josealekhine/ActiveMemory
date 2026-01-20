package main

import (
	"fmt"
	"os"

	"github.com/josealekhine/ActiveMemory/internal/cli"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "amem",
	Short: "Active Memory - persistent context for AI coding assistants",
	Long: `Active Memory (amem) maintains persistent context files that help
AI coding assistants understand your project's architecture, conventions,
decisions, and current tasks.

Use 'amem init' to create a .context/ directory in your project,
then use 'amem status', 'amem load', and 'amem agent' to work with context.`,
	Version: Version,
}

func init() {
	rootCmd.AddCommand(cli.InitCmd())
	rootCmd.AddCommand(cli.StatusCmd())
	rootCmd.AddCommand(cli.LoadCmd())
	rootCmd.AddCommand(cli.AddCmd())
	rootCmd.AddCommand(cli.CompleteCmd())
	rootCmd.AddCommand(cli.AgentCmd())
	rootCmd.AddCommand(cli.DriftCmd())
	rootCmd.AddCommand(cli.SyncCmd())
	rootCmd.AddCommand(cli.CompactCmd())
	rootCmd.AddCommand(cli.WatchCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
