//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestTruncate tests the truncate helper function.
func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a longer string", 10, "this is..."},
		{"", 10, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestParseIndex tests the parseIndex helper function.
func TestParseIndex(t *testing.T) {
	tests := []struct {
		input     string
		expected  int
		expectErr bool
	}{
		{"1", 1, false},
		{"10", 10, false},
		{"100", 100, false},
		{"0", 0, true},  // index must be positive
		{"-1", 0, true}, // index must be positive
		{"abc", 0, true},
		{"", 0, true},
		{"1.5", 1, false}, // Sscanf stops at the decimal
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseIndex(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("parseIndex(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("parseIndex(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("parseIndex(%q) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestCleanInsight tests the cleanInsight helper function.
func TestCleanInsight(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple text", "simple text"},
		{"  trimmed  ", "trimmed"},
		{"ends with period.", "ends with period"},
		{"ends with comma,", "ends with comma"},
		{"ends with multiple...", "ends with multiple"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := cleanInsight(tt.input)
			if result != tt.expected {
				t.Errorf("cleanInsight(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractTextContent tests the extractTextContent helper function.
func TestExtractTextContent(t *testing.T) {
	tests := []struct {
		name     string
		entry    transcriptEntry
		expected []string
	}{
		{
			name: "string content",
			entry: transcriptEntry{
				Message: transcriptMsg{
					Content: "Hello world",
				},
			},
			expected: []string{"Hello world"},
		},
		{
			name: "array content with text",
			entry: transcriptEntry{
				Message: transcriptMsg{
					Content: []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "First text",
						},
						map[string]interface{}{
							"type": "text",
							"text": "Second text",
						},
					},
				},
			},
			expected: []string{"First text", "Second text"},
		},
		{
			name: "array content with thinking",
			entry: transcriptEntry{
				Message: transcriptMsg{
					Content: []interface{}{
						map[string]interface{}{
							"type":     "thinking",
							"thinking": "Some thinking",
						},
					},
				},
			},
			expected: []string{"Some thinking"},
		},
		{
			name: "empty content",
			entry: transcriptEntry{
				Message: transcriptMsg{
					Content: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTextContent(tt.entry)
			if len(result) != len(tt.expected) {
				t.Errorf("extractTextContent() returned %d items, want %d", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("extractTextContent()[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

// TestSessionCommands tests the session subcommands.
func TestSessionCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test session save
	t.Run("session save", func(t *testing.T) {
		sessionCmd := Cmd()
		sessionCmd.SetArgs([]string{"save", "test-topic"})
		if err := sessionCmd.Execute(); err != nil {
			t.Fatalf("session save failed: %v", err)
		}

		// Verify session file was created
		entries, err := os.ReadDir(".context/sessions")
		if err != nil {
			t.Fatalf("failed to read sessions dir: %v", err)
		}
		found := false
		for _, e := range entries {
			if strings.Contains(e.Name(), "test-topic") {
				found = true
				break
			}
		}
		if !found {
			t.Error("session file was not created")
		}
	})

	// Test session list
	t.Run("session list", func(t *testing.T) {
		sessionCmd := Cmd()
		sessionCmd.SetArgs([]string{"list"})
		if err := sessionCmd.Execute(); err != nil {
			t.Fatalf("session list failed: %v", err)
		}
	})

	// Test session load
	t.Run("session load", func(t *testing.T) {
		sessionCmd := Cmd()
		sessionCmd.SetArgs([]string{"load", "1"})
		if err := sessionCmd.Execute(); err != nil {
			t.Fatalf("session load failed: %v", err)
		}
	})
}

// TestSessionParse tests the session parse command with a jsonl file.
func TestSessionParse(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-parse-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a test jsonl file
	jsonlContent := `{"type":"user","message":{"role":"user","content":"Hello"},"timestamp":"2025-01-21T10:00:00Z"}
{"type":"assistant","message":{"role":"assistant","content":"Hi there!"},"timestamp":"2025-01-21T10:00:05Z"}
`
	jsonlPath := filepath.Join(tmpDir, "test-transcript.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	// Test session parse
	sessionCmd := Cmd()
	sessionCmd.SetArgs([]string{"parse", jsonlPath})
	if err := sessionCmd.Execute(); err != nil {
		t.Fatalf("session parse failed: %v", err)
	}
}

// TestSessionParseWithExtract tests the session parse command with --extract flag.
func TestSessionParseWithExtract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-parse-extract-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a test jsonl file with content that should trigger extraction
	jsonlContent := `{"type":"assistant","message":{"role":"assistant","content":"We decided to use PostgreSQL for the database. I learned that connection pooling is important."},"timestamp":"2025-01-21T10:00:00Z"}
`
	jsonlPath := filepath.Join(tmpDir, "test-transcript.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	// Test session parse with --extract
	sessionCmd := Cmd()
	sessionCmd.SetArgs([]string{"parse", jsonlPath, "--extract"})
	if err := sessionCmd.Execute(); err != nil {
		t.Fatalf("session parse --extract failed: %v", err)
	}
}

// TestSessionParseWithOutput tests the session parse command with --output flag.
func TestSessionParseWithOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-parse-output-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a test jsonl file
	jsonlContent := `{"type":"user","message":{"role":"user","content":"Hello"},"timestamp":"2025-01-21T10:00:00Z"}
`
	jsonlPath := filepath.Join(tmpDir, "test-transcript.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.md")

	// Test session parse with --output
	sessionCmd := Cmd()
	sessionCmd.SetArgs([]string{"parse", jsonlPath, "--output", outputPath})
	if err := sessionCmd.Execute(); err != nil {
		t.Fatalf("session parse --output failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}
