package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agnath18/lumo/pkg/config"
)

// TestDefaultConfig tests that the default configuration is created correctly
func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	// Check default values
	if cfg.AIProvider != "gemini" {
		t.Errorf("Expected default AIProvider to be 'gemini', got '%s'", cfg.AIProvider)
	}

	if cfg.GeminiModel != "gemini-2.0-flash-lite" {
		t.Errorf("Expected default GeminiModel to be 'gemini-2.0-flash-lite', got '%s'", cfg.GeminiModel)
	}

	if cfg.OpenAIModel != "gpt-3.5-turbo" {
		t.Errorf("Expected default OpenAIModel to be 'gpt-3.5-turbo', got '%s'", cfg.OpenAIModel)
	}

	if cfg.MaxHistorySize != 1000 {
		t.Errorf("Expected default MaxHistorySize to be 1000, got %d", cfg.MaxHistorySize)
	}

	if !cfg.EnableLogging {
		t.Error("Expected default EnableLogging to be true")
	}

	if cfg.EnableShellInInteractive {
		t.Error("Expected default EnableShellInInteractive to be false")
	}

	if !cfg.EnableAgentMode {
		t.Error("Expected default EnableAgentMode to be true")
	}

	if !cfg.EnableAgentREPL {
		t.Error("Expected default EnableAgentREPL to be true")
	}

	if !cfg.AgentConfirmBeforeExecution {
		t.Error("Expected default AgentConfirmBeforeExecution to be true")
	}

	if cfg.AgentMaxSteps != 10 {
		t.Errorf("Expected default AgentMaxSteps to be 10, got %d", cfg.AgentMaxSteps)
	}

	if cfg.AgentSafetyLevel != "medium" {
		t.Errorf("Expected default AgentSafetyLevel to be 'medium', got '%s'", cfg.AgentSafetyLevel)
	}

	if !cfg.EnablePipeProcessing {
		t.Error("Expected default EnablePipeProcessing to be true")
	}

	if !cfg.EnableSystemHealth {
		t.Error("Expected default EnableSystemHealth to be true")
	}

	if !cfg.EnableSystemReport {
		t.Error("Expected default EnableSystemReport to be true")
	}

	if cfg.Debug {
		t.Error("Expected default Debug to be false")
	}
}

// TestSaveAndLoad tests saving and loading configuration
func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lumo-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up a temporary home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create a config with non-default values
	cfg := config.DefaultConfig()
	cfg.AIProvider = "openai"
	cfg.GeminiAPIKey = "test-gemini-key"
	cfg.OpenAIAPIKey = "test-openai-key"
	cfg.MaxHistorySize = 500
	cfg.EnableLogging = false
	cfg.EnableShellInInteractive = true
	cfg.EnableAgentMode = false
	cfg.AgentMaxSteps = 5
	cfg.Debug = true

	// Save the config
	err = cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check that the config file was created
	configPath := filepath.Join(tempDir, ".config", "lumo", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Load a new config
	newCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that the values match
	if newCfg.AIProvider != cfg.AIProvider {
		t.Errorf("Expected AIProvider to be '%s', got '%s'", cfg.AIProvider, newCfg.AIProvider)
	}

	if newCfg.GeminiAPIKey != cfg.GeminiAPIKey {
		t.Errorf("Expected GeminiAPIKey to be '%s', got '%s'", cfg.GeminiAPIKey, newCfg.GeminiAPIKey)
	}

	if newCfg.OpenAIAPIKey != cfg.OpenAIAPIKey {
		t.Errorf("Expected OpenAIAPIKey to be '%s', got '%s'", cfg.OpenAIAPIKey, newCfg.OpenAIAPIKey)
	}

	if newCfg.MaxHistorySize != cfg.MaxHistorySize {
		t.Errorf("Expected MaxHistorySize to be %d, got %d", cfg.MaxHistorySize, newCfg.MaxHistorySize)
	}

	if newCfg.EnableLogging != cfg.EnableLogging {
		t.Errorf("Expected EnableLogging to be %v, got %v", cfg.EnableLogging, newCfg.EnableLogging)
	}

	if newCfg.EnableShellInInteractive != cfg.EnableShellInInteractive {
		t.Errorf("Expected EnableShellInInteractive to be %v, got %v", cfg.EnableShellInInteractive, newCfg.EnableShellInInteractive)
	}

	if newCfg.EnableAgentMode != cfg.EnableAgentMode {
		t.Errorf("Expected EnableAgentMode to be %v, got %v", cfg.EnableAgentMode, newCfg.EnableAgentMode)
	}

	if newCfg.AgentMaxSteps != cfg.AgentMaxSteps {
		t.Errorf("Expected AgentMaxSteps to be %d, got %d", cfg.AgentMaxSteps, newCfg.AgentMaxSteps)
	}

	if newCfg.Debug != cfg.Debug {
		t.Errorf("Expected Debug to be %v, got %v", cfg.Debug, newCfg.Debug)
	}
}

// TestEnvironmentVariables tests that environment variables override config values
func TestEnvironmentVariables(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lumo-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up a temporary home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Set environment variables
	os.Setenv("GEMINI_API_KEY", "env-gemini-key")
	os.Setenv("OPENAI_API_KEY", "env-openai-key")
	defer os.Unsetenv("GEMINI_API_KEY")
	defer os.Unsetenv("OPENAI_API_KEY")

	// Create a config with different API keys
	cfg := config.DefaultConfig()
	cfg.GeminiAPIKey = "file-gemini-key"
	cfg.OpenAIAPIKey = "file-openai-key"

	// Save the config
	err = cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load a new config
	newCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that environment variables override file values
	if newCfg.GeminiAPIKey != "env-gemini-key" {
		t.Errorf("Expected GeminiAPIKey to be 'env-gemini-key', got '%s'", newCfg.GeminiAPIKey)
	}

	if newCfg.OpenAIAPIKey != "env-openai-key" {
		t.Errorf("Expected OpenAIAPIKey to be 'env-openai-key', got '%s'", newCfg.OpenAIAPIKey)
	}
}
