package executor

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/agnath18K/lumo/pkg/connect"
	"github.com/agnath18K/lumo/pkg/nlp"
	"github.com/agnath18K/lumo/pkg/utils"
)

// executeConnectCommand handles file transfer connections
func (e *Executor) executeConnectCommand(cmd *nlp.Command) (*Result, error) {
	// Parse the intent
	intent := strings.TrimSpace(cmd.Intent)

	// Default values
	var downloadPath string
	port := 8080
	useChunked := false

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

		// Check for chunked transfer option
		if arg == "--chunked" || arg == "-c" {
			useChunked = true
		}
	}

	// Create a connect manager with the specified options
	connectManager := connect.NewConnectManager(downloadPath, port, useChunked)

	// Check if we're in receive mode
	if strings.Contains(intent, "--receive") || strings.Contains(intent, "-r") {
		// Start a WebSocket server to receive files
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the receiver in the current goroutine
		err := connectManager.StartReceiver(ctx)
		if err != nil {
			// Check if it's a port conflict error
			if strings.Contains(err.Error(), "port") && strings.Contains(err.Error(), "already in use") {
				return &Result{
					Output: fmt.Sprintf("Error: %v\n\n"+
						"This could be due to:\n"+
						"1. Another Lumo connect session running\n"+
						"2. The Lumo server using this port\n"+
						"3. Another application using this port\n\n"+
						"Try using a different port with: lumo connect --receive --port <port>\n"+
						"%s", err, utils.GetPortRangeMessage("connect")),
					IsError:    true,
					CommandRun: cmd.RawInput,
				}, nil
			}

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

	// Check for discovery mode
	if strings.Contains(intent, "--discover") || strings.Contains(intent, "-d") {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Discover services
		services, err := connectManager.DiscoverServices(ctx)
		if err != nil {
			return &Result{
				Output:     fmt.Sprintf("Error discovering services: %v", err),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

		// Print discovered services
		var output strings.Builder
		if len(services) == 0 {
			output.WriteString("No Lumo Connect services found on the network.\n")
			output.WriteString("Make sure other devices are running 'lumo connect --receive'.\n")
		} else {
			output.WriteString(fmt.Sprintf("Found %d Lumo Connect services on the network:\n\n", len(services)))
			for i, service := range services {
				output.WriteString(fmt.Sprintf("%d. %s\n", i+1, service.Name))
				output.WriteString(fmt.Sprintf("   IP: %s\n", service.IP))
				output.WriteString(fmt.Sprintf("   Port: %d\n", service.Port))
				if username, ok := service.Info["username"]; ok {
					output.WriteString(fmt.Sprintf("   User: %s\n", username))
				}
				output.WriteString("\n")
			}
			output.WriteString("To connect to a service, use: lumo connect <ip-address>[:<port>]\n")
		}

		return &Result{
			Output:     output.String(),
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
  lumo connect --discover, -d            Discover Lumo Connect services on the network
  lumo connect <peer-ip> [options]       Connect to a peer to send and receive files

Options:
  --port, -p <port>            Specify the port to use (default: 8080)
  --path, -d <directory>       Specify where to save received files (default: ~/Downloads)
  --chunked, -c                Use chunked transfer for all files (better for large files)
  --help, -h                   Show this help message

Examples:
  lumo connect --receive                 Start a server on port 8080
  lumo connect --receive --port 9000     Start a server on port 9000
  lumo connect --receive --path /tmp     Save received files to /tmp
  lumo connect --discover                Discover available Lumo Connect services
  lumo connect 192.168.1.5              Connect to peer at 192.168.1.5:8080
  lumo connect 192.168.1.5:9000         Connect to peer at 192.168.1.5:9000
  lumo connect 192.168.1.5 --path /tmp  Connect and save files to /tmp
  lumo connect 192.168.1.5 --chunked    Connect and use chunked transfer for all files

Notes:
  - Both sides can send and receive files simultaneously
  - Drag and drop files into the terminal to send them
  - Type 'select' to open a file browser
  - Press Ctrl+C to stop the connection
  - Files larger than 10MB automatically use chunked transfer
  - Use --chunked option for better performance with large files
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
		// Provide more helpful error messages
		if strings.Contains(err.Error(), "connection refused") {
			return &Result{
				Output: fmt.Sprintf("Error: Could not connect to %s:%d\n\n"+
					"Possible reasons:\n"+
					"1. No Lumo connect server is running at that address\n"+
					"2. A firewall is blocking the connection\n"+
					"3. The port number is incorrect\n\n"+
					"Try running 'lumo connect --discover' to find available services.",
					peerIP, peerPort),
				IsError:    true,
				CommandRun: cmd.RawInput,
			}, nil
		}

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
