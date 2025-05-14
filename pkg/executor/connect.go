package executor

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/agnath18/lumo/pkg/connect"
	"github.com/agnath18/lumo/pkg/nlp"
)

// executeConnectCommand handles file transfer connections
func (e *Executor) executeConnectCommand(cmd *nlp.Command) (*Result, error) {
	// Parse the intent
	intent := strings.TrimSpace(cmd.Intent)

	// Default values
	var downloadPath string
	port := 8080

	// Parse options
	args := strings.Fields(intent)
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Check for port option
		if arg == "--port" || arg == "-p" {
			if i+1 < len(args) {
				portNum, err := strconv.Atoi(args[i+1])
				if err == nil && portNum > 0 && portNum < 65536 {
					port = portNum
					i++ // Skip the next argument
				}
			}
		}

		// Check for download path option
		if arg == "--path" || arg == "-d" {
			if i+1 < len(args) {
				downloadPath = args[i+1]
				i++ // Skip the next argument
			}
		}
	}

	// Create a connect manager with the specified options
	connectManager := connect.NewConnectManager(downloadPath, port)

	// Check if we're in receive mode
	if strings.Contains(intent, "--receive") || strings.Contains(intent, "-r") {
		// Start a WebSocket server to receive files
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the receiver in the current goroutine
		err := connectManager.StartReceiver(ctx)
		if err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error starting receiver: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		return &Result{
			Output:     "Receiver stopped",
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// If we're here, we're in send mode
	// Check if we have a peer IP
	if len(args) == 0 || (len(args) > 0 && (args[0] == "--help" || args[0] == "-h")) {
		return &Result{
			Output: `
🔌 Lumo Connect - Duplex File Transfer Over LAN

Usage:
  lumo connect --receive [options]       Start a server to send and receive files
  lumo connect <peer-ip> [options]       Connect to a peer to send and receive files

Options:
  --port, -p <port>            Specify the port to use (default: 8080)
  --path, -d <directory>       Specify where to save received files (default: ~/Downloads)
  --help, -h                   Show this help message

Examples:
  lumo connect --receive                 Start a server on port 8080
  lumo connect --receive --port 9000     Start a server on port 9000
  lumo connect --receive --path /tmp     Save received files to /tmp
  lumo connect 192.168.1.5              Connect to peer at 192.168.1.5:8080
  lumo connect 192.168.1.5:9000         Connect to peer at 192.168.1.5:9000
  lumo connect 192.168.1.5 --path /tmp  Connect and save files to /tmp

Notes:
  - Both sides can send and receive files simultaneously
  - Drag and drop files into the terminal to send them
  - Type 'select' to open a file browser
  - Press Ctrl+C to stop the connection
`,
			IsError:    false,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Extract the peer IP from the arguments
	peerIP := args[0]

	// Skip if it's an option
	if strings.HasPrefix(peerIP, "-") {
		return &Result{
			Output:     "Invalid command. Use 'lumo connect --help' for usage information.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Check if the peer IP includes a port
	peerPort := port
	if strings.Contains(peerIP, ":") {
		parts := strings.Split(peerIP, ":")
		peerIP = parts[0]

		// Parse the port
		if len(parts) > 1 {
			portNum, err := strconv.Atoi(parts[1])
			if err == nil && portNum > 0 && portNum < 65536 {
				peerPort = portNum
			}
		}
	}

	// Validate the IP address
	if net.ParseIP(peerIP) == nil {
		return &Result{
			Output:     fmt.Sprintf("Invalid IP address: %s", peerIP),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Connect to the peer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := connectManager.ConnectToPeer(ctx, peerIP, peerPort)
	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Error connecting to peer: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	return &Result{
		Output:     "Connection closed",
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}
