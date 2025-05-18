package connect

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/agnath18K/lumo/pkg/discovery"
	"github.com/agnath18K/lumo/pkg/utils"
	"github.com/gorilla/websocket"
)

// FileTransferMessage represents a message for file transfer
type FileTransferMessage struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
	Size     int64  `json:"size,omitempty"`
	Content  []byte `json:"content,omitempty"`
	Progress int    `json:"progress,omitempty"` // Progress percentage (0-100)
}

// ConnectManager handles WebSocket connections for file transfers
type ConnectManager struct {
	upgrader     websocket.Upgrader
	server       *http.Server
	mode         string // "server", "client", or "duplex"
	downloadPath string // Custom download path
	port         int    // Custom port
	discoverer   discovery.Discoverer
	advertised   bool
	useChunked   bool // Whether to use chunked transfer for all files
}

// GetPort returns the current port
func (m *ConnectManager) GetPort() int {
	return m.port
}

// NewConnectManager creates a new connect manager
func NewConnectManager(downloadPath string, port int, useChunked ...bool) *ConnectManager {
	// Set default values if not provided
	if downloadPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			downloadPath = filepath.Join(homeDir, "Downloads")
		} else {
			downloadPath = "."
		}
	}

	if port <= 0 {
		port = 8080 // Default port
	}

	// Check if useChunked is provided
	chunkedTransfer := false
	if len(useChunked) > 0 {
		chunkedTransfer = useChunked[0]
	}

	// Create a new discoverer
	discoverer := discovery.NewDiscoverer()

	return &ConnectManager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections
			},
		},
		mode:         "duplex", // Default to duplex mode
		downloadPath: downloadPath,
		port:         port,
		discoverer:   discoverer,
		advertised:   false,
		useChunked:   chunkedTransfer,
	}
}

