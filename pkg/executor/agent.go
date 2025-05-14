package executor

import (
	"context"
)

// AgentInterface defines the interface for agent implementations
type AgentInterface interface {
	// Execute processes a task and executes the necessary commands
	Execute(ctx context.Context, taskDescription string) (*Result, error)
}
