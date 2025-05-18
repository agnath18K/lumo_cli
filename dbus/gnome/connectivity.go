package gnome

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/agnath18K/lumo/internal/core"
)

// GNOME network-related DBus service names and interfaces
const (
	// NetworkManager is the NetworkManager service
	NetworkManager = "org.freedesktop.NetworkManager"
	// NetworkManagerPath is the NetworkManager object path
	NetworkManagerPath = "/org/freedesktop/NetworkManager"
	// NetworkManagerInterface is the NetworkManager interface
	NetworkManagerInterface = "org.freedesktop.NetworkManager"

	// Bluetooth is the Bluetooth service
	Bluetooth = "org.bluez"
	// BluetoothPath is the Bluetooth object path
	BluetoothPath = "/org/bluez"
	// BluetoothInterface is the Bluetooth interface
	BluetoothInterface = "org.bluez.Manager1"

	// RfKill is the RfKill service (for airplane mode)
	RfKill = "org.gnome.SettingsDaemon.Rfkill"
	// RfKillPath is the RfKill object path
	RfKillPath = "/org/gnome/SettingsDaemon/Rfkill"
	// RfKillInterface is the RfKill interface
	RfKillInterface = "org.gnome.SettingsDaemon.Rfkill"
)

// executeConnectivityCommand executes a connectivity management command
func (e *Environment) executeConnectivityCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Action {
	case "list-devices":
		devices, err := e.GetNetworkDevices(ctx)
		if err != nil {
			return nil, err
		}
		var output strings.Builder
		output.WriteString("Network devices:\n")
		for _, device := range devices {
			status := "Disabled"
			if device.Enabled {
				if device.Connected {
					status = "Connected"
				} else {
					status = "Enabled"
				}
			}
			output.WriteString(fmt.Sprintf("- %s (%s): %s\n", device.Name, device.Type, status))
		}
		return &core.Result{
			Output:  output.String(),
			Success: true,
			Data: map[string]interface{}{
				"devices": devices,
			},
		}, nil
	case "enable-wifi":
		if err := e.EnableWifi(ctx); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  "WiFi enabled",
			Success: true,
		}, nil
	case "disable-wifi":
		if err := e.DisableWifi(ctx); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  "WiFi disabled",
			Success: true,
		}, nil
	case "wifi-status":
		enabled, err := e.GetWifiStatus(ctx)
		if err != nil {
			return nil, err
		}
		status := "disabled"
		if enabled {
			status = "enabled"
		}
		return &core.Result{
			Output:  fmt.Sprintf("WiFi is %s", status),
			Success: true,
			Data: map[string]interface{}{
				"enabled": enabled,
			},
		}, nil
	case "enable-bluetooth":
		if err := e.EnableBluetooth(ctx); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  "Bluetooth enabled",
			Success: true,
		}, nil
	case "disable-bluetooth":
		if err := e.DisableBluetooth(ctx); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  "Bluetooth disabled",
			Success: true,
		}, nil
	case "bluetooth-status":
		enabled, err := e.GetBluetoothStatus(ctx)
		if err != nil {
			return nil, err
		}
		status := "disabled"
		if enabled {
			status = "enabled"
		}
		return &core.Result{
			Output:  fmt.Sprintf("Bluetooth is %s", status),
			Success: true,
			Data: map[string]interface{}{
				"enabled": enabled,
			},
		}, nil
	case "enable-airplane-mode":
		if err := e.SetAirplaneMode(ctx, true); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  "Airplane mode enabled",
			Success: true,
		}, nil
	case "disable-airplane-mode":
		if err := e.SetAirplaneMode(ctx, false); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  "Airplane mode disabled",
			Success: true,
		}, nil
	case "airplane-mode-status":
		enabled, err := e.GetAirplaneMode(ctx)
		if err != nil {
			return nil, err
		}
		status := "disabled"
		if enabled {
			status = "enabled"
		}
		return &core.Result{
			Output:  fmt.Sprintf("Airplane mode is %s", status),
			Success: true,
			Data: map[string]interface{}{
				"enabled": enabled,
			},
		}, nil
	case "enable-hotspot":
		ssid := cmd.Target
		password := ""
		if passwordVal, ok := cmd.Arguments["password"]; ok {
			if passwordStr, ok := passwordVal.(string); ok {
				password = passwordStr
			}
		}
		if ssid == "" {
			return nil, fmt.Errorf("SSID is required")
		}
		if err := e.EnableHotspot(ctx, ssid, password); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("WiFi hotspot enabled with SSID: %s", ssid),
			Success: true,
		}, nil
	case "disable-hotspot":
		if err := e.DisableHotspot(ctx); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  "WiFi hotspot disabled",
			Success: true,
		}, nil
	case "hotspot-status":
		enabled, info, err := e.GetHotspotStatus(ctx)
		if err != nil {
			return nil, err
		}
		status := "disabled"
		if enabled {
			status = "enabled"
		}
		return &core.Result{
			Output:  fmt.Sprintf("WiFi hotspot is %s", status),
			Success: true,
			Data: map[string]interface{}{
				"enabled": enabled,
				"info":    info,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported connectivity action: %s", cmd.Action)
	}
}

