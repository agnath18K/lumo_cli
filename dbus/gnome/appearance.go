package gnome

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/agnath18K/lumo/internal/core"
)

// GNOME appearance-related DBus service names and interfaces
const (
	// GSettings is the GSettings service
	GSettings = "org.gnome.Settings"
	// GSettingsPath is the GSettings object path
	GSettingsPath = "/org/gnome/Settings"
	// GSettingsInterface is the GSettings interface
	GSettingsInterface = "org.gnome.Settings"

	// GSettingsSchemaDesktopInterface is the schema for desktop interface settings
	GSettingsSchemaDesktopInterface = "org.gnome.desktop.interface"
	// GSettingsSchemaDesktopBackground is the schema for desktop background settings
	GSettingsSchemaDesktopBackground = "org.gnome.desktop.background"
)

// executeAppearanceCommand executes an appearance management command
func (e *Environment) executeAppearanceCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Action {
	case "set-theme":
		theme := cmd.Target
		if theme == "" {
			return nil, fmt.Errorf("theme name is required")
		}
		if err := e.SetGtkTheme(ctx, theme); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set GTK theme to: %s", theme),
			Success: true,
		}, nil
	case "set-dark-mode":
		// Convert target to boolean
		enable := true
		if cmd.Target == "false" || cmd.Target == "off" || cmd.Target == "0" {
			enable = false
		}

		// Set color scheme based on dark mode preference
		colorScheme := "prefer-dark"
		if !enable {
			colorScheme = "prefer-light"
		}

		// Use gsettings to set the color scheme
		if err := e.setGSetting(GSettingsSchemaDesktopInterface, "color-scheme", colorScheme); err != nil {
			return nil, err
		}

		return &core.Result{
			Output:  fmt.Sprintf("Set dark mode to: %v", enable),
			Success: true,
		}, nil
	case "set-background":
		imagePath := cmd.Target
		if imagePath == "" {
			return nil, fmt.Errorf("background image path is required")
		}
		if err := e.SetDesktopBackground(ctx, imagePath); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set desktop background to: %s", imagePath),
			Success: true,
		}, nil
	case "set-accent-color":
		color := cmd.Target
		if color == "" {
			return nil, fmt.Errorf("accent color is required")
		}
		if err := e.SetAccentColor(ctx, color); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set accent color to: %s", color),
			Success: true,
		}, nil
	case "set-icon-theme":
		theme := cmd.Target
		if theme == "" {
			return nil, fmt.Errorf("icon theme name is required")
		}
		if err := e.SetIconTheme(ctx, theme); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set icon theme to: %s", theme),
			Success: true,
		}, nil
	case "get-theme":
		theme, err := e.GetCurrentTheme(ctx)
		if err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Current GTK theme: %s", theme),
			Success: true,
			Data: map[string]any{
				"theme": theme,
			},
		}, nil
	case "get-background":
		background, err := e.GetCurrentBackground(ctx)
		if err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Current desktop background: %s", background),
			Success: true,
			Data: map[string]any{
				"background": background,
			},
		}, nil
	case "get-icon-theme":
		theme, err := e.GetCurrentIconTheme(ctx)
		if err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Current icon theme: %s", theme),
			Success: true,
			Data: map[string]any{
				"icon_theme": theme,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported appearance action: %s", cmd.Action)
	}
}

// SetGtkTheme sets the GTK theme
func (e *Environment) SetGtkTheme(ctx context.Context, theme string) error {
	// Use gsettings to set the GTK theme
	if err := e.setGSetting(GSettingsSchemaDesktopInterface, "gtk-theme", theme); err != nil {
		return fmt.Errorf("failed to set GTK theme: %w", err)
	}
	return nil
}

