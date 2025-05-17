package executor

import (
	"fmt"

	"github.com/agnath18K/lumo/pkg/ai"
	"github.com/agnath18K/lumo/pkg/create"
	"github.com/agnath18K/lumo/pkg/nlp"
)

// executeCreateCommand executes a project creation command
func (e *Executor) executeCreateCommand(cmd *nlp.Command) (*Result, error) {
	// Check if API keys are configured and run setup if needed
	if (e.config.AIProvider == "gemini" && e.config.GeminiAPIKey == "") ||
		(e.config.AIProvider == "openai" && e.config.OpenAIAPIKey == "") {

		// Run interactive setup
		setupPerformed, err := e.apiSetup.CheckAndSetupAPIKeys()
		if err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error during API key setup: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		if setupPerformed {
			// Reinitialize the AI client with the new API key
			if e.config.AIProvider == "gemini" {
				e.aiClient = ai.NewGeminiClient(e.config.GeminiAPIKey, e.config.GeminiModel)
			} else {
				e.aiClient = ai.NewOpenAIClient(e.config.OpenAIAPIKey, e.config.OpenAIModel)
			}
		} else {
			// Setup was not completed successfully
			return &Result{
				Output:     "Error: No API key configured for " + e.config.AIProvider + ". Please set the API key in the configuration or environment variables.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}
	}

	// Create a project generator
	generator := create.NewGenerator(e.aiClient)

	// Execute the create command
	output, err := generator.Execute(cmd.Intent)
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Project Creation Error: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	return &Result{
		Output:     output,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}
