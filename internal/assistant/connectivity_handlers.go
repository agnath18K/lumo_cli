package assistant

import (
	"strings"

	"github.com/agnath18K/lumo/internal/core"
)

// handleListNetworkDevices handles the "list network devices" command
func (p *Processor) handleListNetworkDevices(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "list-devices",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleEnableWifi handles the "enable wifi" command
func (p *Processor) handleEnableWifi(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "enable-wifi",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleDisableWifi handles the "disable wifi" command
func (p *Processor) handleDisableWifi(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "disable-wifi",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleWifiStatus handles the "wifi status" command
func (p *Processor) handleWifiStatus(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "wifi-status",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleEnableBluetooth handles the "enable bluetooth" command
func (p *Processor) handleEnableBluetooth(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "enable-bluetooth",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleDisableBluetooth handles the "disable bluetooth" command
func (p *Processor) handleDisableBluetooth(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "disable-bluetooth",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleBluetoothStatus handles the "bluetooth status" command
func (p *Processor) handleBluetoothStatus(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "bluetooth-status",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleEnableAirplaneMode handles the "enable airplane mode" command
func (p *Processor) handleEnableAirplaneMode(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "enable-airplane-mode",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleDisableAirplaneMode handles the "disable airplane mode" command
func (p *Processor) handleDisableAirplaneMode(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "disable-airplane-mode",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleAirplaneModeStatus handles the "airplane mode status" command
func (p *Processor) handleAirplaneModeStatus(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "airplane-mode-status",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleEnableHotspot handles the "enable hotspot" command
func (p *Processor) handleEnableHotspot(input string) (*core.Command, error) {
	// Extract SSID and password from the input
	ssid := ""
	password := ""

	// Look for SSID in the input
	if strings.Contains(input, "ssid") || strings.Contains(input, "name") {
		parts := strings.Split(input, "ssid")
		if len(parts) > 1 {
			ssidPart := parts[1]
			ssidPart = strings.TrimSpace(ssidPart)
			if strings.HasPrefix(ssidPart, ":") || strings.HasPrefix(ssidPart, "=") {
				ssidPart = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(ssidPart, ":"), "="))
			}
			// Extract the SSID (may be in quotes)
			if strings.Contains(ssidPart, "'") {
				ssidParts := strings.Split(ssidPart, "'")
				if len(ssidParts) > 1 {
					ssid = ssidParts[1]
				}
			} else if strings.Contains(ssidPart, "\"") {
				ssidParts := strings.Split(ssidPart, "\"")
				if len(ssidParts) > 1 {
					ssid = ssidParts[1]
				}
			} else {
				// Take the first word as the SSID
				words := strings.Fields(ssidPart)
				if len(words) > 0 {
					ssid = words[0]
				}
			}
		}
	} else if strings.Contains(input, "name") {
		parts := strings.Split(input, "name")
		if len(parts) > 1 {
			namePart := parts[1]
			namePart = strings.TrimSpace(namePart)
			if strings.HasPrefix(namePart, ":") || strings.HasPrefix(namePart, "=") {
				namePart = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(namePart, ":"), "="))
			}
			// Extract the name (may be in quotes)
			if strings.Contains(namePart, "'") {
				nameParts := strings.Split(namePart, "'")
				if len(nameParts) > 1 {
					ssid = nameParts[1]
				}
			} else if strings.Contains(namePart, "\"") {
				nameParts := strings.Split(namePart, "\"")
				if len(nameParts) > 1 {
					ssid = nameParts[1]
				}
			} else {
				// Take the first word as the name
				words := strings.Fields(namePart)
				if len(words) > 0 {
					ssid = words[0]
				}
			}
		}
	}

	// Look for password in the input
	if strings.Contains(input, "password") {
		parts := strings.Split(input, "password")
		if len(parts) > 1 {
			passwordPart := parts[1]
			passwordPart = strings.TrimSpace(passwordPart)
			if strings.HasPrefix(passwordPart, ":") || strings.HasPrefix(passwordPart, "=") {
				passwordPart = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(passwordPart, ":"), "="))
			}
			// Extract the password (may be in quotes)
			if strings.Contains(passwordPart, "'") {
				passwordParts := strings.Split(passwordPart, "'")
				if len(passwordParts) > 1 {
					password = passwordParts[1]
				}
			} else if strings.Contains(passwordPart, "\"") {
				passwordParts := strings.Split(passwordPart, "\"")
				if len(passwordParts) > 1 {
					password = passwordParts[1]
				}
			} else {
				// Take the first word as the password
				words := strings.Fields(passwordPart)
				if len(words) > 0 {
					password = words[0]
				}
			}
		}
	}

	// If no SSID was found, use a default one
	if ssid == "" {
		ssid = "LumoHotspot"
	}

	// Create the command
	cmd := &core.Command{
		Type:      core.CommandTypeConnectivity,
		Action:    "enable-hotspot",
		Target:    ssid,
		Arguments: make(map[string]interface{}),
		RawInput:  input,
	}

	// Add password if provided
	if password != "" {
		cmd.Arguments["password"] = password
	}

	return cmd, nil
}

// handleDisableHotspot handles the "disable hotspot" command
func (p *Processor) handleDisableHotspot(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "disable-hotspot",
		Target:   "",
		RawInput: input,
	}, nil
}

// handleHotspotStatus handles the "hotspot status" command
func (p *Processor) handleHotspotStatus(input string) (*core.Command, error) {
	return &core.Command{
		Type:     core.CommandTypeConnectivity,
		Action:   "hotspot-status",
		Target:   "",
		RawInput: input,
	}, nil
}