// SetDesktopBackground sets the desktop background image
func (e *Environment) SetDesktopBackground(ctx context.Context, imagePath string) error {
	// Verify the image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("background image does not exist: %s", imagePath)
	}

	// Convert to absolute path if needed
	if !filepath.IsAbs(imagePath) {
		absPath, err := filepath.Abs(imagePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		imagePath = absPath
	}

	// Format the path as a URI
	imageURI := fmt.Sprintf("file://%s", imagePath)

	// Use gsettings to set the desktop background
	if err := e.setGSetting(GSettingsSchemaDesktopBackground, "picture-uri", imageURI); err != nil {
		return fmt.Errorf("failed to set desktop background: %w", err)
	}

	// Also set the dark mode background for GNOME 42+
	if err := e.setGSetting(GSettingsSchemaDesktopBackground, "picture-uri-dark", imageURI); err != nil {
		// This might fail on older GNOME versions, so we'll just log it
		fmt.Printf("Warning: Failed to set dark mode background: %v\n", err)
	}

	return nil
}

// SetAccentColor sets the accent color
func (e *Environment) SetAccentColor(ctx context.Context, color string) error {
	// Use gsettings to set the accent color (GNOME 42+)
	if err := e.setGSetting(GSettingsSchemaDesktopInterface, "accent-color", color); err != nil {
		return fmt.Errorf("failed to set accent color (this may not be supported in your GNOME version): %w", err)
	}
	return nil
}

// SetIconTheme sets the icon theme
func (e *Environment) SetIconTheme(ctx context.Context, theme string) error {
	// Use gsettings to set the icon theme
	if err := e.setGSetting(GSettingsSchemaDesktopInterface, "icon-theme", theme); err != nil {
		return fmt.Errorf("failed to set icon theme: %w", err)
	}
	return nil
}

// GetCurrentTheme gets the current GTK theme
func (e *Environment) GetCurrentTheme(ctx context.Context) (string, error) {
	// Use gsettings to get the current GTK theme
	theme, err := e.getGSetting(GSettingsSchemaDesktopInterface, "gtk-theme")
	if err != nil {
		return "", fmt.Errorf("failed to get current GTK theme: %w", err)
	}
	return theme, nil
}

// GetCurrentBackground gets the current desktop background
func (e *Environment) GetCurrentBackground(ctx context.Context) (string, error) {
	// Use gsettings to get the current desktop background
	background, err := e.getGSetting(GSettingsSchemaDesktopBackground, "picture-uri")
	if err != nil {
		return "", fmt.Errorf("failed to get current desktop background: %w", err)
	}

	// Convert from URI format to path
	background = strings.TrimPrefix(background, "file://")

	return background, nil
}

// GetCurrentIconTheme gets the current icon theme
func (e *Environment) GetCurrentIconTheme(ctx context.Context) (string, error) {
	// Use gsettings to get the current icon theme
	theme, err := e.getGSetting(GSettingsSchemaDesktopInterface, "icon-theme")
	if err != nil {
		return "", fmt.Errorf("failed to get current icon theme: %w", err)
	}
	return theme, nil
}

// setGSetting sets a GSettings value
func (e *Environment) setGSetting(schema, key, value string) error {
	// Use the gsettings command-line tool to set the value
	cmd := fmt.Sprintf("gsettings set %s %s '%s'", schema, key, value)
	output, err := e.runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to set gsettings value: %w (output: %s)", err, output)
	}
	return nil
}

// getGSetting gets a GSettings value
func (e *Environment) getGSetting(schema, key string) (string, error) {
	// Use the gsettings command-line tool to get the value
	cmd := fmt.Sprintf("gsettings get %s %s", schema, key)
	output, err := e.runCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to get gsettings value: %w", err)
	}

	// Clean up the output (remove quotes, newlines, etc.)
	output = strings.TrimSpace(output)
	output = strings.Trim(output, "'\"")

	return output, nil
}

// runCommand runs a shell command and returns its output
func (e *Environment) runCommand(cmd string) (string, error) {
	// Run the command using the system's shell
	command := exec.Command("sh", "-c", cmd)
	output, err := command.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
