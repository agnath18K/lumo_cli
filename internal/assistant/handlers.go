package assistant

import (
	"github.com/agnath18K/lumo/internal/core"
)

// handleCloseWindow handles the "close window" command
func (p *Processor) handleCloseWindow(input string) (*core.Command, error) {
	// Extract the window name
	windowName := extractTarget(input, []string{"close", "window"})

	return &core.Command{
		Type:      core.CommandTypeWindow,
		Action:    "close",
		Target:    windowName,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleMinimizeWindow handles the "minimize window" command
func (p *Processor) handleMinimizeWindow(input string) (*core.Command, error) {
	// Extract the window name
	windowName := extractTarget(input, []string{"minimize", "window"})

	return &core.Command{
		Type:      core.CommandTypeWindow,
		Action:    "minimize",
		Target:    windowName,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleMaximizeWindow handles the "maximize window" command
func (p *Processor) handleMaximizeWindow(input string) (*core.Command, error) {
	// Extract the window name
	windowName := extractTarget(input, []string{"maximize", "window"})

	return &core.Command{
		Type:      core.CommandTypeWindow,
		Action:    "maximize",
		Target:    windowName,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleRestoreWindow handles the "restore window" command
func (p *Processor) handleRestoreWindow(input string) (*core.Command, error) {
	// Extract the window name
	windowName := extractTarget(input, []string{"restore", "window"})

	return &core.Command{
		Type:      core.CommandTypeWindow,
		Action:    "restore",
		Target:    windowName,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleFocusWindow handles the "focus window" command
func (p *Processor) handleFocusWindow(input string) (*core.Command, error) {
	// Extract the window name
	windowName := extractTarget(input, []string{"focus", "window"})

	return &core.Command{
		Type:      core.CommandTypeWindow,
		Action:    "focus",
		Target:    windowName,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleListWindows handles the "list windows" command
func (p *Processor) handleListWindows(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeWindow,
		Action:    "list",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleLaunchApplication handles the "launch application" command
func (p *Processor) handleLaunchApplication(input string) (*core.Command, error) {
	// Extract the application name and arguments
	appName, args := extractApplicationAndArgs(input)

	// Create the command
	cmd := &core.Command{
		Type:      core.CommandTypeApplication,
		Action:    "launch",
		Target:    appName,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}

	// Add arguments if any
	if args != "" {
		cmd.Arguments["args"] = args
	}

	return cmd, nil
}

// handleListApplications handles the "list applications" command
func (p *Processor) handleListApplications(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeApplication,
		Action:    "list",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleShutdownSystem handles the "shutdown system" command
func (p *Processor) handleShutdownSystem(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeSystem,
		Action:    "shutdown",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleRestartSystem handles the "restart system" command
func (p *Processor) handleRestartSystem(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeSystem,
		Action:    "restart",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleLogout handles the "logout" command
func (p *Processor) handleLogout(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeSystem,
		Action:    "logout",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleLockScreen handles the "lock screen" command
func (p *Processor) handleLockScreen(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeSystem,
		Action:    "lock",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleSendNotification handles the "send notification" command
func (p *Processor) handleSendNotification(input string) (*core.Command, error) {
	// Extract the notification summary and body
	summary, body := extractNotificationContent(input)

	// Create the command
	cmd := &core.Command{
		Type:      core.CommandTypeNotification,
		Action:    "send",
		Target:    summary,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}

	// Add body if any
	if body != "" {
		cmd.Arguments["body"] = body
	}

	return cmd, nil
}

// handleCloseNotification handles the "close notification" command
func (p *Processor) handleCloseNotification(input string) (*core.Command, error) {
	// Extract the notification ID
	notificationID := extractTarget(input, []string{"close", "notification"})

	return &core.Command{
		Type:      core.CommandTypeNotification,
		Action:    "close",
		Target:    notificationID,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handlePlayMedia handles the "play media" command
func (p *Processor) handlePlayMedia(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeMedia,
		Action:    "play",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handlePauseMedia handles the "pause media" command
func (p *Processor) handlePauseMedia(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeMedia,
		Action:    "pause",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleStopMedia handles the "stop media" command
func (p *Processor) handleStopMedia(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeMedia,
		Action:    "stop",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handleNextTrack handles the "next track" command
func (p *Processor) handleNextTrack(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeMedia,
		Action:    "next",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}

// handlePreviousTrack handles the "previous track" command
func (p *Processor) handlePreviousTrack(input string) (*core.Command, error) {
	return &core.Command{
		Type:      core.CommandTypeMedia,
		Action:    "previous",
		Target:    "",
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}, nil
}
