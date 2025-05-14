package nlp

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/agnath18/lumo/pkg/config"
)

// Command represents a parsed command with its type and parameters
type Command struct {
	Type       CommandType
	Intent     string
	Parameters map[string]string
	RawInput   string
}

// CommandType represents the type of command to execute
type CommandType int

const (
	// CommandTypeUnknown represents an unknown command
	CommandTypeUnknown CommandType = iota
	// CommandTypeShell represents a shell command
	CommandTypeShell
	// CommandTypeAI represents an AI query
	CommandTypeAI
	// CommandTypeHelp represents a help request
	CommandTypeHelp
	// CommandTypeSystem represents a system command
	CommandTypeSystem
	// CommandTypeAgent represents an agent command
	CommandTypeAgent
	// CommandTypeSystemHealth represents a system health check command
	CommandTypeSystemHealth
	// CommandTypeSystemReport represents a system report command
	CommandTypeSystemReport
	// CommandTypeChat represents a chat conversation command
	CommandTypeChat
	// CommandTypeConfig represents a configuration command
	CommandTypeConfig
	// CommandTypeSpeedTest represents an internet speed test command
	CommandTypeSpeedTest
	// CommandTypeMagic represents a magic command
	CommandTypeMagic
	// CommandTypeClipboard represents a clipboard operation
	CommandTypeClipboard
	// CommandTypeConnect represents a file transfer connection
	CommandTypeConnect
)

// Parser handles natural language parsing
type Parser struct {
	config *config.Config
}

// NewParser creates a new parser instance
func NewParser(cfg *config.Config) *Parser {
	return &Parser{
		config: cfg,
	}
}

