package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// HookCmd returns the hook command.
func HookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook <tool>",
		Short: "Generate AI tool integration configs",
		Long: `Generate configuration and instructions for integrating Active Memory with AI tools.

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

func runHook(cmd *cobra.Command, args []string) error {
	tool := strings.ToLower(args[0])

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	switch tool {
	case "claude-code", "claude":
		fmt.Println(cyan("Claude Code Integration"))
		fmt.Println(cyan("======================="))
		fmt.Println()
		fmt.Println("Add this to your project's CLAUDE.md or system prompt:")
		fmt.Println()
		fmt.Println(green("```markdown"))
		fmt.Print(`## Active Memory Context

Before starting any task, load the project context:

1. Read .context/CONSTITUTION.md — These rules are INVIOLABLE
2. Read .context/TASKS.md — Current work items
3. Read .context/CONVENTIONS.md — Project patterns
4. Read .context/ARCHITECTURE.md — System overview
5. Read .context/DECISIONS.md — Why things are the way they are

When you make changes:
- Add decisions: <context-update type="decision">Your decision</context-update>
- Add tasks: <context-update type="task">New task</context-update>
- Add learnings: <context-update type="learning">What you learned</context-update>
- Complete tasks: <context-update type="complete">task description</context-update>

Run 'ctx agent' for a quick context summary.
`)
		fmt.Println(green("```"))
		fmt.Println()
		fmt.Println("Or use a hook in .claude/settings.json:")
		fmt.Println()
		fmt.Println(green("```json"))
		fmt.Println(`{
  "hooks": {
    "preToolCall": "ctx agent --budget 4000"
  }
}`)
		fmt.Println(green("```"))

	case "cursor":
		fmt.Println(cyan("Cursor IDE Integration"))
		fmt.Println(cyan("======================"))
		fmt.Println()
		fmt.Println("Add to your .cursorrules file:")
		fmt.Println()
		fmt.Println(green("```markdown"))
		fmt.Print(`# Project Context

Always read these files before making changes:
- .context/CONSTITUTION.md (NEVER violate these rules)
- .context/TASKS.md (current work)
- .context/CONVENTIONS.md (how we write code)
- .context/ARCHITECTURE.md (system structure)

Run 'ctx agent' for a context summary.
Run 'ctx drift' to check for stale context.
`)
		fmt.Println(green("```"))

	case "aider":
		fmt.Println(cyan("Aider Integration"))
		fmt.Println(cyan("================="))
		fmt.Println()
		fmt.Println("Add to your .aider.conf.yml:")
		fmt.Println()
		fmt.Println(green("```yaml"))
		fmt.Println(`read:
  - .context/CONSTITUTION.md
  - .context/TASKS.md
  - .context/CONVENTIONS.md
  - .context/ARCHITECTURE.md
  - .context/DECISIONS.md`)
		fmt.Println(green("```"))
		fmt.Println()
		fmt.Println("Or pass context via command line:")
		fmt.Println()
		fmt.Println(green("```bash"))
		fmt.Println(`ctx agent | aider --message "$(cat -)"`)
		fmt.Println(green("```"))

	case "copilot":
		fmt.Println(cyan("GitHub Copilot Integration"))
		fmt.Println(cyan("=========================="))
		fmt.Println()
		fmt.Println("Add to your .github/copilot-instructions.md:")
		fmt.Println()
		fmt.Println(green("```markdown"))
		fmt.Print(`# Project Context

Before generating code, review:
- .context/CONSTITUTION.md for inviolable rules
- .context/CONVENTIONS.md for coding patterns
- .context/ARCHITECTURE.md for system structure

Key conventions:
- [Add your key conventions here]
`)
		fmt.Println(green("```"))

	case "windsurf":
		fmt.Println(cyan("Windsurf Integration"))
		fmt.Println(cyan("===================="))
		fmt.Println()
		fmt.Println("Add to your .windsurfrules file:")
		fmt.Println()
		fmt.Println(green("```markdown"))
		fmt.Print(`# Active Memory Context

Read order for context:
1. .context/CONSTITUTION.md
2. .context/TASKS.md
3. .context/CONVENTIONS.md
4. .context/ARCHITECTURE.md
5. .context/DECISIONS.md

Run 'ctx agent' for AI-ready context packet.
`)
		fmt.Println(green("```"))

	default:
		fmt.Printf("Unknown tool: %s\n\n", tool)
		fmt.Println("Supported tools:")
		fmt.Println("  claude-code  - Anthropic's Claude Code CLI")
		fmt.Println("  cursor       - Cursor IDE")
		fmt.Println("  aider        - Aider AI coding assistant")
		fmt.Println("  copilot      - GitHub Copilot")
		fmt.Println("  windsurf     - Windsurf IDE")
		return fmt.Errorf("unsupported tool: %s", tool)
	}

	return nil
}
