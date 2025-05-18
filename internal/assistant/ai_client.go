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
- appearance (for appearance settings)
- sound (for sound settings)
- connectivity (for network connectivity settings)

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

Valid actions for appearance:
- set-theme (set GTK theme)
- set-dark-mode (enable/disable dark mode)
- set-background (set desktop background)
- set-accent-color (set accent color)
- set-icon-theme (set icon theme)
- get-theme (get current GTK theme)
- get-background (get current desktop background)
- get-icon-theme (get current icon theme)

Valid actions for sound:
- set-volume (set system volume level)
- get-volume (get current system volume level)
- set-mute (set system mute state)
- get-mute (get current system mute state)
- set-input-volume (set microphone volume level)
- get-input-volume (get current microphone volume level)
- set-input-mute (set microphone mute state)
- get-input-mute (get current microphone mute state)
- list-devices (list available sound devices)
- set-default-device (set default sound device)

Valid actions for connectivity:
- list-devices (list all network devices)
- enable-wifi (enable WiFi)
- disable-wifi (disable WiFi)
- wifi-status (get WiFi status)
- enable-bluetooth (enable Bluetooth)
- disable-bluetooth (disable Bluetooth)
- bluetooth-status (get Bluetooth status)
- enable-airplane-mode (enable airplane mode)
- disable-airplane-mode (disable airplane mode)
- airplane-mode-status (get airplane mode status)
- enable-hotspot (enable WiFi hotspot)
- disable-hotspot (disable WiFi hotspot)
- hotspot-status (get WiFi hotspot status)

Examples:
- "Close Firefox window" -> "window:close:firefox"
- "Launch Terminal" -> "application:launch:gnome-terminal"
- "Lock the screen" -> "system:lock:"
- "Send notification Hello World with body This is a test" -> "notification:send:Hello World:body=This is a test"
- "Play media" -> "media:play:"
- "Launch Firefox and maximize it" -> "application:launch:firefox"
- "Set dark mode on" -> "appearance:set-dark-mode:on"
- "Change desktop background to /path/to/image.jpg" -> "appearance:set-background:/path/to/image.jpg"
- "Get current theme" -> "appearance:get-theme:"
- "Set GTK theme to Adwaita-dark" -> "appearance:set-theme:Adwaita-dark"
- "Set volume to 50 percent" -> "sound:set-volume:50"
- "Mute the sound" -> "sound:set-mute:true"
- "Unmute the microphone" -> "sound:set-input-mute:false"
- "Show sound devices" -> "sound:list-devices:"
- "Set microphone volume to 75 percent" -> "sound:set-input-volume:75"
- "Show all network devices" -> "connectivity:list-devices:"
- "Turn on WiFi" -> "connectivity:enable-wifi:"
- "Turn off Bluetooth" -> "connectivity:disable-bluetooth:"
- "Check airplane mode status" -> "connectivity:airplane-mode-status:"
- "Create a WiFi hotspot with name MyHotspot" -> "connectivity:enable-hotspot:MyHotspot"

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
