package gnome

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/agnath18K/lumo/internal/core"
	"github.com/godbus/dbus/v5"
)

// GetWindows returns a list of all windows
func (e *Environment) GetWindows(ctx context.Context) ([]core.Window, error) {
	fmt.Printf("DEBUG: Getting windows using wmctrl command\n")

	// Use wmctrl command to get window list
	cmd := exec.Command("wmctrl", "-l")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl: %v\n", err)
		return nil, fmt.Errorf("failed to list windows: %w", err)
	}

	// Parse the output
	lines := strings.Split(string(output), "\n")
	fmt.Printf("DEBUG: Got %d lines of output from wmctrl\n", len(lines))

	var windows []core.Window
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse the line (format: window_id desktop_id hostname window_title)
		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 4 {
			continue
		}

		// Extract window properties
		id := parts[0]
		title := parts[3]

		// Create window object
		window := core.Window{
			ID:          id,
			Title:       title,
			Application: "", // Not available from wmctrl -l
			Geometry: core.WindowGeometry{
				X:      0, // Not available from wmctrl -l
				Y:      0,
				Width:  0,
				Height: 0,
			},
			State: core.WindowState{
				Maximized:  false, // Not available from wmctrl -l
				Minimized:  false,
				Fullscreen: false,
				Active:     false,
			},
		}

		fmt.Printf("DEBUG: Window: ID=%s, Title=%s\n", window.ID, window.Title)
		windows = append(windows, window)
	}

	fmt.Printf("DEBUG: Parsed %d windows\n", len(windows))
	return windows, nil
}

// GetActiveWindow returns the currently active window
func (e *Environment) GetActiveWindow(ctx context.Context) (*core.Window, error) {
	// Call the DBus method to get the active window
	result, err := e.sessionHandler.Call(
		Shell,
		ShellPath,
		ShellInterface,
		"GetActiveWindow",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get active window: %w", err)
	}

	// Parse the result
	if len(result) > 0 {
		if windowData, ok := result[0].(map[string]dbus.Variant); ok {
			window := &core.Window{
				ID:          windowData["id"].Value().(string),
				Title:       windowData["title"].Value().(string),
				Application: windowData["app_id"].Value().(string),
				Geometry: core.WindowGeometry{
					X:      windowData["x"].Value().(int),
					Y:      windowData["y"].Value().(int),
					Width:  windowData["width"].Value().(int),
					Height: windowData["height"].Value().(int),
				},
				State: core.WindowState{
					Maximized:  windowData["maximized"].Value().(bool),
					Minimized:  windowData["minimized"].Value().(bool),
					Fullscreen: windowData["fullscreen"].Value().(bool),
					Active:     windowData["active"].Value().(bool),
				},
			}
			return window, nil
		}
	}

	return nil, fmt.Errorf("no active window found")
}

// GetRunningApplications returns a list of running applications
func (e *Environment) GetRunningApplications(ctx context.Context) ([]core.Application, error) {
	// Call the DBus method to get all running applications
	result, err := e.sessionHandler.Call(
		Shell,
		ShellPath,
		ShellInterface,
		"ListApplications",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	// Parse the result
	var applications []core.Application
	if len(result) > 0 {
		if appsData, ok := result[0].([]map[string]dbus.Variant); ok {
			for _, appData := range appsData {
				app := core.Application{
					ID:          appData["id"].Value().(string),
					Name:        appData["name"].Value().(string),
					Executable:  appData["executable"].Value().(string),
					DesktopFile: appData["desktop_file"].Value().(string),
					Running:     appData["running"].Value().(bool),
				}
				applications = append(applications, app)
			}
		}
	}

	return applications, nil
}

// LaunchApplication launches an application
func (e *Environment) LaunchApplication(ctx context.Context, appName string, args ...string) error {
	// Try to launch the application using DBus
	_, err := e.sessionHandler.Call(
		AppLauncher,
		AppLauncherPath,
		AppLauncherInterface,
		"LaunchApplication",
		appName,
		args,
	)
	if err != nil {
		// Fallback to using the command line
		cmd := exec.Command(appName, args...)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to launch application: %w", err)
		}
	}

	return nil
}

