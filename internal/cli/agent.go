package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/josealekhine/ActiveMemory/internal/context"
	"github.com/spf13/cobra"
)

var (
	agentBudget int
	agentFormat string
)

// AgentCmd returns the agent command.
func AgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Print AI-ready context packet",
		Long: `Print a concise context packet optimized for AI consumption.

The output is designed to be copy-pasted into an AI chat or piped to a system prompt.
It includes:
  - Constitution rules (NEVER VIOLATE)
  - Current tasks
  - Key conventions
  - Recent decisions

Use --budget to limit token output (default 8000).
Use --format to choose between markdown (md) or JSON output.`,
		RunE: runAgent,
	}

	cmd.Flags().IntVar(&agentBudget, "budget", 8000, "Token budget for context packet")
	cmd.Flags().StringVar(&agentFormat, "format", "md", "Output format: md or json")

	return cmd
}

// AgentPacket represents the JSON output format for agent command.
type AgentPacket struct {
	Generated    string   `json:"generated"`
	Budget       int      `json:"budget"`
	TokensUsed   int      `json:"tokens_used"`
	ReadOrder    []string `json:"read_order"`
	Constitution []string `json:"constitution"`
	Tasks        []string `json:"tasks"`
	Conventions  []string `json:"conventions"`
	Decisions    []string `json:"decisions"`
}

func runAgent(cmd *cobra.Command, args []string) error {
	ctx, err := context.Load("")
	if err != nil {
		if _, ok := err.(*context.NotFoundError); ok {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	if agentFormat == "json" {
		return outputAgentJSON(ctx)
	}

	return outputAgentMarkdown(ctx)
}

func outputAgentJSON(ctx *context.Context) error {
	packet := AgentPacket{
		Generated:    time.Now().UTC().Format(time.RFC3339),
		Budget:       agentBudget,
		TokensUsed:   ctx.TotalTokens,
		ReadOrder:    getReadOrder(ctx),
		Constitution: extractConstitutionRules(ctx),
		Tasks:        extractActiveTasks(ctx),
		Conventions:  extractConventions(ctx),
		Decisions:    extractRecentDecisions(ctx, 3),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(packet)
}

func outputAgentMarkdown(ctx *context.Context) error {
	var sb strings.Builder

	timestamp := time.Now().UTC().Format(time.RFC3339)
	sb.WriteString("# Active Memory Context Packet\n")
	sb.WriteString(fmt.Sprintf("Generated: %s | Budget: %d tokens | Used: %d\n\n", timestamp, agentBudget, ctx.TotalTokens))

	// Read order
	sb.WriteString("## Read These Files (in order)\n")
	for i, path := range getReadOrder(ctx) {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, path))
	}
	sb.WriteString("\n")

	// Constitution
	rules := extractConstitutionRules(ctx)
	if len(rules) > 0 {
		sb.WriteString("## Constitution (NEVER VIOLATE)\n")
		for _, rule := range rules {
			sb.WriteString(fmt.Sprintf("- %s\n", rule))
		}
		sb.WriteString("\n")
	}

	// Current tasks
	tasks := extractActiveTasks(ctx)
	if len(tasks) > 0 {
		sb.WriteString("## Current Tasks\n")
		for _, task := range tasks {
			sb.WriteString(fmt.Sprintf("%s\n", task))
		}
		sb.WriteString("\n")
	}

	// Conventions
	conventions := extractConventions(ctx)
	if len(conventions) > 0 {
		sb.WriteString("## Key Conventions\n")
		for _, conv := range conventions {
			sb.WriteString(fmt.Sprintf("- %s\n", conv))
		}
		sb.WriteString("\n")
	}

	// Recent decisions
	decisions := extractRecentDecisions(ctx, 3)
	if len(decisions) > 0 {
		sb.WriteString("## Recent Decisions\n")
		for _, dec := range decisions {
			sb.WriteString(fmt.Sprintf("- %s\n", dec))
		}
		sb.WriteString("\n")
	}

	fmt.Print(sb.String())
	return nil
}

func getReadOrder(ctx *context.Context) []string {
	order := []string{}
	for _, name := range fileReadOrder {
		for _, f := range ctx.Files {
			if f.Name == name && !f.IsEmpty {
				order = append(order, fmt.Sprintf("%s/%s", ctx.Dir, f.Name))
				break
			}
		}
	}
	return order
}

func extractConstitutionRules(ctx *context.Context) []string {
	for _, f := range ctx.Files {
		if f.Name == "CONSTITUTION.md" {
			return extractCheckboxItems(string(f.Content))
		}
	}
	return nil
}

func extractActiveTasks(ctx *context.Context) []string {
	for _, f := range ctx.Files {
		if f.Name == "TASKS.md" {
			return extractUncheckedTasks(string(f.Content))
		}
	}
	return nil
}

func extractConventions(ctx *context.Context) []string {
	for _, f := range ctx.Files {
		if f.Name == "CONVENTIONS.md" {
			return extractBulletItems(string(f.Content), 5)
		}
	}
	return nil
}

func extractRecentDecisions(ctx *context.Context, limit int) []string {
	for _, f := range ctx.Files {
		if f.Name == "DECISIONS.md" {
			return extractDecisionTitles(string(f.Content), limit)
		}
	}
	return nil
}

func extractCheckboxItems(content string) []string {
	re := regexp.MustCompile(`(?m)^-\s*\[[ x]\]\s*(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		items = append(items, strings.TrimSpace(m[1]))
	}
	return items
}

func extractUncheckedTasks(content string) []string {
	re := regexp.MustCompile(`(?m)^-\s*\[\s*\]\s*(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		items = append(items, "- [ ] "+strings.TrimSpace(m[1]))
	}
	return items
}

func extractBulletItems(content string, limit int) []string {
	re := regexp.MustCompile(`(?m)^-\s+(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, limit)
	for i, m := range matches {
		if i >= limit {
			break
		}
		text := strings.TrimSpace(m[1])
		// Skip empty or header-only items
		if text != "" && !strings.HasPrefix(text, "#") {
			items = append(items, text)
		}
	}
	return items
}

func extractDecisionTitles(content string, limit int) []string {
	re := regexp.MustCompile(`(?m)^##\s+\[[\d-]+\]\s*(.+)$`)
	matches := re.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, limit)
	// Get the most recent (last) decisions
	start := len(matches) - limit
	if start < 0 {
		start = 0
	}
	for i := start; i < len(matches); i++ {
		items = append(items, strings.TrimSpace(matches[i][1]))
	}
	return items
}
