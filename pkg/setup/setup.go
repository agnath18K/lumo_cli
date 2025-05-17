package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/agnath18K/lumo/pkg/config"
)

// APIKeySetup handles interactive setup of API keys
type APIKeySetup struct {
	config *config.Config
	reader *bufio.Reader
}

// NewAPIKeySetup creates a new API key setup instance
func NewAPIKeySetup(cfg *config.Config) *APIKeySetup {
	return &APIKeySetup{
		config: cfg,
		reader: bufio.NewReader(os.Stdin),
	}
}

// CheckAndSetupAPIKeys checks if API keys are configured and runs setup if needed
// Returns true if setup was performed, false otherwise
func (s *APIKeySetup) CheckAndSetupAPIKeys() (bool, error) {
	// Check if we need to set up API keys
	needsSetup := false

	// If using Gemini and no API key is configured
	if s.config.AIProvider == "gemini" && s.config.GeminiAPIKey == "" {
		needsSetup = true
	}

	// If using OpenAI and no API key is configured
	if s.config.AIProvider == "openai" && s.config.OpenAIAPIKey == "" {
		needsSetup = true
	}

	// If no setup is needed, return
	if !needsSetup {
		return false, nil
	}

	// Run the interactive setup
	fmt.Println("\n🔑 Looks like you're missing an API key for your AI adventures!")
	fmt.Println("Let's get you set up so we can start having some fun together.")

	// Display disclaimer
	fmt.Println("\n⚠️ DISCLAIMER ⚠️")
	fmt.Println("Lumo is designed to help with basic terminal commands only, not for coding tasks.")
	fmt.Println("Think of me as your friendly terminal sidekick, not your programming buddy!")
	fmt.Println("I'll help you find and understand commands, but I won't write your next app. 😉")

	// Setup based on the current provider
	if s.config.AIProvider == "gemini" {
		if err := s.setupGeminiAPIKey(); err != nil {
			return true, err
		}
	} else if s.config.AIProvider == "openai" {
		if err := s.setupOpenAIAPIKey(); err != nil {
			return true, err
		}
	}

	// Save the updated configuration
	if err := s.config.Save(); err != nil {
		return true, fmt.Errorf("error saving configuration: %w", err)
	}

	fmt.Println("\n✅ Great! Your API key has been saved. Let's get started with your query!")
	fmt.Println("(You can always update your API keys by editing ~/.config/lumo/config.json)")

	return true, nil
}

// setupGeminiAPIKey guides the user through setting up a Gemini API key
func (s *APIKeySetup) setupGeminiAPIKey() error {
	fmt.Println("📡 Setting up Google Gemini API")
	fmt.Println("--------------------------------")
	fmt.Println("To use Gemini, you'll need an API key from Google AI Studio.")
	fmt.Println("You can get one for free at: https://aistudio.google.com/apikey")
	fmt.Println("\nOnce you have your key, paste it here (don't worry, I won't peek! 👀)")

	// Get API key from user
	fmt.Print("\nGemini API Key: ")
	apiKey, err := s.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Trim whitespace and newlines
	apiKey = strings.TrimSpace(apiKey)

	// Validate the API key (basic validation)
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Save the API key to the configuration
	s.config.GeminiAPIKey = apiKey

	return nil
}

// setupOpenAIAPIKey guides the user through setting up an OpenAI API key
func (s *APIKeySetup) setupOpenAIAPIKey() error {
	fmt.Println("🐦 Setting up OpenAI API")
	fmt.Println("-----------------------")
	fmt.Println("To use OpenAI's models, you'll need an API key from OpenAI.")
	fmt.Println("You can get one at: https://platform.openai.com/api-keys")
	fmt.Println("\nOnce you have your key, paste it here (I promise it stays between us! 🤐)")

	// Get API key from user
	fmt.Print("\nOpenAI API Key: ")
	apiKey, err := s.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Trim whitespace and newlines
	apiKey = strings.TrimSpace(apiKey)

	// Validate the API key (basic validation)
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Save the API key to the configuration
	s.config.OpenAIAPIKey = apiKey

	return nil
}

// SwitchProvider allows the user to switch between AI providers
func (s *APIKeySetup) SwitchProvider() error {
	currentProvider := s.config.AIProvider

	fmt.Println("\n🔄 Switch AI Provider")
	fmt.Println("-------------------")
	fmt.Printf("You're currently using: %s\n\n", currentProvider)
	fmt.Println("Available providers:")
	fmt.Println("1. Gemini (Google)")
	fmt.Println("2. OpenAI (GPT)")
	fmt.Print("\nSelect a provider (1-2): ")

	// Get user selection
	selection, err := s.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Trim whitespace and newlines
	selection = strings.TrimSpace(selection)

	// Process selection
	switch selection {
	case "1":
		s.config.AIProvider = "gemini"
		if s.config.GeminiAPIKey == "" {
			if err := s.setupGeminiAPIKey(); err != nil {
				return err
			}
		}
	case "2":
		s.config.AIProvider = "openai"
		if s.config.OpenAIAPIKey == "" {
			if err := s.setupOpenAIAPIKey(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid selection: %s", selection)
	}

	// Save the updated configuration
	if err := s.config.Save(); err != nil {
		return fmt.Errorf("error saving configuration: %w", err)
	}

	fmt.Printf("\n✅ Switched to %s as your AI provider!\n", s.config.AIProvider)

	// Display disclaimer
	fmt.Println("\n⚠️ DISCLAIMER ⚠️")
	fmt.Println("Lumo is designed to help with basic terminal commands only, not for coding tasks.")
	fmt.Println("Think of me as your friendly terminal sidekick, not your programming buddy!")
	fmt.Println("I'll help you find and understand commands, but I won't write your next app. 😉")

	return nil
}
