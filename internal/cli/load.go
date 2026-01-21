package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/josealekhine/ActiveMemory/internal/context"
	"github.com/spf13/cobra"
)

var (
	loadBudget int
	loadRaw    bool
)

// LoadCmd returns the load command.
func LoadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Output assembled context markdown",
		Long: `Load and display the assembled context as it would be provided to an AI.

The context files are assembled in the recommended read order:
  1. CONSTITUTION.md
  2. TASKS.md
  3. CONVENTIONS.md
  4. ARCHITECTURE.md
  5. DECISIONS.md
  6. LEARNINGS.md
  7. GLOSSARY.md
  8. DRIFT.md
  9. AGENT_PLAYBOOK.md

Use --raw to output raw file contents without headers or assembly.
Use --budget to limit output to a specific token count.`,
		RunE: runLoad,
	}

	cmd.Flags().IntVar(&loadBudget, "budget", 8000, "Token budget for assembly")
	cmd.Flags().BoolVar(&loadRaw, "raw", false, "Output raw file contents without assembly")

	return cmd
}

// fileReadOrder defines the priority order for reading context files.
var fileReadOrder = []string{
	"CONSTITUTION.md",
	"TASKS.md",
	"CONVENTIONS.md",
	"ARCHITECTURE.md",
	"DECISIONS.md",
	"LEARNINGS.md",
	"GLOSSARY.md",
	"DRIFT.md",
	"AGENT_PLAYBOOK.md",
}

func runLoad(cmd *cobra.Command, args []string) error {
	ctx, err := context.Load("")
	if err != nil {
		if _, ok := err.(*context.NotFoundError); ok {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	if loadRaw {
		return outputRaw(ctx)
	}

	return outputAssembled(ctx, loadBudget)
}

func outputRaw(ctx *context.Context) error {
	// Sort files by read order
	files := sortByReadOrder(ctx.Files)

	for i, f := range files {
		if i > 0 {
			fmt.Println()
		}
		fmt.Print(string(f.Content))
	}
	return nil
}

func outputAssembled(ctx *context.Context, budget int) error {
	var sb strings.Builder

	// Header
	sb.WriteString("# Active Memory Context\n\n")
	sb.WriteString(fmt.Sprintf("Token Budget: %d | Available: %d\n\n", budget, ctx.TotalTokens))
	sb.WriteString("---\n\n")

	// Sort files by read order
	files := sortByReadOrder(ctx.Files)

	tokensUsed := context.EstimateTokensString(sb.String())

	for _, f := range files {
		// Skip empty files
		if f.IsEmpty {
			continue
		}

		// Check if we have budget for this file
		fileTokens := f.Tokens
		if tokensUsed+fileTokens > budget {
			// Add truncation notice
			sb.WriteString(fmt.Sprintf("\n---\n\n*[Truncated: %s and remaining files excluded due to token budget]*\n", f.Name))
			break
		}

		// Add file section
		sb.WriteString(fmt.Sprintf("## %s\n\n", fileNameToTitle(f.Name)))
		sb.Write(f.Content)
		if !strings.HasSuffix(string(f.Content), "\n") {
			sb.WriteString("\n")
		}
		sb.WriteString("\n---\n\n")

		tokensUsed += fileTokens
	}

	fmt.Print(sb.String())
	return nil
}

func sortByReadOrder(files []context.FileInfo) []context.FileInfo {
	// Create a map for quick priority lookup
	priority := make(map[string]int)
	for i, name := range fileReadOrder {
		priority[name] = i
	}

	// Copy and sort
	sorted := make([]context.FileInfo, len(files))
	copy(sorted, files)

	sort.Slice(sorted, func(i, j int) bool {
		pi, ok := priority[sorted[i].Name]
		if !ok {
			pi = 100
		}
		pj, ok := priority[sorted[j].Name]
		if !ok {
			pj = 100
		}
		return pi < pj
	})

	return sorted
}

func fileNameToTitle(name string) string {
	// Remove .md extension
	name = strings.TrimSuffix(name, ".md")
	// Convert SCREAMING_SNAKE to Title Case
	name = strings.ReplaceAll(name, "_", " ")
	// Title case each word
	words := strings.Fields(name)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, " ")
}
