package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/josealekhine/ActiveMemory/internal/context"
	"github.com/spf13/cobra"
)

var (
	statusJSON    bool
	statusVerbose bool
)

// StatusOutput represents the JSON output format for status command.
type StatusOutput struct {
	ContextDir  string       `json:"context_dir"`
	TotalFiles  int          `json:"total_files"`
	TotalTokens int          `json:"total_tokens"`
	TotalSize   int64        `json:"total_size"`
	Files       []FileStatus `json:"files"`
}

// FileStatus represents a single file's status in JSON output.
type FileStatus struct {
	Name    string `json:"name"`
	Tokens  int    `json:"tokens"`
	Size    int64  `json:"size"`
	IsEmpty bool   `json:"is_empty"`
	Summary string `json:"summary"`
	ModTime string `json:"mod_time"`
}

// StatusCmd returns the status command.
func StatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show context summary with token estimate",
		Long: `Display a summary of the current .context/ directory including:
  - Number of context files
  - Estimated token count
  - Status of each file
  - Recent activity`,
		RunE: runStatus,
	}

	cmd.Flags().BoolVar(&statusJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVarP(&statusVerbose, "verbose", "v", false, "Include file contents summary")

	return cmd
}

func runStatus(cmd *cobra.Command, args []string) error {
	ctx, err := context.Load("")
	if err != nil {
		if _, ok := err.(*context.NotFoundError); ok {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	if statusJSON {
		return outputStatusJSON(ctx)
	}

	return outputStatusText(ctx)
}

func outputStatusJSON(ctx *context.Context) error {
	output := StatusOutput{
		ContextDir:  ctx.Dir,
		TotalFiles:  len(ctx.Files),
		TotalTokens: ctx.TotalTokens,
		TotalSize:   ctx.TotalSize,
		Files:       make([]FileStatus, 0, len(ctx.Files)),
	}

	for _, f := range ctx.Files {
		output.Files = append(output.Files, FileStatus{
			Name:    f.Name,
			Tokens:  f.Tokens,
			Size:    f.Size,
			IsEmpty: f.IsEmpty,
			Summary: f.Summary,
			ModTime: f.ModTime.Format(time.RFC3339),
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func outputStatusText(ctx *context.Context) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(cyan("Active Memory Status"))
	fmt.Println(cyan("===================="))
	fmt.Println()

	fmt.Printf("Context Directory: %s\n", ctx.Dir)
	fmt.Printf("Total Files: %d\n", len(ctx.Files))
	fmt.Printf("Token Estimate: %s tokens\n", formatNumber(ctx.TotalTokens))
	fmt.Println()

	fmt.Println("Files:")

	// Sort files by a logical order
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
		fmt.Printf("  %s %s (%s)\n", indicator, f.Name, status)
	}

	// Recent activity
	fmt.Println()
	fmt.Println("Recent Activity:")
	recentFiles := getRecentFiles(ctx.Files, 3)
	for _, f := range recentFiles {
		ago := formatTimeAgo(f.ModTime)
		fmt.Printf("  - %s modified %s\n", f.Name, ago)
	}

	return nil
}

// sortFilesByPriority sorts files in the recommended read order.
func sortFilesByPriority(files []context.FileInfo) {
	priority := map[string]int{
		"CONSTITUTION.md":   1,
		"TASKS.md":          2,
		"CONVENTIONS.md":    3,
		"ARCHITECTURE.md":   4,
		"DECISIONS.md":      5,
		"LEARNINGS.md":      6,
		"GLOSSARY.md":       7,
		"DRIFT.md":          8,
		"AGENT_PLAYBOOK.md": 9,
	}

	sort.Slice(files, func(i, j int) bool {
		pi, ok := priority[files[i].Name]
		if !ok {
			pi = 100
		}
		pj, ok := priority[files[j].Name]
		if !ok {
			pj = 100
		}
		return pi < pj
	})
}

func getRecentFiles(files []context.FileInfo, n int) []context.FileInfo {
	sorted := make([]context.FileInfo, len(files))
	copy(sorted, files)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ModTime.After(sorted[j].ModTime)
	})
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	return sorted
}

func formatTimeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2, 2006")
	}
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}