// GetNetworkDevices gets a list of available network devices
func (e *Environment) GetNetworkDevices(ctx context.Context) ([]core.NetworkDevice, error) {
	var devices []core.NetworkDevice

	// Get WiFi devices
	wifiEnabled, _ := e.GetWifiStatus(ctx)
	devices = append(devices, core.NetworkDevice{
		ID:        "wifi0",
		Name:      "WiFi",
		Type:      core.NetworkDeviceTypeWifi,
		Enabled:   wifiEnabled,
		Connected: wifiEnabled, // Simplified for now
	})

	// Get Bluetooth devices
	bluetoothEnabled, _ := e.GetBluetoothStatus(ctx)
	devices = append(devices, core.NetworkDevice{
		ID:        "bluetooth0",
		Name:      "Bluetooth",
		Type:      core.NetworkDeviceTypeBluetooth,
		Enabled:   bluetoothEnabled,
		Connected: bluetoothEnabled, // Simplified for now
	})

	// Get Ethernet devices (simplified)
	devices = append(devices, core.NetworkDevice{
		ID:        "ethernet0",
		Name:      "Ethernet",
		Type:      core.NetworkDeviceTypeEthernet,
		Enabled:   true,
		Connected: true, // Simplified for now
	})

	return devices, nil
}

// EnableWifi enables WiFi
func (e *Environment) EnableWifi(ctx context.Context) error {
	// Try using rfkill
	_, err := exec.LookPath("rfkill")
	if err == nil {
		cmd := "rfkill unblock wifi"
		_, err := e.runCommand(cmd)
		if err == nil {
			return nil
		}
	}

	// Try using NetworkManager via DBus
	_, err = e.systemHandler.Call(
		NetworkManager,
		NetworkManagerPath,
		NetworkManagerInterface,
		"SetWirelessEnabled",
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to enable WiFi: %w", err)
	}

	return nil
}

// DisableWifi disables WiFi
func (e *Environment) DisableWifi(ctx context.Context) error {
	// Try using rfkill
	_, err := exec.LookPath("rfkill")
	if err == nil {
		cmd := "rfkill block wifi"
		_, err := e.runCommand(cmd)
		if err == nil {
			return nil
		}
	}

	// Try using NetworkManager via DBus
	_, err = e.systemHandler.Call(
		NetworkManager,
		NetworkManagerPath,
		NetworkManagerInterface,
		"SetWirelessEnabled",
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to disable WiFi: %w", err)
	}

	return nil
}

// GetWifiStatus gets the current WiFi status
func (e *Environment) GetWifiStatus(ctx context.Context) (bool, error) {
	// Try using NetworkManager via DBus
	result, err := e.systemHandler.GetProperty(
		NetworkManager,
		NetworkManagerPath,
		NetworkManagerInterface,
		"WirelessEnabled",
	)
	if err != nil {
		// Fallback to using nmcli
		cmd := "nmcli radio wifi"
		output, err := e.runCommand(cmd)
		if err != nil {
			return false, fmt.Errorf("failed to get WiFi status: %w", err)
		}
		return strings.TrimSpace(output) == "enabled", nil
	}

	enabled, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected result type for WiFi status")
	}

	return enabled, nil
}

// EnableBluetooth enables Bluetooth
func (e *Environment) EnableBluetooth(ctx context.Context) error {
	// Try using rfkill
	_, err := exec.LookPath("rfkill")
	if err == nil {
		cmd := "rfkill unblock bluetooth"
		_, err := e.runCommand(cmd)
		if err == nil {
			return nil
		}
	}

	// Try using bluetoothctl
	_, err = exec.LookPath("bluetoothctl")
	if err == nil {
		cmd := "echo 'power on' | bluetoothctl"
		_, err := e.runCommand(cmd)
		if err == nil {
			return nil
		}
	}

	// Try using DBus
	_, err = e.systemHandler.Call(
		Bluetooth,
		BluetoothPath,
		BluetoothInterface,
		"SetPowered",
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to enable Bluetooth: %w", err)
	}

	return nil
}

