package executor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agnath18K/lumo/pkg/nlp"
)

// handleKeyStatus handles the key status command
func (e *Executor) handleKeyStatus(cmd *nlp.Command) (*Result, error) {
	// Show API key status (not the actual keys)
	geminiStatus := "Not set"
	if e.config.GeminiAPIKey != "" {
		geminiStatus = "Set"
	}

	openaiStatus := "Not set"
	if e.config.OpenAIAPIKey != "" {
		openaiStatus = "Set"
	}

	// Check Ollama connection
	ollamaStatus := "Not connected"
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	_, err := client.Get(e.config.OllamaURL + "/api/tags")
	if err == nil {
		ollamaStatus = "Connected"
	}

	output := fmt.Sprintf(`
╭─────────────────── 🔑 API Key Status ─────────────────────╮

  • Gemini API Key: %s
  • OpenAI API Key: %s
  • Ollama Server: %s (%s)

  Current provider: %s

╰──────────────────────────────────────────────────────────╯
`, geminiStatus, openaiStatus, ollamaStatus, e.config.OllamaURL, e.config.AIProvider)

	return &Result{
		Output:     output,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}
