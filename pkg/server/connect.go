package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/agnath18K/lumo/pkg/auth"
	"github.com/agnath18K/lumo/pkg/connect"
	"github.com/agnath18K/lumo/pkg/discovery"
	"github.com/agnath18K/lumo/pkg/utils"
	"github.com/gorilla/websocket"
)

// ConnectRequest represents a request to connect to a peer
type ConnectRequest struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Path string `json:"path,omitempty"`
}

// ServerRequest represents a request to start a server
type ServerRequest struct {
	Port int    `json:"port"`
	Path string `json:"path,omitempty"`
}

// DiscoverResponse represents a response from the discover endpoint
type DiscoverResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error,omitempty"`
	Devices []discovery.Service `json:"devices,omitempty"`
}

// ConnectResponse represents a response from the connect endpoint
type ConnectResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	IP      string `json:"ip,omitempty"`
	Port    int    `json:"port,omitempty"`
}

// FileTransferMessage represents a message for file transfer
type FileTransferMessage struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
	Size     int64  `json:"size,omitempty"`
	Content  []byte `json:"content,omitempty"`
	Progress int    `json:"progress,omitempty"` // Progress percentage (0-100)
}

var (
	// activeConnectManager is the active connect manager
	activeConnectManager *connect.ConnectManager
	// activeConnectContext is the context for the active connect manager
	activeConnectContext context.Context
	// activeConnectCancel is the cancel function for the active connect context
	activeConnectCancel context.CancelFunc
	// upgrader is the websocket upgrader
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections
		},
	}
	// activeWebSockets is a map of active websocket connections
	activeWebSockets = make(map[*websocket.Conn]bool)
)

