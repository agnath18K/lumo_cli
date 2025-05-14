package utils_test

import (
	"os"
	"testing"
	"time"

	"github.com/agnath18/lumo/pkg/utils"
)

// TestFormatDuration tests the FormatDuration function
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "microseconds",
			duration: 500 * time.Microsecond,
			expected: "500 µs",
		},
		{
			name:     "milliseconds",
			duration: 50 * time.Millisecond,
			expected: "50 ms",
		},
		{
			name:     "seconds",
			duration: 5 * time.Second,
			expected: "5.00 s",
		},
		{
			name:     "minutes and seconds",
			duration: 65 * time.Second,
			expected: "1 m 5 s",
		},
		{
			name:     "multiple minutes",
			duration: 125 * time.Second,
			expected: "2 m 5 s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("FormatDuration(%v) = %s, expected %s", tt.duration, result, tt.expected)
			}
		})
	}
}

// TestTruncateString tests the TruncateString function
func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short string",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "exact length",
			input:    "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "long string",
			input:    "hello world",
			maxLen:   8,
			expected: "hello...",
		},
		{
			name:     "very short max length",
			input:    "hello",
			maxLen:   3,
			expected: "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.TruncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateString(%s, %d) = %s, expected %s", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestIsTerminal tests the IsTerminal function
func TestIsTerminal(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "lumo-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// A regular file should not be a terminal
	if utils.IsTerminal(tmpFile) {
		t.Error("Expected IsTerminal to return false for a regular file")
	}

	// We can't easily test the true case as it requires a real terminal
	// but we can at least test that the function doesn't panic
}

// TestSplitCommandArgs tests the SplitCommandArgs function
func TestSplitCommandArgs(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCmd   string
		expectedArgs  []string
		expectedCount int
	}{
		{
			name:          "simple command",
			input:         "ls",
			expectedCmd:   "ls",
			expectedArgs:  []string{},
			expectedCount: 0,
		},
		{
			name:          "command with args",
			input:         "ls -la",
			expectedCmd:   "ls",
			expectedArgs:  []string{"-la"},
			expectedCount: 1,
		},
		{
			name:          "command with multiple args",
			input:         "grep -r 'test' .",
			expectedCmd:   "grep",
			expectedArgs:  []string{"-r", "'test'", "."},
			expectedCount: 3,
		},
		{
			name:          "empty string",
			input:         "",
			expectedCmd:   "",
			expectedArgs:  nil,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := utils.SplitCommandArgs(tt.input)
			
			if cmd != tt.expectedCmd {
				t.Errorf("SplitCommandArgs(%s) command = %s, expected %s", tt.input, cmd, tt.expectedCmd)
			}
			
			if len(args) != tt.expectedCount {
				t.Errorf("SplitCommandArgs(%s) args count = %d, expected %d", tt.input, len(args), tt.expectedCount)
			}
			
			for i, arg := range args {
				if i < len(tt.expectedArgs) && arg != tt.expectedArgs[i] {
					t.Errorf("SplitCommandArgs(%s) arg[%d] = %s, expected %s", tt.input, i, arg, tt.expectedArgs[i])
				}
			}
		})
	}
}

// TestExpandPath tests the ExpandPath function
func TestExpandPath(t *testing.T) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "absolute path",
			input:       "/usr/bin",
			expected:    "/usr/bin",
			expectError: false,
		},
		{
			name:        "relative path",
			input:       "config/file.txt",
			expected:    "config/file.txt",
			expectError: false,
		},
		{
			name:        "home directory",
			input:       "~",
			expected:    homeDir,
			expectError: false,
		},
		{
			name:        "path in home directory",
			input:       "~/config/file.txt",
			expected:    homeDir + "/config/file.txt",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ExpandPath(tt.input)
			
			if tt.expectError && err == nil {
				t.Errorf("ExpandPath(%s) expected error, got nil", tt.input)
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("ExpandPath(%s) unexpected error: %v", tt.input, err)
			}
			
			if result != tt.expected {
				t.Errorf("ExpandPath(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCleanMarkdown tests the CleanMarkdown function
func TestCleanMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bold text",
			input:    "This is **bold** text",
			expected: "This is bold text",
		},
		{
			name:     "italic text",
			input:    "This is *italic* text",
			expected: "This is italic text",
		},
		{
			name:     "inline code",
			input:    "This is `code` text",
			expected: "This is code text",
		},
		{
			name:     "bullet points",
			input:    "* Item 1\n* Item 2",
			expected: "- Item 1\n- Item 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CleanMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("CleanMarkdown(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}