// DisableBluetooth disables Bluetooth
func (e *Environment) DisableBluetooth(ctx context.Context) error {
	// Try using rfkill
	_, err := exec.LookPath("rfkill")
	if err == nil {
		cmd := "rfkill block bluetooth"
		_, err := e.runCommand(cmd)
		if err == nil {
			return nil
		}
	}

	// Try using bluetoothctl
	_, err = exec.LookPath("bluetoothctl")
	if err == nil {
		cmd := "echo 'power off' | bluetoothctl"
		_, err := e.runCommand(cmd)
		if err == nil {
			return nil
		}
	}

	// Try using DBus
	_, err = e.systemHandler.Call(
		Bluetooth,
		BluetoothPath,
		BluetoothInterface,
		"SetPowered",
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to disable Bluetooth: %w", err)
	}

	return nil
}

// GetBluetoothStatus gets the current Bluetooth status
func (e *Environment) GetBluetoothStatus(ctx context.Context) (bool, error) {
	// Try using bluetoothctl
	_, err := exec.LookPath("bluetoothctl")
	if err == nil {
		cmd := "bluetoothctl show | grep 'Powered:' | awk '{print $2}'"
		output, err := e.runCommand(cmd)
		if err == nil {
			return strings.TrimSpace(output) == "yes", nil
		}
	}

	// Try using DBus
	result, err := e.systemHandler.GetProperty(
		Bluetooth,
		BluetoothPath,
		BluetoothInterface,
		"Powered",
	)
	if err != nil {
		// If both methods fail, assume Bluetooth is disabled
		return false, nil
	}

	enabled, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected result type for Bluetooth status")
	}

	return enabled, nil
}

// SetAirplaneMode sets the airplane mode state
func (e *Environment) SetAirplaneMode(ctx context.Context, enabled bool) error {
	// Try using rfkill
	_, err := exec.LookPath("rfkill")
	if err == nil {
		var cmd string
		if enabled {
			cmd = "rfkill block all"
		} else {
			cmd = "rfkill unblock all"
		}
		_, err := e.runCommand(cmd)
		if err == nil {
			return nil
		}
	}

	// Try using GNOME settings daemon via DBus
	_, err = e.sessionHandler.Call(
		RfKill,
		RfKillPath,
		RfKillInterface,
		"SetAirplaneMode",
		enabled,
	)
	if err != nil {
		return fmt.Errorf("failed to set airplane mode: %w", err)
	}

	return nil
}

// GetAirplaneMode gets the current airplane mode state
func (e *Environment) GetAirplaneMode(ctx context.Context) (bool, error) {
	// Try using GNOME settings daemon via DBus
	result, err := e.sessionHandler.GetProperty(
		RfKill,
		RfKillPath,
		RfKillInterface,
		"AirplaneMode",
	)
	if err != nil {
		// Fallback to checking if both WiFi and Bluetooth are disabled
		wifiEnabled, _ := e.GetWifiStatus(ctx)
		bluetoothEnabled, _ := e.GetBluetoothStatus(ctx)
		return !wifiEnabled && !bluetoothEnabled, nil
	}

	enabled, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected result type for airplane mode status")
	}

	return enabled, nil
}

// EnableHotspot enables WiFi hotspot
func (e *Environment) EnableHotspot(ctx context.Context, ssid, password string) error {
	// Try using nmcli
	_, err := exec.LookPath("nmcli")
	if err == nil {
		var cmd string
		if password != "" {
			cmd = fmt.Sprintf("nmcli device wifi hotspot ssid '%s' password '%s'", ssid, password)
		} else {
			cmd = fmt.Sprintf("nmcli device wifi hotspot ssid '%s'", ssid)
		}
		_, err := e.runCommand(cmd)
		if err != nil {
			return fmt.Errorf("failed to enable WiFi hotspot: %w", err)
		}
		return nil
	}

	// If nmcli is not available, return an error
	return fmt.Errorf("failed to enable WiFi hotspot: nmcli not found")
}

// DisableHotspot disables WiFi hotspot
func (e *Environment) DisableHotspot(ctx context.Context) error {
	// Try using nmcli
	_, err := exec.LookPath("nmcli")
	if err == nil {
		// Get the hotspot connection name
		cmd := "nmcli -t -f NAME,TYPE connection show | grep ':wifi' | grep -i 'hotspot' | cut -d: -f1"
		output, err := e.runCommand(cmd)
		if err != nil {
			return fmt.Errorf("failed to find WiFi hotspot connection: %w", err)
		}

		// If a hotspot connection is found, delete it
		if output != "" {
			cmd = fmt.Sprintf("nmcli connection down '%s'", strings.TrimSpace(output))
			_, err = e.runCommand(cmd)
			if err != nil {
				return fmt.Errorf("failed to disable WiFi hotspot: %w", err)
			}
			return nil
		}
	}

	// If nmcli is not available or no hotspot connection is found, return an error
	return fmt.Errorf("failed to disable WiFi hotspot: no active hotspot found")
}

