package ai

import (
	"context"
)

// Client defines the interface for AI service clients
type Client interface {
	// Query sends a query to the AI service and returns the response
	Query(query string) (string, error)

	// GetCompletion sends a prompt to the AI service and returns the completion
	// This is used by the agent to generate plans
	GetCompletion(ctx context.Context, prompt string) (string, error)
}

// ChatClient extends the Client interface with chat-specific methods
type ChatClient interface {
	Client

	// ProcessChatMessage processes a chat message with conversation history
	// and returns the AI response
	ProcessChatMessage(ctx context.Context, conversation string) (string, error)
}

// Provider represents the type of AI provider
type Provider string

const (
	// ProviderGemini represents Google's Gemini AI
	ProviderGemini Provider = "gemini"
	// ProviderOpenAI represents OpenAI's GPT
	ProviderOpenAI Provider = "openai"
)