// Parse processes natural language input and returns a structured command
func (p *Parser) Parse(input string) (*Command, error) {
	// Trim input
	input = strings.TrimSpace(input)

	// Create a new command with the raw input
	cmd := &Command{
		RawInput:   input,
		Parameters: make(map[string]string),
	}

	// Check for help command
	if input == "help" {
		cmd.Type = CommandTypeHelp
		cmd.Intent = "help"
		return cmd, nil
	}

	// Check for shell command prefix
	if strings.HasPrefix(input, "lumo:") || strings.HasPrefix(input, "shell:") {
		// Check if we're in interactive mode and shell commands are disabled
		args := os.Args
		isInteractiveMode := len(args) <= 1 || input != strings.Join(args[1:], " ")

		if isInteractiveMode && !p.config.EnableShellInInteractive {
			// Shell commands are disabled in interactive mode
			cmd.Type = CommandTypeAI
			cmd.Intent = input
			return cmd, nil
		}

		// Process as shell command
		cmd.Type = CommandTypeShell
		if strings.HasPrefix(input, "lumo:") {
			cmd.Intent = strings.TrimSpace(input[5:])
		} else {
			cmd.Intent = strings.TrimSpace(input[6:])
		}
		return cmd, nil
	}

	// Check for AI query prefix
	if strings.HasPrefix(input, "ask:") || strings.HasPrefix(input, "ai:") {
		cmd.Type = CommandTypeAI
		if strings.HasPrefix(input, "ask:") {
			cmd.Intent = strings.TrimSpace(input[4:])
		} else {
			cmd.Intent = strings.TrimSpace(input[3:])
		}
		return cmd, nil
	}

	// Check for agent command prefix
	if strings.HasPrefix(input, "auto:") || strings.HasPrefix(input, "agent:") {
		cmd.Type = CommandTypeAgent
		if strings.HasPrefix(input, "auto:") {
			cmd.Intent = strings.TrimSpace(input[5:])
		} else {
			cmd.Intent = strings.TrimSpace(input[6:])
		}
		return cmd, nil
	}

	// Check for system health command prefix
	if strings.HasPrefix(input, "health:") || strings.HasPrefix(input, "syshealth:") {
		cmd.Type = CommandTypeSystemHealth
		if strings.HasPrefix(input, "health:") {
			cmd.Intent = strings.TrimSpace(input[7:])
		} else {
			cmd.Intent = strings.TrimSpace(input[10:])
		}
		return cmd, nil
	}

	// Check for system report command prefix
	if strings.HasPrefix(input, "report:") || strings.HasPrefix(input, "sysreport:") {
		cmd.Type = CommandTypeSystemReport
		if strings.HasPrefix(input, "report:") {
			cmd.Intent = strings.TrimSpace(input[7:])
		} else {
			cmd.Intent = strings.TrimSpace(input[10:])
		}
		return cmd, nil
	}

	// Check for chat command prefix
	if strings.HasPrefix(input, "chat:") || strings.HasPrefix(input, "talk:") {
		cmd.Type = CommandTypeChat
		if strings.HasPrefix(input, "chat:") {
			cmd.Intent = strings.TrimSpace(input[5:])
		} else {
			cmd.Intent = strings.TrimSpace(input[5:])
		}
		return cmd, nil
	}

	// Check for config command prefix
	if strings.HasPrefix(input, "config:") {
		cmd.Type = CommandTypeConfig
		cmd.Intent = strings.TrimSpace(input[7:])
		return cmd, nil
	}

	// Check for speed test command prefix
	if strings.HasPrefix(input, "speed:") || strings.HasPrefix(input, "speedtest:") || strings.HasPrefix(input, "speed-test:") {
		cmd.Type = CommandTypeSpeedTest
		if strings.HasPrefix(input, "speed:") {
			cmd.Intent = strings.TrimSpace(input[6:])
		} else if strings.HasPrefix(input, "speedtest:") {
			cmd.Intent = strings.TrimSpace(input[10:])
		} else {
			cmd.Intent = strings.TrimSpace(input[11:])
		}
		return cmd, nil
	}

	// Check for magic command prefix
	if strings.HasPrefix(input, "magic:") {
		cmd.Type = CommandTypeMagic
		cmd.Intent = strings.TrimSpace(input[6:])
		return cmd, nil
	}

	// Check for clipboard command
	if input == "clipboard" || strings.HasPrefix(input, "clipboard ") {
		cmd.Type = CommandTypeClipboard
		cmd.Intent = strings.TrimSpace(strings.TrimPrefix(input, "clipboard"))
		return cmd, nil
	}

	// Check for connect command
	if input == "connect" || strings.HasPrefix(input, "connect ") {
		cmd.Type = CommandTypeConnect
		cmd.Intent = strings.TrimSpace(strings.TrimPrefix(input, "connect"))
		return cmd, nil
	}

	// Check if this is a command-line argument (first argument is the program name)
	args := os.Args
	if len(args) > 1 && input == strings.Join(args[1:], " ") {
		// Check if the input looks like a natural language query
		// Natural language queries typically:
		// 1. Start with a capital letter (questions, sentences)
		// 2. Contain multiple words with spaces
		// 3. End with a question mark (for questions)
		// 4. Don't contain shell command special characters like |, >, <, etc.

		// If it looks like a natural language query, treat it as an AI query
		if isNaturalLanguageQuery(input) {
			cmd.Type = CommandTypeAI
			cmd.Intent = input
			return cmd, nil
		}

		// Otherwise, treat it as a shell command
		cmd.Type = CommandTypeShell
		cmd.Intent = input
		return cmd, nil
	}

	// Check if this looks like a speed test query
	if isSpeedTestQuery(input) {
		cmd.Type = CommandTypeSpeedTest
		cmd.Intent = input
		return cmd, nil
	}

	// Check if this looks like a task that should use agent mode
	if isAgentTask(input) && p.config.EnableAgentMode {
		cmd.Type = CommandTypeAgent
		cmd.Intent = input
	} else {
		// Default to AI query for natural language processing
		cmd.Type = CommandTypeAI
		cmd.Intent = input
	}

	return cmd, nil
}

