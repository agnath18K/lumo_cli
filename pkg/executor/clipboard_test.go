package executor

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/agnath18/lumo/pkg/clipboard"
	"github.com/agnath18/lumo/pkg/config"
	"github.com/agnath18/lumo/pkg/nlp"
)

// MockClipboardProvider is a mock implementation of the clipboard.ClipboardProvider interface
type MockClipboardProvider struct {
	content string
	err     error
}

// ReadAll is a mock implementation of ClipboardProvider.ReadAll
func (m *MockClipboardProvider) ReadAll() (string, error) {
	return m.content, m.err
}

// WriteAll is a mock implementation of ClipboardProvider.WriteAll
func (m *MockClipboardProvider) WriteAll(text string) error {
	if m.err != nil {
		return m.err
	}
	m.content = text
	return nil
}

// ErrUnavailable is a mock error for clipboard unavailability
var ErrUnavailable = errors.New("no clipboard utilities available")

func TestExecutor_ExecuteClipboardCommand(t *testing.T) {
	// Create a mock clipboard provider
	mockProvider := &MockClipboardProvider{content: "test clipboard content"}

	// Create a clipboard with the mock provider
	clipboardInstance := clipboard.NewClipboardWithProvider(mockProvider)

	// Create a new executor
	cfg := &config.Config{}
	exec := &Executor{
		config:    cfg,
		clipboard: clipboardInstance,
	}

	// Test getting clipboard content
	cmd := &nlp.Command{
		Type:       nlp.CommandTypeClipboard,
		Intent:     "",
		Parameters: make(map[string]string),
		RawInput:   "clipboard",
	}

	// Execute the command
	result, err := exec.ExecuteWithReader(cmd, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Output != "test clipboard content" {
		t.Errorf("Expected 'test clipboard content', got '%s'", result.Output)
	}
	if result.IsError {
		t.Errorf("Expected IsError to be false, got true")
	}

	// Test setting clipboard content
	cmd = &nlp.Command{
		Type:       nlp.CommandTypeClipboard,
		Intent:     "new clipboard content",
		Parameters: make(map[string]string),
		RawInput:   "clipboard new clipboard content",
	}

	// Execute the command
	result, err = exec.ExecuteWithReader(cmd, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(result.Output, "Copied to clipboard: new clipboard content") {
		t.Errorf("Expected output to contain 'Copied to clipboard: new clipboard content', got '%s'", result.Output)
	}
	if result.IsError {
		t.Errorf("Expected IsError to be false, got true")
	}
	if mockProvider.content != "new clipboard content" {
		t.Errorf("Expected clipboard content to be 'new clipboard content', got '%s'", mockProvider.content)
	}

	// Test setting clipboard content from piped input
	cmd = &nlp.Command{
		Type:       nlp.CommandTypeClipboard,
		Intent:     "",
		Parameters: make(map[string]string),
		RawInput:   "clipboard",
	}

	// Create a reader with piped input
	reader := bytes.NewBufferString("piped content")

	// Execute the command with piped input
	result, err = exec.ExecuteWithReader(cmd, reader)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(result.Output, "Copied to clipboard: piped content") {
		t.Errorf("Expected output to contain 'Copied to clipboard: piped content', got '%s'", result.Output)
	}
	if result.IsError {
		t.Errorf("Expected IsError to be false, got true")
	}
	if mockProvider.content != "piped content" {
		t.Errorf("Expected clipboard content to be 'piped content', got '%s'", mockProvider.content)
	}

	// Test appending to clipboard
	cmd = &nlp.Command{
		Type:       nlp.CommandTypeClipboard,
		Intent:     "append more text",
		Parameters: make(map[string]string),
		RawInput:   "clipboard append more text",
	}

	// Execute the command
	result, err = exec.ExecuteWithReader(cmd, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(result.Output, "Appended to clipboard: more text") {
		t.Errorf("Expected output to contain 'Appended to clipboard: more text', got '%s'", result.Output)
	}
	if result.IsError {
		t.Errorf("Expected IsError to be false, got true")
	}
	expected := "piped content\nmore text"
	if mockProvider.content != expected {
		t.Errorf("Expected clipboard content to be '%s', got '%s'", expected, mockProvider.content)
	}

	// Test appending from piped input
	cmd = &nlp.Command{
		Type:       nlp.CommandTypeClipboard,
		Intent:     "append ",
		Parameters: make(map[string]string),
		RawInput:   "clipboard append",
	}

	// Create a reader with piped input for append
	appendReader := bytes.NewBufferString("piped append")

	// Execute the command with piped input
	result, err = exec.ExecuteWithReader(cmd, appendReader)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(result.Output, "Appended to clipboard: piped append") {
		t.Errorf("Expected output to contain 'Appended to clipboard: piped append', got '%s'", result.Output)
	}
	if result.IsError {
		t.Errorf("Expected IsError to be false, got true")
	}
	expected = "piped content\nmore text\npiped append"
	if mockProvider.content != expected {
		t.Errorf("Expected clipboard content to be '%s', got '%s'", expected, mockProvider.content)
	}

	// Test clearing clipboard
	cmd = &nlp.Command{
		Type:       nlp.CommandTypeClipboard,
		Intent:     "clear",
		Parameters: make(map[string]string),
		RawInput:   "clipboard clear",
	}

	// Execute the command
	result, err = exec.ExecuteWithReader(cmd, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Output != "Clipboard cleared" {
		t.Errorf("Expected 'Clipboard cleared', got '%s'", result.Output)
	}
	if result.IsError {
		t.Errorf("Expected IsError to be false, got true")
	}
	if mockProvider.content != "" {
		t.Errorf("Expected empty clipboard content, got '%s'", mockProvider.content)
	}
}
