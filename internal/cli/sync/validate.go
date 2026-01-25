//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// checkPackageFiles detects package manager files without dependency
// documentation.
//
// Checks for common package files (package.json, go.mod, etc.) and suggests
// documenting dependencies if no DEPENDENCIES.md exists or ARCHITECTURE.md
// doesn't mention dependencies.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: Suggested actions for undocumented dependencies
func checkPackageFiles(ctx *context.Context) []Action {
	var actions []Action

	for file, desc := range config.Packages {
		if _, err := os.Stat(file); err == nil {
			// File exists, check if we have DEPENDENCIES.md or similar
			hasDepsDoc := false
			for _, f := range ctx.Files {
				if f.Name == config.FilenameDependency ||
					strings.Contains(strings.ToLower(string(f.Content)),
						"dependencies",
					) {
					hasDepsDoc = true
					break
				}
			}

			if !hasDepsDoc {
				actions = append(actions, Action{
					Type: "DEPS",
					File: "ARCHITECTURE.md",
					Description: fmt.Sprintf(
						"Found %s (%s) but no dependency documentation", file, desc,
					),
					Suggestion: "Consider documenting key dependencies " +
						"in ARCHITECTURE.md or create DEPENDENCIES.md",
				})
			}
		}
	}

	return actions
}

// checkConfigFiles detects config files not documented in CONVENTIONS.md.
//
// Scans for common configuration files (.eslintrc, .prettierrc, tsconfig.json,
// etc.) and suggests documenting them if CONVENTIONS.md doesn't mention the
// related topic.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: Suggested actions for undocumented configurations
func checkConfigFiles(ctx *context.Context) []Action {
	var actions []Action

	for _, cfg := range config.Patterns {
		matches, _ := filepath.Glob(cfg.Pattern)
		if len(matches) > 0 {
			// Check if CONVENTIONS.md mentions this
			var convContent string
			for _, f := range ctx.Files {
				if f.Name == config.FilenameConvention {
					convContent = strings.ToLower(string(f.Content))
					break
				}
			}

			keyword := strings.ToLower(strings.TrimPrefix(cfg.Pattern, "."))
			keyword = strings.TrimSuffix(keyword, "*")
			if convContent == "" || !strings.Contains(convContent, keyword) {
				actions = append(actions, Action{
					Type: "CONFIG",
					File: "CONVENTIONS.md",
					Description: fmt.Sprintf(
						"Found %s but %s not documented", matches[0], cfg.Topic,
					),
					Suggestion: fmt.Sprintf("Document %s in CONVENTIONS.md", cfg.Topic),
				})
			}
		}
	}

	return actions
}

// checkNewDirectories detects important directories not in ARCHITECTURE.md.
//
// Scans top-level directories for common code directories (src, lib, cmd, etc.)
// and suggests documenting them if ARCHITECTURE.md doesn't mention them.
// Skips hidden directories and common non-code directories (node_modules,
// vendor, dist, build).
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: Suggested actions for undocumented directories
func checkNewDirectories(ctx *context.Context) []Action {
	var actions []Action

	// Get ARCHITECTURE.md content
	var archContent string
	for _, f := range ctx.Files {
		if f.Name == config.FilenameArchitecture {
			archContent = strings.ToLower(string(f.Content))
			break
		}
	}

	// Scan top-level directories
	entries, err := os.ReadDir(".")
	if err != nil {
		return actions
	}

	importantDirs := []string{
		"src", "lib", "pkg", "internal",
		"cmd", "api", "web", "app", "services", "components",
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip hidden directories and common non-code directories
		if strings.HasPrefix(name, ".") ||
			name == "node_modules" ||
			name == "vendor" ||
			name == "dist" ||
			name == "build" {
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
			actions = append(actions, Action{
				Type: "NEW_DIR",
				File: "ARCHITECTURE.md",
				Description: fmt.Sprintf(
					"Directory '%s/' exists but not documented", name,
				),
				Suggestion: fmt.Sprintf(
					"Add '%s/' to ARCHITECTURE.md with description", name,
				),
			})
		}
	}

	return actions
}
