package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/agnath18/lumo/pkg/ai"
	"github.com/agnath18/lumo/pkg/chat"
	"github.com/agnath18/lumo/pkg/clipboard"
	"github.com/agnath18/lumo/pkg/config"
	"github.com/agnath18/lumo/pkg/magic"
	"github.com/agnath18/lumo/pkg/nlp"
	"github.com/agnath18/lumo/pkg/setup"
	"github.com/agnath18/lumo/pkg/system"
	"github.com/agnath18/lumo/pkg/utils"
)

// Result represents the output of a command execution
type Result struct {
	Output     string
	IsError    bool
	CommandRun string
}

// Executor handles command execution
type Executor struct {
	config      *config.Config
	aiClient    ai.Client
	apiSetup    *setup.APIKeySetup
	agent       AgentInterface
	chatManager *chat.Manager
	magic       *magic.Magic
	clipboard   *clipboard.Clipboard
}

// NewExecutor creates a new executor instance
func NewExecutor(cfg *config.Config) *Executor {
	// Create AI client based on configuration
	var aiClient ai.Client
	switch cfg.AIProvider {
	case "gemini":
		aiClient = ai.NewGeminiClient(cfg.GeminiAPIKey, cfg.GeminiModel)
	case "ollama":
		aiClient = ai.NewOllamaClient(cfg.OllamaURL, cfg.OllamaModel)
	default: // Default to OpenAI
		aiClient = ai.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.OpenAIModel)
	}

	// Create a chat manager
	chatManager := chat.NewManager(aiClient, 5, 20)

	return &Executor{
		config:      cfg,
		aiClient:    aiClient,
		apiSetup:    setup.NewAPIKeySetup(cfg),
		chatManager: chatManager,
		// The agent will be set later by the agent package
		agent: nil,
		// Initialize the magic handler
		magic: magic.NewMagic(),
		// Initialize the clipboard handler
		clipboard: clipboard.NewClipboard(),
	}
}

// SetAgent sets the agent implementation
func (e *Executor) SetAgent(agent AgentInterface) {
	e.agent = agent
}

// Execute processes a command and returns the result
func (e *Executor) Execute(cmd *nlp.Command) (*Result, error) {
	return e.ExecuteWithReader(cmd, nil)
}

