package daemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
	"github.com/agnath18K/lumo/pkg/server"
)

const (
	// PidFileName is the name of the PID file
	PidFileName = "lumo-server.pid"
	// LogFileName is the name of the log file
	LogFileName = "lumo-server.log"
)

// Daemon represents a background daemon process
type Daemon struct {
	config *config.Config
}

// New creates a new daemon instance
func New(cfg *config.Config) *Daemon {
	return &Daemon{
		config: cfg,
	}
}

// GetPidFilePath returns the path to the PID file
func (d *Daemon) GetPidFilePath() string {
	// Use the user's home directory for the PID file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to /tmp if we can't get the home directory
		return filepath.Join("/tmp", PidFileName)
	}
	return filepath.Join(homeDir, ".lumo", PidFileName)
}

// GetLogFilePath returns the path to the log file
func (d *Daemon) GetLogFilePath() string {
	// Use the user's home directory for the log file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to /tmp if we can't get the home directory
		return filepath.Join("/tmp", LogFileName)
	}
	return filepath.Join(homeDir, ".lumo", LogFileName)
}

// IsRunning checks if the daemon is already running
func (d *Daemon) IsRunning() (bool, int, error) {
	// Check if the PID file exists
	pidFile := d.GetPidFilePath()
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return false, 0, nil
	}

	// Read the PID from the file
	pidBytes, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return false, 0, fmt.Errorf("failed to read PID file: %w", err)
	}

	// Parse the PID
	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil {
		return false, 0, fmt.Errorf("failed to parse PID: %w", err)
	}

	// Check if the process is running
	process, err := os.FindProcess(pid)
	if err != nil {
		// On Unix systems, FindProcess never returns an error
		return false, 0, nil
	}

	// Send a signal 0 to the process to check if it's running
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// Process is not running, clean up the PID file
		os.Remove(pidFile)
		return false, 0, nil
	}

	return true, pid, nil
}

// Start starts the daemon
func (d *Daemon) Start() error {
	// Check if the daemon is already running
	running, pid, err := d.IsRunning()
	if err != nil {
		return fmt.Errorf("failed to check if daemon is running: %w", err)
	}
	if running {
		return fmt.Errorf("daemon is already running with PID %d", pid)
	}

	// Create the .lumo directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err == nil {
		lumoDir := filepath.Join(homeDir, ".lumo")
		os.MkdirAll(lumoDir, 0755)
	}

	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Open the log file
	logFile, err := os.OpenFile(d.GetLogFilePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create a new command to run the server in daemon mode
	cmd := exec.Command(execPath, "server:daemon")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Create a new session
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Write the PID to the PID file
	pidFile := d.GetPidFilePath()
	if err := ioutil.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		logFile.Close()
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	// Close the log file
	logFile.Close()

	log.Printf("Daemon started with PID %d", cmd.Process.Pid)
	return nil
}

// Stop stops the daemon
func (d *Daemon) Stop() error {
	// Check if the daemon is running
	running, pid, err := d.IsRunning()
	if err != nil {
		return fmt.Errorf("failed to check if daemon is running: %w", err)
	}
	if !running {
		return fmt.Errorf("daemon is not running")
	}

	// Find the process
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	// Send a SIGTERM signal to the process
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to process: %w", err)
	}

	// Remove the PID file
	pidFile := d.GetPidFilePath()
	if err := os.Remove(pidFile); err != nil {
		return fmt.Errorf("failed to remove PID file: %w", err)
	}

	log.Printf("Daemon stopped")
	return nil
}

// Status returns the status of the daemon
func (d *Daemon) Status() (bool, int, error) {
	return d.IsRunning()
}

// RunServer runs the server in daemon mode
func (d *Daemon) RunServer(exec *executor.Executor) error {
	// This function is called by the daemon process
	if !d.config.ServerQuietOutput {
		log.Printf("Starting Lumo server in daemon mode on port %d", d.config.ServerPort)
	}

	// Create a new server in daemon mode
	srv := server.NewDaemon(d.config, exec)

	// Start the server (this will block in daemon mode)
	return srv.Start()
}
