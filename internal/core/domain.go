package core

// CommandType represents the type of desktop command
type CommandType string

const (
	// CommandTypeWindow represents window management commands
	CommandTypeWindow CommandType = "window"
	// CommandTypeApplication represents application management commands
	CommandTypeApplication CommandType = "application"
	// CommandTypeSystem represents system-level commands
	CommandTypeSystem CommandType = "system"
	// CommandTypeNotification represents notification commands
	CommandTypeNotification CommandType = "notification"
	// CommandTypeMedia represents media control commands
	CommandTypeMedia CommandType = "media"
	// CommandTypeAppearance represents appearance management commands
	CommandTypeAppearance CommandType = "appearance"
	// CommandTypeSound represents sound settings commands
	CommandTypeSound CommandType = "sound"
	// CommandTypeConnectivity represents network connectivity commands
	CommandTypeConnectivity CommandType = "connectivity"
)

// Command represents a desktop command to be executed
type Command struct {
	// Type is the type of command
	Type CommandType
	// Action is the specific action to perform
	Action string
	// Target is the target of the action (e.g., window name, application name)
	Target string
	// Arguments are additional parameters for the command
	Arguments map[string]interface{}
	// RawInput is the original natural language input
	RawInput string
}

// Result represents the result of executing a desktop command
type Result struct {
	// Output is the textual output of the command
	Output string
	// Success indicates whether the command was successful
	Success bool
	// Error is the error message if the command failed
	Error string
	// Data contains any structured data returned by the command
	Data map[string]interface{}
}

// Capability represents a capability of a desktop environment
type Capability string

const (
	// CapabilityWindowManagement represents window management capabilities
	CapabilityWindowManagement Capability = "window_management"
	// CapabilityApplicationLaunch represents application launching capabilities
	CapabilityApplicationLaunch Capability = "application_launch"
	// CapabilityNotifications represents notification capabilities
	CapabilityNotifications Capability = "notifications"
	// CapabilityMediaControl represents media control capabilities
	CapabilityMediaControl Capability = "media_control"
	// CapabilityScreenshot represents screenshot capabilities
	CapabilityScreenshot Capability = "screenshot"
	// CapabilityClipboard represents clipboard management capabilities
	CapabilityClipboard Capability = "clipboard"
	// CapabilityAppearanceManagement represents appearance management capabilities
	CapabilityAppearanceManagement Capability = "appearance_management"
	// CapabilitySoundManagement represents sound settings management capabilities
	CapabilitySoundManagement Capability = "sound_management"
	// CapabilityConnectivityManagement represents network connectivity management capabilities
	CapabilityConnectivityManagement Capability = "connectivity_management"
)

// Window represents a desktop window
type Window struct {
	// ID is the unique identifier for the window
	ID string
	// Title is the window title
	Title string
	// Application is the application that owns the window
	Application string
	// Geometry contains the window's position and size
	Geometry WindowGeometry
	// State contains the window's state (maximized, minimized, etc.)
	State WindowState
}

// WindowGeometry represents the position and size of a window
type WindowGeometry struct {
	X      int
	Y      int
	Width  int
	Height int
}

// WindowState represents the state of a window
type WindowState struct {
	Maximized  bool
	Minimized  bool
	Fullscreen bool
	Active     bool
}

// Application represents a desktop application
type Application struct {
	// ID is the unique identifier for the application
	ID string
	// Name is the application name
	Name string
	// Executable is the path to the application executable
	Executable string
	// DesktopFile is the path to the application's desktop file
	DesktopFile string
	// Running indicates whether the application is currently running
	Running bool
}

// Notification represents a desktop notification
type Notification struct {
	// ID is the unique identifier for the notification
	ID uint32
	// Summary is the notification summary
	Summary string
	// Body is the notification body
	Body string
	// Icon is the notification icon
	Icon string
	// Actions are the available actions for the notification
	Actions []string
	// Hints are additional hints for the notification
	Hints map[string]interface{}
	// Timeout is the notification timeout in milliseconds
	Timeout int32
}

// SoundDevice represents a sound device (input or output)
type SoundDevice struct {
	// ID is the unique identifier for the sound device
	ID string
	// Name is the human-readable name of the device
	Name string
	// Description is a description of the device
	Description string
	// IsInput indicates whether this is an input device (microphone)
	IsInput bool
	// IsDefault indicates whether this is the default device
	IsDefault bool
	// Volume is the current volume level (0-100)
	Volume int
	// Muted indicates whether the device is muted
	Muted bool
}

// NetworkDeviceType represents the type of network device
type NetworkDeviceType string

const (
	// NetworkDeviceTypeWifi represents a WiFi device
	NetworkDeviceTypeWifi NetworkDeviceType = "wifi"
	// NetworkDeviceTypeBluetooth represents a Bluetooth device
	NetworkDeviceTypeBluetooth NetworkDeviceType = "bluetooth"
	// NetworkDeviceTypeEthernet represents an Ethernet device
	NetworkDeviceTypeEthernet NetworkDeviceType = "ethernet"
	// NetworkDeviceTypeHotspot represents a WiFi hotspot
	NetworkDeviceTypeHotspot NetworkDeviceType = "hotspot"
)

// NetworkDevice represents a network device (WiFi, Bluetooth, Ethernet, etc.)
type NetworkDevice struct {
	// ID is the unique identifier for the network device
	ID string
	// Name is the human-readable name of the device
	Name string
	// Type is the type of network device
	Type NetworkDeviceType
	// Enabled indicates whether the device is enabled
	Enabled bool
	// Connected indicates whether the device is connected
	Connected bool
	// Address is the device address (MAC address, IP address, etc.)
	Address string
	// Properties contains additional device-specific properties
	Properties map[string]interface{}
}
