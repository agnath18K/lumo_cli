package executor

import (
	"context"
	"fmt"

	"github.com/agnath18K/lumo/dbus/gnome"
	"github.com/agnath18K/lumo/internal/assistant"
	"github.com/agnath18K/lumo/internal/core"
	"github.com/agnath18K/lumo/internal/desktop"
	"github.com/agnath18K/lumo/pkg/nlp"
)

// executeDesktopCommand executes a desktop command
func (e *Executor) executeDesktopCommand(cmd *nlp.Command) (*Result, error) {
	// Create a desktop environment factory
	factory := desktop.NewFactory()

	// Register desktop environments
	registerDesktopEnvironments(factory)

	// Create a desktop assistant with AI capabilities
	var desktopAssistant *assistant.Assistant
	if e.aiClient != nil {
		// Create an AI client for the desktop assistant
		fmt.Println("DEBUG: AI client is available, creating AI-enabled desktop assistant")
		aiClient := assistant.NewAIClient(e.aiClient)
		desktopAssistant = assistant.NewAssistantWithAI(factory, aiClient)
	} else {
		// Create a regular desktop assistant without AI
		fmt.Println("DEBUG: AI client is not available, creating regular desktop assistant")
		desktopAssistant = assistant.NewAssistant(factory)
	}

	// Create a context
	ctx := context.Background()

	// Process the command
	result, err := desktopAssistant.ProcessCommand(ctx, cmd.Intent)
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Desktop Error: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Format the result
	output := formatDesktopResult(result)

	return &Result{
		Output:     output,
		IsError:    !result.Success,
		CommandRun: cmd.RawInput,
	}, nil
}

// registerDesktopEnvironments registers all available desktop environments
func registerDesktopEnvironments(factory *desktop.Factory) {
	// Register GNOME environment
	gnomeEnv, err := createGnomeEnvironment()
	if err == nil {
		factory.RegisterEnvironment(gnomeEnv)
	}

	// Register KDE environment
	kdeEnv, err := createKdeEnvironment()
	if err == nil {
		factory.RegisterEnvironment(kdeEnv)
	}

	// Register XFCE environment
	xfceEnv, err := createXfceEnvironment()
	if err == nil {
		factory.RegisterEnvironment(xfceEnv)
	}
}

// createGnomeEnvironment creates a GNOME desktop environment
func createGnomeEnvironment() (core.DesktopEnvironment, error) {
	// Import the GNOME package dynamically to avoid circular imports
	gnomeEnv, err := createEnvironment("gnome")
	if err != nil {
		return nil, err
	}
	return gnomeEnv, nil
}

// createKdeEnvironment creates a KDE desktop environment
func createKdeEnvironment() (core.DesktopEnvironment, error) {
	// Import the KDE package dynamically to avoid circular imports
	kdeEnv, err := createEnvironment("kde")
	if err != nil {
		return nil, err
	}
	return kdeEnv, nil
}

// createXfceEnvironment creates an XFCE desktop environment
func createXfceEnvironment() (core.DesktopEnvironment, error) {
	// Import the XFCE package dynamically to avoid circular imports
	xfceEnv, err := createEnvironment("xfce")
	if err != nil {
		return nil, err
	}
	return xfceEnv, nil
}

// createEnvironment creates a desktop environment by name
func createEnvironment(name string) (core.DesktopEnvironment, error) {
	switch name {
	case "gnome":
		// Import the GNOME package
		gnomeEnv, err := createGnomeEnvironmentImpl()
		if err != nil {
			return nil, err
		}
		return gnomeEnv, nil
	case "kde":
		// Import the KDE package
		kdeEnv, err := createKdeEnvironmentImpl()
		if err != nil {
			return nil, err
		}
		return kdeEnv, nil
	case "xfce":
		// Import the XFCE package
		xfceEnv, err := createXfceEnvironmentImpl()
		if err != nil {
			return nil, err
		}
		return xfceEnv, nil
	default:
		return nil, fmt.Errorf("unknown desktop environment: %s", name)
	}
}

// createGnomeEnvironmentImpl creates a GNOME desktop environment implementation
func createGnomeEnvironmentImpl() (core.DesktopEnvironment, error) {
	// Import the GNOME package dynamically
	gnomeEnv, err := createGnomeEnvironmentDynamic()
	if err != nil {
		return nil, err
	}
	return gnomeEnv, nil
}

// createGnomeEnvironmentDynamic creates a GNOME desktop environment dynamically
func createGnomeEnvironmentDynamic() (core.DesktopEnvironment, error) {
	// Import the GNOME package
	gnomeEnv, err := createGnomeEnvironmentFromPackage()
	if err != nil {
		return nil, err
	}
	return gnomeEnv, nil
}

// createGnomeEnvironmentFromPackage creates a GNOME desktop environment from the package
func createGnomeEnvironmentFromPackage() (core.DesktopEnvironment, error) {
	// Import the GNOME package
	// This is where we would import the GNOME package and create a GNOME environment
	// For now, we'll use a direct import
	gnomeEnv, err := gnome.NewEnvironment()
	if err != nil {
		return nil, err
	}
	return gnomeEnv, nil
}

// createKdeEnvironmentImpl creates a KDE desktop environment implementation
func createKdeEnvironmentImpl() (core.DesktopEnvironment, error) {
	// Import the KDE package dynamically
	kdeEnv, err := createKdeEnvironmentDynamic()
	if err != nil {
		return nil, err
	}
	return kdeEnv, nil
}

// createKdeEnvironmentDynamic creates a KDE desktop environment dynamically
func createKdeEnvironmentDynamic() (core.DesktopEnvironment, error) {
	// Import the KDE package
	kdeEnv, err := createKdeEnvironmentFromPackage()
	if err != nil {
		return nil, err
	}
	return kdeEnv, nil
}

// createKdeEnvironmentFromPackage creates a KDE desktop environment from the package
func createKdeEnvironmentFromPackage() (core.DesktopEnvironment, error) {
	return nil, fmt.Errorf("not implemented")
}

// createXfceEnvironmentImpl creates an XFCE desktop environment implementation
func createXfceEnvironmentImpl() (core.DesktopEnvironment, error) {
	// Import the XFCE package dynamically
	xfceEnv, err := createXfceEnvironmentDynamic()
	if err != nil {
		return nil, err
	}
	return xfceEnv, nil
}

// createXfceEnvironmentDynamic creates an XFCE desktop environment dynamically
func createXfceEnvironmentDynamic() (core.DesktopEnvironment, error) {
	// Import the XFCE package
	return nil, fmt.Errorf("not implemented")
}

// formatDesktopResult formats a desktop command result
func formatDesktopResult(result *core.Result) string {
	// Format the result
	if result.Success {
		return fmt.Sprintf("✅ %s", result.Output)
	}
	return fmt.Sprintf("❌ %s", result.Error)
}