// StartReceiver starts a WebSocket server to receive files
func (m *ConnectManager) StartReceiver(ctx context.Context) error {
	// Set mode to server or duplex
	if m.mode != "duplex" {
		m.mode = "server"
	}

	// Create a new HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", m.handleWebSocket)

	// Get system information
	localIP, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP: %w", err)
	}

	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}

	// Check if the port is available
	if !utils.IsPortAvailable(m.port) {
		// Try to find an available port
		newPort, err := utils.FindAvailablePort(m.port, 100)
		if err != nil {
			return fmt.Errorf("port %d is already in use and no alternative ports are available: %w", m.port, err)
		}

		// Log the port change
		log.Printf("Port %d is already in use. Using port %d instead.", m.port, newPort)
		log.Printf("This could be due to another Lumo connect session or a Lumo server using this port.")
		log.Printf("To avoid this in the future, specify a different port with --port option.")
		log.Printf("%s", utils.GetPortRangeMessage("connect"))

		// Update the port
		m.port = newPort
	}

	// Create server
	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", m.port),
		Handler: mux,
	}

	// Start the discovery service
	if err := m.discoverer.Start(ctx); err != nil {
		log.Printf("Warning: Failed to start discovery service: %v", err)
	}

	// Advertise the service
	info := map[string]string{
		"hostname": hostname,
		"username": username,
		"mode":     m.mode,
	}
	if err := m.discoverer.Advertise(ctx, "Lumo Connect", m.port, info); err != nil {
		log.Printf("Warning: Failed to advertise service: %v", err)
	} else {
		m.advertised = true
	}

	// Print fancy header
	printFancyHeader()

	// Print connection information with colored output
	fmt.Printf("\033[1;36m") // Cyan color
	fmt.Printf("┌─────────────────────────────────────────────────┐\n")
	fmt.Printf("│ 🔌 \033[1;97mLumo Connect\033[1;36m                               │\n")
	fmt.Printf("├─────────────────────────────────────────────────┤\n")
	fmt.Printf("│ \033[1;97mStatus:\033[1;36m Server running                        │\n")
	fmt.Printf("│ \033[1;97mMode:\033[1;36m %s                                  │\n", m.mode)
	fmt.Printf("│ \033[1;97mIP Address:\033[1;36m %-33s │\n", localIP)
	fmt.Printf("│ \033[1;97mPort:\033[1;36m %-39d │\n", m.port)
	fmt.Printf("│ \033[1;97mHostname:\033[1;36m %-35s │\n", hostname)
	fmt.Printf("│ \033[1;97mUser:\033[1;36m %-39s │\n", username)
	fmt.Printf("│ \033[1;97mDownload Path:\033[1;36m %-30s │\n", m.downloadPath)
	fmt.Printf("│ \033[1;97mDiscoverable:\033[1;36m %-32v │\n", m.advertised)
	fmt.Printf("└─────────────────────────────────────────────────┘\n\n")

	if m.mode == "duplex" {
		fmt.Printf("📤 \033[1;97mYou can send files by:\033[1;36m\n")
		fmt.Printf("   • Dragging files into the terminal\n")
		fmt.Printf("   • Typing the full path to a file\n")
		fmt.Printf("   • Typing 'select' to open a file browser\n\n")
	}

	fmt.Printf("⏳ \033[1;97mWaiting for connections...\033[1;36m\n")
	fmt.Printf("🛑 \033[1;97mPress Ctrl+C to stop\033[1;36m\n\n")
	fmt.Printf("\033[0m") // Reset color

	// Start server in a goroutine
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
			if os.IsPermission(err) {
				log.Printf("This may be due to insufficient permissions to bind to port %d.", m.port)
				log.Printf("Try using a port number above 1024 with: --port <port>")
			}
		}
	}()

	// If in duplex mode, start reading from stdin for file paths
	if m.mode == "duplex" {
		go m.readStdinForFilePaths(nil) // nil connection means we'll send to any connected client
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Stop advertising the service
	if m.advertised {
		if err := m.discoverer.StopAdvertising(); err != nil {
			log.Printf("Warning: Failed to stop advertising service: %v", err)
		}
		m.advertised = false
	}

	// Stop the discovery service
	if err := m.discoverer.Stop(); err != nil {
		log.Printf("Warning: Failed to stop discovery service: %v", err)
	}

	return m.server.Shutdown(context.Background())
}

// printFancyHeader prints a fancy ASCII art header
func printFancyHeader() {
	header := `
 ██╗     ██╗   ██╗███╗   ███╗ ██████╗      ██████╗ ██████╗ ███╗   ██╗███╗   ██╗███████╗ ██████╗████████╗
 ██║     ██║   ██║████╗ ████║██╔═══██╗    ██╔════╝██╔═══██╗████╗  ██║████╗  ██║██╔════╝██╔════╝╚══██╔══╝
 ██║     ██║   ██║██╔████╔██║██║   ██║    ██║     ██║   ██║██╔██╗ ██║██╔██╗ ██║█████╗  ██║        ██║
 ██║     ██║   ██║██║╚██╔╝██║██║   ██║    ██║     ██║   ██║██║╚██╗██║██║╚██╗██║██╔══╝  ██║        ██║
 ███████╗╚██████╔╝██║ ╚═╝ ██║╚██████╔╝    ╚██████╗╚██████╔╝██║ ╚████║██║ ╚████║███████╗╚██████╗   ██║
 ╚══════╝ ╚═════╝ ╚═╝     ╚═╝ ╚═════╝      ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝ ╚═════╝   ╚═╝

`
	fmt.Printf("\033[1;36m%s\033[0m\n", header)
}

