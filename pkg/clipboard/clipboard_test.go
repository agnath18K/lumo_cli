package clipboard

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// MockClipboardProvider is a mock implementation of the ClipboardProvider interface
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
var ErrUnavailable = errors.New("No clipboard utilities available")

func TestClipboard_Execute_GetContent(t *testing.T) {
	// Create a mock provider with test content
	mockProvider := &MockClipboardProvider{content: "test content"}

	// Create a clipboard with the mock provider
	c := NewClipboardWithProvider(mockProvider)

	// Test getting content
	result, err := c.Execute("", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "test content" {
		t.Errorf("Expected 'test content', got '%s'", result)
	}
}

func TestClipboard_Execute_SetContent(t *testing.T) {
	// Create a mock provider
	mockProvider := &MockClipboardProvider{}

	// Create a clipboard with the mock provider
	c := NewClipboardWithProvider(mockProvider)

	// Test setting content
	result, err := c.Execute("new content", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(result, "Copied to clipboard: new content") {
		t.Errorf("Expected result to contain 'Copied to clipboard: new content', got '%s'", result)
	}
	if mockProvider.content != "new content" {
		t.Errorf("Expected clipboard content to be 'new content', got '%s'", mockProvider.content)
	}
}

func TestClipboard_Execute_PipedInput(t *testing.T) {
	// Create a mock provider
	mockProvider := &MockClipboardProvider{}

	// Create a clipboard with the mock provider
	c := NewClipboardWithProvider(mockProvider)

	// Create a reader with piped input
	reader := bytes.NewBufferString("piped content")

	// Test setting content from piped input
	result, err := c.Execute("", reader)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(result, "Copied to clipboard: piped content") {
		t.Errorf("Expected result to contain 'Copied to clipboard: piped content', got '%s'", result)
	}
	if mockProvider.content != "piped content" {
		t.Errorf("Expected clipboard content to be 'piped content', got '%s'", mockProvider.content)
	}
}

func TestClipboard_Execute_AppendContent(t *testing.T) {
	// Create a mock provider with initial content
	mockProvider := &MockClipboardProvider{content: "initial content"}

	// Create a clipboard with the mock provider
	c := NewClipboardWithProvider(mockProvider)

	// Test appending content
	result, err := c.Execute("append additional content", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(result, "Appended to clipboard: additional content") {
		t.Errorf("Expected result to contain 'Appended to clipboard: additional content', got '%s'", result)
	}

	// Check that the content was appended correctly
	expected := "initial content\nadditional content"
	if mockProvider.content != expected {
		t.Errorf("Expected clipboard content to be '%s', got '%s'", expected, mockProvider.content)
	}

	// Test appending to empty clipboard
	mockProvider.content = ""
	result, err = c.Execute("append new content", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if mockProvider.content != "new content" {
		t.Errorf("Expected clipboard content to be 'new content', got '%s'", mockProvider.content)
	}
}

func TestClipboard_Execute_ClearContent(t *testing.T) {
	// Create a mock provider with initial content
	mockProvider := &MockClipboardProvider{content: "content to clear"}

	// Create a clipboard with the mock provider
	c := NewClipboardWithProvider(mockProvider)

	// Test clearing content
	result, err := c.Execute("clear", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "Clipboard cleared" {
		t.Errorf("Expected 'Clipboard cleared', got '%s'", result)
	}

	// Check that the content was cleared
	if mockProvider.content != "" {
		t.Errorf("Expected empty clipboard content, got '%s'", mockProvider.content)
	}
}

func TestClipboard_Execute_Error(t *testing.T) {
	// Create a mock provider with an error
	mockProvider := &MockClipboardProvider{err: ErrUnavailable}

	// Create a clipboard with the mock provider
	c := NewClipboardWithProvider(mockProvider)

	// Test getting content with an error
	_, err := c.Execute("", nil)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "clipboard utilities not available") {
		t.Errorf("Expected error to contain 'clipboard utilities not available', got '%s'", err.Error())
	}

	// Test setting content with an error
	_, err = c.Execute("new content", nil)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "clipboard utilities not available") {
		t.Errorf("Expected error to contain 'clipboard utilities not available', got '%s'", err.Error())
	}

	// Test appending content with an error
	_, err = c.Execute("append content", nil)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "clipboard utilities not available") {
		t.Errorf("Expected error to contain 'clipboard utilities not available', got '%s'", err.Error())
	}

	// Test clearing content with an error
	_, err = c.Execute("clear", nil)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "clipboard utilities not available") {
		t.Errorf("Expected error to contain 'clipboard utilities not available', got '%s'", err.Error())
	}
}
