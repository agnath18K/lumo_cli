package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
	"github.com/agnath18K/lumo/pkg/nlp"
)

// Server represents the REST API server for Lumo
type Server struct {
	config   *config.Config
	executor *executor.Executor
	server   *http.Server
	isDaemon bool
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

// New creates a new REST server instance
func New(cfg *config.Config, exec *executor.Executor) *Server {
	return &Server{
		config:   cfg,
		executor: exec,
		isDaemon: false,
	}
}

// NewDaemon creates a new REST server instance in daemon mode
func NewDaemon(cfg *config.Config, exec *executor.Executor) *Server {
	return &Server{
		config:   cfg,
		executor: exec,
		isDaemon: true,
	}
}

// Start starts the REST server
func (s *Server) Start() error {
	// Create a new router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/v1/execute", s.handleExecute)
	mux.HandleFunc("/api/v1/status", s.handleStatus)

	// Add a simple ping endpoint for testing
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Create the server
	s.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.config.ServerPort),
		Handler: mux,
	}

	// If running in daemon mode, start the server in the main goroutine
	if s.isDaemon {
		if !s.config.ServerQuietOutput {
			log.Printf("Starting Lumo REST server in daemon mode on port %d", s.config.ServerPort)
		}
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if !s.config.ServerQuietOutput {
				log.Printf("Error starting server: %v", err)
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
			}
		}
	}()

	// Add a small delay to allow the server to start
	time.Sleep(100 * time.Millisecond)

	// Test if the server is running
	_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/ping", s.config.ServerPort))
	if err != nil {
		if !s.config.ServerQuietOutput {
			log.Printf("Warning: Server may not be running correctly: %v", err)
		}
	} else {
		if !s.config.ServerQuietOutput {
			log.Printf("Server is running and responding to requests")
		}
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
		Version: "1.0.1", // This should be dynamically fetched from version package
		Uptime:  "N/A",   // This could be calculated if we track server start time
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