// ConnectToPeer connects to a peer to send files
func (m *ConnectManager) ConnectToPeer(ctx context.Context, peerIP string, peerPort int) error {
	// Set mode to client or duplex
	if m.mode != "duplex" {
		m.mode = "client"
	}

	// Create WebSocket URL
	url := fmt.Sprintf("ws://%s:%d/ws", peerIP, peerPort)

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to peer: %w", err)
	}
	defer conn.Close()

	// Get system information
	localIP, _ := getLocalIP()
	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}

	// Print fancy header
	printFancyHeader()

	// Print connection information with colored output
	fmt.Printf("\033[1;32m") // Green color
	fmt.Printf("┌─────────────────────────────────────────────────┐\n")
	fmt.Printf("│ 🔌 \033[1;97mLumo Connect\033[1;32m                               │\n")
	fmt.Printf("├─────────────────────────────────────────────────┤\n")
	fmt.Printf("│ \033[1;97mStatus:\033[1;32m Connected to peer                     │\n")
	fmt.Printf("│ \033[1;97mMode:\033[1;32m %s                                  │\n", m.mode)
	fmt.Printf("│ \033[1;97mLocal IP:\033[1;32m %-35s │\n", localIP)
	fmt.Printf("│ \033[1;97mPeer IP:\033[1;32m %-36s │\n", peerIP)
	fmt.Printf("│ \033[1;97mPeer Port:\033[1;32m %-34d │\n", peerPort)
	fmt.Printf("│ \033[1;97mHostname:\033[1;32m %-35s │\n", hostname)
	fmt.Printf("│ \033[1;97mUser:\033[1;32m %-39s │\n", username)
	fmt.Printf("│ \033[1;97mDownload Path:\033[1;32m %-30s │\n", m.downloadPath)
	fmt.Printf("└─────────────────────────────────────────────────┘\n\n")

	fmt.Printf("📤 \033[1;97mYou can send files by:\033[1;32m\n")
	fmt.Printf("   • Dragging files into the terminal\n")
	fmt.Printf("   • Typing the full path to a file\n")
	fmt.Printf("   • Typing 'select' to open a file browser\n\n")

	fmt.Printf("📥 \033[1;97mReceived files will be saved to:\033[1;32m %s\n\n", m.downloadPath)
	fmt.Printf("🛑 \033[1;97mPress Ctrl+C to disconnect\033[1;32m\n\n")
	fmt.Printf("\033[0m") // Reset color

	// Start a goroutine to read messages from the WebSocket
	go func() {
		for {
			var msg FileTransferMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				return
			}

			// Handle received message
			if msg.Type == "ack" {
				fmt.Printf("\033[1;32m✅ File %s received by peer\033[0m\n", msg.Filename)
			} else if msg.Type == "file" {
				// Save the file
				filename := m.saveFile(msg.Filename, msg.Content)

				// Send acknowledgment
				ack := FileTransferMessage{
					Type:     "ack",
					Filename: msg.Filename,
				}
				if err := conn.WriteJSON(ack); err != nil {
					log.Printf("Error sending acknowledgment: %v", err)
				}

				// Format file size
				sizeStr := formatFileSize(int64(len(msg.Content)))
				fmt.Printf("\033[1;36m📥 Received file: %s (%s)\033[0m\n", filename, sizeStr)
			}
		}
	}()

	// Read from stdin for file paths
	return m.readStdinForFilePaths(conn)
}

