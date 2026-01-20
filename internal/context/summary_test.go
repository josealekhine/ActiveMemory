package context

import "testing"

func TestGenerateSummary(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  []byte
		expected string
	}{
		{
			name:     "constitution with invariants",
			filename: "CONSTITUTION.md",
			content:  []byte("# Constitution\n\n- [ ] Rule 1\n- [ ] Rule 2\n- [x] Rule 3\n"),
			expected: "3 invariants",
		},
		{
			name:     "constitution empty",
			filename: "CONSTITUTION.md",
			content:  []byte("# Constitution\n"),
			expected: "loaded",
		},
		{
			name:     "tasks mixed",
			filename: "TASKS.md",
			content:  []byte("# Tasks\n\n- [ ] Task 1\n- [ ] Task 2\n- [x] Done\n"),
			expected: "2 active, 1 completed",
		},
		{
			name:     "tasks only active",
			filename: "TASKS.md",
			content:  []byte("# Tasks\n\n- [ ] Task 1\n- [ ] Task 2\n"),
			expected: "2 active",
		},
		{
			name:     "tasks only completed",
			filename: "TASKS.md",
			content:  []byte("# Tasks\n\n- [x] Done 1\n- [x] Done 2\n"),
			expected: "2 completed",
		},
		{
			name:     "tasks empty",
			filename: "TASKS.md",
			content:  []byte("# Tasks\n"),
			expected: "empty",
		},
		{
			name:     "decisions multiple",
			filename: "DECISIONS.md",
			content:  []byte("# Decisions\n\n## 2024-01-15 First\n\nContent\n\n## 2024-01-16 Second\n\nContent\n"),
			expected: "2 decisions",
		},
		{
			name:     "decisions single",
			filename: "DECISIONS.md",
			content:  []byte("# Decisions\n\n## One decision\n\nContent\n"),
			expected: "1 decision",
		},
		{
			name:     "decisions empty",
			filename: "DECISIONS.md",
			content:  []byte("# Decisions\n"),
			expected: "empty",
		},
		{
			name:     "glossary multiple",
			filename: "GLOSSARY.md",
			content:  []byte("# Glossary\n\n- **Term1** - Definition 1\n- **Term2** - Definition 2\n"),
			expected: "2 terms",
		},
		{
			name:     "glossary single",
			filename: "GLOSSARY.md",
			content:  []byte("# Glossary\n\n**SingleTerm** - Definition\n"),
			expected: "1 term",
		},
		{
			name:     "glossary empty",
			filename: "GLOSSARY.md",
			content:  []byte("# Glossary\n"),
			expected: "empty",
		},
		{
			name:     "unknown file with content",
			filename: "OTHER.md",
			content:  []byte("# Other\n\nSome content here\n"),
			expected: "loaded",
		},
		{
			name:     "unknown file empty",
			filename: "OTHER.md",
			content:  []byte(""),
			expected: "empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSummary(tt.filename, tt.content)
			if result != tt.expected {
				t.Errorf("generateSummary(%q) = %q, want %q", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestSummarizeConstitution(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected string
	}{
		{
			name:     "no checkboxes",
			content:  []byte("# Constitution\n\nNo rules here\n"),
			expected: "loaded",
		},
		{
			name:     "with checkboxes",
			content:  []byte("- [ ] Rule 1\n- [x] Rule 2\n"),
			expected: "2 invariants",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := summarizeConstitution(tt.content)
			if result != tt.expected {
				t.Errorf("summarizeConstitution() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSummarizeTasks(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected string
	}{
		{
			name:     "no tasks",
			content:  []byte("# Tasks\n"),
			expected: "empty",
		},
		{
			name:     "only active",
			content:  []byte("- [ ] Task 1\n- [ ] Task 2\n"),
			expected: "2 active",
		},
		{
			name:     "only completed",
			content:  []byte("- [x] Done 1\n"),
			expected: "1 completed",
		},
		{
			name:     "mixed",
			content:  []byte("- [ ] Task\n- [x] Done\n- [x] Also done\n"),
			expected: "1 active, 2 completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := summarizeTasks(tt.content)
			if result != tt.expected {
				t.Errorf("summarizeTasks() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSummarizeDecisions(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected string
	}{
		{
			name:     "no decisions",
			content:  []byte("# Decisions\n"),
			expected: "empty",
		},
		{
			name:     "one decision",
			content:  []byte("## First\n\nContent\n"),
			expected: "1 decision",
		},
		{
			name:     "multiple decisions",
			content:  []byte("## First\n\n## Second\n\n## Third\n"),
			expected: "3 decisions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := summarizeDecisions(tt.content)
			if result != tt.expected {
				t.Errorf("summarizeDecisions() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSummarizeGlossary(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected string
	}{
		{
			name:     "no terms",
			content:  []byte("# Glossary\n"),
			expected: "empty",
		},
		{
			name:     "one term",
			content:  []byte("**Term** - Definition\n"),
			expected: "1 term",
		},
		{
			name:     "multiple terms with list",
			content:  []byte("- **Term1** - Def 1\n- **Term2** - Def 2\n"),
			expected: "2 terms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := summarizeGlossary(tt.content)
			if result != tt.expected {
				t.Errorf("summarizeGlossary() = %q, want %q", result, tt.expected)
			}
		})
	}
}
