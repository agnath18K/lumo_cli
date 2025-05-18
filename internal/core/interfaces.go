package core

import "context"

// DesktopEnvironment represents a desktop environment
type DesktopEnvironment interface {
	// Name returns the name of the desktop environment
	Name() string

	// IsAvailable checks if this desktop environment is available on the system
	IsAvailable() bool

	// ExecuteCommand executes a desktop command
	ExecuteCommand(ctx context.Context, cmd *Command) (*Result, error)

	// GetCapabilities returns the capabilities of this desktop environment
	GetCapabilities() []Capability

	// GetWindows returns a list of all windows
	GetWindows(ctx context.Context) ([]Window, error)

	// GetActiveWindow returns the currently active window
	GetActiveWindow(ctx context.Context) (*Window, error)

	// GetRunningApplications returns a list of running applications
	GetRunningApplications(ctx context.Context) ([]Application, error)

	// LaunchApplication launches an application
	LaunchApplication(ctx context.Context, appName string, args ...string) error

	// CloseWindow closes a window
	CloseWindow(ctx context.Context, windowID string) error

	// MinimizeWindow minimizes a window
	MinimizeWindow(ctx context.Context, windowID string) error

	// MaximizeWindow maximizes a window
	MaximizeWindow(ctx context.Context, windowID string) error

	// RestoreWindow restores a window
	RestoreWindow(ctx context.Context, windowID string) error

	// MoveWindow moves a window to a new position
	MoveWindow(ctx context.Context, windowID string, x, y int) error

	// ResizeWindow resizes a window
	ResizeWindow(ctx context.Context, windowID string, width, height int) error

	// FocusWindow focuses a window
	FocusWindow(ctx context.Context, windowID string) error

	// ShowDesktop shows the desktop
	ShowDesktop(ctx context.Context) error

	// SendNotification sends a notification
	SendNotification(ctx context.Context, summary, body, icon string) (uint32, error)

	// CloseNotification closes a notification
	CloseNotification(ctx context.Context, id uint32) error

	// TakeScreenshot takes a screenshot
	TakeScreenshot(ctx context.Context, fullScreen bool, delay int) (string, error)

	// GetClipboardText gets the text from the clipboard
	GetClipboardText(ctx context.Context) (string, error)

	// SetClipboardText sets the text in the clipboard
	SetClipboardText(ctx context.Context, text string) error

	// SetGtkTheme sets the GTK theme (light/dark mode)
	SetGtkTheme(ctx context.Context, theme string) error

	// SetDesktopBackground sets the desktop background image
	SetDesktopBackground(ctx context.Context, imagePath string) error

	// SetAccentColor sets the accent color if supported
	SetAccentColor(ctx context.Context, color string) error

	// SetIconTheme sets the icon theme
	SetIconTheme(ctx context.Context, theme string) error

	// GetCurrentTheme gets the current GTK theme
	GetCurrentTheme(ctx context.Context) (string, error)

	// GetCurrentBackground gets the current desktop background
	GetCurrentBackground(ctx context.Context) (string, error)

	// GetCurrentIconTheme gets the current icon theme
	GetCurrentIconTheme(ctx context.Context) (string, error)

	// SetVolume sets the system volume level (0-100)
	SetVolume(ctx context.Context, level int) error

	// GetVolume gets the current system volume level (0-100)
	GetVolume(ctx context.Context) (int, error)

	// SetMute sets the system mute state
	SetMute(ctx context.Context, mute bool) error

	// GetMute gets the current system mute state
	GetMute(ctx context.Context) (bool, error)

	// SetInputVolume sets the microphone volume level (0-100)
	SetInputVolume(ctx context.Context, level int) error

	// GetInputVolume gets the current microphone volume level (0-100)
	GetInputVolume(ctx context.Context) (int, error)

	// SetInputMute sets the microphone mute state
	SetInputMute(ctx context.Context, mute bool) error

	// GetInputMute gets the current microphone mute state
	GetInputMute(ctx context.Context) (bool, error)

	// GetSoundDevices gets a list of available sound devices
	GetSoundDevices(ctx context.Context) ([]SoundDevice, error)

	// SetDefaultSoundDevice sets the default sound device
	SetDefaultSoundDevice(ctx context.Context, deviceID string) error

	// GetNetworkDevices gets a list of available network devices
	GetNetworkDevices(ctx context.Context) ([]NetworkDevice, error)

	// EnableNetworkDevice enables a network device
	EnableNetworkDevice(ctx context.Context, deviceID string) error

	// DisableNetworkDevice disables a network device
	DisableNetworkDevice(ctx context.Context, deviceID string) error

	// ConnectNetworkDevice connects to a network device
	ConnectNetworkDevice(ctx context.Context, deviceID string, params map[string]interface{}) error

	// DisconnectNetworkDevice disconnects from a network device
	DisconnectNetworkDevice(ctx context.Context, deviceID string) error

	// SetAirplaneMode sets the airplane mode state
	SetAirplaneMode(ctx context.Context, enabled bool) error

	// GetAirplaneMode gets the current airplane mode state
	GetAirplaneMode(ctx context.Context) (bool, error)

	// EnableWifi enables WiFi
	EnableWifi(ctx context.Context) error

	// DisableWifi disables WiFi
	DisableWifi(ctx context.Context) error

	// GetWifiStatus gets the current WiFi status
	GetWifiStatus(ctx context.Context) (bool, error)

	// EnableBluetooth enables Bluetooth
	EnableBluetooth(ctx context.Context) error

	// DisableBluetooth disables Bluetooth
	DisableBluetooth(ctx context.Context) error

	// GetBluetoothStatus gets the current Bluetooth status
	GetBluetoothStatus(ctx context.Context) (bool, error)

	// EnableHotspot enables WiFi hotspot
	EnableHotspot(ctx context.Context, ssid, password string) error

	// DisableHotspot disables WiFi hotspot
	DisableHotspot(ctx context.Context) error

	// GetHotspotStatus gets the current WiFi hotspot status
	GetHotspotStatus(ctx context.Context) (bool, map[string]interface{}, error)
}