// GetHotspotStatus gets the current WiFi hotspot status
func (e *Environment) GetHotspotStatus(ctx context.Context) (bool, map[string]interface{}, error) {
	// Try using nmcli
	_, err := exec.LookPath("nmcli")
	if err == nil {
		// Check if there's an active hotspot connection
		cmd := "nmcli -t -f NAME,TYPE,ACTIVE connection show | grep ':wifi:yes' | grep -i 'hotspot' | cut -d: -f1"
		output, err := e.runCommand(cmd)
		if err != nil {
			return false, nil, fmt.Errorf("failed to check WiFi hotspot status: %w", err)
		}

		// If a hotspot connection is found, get its details
		if output != "" {
			hotspotName := strings.TrimSpace(output)
			cmd = fmt.Sprintf("nmcli -t -f 802-11-wireless.ssid connection show '%s'", hotspotName)
			ssidOutput, _ := e.runCommand(cmd)
			ssid := strings.TrimSpace(strings.TrimPrefix(ssidOutput, "802-11-wireless.ssid:"))

			return true, map[string]interface{}{
				"name": hotspotName,
				"ssid": ssid,
			}, nil
		}

		return false, nil, nil
	}

	// If nmcli is not available, return an error
	return false, nil, fmt.Errorf("failed to check WiFi hotspot status: nmcli not found")
}

// EnableNetworkDevice enables a network device
func (e *Environment) EnableNetworkDevice(ctx context.Context, deviceID string) error {
	// Get the device type
	devices, err := e.GetNetworkDevices(ctx)
	if err != nil {
		return fmt.Errorf("failed to get network devices: %w", err)
	}

	// Find the device
	var device *core.NetworkDevice
	for i, d := range devices {
		if d.ID == deviceID {
			device = &devices[i]
			break
		}
	}

	if device == nil {
		return fmt.Errorf("network device not found: %s", deviceID)
	}

	// Enable the device based on its type
	switch device.Type {
	case core.NetworkDeviceTypeWifi:
		return e.EnableWifi(ctx)
	case core.NetworkDeviceTypeBluetooth:
		return e.EnableBluetooth(ctx)
	case core.NetworkDeviceTypeEthernet:
		// Ethernet is typically always enabled
		return nil
	case core.NetworkDeviceTypeHotspot:
		// Hotspot requires SSID and password
		return fmt.Errorf("enabling hotspot requires SSID and password")
	default:
		return fmt.Errorf("unsupported network device type: %s", device.Type)
	}
}

// DisableNetworkDevice disables a network device
func (e *Environment) DisableNetworkDevice(ctx context.Context, deviceID string) error {
	// Get the device type
	devices, err := e.GetNetworkDevices(ctx)
	if err != nil {
		return fmt.Errorf("failed to get network devices: %w", err)
	}

	// Find the device
	var device *core.NetworkDevice
	for i, d := range devices {
		if d.ID == deviceID {
			device = &devices[i]
			break
		}
	}

	if device == nil {
		return fmt.Errorf("network device not found: %s", deviceID)
	}

	// Disable the device based on its type
	switch device.Type {
	case core.NetworkDeviceTypeWifi:
		return e.DisableWifi(ctx)
	case core.NetworkDeviceTypeBluetooth:
		return e.DisableBluetooth(ctx)
	case core.NetworkDeviceTypeEthernet:
		// Disabling Ethernet is not typically supported
		return fmt.Errorf("disabling Ethernet is not supported")
	case core.NetworkDeviceTypeHotspot:
		return e.DisableHotspot(ctx)
	default:
		return fmt.Errorf("unsupported network device type: %s", device.Type)
	}
}

// ConnectNetworkDevice connects to a network device
func (e *Environment) ConnectNetworkDevice(ctx context.Context, deviceID string, params map[string]interface{}) error {
	// This is a simplified implementation
	// In a real implementation, you would need to handle different types of connections
	// For now, we'll just enable the device
	return e.EnableNetworkDevice(ctx, deviceID)
}

// DisconnectNetworkDevice disconnects from a network device
func (e *Environment) DisconnectNetworkDevice(ctx context.Context, deviceID string) error {
	// This is a simplified implementation
	// In a real implementation, you would need to handle different types of connections
	// For now, we'll just disable the device
	return e.DisableNetworkDevice(ctx, deviceID)
}
