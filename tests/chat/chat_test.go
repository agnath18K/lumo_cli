package chat_test

import (
	"context"
	"testing"

	"github.com/agnath18/lumo/pkg/chat"
)

// TestConversation tests the basic conversation functionality
func TestConversation(t *testing.T) {
	// Create a new conversation
	conv := chat.NewConversation("Test system message", 10)

	// Check that the conversation was created
	if conv == nil {
		t.Fatal("Expected conversation to be created, got nil")
	}

	// Add messages to the conversation
	conv.AddUserMessage("Hello")
	conv.AddAssistantMessage("Hi there!")
	conv.AddUserMessage("How are you?")
	conv.AddAssistantMessage("I'm doing well, thanks for asking!")

	// Check the number of messages
	messages := conv.GetMessages()
	if len(messages) != 5 { // 1 system + 2 user + 2 assistant
		t.Errorf("Expected 5 messages, got %d", len(messages))
	}

	// Check the last user message
	lastUserMsg, found := conv.GetLastUserMessage()
	if !found {
		t.Fatal("Expected to find last user message")
	}
	if lastUserMsg.Content != "How are you?" {
		t.Errorf("Expected last user message to be 'How are you?', got '%s'", lastUserMsg.Content)
	}

	// Check the last assistant message
	lastAssistantMsg, found := conv.GetLastAssistantMessage()
	if !found {
		t.Fatal("Expected to find last assistant message")
	}
	if lastAssistantMsg.Content != "I'm doing well, thanks for asking!" {
		t.Errorf("Expected last assistant message to be 'I'm doing well, thanks for asking!', got '%s'", lastAssistantMsg.Content)
	}

	// Test clearing the conversation
	conv.Clear()
	messages = conv.GetMessages()
	if len(messages) != 1 { // Only system message should remain
		t.Errorf("Expected 1 message after clearing, got %d", len(messages))
	}
	if messages[0].Role != chat.RoleSystem {
		t.Errorf("Expected remaining message to be system message, got %s", messages[0].Role)
	}
}

// TestConversationTrimming tests that conversations are trimmed when they exceed the maximum size
func TestConversationTrimming(t *testing.T) {
	// Create a new conversation with a small maximum size
	maxSize := 5
	conv := chat.NewConversation("Test system message", maxSize)

	// Add more messages than the maximum size
	for i := 0; i < 10; i++ {
		conv.AddUserMessage("User message")
		conv.AddAssistantMessage("Assistant message")
	}

	// Check that the conversation was trimmed
	messages := conv.GetMessages()
	if len(messages) > maxSize {
		t.Errorf("Expected at most %d messages, got %d", maxSize, len(messages))
	}
}

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

func (m *MockAIClient) ProcessChatMessage(ctx context.Context, conversation string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

// TestChatManager tests the chat manager functionality
func TestChatManager(t *testing.T) {
	// Create a mock AI client
	mockClient := &MockAIClient{
		response: "This is a test response",
	}

	// Create a chat manager
	manager := chat.NewManager(mockClient, 3, 10)

	// Check that the manager was created
	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}

	// Start a new conversation
	conv := manager.StartNewConversation()
	if conv == nil {
		t.Fatal("Expected conversation to be created, got nil")
	}

	// Get the active conversation
	activeConv := manager.GetActiveConversation()
	if activeConv == nil {
		t.Fatal("Expected active conversation, got nil")
	}
	if activeConv.ID != conv.ID {
		t.Errorf("Expected active conversation ID to be %s, got %s", conv.ID, activeConv.ID)
	}

	// Process a message
	ctx := context.Background()
	response, err := manager.ProcessMessage(ctx, "Hello")
	if err != nil {
		t.Fatalf("Error processing message: %v", err)
	}
	if response != "This is a test response" {
		t.Errorf("Expected response to be 'This is a test response', got '%s'", response)
	}

	// Check that the message was added to the conversation
	messages := activeConv.GetMessages()
	if len(messages) != 3 { // 1 system + 1 user + 1 assistant
		t.Errorf("Expected 3 messages, got %d", len(messages))
	}

	// Start more conversations than the maximum
	for i := 0; i < 5; i++ {
		manager.StartNewConversation()
	}

	// Check that the number of conversations is limited
	convs := manager.ListConversations()
	if len(convs) > 3 {
		t.Errorf("Expected at most 3 conversations, got %d", len(convs))
	}

	// Delete a conversation
	if !manager.DeleteConversation(convs[0]) {
		t.Errorf("Failed to delete conversation %s", convs[0])
	}

	// Check that the conversation was deleted
	newConvs := manager.ListConversations()
	if len(newConvs) != len(convs)-1 {
		t.Errorf("Expected %d conversations after deletion, got %d", len(convs)-1, len(newConvs))
	}

	// Clear all conversations
	manager.ClearAllConversations()
	finalConvs := manager.ListConversations()
	if len(finalConvs) != 0 {
		t.Errorf("Expected 0 conversations after clearing, got %d", len(finalConvs))
	}
}
