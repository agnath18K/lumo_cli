package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/agnath18K/lumo/pkg/ai"
	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/nlp"
)

// getCurrentModel returns the current model based on the provider
func getCurrentModel(cfg *config.Config) string {
	if cfg.AIProvider == "gemini" {
		return cfg.GeminiModel
	}
	return cfg.OpenAIModel
}

// executeConfigCommand handles configuration commands
func (e *Executor) executeConfigCommand(cmd *nlp.Command) (*Result, error) {
	// Split the command into parts
	parts := strings.Fields(cmd.Intent)
	if len(parts) == 0 {
		return &Result{
			Output: `
╭─────────────────── 🔧 Lumo Configuration ─────────────────╮

  Commands:
   • config:provider list           List available AI providers
   • config:provider show           Show current AI provider
   • config:provider set <provider> Set AI provider (gemini/openai/ollama)

   • config:model list              List available models
   • config:model show              Show current model
   • config:model set <model>       Set model for current provider

   • config:key show                Show current API key status
   • config:key set <provider> <key> Set API key for provider
   • config:key remove <provider>   Remove API key for provider

   • config:ollama show             Show current Ollama URL
   • config:ollama set <url>        Set Ollama URL
   • config:ollama test             Test connection to Ollama server

   • config:mode show               Show current input mode
   • config:mode ai                 Set AI-first mode (default)
   • config:mode command            Set command-first mode

   • config:server show             Show current server settings
   • config:server quiet on/off     Enable/disable server log messages

╰──────────────────────────────────────────────────────────╯
`,
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Handle different configuration commands
	switch parts[0] {
	case "provider":
		return e.handleProviderConfig(parts[1:], cmd)
	case "model":
		return e.handleModelConfig(parts[1:], cmd)
	case "key":
		return e.handleKeyConfig(parts[1:], cmd)
	case "ollama":
		return e.handleOllamaConfig(parts[1:], cmd)
	case "mode":
		return e.handleModeConfig(parts[1:], cmd)
	case "server":
		return e.handleServerConfig(parts[1:], cmd)
	default:
		return &Result{
			Output:     fmt.Sprintf("Unknown configuration command: %s\nUse 'config:' for help.", parts[0]),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}

// handleProviderConfig handles provider configuration commands
func (e *Executor) handleProviderConfig(args []string, cmd *nlp.Command) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			Output:     "Missing provider command. Use 'list', 'show', or 'set'.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	switch args[0] {
	case "list":
		// List available providers
		output := `
╭─────────────── 🐦 Available AI Providers ───────────────╮

  • gemini  (Google's Gemini AI models)
  • openai  (OpenAI's GPT models)
  • ollama  (Local Ollama models)

  Current provider: ` + e.config.AIProvider + `

╰──────────────────────────────────────────────────────────╯
`
		return &Result{
			Output:     output,
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "show":
		// Show current provider
		return &Result{
			Output:     fmt.Sprintf("Current AI provider: %s", e.config.AIProvider),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "set":
		// Set provider
		if len(args) < 2 {
			return &Result{
				Output:     "Missing provider name. Use 'gemini', 'openai', or 'ollama'.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		provider := strings.ToLower(args[1])
		if provider != "gemini" && provider != "openai" && provider != "ollama" {
			return &Result{
				Output:     fmt.Sprintf("Invalid provider: %s. Use 'gemini', 'openai', or 'ollama'.", provider),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Check if API key is set for the provider (not needed for Ollama)
		if provider == "gemini" && e.config.GeminiAPIKey == "" {
			return &Result{
				Output:     "No API key set for Gemini. Please set an API key first with 'config:key set gemini <key>'.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}
		if provider == "openai" && e.config.OpenAIAPIKey == "" {
			return &Result{
				Output:     "No API key set for OpenAI. Please set an API key first with 'config:key set openai <key>'.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}
		// Ollama doesn't need an API key, but we should check if the URL is accessible
		if provider == "ollama" {
			// Try to connect to the Ollama server
			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			_, err := client.Get(e.config.OllamaURL + "/api/tags")
			if err != nil {
				return &Result{
					Output:     fmt.Sprintf("Cannot connect to Ollama server at %s. Please make sure Ollama is running and accessible.", e.config.OllamaURL),
					IsError:    true,
					CommandRun: cmd.RawInput,
				}, nil
			}
		}

		// Set the provider
		e.config.AIProvider = provider

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Reinitialize the AI client with the new provider
		switch provider {
		case "gemini":
			e.aiClient = ai.NewGeminiClient(e.config.GeminiAPIKey, e.config.GeminiModel)
		case "ollama":
			e.aiClient = ai.NewOllamaClient(e.config.OllamaURL, e.config.OllamaModel)
		default: // Default to OpenAI
			e.aiClient = ai.NewOpenAIClient(e.config.OpenAIAPIKey, e.config.OpenAIModel)
		}

		return &Result{
			Output:     fmt.Sprintf("AI provider set to: %s", provider),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	default:
		return &Result{
			Output:     fmt.Sprintf("Unknown provider command: %s. Use 'show' or 'set'.", args[0]),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}

// handleModelConfig handles model configuration commands
func (e *Executor) handleModelConfig(args []string, cmd *nlp.Command) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			Output:     "Missing model command. Use 'list', 'show', or 'set'.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	switch args[0] {
	case "list":
		// Use the dedicated model list handler
		return e.handleModelList(cmd)
	case "unused":
		// This case is never used, just here to keep the old code structure
		var output string
		if false {
			output = `
╭─────────────── 🐦 Available Gemini Models ───────────────╮

  • gemini-2.0-flash-lite  (Fast, efficient for most queries)
  • gemini-2.0-flash       (Balanced performance and quality)
  • gemini-2.0-pro         (High quality, more capabilities)

  Current model: ` + e.config.GeminiModel + `

╰──────────────────────────────────────────────────────────╯
`
		} else {
			output = `
╭─────────────── 🐦 Available OpenAI Models ───────────────╮

  • gpt-3.5-turbo          (Fast, cost-effective)
  • gpt-4o                 (Advanced capabilities, slower)
  • gpt-4o-mini            (Balanced performance and quality)

  Current model: ` + e.config.OpenAIModel + `

╰──────────────────────────────────────────────────────────╯
`
		}
		return &Result{
			Output:     output,
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "show":
		// Show current model
		var currentModel string
		switch e.config.AIProvider {
		case "gemini":
			currentModel = e.config.GeminiModel
		case "ollama":
			currentModel = e.config.OllamaModel
		default: // OpenAI
			currentModel = e.config.OpenAIModel
		}
		return &Result{
			Output:     fmt.Sprintf("Current %s model: %s", e.config.AIProvider, currentModel),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "set":
		// Set model
		if len(args) < 2 {
			return &Result{
				Output:     "Missing model name. Use 'config:model list' to see available models.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		model := args[1]

		// Validate model based on provider
		switch e.config.AIProvider {
		case "gemini":
			validModels := []string{"gemini-2.0-flash-lite", "gemini-2.0-flash", "gemini-2.0-pro"}
			isValid := false
			for _, validModel := range validModels {
				if model == validModel {
					isValid = true
					break
				}
			}
			if !isValid {
				return &Result{
					Output:     fmt.Sprintf("Invalid Gemini model: %s. Use 'config:model list' to see available models.", model),
					IsError:    true,
					CommandRun: cmd.RawInput,
				}, nil
			}

			// Set the model
			e.config.GeminiModel = model

			// Reinitialize the AI client with the new model
			e.aiClient = ai.NewGeminiClient(e.config.GeminiAPIKey, e.config.GeminiModel)

		case "ollama":
			// For Ollama, we need to check if the model exists
			ollamaClient := ai.NewOllamaClient(e.config.OllamaURL, e.config.OllamaModel)
			models, err := ollamaClient.ListModels()
			if err != nil {
				return &Result{
					Output:     fmt.Sprintf("Error checking Ollama models: %v", err),
					IsError:    true,
					CommandRun: cmd.RawInput,
				}, nil
			}

			isValid := false
			for _, validModel := range models {
				if model == validModel {
					isValid = true
					break
				}
			}

			if !isValid {
				return &Result{
					Output:     fmt.Sprintf("Invalid or unavailable Ollama model: %s. Use 'config:model list' to see available models or 'ollama pull %s' to download it.", model, model),
					IsError:    true,
					CommandRun: cmd.RawInput,
				}, nil
			}

			// Set the model
			e.config.OllamaModel = model

			// Reinitialize the AI client with the new model
			e.aiClient = ai.NewOllamaClient(e.config.OllamaURL, e.config.OllamaModel)

		default: // OpenAI
			validModels := []string{"gpt-3.5-turbo", "gpt-4o", "gpt-4o-mini"}
			isValid := false
			for _, validModel := range validModels {
				if model == validModel {
					isValid = true
					break
				}
			}
			if !isValid {
				return &Result{
					Output:     fmt.Sprintf("Invalid OpenAI model: %s. Use 'config:model list' to see available models.", model),
					IsError:    true,
					CommandRun: cmd.RawInput,
				}, nil
			}

			// Set the model
			e.config.OpenAIModel = model

			// Reinitialize the AI client with the new model
			e.aiClient = ai.NewOpenAIClient(e.config.OpenAIAPIKey, e.config.OpenAIModel)
		}

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     fmt.Sprintf("%s model set to: %s", e.config.AIProvider, model),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	default:
		return &Result{
			Output:     fmt.Sprintf("Unknown model command: %s. Use 'list', 'show', or 'set'.", args[0]),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}

// handleOllamaConfig handles Ollama URL configuration commands
func (e *Executor) handleOllamaConfig(args []string, cmd *nlp.Command) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			Output:     "Missing Ollama command. Use 'show', 'set', or 'test'.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	switch args[0] {
	case "show":
		// Show current Ollama URL
		return &Result{
			Output:     fmt.Sprintf("Current Ollama URL: %s", e.config.OllamaURL),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "set":
		// Set Ollama URL
		if len(args) < 2 {
			return &Result{
				Output:     "Missing URL. Usage: config:ollama set <url>",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		url := args[1]

		// Basic URL validation
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return &Result{
				Output:     "Invalid URL format. URL must start with http:// or https://",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Remove trailing slash if present
		url = strings.TrimSuffix(url, "/")

		// Set the URL
		e.config.OllamaURL = url

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// If Ollama is the current provider, reinitialize the client
		if e.config.AIProvider == "ollama" {
			e.aiClient = ai.NewOllamaClient(e.config.OllamaURL, e.config.OllamaModel)
		}

		return &Result{
			Output:     fmt.Sprintf("Ollama URL set to: %s", url),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "test":
		// Test connection to Ollama server
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(e.config.OllamaURL + "/api/tags")

		if err != nil {
			return &Result{
				Output:     fmt.Sprintf("❌ Cannot connect to Ollama server at %s\nError: %v", e.config.OllamaURL, err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return &Result{
				Output:     fmt.Sprintf("❌ Ollama server returned status code %d", resp.StatusCode),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Try to parse the response to get available models
		var modelsResponse struct {
			Models []struct {
				Name string `json:"name"`
			} `json:"models"`
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return &Result{
				Output:     fmt.Sprintf("✅ Connected to Ollama server at %s, but couldn't read response: %v", e.config.OllamaURL, err),
				IsError:    false,
				CommandRun: cmd.RawInput,
			}, nil
		}

		err = json.Unmarshal(body, &modelsResponse)
		if err != nil {
			return &Result{
				Output:     fmt.Sprintf("✅ Connected to Ollama server at %s, but couldn't parse response: %v", e.config.OllamaURL, err),
				IsError:    false,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Count available models
		modelCount := len(modelsResponse.Models)

		return &Result{
			Output:     fmt.Sprintf("✅ Successfully connected to Ollama server at %s\nFound %d available models", e.config.OllamaURL, modelCount),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	default:
		return &Result{
			Output:     fmt.Sprintf("Unknown Ollama command: %s. Use 'show', 'set', or 'test'.", args[0]),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}

// handleModeConfig handles input mode configuration commands
func (e *Executor) handleModeConfig(args []string, cmd *nlp.Command) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			Output:     "Missing mode command. Use 'show', 'ai', or 'command'.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	switch args[0] {
	case "show":
		// Show current mode
		modeStr := "AI-first"
		if e.config.CommandFirstMode {
			modeStr = "Command-first"
		}

		output := fmt.Sprintf(`
╭─────────────────── 🔧 Input Mode ─────────────────────────╮

  Current input mode: %s

  • AI-first mode: Treats all input as AI queries by default
    unless it starts with a specific command prefix.

  • Command-first mode: Treats input as shell commands if it
    looks like a command, otherwise as an AI query.

╰──────────────────────────────────────────────────────────╯
`, modeStr)

		return &Result{
			Output:     output,
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil

	case "ai":
		// Set AI-first mode
		e.config.CommandFirstMode = false

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     "Input mode set to AI-first. Lumo will now treat all input as AI queries by default unless it starts with a specific command prefix.",
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil

	case "command":
		// Set Command-first mode
		e.config.CommandFirstMode = true

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     "Input mode set to Command-first. Lumo will now treat input as shell commands if it looks like a command, otherwise as an AI query.",
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil

	default:
		return &Result{
			Output:     fmt.Sprintf("Unknown mode command: %s. Use 'show', 'ai', or 'command'.", args[0]),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}

// handleServerConfig handles server configuration commands
func (e *Executor) handleServerConfig(args []string, cmd *nlp.Command) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			Output:     "Missing server command. Use 'show', 'enable', 'disable', or 'quiet'.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	switch args[0] {
	case "show":
		// Show current server settings
		enabledStr := "Disabled"
		if e.config.EnableServer {
			enabledStr = "Enabled"
		}

		quietStr := "Disabled"
		if e.config.ServerQuietOutput {
			quietStr = "Enabled"
		}

		output := fmt.Sprintf(`
╭─────────────────── 🖥️ Server Settings ───────────────────╮

  • Server Status: %s
  • Server Port: %d
  • Quiet Output: %s

  Configure these settings in ~/.config/lumo/config.json
  or use the commands below.

╰──────────────────────────────────────────────────────────╯
`, enabledStr, e.config.ServerPort, quietStr)
		return &Result{
			Output:     output,
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "enable":
		// Enable the server
		e.config.EnableServer = true

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     "Server has been enabled. The REST server will now start when Lumo is executed.\n\nNOTE: The server will be accessible on port " + fmt.Sprintf("%d", e.config.ServerPort) + ". Make sure your firewall is configured appropriately.",
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil

	case "disable":
		// Disable the server
		e.config.EnableServer = false

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     "Server has been disabled. The REST server will not start when Lumo is executed.",
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil

	case "quiet":
		// Set quiet mode
		if len(args) < 2 {
			return &Result{
				Output:     "Missing argument. Use 'on' or 'off'.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		switch args[1] {
		case "on":
			// Enable quiet mode
			e.config.ServerQuietOutput = true
		case "off":
			// Disable quiet mode
			e.config.ServerQuietOutput = false
		default:
			return &Result{
				Output:     fmt.Sprintf("Invalid argument: %s. Use 'on' or 'off'.", args[1]),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		quietStr := "disabled"
		if e.config.ServerQuietOutput {
			quietStr = "enabled"
		}

		return &Result{
			Output:     fmt.Sprintf("Server quiet output %s. Server log messages will %s be displayed.", quietStr, map[bool]string{true: "not", false: ""}[e.config.ServerQuietOutput]),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	default:
		return &Result{
			Output:     fmt.Sprintf("Unknown server command: %s. Use 'show', 'enable', 'disable', or 'quiet'.", args[0]),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}

// handleKeyConfig handles API key configuration commands
func (e *Executor) handleKeyConfig(args []string, cmd *nlp.Command) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			Output:     "Missing key command. Use 'show', 'set', or 'remove'.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	switch args[0] {
	case "show":
		// Use the dedicated key status handler
		return e.handleKeyStatus(cmd)
	case "unused_show":
		// This case is never used, just here to keep the old code structure
		geminiStatus := "Not set"
		if e.config.GeminiAPIKey != "" {
			geminiStatus = "Set"
		}

		openaiStatus := "Not set"
		if e.config.OpenAIAPIKey != "" {
			openaiStatus = "Set"
		}

		output := fmt.Sprintf(`
╭─────────────────── 🔑 API Key Status ─────────────────────╮

  • Gemini API Key: %s
  • OpenAI API Key: %s

  Current provider: %s

╰──────────────────────────────────────────────────────────╯
`, geminiStatus, openaiStatus, e.config.AIProvider)

		return &Result{
			Output:     output,
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "set":
		// Set API key
		if len(args) < 2 {
			return &Result{
				Output:     "Missing provider name. Use 'gemini' or 'openai'. Note: Ollama doesn't require an API key.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		if len(args) < 3 {
			return &Result{
				Output:     "Missing API key. Usage: config:key set <provider> <key>",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		provider := strings.ToLower(args[1])
		apiKey := args[2]

		if provider != "gemini" && provider != "openai" {
			return &Result{
				Output:     fmt.Sprintf("Invalid provider: %s. Use 'gemini' or 'openai'.", provider),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Set the API key
		if provider == "gemini" {
			e.config.GeminiAPIKey = apiKey

			// If this is the current provider, reinitialize the client
			if e.config.AIProvider == "gemini" {
				e.aiClient = ai.NewGeminiClient(e.config.GeminiAPIKey, e.config.GeminiModel)
			}
		} else {
			e.config.OpenAIAPIKey = apiKey

			// If this is the current provider, reinitialize the client
			if e.config.AIProvider == "openai" {
				e.aiClient = ai.NewOpenAIClient(e.config.OpenAIAPIKey, e.config.OpenAIModel)
			}
		}

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     fmt.Sprintf("%s API key has been set.", provider),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	case "remove":
		// Remove API key
		if len(args) < 2 {
			return &Result{
				Output:     "Missing provider name. Use 'gemini' or 'openai'. Note: Ollama doesn't require an API key.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		provider := strings.ToLower(args[1])

		if provider != "gemini" && provider != "openai" {
			return &Result{
				Output:     fmt.Sprintf("Invalid provider: %s. Use 'gemini' or 'openai'.", provider),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Check if this is the current provider
		if e.config.AIProvider == provider {
			return &Result{
				Output:     "Cannot remove API key for the current provider. Switch providers first with 'config:provider set'.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Remove the API key
		if provider == "gemini" {
			e.config.GeminiAPIKey = ""
		} else {
			e.config.OpenAIAPIKey = ""
		}

		// Save the configuration
		if err := e.config.Save(); err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error saving configuration: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     fmt.Sprintf("%s API key has been removed.", provider),
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	default:
		return &Result{
			Output:     fmt.Sprintf("Unknown key command: %s. Use 'show', 'set', or 'remove'.", args[0]),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}
