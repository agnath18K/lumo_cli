package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/agnath18K/lumo/pkg/agent"
	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/daemon"
	"github.com/agnath18K/lumo/pkg/executor"
	"github.com/agnath18K/lumo/pkg/nlp"
	"github.com/agnath18K/lumo/pkg/pipe"
	"github.com/agnath18K/lumo/pkg/server"
	"github.com/agnath18K/lumo/pkg/terminal"
	"github.com/agnath18K/lumo/pkg/utils"
	"github.com/agnath18K/lumo/pkg/version"
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

	// Check for server daemon commands
	if len(os.Args) > 1 {
		// Handle server daemon commands
		if os.Args[1] == "server:start" {
			// Start the server daemon
			d := daemon.New(cfg)
			if err := d.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Error starting server daemon: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Server daemon started")
			os.Exit(0)
		} else if os.Args[1] == "server:stop" {
			// Stop the server daemon
			d := daemon.New(cfg)
			if err := d.Stop(); err != nil {
				fmt.Fprintf(os.Stderr, "Error stopping server daemon: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Server daemon stopped")
			os.Exit(0)
		} else if os.Args[1] == "server:status" {
			// Check server daemon status
			d := daemon.New(cfg)
			running, pid, err := d.Status()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error checking server daemon status: %v\n", err)
				os.Exit(1)
			}
			if running {
				fmt.Printf("Server daemon is running with PID %d\n", pid)
			} else {
				fmt.Println("Server daemon is not running")
			}
			os.Exit(0)
		} else if os.Args[1] == "server:daemon" {
			// This is the daemon process
			d := daemon.New(cfg)
			if err := d.RunServer(exec); err != nil {
				fmt.Fprintf(os.Stderr, "Error running server daemon: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	// Check if a server daemon is already running
	d := daemon.New(cfg)
	running, _, err := d.IsRunning()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking if server daemon is running: %v\n", err)
	}

	// Start the REST server if enabled and not already running as a daemon
	var srv *server.Server
	if cfg.EnableServer && !running {
		srv = server.New(cfg, exec)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting REST server: %v\n", err)
			// Continue execution even if server fails to start
		} else {
			// Set up signal handling for graceful shutdown
			setupSignalHandling(srv)

			// Notify the user that the server is running
			if !cfg.ServerQuietOutput {
				fmt.Fprintf(os.Stderr, "\nNOTE: Lumo REST server is running on port %d\n", cfg.ServerPort)
				fmt.Fprintf(os.Stderr, "To disable the server, run: lumo config:server disable\n\n")
			}
		}
	}

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
				"speed:", "speedtest:", "speed-test:", "magic:", "clipboard", "connect", "create", "server:"} {
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
		if strings.HasPrefix(command, "shell:") {
			// Handle shell commands (ONLY with shell: prefix)
			intent := strings.TrimSpace(command[6:])
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
		} else if strings.HasPrefix(command, "server:") {
			// Handle server commands
			intent := strings.TrimSpace(command[7:])
			if intent == "start" {
				// Start the server daemon
				d := daemon.New(cfg)
				if err := d.Start(); err != nil {
					fmt.Fprintf(os.Stderr, "Error starting server daemon: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("Server daemon started")
			} else if intent == "stop" {
				// Stop the server daemon
				d := daemon.New(cfg)
				if err := d.Stop(); err != nil {
					fmt.Fprintf(os.Stderr, "Error stopping server daemon: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("Server daemon stopped")
			} else if intent == "status" {
				// Check server daemon status
				d := daemon.New(cfg)
				running, pid, err := d.Status()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error checking server daemon status: %v\n", err)
					os.Exit(1)
				}
				if running {
					fmt.Printf("Server daemon is running with PID %d\n", pid)
				} else {
					fmt.Println("Server daemon is not running")
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unknown server command: %s\n", intent)
				fmt.Println("Available commands: server:start, server:stop, server:status")
				os.Exit(1)
			}
		} else if strings.HasPrefix(command, "lumo:") {
			// Legacy "lumo:" prefix is now treated as an AI query for safety
			intent := strings.TrimSpace(command[5:])
			cmd := &nlp.Command{
				Type:       nlp.CommandTypeAI,
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

// setupSignalHandling sets up signal handling for graceful shutdown
func setupSignalHandling(srv *server.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if !srv.GetConfig().ServerQuietOutput {
			log.Println("Shutting down REST server...")
		}
		if err := srv.Stop(); err != nil {
			if !srv.GetConfig().ServerQuietOutput {
				log.Printf("Error stopping server: %v", err)
			}
		}
		if !srv.GetConfig().ServerQuietOutput {
			log.Println("Server stopped")
		}
	}()
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
