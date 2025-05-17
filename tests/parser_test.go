package tests

import (
	"testing"

	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/nlp"
)

// TestParserCommandTypeDetection tests the parser's ability to detect different command types
func TestParserCommandTypeDetection(t *testing.T) {
	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
	}

	// Create a parser instance
	parser := nlp.NewParser(cfg)

	// Test cases
	testCases := []struct {
		input        string
		expectedType nlp.CommandType
		description  string
	}{
		// Shell commands
		{"lumo:ls -la", nlp.CommandTypeShell, "Shell command with lumo: prefix"},
		{"shell:echo hello", nlp.CommandTypeShell, "Shell command with shell: prefix"},

		// AI queries
		{"how do I find large files?", nlp.CommandTypeAI, "Natural language query"},
		{"ask:what is Linux?", nlp.CommandTypeAI, "AI query with ask: prefix"},

		// Help commands
		{"help", nlp.CommandTypeHelp, "Help command"},

		// Agent commands
		{"agent:create a backup", nlp.CommandTypeAgent, "Agent command with agent: prefix"},
		{"auto:install nodejs", nlp.CommandTypeAgent, "Agent command with auto: prefix"},

		// System health commands
		{"health:cpu", nlp.CommandTypeSystemHealth, "System health command with health: prefix"},
		{"syshealth:memory", nlp.CommandTypeSystemHealth, "System health command with syshealth: prefix"},

		// System report commands
		{"report:full", nlp.CommandTypeSystemReport, "System report command with report: prefix"},
		{"sysreport:disk", nlp.CommandTypeSystemReport, "System report command with sysreport: prefix"},

		// Chat commands
		{"chat:hello there", nlp.CommandTypeChat, "Chat command with chat: prefix"},
		{"talk:how are you?", nlp.CommandTypeChat, "Chat command with talk: prefix"},

		// Config commands
		{"config:provider list", nlp.CommandTypeConfig, "Config command with config: prefix"},

		// Speed test commands
		{"speed:", nlp.CommandTypeSpeedTest, "Speed test command with speed: prefix"},
		{"speedtest:download", nlp.CommandTypeSpeedTest, "Speed test command with speedtest: prefix"},
		{"speed-test:upload", nlp.CommandTypeSpeedTest, "Speed test command with speed-test: prefix"},

		// Magic commands
		{"magic:dance", nlp.CommandTypeMagic, "Magic command with magic: prefix"},

		// Clipboard commands
		{"clipboard", nlp.CommandTypeClipboard, "Clipboard command"},
		{"clipboard hello world", nlp.CommandTypeClipboard, "Clipboard command with content"},

		// Connect commands
		{"connect", nlp.CommandTypeConnect, "Connect command"},
		{"connect 192.168.1.1", nlp.CommandTypeConnect, "Connect command with IP"},

		// Create commands
		{"create", nlp.CommandTypeCreate, "Create command"},
		{"create:flutter app", nlp.CommandTypeCreate, "Create command with project type"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cmd, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			if cmd.Type != tc.expectedType {
				t.Errorf("Expected command type %v, got %v", tc.expectedType, cmd.Type)
			}
		})
	}
}

// TestParserCommandFirstMode tests the parser's behavior in command-first mode
func TestParserCommandFirstMode(t *testing.T) {
	t.Skip("Skipping test that requires command-first mode detection")
}

// TestParserPrefixHandling tests the parser's handling of command prefixes
func TestParserPrefixHandling(t *testing.T) {
	// Create a default config for testing
	cfg := &config.Config{
		EnableShellInInteractive: true,
		CommandFirstMode:         false, // AI-first mode
	}

	// Create a parser instance
	parser := nlp.NewParser(cfg)

	// Test cases
	testCases := []struct {
		input          string
		expectedIntent string
		description    string
	}{
		{"lumo:ls -la", "ls -la", "Shell command with lumo: prefix"},
		{"shell:echo hello", "echo hello", "Shell command with shell: prefix"},
		{"health:cpu", "cpu", "System health command with health: prefix"},
		{"syshealth:memory", "memory", "System health command with syshealth: prefix"},
		{"report:full", "full", "System report command with report: prefix"},
		{"sysreport:disk", "disk", "System report command with sysreport: prefix"},
		{"config:provider list", "provider list", "Config command with config: prefix"},
		{"speed:download", "download", "Speed test command with speed: prefix"},
		{"speedtest:upload", "upload", "Speed test command with speedtest: prefix"},
		{"speed-test:ping", "ping", "Speed test command with speed-test: prefix"},
		{"magic:dance", "dance", "Magic command with magic: prefix"},
		{"create:flutter app", "flutter app", "Create command with create: prefix"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cmd, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			if cmd.Intent != tc.expectedIntent {
				t.Errorf("Expected intent '%s', got '%s'", tc.expectedIntent, cmd.Intent)
			}
		})
	}
}
