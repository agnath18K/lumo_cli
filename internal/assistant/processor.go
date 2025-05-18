package assistant

import (
	"fmt"
	"strings"

	"github.com/agnath18K/lumo/internal/core"
)

// Processor processes natural language commands
type Processor struct {
	// commandPatterns maps command patterns to handlers
	commandPatterns map[string]commandHandler
	// aiEnabled indicates whether AI processing is enabled
	aiEnabled bool
	// aiClient is the AI client for processing complex commands
	aiClient AIClient
}

// AIClient is an interface for AI processing
type AIClient interface {
	// ProcessNLP processes a natural language command using AI
	ProcessNLP(input string) (string, error)
}

// commandHandler is a function that handles a command
type commandHandler func(input string) (*core.Command, error)

// NewProcessor creates a new natural language processor
func NewProcessor() *Processor {
	p := &Processor{
		commandPatterns: make(map[string]commandHandler),
		aiEnabled:       false,
		aiClient:        nil,
	}

	// Register command patterns
	p.registerCommandPatterns()

	return p
}

// NewProcessorWithAI creates a new natural language processor with AI capabilities
func NewProcessorWithAI(aiClient AIClient) *Processor {
	p := &Processor{
		commandPatterns: make(map[string]commandHandler),
		aiEnabled:       true,
		aiClient:        aiClient,
	}

	// Register command patterns
	p.registerCommandPatterns()

	return p
}

// registerCommandPatterns registers command patterns
func (p *Processor) registerCommandPatterns() {
	// Window commands
	p.commandPatterns["close window"] = p.handleCloseWindow
	p.commandPatterns["minimize window"] = p.handleMinimizeWindow
	p.commandPatterns["maximize window"] = p.handleMaximizeWindow
	p.commandPatterns["restore window"] = p.handleRestoreWindow
	p.commandPatterns["focus window"] = p.handleFocusWindow
	p.commandPatterns["list windows"] = p.handleListWindows

	// Application commands
	p.commandPatterns["launch application"] = p.handleLaunchApplication
	p.commandPatterns["open application"] = p.handleLaunchApplication
	p.commandPatterns["start application"] = p.handleLaunchApplication
	p.commandPatterns["list applications"] = p.handleListApplications

	// System commands
	p.commandPatterns["shutdown system"] = p.handleShutdownSystem
	p.commandPatterns["restart system"] = p.handleRestartSystem
	p.commandPatterns["logout"] = p.handleLogout
	p.commandPatterns["lock screen"] = p.handleLockScreen

	// Notification commands
	p.commandPatterns["send notification"] = p.handleSendNotification
	p.commandPatterns["close notification"] = p.handleCloseNotification

	// Media commands
	p.commandPatterns["play media"] = p.handlePlayMedia
	p.commandPatterns["pause media"] = p.handlePauseMedia
	p.commandPatterns["stop media"] = p.handleStopMedia
	p.commandPatterns["next track"] = p.handleNextTrack
	p.commandPatterns["previous track"] = p.handlePreviousTrack

	// Connectivity commands
	p.commandPatterns["list network devices"] = p.handleListNetworkDevices
	p.commandPatterns["enable wifi"] = p.handleEnableWifi
	p.commandPatterns["disable wifi"] = p.handleDisableWifi
	p.commandPatterns["wifi status"] = p.handleWifiStatus
	p.commandPatterns["enable bluetooth"] = p.handleEnableBluetooth
	p.commandPatterns["disable bluetooth"] = p.handleDisableBluetooth
	p.commandPatterns["bluetooth status"] = p.handleBluetoothStatus
	p.commandPatterns["enable airplane mode"] = p.handleEnableAirplaneMode
	p.commandPatterns["disable airplane mode"] = p.handleDisableAirplaneMode
	p.commandPatterns["airplane mode status"] = p.handleAirplaneModeStatus
	p.commandPatterns["enable hotspot"] = p.handleEnableHotspot
	p.commandPatterns["disable hotspot"] = p.handleDisableHotspot
	p.commandPatterns["hotspot status"] = p.handleHotspotStatus
}

