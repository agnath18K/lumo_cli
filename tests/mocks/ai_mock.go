package mocks

import (
	"context"
	"errors"
)

// MockAIClient is a comprehensive mock implementation of the ai.Client interface
type MockAIClient struct {
	QueryResponse        string
	QueryError           error
	CompletionResponse   string
	CompletionError      error
	QueryCalls           []string
	CompletionCalls      []string
	ProcessChatResponse  string
	ProcessChatError     error
	ProcessChatCalls     []string
	ListModelsResponse   []string
	ListModelsError      error
	ShouldFailWithStatus int
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

// ProcessChatMessage records the call and returns the mock response or error
// This is for clients that implement the ChatClient interface
func (m *MockAIClient) ProcessChatMessage(ctx context.Context, conversation string) (string, error) {
	m.ProcessChatCalls = append(m.ProcessChatCalls, conversation)
	return m.ProcessChatResponse, m.ProcessChatError
}

// ListModels returns a mock list of models
// This is for clients that implement model listing (like Ollama)
func (m *MockAIClient) ListModels() ([]string, error) {
	return m.ListModelsResponse, m.ListModelsError
}

// NewMockAIClient creates a new mock AI client with default success responses
func NewMockAIClient() *MockAIClient {
	return &MockAIClient{
		QueryResponse:       "This is a mock AI response",
		CompletionResponse:  "This is a mock completion response",
		ProcessChatResponse: "This is a mock chat response",
		ListModelsResponse:  []string{"model1", "model2", "model3"},
	}
}

// NewMockAIClientWithError creates a new mock AI client that returns errors
func NewMockAIClientWithError(errMsg string) *MockAIClient {
	err := errors.New(errMsg)
	return &MockAIClient{
		QueryError:       err,
		CompletionError:  err,
		ProcessChatError: err,
		ListModelsError:  err,
	}
}

// NewMockAIClientWithCustomResponses creates a new mock AI client with custom responses
func NewMockAIClientWithCustomResponses(queryResp, completionResp, chatResp string) *MockAIClient {
	return &MockAIClient{
		QueryResponse:       queryResp,
		CompletionResponse:  completionResp,
		ProcessChatResponse: chatResp,
		ListModelsResponse:  []string{"model1", "model2", "model3"},
	}
}

// Reset clears the recorded calls
func (m *MockAIClient) Reset() {
	m.QueryCalls = nil
	m.CompletionCalls = nil
	m.ProcessChatCalls = nil
}
