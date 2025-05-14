package nlp_test

import (
	"testing"

	"github.com/agnath18/lumo/pkg/config"
	"github.com/agnath18/lumo/pkg/nlp"
)

// TestParseHelp tests parsing the help command
func TestParseHelp(t *testing.T) {
	cfg := config.DefaultConfig()
	parser := nlp.NewParser(cfg)

	cmd, err := parser.Parse("help")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if cmd.Type != nlp.CommandTypeHelp {
		t.Errorf("Expected command type to be CommandTypeHelp, got %v", cmd.Type)
	}

	if cmd.Intent != "help" {
		t.Errorf("Expected intent to be 'help', got '%s'", cmd.Intent)
	}
}

// TestParseShellCommands tests parsing shell commands
func TestParseShellCommands(t *testing.T) {
	// Create a config with shell commands enabled in interactive mode
	cfg := config.DefaultConfig()
	cfg.EnableShellInInteractive = true
	parser := nlp.NewParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected nlp.CommandType
		intent   string
	}{
		{
			name:     "lumo prefix",
			input:    "lumo:ls -la",
			expected: nlp.CommandTypeShell,
			intent:   "ls -la",
		},
		{
			name:     "shell prefix",
			input:    "shell:grep -r 'test' .",
			expected: nlp.CommandTypeShell,
			intent:   "grep -r 'test' .",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}

			if cmd.Type != tt.expected {
				t.Errorf("Expected command type to be %v, got %v", tt.expected, cmd.Type)
			}

			if cmd.Intent != tt.intent {
				t.Errorf("Expected intent to be '%s', got '%s'", tt.intent, cmd.Intent)
			}
		})
	}
}

// TestParseAgentCommands tests parsing agent commands
func TestParseAgentCommands(t *testing.T) {
	cfg := config.DefaultConfig()
	parser := nlp.NewParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected nlp.CommandType
		intent   string
	}{
		{
			name:     "auto prefix",
			input:    "auto:create a backup of my documents",
			expected: nlp.CommandTypeAgent,
			intent:   "create a backup of my documents",
		},
		{
			name:     "agent prefix",
			input:    "agent:find large files in the current directory",
			expected: nlp.CommandTypeAgent,
			intent:   "find large files in the current directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}

			if cmd.Type != tt.expected {
				t.Errorf("Expected command type to be %v, got %v", tt.expected, cmd.Type)
			}

			if cmd.Intent != tt.intent {
				t.Errorf("Expected intent to be '%s', got '%s'", tt.intent, cmd.Intent)
			}
		})
	}
}

// TestParseAICommands tests parsing AI commands
func TestParseAICommands(t *testing.T) {
	cfg := config.DefaultConfig()
	parser := nlp.NewParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected nlp.CommandType
		intent   string
	}{
		{
			name:     "ask prefix",
			input:    "ask:What is the capital of France?",
			expected: nlp.CommandTypeAI,
			intent:   "What is the capital of France?",
		},
		{
			name:     "ai prefix",
			input:    "ai:Explain quantum computing",
			expected: nlp.CommandTypeAI,
			intent:   "Explain quantum computing",
		},
		{
			name:     "no prefix",
			input:    "What is the capital of France?",
			expected: nlp.CommandTypeAI,
			intent:   "What is the capital of France?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}

			if cmd.Type != tt.expected {
				t.Errorf("Expected command type to be %v, got %v", tt.expected, cmd.Type)
			}

			if cmd.Intent != tt.intent {
				t.Errorf("Expected intent to be '%s', got '%s'", tt.intent, cmd.Intent)
			}
		})
	}
}

// TestParseSystemCommands tests parsing system commands
func TestParseSystemCommands(t *testing.T) {
	cfg := config.DefaultConfig()
	parser := nlp.NewParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected nlp.CommandType
		intent   string
	}{
		{
			name:     "health prefix",
			input:    "health:check",
			expected: nlp.CommandTypeSystemHealth,
			intent:   "check",
		},
		{
			name:     "syshealth prefix",
			input:    "syshealth:full",
			expected: nlp.CommandTypeSystemHealth,
			intent:   "full",
		},
		{
			name:     "report prefix",
			input:    "report:generate",
			expected: nlp.CommandTypeSystemReport,
			intent:   "generate",
		},
		{
			name:     "sysreport prefix",
			input:    "sysreport:full",
			expected: nlp.CommandTypeSystemReport,
			intent:   "full",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}

			if cmd.Type != tt.expected {
				t.Errorf("Expected command type to be %v, got %v", tt.expected, cmd.Type)
			}

			if cmd.Intent != tt.intent {
				t.Errorf("Expected intent to be '%s', got '%s'", tt.intent, cmd.Intent)
			}
		})
	}
}

// TestParseExitCommands tests parsing exit commands
func TestParseExitCommands(t *testing.T) {
	cfg := config.DefaultConfig()
	parser := nlp.NewParser(cfg)

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "exit",
			input: "exit",
		},
		{
			name:  "quit",
			input: "quit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}

			// Based on the parser implementation, exit/quit commands are treated as AI queries
			// since there's no specific handling for them in the parser
			if cmd.Type != nlp.CommandTypeAI {
				t.Errorf("Expected command type to be CommandTypeAI, got %v", cmd.Type)
			}

			if cmd.Intent != tt.input {
				t.Errorf("Expected intent to be '%s', got '%s'", tt.input, cmd.Intent)
			}
		})
	}
}

// TestParseSpeedTestCommands tests parsing speed test commands
func TestParseSpeedTestCommands(t *testing.T) {
	cfg := config.DefaultConfig()
	parser := nlp.NewParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected nlp.CommandType
		intent   string
	}{
		{
			name:     "speed prefix",
			input:    "speed:test",
			expected: nlp.CommandTypeSpeedTest,
			intent:   "test",
		},
		{
			name:     "speedtest prefix",
			input:    "speedtest:download",
			expected: nlp.CommandTypeSpeedTest,
			intent:   "download",
		},
		{
			name:     "speed-test prefix",
			input:    "speed-test:upload",
			expected: nlp.CommandTypeSpeedTest,
			intent:   "upload",
		},
		{
			name:     "empty speed command",
			input:    "speed:",
			expected: nlp.CommandTypeSpeedTest,
			intent:   "",
		},
		{
			name:     "natural language speed query",
			input:    "check my internet speed",
			expected: nlp.CommandTypeSpeedTest,
			intent:   "check my internet speed",
		},
		{
			name:     "natural language download speed query",
			input:    "how fast is my download speed",
			expected: nlp.CommandTypeSpeedTest,
			intent:   "how fast is my download speed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}

			if cmd.Type != tt.expected {
				t.Errorf("Expected command type to be %v, got %v", tt.expected, cmd.Type)
			}

			if cmd.Intent != tt.intent {
				t.Errorf("Expected intent to be '%s', got '%s'", tt.intent, cmd.Intent)
			}
		})
	}
}