// isNaturalLanguageQuery determines if a string is likely to be a natural language query
// rather than a shell command
func isNaturalLanguageQuery(input string) bool {
	// Trim the input
	input = strings.TrimSpace(input)

	// If the input is empty, it's not a natural language query
	if input == "" {
		return false
	}

	// Check if the input is quoted (starts and ends with quotes)
	// This is a strong indicator that it's a natural language query
	if (strings.HasPrefix(input, "\"") && strings.HasSuffix(input, "\"")) ||
		(strings.HasPrefix(input, "'") && strings.HasSuffix(input, "'")) {
		return true
	}

	// Check if the input starts with a capital letter (common for questions and sentences)
	if len(input) > 0 && input[0] >= 'A' && input[0] <= 'Z' {
		return true
	}

	// Check if the input ends with a question mark (common for questions)
	if strings.HasSuffix(input, "?") {
		return true
	}

	// Check if the input contains shell command special characters
	shellSpecialChars := []string{"|", ">", "<", ";", "&", "$", "(", ")", "[", "]", "{", "}", "`"}
	for _, char := range shellSpecialChars {
		if strings.Contains(input, char) {
			return false
		}
	}

	// Check for common greetings and conversational phrases
	commonPhrases := []string{"hello", "hi", "hey", "greetings", "thanks", "thank you", "goodbye", "bye"}
	lowerInput := strings.ToLower(input)
	for _, phrase := range commonPhrases {
		if lowerInput == phrase || strings.HasPrefix(lowerInput, phrase+" ") {
			return true
		}
	}

	// Check if the input contains multiple words (common for natural language)
	words := strings.Fields(input)
	if len(words) >= 3 {
		return true
	}

	// Check for common question words at the beginning
	questionWords := []string{"what", "who", "where", "when", "why", "how", "is", "are", "can", "could", "would", "should", "do", "does", "did"}
	if len(words) > 0 {
		firstWord := strings.ToLower(words[0])
		for _, word := range questionWords {
			if firstWord == word {
				return true
			}
		}
	}

	// Check if the input is a single word that doesn't exist as a command
	// This is a heuristic to prevent treating unknown commands as shell commands
	if len(words) == 1 {
		// Check if the word exists as an executable in PATH
		_, err := exec.LookPath(words[0])
		if err != nil {
			// If the command doesn't exist, treat it as a natural language query
			return true
		}
	}

	// If none of the above conditions are met, it's likely not a natural language query
	return false
}

// isAgentTask determines if a query is likely to be a task for the agent
// rather than a simple question for the AI
func isAgentTask(input string) bool {
	// Convert to lowercase for case-insensitive matching
	lowerInput := strings.ToLower(input)

	// Check for action verbs and task-like phrases at the beginning
	actionPrefixes := []string{
		"create", "find", "list", "show", "get", "make", "setup",
		"install", "configure", "backup", "search", "organize",
		"clean", "delete", "remove", "update", "check", "analyze",
		"how to", "how do i", "can you", "please", "help me",
	}

	for _, prefix := range actionPrefixes {
		if strings.HasPrefix(lowerInput, prefix) {
			return true
		}
	}

	// Check for task-like keywords anywhere in the input
	taskKeywords := []string{
		"directory", "folder", "file", "system", "command", "script",
		"backup", "install", "config", "setup", "database", "server",
		"process", "service", "network", "disk", "memory", "cpu",
	}

	for _, keyword := range taskKeywords {
		if strings.Contains(lowerInput, keyword) {
			return true
		}
	}

	// Use regex to check for common task patterns
	taskPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\bfind\s+all\b`),
		regexp.MustCompile(`\bcreate\s+a\b`),
		regexp.MustCompile(`\bshow\s+me\b`),
		regexp.MustCompile(`\bhelp\s+me\s+with\b`),
		regexp.MustCompile(`\bhow\s+can\s+I\b`),
		regexp.MustCompile(`\bwhat's\s+the\s+best\s+way\s+to\b`),
	}

	for _, pattern := range taskPatterns {
		if pattern.MatchString(lowerInput) {
			return true
		}
	}

	return false
}

// isSpeedTestQuery determines if a query is related to internet speed testing
func isSpeedTestQuery(input string) bool {
	// Convert to lowercase for case-insensitive matching
	lowerInput := strings.ToLower(input)

	// Check for speed test related keywords and phrases
	speedTestKeywords := []string{
		"internet speed", "connection speed", "bandwidth", "speed test",
		"how fast is my internet", "check my internet speed", "test my connection",
		"network speed", "download speed", "upload speed", "internet performance",
	}

	for _, keyword := range speedTestKeywords {
		if strings.Contains(lowerInput, keyword) {
			return true
		}
	}

	return false
}
