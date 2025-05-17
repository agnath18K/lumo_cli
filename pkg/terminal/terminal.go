package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
)

// Terminal handles terminal interaction
type Terminal struct {
	config         *config.Config
	commandHistory []string
	historyFile    string
}

// NewTerminal creates a new terminal instance
func NewTerminal(cfg *config.Config) *Terminal {
	// Set history file path
	homeDir, err := os.UserHomeDir()
	historyFile := ".lumo_history"
	if err == nil {
		historyFile = homeDir + "/.lumo_history"
	}

	return &Terminal{
		config:         cfg,
		commandHistory: []string{},
		historyFile:    historyFile,
	}
}

// StartInteractiveMode starts an interactive terminal session
func (t *Terminal) StartInteractiveMode(handler func(string)) {
	// Load command history
	t.loadHistory()

	// Create a scanner for reading input
	scanner := bufio.NewScanner(os.Stdin)

	// Display prompt and read input in a loop
	for {
		fmt.Print("lumo> ")
		if !scanner.Scan() {
			break
		}

		// Get input and trim whitespace
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Add command to history
		t.addToHistory(input)

		// Process the command
		handler(input)
	}

	// Save history before exiting
	t.saveHistory()
}

// Display shows the result of a command execution
func (t *Terminal) Display(result *executor.Result) {
	if result.IsError {
		fmt.Fprintf(os.Stderr, "Error: %s\n", result.Output)
	} else {
		fmt.Println(result.Output)
	}
}

// addToHistory adds a command to the history
func (t *Terminal) addToHistory(cmd string) {
	t.commandHistory = append(t.commandHistory, cmd)

	// Trim history if it exceeds the maximum size
	maxHistory := t.config.MaxHistorySize
	if maxHistory <= 0 {
		maxHistory = 1000 // Default to 1000 entries
	}

	if len(t.commandHistory) > maxHistory {
		t.commandHistory = t.commandHistory[len(t.commandHistory)-maxHistory:]
	}
}

// loadHistory loads command history from file
func (t *Terminal) loadHistory() {
	file, err := os.Open(t.historyFile)
	if err != nil {
		return // File doesn't exist or can't be opened
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd != "" {
			t.commandHistory = append(t.commandHistory, cmd)
		}
	}
}

// saveHistory saves command history to file
func (t *Terminal) saveHistory() {
	file, err := os.Create(t.historyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving history: %v\n", err)
		return
	}
	defer file.Close()

	for _, cmd := range t.commandHistory {
		fmt.Fprintln(file, cmd)
	}
}

// GetCommandHistory returns the command history
func (t *Terminal) GetCommandHistory() []string {
	return t.commandHistory
}

// LogCommand logs a command and its result
func (t *Terminal) LogCommand(cmd string, result *executor.Result, duration time.Duration) {
	if !t.config.EnableLogging {
		return
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logs directory: %v\n", err)
		return
	}

	// Open log file
	logFile := fmt.Sprintf("logs/lumo_%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	// Write log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	status := "SUCCESS"
	if result.IsError {
		status = "ERROR"
	}

	logEntry := fmt.Sprintf("[%s] CMD: %s | STATUS: %s | DURATION: %v\n",
		timestamp, cmd, status, duration)

	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
	}
}
