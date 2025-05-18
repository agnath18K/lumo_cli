package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/godbus/dbus/v5"
)

// DBusType represents the type of DBus connection
type DBusType int

const (
	// DBusTypeSystem represents the system DBus
	DBusTypeSystem DBusType = iota
	// DBusTypeSession represents the session DBus
	DBusTypeSession
)

// NewDBusConnection creates a new DBus connection
func NewDBusConnection(dbusType DBusType) (DBusConnection, error) {
	var conn *dbus.Conn
	var err error

	// Check if DBUS_SESSION_BUS_ADDRESS environment variable is set
	if dbusType == DBusTypeSession {
		sessionAddr := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
		if sessionAddr == "" {
			fmt.Printf("DEBUG: DBUS_SESSION_BUS_ADDRESS environment variable is not set\n")
		} else {
			fmt.Printf("DEBUG: DBUS_SESSION_BUS_ADDRESS = %s\n", sessionAddr)
		}
	}

	// Check if DISPLAY environment variable is set (needed for X11 applications)
	display := os.Getenv("DISPLAY")
	if display == "" {
		fmt.Printf("DEBUG: DISPLAY environment variable is not set\n")
	} else {
		fmt.Printf("DEBUG: DISPLAY = %s\n", display)
	}

	switch dbusType {
	case DBusTypeSystem:
		fmt.Printf("DEBUG: Connecting to system DBus...\n")
		conn, err = dbus.SystemBus()
	case DBusTypeSession:
		fmt.Printf("DEBUG: Connecting to session DBus...\n")
		conn, err = dbus.SessionBus()
	default:
		return nil, fmt.Errorf("invalid DBus type: %d", dbusType)
	}

	if err != nil {
		fmt.Printf("DEBUG: DBus connection error: %v\n", err)
		return nil, fmt.Errorf("failed to connect to DBus: %w", err)
	}

	fmt.Printf("DEBUG: Successfully connected to DBus\n")
	return &dbusConnection{conn: conn}, nil
}

// dbusConnection implements the DBusConnection interface
type dbusConnection struct {
	conn *dbus.Conn
}

// GetConn returns the underlying dbus.Conn
func (c *dbusConnection) GetConn() *dbus.Conn {
	return c.conn
}

// Object returns a DBus object
func (c *dbusConnection) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
	return c.conn.Object(dest, path)
}

// Signal returns a channel for receiving signals
func (c *dbusConnection) Signal() chan *dbus.Signal {
	ch := make(chan *dbus.Signal)
	c.conn.Signal(ch)
	return ch
}

// AddMatchSignal adds a match rule for a signal
func (c *dbusConnection) AddMatchSignal(options ...dbus.MatchOption) error {
	return c.conn.AddMatchSignal(options...)
}

// RemoveMatchSignal removes a match rule for a signal
func (c *dbusConnection) RemoveMatchSignal(options ...dbus.MatchOption) error {
	return c.conn.RemoveMatchSignal(options...)
}

// Close closes the connection
func (c *dbusConnection) Close() error {
	return c.conn.Close()
}

// DetectDesktopEnvironment detects the current desktop environment
func DetectDesktopEnvironment() string {
	// Check XDG_CURRENT_DESKTOP environment variable
	if desktop := os.Getenv("XDG_CURRENT_DESKTOP"); desktop != "" {
		return strings.ToLower(desktop)
	}

	// Check DESKTOP_SESSION environment variable
	if session := os.Getenv("DESKTOP_SESSION"); session != "" {
		return strings.ToLower(session)
	}

	// Check for common desktop environment processes
	desktops := map[string][]string{
		"gnome":    {"gnome-shell", "gnome-session"},
		"kde":      {"plasmashell", "kwin"},
		"xfce":     {"xfce4-session", "xfwm4"},
		"mate":     {"mate-session", "marco"},
		"lxde":     {"lxsession", "openbox"},
		"cinnamon": {"cinnamon", "cinnamon-session"},
	}

	for desktop, processes := range desktops {
		for _, process := range processes {
			if isProcessRunning(process) {
				return desktop
			}
		}
	}

	// Default to "unknown"
	return "unknown"
}

// isProcessRunning checks if a process is running
func isProcessRunning(processName string) bool {
	cmd := exec.Command("pgrep", processName)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// IsDBusServiceAvailable checks if a DBus service is available
func IsDBusServiceAvailable(conn DBusConnection, service string) bool {
	fmt.Printf("DEBUG: Checking if DBus service %s is available...\n", service)

	obj := conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus")
	call := obj.Call("org.freedesktop.DBus.NameHasOwner", 0, service)
	if call.Err != nil {
		fmt.Printf("DEBUG: Error checking if service %s is available: %v\n", service, call.Err)
		return false
	}

	var hasOwner bool
	if err := call.Store(&hasOwner); err != nil {
		fmt.Printf("DEBUG: Error storing result: %v\n", err)
		return false
	}

	fmt.Printf("DEBUG: Service %s available: %v\n", service, hasOwner)
	return hasOwner
}

// ListDBusServices lists all available DBus services
func ListDBusServices(conn DBusConnection) ([]string, error) {
	obj := conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus")
	call := obj.Call("org.freedesktop.DBus.ListNames", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var services []string
	if err := call.Store(&services); err != nil {
		return nil, err
	}

	return services, nil
}
