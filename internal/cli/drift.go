package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/josealekhine/ActiveMemory/internal/context"
	"github.com/josealekhine/ActiveMemory/internal/drift"
	"github.com/spf13/cobra"
)

var (
	driftJSON bool
	driftFix  bool
)

// DriftCmd returns the drift command.
func DriftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect stale or invalid context",
		Long: `Run drift detection to find stale paths, broken references, and constitution violations.

Checks performed:
  - Path references in ARCHITECTURE.md and CONVENTIONS.md exist
  - Staleness indicators (many completed tasks)
  - Constitution rule violations (potential secrets)
  - Required files are present

Use --json for machine-readable output.`,
		RunE: runDrift,
	}

	cmd.Flags().BoolVar(&driftJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&driftFix, "fix", false, "Auto-fix simple issues (not yet implemented)")

	return cmd
}

// DriftJSONOutput represents the JSON output format.
type DriftJSONOutput struct {
	Timestamp  string        `json:"timestamp"`
	Status     string        `json:"status"`
	Warnings   []drift.Issue `json:"warnings"`
	Violations []drift.Issue `json:"violations"`
	Passed     []string      `json:"passed"`
}

func runDrift(cmd *cobra.Command, args []string) error {
	ctx, err := context.Load("")
	if err != nil {
		if _, ok := err.(*context.NotFoundError); ok {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	report := drift.Detect(ctx)

	if driftJSON {
		return outputDriftJSON(report)
	}

	return outputDriftText(report)
}

func outputDriftJSON(report *drift.Report) error {
	output := DriftJSONOutput{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Status:     report.Status(),
		Warnings:   report.Warnings,
		Violations: report.Violations,
		Passed:     report.Passed,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func outputDriftText(report *drift.Report) error {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(cyan("Drift Detection Report"))
	fmt.Println(cyan("======================"))
	fmt.Println()

	// Violations
	if len(report.Violations) > 0 {
		fmt.Printf("%s VIOLATIONS (%d)\n\n", red("❌"), len(report.Violations))
		for _, v := range report.Violations {
			if v.Line > 0 {
				fmt.Printf("  - %s:%d %s", v.File, v.Line, v.Message)
			} else {
				fmt.Printf("  - %s: %s", v.File, v.Message)
			}
			if v.Rule != "" {
				fmt.Printf(" (rule: %s)", v.Rule)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Warnings
	if len(report.Warnings) > 0 {
		fmt.Printf("%s WARNINGS (%d)\n\n", yellow("⚠️ "), len(report.Warnings))

		// Group by type
		pathRefs := []drift.Issue{}
		staleness := []drift.Issue{}
		other := []drift.Issue{}

		for _, w := range report.Warnings {
			switch w.Type {
			case "dead_path":
				pathRefs = append(pathRefs, w)
			case "staleness":
				staleness = append(staleness, w)
			default:
				other = append(other, w)
			}
		}

		if len(pathRefs) > 0 {
			fmt.Println("  Path References:")
			for _, w := range pathRefs {
				fmt.Printf("  - %s:%d references '%s' (not found)\n", w.File, w.Line, w.Path)
			}
			fmt.Println()
		}

		if len(staleness) > 0 {
			fmt.Println("  Staleness:")
			for _, w := range staleness {
				fmt.Printf("  - %s %s\n", w.File, w.Message)
			}
			fmt.Println()
		}

		if len(other) > 0 {
			fmt.Println("  Other:")
			for _, w := range other {
				fmt.Printf("  - %s: %s\n", w.File, w.Message)
			}
			fmt.Println()
		}
	}

	// Passed
	if len(report.Passed) > 0 {
		fmt.Printf("%s PASSED (%d)\n", green("✅"), len(report.Passed))
		for _, p := range report.Passed {
			fmt.Printf("  - %s\n", formatCheckName(p))
		}
		fmt.Println()
	}

	// Summary
	status := report.Status()
	switch status {
	case "violation":
		fmt.Printf("\nStatus: %s — Constitution violations detected\n", red("VIOLATION"))
		return fmt.Errorf("drift detection found violations")
	case "warning":
		fmt.Printf("\nStatus: %s — Issues detected that should be addressed\n", yellow("WARNING"))
	default:
		fmt.Printf("\nStatus: %s — No drift detected\n", green("OK"))
	}

	return nil
}

func formatCheckName(name string) string {
	switch name {
	case "path_references":
		return "Path references are valid"
	case "staleness_check":
		return "No staleness indicators"
	case "constitution_check":
		return "Constitution rules respected"
	case "required_files":
		return "All required files present"
	default:
		return name
	}
}
