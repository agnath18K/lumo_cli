package chat

import (
	"fmt"
	"time"
)

// MessageRole represents the role of a message in a conversation
type MessageRole string

const (
	// RoleSystem represents a system message
	RoleSystem MessageRole = "system"
	// RoleUser represents a user message
	RoleUser MessageRole = "user"
	// RoleAssistant represents an assistant message
	RoleAssistant MessageRole = "assistant"
)

// Message represents a single message in a conversation
type Message struct {
	Role      MessageRole
	Content   string
	Timestamp time.Time
}

// Conversation represents a chat conversation with history
type Conversation struct {
	ID       string
	Messages []Message
	MaxSize  int
}

// NewConversation creates a new conversation with the given system message
func NewConversation(systemMessage string, maxSize int) *Conversation {
	// If maxSize is not specified, use a default value
	if maxSize <= 0 {
		maxSize = 20
	}

	// Create a new conversation with a system message
	conversation := &Conversation{
		ID:       fmt.Sprintf("conv_%d", time.Now().UnixNano()),
		Messages: []Message{},
		MaxSize:  maxSize,
	}

	// Add the system message if provided
	if systemMessage != "" {
		conversation.AddSystemMessage(systemMessage)
	}

	return conversation
}

// AddSystemMessage adds a system message to the conversation
func (c *Conversation) AddSystemMessage(content string) {
	c.Messages = append(c.Messages, Message{
		Role:      RoleSystem,
		Content:   content,
		Timestamp: time.Now(),
	})

	// Trim the conversation if it exceeds the maximum size
	c.trimIfNeeded()
}

// AddUserMessage adds a user message to the conversation
func (c *Conversation) AddUserMessage(content string) {
	c.Messages = append(c.Messages, Message{
		Role:      RoleUser,
		Content:   content,
		Timestamp: time.Now(),
	})

	// Trim the conversation if it exceeds the maximum size
	c.trimIfNeeded()
}

// AddAssistantMessage adds an assistant message to the conversation
func (c *Conversation) AddAssistantMessage(content string) {
	c.Messages = append(c.Messages, Message{
		Role:      RoleAssistant,
		Content:   content,
		Timestamp: time.Now(),
	})

	// Trim the conversation if it exceeds the maximum size
	c.trimIfNeeded()
}

// GetMessages returns all messages in the conversation
func (c *Conversation) GetMessages() []Message {
	return c.Messages
}

// GetLastUserMessage returns the last user message in the conversation
func (c *Conversation) GetLastUserMessage() (Message, bool) {
	// Iterate through messages in reverse order
	for i := len(c.Messages) - 1; i >= 0; i-- {
		if c.Messages[i].Role == RoleUser {
			return c.Messages[i], true
		}
	}

	// No user message found
	return Message{}, false
}

// GetLastAssistantMessage returns the last assistant message in the conversation
func (c *Conversation) GetLastAssistantMessage() (Message, bool) {
	// Iterate through messages in reverse order
	for i := len(c.Messages) - 1; i >= 0; i-- {
		if c.Messages[i].Role == RoleAssistant {
			return c.Messages[i], true
		}
	}

	// No assistant message found
	return Message{}, false
}

// Clear clears all messages in the conversation except for system messages
func (c *Conversation) Clear() {
	// Keep only system messages
	var systemMessages []Message
	for _, msg := range c.Messages {
		if msg.Role == RoleSystem {
			systemMessages = append(systemMessages, msg)
		}
	}

	c.Messages = systemMessages
}

// trimIfNeeded trims the conversation if it exceeds the maximum size
// It keeps the system messages and the most recent messages
func (c *Conversation) trimIfNeeded() {
	// If the conversation is within the limit, do nothing
	if len(c.Messages) <= c.MaxSize {
		return
	}

	// Separate system messages from other messages
	var systemMessages []Message
	var otherMessages []Message

	for _, msg := range c.Messages {
		if msg.Role == RoleSystem {
			systemMessages = append(systemMessages, msg)
		} else {
			otherMessages = append(otherMessages, msg)
		}
	}

	// Calculate how many non-system messages to keep
	keepCount := c.MaxSize - len(systemMessages)
	if keepCount < 0 {
		keepCount = 0
	}

	// Keep only the most recent non-system messages
	if len(otherMessages) > keepCount {
		otherMessages = otherMessages[len(otherMessages)-keepCount:]
	}

	// Combine system messages and kept non-system messages
	c.Messages = append(systemMessages, otherMessages...)
}
