//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"encoding/json"
	"time"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// outputStatusJSON writes context status as JSON to the command output.
//
// When verbose is true, includes content previews for each file.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context to display
//   - verbose: If true, include file content previews
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func outputStatusJSON(
	cmd *cobra.Command, ctx *context.Context, verbose bool,
) error {
	output := Output{
		ContextDir:  ctx.Dir,
		TotalFiles:  len(ctx.Files),
		TotalTokens: ctx.TotalTokens,
		TotalSize:   ctx.TotalSize,
		Files:       make([]FileStatus, 0, len(ctx.Files)),
	}

	for _, f := range ctx.Files {
		fs := FileStatus{
			Name:    f.Name,
			Tokens:  f.Tokens,
			Size:    f.Size,
			IsEmpty: f.IsEmpty,
			Summary: f.Summary,
			ModTime: f.ModTime.Format(time.RFC3339),
		}
		if verbose && !f.IsEmpty {
			fs.Preview = getContentPreview(string(f.Content), 5)
		}
		output.Files = append(output.Files, fs)
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

// outputStatusText writes context status as formatted text to the command output.
//
// Displays a summary including file count, token estimate, file list with
// status indicators, and recent activity. When verbose is true, includes
// token counts, file sizes, and content previews for each file.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context to display
//   - verbose: If true, include detailed info and content previews
//
// Returns:
//   - error: Always nil (included for interface consistency)
func outputStatusText(
	cmd *cobra.Command, ctx *context.Context, verbose bool,
) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	cmd.Println(cyan("Context Status"))
	cmd.Println(cyan("===================="))
	cmd.Println()

	cmd.Printf("Context Directory: %s\n", ctx.Dir)
	cmd.Printf("Total Files: %d\n", len(ctx.Files))
	cmd.Printf("Token Estimate: %s tokens\n", formatNumber(ctx.TotalTokens))
	cmd.Println()

	cmd.Println("Files:")

	// Sort files in a logical order
	sortedFiles := make([]context.FileInfo, len(ctx.Files))
	copy(sortedFiles, ctx.Files)
	sortFilesByPriority(sortedFiles)

	for _, f := range sortedFiles {
		var status string
		var indicator string
		if f.IsEmpty {
			indicator = yellow("○")
			status = yellow("empty")
		} else {
			indicator = green("✓")
			status = f.Summary
		}

		if verbose {
			// Verbose: show tokens and size
			cmd.Printf("  %s %s (%s) [%s tokens, %s]\n",
				indicator, f.Name, status,
				formatNumber(f.Tokens), formatBytes(f.Size))

			// Show content preview for non-empty files
			if !f.IsEmpty {
				preview := getContentPreview(string(f.Content), 3)
				for _, line := range preview {
					cmd.Printf("      %s\n", dim(line))
				}
			}
		} else {
			cmd.Printf("  %s %s (%s)\n", indicator, f.Name, status)
		}
	}

	// Recent activity
	cmd.Println()
	cmd.Println("Recent Activity:")
	recentFiles := getRecentFiles(ctx.Files, 3)
	for _, f := range recentFiles {
		ago := formatTimeAgo(f.ModTime)
		cmd.Printf("  - %s modified %s\n", f.Name, ago)
	}

	return nil
}
