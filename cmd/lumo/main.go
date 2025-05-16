package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agnath18/lumo/pkg/agent"
	"github.com/agnath18/lumo/pkg/config"
	"github.com/agnath18/lumo/pkg/executor"
	"github.com/agnath18/lumo/pkg/nlp"
	"github.com/agnath18/lumo/pkg/pipe"
	"github.com/agnath18/lumo/pkg/terminal"
	"github.com/agnath18/lumo/pkg/utils"
	"github.com/agnath18/lumo/pkg/version"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize components
	parser := nlp.NewParser(cfg)
	exec := executor.NewExecutor(cfg)
	term := terminal.NewTerminal(cfg)

	// Initialize agent
	_ = agent.Initialize(cfg, exec)

	// Check if input is being piped
	stat, _ := os.Stdin.Stat()
	isPiped := (stat.Mode() & os.ModeCharDevice) == 0

	if isPiped && cfg.EnablePipeProcessing {
		// Process piped input
		processPipedInput(exec, term)
	} else if len(os.Args) > 1 {
		// Check for version flag
		if os.Args[1] == "--version" || os.Args[1] == "-v" || os.Args[1] == "version" {
			version.PrintVersion()
			os.Exit(0)
		}

		// Check for help flag
		if os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help" {
			// Display help message
			helpCmd := &nlp.Command{
				Type:       nlp.CommandTypeHelp,
				Intent:     "help",
				Parameters: make(map[string]string),
				RawInput:   "help",
			}
			result, err := exec.Execute(helpCmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying help: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
			os.Exit(0)
		}

		// Process command from arguments
		// Join arguments with spaces, preserving quotes if present
		command := strings.Join(os.Args[1:], " ")

		// In AI-first mode (default), we don't need special handling for quoted strings
		// as everything will be treated as an AI query by default unless it has a specific prefix
		// or is a single executable command in command-first mode.

		// However, we still want to handle the case where a command might be a quoted string
		// that was split by the shell, for better user experience
		if len(os.Args) > 2 && !cfg.CommandFirstMode {
			// If we have multiple arguments and none of them start with a prefix like "lumo:" or "shell:",
			// it might be a quoted string that was split
			hasPrefix := false
			for _, prefix := range []string{"lumo:", "shell:", "ask:", "ai:", "auto:", "agent:",
				"health:", "syshealth:", "report:", "sysreport:", "chat:", "talk:", "config:",
				"speed:", "speedtest:", "speed-test:", "magic:", "clipboard", "connect"} {
				if strings.HasPrefix(command, prefix) {
					hasPrefix = true
					break
				}
			}

			if !hasPrefix {
				// In AI-first mode, treat it as an AI query by default
				cmd := &nlp.Command{
					Type:       nlp.CommandTypeAI,
					Intent:     command,
					Parameters: make(map[string]string),
					RawInput:   command,
				}
				result, err := exec.Execute(cmd)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
					os.Exit(1)
				}
				term.Display(result)
				os.Exit(0)
			}
		}

		// Special handling for commands with specific prefixes
		if strings.HasPrefix(command, "lumo:") || strings.HasPrefix(command, "shell:") {
			// Handle shell commands
			var intent string
			if strings.HasPrefix(command, "lumo:") {
				intent = strings.TrimSpace(command[5:])
			} else {
				intent = strings.TrimSpace(command[6:])
			}
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeShell,
				Intent:     intent,
				Parameters: make(map[string]string),
				RawInput:   command,
			}
			result, err := exec.Execute(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
		} else if strings.HasPrefix(command, "health:") || strings.HasPrefix(command, "syshealth:") {
			// Handle system health commands
			var intent string
			if strings.HasPrefix(command, "health:") {
				intent = strings.TrimSpace(command[7:])
			} else {
				intent = strings.TrimSpace(command[10:])
			}
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeSystemHealth,
				Intent:     intent,
				Parameters: make(map[string]string),
				RawInput:   command,
			}
			result, err := exec.Execute(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
		} else if strings.HasPrefix(command, "report:") || strings.HasPrefix(command, "sysreport:") {
			// Handle system report commands
			var intent string
			if strings.HasPrefix(command, "report:") {
				intent = strings.TrimSpace(command[7:])
			} else {
				intent = strings.TrimSpace(command[10:])
			}
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeSystemReport,
				Intent:     intent,
				Parameters: make(map[string]string),
				RawInput:   command,
			}
			result, err := exec.Execute(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
		} else if strings.HasPrefix(command, "chat:") || strings.HasPrefix(command, "talk:") || command == "chat" || command == "talk" {
			// Handle chat commands
			var intent string
			var rawInput string

			if command == "chat" || command == "talk" {
				// Empty chat command to start REPL mode
				intent = ""
				rawInput = command + ":"
			} else if strings.HasPrefix(command, "chat:") {
				intent = strings.TrimSpace(command[5:])
				rawInput = command
			} else {
				intent = strings.TrimSpace(command[5:])
				rawInput = command
			}

			cmd := &nlp.Command{
				Type:       nlp.CommandTypeChat,
				Intent:     intent,
				Parameters: make(map[string]string),
				RawInput:   rawInput,
			}
			result, err := exec.Execute(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
		} else if strings.HasPrefix(command, "config:") {
			// Handle configuration commands
			intent := strings.TrimSpace(command[7:])
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeConfig,
				Intent:     intent,
				Parameters: make(map[string]string),
				RawInput:   command,
			}
			result, err := exec.Execute(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
		} else if strings.HasPrefix(command, "magic:") {
			// Handle magic commands
			intent := strings.TrimSpace(command[6:])
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeMagic,
				Intent:     intent,
				Parameters: make(map[string]string),
				RawInput:   command,
			}
			result, err := exec.Execute(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
		} else if command == "clipboard" || strings.HasPrefix(command, "clipboard ") {
			// Handle clipboard commands
			intent := ""
			if strings.HasPrefix(command, "clipboard ") {
				intent = strings.TrimSpace(command[10:])
			}
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeClipboard,
				Intent:     intent,
				Parameters: make(map[string]string),
				RawInput:   command,
			}
			result, err := exec.Execute(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
			term.Display(result)
		} else {
			processCommand(command, parser, exec, term)
		}
	} else {
		// Display welcome message when run without arguments
		result, err := exec.ShowWelcome()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error displaying welcome message: %v\n", err)
			os.Exit(1)
		}
		term.Display(result)
	}
}

