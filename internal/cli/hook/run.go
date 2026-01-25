//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runHook executes the hook command logic.
//
// Outputs integration instructions and configuration snippets for the
// specified AI tool.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Command arguments; args[0] is the tool name
//
// Returns:
//   - error: Non-nil if the tool is not supported
func runHook(cmd *cobra.Command, args []string) error {
	tool := strings.ToLower(args[0])

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	switch tool {
	case "claude-code", "claude":
		cmd.Println(cyan("Claude Code Integration"))
		cmd.Println(cyan("======================="))
		cmd.Println()
		cmd.Println("Add this to your project's CLAUDE.md or system prompt:")
		cmd.Println()
		cmd.Println(green("```markdown"))
		cmd.Print(`## Context

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
		cmd.Println(green("```"))
		cmd.Println()
		cmd.Println("Or use a hook in .claude/settings.json:")
		cmd.Println()
		cmd.Println(green("```json"))
		cmd.Println(`{
  "hooks": {
    "preToolCall": "ctx agent --budget 4000"
  }
}`)
		cmd.Println(green("```"))

	case "cursor":
		cmd.Println(cyan("Cursor IDE Integration"))
		cmd.Println(cyan("======================"))
		cmd.Println()
		cmd.Println("Add to your .cursorrules file:")
		cmd.Println()
		cmd.Println(green("```markdown"))
		cmd.Print(`# Project Context

Always read these files before making changes:
- .context/CONSTITUTION.md (NEVER violate these rules)
- .context/TASKS.md (current work)
- .context/CONVENTIONS.md (how we write code)
- .context/ARCHITECTURE.md (system structure)

Run 'ctx agent' for a context summary.
Run 'ctx drift' to check for stale context.
`)
		cmd.Println(green("```"))

	case "aider":
		cmd.Println(cyan("Aider Integration"))
		cmd.Println(cyan("================="))
		cmd.Println()
		cmd.Println("Add to your .aider.conf.yml:")
		cmd.Println()
		cmd.Println(green("```yaml"))
		cmd.Println(`read:
  - .context/CONSTITUTION.md
  - .context/TASKS.md
  - .context/CONVENTIONS.md
  - .context/ARCHITECTURE.md
  - .context/DECISIONS.md`)
		cmd.Println(green("```"))
		cmd.Println()
		cmd.Println("Or pass context via command line:")
		cmd.Println()
		cmd.Println(green("```bash"))
		cmd.Println(`ctx agent | aider --message "$(cat -)"`)
		cmd.Println(green("```"))

	case "copilot":
		cmd.Println(cyan("GitHub Copilot Integration"))
		cmd.Println(cyan("=========================="))
		cmd.Println()
		cmd.Println("Add to your .github/copilot-instructions.md:")
		cmd.Println()
		cmd.Println(green("```markdown"))
		cmd.Print(`# Project Context

Before generating code, review:
- .context/CONSTITUTION.md for inviolable rules
- .context/CONVENTIONS.md for coding patterns
- .context/ARCHITECTURE.md for system structure

Key conventions:
- [Add your key conventions here]
`)
		cmd.Println(green("```"))

	case "windsurf":
		cmd.Println(cyan("Windsurf Integration"))
		cmd.Println(cyan("===================="))
		cmd.Println()
		cmd.Println("Add to your .windsurfrules file:")
		cmd.Println()
		cmd.Println(green("```markdown"))
		cmd.Print(`# Context

Read order for context:
1. .context/CONSTITUTION.md
2. .context/TASKS.md
3. .context/CONVENTIONS.md
4. .context/ARCHITECTURE.md
5. .context/DECISIONS.md

Run 'ctx agent' for AI-ready context packet.
`)
		cmd.Println(green("```"))

	default:
		cmd.Printf("Unknown tool: %s\n\n", tool)
		cmd.Println("Supported tools:")
		cmd.Println("  claude-code  - Anthropic's Claude Code CLI")
		cmd.Println("  cursor       - Cursor IDE")
		cmd.Println("  aider        - Aider AI coding assistant")
		cmd.Println("  copilot      - GitHub Copilot")
		cmd.Println("  windsurf     - Windsurf IDE")
		return fmt.Errorf("unsupported tool: %s", tool)
	}

	return nil
}
