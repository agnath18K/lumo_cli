package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/agnath18K/lumo/pkg/assets"
	"github.com/agnath18K/lumo/pkg/auth"
	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
	"github.com/agnath18K/lumo/pkg/nlp"
	"github.com/agnath18K/lumo/pkg/utils"
	"github.com/agnath18K/lumo/pkg/version"
)

// Server represents the REST API server for Lumo
type Server struct {
	config        *config.Config
	executor      *executor.Executor
	server        *http.Server
	isDaemon      bool
	authenticator *auth.Authenticator
}

// CommandRequest represents a request to execute a command
type CommandRequest struct {
	Command string            `json:"command"`
	Type    string            `json:"type,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
}

// CommandResponse represents the response from executing a command
type CommandResponse struct {
	Success    bool   `json:"success"`
	Output     string `json:"output"`
	CommandRun string `json:"command_run"`
	Error      string `json:"error,omitempty"`
}

// StatusResponse represents the server status response
type StatusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Username     string `json:"username"`
	ExpiresIn    int    `json:"expires_in"` // Seconds until token expires
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// New creates a new REST server instance
func New(cfg *config.Config, exec *executor.Executor) *Server {
	// Create the authenticator
	homeDir, err := os.UserHomeDir()
	credentialsDir := filepath.Join(homeDir, ".config", "lumo")
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		credentialsDir = ".config/lumo"
	}

	authenticator, err := auth.NewAuthenticator(cfg.JWTSecret, credentialsDir)
	if err != nil {
		log.Printf("Error creating authenticator: %v", err)
	}

	return &Server{
		config:        cfg,
		executor:      exec,
		isDaemon:      false,
		authenticator: authenticator,
	}
}

// NewDaemon creates a new REST server instance in daemon mode
func NewDaemon(cfg *config.Config, exec *executor.Executor) *Server {
	// Create the authenticator
	homeDir, err := os.UserHomeDir()
	credentialsDir := filepath.Join(homeDir, ".config", "lumo")
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		credentialsDir = ".config/lumo"
	}

	authenticator, err := auth.NewAuthenticator(cfg.JWTSecret, credentialsDir)
	if err != nil {
		log.Printf("Error creating authenticator: %v", err)
	}

	return &Server{
		config:        cfg,
		executor:      exec,
		isDaemon:      true,
		authenticator: authenticator,
	}
}

// Start starts the REST server
func (s *Server) Start() error {
	// Initialize the authenticator
	if err := s.authenticator.InitializeCredentialsStore(); err != nil {
		log.Printf("Error initializing credentials store: %v", err)
	}

	// Check if we need to create a default user
	hasUsers, err := s.authenticator.HasUsers()
	if err != nil {
		log.Printf("Error checking for users: %v", err)
	} else if !hasUsers {
		// Create a default user
		defaultUsername := "admin"
		defaultPassword := "lumo"
		if err := s.authenticator.AddUser(defaultUsername, defaultPassword); err != nil {
			log.Printf("Error creating default user: %v", err)
		} else {
			log.Printf("Created default user '%s' with password '%s'", defaultUsername, defaultPassword)
			log.Printf("Please change this password immediately using the web interface or API")
		}
	}

	// Create a new router
	mux := http.NewServeMux()

	// Create a middleware chain
	var handler http.Handler = mux
	if s.config.EnableAuth {
		handler = s.AuthMiddleware(mux)
	}

	// Register API routes
	mux.HandleFunc("/api/v1/execute", s.handleExecute)
	mux.HandleFunc("/api/v1/status", s.handleStatus)

	// Register authentication routes
	mux.HandleFunc("/api/v1/auth/login", s.handleLogin)
	mux.HandleFunc("/api/v1/auth/refresh", s.handleRefreshToken)
	mux.HandleFunc("/api/v1/auth/change-password", s.handleChangePassword)

	// Register Connect API routes
	mux.HandleFunc("/api/v1/connect/discover", s.handleConnectDiscover)
	mux.HandleFunc("/api/v1/connect/start-server", s.handleConnectStartServer)
	mux.HandleFunc("/api/v1/connect/connect-to-peer", s.handleConnectToPeer)
	mux.HandleFunc("/api/v1/connect/disconnect", s.handleConnectDisconnect)
	mux.HandleFunc("/api/v1/connect/send-file", s.handleConnectSendFile)
	mux.HandleFunc("/api/v1/connect/ws", s.handleConnectWebSocket)

	// Register Chunked File Transfer API routes
	mux.HandleFunc("/api/v1/connect/upload/init", s.handleInitUpload)
	mux.HandleFunc("/api/v1/connect/upload/chunk", s.handleUploadChunk)
	mux.HandleFunc("/api/v1/connect/upload/complete", s.handleCompleteUpload)

	// Add a simple ping endpoint for testing
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Try to use embedded web files first
	embeddedFS := assets.GetWebFileSystem()

	// Check if we have embedded web files
	hasEmbedded := assets.HasEmbeddedWebFiles()

	// Always log whether we have embedded files or not
	log.Printf("Has embedded web files: %v", hasEmbedded)

	if hasEmbedded {
		// Create a file server handler for the embedded web files
		fs := http.FileServer(embeddedFS)

		// Register the file server handler for the root path
		mux.Handle("/", fs)

		if !s.config.ServerQuietOutput {
			log.Printf("Serving web client from embedded files")
		}
	} else {
		// Fall back to the web directory if embedded files are not available
		// Check if the web directory exists
		if _, err := os.Stat("web/static"); !os.IsNotExist(err) {
			// Create a file server handler for the web/static directory
			fs := http.FileServer(http.Dir("web/static"))

			// Register the file server handler for the root path
			mux.Handle("/", fs)

			if !s.config.ServerQuietOutput {
				log.Printf("Serving web client from web/static directory")
			}
		} else {
			if !s.config.ServerQuietOutput {
				log.Printf("Web client files not found, web interface will not be available")
			}
		}
	}

	// Check if the port is available
	if !utils.IsPortAvailable(s.config.ServerPort) {
		// Try to find an available port
		newPort, err := utils.FindAvailablePort(s.config.ServerPort, 100)
		if err != nil {
			return fmt.Errorf("port %d is already in use and no alternative ports are available: %w", s.config.ServerPort, err)
		}

		// Log the port change
		if !s.config.ServerQuietOutput {
			log.Printf("Port %d is already in use. Using port %d instead.", s.config.ServerPort, newPort)
			log.Printf("This could be due to another Lumo server instance or a Lumo connect session using this port.")
			log.Printf("To avoid this in the future, configure a different port with: lumo config:server port <port>")
			log.Printf("%s", utils.GetPortRangeMessage("server"))
		}

		// Update the port
		s.config.ServerPort = newPort
	}

	// Create the server
	s.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.config.ServerPort),
		Handler: handler,
	}

	// If running in daemon mode, start the server in the main goroutine
	if s.isDaemon {
		if !s.config.ServerQuietOutput {
			log.Printf("Starting Lumo REST server in daemon mode on port %d", s.config.ServerPort)
		}
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if !s.config.ServerQuietOutput {
				log.Printf("Error starting server: %v", err)
				if os.IsPermission(err) {
					log.Printf("This may be due to insufficient permissions to bind to port %d.", s.config.ServerPort)
					log.Printf("Try using a port number above 1024 with: lumo config:server port <port>")
				}
			}
			return err
		}
		return nil
	}

	// Otherwise, start the server in a goroutine
	go func() {
		if !s.config.ServerQuietOutput {
			log.Printf("Starting Lumo REST server on port %d", s.config.ServerPort)
		}
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if !s.config.ServerQuietOutput {
				log.Printf("Error starting server: %v", err)
				if os.IsPermission(err) {
					log.Printf("This may be due to insufficient permissions to bind to port %d.", s.config.ServerPort)
					log.Printf("Try using a port number above 1024 with: lumo config:server port <port>")
				}
			}
		}
	}()

	// Add a small delay to allow the server to start
	time.Sleep(100 * time.Millisecond)

	// Test if the server is running
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/ping", s.config.ServerPort))
	if err != nil {
		if !s.config.ServerQuietOutput {
			log.Printf("Warning: Server may not be running correctly: %v", err)
		}
	} else {
		if !s.config.ServerQuietOutput {
			log.Printf("Server is running and responding to requests")
		}
		resp.Body.Close()
	}

	return nil
}

// Stop stops the REST server
func (s *Server) Stop() error {
	if s.server != nil {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown the server
		return s.server.Shutdown(ctx)
	}
	return nil
}

// GetConfig returns the server's configuration
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// handleExecute handles the /api/v1/execute endpoint
func (s *Server) handleExecute(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.Command == "" {
		http.Error(w, "Command is required", http.StatusBadRequest)
		return
	}

	// Create a command based on the request
	var cmd *nlp.Command

	// If type is specified, use it
	if req.Type != "" {
		cmdType := mapStringToCommandType(req.Type)
		cmd = &nlp.Command{
			Type:       cmdType,
			Intent:     req.Command,
			Parameters: req.Params,
			RawInput:   req.Command,
		}
	} else {
		// Otherwise, parse the command
		parser := nlp.NewParser(s.config)
		var err error
		cmd, err = parser.Parse(req.Command)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing command: %v", err), http.StatusBadRequest)
			return
		}
	}

	// Execute the command
	result, err := s.executor.Execute(cmd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing command: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the response
	resp := CommandResponse{
		Success:    !result.IsError,
		Output:     result.Output,
		CommandRun: result.CommandRun,
	}

	if result.IsError {
		resp.Error = result.Output
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Set the status code
	if result.IsError {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Write the response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleStatus handles the /api/v1/status endpoint
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create the response
	resp := StatusResponse{
		Status:  "running",
		Version: version.GetShortVersion(), // Dynamically fetch from version package
		Uptime:  "N/A",                     // This could be calculated if we track server start time
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// mapStringToCommandType maps a string to a CommandType
func mapStringToCommandType(cmdType string) nlp.CommandType {
	switch cmdType {
	case "shell":
		return nlp.CommandTypeShell
	case "ai":
		return nlp.CommandTypeAI
	case "agent":
		return nlp.CommandTypeAgent
	case "chat":
		return nlp.CommandTypeChat
	case "system_health":
		return nlp.CommandTypeSystemHealth
	case "system_report":
		return nlp.CommandTypeSystemReport
	case "help":
		return nlp.CommandTypeHelp
	case "config":
		return nlp.CommandTypeConfig
	case "speed_test":
		return nlp.CommandTypeSpeedTest
	case "magic":
		return nlp.CommandTypeMagic
	case "clipboard":
		return nlp.CommandTypeClipboard
	case "connect":
		return nlp.CommandTypeConnect
	case "create":
		return nlp.CommandTypeCreate
	case "desktop":
		return nlp.CommandTypeDesktop
	default:
		return nlp.CommandTypeAI
	}
}
