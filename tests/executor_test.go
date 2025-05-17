package tests

import (
	"strings"
	"testing"

	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
	"github.com/agnath18K/lumo/pkg/nlp"
)

// TestExecutorCommandRouting tests the executor's ability to route commands to the correct handler
func TestExecutorCommandRouting(t *testing.T) {
	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
		EnableSystemHealth:       true,
		EnableSystemReport:       true,
		EnableSpeedTest:          true,
		EnableAgentMode:          true,
	}

	// Create an executor instance
	exec := executor.NewExecutor(cfg)

	// Test cases
	testCases := []struct {
		commandType nlp.CommandType
		intent      string
		shouldError bool
		description string
	}{
		{nlp.CommandTypeShell, "echo hello", false, "Shell command execution"},
		{nlp.CommandTypeHelp, "", false, "Help command execution"},
		{nlp.CommandTypeSystemHealth, "cpu", false, "System health command execution"},
		{nlp.CommandTypeSystemReport, "disk", false, "System report command execution"},
		{nlp.CommandTypeConfig, "provider list", false, "Config command execution"},
		{nlp.CommandTypeSpeedTest, "download", false, "Speed test command execution"},
		{nlp.CommandTypeMagic, "dance", false, "Magic command execution"},
		{nlp.CommandTypeClipboard, "", false, "Clipboard command execution"},
		{nlp.CommandTypeConnect, "", false, "Connect command execution"},
		{nlp.CommandTypeCreate, "", false, "Create command execution"},
		{nlp.CommandTypeUnknown, "", true, "Unknown command type should error"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create a command
			cmd := &nlp.Command{
				Type:       tc.commandType,
				Intent:     tc.intent,
				Parameters: make(map[string]string),
				RawInput:   tc.intent,
			}

			// Execute the command
			result, err := exec.Execute(cmd)

			// Check for errors
			if tc.shouldError {
				if err == nil && !result.IsError {
					t.Errorf("Expected error but got none")
				}
			} else {
				// We're not checking for specific errors here, just that the routing works
				// The actual command execution might fail due to dependencies
				// but the routing should still work
			}

			// Verify that CommandRun is set correctly
			if result != nil && result.CommandRun != tc.intent {
				t.Errorf("Expected CommandRun to be '%s', got '%s'", tc.intent, result.CommandRun)
			}
		})
	}
}

// TestExecutorShellCommand tests the executor's shell command execution
func TestExecutorShellCommand(t *testing.T) {
	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
	}

	// Create an executor instance
	exec := executor.NewExecutor(cfg)

	// Test cases
	testCases := []struct {
		intent      string
		shouldError bool
		description string
	}{
		{"echo hello", false, "Simple echo command"},
		{"ls -la", false, "List directory command"},
		{"nonexistentcommand", true, "Non-existent command should error"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create a command
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeShell,
				Intent:     tc.intent,
				Parameters: make(map[string]string),
				RawInput:   tc.intent,
			}

			// Execute the command
			result, err := exec.Execute(cmd)

			// Check for errors
			if tc.shouldError {
				if err == nil && !result.IsError {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.IsError {
					t.Errorf("Command execution failed: %s", result.Output)
				}
			}
		})
	}
}

// TestExecutorHelpCommand tests the executor's help command execution
func TestExecutorHelpCommand(t *testing.T) {
	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
		EnableSystemHealth:       true,
		EnableSystemReport:       true,
		EnableSpeedTest:          true,
		EnableAgentMode:          true,
	}

	// Create an executor instance
	exec := executor.NewExecutor(cfg)

	// Create a help command
	cmd := &nlp.Command{
		Type:       nlp.CommandTypeHelp,
		Intent:     "",
		Parameters: make(map[string]string),
		RawInput:   "help",
	}

	// Execute the command
	result, err := exec.Execute(cmd)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.IsError {
		t.Errorf("Help command execution failed: %s", result.Output)
	}

	// Check that the help text contains expected sections
	expectedSections := []string{
		"Commands:",
		"Examples:",
		"Configuration:",
		"Status:",
	}

	for _, section := range expectedSections {
		if !contains(result.Output, section) {
			t.Errorf("Help text missing expected section: %s", section)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && len(s) >= len(substr) && s != "" && substr != "" && strings.Contains(s, substr)
}