// Process processes a natural language command
func (p *Processor) Process(input string) (*core.Command, error) {
	fmt.Printf("DEBUG: Processing command: %s\n", input)

	// If AI is enabled, try to use AI first for complex queries
	if p.aiEnabled && p.aiClient != nil && (len(input) > 15 ||
		strings.Contains(strings.ToLower(input), "can you") ||
		strings.Contains(strings.ToLower(input), "would like") ||
		strings.Contains(strings.ToLower(input), "please") ||
		strings.Contains(strings.ToLower(input), "could you") ||
		strings.Contains(strings.ToLower(input), "i want")) {
		fmt.Printf("DEBUG: Complex query detected, trying AI processing first\n")
		cmd, err := p.processWithAI(input)
		if err == nil {
			fmt.Printf("DEBUG: AI processing successful\n")
			return cmd, nil
		}
		fmt.Printf("DEBUG: AI processing failed: %v, falling back to pattern matching\n", err)
	}

	// Normalize the input
	normalizedInput := strings.ToLower(strings.TrimSpace(input))
	fmt.Printf("DEBUG: Normalized input: %s\n", normalizedInput)

	// Try to match the input to a command pattern
	for pattern, handler := range p.commandPatterns {
		if strings.Contains(normalizedInput, pattern) {
			fmt.Printf("DEBUG: Found matching pattern: %s\n", pattern)
			cmd, err := handler(normalizedInput)
			if err != nil {
				fmt.Printf("DEBUG: Error handling command: %v\n", err)
				return nil, err
			}
			fmt.Printf("DEBUG: Command processed: Type=%s, Action=%s, Target=%s\n", cmd.Type, cmd.Action, cmd.Target)
			return cmd, nil
		}
	}

	fmt.Printf("DEBUG: No pattern match found, trying to infer command\n")
	// If no pattern matches, try to infer the command
	cmd, err := p.inferCommand(normalizedInput)
	if err != nil {
		fmt.Printf("DEBUG: Failed to infer command: %v\n", err)

		// If AI is enabled, try to use AI to process the command
		if p.aiEnabled && p.aiClient != nil {
			fmt.Printf("DEBUG: AI is enabled and client is available, trying AI processing\n")
			return p.processWithAI(input)
		} else {
			fmt.Printf("DEBUG: AI is not enabled or client is not available. aiEnabled=%v, aiClient=%v\n", p.aiEnabled, p.aiClient != nil)
		}
	} else {
		// Check if the target looks like a sentence (more than 3 words)
		words := strings.Fields(cmd.Target)
		if len(words) > 3 && p.aiEnabled && p.aiClient != nil {
			fmt.Printf("DEBUG: Target looks like a sentence, trying AI processing\n")
			return p.processWithAI(input)
		}

		fmt.Printf("DEBUG: Command inferred: Type=%s, Action=%s, Target=%s\n", cmd.Type, cmd.Action, cmd.Target)
	}
	return cmd, err
}

// processWithAI processes a command using AI
func (p *Processor) processWithAI(input string) (*core.Command, error) {
	fmt.Printf("DEBUG: Processing with AI: %s\n", input)

	// Use AI to process the command
	aiResult, err := p.aiClient.ProcessNLP(input)
	if err != nil {
		fmt.Printf("DEBUG: AI processing error: %v\n", err)
		return nil, fmt.Errorf("AI processing error: %w", err)
	}

	fmt.Printf("DEBUG: AI result: %s\n", aiResult)

	// Parse the AI result to extract the command
	// The AI result should be in the format: "TYPE:ACTION:TARGET[:ARG1=VAL1,ARG2=VAL2,...]"
	parts := strings.Split(aiResult, ":")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid AI result format: %s", aiResult)
	}

	// Extract the command type, action, and target
	cmdType := strings.TrimSpace(parts[0])
	action := strings.TrimSpace(parts[1])
	target := strings.TrimSpace(parts[2])

	// Create the command
	cmd := &core.Command{
		Type:      core.CommandType(cmdType),
		Action:    action,
		Target:    target,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}

	// Extract arguments if any
	if len(parts) > 3 {
		argStr := strings.TrimSpace(parts[3])
		args := strings.Split(argStr, ",")
		for _, arg := range args {
			keyVal := strings.Split(arg, "=")
			if len(keyVal) == 2 {
				key := strings.TrimSpace(keyVal[0])
				val := strings.TrimSpace(keyVal[1])
				cmd.Arguments[key] = val
			}
		}
	}

	fmt.Printf("DEBUG: AI command processed: Type=%s, Action=%s, Target=%s\n", cmd.Type, cmd.Action, cmd.Target)
	return cmd, nil
}

