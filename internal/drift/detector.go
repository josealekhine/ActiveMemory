//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Packages drift provides functionality for detecting stale or invalid context.
package drift

import (
	"os"
	"regexp"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// Issue represents a detected drift issue.
type Issue struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
	Rule    string `json:"rule,omitempty"`
}

// Report represents the complete drift detection report.
type Report struct {
	Warnings   []Issue  `json:"warnings"`
	Violations []Issue  `json:"violations"`
	Passed     []string `json:"passed"`
}

// Status returns the overall status of the report.
func (r *Report) Status() string {
	if len(r.Violations) > 0 {
		return "violation"
	}
	if len(r.Warnings) > 0 {
		return "warning"
	}
	return "ok"
}

// Detect runs all drift detection checks on the given context.
func Detect(ctx *context.Context) *Report {
	report := &Report{
		Warnings:   []Issue{},
		Violations: []Issue{},
		Passed:     []string{},
	}

	// Check path references in context files
	checkPathReferences(ctx, report)

	// Check for staleness indicators
	checkStaleness(ctx, report)

	// Check constitution rules (basic heuristics)
	checkConstitution(ctx, report)

	// Check for empty required files
	checkRequiredFiles(ctx, report)

	return report
}

func checkPathReferences(ctx *context.Context, report *Report) {
	// Pattern to match file paths in Markdown (backticks or code blocks)
	pathPattern := regexp.MustCompile("`([^`]+\\.[a-zA-Z]{1,5})`")

	foundDeadPaths := false

	for _, f := range ctx.Files {
		if f.Name != config.FilenameArchitecture && f.Name != config.FilenameConvention {
			continue
		}

		lines := strings.Split(string(f.Content), "\n")
		for lineNum, line := range lines {
			matches := pathPattern.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				path := m[1]
				// Skip URLs and common non-file patterns
				if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "//") {
					continue
				}
				// Skip template patterns
				if strings.Contains(path, "{") || strings.Contains(path, "*") {
					continue
				}
				// Check if the file exists
				if _, err := os.Stat(path); os.IsNotExist(err) {
					report.Warnings = append(report.Warnings, Issue{
						File:    f.Name,
						Line:    lineNum + 1,
						Type:    "dead_path",
						Message: "references path that does not exist",
						Path:    path,
					})
					foundDeadPaths = true
				}
			}
		}
	}

	if !foundDeadPaths {
		report.Passed = append(report.Passed, "path_references")
	}
}

func checkStaleness(ctx *context.Context, report *Report) {
	staleness := false

	for _, f := range ctx.Files {
		if f.Name == config.FilenameTask {
			// Count completed tasks
			completedCount := strings.Count(string(f.Content), "- [x]")
			if completedCount > 10 {
				report.Warnings = append(report.Warnings, Issue{
					File:    f.Name,
					Type:    "staleness",
					Message: "has many completed items (consider archiving)",
					Path:    "",
				})
				staleness = true
			}
		}
	}

	if !staleness {
		report.Passed = append(report.Passed, "staleness_check")
	}
}

func checkConstitution(_ *context.Context, report *Report) {
	// Basic heuristic checks for constitution violations
	// Check for potential secrets in common config files

	secretPatterns := []string{
		".env",
		"credentials",
		"secret",
		"api_key",
		"apikey",
		"password",
	}

	// Look for common secret file patterns in the working directory
	entries, err := os.ReadDir(".")
	if err != nil {
		return
	}

	foundViolation := false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		for _, pattern := range secretPatterns {
			if strings.Contains(name, pattern) && !strings.HasSuffix(name, ".example") && !strings.HasSuffix(name, ".sample") {
				// Check if it contains actual content (not just template)
				content, err := os.ReadFile(entry.Name())
				if err != nil {
					continue
				}
				if len(content) > 0 && !isTemplateFile(content) {
					report.Violations = append(report.Violations, Issue{
						File:    entry.Name(),
						Type:    "potential_secret",
						Message: "may contain secrets (constitution violation)",
						Rule:    "no_secrets",
					})
					foundViolation = true
				}
			}
		}
	}

	if !foundViolation {
		report.Passed = append(report.Passed, "constitution_check")
	}
}

func checkRequiredFiles(ctx *context.Context, report *Report) {
	allPresent := true

	existingFiles := make(map[string]bool)
	for _, f := range ctx.Files {
		existingFiles[f.Name] = true
	}

	for _, name := range config.RequiredFiles {
		if !existingFiles[name] {
			report.Warnings = append(report.Warnings, Issue{
				File:    name,
				Type:    "missing_file",
				Message: "required context file is missing",
			})
			allPresent = false
		}
	}

	if allPresent {
		report.Passed = append(report.Passed, "required_files")
	}
}

func isTemplateFile(content []byte) bool {
	s := string(content)
	// Check for common template markers
	templateMarkers := []string{
		"YOUR_",
		"<your",
		"{{",
		"REPLACE_",
		"TODO:",
		"CHANGEME",
	}
	for _, marker := range templateMarkers {
		if strings.Contains(strings.ToUpper(s), marker) {
			return true
		}
	}
	return false
}