// handleConnectDiscover handles the /api/v1/connect/discover endpoint
func (s *Server) handleConnectDiscover(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a connect manager
	connectManager := connect.NewConnectManager("", 0)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Discover services
	services, err := connectManager.DiscoverServices(ctx)
	if err != nil {
		// Return error response
		response := DiscoverResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to discover services: %v", err),
			Devices: []discovery.Service{}, // Always include an empty array
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return success response
	response := DiscoverResponse{
		Success: true,
		Devices: services,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleConnectStartServer handles the /api/v1/connect/start-server endpoint
func (s *Server) handleConnectStartServer(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request ServerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if there's already an active connection
	if activeConnectManager != nil {
		// Stop the active connection
		if activeConnectCancel != nil {
			activeConnectCancel()
		}
		activeConnectManager = nil
		activeConnectCancel = nil
	}

	// Create a connect manager
	connectManager := connect.NewConnectManager(request.Path, request.Port)

	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start the server in a goroutine
	var startErr error
	errChan := make(chan error, 1)

	go func() {
		err := connectManager.StartReceiver(ctx)
		if err != nil {
			log.Printf("Error starting server: %v", err)
			errChan <- err
		}
	}()

	// Wait a bit for the server to start or for an error
	select {
	case startErr = <-errChan:
		// Handle error
		if startErr != nil {
			// Cancel the context to clean up resources
			cancel()

			// Return error response
			response := ConnectResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to start connect server: %v", startErr),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	case <-time.After(500 * time.Millisecond):
		// Continue with normal operation
	}

	// Store the connect manager and context
	activeConnectManager = connectManager
	activeConnectContext = ctx
	activeConnectCancel = cancel

	// Get local IP
	localIP, err := getLocalIP()
	if err != nil {
		localIP = "127.0.0.1"
	}

	// Check if the port might have changed due to port conflict
	actualPort := connectManager.GetPort()
	if actualPort != request.Port {
		log.Printf("Port %d was already in use. Using port %d instead.", request.Port, actualPort)
		log.Printf("This could be due to another Lumo server instance or a Lumo connect session using this port.")
		log.Printf("To avoid this in the future, configure a different port with: lumo config:server port <port>")
		log.Printf("%s", utils.GetPortRangeMessage("server"))
	}

	// Return success response
	response := ConnectResponse{
		Success: true,
		IP:      localIP,
		Port:    actualPort, // Use the actual port which might have changed
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleConnectToPeer handles the /api/v1/connect/connect-to-peer endpoint
func (s *Server) handleConnectToPeer(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request ConnectRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if there's already an active connection
	if activeConnectManager != nil {
		// Stop the active connection
		if activeConnectCancel != nil {
			activeConnectCancel()
		}
		activeConnectManager = nil
		activeConnectCancel = nil
	}

	// Create a connect manager
	connectManager := connect.NewConnectManager(request.Path, 0)

	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Connect to the peer in a goroutine
	go func() {
		err := connectManager.ConnectToPeer(ctx, request.IP, request.Port)
		if err != nil {
			log.Printf("Error connecting to peer: %v", err)
		}
	}()

	// Wait a bit for the connection to establish
	time.Sleep(500 * time.Millisecond)

	// Store the connect manager and context
	activeConnectManager = connectManager
	activeConnectContext = ctx
	activeConnectCancel = cancel

	// Return success response
	response := ConnectResponse{
		Success: true,
		IP:      request.IP,
		Port:    request.Port,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleConnectDisconnect handles the /api/v1/connect/disconnect endpoint
func (s *Server) handleConnectDisconnect(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if there's an active connection
	if activeConnectManager == nil {
		// Return error response
		response := ConnectResponse{
			Success: false,
			Error:   "No active connection",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Stop the active connection
	if activeConnectCancel != nil {
		activeConnectCancel()
	}
	activeConnectManager = nil
	activeConnectCancel = nil

	// Return success response
	response := ConnectResponse{
		Success: true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleConnectSendFile handles the /api/v1/connect/send-file endpoint
func (s *Server) handleConnectSendFile(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if there's an active connection
	if activeConnectManager == nil {
		// Return error response
		response := ConnectResponse{
			Success: false,
			Error:   "No active connection",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(100 << 20) // 100 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "lumo-connect-*")
	if err != nil {
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to copy file", http.StatusInternalServerError)
		return
	}

	// Rewind the temporary file to the beginning
	if _, err := tempFile.Seek(0, 0); err != nil {
		http.Error(w, "Failed to rewind temporary file", http.StatusInternalServerError)
		return
	}

	// Get the filename from the form
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file header", http.StatusBadRequest)
		return
	}

	// Read the file content
	fileContent, err := io.ReadAll(tempFile)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// Send the file to all connected clients
	// In a real implementation, this would use the connect manager to send the file
	// For now, we'll just broadcast a message to all connected WebSockets
	message := connect.FileTransferMessage{
		Type:     "file",
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		Content:  fileContent,
	}

	// Broadcast to all connected WebSockets
	for conn := range activeWebSockets {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error sending file to WebSocket: %v", err)
		}
	}

	// Return success response
	response := ConnectResponse{
		Success: true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleConnectWebSocket handles the /api/v1/connect/ws endpoint
func (s *Server) handleConnectWebSocket(w http.ResponseWriter, r *http.Request) {
	// Skip authentication if it's disabled
	if !s.config.EnableAuth {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading connection: %v", err)
			return
		}

		// Register connection
		activeWebSockets[conn] = true

		// Ensure connection is removed when closed
		defer func() {
			conn.Close()
			delete(activeWebSockets, conn)
		}()

		// Handle WebSocket connection
		for {
			// Just keep the connection open
			_, _, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				break
			}
		}
		return
	}

	// Get token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	// Validate the token
	claims, err := s.authenticator.ValidateToken(token)
	if err != nil {
		if err == auth.ErrTokenExpired {
			http.Error(w, "Token expired", http.StatusUnauthorized)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
		return
	}

	// Add the username to the request context
	ctx := context.WithValue(r.Context(), userContextKey, claims.Username)
	r = r.WithContext(ctx)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Register connection
	activeWebSockets[conn] = true

	// Ensure connection is removed when closed
	defer func() {
		conn.Close()
		delete(activeWebSockets, conn)
	}()

	// Handle WebSocket connection
	for {
		// Just keep the connection open
		_, _, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
	}
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
