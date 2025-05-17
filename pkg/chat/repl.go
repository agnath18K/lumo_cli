package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agnath18K/lumo/pkg/ai"
	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/utils"
)

// REPL handles the interactive chat REPL mode
type REPL struct {
	config     *config.Config
	manager    *Manager
	reader     *bufio.Reader
	aiClient   ai.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewREPL creates a new REPL instance
func NewREPL(cfg *config.Config, manager *Manager, aiClient ai.Client) *REPL {
	ctx, cancel := context.WithCancel(context.Background())
	return &REPL{
		config:     cfg,
		manager:    manager,
		reader:     bufio.NewReader(os.Stdin),
		aiClient:   aiClient,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Start starts the REPL loop
func (r *REPL) Start() (string, error) {
	// Display welcome message
	r.displayWelcome()

	// Get or create the active conversation
	conv := r.manager.GetActiveConversation()

	// Main REPL loop
	for {
		// Display prompt
		fmt.Print("\nchat> ")
		input, err := r.reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Parse the command
		parts := strings.SplitN(input, " ", 2)
		cmd := strings.ToLower(parts[0])
		args := ""
		if len(parts) > 1 {
			args = strings.TrimSpace(parts[1])
		}

		// Process the command
		switch cmd {
		case "exit", "quit":
			// Exit the REPL
			fmt.Println("Exiting chat mode.")
			return "Chat session ended.", nil

		case "help":
			// Display help
			r.displayHelp()

		case "clear":
			// Clear the conversation history
			conv.Clear()
			fmt.Println("Conversation history cleared.")

		case "history":
			// Display conversation history
			r.displayHistory(conv)

		case "new":
			// Start a new conversation
			conv = r.manager.StartNewConversation()
			fmt.Println("Started a new conversation.")

		case "list":
			// List all conversations
			r.listConversations()

		case "switch":
			// Switch to another conversation
			if args == "" {
				fmt.Println("Error: Conversation ID required.")
				continue
			}
			if r.manager.SetActiveConversation(args) {
				conv = r.manager.GetConversation(args)
				fmt.Printf("Switched to conversation %s.\n", args)
			} else {
				fmt.Printf("Error: Conversation %s not found.\n", args)
			}

		case "delete":
			// Delete a conversation
			if args == "" {
				fmt.Println("Error: Conversation ID required.")
				continue
			}
			if r.manager.DeleteConversation(args) {
				fmt.Printf("Deleted conversation %s.\n", args)
				// If we deleted the active conversation, get a new one
				conv = r.manager.GetActiveConversation()
			} else {
				fmt.Printf("Error: Conversation %s not found.\n", args)
			}

		default:
			// Treat as a message to the AI
			fmt.Println(ai.ThinkingIndicator)

			// Add the user message to the conversation
			conv.AddUserMessage(input)

			// Process the message
			response, err := r.manager.ProcessMessage(r.ctx, input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			// Display the response without box formatting
			// Clean up markdown formatting for better terminal display
			cleanResponse := utils.CleanMarkdown(response)
			fmt.Println("\n" + cleanResponse)
		}
	}
}

// displayWelcome displays a welcome message
func (r *REPL) displayWelcome() {
	fmt.Println("\n🗣️  Welcome to Lumo Chat Mode!")
	fmt.Println("Type your message or use one of the commands below.")
	fmt.Println("Type 'help' for more information.")
}

// displayHelp displays help information
func (r *REPL) displayHelp() {
	fmt.Println("\n📚 Chat Mode Help:")
	fmt.Println("  <message>            - Send a message to the AI")
	fmt.Println("  help                 - Show this help message")
	fmt.Println("  clear                - Clear the conversation history")
	fmt.Println("  history              - Display conversation history")
	fmt.Println("  new                  - Start a new conversation")
	fmt.Println("  list                 - List all conversations")
	fmt.Println("  switch <id>          - Switch to another conversation")
	fmt.Println("  delete <id>          - Delete a conversation")
	fmt.Println("  exit, quit           - Exit chat mode")
}

// displayHistory displays the conversation history
func (r *REPL) displayHistory(conv *Conversation) {
	messages := conv.GetMessages()

	// Skip the system message
	startIdx := 0
	for i, msg := range messages {
		if msg.Role != RoleSystem {
			startIdx = i
			break
		}
	}

	if startIdx >= len(messages) {
		fmt.Println("No messages in the conversation yet.")
		return
	}

	fmt.Println("\n📜 Conversation History:")
	fmt.Println("─────────────────────────────────────────")

	for i := startIdx; i < len(messages); i++ {
		msg := messages[i]
		timestamp := msg.Timestamp.Format("15:04:05")

		switch msg.Role {
		case RoleUser:
			fmt.Printf("🧑 [%s] You: %s\n\n", timestamp, msg.Content)
		case RoleAssistant:
			fmt.Printf("🐦 [%s] Lumo: %s\n\n", timestamp, msg.Content)
		}
	}
}

// listConversations lists all conversations
func (r *REPL) listConversations() {
	convIDs := r.manager.ListConversations()
	activeConv := r.manager.GetActiveConversation()

	if len(convIDs) == 0 {
		fmt.Println("No conversations found.")
		return
	}

	fmt.Println("\n💬 Conversations:")
	fmt.Println("─────────────────────────────────────────")

	for _, id := range convIDs {
		conv := r.manager.GetConversation(id)
		if conv == nil {
			continue
		}

		// Get the first user message as a preview
		var preview string
		for _, msg := range conv.GetMessages() {
			if msg.Role == RoleUser {
				preview = msg.Content
				if len(preview) > 30 {
					preview = preview[:27] + "..."
				}
				break
			}
		}

		// Mark the active conversation
		activeMarker := " "
		if activeConv != nil && id == activeConv.ID {
			activeMarker = "*"
		}

		// Get the creation time
		var creationTime time.Time
		if len(conv.Messages) > 0 {
			creationTime = conv.Messages[0].Timestamp
		} else {
			creationTime = time.Now()
		}

		// Format the creation time
		timeAgo := utils.FormatTimeAgo(creationTime)

		fmt.Printf("%s %s - %s (%s)\n", activeMarker, id, preview, timeAgo)
	}
}