// readStdinForFilePaths reads file paths from stdin and sends files
// If conn is nil, it will send to all connected clients (server mode)
func (m *ConnectManager) readStdinForFilePaths(conn *websocket.Conn) error {
	// Print instructions for manual file entry
	fmt.Printf("\033[1;33mℹ️ You can type the full path to a file and press Enter\033[0m\n")
	fmt.Printf("\033[1;33mℹ️ Type 'select' to open a file browser\033[0m\n")

	// Read from stdin for file paths
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		filePath := scanner.Text()

		// Handle special formats from drag-and-drop
		// Some terminals prefix with "file://" or have URL encoding
		if strings.HasPrefix(filePath, "file://") {
			filePath = strings.TrimPrefix(filePath, "file://")
		}

		// Trim any quotes or whitespace that might be around the path
		filePath = strings.Trim(filePath, "\"' \t\n\r")

		// Skip empty lines
		if filePath == "" {
			continue
		}

		// Check for special commands
		if filePath == "select" {
			// Open a file dialog using zenity if available
			selectedFile, err := openFileDialog()
			if err != nil {
				fmt.Printf("\033[1;31m❌ Error opening file dialog: %v\033[0m\n", err)
				fmt.Printf("\033[1;33mℹ️ Try dragging and dropping a file instead\033[0m\n")
			} else if selectedFile != "" {
				// Try to send the selected file
				if conn != nil {
					// Send to specific connection
					err := m.sendFile(conn, selectedFile)
					if err != nil {
						fmt.Printf("\033[1;31m❌ Error sending file: %v\033[0m\n", err)
					}
				} else {
					// Send to all connected clients
					m.sendFileToAllClients(selectedFile)
				}
			}
			continue
		}

		// Check if this looks like a file path
		if strings.Contains(filePath, "/") || strings.Contains(filePath, "\\") {
			// Check if the file exists
			if _, err := os.Stat(filePath); err == nil {
				// Try to send the file
				if conn != nil {
					// Send to specific connection
					err := m.sendFile(conn, filePath)
					if err != nil {
						fmt.Printf("\033[1;31m❌ Error sending file: %v\033[0m\n", err)
					}
				} else {
					// Send to all connected clients
					m.sendFileToAllClients(filePath)
				}
			} else {
				fmt.Printf("\033[1;33m⚠️ File not found: %s\033[0m\n", filePath)
				fmt.Printf("\033[1;33mℹ️ Make sure to provide the full path to the file\033[0m\n")
				fmt.Printf("\033[1;33mℹ️ Type 'select' to open a file browser\033[0m\n")
			}
		} else {
			// Print a message to remind the user to drag and drop files
			fmt.Printf("\033[1;33mℹ️ Drag and drop a file into the terminal or type the full path\033[0m\n")
			fmt.Printf("\033[1;33mℹ️ Type 'select' to open a file browser\033[0m\n")
		}
	}

	return nil
}

// Global variable to store active connections
var activeConnections = make(map[*websocket.Conn]bool)
var connectionsMutex = &sync.Mutex{}

// handleWebSocket handles WebSocket connections
func (m *ConnectManager) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Register connection
	connectionsMutex.Lock()
	activeConnections[conn] = true
	connectionsMutex.Unlock()

	// Ensure connection is removed when closed
	defer func() {
		conn.Close()
		connectionsMutex.Lock()
		delete(activeConnections, conn)
		connectionsMutex.Unlock()
	}()

	// Get client IP
	clientIP := r.RemoteAddr
	fmt.Printf("\033[1;36m🔗 New connection from %s\033[0m\n", clientIP)

	// Handle WebSocket connection
	for {
		var msg FileTransferMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Handle file transfer message
		if msg.Type == "file" {
			// Save the file
			filename := m.saveFile(msg.Filename, msg.Content)

			// Send acknowledgment
			ack := FileTransferMessage{
				Type:     "ack",
				Filename: msg.Filename,
			}
			if err := conn.WriteJSON(ack); err != nil {
				log.Printf("Error sending acknowledgment: %v", err)
			}

			// Format file size
			sizeStr := formatFileSize(int64(len(msg.Content)))
			fmt.Printf("\033[1;36m📥 Received file: %s (%s)\033[0m\n", filename, sizeStr)
		}
	}
}

