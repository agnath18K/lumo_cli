package desktop

import (
	"context"
	"fmt"

	"github.com/agnath18K/lumo/internal/core"
)

// BaseEnvironment provides a base implementation of the core.DesktopEnvironment interface
type BaseEnvironment struct {
	// name is the name of the desktop environment
	name string
	// capabilities is a list of capabilities supported by this environment
	capabilities []core.Capability
	// dbusHandler is the DBus handler for this environment
	dbusHandler core.DBusHandler
}

// NewBaseEnvironment creates a new base desktop environment
func NewBaseEnvironment(name string, capabilities []core.Capability, dbusHandler core.DBusHandler) *BaseEnvironment {
	return &BaseEnvironment{
		name:         name,
		capabilities: capabilities,
		dbusHandler:  dbusHandler,
	}
}

// Name returns the name of the desktop environment
func (e *BaseEnvironment) Name() string {
	return e.name
}

// IsAvailable checks if this desktop environment is available on the system
func (e *BaseEnvironment) IsAvailable() bool {
	// This should be overridden by specific implementations
	return false
}

// ExecuteCommand executes a desktop command
func (e *BaseEnvironment) ExecuteCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	// This should be overridden by specific implementations
	return nil, fmt.Errorf("not implemented")
}

// GetCapabilities returns the capabilities of this desktop environment
func (e *BaseEnvironment) GetCapabilities() []core.Capability {
	return e.capabilities
}

// GetWindows returns a list of all windows
func (e *BaseEnvironment) GetWindows(ctx context.Context) ([]core.Window, error) {
	// This should be overridden by specific implementations
	return nil, fmt.Errorf("not implemented")
}

// GetActiveWindow returns the currently active window
func (e *BaseEnvironment) GetActiveWindow(ctx context.Context) (*core.Window, error) {
	// This should be overridden by specific implementations
	return nil, fmt.Errorf("not implemented")
}

// GetRunningApplications returns a list of running applications
func (e *BaseEnvironment) GetRunningApplications(ctx context.Context) ([]core.Application, error) {
	// This should be overridden by specific implementations
	return nil, fmt.Errorf("not implemented")
}

// LaunchApplication launches an application
func (e *BaseEnvironment) LaunchApplication(ctx context.Context, appName string, args ...string) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// CloseWindow closes a window
func (e *BaseEnvironment) CloseWindow(ctx context.Context, windowID string) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// MinimizeWindow minimizes a window
func (e *BaseEnvironment) MinimizeWindow(ctx context.Context, windowID string) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// MaximizeWindow maximizes a window
func (e *BaseEnvironment) MaximizeWindow(ctx context.Context, windowID string) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// RestoreWindow restores a window
func (e *BaseEnvironment) RestoreWindow(ctx context.Context, windowID string) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// MoveWindow moves a window to a new position
func (e *BaseEnvironment) MoveWindow(ctx context.Context, windowID string, x, y int) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// ResizeWindow resizes a window
func (e *BaseEnvironment) ResizeWindow(ctx context.Context, windowID string, width, height int) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// FocusWindow focuses a window
func (e *BaseEnvironment) FocusWindow(ctx context.Context, windowID string) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// ShowDesktop shows the desktop
func (e *BaseEnvironment) ShowDesktop(ctx context.Context) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// SendNotification sends a notification
func (e *BaseEnvironment) SendNotification(ctx context.Context, summary, body, icon string) (uint32, error) {
	// This should be overridden by specific implementations
	return 0, fmt.Errorf("not implemented")
}

// CloseNotification closes a notification
func (e *BaseEnvironment) CloseNotification(ctx context.Context, id uint32) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}

// TakeScreenshot takes a screenshot
func (e *BaseEnvironment) TakeScreenshot(ctx context.Context, fullScreen bool, delay int) (string, error) {
	// This should be overridden by specific implementations
	return "", fmt.Errorf("not implemented")
}

// GetClipboardText gets the text from the clipboard
func (e *BaseEnvironment) GetClipboardText(ctx context.Context) (string, error) {
	// This should be overridden by specific implementations
	return "", fmt.Errorf("not implemented")
}

// SetClipboardText sets the text in the clipboard
func (e *BaseEnvironment) SetClipboardText(ctx context.Context, text string) error {
	// This should be overridden by specific implementations
	return fmt.Errorf("not implemented")
}
