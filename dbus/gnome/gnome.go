package gnome

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/agnath18K/lumo/dbus/common"
	"github.com/agnath18K/lumo/internal/core"
	"github.com/agnath18K/lumo/internal/desktop"
)

// Environment implements the core.DesktopEnvironment interface for GNOME
type Environment struct {
	*desktop.BaseEnvironment
	sessionHandler core.DBusHandler
	systemHandler  core.DBusHandler
	// Keep a reference to the connections to prevent them from being closed
	sessionConn common.DBusConnection
	systemConn  common.DBusConnection
}

// NewEnvironment creates a new GNOME desktop environment
func NewEnvironment() (*Environment, error) {
	// Create session DBus connection
	sessionConn, err := common.NewDBusConnection(common.DBusTypeSession)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session DBus: %w", err)
	}

	// Create system DBus connection
	systemConn, err := common.NewDBusConnection(common.DBusTypeSystem)
	if err != nil {
		// Close the session connection if system connection fails
		sessionConn.Close()
		return nil, fmt.Errorf("failed to connect to system DBus: %w", err)
	}

	// Create DBus handlers
	sessionHandler := common.NewDBusHandler(sessionConn)
	systemHandler := common.NewDBusHandler(systemConn)

	// Define capabilities
	capabilities := []core.Capability{
		core.CapabilityWindowManagement,
		core.CapabilityApplicationLaunch,
		core.CapabilityNotifications,
		core.CapabilityMediaControl,
		core.CapabilityScreenshot,
		core.CapabilityClipboard,
		core.CapabilityAppearanceManagement,
		core.CapabilitySoundManagement,
	}

	// Create base environment
	baseEnv := desktop.NewBaseEnvironment("gnome", capabilities, sessionHandler)

	return &Environment{
		BaseEnvironment: baseEnv,
		sessionHandler:  sessionHandler,
		systemHandler:   systemHandler,
		sessionConn:     sessionConn,
		systemConn:      systemConn,
	}, nil
}

// IsAvailable checks if GNOME is available on the system
func (e *Environment) IsAvailable() bool {
	fmt.Printf("DEBUG: Checking if GNOME is available...\n")

	// Use the existing session connection instead of creating a new one
	if e.sessionConn == nil {
		fmt.Printf("DEBUG: Session connection is nil, creating a new one\n")
		conn, err := common.NewDBusConnection(common.DBusTypeSession)
		if err != nil {
			fmt.Printf("DEBUG: Failed to connect to session DBus: %v\n", err)
			return false
		}
		// Store the connection for future use
		e.sessionConn = conn
	}

	// Check if GNOME Shell service is available
	available := common.IsDBusServiceAvailable(e.sessionConn, Shell)
	if available {
		fmt.Printf("DEBUG: GNOME Shell service is available\n")
	} else {
		fmt.Printf("DEBUG: GNOME Shell service is not available\n")
	}

	return available
}

// ExecuteCommand executes a desktop command
func (e *Environment) ExecuteCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Type {
	case core.CommandTypeWindow:
		return e.executeWindowCommand(ctx, cmd)
	case core.CommandTypeApplication:
		return e.executeApplicationCommand(ctx, cmd)
	case core.CommandTypeSystem:
		return e.executeSystemCommand(ctx, cmd)
	case core.CommandTypeNotification:
		return e.executeNotificationCommand(ctx, cmd)
	case core.CommandTypeMedia:
		return e.executeMediaCommand(ctx, cmd)
	case core.CommandTypeAppearance:
		return e.executeAppearanceCommand(ctx, cmd)
	case core.CommandTypeSound:
		return e.executeSoundCommand(ctx, cmd)
	default:
		return nil, fmt.Errorf("unsupported command type: %s", cmd.Type)
	}
}

// executeWindowCommand executes a window management command
func (e *Environment) executeWindowCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Action {
	case "close":
		if err := e.CloseWindow(ctx, cmd.Target); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Closed window: %s", cmd.Target),
			Success: true,
		}, nil
	case "minimize":
		if err := e.MinimizeWindow(ctx, cmd.Target); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Minimized window: %s", cmd.Target),
			Success: true,
		}, nil
	case "maximize":
		if err := e.MaximizeWindow(ctx, cmd.Target); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Maximized window: %s", cmd.Target),
			Success: true,
		}, nil
	case "restore":
		if err := e.RestoreWindow(ctx, cmd.Target); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Restored window: %s", cmd.Target),
			Success: true,
		}, nil
	case "focus":
		if err := e.FocusWindow(ctx, cmd.Target); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Focused window: %s", cmd.Target),
			Success: true,
		}, nil
	case "list":
		windows, err := e.GetWindows(ctx)
		if err != nil {
			return nil, err
		}
		var output strings.Builder
		output.WriteString("Windows:\n")
		for _, window := range windows {
			output.WriteString(fmt.Sprintf("- %s (%s)\n", window.Title, window.Application))
		}
		return &core.Result{
			Output:  output.String(),
			Success: true,
			Data: map[string]interface{}{
				"windows": windows,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported window action: %s", cmd.Action)
	}
}

// executeApplicationCommand executes an application management command
func (e *Environment) executeApplicationCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Action {
	case "launch":
		args := []string{}
		if argsVal, ok := cmd.Arguments["args"]; ok {
			if argsStr, ok := argsVal.(string); ok {
				args = strings.Fields(argsStr)
			} else if argsSlice, ok := argsVal.([]string); ok {
				args = argsSlice
			}
		}
		if err := e.LaunchApplication(ctx, cmd.Target, args...); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Launched application: %s", cmd.Target),
			Success: true,
		}, nil
	case "list":
		apps, err := e.GetRunningApplications(ctx)
		if err != nil {
			return nil, err
		}
		var output strings.Builder
		output.WriteString("Running applications:\n")
		for _, app := range apps {
			output.WriteString(fmt.Sprintf("- %s\n", app.Name))
		}
		return &core.Result{
			Output:  output.String(),
			Success: true,
			Data: map[string]interface{}{
				"applications": apps,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported application action: %s", cmd.Action)
	}
}

