// Package claude provides Claude Code integration templates and utilities.
package claude

import (
	"embed"
	"fmt"
	"runtime"
	"strings"
)

//go:embed auto-save-session.sh
var FS embed.FS

// GetAutoSaveScript returns the auto-save session script with the binary path substituted.
func GetAutoSaveScript(projectDir string) ([]byte, error) {
	content, err := FS.ReadFile("auto-save-session.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to read auto-save-session.sh: %w", err)
	}

	binaryPath := GetBinaryPath(projectDir)
	result := strings.ReplaceAll(string(content), "{{CTX_BINARY_PATH}}", binaryPath)

	return []byte(result), nil
}

// GetBinaryPath returns the appropriate ctx binary path for the current platform.
// It constructs the path based on OS and architecture.
func GetBinaryPath(projectDir string) string {
	// Build binary name based on platform
	binaryName := fmt.Sprintf("ctx-%s-%s", runtime.GOOS, runtime.GOARCH)

	// Use relative path from project directory
	if projectDir != "" {
		return fmt.Sprintf("%s/dist/%s", projectDir, binaryName)
	}

	// Fallback to just the binary name in dist/
	return fmt.Sprintf("./dist/%s", binaryName)
}

// GetBinaryName returns just the binary filename for the current platform.
func GetBinaryName() string {
	return fmt.Sprintf("ctx-%s-%s", runtime.GOOS, runtime.GOARCH)
}

// SettingsHooks represents the hooks section of settings.local.json
type SettingsHooks struct {
	PreToolUse []HookMatcher `json:"PreToolUse,omitempty"`
	SessionEnd []HookMatcher `json:"SessionEnd,omitempty"`
}

// HookMatcher represents a hook matcher with optional pattern
type HookMatcher struct {
	Matcher string `json:"matcher,omitempty"`
	Hooks   []Hook `json:"hooks"`
}

// Hook represents a single hook command
type Hook struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// Settings represents the full settings.local.json structure
type Settings struct {
	Hooks       SettingsHooks          `json:"hooks,omitempty"`
	Permissions map[string]interface{} `json:"permissions,omitempty"`
}

// CreateDefaultHooks returns the default ctx hooks configuration.
func CreateDefaultHooks(projectDir string) SettingsHooks {
	binaryPath := GetBinaryPath(projectDir)
	hooksDir := ".claude/hooks"
	if projectDir != "" {
		hooksDir = fmt.Sprintf("%s/.claude/hooks", projectDir)
	}

	return SettingsHooks{
		PreToolUse: []HookMatcher{
			{
				Matcher: ".*",
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s agent --budget 4000 2>/dev/null || true", binaryPath),
					},
				},
			},
		},
		SessionEnd: []HookMatcher{
			{
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/auto-save-session.sh", hooksDir),
					},
				},
			},
		},
	}
}
