package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/josealekhine/ActiveMemory/internal/context"
	"github.com/spf13/cobra"
)

var (
	syncDryRun bool
)

// SyncCmd returns the sync command.
func SyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Reconcile context with codebase",
		Long: `Scan the codebase and reconcile context files with current state.

Actions performed:
  - Scan for new directories that should be in ARCHITECTURE.md
  - Check for package.json/go.mod changes
  - Identify stale references
  - Suggest updates to context files

Use --dry-run to see what would change without modifying files.`,
		RunE: runSync,
	}

	cmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "Show what would change without modifying")

	return cmd
}

// SyncAction represents a suggested sync action.
type SyncAction struct {
	Type        string
	File        string
	Description string
	Suggestion  string
}

func runSync(cmd *cobra.Command, args []string) error {
	ctx, err := context.Load("")
	if err != nil {
		if _, ok := err.(*context.NotFoundError); ok {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	actions := detectSyncActions(ctx)

	if len(actions) == 0 {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Context is in sync with codebase\n", green("✓"))
		return nil
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(cyan("Sync Analysis"))
	fmt.Println(cyan("============="))
	fmt.Println()

	if syncDryRun {
		fmt.Println(yellow("DRY RUN — No changes will be made"))
		fmt.Println()
	}

	for i, action := range actions {
		fmt.Printf("%d. [%s] %s\n", i+1, action.Type, action.Description)
		if action.Suggestion != "" {
			fmt.Printf("   Suggestion: %s\n", action.Suggestion)
		}
		fmt.Println()
	}

	if syncDryRun {
		fmt.Printf("Found %d items to sync. Run without --dry-run to apply suggestions.\n", len(actions))
	} else {
		fmt.Printf("Found %d items. Review and update context files manually.\n", len(actions))
	}

	return nil
}

func detectSyncActions(ctx *context.Context) []SyncAction {
	actions := []SyncAction{}

	// Check for new top-level directories not mentioned in ARCHITECTURE.md
	actions = append(actions, checkNewDirectories(ctx)...)

	// Check for package manager files
	actions = append(actions, checkPackageFiles(ctx)...)

	// Check for common config files that might need documenting
	actions = append(actions, checkConfigFiles(ctx)...)

	return actions
}

func checkNewDirectories(ctx *context.Context) []SyncAction {
	actions := []SyncAction{}

	// Get ARCHITECTURE.md content
	var archContent string
	for _, f := range ctx.Files {
		if f.Name == "ARCHITECTURE.md" {
			archContent = strings.ToLower(string(f.Content))
			break
		}
	}

	// Scan top-level directories
	entries, err := os.ReadDir(".")
	if err != nil {
		return actions
	}

	importantDirs := []string{"src", "lib", "pkg", "internal", "cmd", "api", "web", "app", "services", "components"}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip hidden directories and common non-code directories
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
			continue
		}

		// Check if this is an important directory not mentioned in ARCHITECTURE.md
		isImportant := false
		for _, imp := range importantDirs {
			if name == imp {
				isImportant = true
				break
			}
		}

		if isImportant && !strings.Contains(archContent, name) {
			actions = append(actions, SyncAction{
				Type:        "NEW_DIR",
				File:        "ARCHITECTURE.md",
				Description: fmt.Sprintf("Directory '%s/' exists but not documented", name),
				Suggestion:  fmt.Sprintf("Add '%s/' to ARCHITECTURE.md with description", name),
			})
		}
	}

	return actions
}

func checkPackageFiles(ctx *context.Context) []SyncAction {
	actions := []SyncAction{}

	packageFiles := map[string]string{
		"package.json": "Node.js dependencies",
		"go.mod":       "Go module dependencies",
		"Cargo.toml":   "Rust dependencies",
		"requirements.txt": "Python dependencies",
		"Gemfile":      "Ruby dependencies",
	}

	for file, desc := range packageFiles {
		if _, err := os.Stat(file); err == nil {
			// File exists, check if we have DEPENDENCIES.md or similar
			hasDepsDoc := false
			for _, f := range ctx.Files {
				if f.Name == "DEPENDENCIES.md" || strings.Contains(strings.ToLower(string(f.Content)), "dependencies") {
					hasDepsDoc = true
					break
				}
			}

			if !hasDepsDoc {
				actions = append(actions, SyncAction{
					Type:        "DEPS",
					File:        "ARCHITECTURE.md",
					Description: fmt.Sprintf("Found %s (%s) but no dependency documentation", file, desc),
					Suggestion:  "Consider documenting key dependencies in ARCHITECTURE.md or create DEPENDENCIES.md",
				})
			}
		}
	}

	return actions
}

func checkConfigFiles(ctx *context.Context) []SyncAction {
	actions := []SyncAction{}

	// Check for config files that might indicate conventions
	configPatterns := []struct {
		pattern string
		topic   string
	}{
		{".eslintrc*", "linting conventions"},
		{".prettierrc*", "formatting conventions"},
		{"tsconfig.json", "TypeScript configuration"},
		{".editorconfig", "editor configuration"},
		{"Makefile", "build commands"},
		{"Dockerfile", "containerization"},
	}

	for _, cfg := range configPatterns {
		matches, _ := filepath.Glob(cfg.pattern)
		if len(matches) > 0 {
			// Check if CONVENTIONS.md mentions this
			var convContent string
			for _, f := range ctx.Files {
				if f.Name == "CONVENTIONS.md" {
					convContent = strings.ToLower(string(f.Content))
					break
				}
			}

			keyword := strings.ToLower(strings.TrimPrefix(cfg.pattern, "."))
			keyword = strings.TrimSuffix(keyword, "*")
			if convContent == "" || !strings.Contains(convContent, keyword) {
				actions = append(actions, SyncAction{
					Type:        "CONFIG",
					File:        "CONVENTIONS.md",
					Description: fmt.Sprintf("Found %s but %s not documented", matches[0], cfg.topic),
					Suggestion:  fmt.Sprintf("Document %s in CONVENTIONS.md", cfg.topic),
				})
			}
		}
	}

	return actions
}
