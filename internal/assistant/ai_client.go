package assistant

import (
	"context"
	"fmt"
	"strings"

	"github.com/agnath18K/lumo/pkg/ai"
)

// AIClientImpl implements the AIClient interface using the AI service
type AIClientImpl struct {
	// aiClient is the AI client
	aiClient ai.Client
}

// NewAIClient creates a new AI client
func NewAIClient(aiClient ai.Client) *AIClientImpl {
	return &AIClientImpl{
		aiClient: aiClient,
	}
}

// ProcessNLP processes a natural language command using AI
func (c *AIClientImpl) ProcessNLP(input string) (string, error) {
	// Create a context
	ctx := context.Background()

	// Create a prompt for the AI to process the command
	prompt := fmt.Sprintf(`
You are an AI assistant that helps process desktop commands.
Convert the following natural language command into a structured format.

Command: %s

The output should be in the format: "TYPE:ACTION:TARGET[:ARG1=VAL1,ARG2=VAL2,...]"

If the command involves multiple actions, only return the first action. Do not return multiple commands.

Valid command types:
- window (for window operations)
- application (for application operations)
- system (for system operations)
- notification (for notification operations)
- media (for media operations)

Valid actions for window:
- close (close a window)
- minimize (minimize a window)
- maximize (maximize a window)
- restore (restore a window)
- focus (focus a window)
- list (list all windows)

Valid actions for application:
- launch (launch an application)
- list (list all applications)

Valid actions for system:
- shutdown (shutdown the system)
- restart (restart the system)
- logout (logout the user)
- lock (lock the screen)

Valid actions for notification:
- send (send a notification)
- close (close a notification)

Valid actions for media:
- play (play media)
- pause (pause media)
- stop (stop media)
- next (next track)
- previous (previous track)

Examples:
- "Close Firefox window" -> "window:close:firefox"
- "Launch Terminal" -> "application:launch:gnome-terminal"
- "Lock the screen" -> "system:lock:"
- "Send notification Hello World with body This is a test" -> "notification:send:Hello World:body=This is a test"
- "Play media" -> "media:play:"
- "Launch Firefox and maximize it" -> "application:launch:firefox"

Only output the structured format, nothing else. Do not include newlines or multiple commands.
`, input)

	// Get completion from AI
	completion, err := c.aiClient.GetCompletion(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to get completion: %w", err)
	}

	// Clean up the completion
	result := strings.TrimSpace(completion)

	// Remove quotes if present
	result = strings.Trim(result, "\"'")

	// If there are multiple lines, only take the first line
	if strings.Contains(result, "\n") {
		result = strings.Split(result, "\n")[0]
	}

	return result, nil
}
