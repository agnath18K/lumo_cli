package chat_test

import (
	"testing"

	"github.com/agnath18/lumo/pkg/chat"
	"github.com/agnath18/lumo/pkg/config"
)

// MockREPLReader is a mock implementation of the bufio.Reader for testing
type MockREPLReader struct {
	inputs []string
	index  int
}

// NewMockREPLReader creates a new mock reader with the given inputs
func NewMockREPLReader(inputs []string) *MockREPLReader {
	return &MockREPLReader{
		inputs: inputs,
		index:  0,
	}
}

// ReadString simulates reading a string from the input
func (m *MockREPLReader) ReadString(delim byte) (string, error) {
	if m.index >= len(m.inputs) {
		return "", nil
	}

	input := m.inputs[m.index]
	m.index++
	return input + string(delim), nil
}

// TestREPL tests the REPL functionality
func TestREPL(t *testing.T) {
	// Create a mock AI client
	mockClient := &MockAIClient{
		response: "This is a test response",
	}

	// Create a config
	cfg := config.DefaultConfig()
	cfg.EnableChatREPL = true

	// Create a chat manager
	manager := chat.NewManager(mockClient, 3, 10)

	// Create a REPL with a mock reader
	repl := chat.NewREPL(cfg, manager, mockClient)

	// Test the REPL functionality
	// This is a basic test since we can't fully test the interactive REPL
	// without mocking the reader and writer
	if repl == nil {
		t.Fatal("Expected REPL to be created, got nil")
	}
}