// sendFileToAllClients sends a file to all connected clients
func (m *ConnectManager) sendFileToAllClients(filePath string) {
	// Get the number of active connections
	connectionsMutex.Lock()
	numConnections := len(activeConnections)
	connectionsMutex.Unlock()

	// Check if there are any connections
	if numConnections == 0 {
		fmt.Printf("\033[1;33m⚠️ No connected clients to send file to\033[0m\n")
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("\033[1;31m❌ Error opening file: %v\033[0m\n", err)
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("\033[1;31m❌ Error getting file info: %v\033[0m\n", err)
		return
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		fmt.Printf("\033[1;31m❌ Not a regular file\033[0m\n")
		return
	}

	// Get base filename
	filename := filepath.Base(filePath)

	// Format file size
	sizeStr := formatFileSize(fileInfo.Size())
	fmt.Printf("\033[1;32m📤 Sending file: %s (%s) to %d clients...\033[0m\n", filename, sizeStr, numConnections)

	// Check if we should use chunked transfer
	if m.useChunked || fileInfo.Size() > 10*1024*1024 { // Use chunked if explicitly requested or file is larger than 10MB
		// For large files, use chunked transfer
		fmt.Printf("\033[1;33mℹ️ Large file detected. Using chunked transfer...\033[0m\n")

		// Get local IP
		localIP, err := getLocalIP()
		if err != nil {
			localIP = "localhost"
		}

		// Create a chunked client
		client := NewChunkedClient(fmt.Sprintf("http://%s:7531", localIP), m.downloadPath, DefaultChunkSize)

		// Upload the file
		resultPath, err := client.UploadFile(filePath, nil)
		if err != nil {
			fmt.Printf("\033[1;31m❌ Error uploading file using chunked transfer: %v\033[0m\n", err)
			return
		}

		fmt.Printf("\033[1;32m📤 File uploaded successfully to: %s\033[0m\n", resultPath)
		return
	}

	// For small files, use WebSocket transfer
	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("\033[1;31m❌ Error reading file: %v\033[0m\n", err)
		return
	}

	// Create file transfer message
	msg := FileTransferMessage{
		Type:     "file",
		Filename: filename,
		Size:     fileInfo.Size(),
		Content:  content,
	}

	// Send to all connections
	connectionsMutex.Lock()
	for conn := range activeConnections {
		// Send the message
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Printf("\033[1;31m❌ Error sending file to a client: %v\033[0m\n", err)
			continue
		}
	}
	connectionsMutex.Unlock()

	fmt.Printf("\033[1;32m📤 File sent to all connected clients!\033[0m\n")
}

// sendFile sends a file over WebSocket
func (m *ConnectManager) sendFile(conn *websocket.Conn, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("not a regular file")
	}

	// Get base filename
	filename := filepath.Base(filePath)

	// Format file size
	sizeStr := formatFileSize(fileInfo.Size())
	fmt.Printf("\033[1;32m📤 Sending file: %s (%s)...\033[0m\n", filename, sizeStr)

	// Check if we should use chunked transfer
	if m.useChunked || fileInfo.Size() > 10*1024*1024 { // Use chunked if explicitly requested or file is larger than 10MB
		// For large files, use chunked transfer
		fmt.Printf("\033[1;33mℹ️ Large file detected. Using chunked transfer...\033[0m\n")

		// Get peer IP and port from WebSocket connection
		peerAddr := conn.RemoteAddr().String()
		peerIP, _, err := net.SplitHostPort(peerAddr)
		if err != nil {
			// If we can't get the peer IP from the connection, try to get it from the local connection
			peerIP, err = getLocalIP()
			if err != nil {
				peerIP = "localhost"
			}
		}

		// Create a chunked client
		client := NewChunkedClient(fmt.Sprintf("http://%s:7531", peerIP), m.downloadPath, DefaultChunkSize)

		// Upload the file
		resultPath, err := client.UploadFile(filePath, nil)
		if err != nil {
			return fmt.Errorf("failed to upload file using chunked transfer: %w", err)
		}

		fmt.Printf("\033[1;32m📤 File uploaded successfully to: %s\033[0m\n", resultPath)
		return nil
	}

	// For small files, use WebSocket transfer
	// Show progress bar
	fmt.Printf("\033[1;32m[                    ] 0%%\033[0m")
	fmt.Printf("\r")

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create file transfer message
	msg := FileTransferMessage{
		Type:     "file",
		Filename: filename,
		Size:     fileInfo.Size(),
		Content:  content,
	}

	// Send the message
	if err := conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("failed to send file: %w", err)
	}

	// Update progress bar to 100%
	fmt.Printf("\033[1;32m[====================] 100%%\033[0m\n")
	fmt.Printf("\033[1;32m📤 File sent successfully!\033[0m\n")
	return nil
}

