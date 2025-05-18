package common

import (
	"fmt"

	"github.com/agnath18K/lumo/internal/core"
	"github.com/godbus/dbus/v5"
)

// DBusConnection represents a connection to the DBus
type DBusConnection interface {
	// Object returns a DBus object
	Object(dest string, path dbus.ObjectPath) dbus.BusObject

	// Signal returns a channel for receiving signals
	Signal() chan *dbus.Signal

	// AddMatchSignal adds a match rule for a signal
	AddMatchSignal(options ...dbus.MatchOption) error

	// RemoveMatchSignal removes a match rule for a signal
	RemoveMatchSignal(options ...dbus.MatchOption) error

	// Close closes the connection
	Close() error

	// GetConn returns the underlying dbus.Conn
	GetConn() *dbus.Conn
}

// DBusHandler implements the core.DBusHandler interface
type DBusHandler struct {
	// conn is the DBus connection
	conn DBusConnection
}

// NewDBusHandler creates a new DBus handler
func NewDBusHandler(conn DBusConnection) *DBusHandler {
	return &DBusHandler{
		conn: conn,
	}
}

// Connect connects to the DBus
func (h *DBusHandler) Connect() error {
	// Connection is already established in the constructor
	return nil
}

// Disconnect disconnects from the DBus
func (h *DBusHandler) Disconnect() error {
	if h.conn != nil {
		return h.conn.Close()
	}
	return nil
}

// Call calls a DBus method
func (h *DBusHandler) Call(service, objectPath, interfaceName, method string, args ...interface{}) ([]interface{}, error) {
	if h.conn == nil {
		return nil, fmt.Errorf("DBus connection is nil")
	}
	obj := h.conn.Object(service, dbus.ObjectPath(objectPath))
	call := obj.Call(interfaceName+"."+method, 0, args...)
	if call.Err != nil {
		return nil, call.Err
	}
	return call.Body, nil
}

// GetProperty gets a DBus property
func (h *DBusHandler) GetProperty(service, objectPath, interfaceName, property string) (interface{}, error) {
	if h.conn == nil {
		return nil, fmt.Errorf("DBus connection is nil")
	}
	obj := h.conn.Object(service, dbus.ObjectPath(objectPath))
	variant, err := obj.GetProperty(interfaceName + "." + property)
	if err != nil {
		return nil, err
	}
	return variant.Value(), nil
}

// SetProperty sets a DBus property
func (h *DBusHandler) SetProperty(service, objectPath, interfaceName, property string, value interface{}) error {
	if h.conn == nil {
		return fmt.Errorf("DBus connection is nil")
	}
	obj := h.conn.Object(service, dbus.ObjectPath(objectPath))
	return obj.SetProperty(interfaceName+"."+property, dbus.MakeVariant(value))
}

// AddMatch adds a match rule
func (h *DBusHandler) AddMatch(rule string) error {
	if h.conn == nil {
		return fmt.Errorf("DBus connection is nil")
	}
	// Use the correct match option
	return h.conn.AddMatchSignal(dbus.WithMatchOption("type", "signal"), dbus.WithMatchOption("match", rule))
}

// RemoveMatch removes a match rule
func (h *DBusHandler) RemoveMatch(rule string) error {
	if h.conn == nil {
		return fmt.Errorf("DBus connection is nil")
	}
	// Use the correct match option
	return h.conn.RemoveMatchSignal(dbus.WithMatchOption("type", "signal"), dbus.WithMatchOption("match", rule))
}

// Signal returns a channel for receiving signals
func (h *DBusHandler) Signal() <-chan *core.DBusSignal {
	// Create a channel for core.DBusSignal
	ch := make(chan *core.DBusSignal)

	if h.conn == nil {
		// Close the channel immediately if connection is nil
		close(ch)
		return ch
	}

	// Get the dbus.Signal channel
	dbusSignals := h.conn.Signal()

	// Convert dbus.Signal to core.DBusSignal
	go func() {
		for signal := range dbusSignals {
			ch <- &core.DBusSignal{
				Path: string(signal.Path),
				Name: signal.Name,
				Body: signal.Body,
			}
		}
		close(ch)
	}()

	return ch
}