// ExecuteWithReader executes a command with an optional reader for piped input
func (e *Executor) ExecuteWithReader(cmd *nlp.Command, reader io.Reader) (*Result, error) {
	switch cmd.Type {
	case nlp.CommandTypeShell:
		return e.executeShellCommand(cmd)
	case nlp.CommandTypeAI:
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
		return e.executeAIQuery(cmd)
	case nlp.CommandTypeChat:
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
		return e.executeChatCommand(cmd)
	case nlp.CommandTypeAgent:
		// Check if agent is initialized
		if e.agent == nil {
			return &Result{
				Output:     "Agent mode is not available. Please initialize the agent first.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Check if agent mode is enabled
		if !e.config.EnableAgentMode {
			return &Result{
				Output:     "Agent mode is disabled. Enable it in the configuration file.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

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
		return e.executeAgentCommand(cmd)
	case nlp.CommandTypeSystemHealth:
		// Check if system health is enabled
		if !e.config.EnableSystemHealth {
			return &Result{
				Output:     "System health checks are disabled. Enable them in the configuration file.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}
		return e.executeSystemHealthCheck(cmd)
	case nlp.CommandTypeSystemReport:
		// Check if system report is enabled
		if !e.config.EnableSystemReport {
			return &Result{
				Output:     "System reports are disabled. Enable them in the configuration file.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}
		return e.executeSystemReport(cmd)
	case nlp.CommandTypeHelp:
		return e.showHelp(cmd)
	case nlp.CommandTypeConfig:
		return e.executeConfigCommand(cmd)
	case nlp.CommandTypeSpeedTest:
		// Check if speed test is enabled
		if !e.config.EnableSpeedTest {
			return &Result{
				Output:     "Speed test is disabled. Enable it in the configuration file.",
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}
		return e.executeSpeedTest(cmd)
	case nlp.CommandTypeMagic:
		// Execute magic command
		return e.executeMagicCommand(cmd)
	case nlp.CommandTypeClipboard:
		// Execute clipboard command
		return e.executeClipboardCommand(cmd, reader)
	case nlp.CommandTypeConnect:
		// Execute connect command
		return e.executeConnectCommand(cmd)
	default:
		return &Result{
			Output:     "Unknown command type",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}
}

// executeShellCommand runs a shell command
func (e *Executor) executeShellCommand(cmd *nlp.Command) (*Result, error) {
	// Split the command into parts
	parts := strings.Fields(cmd.Intent)
	if len(parts) == 0 {
		return &Result{
			Output:     "Empty command",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Check if the command exists before trying to execute it
	_, err := exec.LookPath(parts[0])
	if err != nil {
		// Command doesn't exist, provide a helpful error message
		suggestion := ""
		if len(parts) > 1 {
			// If there are multiple words, suggest using it as an AI query
			suggestion = fmt.Sprintf("\n\nDid you mean to ask AI about \"%s\"? Try: lumo ask:\"%s\"", cmd.Intent, cmd.Intent)
		}

		return &Result{
			Output:     fmt.Sprintf("Error: exec: \"%s\": executable file not found in $PATH%s", parts[0], suggestion),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Create the command
	shellCmd := exec.Command(parts[0], parts[1:]...)

	// Run the command and capture output
	output, err := shellCmd.CombinedOutput()

	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Error: %v\n%s", err, string(output)),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	return &Result{
		Output:     string(output),
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}

// executeAIQuery sends a query to the AI service
func (e *Executor) executeAIQuery(cmd *nlp.Command) (*Result, error) {
	// Check internet connectivity for cloud-based providers
	if (e.config.AIProvider == "gemini" || e.config.AIProvider == "openai") && !utils.CheckInternetConnectivity() {
		// We're offline and using a cloud provider

		// Check if Ollama is available locally
		ollamaAvailable := e.isOllamaAvailable()

		// Use the new function for a more humorous offline warning without a box
		return &Result{
			Output:     utils.FormatOfflineWarning(e.config.AIProvider, ollamaAvailable, false),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Proceed with the query
	response, err := e.aiClient.Query(cmd.Intent)
	if err != nil {
		// Check if the error might be due to connectivity issues
		if !utils.CheckInternetConnectivity() && (e.config.AIProvider == "gemini" || e.config.AIProvider == "openai") {
			// We're offline and using a cloud provider
			ollamaAvailable := e.isOllamaAvailable()

			// Use the new function for a more humorous offline warning without a box
			return &Result{
				Output:     "Error: " + err.Error() + "\n\n" + utils.FormatOfflineWarning(e.config.AIProvider, ollamaAvailable, false),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Regular error handling
		return &Result{
			Output:     fmt.Sprintf("AI Error: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Clean up markdown formatting for better terminal display
	cleanResponse := utils.CleanMarkdown(response)

	// Check if the response already has a box format (either style)
	hasBox := (strings.Contains(cleanResponse, "┌") && strings.Contains(cleanResponse, "┐") &&
		strings.Contains(cleanResponse, "└") && strings.Contains(cleanResponse, "┘")) ||
		(strings.Contains(cleanResponse, "╭") && strings.Contains(cleanResponse, "╮") &&
			strings.Contains(cleanResponse, "╰") && strings.Contains(cleanResponse, "╯"))

	// If the response doesn't already have a box, add one
	if !hasBox {
		// Add a box around the response for consistent display
		title := "🐦 Lumo"
		cleanResponse = utils.FormatWithBox(cleanResponse, title)
	}

	return &Result{
		Output:     cleanResponse,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}

// executeChatCommand processes a chat message and returns the AI response
func (e *Executor) executeChatCommand(cmd *nlp.Command) (*Result, error) {
	// Check if chat REPL mode is enabled
	if e.config.EnableChatREPL && cmd.Intent == "" {
		// Start REPL mode
		return e.startChatREPL()
	}

	// Check internet connectivity for cloud-based providers
	if (e.config.AIProvider == "gemini" || e.config.AIProvider == "openai") && !utils.CheckInternetConnectivity() {
		// We're offline and using a cloud provider

		// Check if Ollama is available locally
		ollamaAvailable := e.isOllamaAvailable()

		// Use the new function for a more humorous offline warning without a box
		return &Result{
			Output:     utils.FormatOfflineWarning(e.config.AIProvider, ollamaAvailable, false),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Create a context
	ctx := context.Background()

	// Process the message using the chat manager
	response, err := e.chatManager.ProcessMessage(ctx, cmd.Intent)
	if err != nil {
		// Check if the error might be due to connectivity issues
		if !utils.CheckInternetConnectivity() && (e.config.AIProvider == "gemini" || e.config.AIProvider == "openai") {
			// We're offline and using a cloud provider
			ollamaAvailable := e.isOllamaAvailable()

			// Use the new function for a more humorous offline warning without a box
			return &Result{
				Output:     "Error: " + err.Error() + "\n\n" + utils.FormatOfflineWarning(e.config.AIProvider, ollamaAvailable, false),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Regular error handling
		return &Result{
			Output:     fmt.Sprintf("Chat Error: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Clean up markdown formatting for better terminal display
	cleanResponse := utils.CleanMarkdown(response)

	// Check if the response already has a box format (either style)
	hasBox := (strings.Contains(cleanResponse, "┌") && strings.Contains(cleanResponse, "┐") &&
		strings.Contains(cleanResponse, "└") && strings.Contains(cleanResponse, "┘")) ||
		(strings.Contains(cleanResponse, "╭") && strings.Contains(cleanResponse, "╮") &&
			strings.Contains(cleanResponse, "╰") && strings.Contains(cleanResponse, "╯"))

	// If the response doesn't already have a box, add one
	// This is only for single chat commands, not for REPL mode
	if !hasBox {
		// Add a box around the response for consistent display
		title := "🐦 Lumo Chat"
		cleanResponse = utils.FormatWithBox(cleanResponse, title)
	}

	return &Result{
		Output:     cleanResponse,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}

// startChatREPL starts the chat REPL mode
func (e *Executor) startChatREPL() (*Result, error) {
	// Create a new REPL instance
	repl := chat.NewREPL(e.config, e.chatManager, e.aiClient)

	// Start the REPL loop
	output, err := repl.Start()
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Chat REPL Error: %v", err),
			IsError:    true,
			CommandRun: "chat:",
		}, nil
	}

	return &Result{
		Output:     output,
		IsError:    false,
		CommandRun: "chat:",
	}, nil
}

// executeAgentCommand executes a command using the agent
func (e *Executor) executeAgentCommand(cmd *nlp.Command) (*Result, error) {
	// Check internet connectivity for cloud-based providers
	if (e.config.AIProvider == "gemini" || e.config.AIProvider == "openai") && !utils.CheckInternetConnectivity() {
		// We're offline and using a cloud provider

		// Check if Ollama is available locally
		ollamaAvailable := e.isOllamaAvailable()

		// Use the new function for a more humorous offline warning without a box
		return &Result{
			Output:     utils.FormatOfflineWarning(e.config.AIProvider, ollamaAvailable, true),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Create a context
	ctx := context.Background()

	// Execute the command using the agent
	result, err := e.agent.Execute(ctx, cmd.Intent)

	// Check if the error might be due to connectivity issues
	if err != nil && !utils.CheckInternetConnectivity() && (e.config.AIProvider == "gemini" || e.config.AIProvider == "openai") {
		// We're offline and using a cloud provider
		ollamaAvailable := e.isOllamaAvailable()

		// Use the new function for a more humorous offline warning without a box
		return &Result{
			Output:     "Error: " + err.Error() + "\n\n" + utils.FormatOfflineWarning(e.config.AIProvider, ollamaAvailable, true),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	return result, err
}

// executeSystemHealthCheck performs a system health check
func (e *Executor) executeSystemHealthCheck(cmd *nlp.Command) (*Result, error) {
	// Create a health checker
	healthChecker := system.NewHealthChecker()

	// Perform health check
	healthResult, err := healthChecker.CheckHealth()
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Error performing health check: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Format the health check result
	formattedResult := system.FormatHealthCheck(healthResult)

	return &Result{
		Output:     formattedResult,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}

// executeSystemReport generates a system report
func (e *Executor) executeSystemReport(cmd *nlp.Command) (*Result, error) {
	// Create a report generator
	reportGenerator := system.NewReportGenerator()

	// Generate report
	report, err := reportGenerator.GenerateReport()
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Error generating system report: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Format the report
	formattedReport := system.FormatSystemReport(report)

	return &Result{
		Output:     formattedReport,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}

// showHelp displays help information
func (e *Executor) showHelp(cmd *nlp.Command) (*Result, error) {
	shellStatus := "DISABLED"
	if e.config.EnableShellInInteractive {
		shellStatus = "ENABLED"
	}

	agentStatus := "DISABLED"
	if e.config.EnableAgentMode {
		agentStatus = "ENABLED"
	}

	replStatus := "DISABLED"
	if e.config.EnableAgentREPL {
		replStatus = "ENABLED"
	}

	pipeStatus := "DISABLED"
	if e.config.EnablePipeProcessing {
		pipeStatus = "ENABLED"
	}

	healthStatus := "DISABLED"
	if e.config.EnableSystemHealth {
		healthStatus = "ENABLED"
	}

	reportStatus := "DISABLED"
	if e.config.EnableSystemReport {
		reportStatus = "ENABLED"
	}

	// Get chat REPL status
	chatReplStatus := "Disabled"
	if e.config.EnableChatREPL {
		chatReplStatus = "Enabled"
	}

	// Get speed test status
	speedTestStatus := "DISABLED"
	if e.config.EnableSpeedTest {
		speedTestStatus = "ENABLED"
	}

	helpText := fmt.Sprintf(`
╭──────────────────── 🐦 Lumo CLI Assistant ──────────────────────╮

  Commands:
   • ask:<query>                Ask the AI a question
   • chat:<message>             Start or continue a conversation
   • chat                       Start interactive chat mode
   • lumo:<command>             Run shell command [%s]
   • shell:<command>            Run shell command [%s]
   • auto:<task>                Use agent mode [%s]
   • agent:<task>               Use agent mode [%s]
   • health:<options>           Check system health [%s]
   • syshealth:<options>        Check system health [%s]
   • report:<options>           Generate system report [%s]
   • sysreport:<options>        Generate system report [%s]
   • speed:<options>            Run internet speed test [%s]
   • magic:<command>            Run fun magic commands
   • clipboard                  Show clipboard contents
   • clipboard <text>           Copy text to clipboard
   • clipboard append <text>    Append text to clipboard
   • clipboard clear            Clear clipboard contents
   • connect --receive [options]  Start a server to send/receive files
   • connect <peer-ip> [options]  Connect to peer to send/receive files
   • connect --help              Show connect command options
   • config:<options>           Configure Lumo settings
   • version, -v, --version     Show version information
   • help, -h, --help           Show this help

  Examples:
   • lumo "how to find large files"
   • chat:Tell me about Linux
   • chat                       Start interactive chat session
   • lumo:ls -la
   • auto:"create a backup of my documents"
   • magic:dance                Show a fun dance animation
   • clipboard                  Show current clipboard contents
   • clipboard "Hello World"    Copy text to clipboard
   • clipboard append "More"    Append text to clipboard
   • clipboard clear            Clear clipboard contents
   • echo "text" | clipboard    Copy piped text to clipboard
   • echo "more" | clipboard append  Append piped text to clipboard
   • connect --receive          Start a server on port 8080
   • connect --receive --port 9000  Start a server on port 9000
   • connect 192.168.1.5        Connect to peer at 192.168.1.5:8080
   • speed:                     Run a full internet speed test
   • speed:download             Test download speed only
   • cat file.txt | lumo        Analyze piped content
   • config:model list          List available AI models
   • config:key show            Show API key status
   • version                    Show version information

  Configuration:
   • config:provider list       List available AI providers
   • config:provider show       Show current AI provider
   • config:provider set <name> Set AI provider (gemini/openai/ollama)
   • config:model list          List available models
   • config:model set <name>    Set model for current provider
   • config:key set <prov> <key> Set API key for provider
   • config:ollama show         Show current Ollama URL
   • config:ollama set <url>    Set Ollama URL
   • config:ollama test         Test connection to Ollama server

  Status:
   • Shell in interactive mode: %s
   • Agent mode: %s
   • Agent REPL mode: %s
   • Chat REPL mode: %s
   • Pipe processing: %s
   • System health checks: %s
   • System reports: %s
   • Speed test: %s
   • Current AI provider: %s
   • Current model: %s

  API Keys:
   • Gemini: https://aistudio.google.com/apikey
   • OpenAI: https://platform.openai.com/api-keys
   • Ollama: http://localhost:11434 (default local URL)

  ⚠️  DISCLAIMERS:
   • For basic terminal help only, not coding tasks
   • Agent mode executes commands - ALWAYS review plans!
   • Use 'ask:' instead of 'auto:' for safer operation
   • Offline mode available with Ollama (config:provider set ollama)

╰─────────────────────────────────────────────────────────────────────╯
`, shellStatus, shellStatus, agentStatus, agentStatus, healthStatus, healthStatus, reportStatus, reportStatus, speedTestStatus, shellStatus, agentStatus, replStatus, chatReplStatus, pipeStatus, healthStatus, reportStatus, speedTestStatus, e.config.AIProvider, getCurrentModel(e.config))

	return &Result{
		Output:     helpText,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}

// GetConfig returns the executor's configuration
func (e *Executor) GetConfig() *config.Config {
	return e.config
}

// GetAIClient returns the executor's AI client
func (e *Executor) GetAIClient() ai.Client {
	return e.aiClient
}

// ShowWelcome displays a minimal welcome message
func (e *Executor) ShowWelcome() (*Result, error) {
	welcomeText := `
╭─────────────────── 🐦 Lumo CLI Assistant ─────────────────╮

  Welcome to Lumo! Type your query or use a command prefix.

  Examples:
   • lumo "how to find large files"
   • lumo chat:Tell me about Linux
   • lumo auto:"create a backup of my documents"
   • lumo connect --receive

  Type 'help' for full documentation and available commands.

╰───────────────────────────────────────────────────────────╯
`
	return &Result{
		Output:     welcomeText,
		IsError:    false,
		CommandRun: "welcome",
	}, nil
}

// executeMagicCommand executes a magic command
func (e *Executor) executeMagicCommand(cmd *nlp.Command) (*Result, error) {
	// Execute the magic command
	output, err := e.magic.Execute(cmd.Intent)
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Magic Error: %v", err),
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

// executeClipboardCommand executes a clipboard command
func (e *Executor) executeClipboardCommand(cmd *nlp.Command, reader io.Reader) (*Result, error) {
	// Execute the clipboard command
	output, err := e.clipboard.Execute(cmd.Intent, reader)
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Clipboard Error: %v", err),
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

// isOllamaAvailable checks if Ollama is available locally
func (e *Executor) isOllamaAvailable() bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	_, err := client.Get(e.config.OllamaURL + "/api/tags")
	return err == nil
}
