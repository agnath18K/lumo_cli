package gnome

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ShowDesktop shows the desktop
func (e *Environment) ShowDesktop(ctx context.Context) error {
	// Call the DBus method to show the desktop
	_, err := e.sessionHandler.Call(
		Shell,
		ShellPath,
		ShellInterface,
		"ShowDesktop",
	)
	if err != nil {
		return fmt.Errorf("failed to show desktop: %w", err)
	}

	return nil
}

// SendNotification sends a notification
func (e *Environment) SendNotification(ctx context.Context, summary, body, icon string) (uint32, error) {
	// Call the DBus method to send a notification
	result, err := e.sessionHandler.Call(
		Notifications,
		NotificationsPath,
		NotificationsInterface,
		"Notify",
		"Lumo",                   // Application name
		uint32(0),                // Replaces ID (0 = new notification)
		icon,                     // Icon
		summary,                  // Summary
		body,                     // Body
		[]string{},               // Actions
		map[string]interface{}{}, // Hints
		int32(5000),              // Timeout (5 seconds)
	)
	if err != nil {
		return 0, fmt.Errorf("failed to send notification: %w", err)
	}

	// Parse the result
	if len(result) > 0 {
		if id, ok := result[0].(uint32); ok {
			return id, nil
		}
	}

	return 0, fmt.Errorf("failed to get notification ID")
}

// CloseNotification closes a notification
func (e *Environment) CloseNotification(ctx context.Context, id uint32) error {
	// Call the DBus method to close the notification
	_, err := e.sessionHandler.Call(
		Notifications,
		NotificationsPath,
		NotificationsInterface,
		"CloseNotification",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to close notification: %w", err)
	}

	return nil
}

// TakeScreenshot takes a screenshot
func (e *Environment) TakeScreenshot(ctx context.Context, fullScreen bool, delay int) (string, error) {
	// Create a temporary file to store the screenshot
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create a timestamp for the filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("screenshot-%s.png", timestamp)
	screenshotDir := filepath.Join(homeDir, "Pictures")
	screenshotPath := filepath.Join(screenshotDir, filename)

	// Ensure the directory exists
	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Try to use the DBus method to take a screenshot
	var result []interface{}
	if fullScreen {
		result, err = e.sessionHandler.Call(
			Screenshot,
			ScreenshotPath,
			ScreenshotInterface,
			"Screenshot",
			true,  // Include cursor
			delay, // Delay in seconds
			screenshotPath,
		)
	} else {
		result, err = e.sessionHandler.Call(
			Screenshot,
			ScreenshotPath,
			ScreenshotInterface,
			"ScreenshotArea",
			0, 0, -1, -1, // x, y, width, height (-1 = full screen)
			true,  // Include cursor
			delay, // Delay in seconds
			screenshotPath,
		)
	}

	if err != nil {
		// Fallback to using the command line
		var cmd *exec.Cmd
		if fullScreen {
			cmd = exec.Command("gnome-screenshot", "-f", screenshotPath)
		} else {
			cmd = exec.Command("gnome-screenshot", "-a", "-f", screenshotPath)
		}

		if delay > 0 {
			cmd.Args = append(cmd.Args, "-d", fmt.Sprintf("%d", delay))
		}

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to take screenshot: %w", err)
		}
	} else {
		// Parse the result
		if len(result) > 0 {
			if success, ok := result[0].(bool); ok && !success {
				return "", fmt.Errorf("failed to take screenshot")
			}
		}
	}

	return screenshotPath, nil
}

// GetClipboardText gets the text from the clipboard
func (e *Environment) GetClipboardText(ctx context.Context) (string, error) {
	// Try to use the DBus method to get the clipboard text
	result, err := e.sessionHandler.Call(
		Clipboard,
		ClipboardPath,
		ClipboardInterface,
		"GetClipboardText",
	)
	if err != nil {
		// Fallback to using the command line
		cmd := exec.Command("xclip", "-o", "-selection", "clipboard")
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to get clipboard text: %w", err)
		}
		return string(output), nil
	}

	// Parse the result
	if len(result) > 0 {
		if text, ok := result[0].(string); ok {
			return text, nil
		}
	}

	return "", fmt.Errorf("failed to get clipboard text")
}

// SetClipboardText sets the text in the clipboard
func (e *Environment) SetClipboardText(ctx context.Context, text string) error {
	// Try to use the DBus method to set the clipboard text
	_, err := e.sessionHandler.Call(
		Clipboard,
		ClipboardPath,
		ClipboardInterface,
		"SetClipboardText",
		text,
	)
	if err != nil {
		// Fallback to using the command line
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set clipboard text: %w", err)
		}
	}

	return nil
}
