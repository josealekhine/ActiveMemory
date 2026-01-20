package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "context-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a .context directory
	ctxDir := filepath.Join(tmpDir, ".context")
	if err := os.Mkdir(ctxDir, 0755); err != nil {
		t.Fatalf("failed to create .context dir: %v", err)
	}

	tests := []struct {
		name     string
		dir      string
		expected bool
	}{
		{
			name:     "existing directory",
			dir:      ctxDir,
			expected: true,
		},
		{
			name:     "non-existing directory",
			dir:      filepath.Join(tmpDir, "nonexistent"),
			expected: false,
		},
		{
			name:     "file not directory",
			dir:      filepath.Join(tmpDir, "file.txt"),
			expected: false,
		},
	}

	// Create a file for the "file not directory" test
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Exists(tt.dir)
			if result != tt.expected {
				t.Errorf("Exists(%q) = %v, want %v", tt.dir, result, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "context-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a .context directory
	ctxDir := filepath.Join(tmpDir, ".context")
	if err := os.Mkdir(ctxDir, 0755); err != nil {
		t.Fatalf("failed to create .context dir: %v", err)
	}

	// Create some test files
	files := map[string]string{
		"CONSTITUTION.md": "# Constitution\n\n- [ ] Never break the build\n- [ ] Always write tests\n",
		"TASKS.md":        "# Tasks\n\n- [ ] Implement feature A\n- [x] Setup project\n",
		"DECISIONS.md":    "# Decisions\n\n## 2024-01-15 Use PostgreSQL\n\nWe decided to use PostgreSQL.\n",
	}

	for name, content := range files {
		path := filepath.Join(ctxDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
	}

	// Test loading the context
	ctx, err := Load(ctxDir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if ctx.Dir != ctxDir {
		t.Errorf("ctx.Dir = %q, want %q", ctx.Dir, ctxDir)
	}

	if len(ctx.Files) != 3 {
		t.Errorf("len(ctx.Files) = %d, want 3", len(ctx.Files))
	}

	if ctx.TotalTokens == 0 {
		t.Error("ctx.TotalTokens should be > 0")
	}

	if ctx.TotalSize == 0 {
		t.Error("ctx.TotalSize should be > 0")
	}
}

func TestLoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/path/.context")
	if err == nil {
		t.Error("Load() should return error for non-existent directory")
	}

	_, ok := err.(*NotFoundError)
	if !ok {
		t.Errorf("error should be *NotFoundError, got %T", err)
	}
}

func TestIsEffectivelyEmpty(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "truly empty",
			content:  []byte{},
			expected: true,
		},
		{
			name:     "only whitespace",
			content:  []byte("   \n\n  "),
			expected: true,
		},
		{
			name:     "only headers",
			content:  []byte("# Header\n\n## Another Header\n"),
			expected: true,
		},
		{
			name:     "has content",
			content:  []byte("# Header\n\nThis is actual content here.\n"),
			expected: false,
		},
		{
			name:     "has list items",
			content:  []byte("# Tasks\n\n- [ ] Do something important\n"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isEffectivelyEmpty(tt.content)
			if result != tt.expected {
				t.Errorf("isEffectivelyEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Dir: "/test/path"}
	expected := "context directory not found: /test/path"
	if err.Error() != expected {
		t.Errorf("NotFoundError.Error() = %q, want %q", err.Error(), expected)
	}
}