// inferCommand tries to infer the command from the input
func (p *Processor) inferCommand(input string) (*core.Command, error) {
	fmt.Printf("DEBUG: Inferring command from: %s\n", input)

	// Check for window commands
	if strings.Contains(input, "close") && (strings.Contains(input, "window") || strings.Contains(input, "app")) {
		return p.handleCloseWindow(input)
	}
	if strings.Contains(input, "minimize") && (strings.Contains(input, "window") || strings.Contains(input, "app")) {
		return p.handleMinimizeWindow(input)
	}
	if strings.Contains(input, "maximize") && (strings.Contains(input, "window") || strings.Contains(input, "app")) {
		return p.handleMaximizeWindow(input)
	}
	if strings.Contains(input, "restore") && (strings.Contains(input, "window") || strings.Contains(input, "app")) {
		return p.handleRestoreWindow(input)
	}
	if strings.Contains(input, "focus") && (strings.Contains(input, "window") || strings.Contains(input, "app")) {
		return p.handleFocusWindow(input)
	}
	if strings.Contains(input, "list") && strings.Contains(input, "window") {
		return p.handleListWindows(input)
	}

	// Check for application commands - more flexible patterns
	if strings.Contains(input, "launch") || strings.Contains(input, "open") || strings.Contains(input, "start") ||
		strings.Contains(input, "run") {
		// This is likely a launch application command even if "application" is not explicitly mentioned
		fmt.Printf("DEBUG: Detected launch application command\n")
		return p.handleLaunchApplication(input)
	}

	if strings.Contains(input, "list") && (strings.Contains(input, "app") || strings.Contains(input, "application") || strings.Contains(input, "program")) {
		return p.handleListApplications(input)
	}

	// Check for system commands
	if strings.Contains(input, "shutdown") || strings.Contains(input, "turn off") || strings.Contains(input, "power off") {
		return p.handleShutdownSystem(input)
	}
	if strings.Contains(input, "restart") || strings.Contains(input, "reboot") {
		return p.handleRestartSystem(input)
	}
	if strings.Contains(input, "logout") || strings.Contains(input, "log out") || strings.Contains(input, "sign out") {
		return p.handleLogout(input)
	}
	if strings.Contains(input, "lock") && strings.Contains(input, "screen") {
		return p.handleLockScreen(input)
	}

	// Check for notification commands
	if strings.Contains(input, "send") && strings.Contains(input, "notification") {
		return p.handleSendNotification(input)
	}
	if strings.Contains(input, "close") && strings.Contains(input, "notification") {
		return p.handleCloseNotification(input)
	}

	// Check for media commands
	if strings.Contains(input, "play") && (strings.Contains(input, "media") || strings.Contains(input, "music") || strings.Contains(input, "song")) {
		return p.handlePlayMedia(input)
	}
	if strings.Contains(input, "pause") && (strings.Contains(input, "media") || strings.Contains(input, "music") || strings.Contains(input, "song")) {
		return p.handlePauseMedia(input)
	}
	if strings.Contains(input, "stop") && (strings.Contains(input, "media") || strings.Contains(input, "music") || strings.Contains(input, "song")) {
		return p.handleStopMedia(input)
	}
	if (strings.Contains(input, "next") || strings.Contains(input, "skip")) && (strings.Contains(input, "track") || strings.Contains(input, "song")) {
		return p.handleNextTrack(input)
	}
	if strings.Contains(input, "previous") && (strings.Contains(input, "track") || strings.Contains(input, "song")) {
		return p.handlePreviousTrack(input)
	}

	// Special cases for common applications
	if strings.Contains(input, "terminal") || strings.Contains(input, "console") {
		fmt.Printf("DEBUG: Special case: terminal command detected\n")
		return p.handleLaunchApplication("launch application terminal")
	}
	if strings.Contains(input, "firefox") || strings.Contains(input, "browser") {
		fmt.Printf("DEBUG: Special case: browser command detected\n")
		return p.handleLaunchApplication("launch application firefox")
	}
	if strings.Contains(input, "chrome") {
		fmt.Printf("DEBUG: Special case: chrome command detected\n")
		return p.handleLaunchApplication("launch application chrome")
	}

	// Check for connectivity commands
	if strings.Contains(input, "list") && (strings.Contains(input, "network") || strings.Contains(input, "device")) {
		return p.handleListNetworkDevices(input)
	}
	if (strings.Contains(input, "enable") || strings.Contains(input, "turn on")) && strings.Contains(input, "wifi") {
		return p.handleEnableWifi(input)
	}
	if (strings.Contains(input, "disable") || strings.Contains(input, "turn off")) && strings.Contains(input, "wifi") {
		return p.handleDisableWifi(input)
	}
	if strings.Contains(input, "status") && strings.Contains(input, "wifi") {
		return p.handleWifiStatus(input)
	}
	if (strings.Contains(input, "enable") || strings.Contains(input, "turn on")) && strings.Contains(input, "bluetooth") {
		return p.handleEnableBluetooth(input)
	}
	if (strings.Contains(input, "disable") || strings.Contains(input, "turn off")) && strings.Contains(input, "bluetooth") {
		return p.handleDisableBluetooth(input)
	}
	if strings.Contains(input, "status") && strings.Contains(input, "bluetooth") {
		return p.handleBluetoothStatus(input)
	}
	if (strings.Contains(input, "enable") || strings.Contains(input, "turn on")) && strings.Contains(input, "airplane") {
		return p.handleEnableAirplaneMode(input)
	}
	if (strings.Contains(input, "disable") || strings.Contains(input, "turn off")) && strings.Contains(input, "airplane") {
		return p.handleDisableAirplaneMode(input)
	}
	if strings.Contains(input, "status") && strings.Contains(input, "airplane") {
		return p.handleAirplaneModeStatus(input)
	}
	if (strings.Contains(input, "enable") || strings.Contains(input, "turn on") || strings.Contains(input, "create")) && strings.Contains(input, "hotspot") {
		return p.handleEnableHotspot(input)
	}
	if (strings.Contains(input, "disable") || strings.Contains(input, "turn off")) && strings.Contains(input, "hotspot") {
		return p.handleDisableHotspot(input)
	}
	if strings.Contains(input, "status") && strings.Contains(input, "hotspot") {
		return p.handleHotspotStatus(input)
	}

	// If no command can be inferred, return an error
	return nil, fmt.Errorf("could not understand command: %s", input)
}
