package context

import "testing"

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected int
	}{
		{
			name:     "empty content",
			content:  []byte{},
			expected: 0,
		},
		{
			name:     "single character",
			content:  []byte("a"),
			expected: 1, // (1+3)/4 = 1
		},
		{
			name:     "four characters",
			content:  []byte("test"),
			expected: 1, // (4+3)/4 = 1
		},
		{
			name:     "five characters",
			content:  []byte("tests"),
			expected: 2, // (5+3)/4 = 2
		},
		{
			name:     "typical sentence",
			content:  []byte("This is a test sentence with multiple words."),
			expected: 11, // (44+3)/4 = 11
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EstimateTokens(tt.content)
			if result != tt.expected {
				t.Errorf("EstimateTokens() = %d, want %d (len=%d)", result, tt.expected, len(tt.content))
			}
		})
	}
}

func TestEstimateTokensString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "short string",
			input:    "hello",
			expected: 2, // (5+3)/4 = 2
		},
		{
			name:     "longer string",
			input:    "Hello, world! This is a test.",
			expected: 8, // (29+3)/4 = 8
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EstimateTokensString(tt.input)
			if result != tt.expected {
				t.Errorf("EstimateTokensString(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}
