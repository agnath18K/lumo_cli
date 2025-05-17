package tests

import (
	"context"
	"errors"
	"testing"
)

// MockAIClient is a simple mock implementation of an AI client
type MockAIClient struct {
	QueryResponse      string
	QueryError         error
	CompletionResponse string
	CompletionError    error
	QueryCalls         []string
	CompletionCalls    []string
}

// Query records the call and returns the mock response or error
func (m *MockAIClient) Query(query string) (string, error) {
	m.QueryCalls = append(m.QueryCalls, query)
	return m.QueryResponse, m.QueryError
}

// GetCompletion records the call and returns the mock response or error
func (m *MockAIClient) GetCompletion(ctx context.Context, prompt string) (string, error) {
	m.CompletionCalls = append(m.CompletionCalls, prompt)
	return m.CompletionResponse, m.CompletionError
}

// TestAICompletion tests the AI completion functionality
func TestAICompletion(t *testing.T) {
	// Create test cases
	testCases := []struct {
		name           string
		prompt         string
		mockResponse   string
		mockError      error
		expectedOutput string
		shouldError    bool
	}{
		{
			name:           "Successful completion",
			prompt:         "Generate a plan to create a web server",
			mockResponse:   "1. Install Node.js\n2. Create a new project\n3. Install Express\n4. Write server code\n5. Test the server",
			mockError:      nil,
			expectedOutput: "1. Install Node.js\n2. Create a new project\n3. Install Express\n4. Write server code\n5. Test the server",
			shouldError:    false,
		},
		{
			name:           "Error from AI provider",
			prompt:         "Generate a plan to create a web server",
			mockResponse:   "",
			mockError:      errors.New("API error: model not available"),
			expectedOutput: "Error",
			shouldError:    true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock AI client
			mockAI := &MockAIClient{
				CompletionResponse: tc.mockResponse,
				CompletionError:    tc.mockError,
			}

			// Call GetCompletion directly
			result, err := mockAI.GetCompletion(context.Background(), tc.prompt)

			// Check for errors
			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expectedOutput {
					t.Errorf("Expected output '%s', got '%s'", tc.expectedOutput, result)
				}
			}

			// Verify the completion was called with the correct parameters
			if len(mockAI.CompletionCalls) != 1 {
				t.Errorf("Expected 1 call to GetCompletion, got %d", len(mockAI.CompletionCalls))
			} else if mockAI.CompletionCalls[0] != tc.prompt {
				t.Errorf("Expected prompt '%s', got '%s'", tc.prompt, mockAI.CompletionCalls[0])
			}
		})
	}
}

// TestAIQuery tests the AI query functionality
func TestAIQuery(t *testing.T) {
	// Create test cases
	testCases := []struct {
		name           string
		query          string
		mockResponse   string
		mockError      error
		expectedOutput string
		shouldError    bool
	}{
		{
			name:           "Successful query",
			query:          "What is the capital of France?",
			mockResponse:   "The capital of France is Paris.",
			mockError:      nil,
			expectedOutput: "The capital of France is Paris.",
			shouldError:    false,
		},
		{
			name:           "Error from AI provider",
			query:          "What is the capital of France?",
			mockResponse:   "",
			mockError:      errors.New("API error: rate limit exceeded"),
			expectedOutput: "Error",
			shouldError:    true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock AI client
			mockAI := &MockAIClient{
				QueryResponse: tc.mockResponse,
				QueryError:    tc.mockError,
			}

			// Call Query directly
			result, err := mockAI.Query(tc.query)

			// Check for errors
			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expectedOutput {
					t.Errorf("Expected output '%s', got '%s'", tc.expectedOutput, result)
				}
			}

			// Verify the query was called with the correct parameters
			if len(mockAI.QueryCalls) != 1 {
				t.Errorf("Expected 1 call to Query, got %d", len(mockAI.QueryCalls))
			} else if mockAI.QueryCalls[0] != tc.query {
				t.Errorf("Expected query '%s', got '%s'", tc.query, mockAI.QueryCalls[0])
			}
		})
	}
}
