package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// generateFlutterProject creates a new Flutter project
func generateFlutterProject(stateManagement string, options map[string]string) (string, error) {
	// Get project name from options or use a default
	projectName := options["name"]
	if projectName == "" {
		projectName = "my_flutter_app"
	}
	
	// Check if Flutter is installed
	if err := checkFlutterInstalled(); err != nil {
		return "", err
	}
	
	// Create the project using Flutter CLI
	if err := createBaseFlutterProject(projectName); err != nil {
		return "", err
	}
	
	// Set up the project structure based on state management
	switch strings.ToLower(stateManagement) {
	case "bloc":
		if err := setupBlocArchitecture(projectName); err != nil {
			return "", err
		}
	case "provider":
		if err := setupProviderArchitecture(projectName); err != nil {
			return "", err
		}
	case "riverpod":
		if err := setupRiverpodArchitecture(projectName); err != nil {
			return "", err
		}
	default:
		// Default to a basic MVVM structure
		if err := setupBasicMVVMArchitecture(projectName); err != nil {
			return "", err
		}
	}
	
	return fmt.Sprintf("✅ Flutter project '%s' created successfully with %s architecture!", 
		projectName, 
		getArchitectureName(stateManagement)), nil
}

// checkFlutterInstalled verifies that Flutter is installed
func checkFlutterInstalled() error {
	cmd := exec.Command("flutter", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Flutter is not installed or not in PATH. Please install Flutter first: https://flutter.dev/docs/get-started/install")
	}
	return nil
}

// createBaseFlutterProject creates a new Flutter project using the Flutter CLI
func createBaseFlutterProject(name string) error {
	cmd := exec.Command("flutter", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// setupBlocArchitecture sets up a Flutter project with BLoC architecture
func setupBlocArchitecture(projectPath string) error {
	// Create directories for BLoC architecture
	dirs := []string{
		"lib/blocs",
		"lib/repositories",
		"lib/models",
		"lib/screens",
		"lib/widgets",
		"lib/services",
		"lib/utils",
	}
	
	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}
	
	// Add BLoC dependencies to pubspec.yaml
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if err := addDependenciesToPubspec(pubspecPath, []string{
		"flutter_bloc: ^8.1.3",
		"equatable: ^2.0.5",
	}); err != nil {
		return err
	}
	
	// Create sample BLoC files
	if err := createSampleBlocFiles(projectPath); err != nil {
		return err
	}
	
	// Update main.dart to use BLoC
	if err := updateMainDartForBloc(projectPath); err != nil {
		return err
	}
	
	return nil
}

// setupProviderArchitecture sets up a Flutter project with Provider architecture
func setupProviderArchitecture(projectPath string) error {
	// Create directories for Provider architecture
	dirs := []string{
		"lib/providers",
		"lib/models",
		"lib/screens",
		"lib/widgets",
		"lib/services",
		"lib/utils",
	}
	
	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}
	
	// Add Provider dependencies to pubspec.yaml
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if err := addDependenciesToPubspec(pubspecPath, []string{
		"provider: ^6.0.5",
	}); err != nil {
		return err
	}
	
	// Create sample Provider files
	if err := createSampleProviderFiles(projectPath); err != nil {
		return err
	}
	
	// Update main.dart to use Provider
	if err := updateMainDartForProvider(projectPath); err != nil {
		return err
	}
	
	return nil
}

// setupRiverpodArchitecture sets up a Flutter project with Riverpod architecture
func setupRiverpodArchitecture(projectPath string) error {
	// Create directories for Riverpod architecture
	dirs := []string{
		"lib/providers",
		"lib/models",
		"lib/screens",
		"lib/widgets",
		"lib/services",
		"lib/utils",
	}
	
	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}
	
	// Add Riverpod dependencies to pubspec.yaml
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if err := addDependenciesToPubspec(pubspecPath, []string{
		"flutter_riverpod: ^2.4.0",
		"riverpod_annotation: ^2.1.5",
	}); err != nil {
		return err
	}
	
	// Create sample Riverpod files
	if err := createSampleRiverpodFiles(projectPath); err != nil {
		return err
	}
	
	// Update main.dart to use Riverpod
	if err := updateMainDartForRiverpod(projectPath); err != nil {
		return err
	}
	
	return nil
}

// setupBasicMVVMArchitecture sets up a Flutter project with basic MVVM architecture
func setupBasicMVVMArchitecture(projectPath string) error {
	// Create directories for basic MVVM architecture
	dirs := []string{
		"lib/models",
		"lib/views",
		"lib/viewmodels",
		"lib/services",
		"lib/utils",
	}
	
	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}
	
	// Create sample MVVM files
	if err := createSampleMVVMFiles(projectPath); err != nil {
		return err
	}
	
	// Update main.dart to use MVVM
	if err := updateMainDartForMVVM(projectPath); err != nil {
		return err
	}
	
	return nil
}

// getArchitectureName returns a user-friendly name for the architecture
func getArchitectureName(stateManagement string) string {
	switch strings.ToLower(stateManagement) {
	case "bloc":
		return "BLoC"
	case "provider":
		return "Provider"
	case "riverpod":
		return "Riverpod"
	default:
		return "MVVM"
	}
}