// CloseWindow closes a window
func (e *Environment) CloseWindow(ctx context.Context, windowID string) error {
	fmt.Printf("DEBUG: Closing window using wmctrl: %s\n", windowID)

	// Use wmctrl to close the window
	// The -c option closes the window gracefully
	cmd := exec.Command("wmctrl", "-c", windowID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to close window: %w", err)
	}

	fmt.Printf("DEBUG: Successfully closed window: %s\n", windowID)
	return nil
}

// MinimizeWindow minimizes a window
func (e *Environment) MinimizeWindow(ctx context.Context, windowID string) error {
	fmt.Printf("DEBUG: Minimizing window using wmctrl: %s\n", windowID)

	// Use wmctrl to minimize the window
	// The -r option selects the window, and -b add,hidden adds the hidden state
	cmd := exec.Command("wmctrl", "-r", windowID, "-b", "add,hidden")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to minimize window: %w", err)
	}

	fmt.Printf("DEBUG: Successfully minimized window: %s\n", windowID)
	return nil
}

// MaximizeWindow maximizes a window
func (e *Environment) MaximizeWindow(ctx context.Context, windowID string) error {
	fmt.Printf("DEBUG: Maximizing window using wmctrl: %s\n", windowID)

	// Use wmctrl to maximize the window
	// The -r option selects the window, and -b add,maximized_vert,maximized_horz adds both vertical and horizontal maximization
	cmd := exec.Command("wmctrl", "-r", windowID, "-b", "add,maximized_vert,maximized_horz")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to maximize window: %w", err)
	}

	fmt.Printf("DEBUG: Successfully maximized window: %s\n", windowID)
	return nil
}

// RestoreWindow restores a window
func (e *Environment) RestoreWindow(ctx context.Context, windowID string) error {
	fmt.Printf("DEBUG: Restoring window using wmctrl: %s\n", windowID)

	// Use wmctrl to restore the window
	// First, remove the hidden state to unminimize
	cmd := exec.Command("wmctrl", "-r", windowID, "-b", "remove,hidden")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl to unhide: %v, output: %s\n", err, string(output))
		// Continue anyway, as the window might not be hidden
	}

	// Then, remove maximized state
	cmd = exec.Command("wmctrl", "-r", windowID, "-b", "remove,maximized_vert,maximized_horz")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl to unmaximize: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to restore window: %w", err)
	}

	fmt.Printf("DEBUG: Successfully restored window: %s\n", windowID)
	return nil
}

// MoveWindow moves a window to a new position
func (e *Environment) MoveWindow(ctx context.Context, windowID string, x, y int) error {
	fmt.Printf("DEBUG: Moving window using wmctrl: %s to (%d, %d)\n", windowID, x, y)

	// Use wmctrl to move the window
	// The -e option changes the geometry of the window
	// Format: -e 0,x,y,-1,-1 (0 is gravity, -1 means don't change width/height)
	cmd := exec.Command("wmctrl", "-r", windowID, "-e", fmt.Sprintf("0,%d,%d,-1,-1", x, y))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to move window: %w", err)
	}

	fmt.Printf("DEBUG: Successfully moved window: %s\n", windowID)
	return nil
}

// ResizeWindow resizes a window
func (e *Environment) ResizeWindow(ctx context.Context, windowID string, width, height int) error {
	fmt.Printf("DEBUG: Resizing window using wmctrl: %s to %dx%d\n", windowID, width, height)

	// Use wmctrl to resize the window
	// The -e option changes the geometry of the window
	// Format: -e 0,-1,-1,width,height (0 is gravity, -1 means don't change x/y)
	cmd := exec.Command("wmctrl", "-r", windowID, "-e", fmt.Sprintf("0,-1,-1,%d,%d", width, height))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to resize window: %w", err)
	}

	fmt.Printf("DEBUG: Successfully resized window: %s\n", windowID)
	return nil
}

// FocusWindow focuses a window
func (e *Environment) FocusWindow(ctx context.Context, windowID string) error {
	fmt.Printf("DEBUG: Focusing window using wmctrl: %s\n", windowID)

	// Use wmctrl to focus the window
	// The -a option activates the window by switching to its desktop and raising it
	cmd := exec.Command("wmctrl", "-a", windowID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("DEBUG: Error running wmctrl: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to focus window: %w", err)
	}

	fmt.Printf("DEBUG: Successfully focused window: %s\n", windowID)
	return nil
}