// DesktopFactory creates desktop environment instances
type DesktopFactory interface {
	// DetectEnvironment detects the current desktop environment
	DetectEnvironment() (DesktopEnvironment, error)

	// GetEnvironment gets a specific desktop environment by name
	GetEnvironment(name string) (DesktopEnvironment, error)

	// ListAvailableEnvironments lists all available desktop environments
	ListAvailableEnvironments() []string
}

// DBusHandler handles DBus communication
type DBusHandler interface {
	// Connect connects to the DBus
	Connect() error

	// Disconnect disconnects from the DBus
	Disconnect() error

	// Call calls a DBus method
	Call(service, objectPath, interfaceName, method string, args ...interface{}) ([]interface{}, error)

	// GetProperty gets a DBus property
	GetProperty(service, objectPath, interfaceName, property string) (interface{}, error)

	// SetProperty sets a DBus property
	SetProperty(service, objectPath, interfaceName, property string, value interface{}) error

	// AddMatch adds a match rule
	AddMatch(rule string) error

	// RemoveMatch removes a match rule
	RemoveMatch(rule string) error

	// Signal returns a channel for receiving signals
	Signal() <-chan *DBusSignal
}

// DBusSignal represents a DBus signal
type DBusSignal struct {
	// Path is the object path the signal was emitted from
	Path string
	// Name is the name of the signal
	Name string
	// Body is the body of the signal
	Body []interface{}
}

// Assistant processes natural language commands
type Assistant interface {
	// ProcessCommand processes a natural language command
	ProcessCommand(ctx context.Context, input string) (*Result, error)

	// GetSupportedCommands returns a list of supported commands
	GetSupportedCommands() []string

	// GetCommandExamples returns examples of supported commands
	GetCommandExamples() []string
}
