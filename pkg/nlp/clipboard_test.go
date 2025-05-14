package nlp

import (
	"testing"

	"github.com/agnath18/lumo/pkg/config"
)

func TestParser_Parse_Clipboard(t *testing.T) {
	cfg := &config.Config{}
	parser := NewParser(cfg)

	// Test clipboard command with no arguments
	cmd, err := parser.Parse("clipboard")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cmd.Type != CommandTypeClipboard {
		t.Errorf("Expected CommandTypeClipboard, got %v", cmd.Type)
	}
	if cmd.Intent != "" {
		t.Errorf("Expected empty intent, got '%s'", cmd.Intent)
	}

	// Test clipboard command with arguments
	cmd, err = parser.Parse("clipboard Hello World")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cmd.Type != CommandTypeClipboard {
		t.Errorf("Expected CommandTypeClipboard, got %v", cmd.Type)
	}
	if cmd.Intent != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", cmd.Intent)
	}

	// Test clipboard command with quoted arguments
	cmd, err = parser.Parse("clipboard \"Hello, World!\"")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cmd.Type != CommandTypeClipboard {
		t.Errorf("Expected CommandTypeClipboard, got %v", cmd.Type)
	}
	if cmd.Intent != "\"Hello, World!\"" {
		t.Errorf("Expected '\"Hello, World!\"', got '%s'", cmd.Intent)
	}
}