func processPipedInput(exec *executor.Executor, term *terminal.Terminal) {
	// Record start time for performance measurement
	startTime := time.Now()

	// Check if we have arguments that might indicate a clipboard command
	if len(os.Args) > 1 && (os.Args[1] == "clipboard" || strings.HasPrefix(os.Args[1], "clipboard ")) {
		// This is a clipboard command with piped input
		intent := ""

		// Check if it's an append operation
		if len(os.Args) > 2 && os.Args[2] == "append" {
			intent = "append "
		}

		cmd := &nlp.Command{
			Type:       nlp.CommandTypeClipboard,
			Intent:     intent, // Empty intent means use piped input, "append " means append piped input
			Parameters: make(map[string]string),
			RawInput:   os.Args[1],
		}

		// Execute with stdin as the reader
		result, err := exec.ExecuteWithReader(cmd, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing clipboard command: %v\n", err)
			os.Exit(1)
		}

		// Display the result
		term.Display(result)

		// Calculate execution duration
		duration := time.Since(startTime)

		// Log the command if logging is enabled
		term.LogCommand(os.Args[1], result, duration)

		// Show execution time in debug mode
		if exec.GetConfig().Debug {
			fmt.Printf("Execution time: %s\n", utils.FormatDuration(duration))
		}
		return
	}

	// For non-clipboard commands, process as before
	// Create a pipe processor
	pipeProcessor := pipe.NewProcessor(exec.GetAIClient())

	// Process the piped input
	result, err := pipeProcessor.ProcessInput(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing piped input: %v\n", err)
		os.Exit(1)
	}

	// Create a result object
	execResult := &executor.Result{
		Output:     result,
		IsError:    false,
		CommandRun: "piped input",
	}

	// Calculate execution duration
	duration := time.Since(startTime)

	// Display the result
	term.Display(execResult)

	// Log the command if logging is enabled
	term.LogCommand("piped input", execResult, duration)

	// Show execution time in debug mode
	if exec.GetConfig().Debug {
		fmt.Printf("Execution time: %s\n", utils.FormatDuration(duration))
	}
}

func processCommand(input string, parser *nlp.Parser, exec *executor.Executor, term *terminal.Terminal) {
	// Check for exit commands
	if input == "exit" || input == "quit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	// Check for version command
	if input == "version" || input == "--version" || input == "-v" {
		fmt.Println("Lumo version", version.GetShortVersion())
		os.Exit(0)
	}

	// Record start time for performance measurement
	startTime := time.Now()

	// Parse the natural language input
	cmd, err := parser.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing command: %v\n", err)
		return
	}

	// Execute the command
	result, err := exec.Execute(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		return
	}

	// Calculate execution duration
	duration := time.Since(startTime)

	// Display the result
	term.Display(result)

	// Log the command if logging is enabled
	term.LogCommand(input, result, duration)

	// Show execution time in debug mode
	if exec.GetConfig().Debug {
		fmt.Printf("Execution time: %s\n", utils.FormatDuration(duration))
	}
}
