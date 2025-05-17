package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
	"github.com/agnath18K/lumo/pkg/terminal"
)

// TestTerminalDisplay tests the terminal's display functionality
func TestTerminalDisplay(t *testing.T) {
	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
	}

	// Create a terminal instance
	term := terminal.NewTerminal(cfg)

	// Test cases
	testCases := []struct {
		result      *executor.Result
		description string
	}{
		{
			&executor.Result{
				Output:     "Hello, world!",
				IsError:    false,
				CommandRun: "echo Hello, world!",
			},
			"Simple output",
		},
		{
			&executor.Result{
				Output:     "Error: command not found",
				IsError:    true,
				CommandRun: "nonexistentcommand",
			},
			"Error output",
		},
		{
			&executor.Result{
				Output:     "Multi\nline\noutput",
				IsError:    false,
				CommandRun: "echo -e 'Multi\\nline\\noutput'",
			},
			"Multi-line output",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// We can't easily test the actual display output since it goes to stdout
			// But we can at least ensure the method doesn't panic
			term.Display(tc.result)
		})
	}
}

// TestTerminalLogCommand tests the terminal's command logging functionality
func TestTerminalLogCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lumo-terminal-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a logs directory
	logsDir := filepath.Join(tempDir, "logs")
	err = os.Mkdir(logsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
		EnableLogging:            true,
	}

	// Change the working directory to the logs directory for this test
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// Create a terminal instance
	term := terminal.NewTerminal(cfg)

	// Test command and result
	command := "echo Hello, world!"
	result := &executor.Result{
		Output:     "Hello, world!",
		IsError:    false,
		CommandRun: command,
	}
	duration := 100 * time.Millisecond

	// Log the command
	term.LogCommand(command, result, duration)

	// Check that the log file was created
	files, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatalf("Failed to read logs directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("No log files were created")
	}

	// Check the content of the log file
	logFile := filepath.Join(logsDir, files[0].Name())
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Check that the log contains the command and result
	if !strings.Contains(logContent, command) {
		t.Errorf("Log file does not contain the command")
	}

	if !strings.Contains(logContent, result.Output) {
		t.Errorf("Log file does not contain the command output")
	}

	if !strings.Contains(logContent, "100ms") {
		t.Errorf("Log file does not contain the duration")
	}
}

// TestTerminalDisplayFormatting tests the terminal's display formatting
func TestTerminalDisplayFormatting(t *testing.T) {
	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
	}

	// Create a terminal instance
	term := terminal.NewTerminal(cfg)

	// Test cases
	testCases := []struct {
		result      *executor.Result
		description string
	}{
		{
			&executor.Result{
				Output:     "Hello, world!",
				IsError:    false,
				CommandRun: "echo Hello, world!",
			},
			"Simple output",
		},
		{
			&executor.Result{
				Output:     "Error: command not found",
				IsError:    true,
				CommandRun: "nonexistentcommand",
			},
			"Error output",
		},
		{
			&executor.Result{
				Output:     "Multi\nline\noutput",
				IsError:    false,
				CommandRun: "echo -e 'Multi\\nline\\noutput'",
			},
			"Multi-line output",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// We can't easily test the actual display output since it goes to stdout
			// But we can at least ensure the method doesn't panic
			term.Display(tc.result)
		})
	}
}
