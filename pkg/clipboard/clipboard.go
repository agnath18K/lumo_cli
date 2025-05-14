package clipboard

import (
	"fmt"
	"io"
	"strings"

	"github.com/atotto/clipboard"
)

// ClipboardProvider defines the interface for clipboard operations
type ClipboardProvider interface {
	ReadAll() (string, error)
	WriteAll(text string) error
}

// DefaultClipboardProvider is the default implementation using the clipboard package
type DefaultClipboardProvider struct{}

// ReadAll reads the clipboard content
func (p *DefaultClipboardProvider) ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll writes text to the clipboard
func (p *DefaultClipboardProvider) WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// Clipboard handles clipboard operations
type Clipboard struct {
	provider ClipboardProvider
}

// NewClipboard creates a new Clipboard instance with the default provider
func NewClipboard() *Clipboard {
	return &Clipboard{
		provider: &DefaultClipboardProvider{},
	}
}

// NewClipboardWithProvider creates a new Clipboard instance with a custom provider
func NewClipboardWithProvider(provider ClipboardProvider) *Clipboard {
	return &Clipboard{
		provider: provider,
	}
}

// Execute processes a clipboard command and returns the result
func (c *Clipboard) Execute(command string, reader io.Reader) (string, error) {
	// Check for special commands
	if command == "clear" {
		return c.ClearContent()
	}

	// Check for append mode
	isAppend := false
	if strings.HasPrefix(command, "append ") {
		isAppend = true
		command = strings.TrimPrefix(command, "append ")
	}

	// If we have a reader (piped input), use that content
	if reader != nil {
		content, err := readAllInput(reader)
		if err != nil {
			return "", fmt.Errorf("failed to read piped input: %w", err)
		}
		if isAppend {
			return c.AppendContent(content)
		}
		return c.SetContent(content)
	}

	// If command is empty, get clipboard content
	if command == "" {
		return c.GetContent()
	}

	// Otherwise, use the command as content
	if isAppend {
		return c.AppendContent(command)
	}
	return c.SetContent(command)
}

// GetContent retrieves the current clipboard content
func (c *Clipboard) GetContent() (string, error) {
	content, err := c.provider.ReadAll()
	if err != nil {
		// Check if the error is due to missing clipboard utilities
		if strings.Contains(err.Error(), "No clipboard utilities available") {
			return "", fmt.Errorf("clipboard utilities not available. Please install xsel, xclip, wl-clipboard, or Termux:API")
		}
		return "", fmt.Errorf("failed to read clipboard: %w", err)
	}

	if content == "" {
		return "Clipboard is empty", nil
	}

	return content, nil
}

// SetContent sets the clipboard content
func (c *Clipboard) SetContent(content string) (string, error) {
	err := c.provider.WriteAll(content)
	if err != nil {
		// Check if the error is due to missing clipboard utilities
		if strings.Contains(err.Error(), "No clipboard utilities available") {
			return "", fmt.Errorf("clipboard utilities not available. Please install xsel, xclip, wl-clipboard, or Termux:API")
		}
		return "", fmt.Errorf("failed to write to clipboard: %w", err)
	}

	return fmt.Sprintf("Copied to clipboard: %s", truncateForDisplay(content)), nil
}

// AppendContent appends text to the existing clipboard content
func (c *Clipboard) AppendContent(content string) (string, error) {
	// First, get the current content
	currentContent, err := c.provider.ReadAll()
	if err != nil {
		// Check if the error is due to missing clipboard utilities
		if strings.Contains(err.Error(), "No clipboard utilities available") {
			return "", fmt.Errorf("clipboard utilities not available. Please install xsel, xclip, wl-clipboard, or Termux:API")
		}
		return "", fmt.Errorf("failed to read clipboard: %w", err)
	}

	// Append the new content
	newContent := currentContent
	if currentContent != "" {
		newContent += "\n" + content
	} else {
		newContent = content
	}

	// Write the combined content back to the clipboard
	err = c.provider.WriteAll(newContent)
	if err != nil {
		return "", fmt.Errorf("failed to write to clipboard: %w", err)
	}

	return fmt.Sprintf("Appended to clipboard: %s", truncateForDisplay(content)), nil
}

// ClearContent clears the clipboard content
func (c *Clipboard) ClearContent() (string, error) {
	err := c.provider.WriteAll("")
	if err != nil {
		// Check if the error is due to missing clipboard utilities
		if strings.Contains(err.Error(), "No clipboard utilities available") {
			return "", fmt.Errorf("clipboard utilities not available. Please install xsel, xclip, wl-clipboard, or Termux:API")
		}
		return "", fmt.Errorf("failed to clear clipboard: %w", err)
	}

	return "Clipboard cleared", nil
}

// truncateForDisplay truncates a string for display purposes
func truncateForDisplay(s string) string {
	// Remove newlines for display
	s = strings.ReplaceAll(s, "\n", " ")

	// Truncate if too long
	maxLen := 50
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// readAllInput reads all input from a reader
func readAllInput(reader io.Reader) (string, error) {
	// Read all content from the reader
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
