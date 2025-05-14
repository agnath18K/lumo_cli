package pipe_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agnath18/lumo/pkg/pipe"
)

// MockAIClient is a mock implementation of the ai.Client interface for testing
type MockAIClient struct {
	response string
	err      error
}

func (m *MockAIClient) Query(prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func (m *MockAIClient) GetCompletion(ctx context.Context, prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

// TestProcessInput tests the ProcessInput function
func TestProcessInput(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		mockResponse   string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "simple text",
			input:          "Hello, world!",
			mockResponse:   "This is a simple greeting message.",
			expectedOutput: "This is a simple greeting message.",
			expectError:    false,
		},
		{
			name: "code snippet",
			input: `
func main() {
    fmt.Println("Hello, world!")
}
`,
			mockResponse:   "This is a Go code snippet that prints 'Hello, world!'.",
			expectedOutput: "This is a Go code snippet that prints 'Hello, world!'.",
			expectError:    false,
		},
		{
			name: "command output",
			input: `
file1.txt
file2.txt
file3.txt
`,
			mockResponse:   "This appears to be a list of files.",
			expectedOutput: "This appears to be a list of files.",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock AI client
			mockClient := &MockAIClient{
				response: tt.mockResponse,
			}

			// Create a pipe processor with the mock client
			processor := pipe.NewProcessor(mockClient)

			// Create a buffer with the test input
			inputBuffer := bytes.NewBufferString(tt.input)

			// Process the input
			result, err := processor.ProcessInput(inputBuffer)

			// Check for errors
			if tt.expectError && err == nil {
				t.Error("Expected an error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check the result
			if result != tt.expectedOutput {
				t.Errorf("Expected output '%s', got '%s'", tt.expectedOutput, result)
			}
		})
	}
}

// TestProcessLargeInput tests processing a large input
func TestProcessLargeInput(t *testing.T) {
	// Create a large input (repeated text)
	var inputBuilder strings.Builder
	for i := 0; i < 1000; i++ {
		inputBuilder.WriteString("This is line " + string(rune('A'+i%26)) + "\n")
	}
	largeInput := inputBuilder.String()

	// Create a mock AI client
	mockClient := &MockAIClient{
		response: "This is a large text file with alphabetical lines.",
	}

	// Create a pipe processor with the mock client
	processor := pipe.NewProcessor(mockClient)

	// Create a buffer with the large input
	inputBuffer := bytes.NewBufferString(largeInput)

	// Process the input
	result, err := processor.ProcessInput(inputBuffer)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the result
	expectedOutput := "This is a large text file with alphabetical lines."
	if result != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, result)
	}
}

// TestProcessEmptyInput tests processing an empty input
func TestProcessEmptyInput(t *testing.T) {
	// Create a mock AI client
	mockClient := &MockAIClient{
		response: "This is an empty input.",
	}

	// Create a pipe processor with the mock client
	processor := pipe.NewProcessor(mockClient)

	// Create a buffer with empty input
	inputBuffer := bytes.NewBufferString("")

	// Process the input
	_, err := processor.ProcessInput(inputBuffer)

	// Check for errors - we expect an error for empty input
	if err == nil {
		t.Error("Expected an error for empty input but got nil")
	}

	// Check that the error message is correct
	expectedErrMsg := "empty input"
	if err != nil && !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestProcessInputWithError tests handling errors from the AI client
func TestProcessInputWithError(t *testing.T) {
	// Create a mock AI client that returns an error
	mockClient := &MockAIClient{
		err: fmt.Errorf("API request failed"),
	}

	// Create a pipe processor with the mock client
	processor := pipe.NewProcessor(mockClient)

	// Create a buffer with test input
	inputBuffer := bytes.NewBufferString("Test input")

	// Process the input
	_, err := processor.ProcessInput(inputBuffer)

	// Check for errors
	if err == nil {
		t.Error("Expected an error but got nil")
	}
}