// saveFile saves a file to the downloads directory
func (m *ConnectManager) saveFile(filename string, content []byte) string {
	// Create the download directory if it doesn't exist
	err := os.MkdirAll(m.downloadPath, 0755)
	if err != nil {
		log.Printf("Error creating download directory: %v", err)
		// Fall back to current directory
		m.downloadPath = "."
	}

	// Create timestamp
	timestamp := time.Now().Format("20060102_150405")

	// Create filename with timestamp
	baseFilename := filepath.Base(filename)
	ext := filepath.Ext(baseFilename)
	name := strings.TrimSuffix(baseFilename, ext)
	newFilename := fmt.Sprintf("%s_%s%s", name, timestamp, ext)

	// Create full path
	filePath := filepath.Join(m.downloadPath, newFilename)

	// Write file
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		return filename
	}

	return filePath
}

// getLocalIP returns the local IP address
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "127.0.0.1", nil
}

// formatFileSize formats a file size in bytes to a human-readable string
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// openFileDialog opens a file selection dialog
func openFileDialog() (string, error) {
	// Try to use zenity if available
	cmd := exec.Command("zenity", "--file-selection", "--title=Select a file to send")
	output, err := cmd.Output()
	if err != nil {
		// Try to use kdialog if zenity is not available
		cmd = exec.Command("kdialog", "--getopenfilename", ".", "All Files (*)")
		output, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("no file dialog available (install zenity or kdialog)")
		}
	}

	// Trim newline from output
	return strings.TrimSpace(string(output)), nil
}

// DiscoverServices discovers available Lumo Connect services on the network
func (m *ConnectManager) DiscoverServices(ctx context.Context) ([]discovery.Service, error) {
	// Start the discovery service if not already started
	if m.discoverer == nil {
		m.discoverer = discovery.NewDiscoverer()
		if err := m.discoverer.Start(ctx); err != nil {
			return nil, fmt.Errorf("failed to start discovery service: %w", err)
		}
	}

	// Browse for services
	services, err := m.discoverer.Browse(ctx, discovery.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to browse for services: %w", err)
	}

	return services, nil
}

// PrintDiscoveredServices prints a list of discovered services
func (m *ConnectManager) PrintDiscoveredServices(services []discovery.Service) {
	if len(services) == 0 {
		fmt.Printf("\033[1;33mNo Lumo Connect services found on the network\033[0m\n")
		return
	}

	fmt.Printf("\033[1;36m") // Cyan color
	fmt.Printf("┌─────────────────────────────────────────────────┐\n")
	fmt.Printf("│ 🔍 \033[1;97mDiscovered Lumo Connect Services\033[1;36m             │\n")
	fmt.Printf("├─────────────────────────────────────────────────┤\n")

	for i, service := range services {
		fmt.Printf("│ \033[1;97m%d.\033[1;36m %-45s │\n", i+1, service.Name)
		fmt.Printf("│   \033[1;97mIP:\033[1;36m %-43s │\n", service.IP)
		fmt.Printf("│   \033[1;97mPort:\033[1;36m %-41d │\n", service.Port)
		if username, ok := service.Info["username"]; ok {
			fmt.Printf("│   \033[1;97mUser:\033[1;36m %-41s │\n", username)
		}
		if i < len(services)-1 {
			fmt.Printf("├─────────────────────────────────────────────────┤\n")
		}
	}

	fmt.Printf("└─────────────────────────────────────────────────┘\n")
	fmt.Printf("\033[0m") // Reset color
}