// executeSystemCommand executes a system command
func (e *Environment) executeSystemCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Action {
	case "shutdown":
		// Call the DBus method to shut down the system
		_, err := e.sessionHandler.Call(
			SessionManager,
			SessionManagerPath,
			SessionManagerInterface,
			"Shutdown",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to shut down: %w", err)
		}
		return &core.Result{
			Output:  "System is shutting down",
			Success: true,
		}, nil
	case "restart":
		// Call the DBus method to restart the system
		_, err := e.sessionHandler.Call(
			SessionManager,
			SessionManagerPath,
			SessionManagerInterface,
			"Reboot",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to restart: %w", err)
		}
		return &core.Result{
			Output:  "System is restarting",
			Success: true,
		}, nil
	case "logout":
		// Call the DBus method to log out
		_, err := e.sessionHandler.Call(
			SessionManager,
			SessionManagerPath,
			SessionManagerInterface,
			"Logout",
			uint32(0), // Normal logout
		)
		if err != nil {
			return nil, fmt.Errorf("failed to log out: %w", err)
		}
		return &core.Result{
			Output:  "Logging out",
			Success: true,
		}, nil
	case "lock":
		// Call the DBus method to lock the screen
		_, err := e.sessionHandler.Call(
			"org.gnome.ScreenSaver",
			"/org/gnome/ScreenSaver",
			"org.gnome.ScreenSaver",
			"Lock",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to lock screen: %w", err)
		}
		return &core.Result{
			Output:  "Screen locked",
			Success: true,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported system action: %s", cmd.Action)
	}
}

// executeNotificationCommand executes a notification command
func (e *Environment) executeNotificationCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Action {
	case "send":
		// Get notification parameters
		summary := cmd.Target
		body := ""
		icon := ""

		if bodyVal, ok := cmd.Arguments["body"]; ok {
			if bodyStr, ok := bodyVal.(string); ok {
				body = bodyStr
			}
		}

		if iconVal, ok := cmd.Arguments["icon"]; ok {
			if iconStr, ok := iconVal.(string); ok {
				icon = iconStr
			}
		}

		// Send the notification
		id, err := e.SendNotification(ctx, summary, body, icon)
		if err != nil {
			return nil, err
		}

		return &core.Result{
			Output:  fmt.Sprintf("Notification sent (ID: %d)", id),
			Success: true,
			Data: map[string]interface{}{
				"notification_id": id,
			},
		}, nil
	case "close":
		// Get notification ID
		idStr := cmd.Target
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid notification ID: %s", idStr)
		}

		// Close the notification
		if err := e.CloseNotification(ctx, uint32(id)); err != nil {
			return nil, err
		}

		return &core.Result{
			Output:  fmt.Sprintf("Notification closed (ID: %d)", id),
			Success: true,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported notification action: %s", cmd.Action)
	}
}

// executeMediaCommand executes a media control command
func (e *Environment) executeMediaCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	// Find the active media player
	playerService := ""

	// List DBus services
	conn, err := common.NewDBusConnection(common.DBusTypeSession)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DBus: %w", err)
	}
	defer conn.Close()

	services, err := common.ListDBusServices(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to list DBus services: %w", err)
	}

	// Find a media player service
	for _, service := range services {
		if strings.HasPrefix(service, "org.mpris.MediaPlayer2.") {
			playerService = service
			break
		}
	}

	if playerService == "" {
		return nil, fmt.Errorf("no active media player found")
	}

	// Execute the command
	switch cmd.Action {
	case "play":
		_, err := e.sessionHandler.Call(
			playerService,
			"/org/mpris/MediaPlayer2",
			"org.mpris.MediaPlayer2.Player",
			"Play",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to play media: %w", err)
		}
		return &core.Result{
			Output:  "Media playback started",
			Success: true,
		}, nil
	case "pause":
		_, err := e.sessionHandler.Call(
			playerService,
			"/org/mpris/MediaPlayer2",
			"org.mpris.MediaPlayer2.Player",
			"Pause",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to pause media: %w", err)
		}
		return &core.Result{
			Output:  "Media playback paused",
			Success: true,
		}, nil
	case "stop":
		_, err := e.sessionHandler.Call(
			playerService,
			"/org/mpris/MediaPlayer2",
			"org.mpris.MediaPlayer2.Player",
			"Stop",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to stop media: %w", err)
		}
		return &core.Result{
			Output:  "Media playback stopped",
			Success: true,
		}, nil
	case "next":
		_, err := e.sessionHandler.Call(
			playerService,
			"/org/mpris/MediaPlayer2",
			"org.mpris.MediaPlayer2.Player",
			"Next",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to go to next track: %w", err)
		}
		return &core.Result{
			Output:  "Skipped to next track",
			Success: true,
		}, nil
	case "previous":
		_, err := e.sessionHandler.Call(
			playerService,
			"/org/mpris/MediaPlayer2",
			"org.mpris.MediaPlayer2.Player",
			"Previous",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to go to previous track: %w", err)
		}
		return &core.Result{
			Output:  "Skipped to previous track",
			Success: true,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported media action: %s", cmd.Action)
	}
}
