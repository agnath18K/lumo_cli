package assistant

import (
	"context"
	"fmt"

	"github.com/agnath18K/lumo/internal/core"
)

// Assistant implements the core.Assistant interface
type Assistant struct {
	// factory is the desktop environment factory
	factory core.DesktopFactory
	// processor is the natural language processor
	processor *Processor
}

// NewAssistant creates a new desktop assistant
func NewAssistant(factory core.DesktopFactory) *Assistant {
	return &Assistant{
		factory:   factory,
		processor: NewProcessor(),
	}
}

// NewAssistantWithAI creates a new desktop assistant with AI capabilities
func NewAssistantWithAI(factory core.DesktopFactory, aiClient AIClient) *Assistant {
	return &Assistant{
		factory:   factory,
		processor: NewProcessorWithAI(aiClient),
	}
}

// ProcessCommand processes a natural language command
func (a *Assistant) ProcessCommand(ctx context.Context, input string) (*core.Result, error) {
	// Process the input to extract the command
	cmd, err := a.processor.Process(input)
	if err != nil {
		return nil, fmt.Errorf("failed to process command: %w", err)
	}

	// Get the desktop environment
	env, err := a.factory.DetectEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to detect desktop environment: %w", err)
	}

	// Execute the command
	return env.ExecuteCommand(ctx, cmd)
}

// GetSupportedCommands returns a list of supported commands
func (a *Assistant) GetSupportedCommands() []string {
	return []string{
		"window:close <window>",
		"window:minimize <window>",
		"window:maximize <window>",
		"window:restore <window>",
		"window:focus <window>",
		"window:list",
		"application:launch <app> [args]",
		"application:list",
		"system:shutdown",
		"system:restart",
		"system:logout",
		"system:lock",
		"notification:send <summary> [body] [icon]",
		"notification:close <id>",
		"media:play",
		"media:pause",
		"media:stop",
		"media:next",
		"media:previous",
	}
}

// GetCommandExamples returns examples of supported commands
func (a *Assistant) GetCommandExamples() []string {
	return []string{
		"Close Firefox window",
		"Minimize all windows",
		"Maximize the current window",
		"Show all open windows",
		"Launch Firefox",
		"Open Terminal",
		"List running applications",
		"Lock the screen",
		"Shutdown the computer",
		"Restart the system",
		"Log out",
		"Send a notification with the message 'Hello World'",
		"Play music",
		"Pause media playback",
		"Skip to the next track",
		"Go to the previous song",
	}
}
