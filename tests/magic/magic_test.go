package magic_test

import (
	"strings"
	"testing"

	"github.com/agnath18/lumo/pkg/magic"
)

// TestMagicCreation tests that a new Magic instance can be created
func TestMagicCreation(t *testing.T) {
	m := magic.NewMagic()
	if m == nil {
		t.Fatal("Expected Magic instance to be created, got nil")
	}
}

// TestExecuteUnknownCommand tests that an unknown command returns a helpful message
func TestExecuteUnknownCommand(t *testing.T) {
	m := magic.NewMagic()
	
	result, err := m.Execute("unknown")
	
	// Check that there was no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check that the result contains the expected text
	if !strings.Contains(result, "Unknown magic command") {
		t.Errorf("Expected result to contain 'Unknown magic command', got: %s", result)
	}
	
	// Check that the result contains available commands
	if !strings.Contains(result, "Available magic commands") {
		t.Errorf("Expected result to contain 'Available magic commands', got: %s", result)
	}
}

// TestExecuteDance tests the dance command
func TestExecuteDance(t *testing.T) {
	m := magic.NewMagic()
	
	result, err := m.Execute("dance")
	
	// Check that there was no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check that the result is not empty
	if result == "" {
		t.Error("Expected non-empty result for dance command")
	}
	
	// Check that the result contains dance-related text
	// We can't check for specific animations since they're randomly selected
	dancePhrases := []string{"dance", "Dance", "DANCE", "🎵", "🕺", "💃", "🦜", "🪩", "🤖"}
	foundDancePhrase := false
	
	for _, phrase := range dancePhrases {
		if strings.Contains(result, phrase) {
			foundDancePhrase = true
			break
		}
	}
	
	if !foundDancePhrase {
		t.Errorf("Expected result to contain dance-related text, got: %s", result)
	}
}

// TestCommandCaseInsensitivity tests that commands are case-insensitive
func TestCommandCaseInsensitivity(t *testing.T) {
	m := magic.NewMagic()
	
	// Test with different cases
	commands := []string{"dance", "DANCE", "Dance", "dAnCe"}
	
	for _, cmd := range commands {
		result, err := m.Execute(cmd)
		
		// Check that there was no error
		if err != nil {
			t.Fatalf("Expected no error for command '%s', got %v", cmd, err)
		}
		
		// Check that the result is not empty
		if result == "" {
			t.Errorf("Expected non-empty result for command '%s'", cmd)
		}
		
		// Check that the result contains dance-related text
		dancePhrases := []string{"dance", "Dance", "DANCE", "🎵", "🕺", "💃", "🦜", "🪩", "🤖"}
		foundDancePhrase := false
		
		for _, phrase := range dancePhrases {
			if strings.Contains(result, phrase) {
				foundDancePhrase = true
				break
			}
		}
		
		if !foundDancePhrase {
			t.Errorf("Expected result to contain dance-related text for command '%s', got: %s", cmd, result)
		}
	}
}

// TestCommandWithWhitespace tests that commands with whitespace are handled correctly
func TestCommandWithWhitespace(t *testing.T) {
	m := magic.NewMagic()
	
	// Test with whitespace
	commands := []string{" dance", "dance ", " dance "}
	
	for _, cmd := range commands {
		result, err := m.Execute(cmd)
		
		// Check that there was no error
		if err != nil {
			t.Fatalf("Expected no error for command '%s', got %v", cmd, err)
		}
		
		// Check that the result is not empty
		if result == "" {
			t.Errorf("Expected non-empty result for command '%s'", cmd)
		}
		
		// Check that the result contains dance-related text
		dancePhrases := []string{"dance", "Dance", "DANCE", "🎵", "🕺", "💃", "🦜", "🪩", "🤖"}
		foundDancePhrase := false
		
		for _, phrase := range dancePhrases {
			if strings.Contains(result, phrase) {
				foundDancePhrase = true
				break
			}
		}
		
		if !foundDancePhrase {
			t.Errorf("Expected result to contain dance-related text for command '%s', got: %s", cmd, result)
		}
	}
}
