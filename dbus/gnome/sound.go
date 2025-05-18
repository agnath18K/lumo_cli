package gnome

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/agnath18K/lumo/internal/core"
)

// GNOME sound-related DBus service names and interfaces
const (
	// PulseAudio is the PulseAudio service
	PulseAudio = "org.pulseaudio.Server"
	// PulseAudioPath is the PulseAudio object path
	PulseAudioPath = "/org/pulseaudio/server_lookup1"
	// PulseAudioInterface is the PulseAudio interface
	PulseAudioInterface = "org.pulseaudio.ServerLookup1"

	// GSettingsSchemaSound is the schema for sound settings
	GSettingsSchemaSound = "org.gnome.desktop.sound"
)

// executeSoundCommand executes a sound management command
func (e *Environment) executeSoundCommand(ctx context.Context, cmd *core.Command) (*core.Result, error) {
	switch cmd.Action {
	case "set-volume":
		// Parse volume level
		level, err := parseVolumeLevel(cmd.Target)
		if err != nil {
			return nil, err
		}
		if err := e.SetVolume(ctx, level); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set volume to %d%%", level),
			Success: true,
		}, nil
	case "get-volume":
		volume, err := e.GetVolume(ctx)
		if err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Current volume: %d%%", volume),
			Success: true,
			Data: map[string]any{
				"volume": volume,
			},
		}, nil
	case "set-mute":
		// Parse mute state
		mute := true
		if cmd.Target == "false" || cmd.Target == "off" || cmd.Target == "0" {
			mute = false
		}
		if err := e.SetMute(ctx, mute); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set mute to: %v", mute),
			Success: true,
		}, nil
	case "get-mute":
		mute, err := e.GetMute(ctx)
		if err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Mute state: %v", mute),
			Success: true,
			Data: map[string]any{
				"mute": mute,
			},
		}, nil
	case "set-input-volume":
		// Parse volume level
		level, err := parseVolumeLevel(cmd.Target)
		if err != nil {
			return nil, err
		}
		if err := e.SetInputVolume(ctx, level); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set microphone volume to %d%%", level),
			Success: true,
		}, nil
	case "get-input-volume":
		volume, err := e.GetInputVolume(ctx)
		if err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Current microphone volume: %d%%", volume),
			Success: true,
			Data: map[string]any{
				"input_volume": volume,
			},
		}, nil
	case "set-input-mute":
		// Parse mute state
		mute := true
		if cmd.Target == "false" || cmd.Target == "off" || cmd.Target == "0" {
			mute = false
		}
		if err := e.SetInputMute(ctx, mute); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set microphone mute to: %v", mute),
			Success: true,
		}, nil
	case "get-input-mute":
		mute, err := e.GetInputMute(ctx)
		if err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Microphone mute state: %v", mute),
			Success: true,
			Data: map[string]any{
				"input_mute": mute,
			},
		}, nil
	case "list-devices":
		devices, err := e.GetSoundDevices(ctx)
		if err != nil {
			return nil, err
		}
		var output strings.Builder
		output.WriteString("Sound devices:\n")
		for _, device := range devices {
			deviceType := "Output"
			if device.IsInput {
				deviceType = "Input"
			}
			defaultMark := ""
			if device.IsDefault {
				defaultMark = " (default)"
			}
			output.WriteString(fmt.Sprintf("- %s: %s%s\n", deviceType, device.Name, defaultMark))
		}
		return &core.Result{
			Output:  output.String(),
			Success: true,
			Data: map[string]any{
				"devices": devices,
			},
		}, nil
	case "set-default-device":
		if cmd.Target == "" {
			return nil, fmt.Errorf("device ID is required")
		}
		if err := e.SetDefaultSoundDevice(ctx, cmd.Target); err != nil {
			return nil, err
		}
		return &core.Result{
			Output:  fmt.Sprintf("Set default sound device to: %s", cmd.Target),
			Success: true,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported sound action: %s", cmd.Action)
	}
}

// SetVolume sets the system volume level (0-100)
func (e *Environment) SetVolume(ctx context.Context, level int) error {
	// Ensure level is within valid range
	if level < 0 {
		level = 0
	} else if level > 100 {
		level = 100
	}

	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using gsettings as a fallback
		return e.setVolumeWithGSettings(level)
	}

	// Use pactl to set the volume
	cmd := fmt.Sprintf("pactl set-sink-volume @DEFAULT_SINK@ %d%%", level)
	_, err = e.runCommand(cmd)
	if err != nil {
		// Try using gsettings as a fallback
		return e.setVolumeWithGSettings(level)
	}
	return nil
}

// setVolumeWithGSettings sets the volume using a fallback method
func (e *Environment) setVolumeWithGSettings(level int) error {
	// Try to set the volume using amixer as a fallback
	cmd := fmt.Sprintf("amixer set Master %d%%", level)
	_, err := e.runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to set volume with amixer: %w", err)
	}
	return nil
}

// GetVolume gets the current system volume level (0-100)
func (e *Environment) GetVolume(ctx context.Context) (int, error) {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using gsettings as a fallback
		return e.getVolumeWithGSettings()
	}

	// Use pactl to get the volume
	cmd := "pactl get-sink-volume @DEFAULT_SINK@"
	output, err := e.runCommand(cmd)
	if err != nil {
		// Try using gsettings as a fallback
		return e.getVolumeWithGSettings()
	}

	// Parse the output to extract the volume level
	volume, err := parseVolumeFromPactl(output)
	if err != nil {
		// Try using gsettings as a fallback
		return e.getVolumeWithGSettings()
	}

	return volume, nil
}

// getVolumeWithGSettings gets the volume using a fallback method
func (e *Environment) getVolumeWithGSettings() (int, error) {
	// Try to get the volume using amixer as a fallback
	cmd := "amixer get Master | grep -o '[0-9]*%' | head -1 | tr -d '%'"
	output, err := e.runCommand(cmd)
	if err != nil {
		// If amixer fails, return a default value
		return 50, fmt.Errorf("failed to get volume with amixer: %w", err)
	}

	// Parse the output (should be a percentage)
	output = strings.TrimSpace(output)

	// Convert to int
	volumePercent, err := strconv.Atoi(output)
	if err != nil {
		// If parsing fails, return a default value
		return 50, fmt.Errorf("failed to parse volume from amixer: %w", err)
	}

	return volumePercent, nil
}

// SetMute sets the system mute state
func (e *Environment) SetMute(ctx context.Context, mute bool) error {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using amixer as a fallback
		return e.setMuteWithAmixer(mute)
	}

	// Use pactl to set the mute state
	muteStr := "1"
	if !mute {
		muteStr = "0"
	}
	cmd := fmt.Sprintf("pactl set-sink-mute @DEFAULT_SINK@ %s", muteStr)
	_, err = e.runCommand(cmd)
	if err != nil {
		// Try using amixer as a fallback
		return e.setMuteWithAmixer(mute)
	}
	return nil
}

// setMuteWithAmixer sets the mute state using amixer
func (e *Environment) setMuteWithAmixer(mute bool) error {
	// Use amixer to set the mute state
	muteStr := "mute"
	if !mute {
		muteStr = "unmute"
	}
	cmd := fmt.Sprintf("amixer set Master %s", muteStr)
	_, err := e.runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to set mute with amixer: %w", err)
	}
	return nil
}

// GetMute gets the current system mute state
func (e *Environment) GetMute(ctx context.Context) (bool, error) {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using amixer as a fallback
		return e.getMuteWithAmixer()
	}

	// Use pactl to get the mute state
	cmd := "pactl get-sink-mute @DEFAULT_SINK@"
	output, err := e.runCommand(cmd)
	if err != nil {
		// Try using amixer as a fallback
		return e.getMuteWithAmixer()
	}

	// Parse the output to extract the mute state
	return strings.Contains(output, "yes"), nil
}

// getMuteWithAmixer gets the mute state using amixer
func (e *Environment) getMuteWithAmixer() (bool, error) {
	// Use amixer to get the mute state
	cmd := "amixer get Master | grep -o '\\[on\\]\\|\\[off\\]' | head -1"
	output, err := e.runCommand(cmd)
	if err != nil {
		return false, fmt.Errorf("failed to get mute state with amixer: %w", err)
	}

	// Parse the output to extract the mute state
	return !strings.Contains(output, "on"), nil
}

// SetInputVolume sets the microphone volume level (0-100)
func (e *Environment) SetInputVolume(ctx context.Context, level int) error {
	// Ensure level is within valid range
	if level < 0 {
		level = 0
	} else if level > 100 {
		level = 100
	}

	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using amixer as a fallback
		return e.setInputVolumeWithAmixer(level)
	}

	// Use pactl to set the input volume
	cmd := fmt.Sprintf("pactl set-source-volume @DEFAULT_SOURCE@ %d%%", level)
	_, err = e.runCommand(cmd)
	if err != nil {
		// Try using amixer as a fallback
		return e.setInputVolumeWithAmixer(level)
	}
	return nil
}

// setInputVolumeWithAmixer sets the microphone volume using amixer
func (e *Environment) setInputVolumeWithAmixer(level int) error {
	// Try to set the microphone volume using amixer
	// First try with "Capture" which is common for microphones
	cmd := fmt.Sprintf("amixer set Capture %d%%", level)
	_, err := e.runCommand(cmd)
	if err != nil {
		// If that fails, try with "Mic" which is another common name
		cmd = fmt.Sprintf("amixer set Mic %d%%", level)
		_, err = e.runCommand(cmd)
		if err != nil {
			// If that fails too, try with "Input" as a last resort
			cmd = fmt.Sprintf("amixer set Input %d%%", level)
			_, err = e.runCommand(cmd)
			if err != nil {
				return fmt.Errorf("failed to set microphone volume with amixer: %w", err)
			}
		}
	}
	return nil
}

// GetInputVolume gets the current microphone volume level (0-100)
func (e *Environment) GetInputVolume(ctx context.Context) (int, error) {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using amixer as a fallback
		return e.getInputVolumeWithAmixer()
	}

	// Use pactl to get the input volume
	cmd := "pactl get-source-volume @DEFAULT_SOURCE@"
	output, err := e.runCommand(cmd)
	if err != nil {
		// Try using amixer as a fallback
		return e.getInputVolumeWithAmixer()
	}

	// Parse the output to extract the volume level
	volume, err := parseVolumeFromPactl(output)
	if err != nil {
		// Try using amixer as a fallback
		return e.getInputVolumeWithAmixer()
	}

	return volume, nil
}

// getInputVolumeWithAmixer gets the microphone volume using amixer
func (e *Environment) getInputVolumeWithAmixer() (int, error) {
	// Try to get the microphone volume using amixer
	// First try with "Capture" which is common for microphones
	cmd := "amixer get Capture | grep -o '[0-9]*%' | head -1 | tr -d '%'"
	output, err := e.runCommand(cmd)
	if err == nil && output != "" {
		// Parse the output (should be a percentage)
		output = strings.TrimSpace(output)
		volume, err := strconv.Atoi(output)
		if err == nil {
			return volume, nil
		}
	}

	// If that fails, try with "Mic"
	cmd = "amixer get Mic | grep -o '[0-9]*%' | head -1 | tr -d '%'"
	output, err = e.runCommand(cmd)
	if err == nil && output != "" {
		// Parse the output (should be a percentage)
		output = strings.TrimSpace(output)
		volume, err := strconv.Atoi(output)
		if err == nil {
			return volume, nil
		}
	}

	// If that fails too, try with "Input"
	cmd = "amixer get Input | grep -o '[0-9]*%' | head -1 | tr -d '%'"
	output, err = e.runCommand(cmd)
	if err == nil && output != "" {
		// Parse the output (should be a percentage)
		output = strings.TrimSpace(output)
		volume, err := strconv.Atoi(output)
		if err == nil {
			return volume, nil
		}
	}

	// If all attempts fail, return a default value
	return 50, fmt.Errorf("failed to get microphone volume with amixer")
}

// SetInputMute sets the microphone mute state
func (e *Environment) SetInputMute(ctx context.Context, mute bool) error {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using amixer as a fallback
		return e.setInputMuteWithAmixer(mute)
	}

	// Use pactl to set the input mute state
	muteStr := "1"
	if !mute {
		muteStr = "0"
	}
	cmd := fmt.Sprintf("pactl set-source-mute @DEFAULT_SOURCE@ %s", muteStr)
	_, err = e.runCommand(cmd)
	if err != nil {
		// Try using amixer as a fallback
		return e.setInputMuteWithAmixer(mute)
	}
	return nil
}

// setInputMuteWithAmixer sets the microphone mute state using amixer
func (e *Environment) setInputMuteWithAmixer(mute bool) error {
	// Use amixer to set the microphone mute state
	muteStr := "mute"
	if !mute {
		muteStr = "unmute"
	}

	// Try with "Capture" which is common for microphones
	cmd := fmt.Sprintf("amixer set Capture %s", muteStr)
	_, err := e.runCommand(cmd)
	if err != nil {
		// If that fails, try with "Mic"
		cmd = fmt.Sprintf("amixer set Mic %s", muteStr)
		_, err = e.runCommand(cmd)
		if err != nil {
			// If that fails too, try with "Input"
			cmd = fmt.Sprintf("amixer set Input %s", muteStr)
			_, err = e.runCommand(cmd)
			if err != nil {
				return fmt.Errorf("failed to set microphone mute with amixer: %w", err)
			}
		}
	}
	return nil
}

// GetInputMute gets the current microphone mute state
func (e *Environment) GetInputMute(ctx context.Context) (bool, error) {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using amixer as a fallback
		return e.getInputMuteWithAmixer()
	}

	// Use pactl to get the input mute state
	cmd := "pactl get-source-mute @DEFAULT_SOURCE@"
	output, err := e.runCommand(cmd)
	if err != nil {
		// Try using amixer as a fallback
		return e.getInputMuteWithAmixer()
	}

	// Parse the output to extract the mute state
	return strings.Contains(output, "yes"), nil
}

// getInputMuteWithAmixer gets the microphone mute state using amixer
func (e *Environment) getInputMuteWithAmixer() (bool, error) {
	// Try with "Capture" which is common for microphones
	cmd := "amixer get Capture | grep -o '\\[on\\]\\|\\[off\\]' | head -1"
	output, err := e.runCommand(cmd)
	if err == nil && output != "" {
		return !strings.Contains(output, "on"), nil
	}

	// If that fails, try with "Mic"
	cmd = "amixer get Mic | grep -o '\\[on\\]\\|\\[off\\]' | head -1"
	output, err = e.runCommand(cmd)
	if err == nil && output != "" {
		return !strings.Contains(output, "on"), nil
	}

	// If that fails too, try with "Input"
	cmd = "amixer get Input | grep -o '\\[on\\]\\|\\[off\\]' | head -1"
	output, err = e.runCommand(cmd)
	if err == nil && output != "" {
		return !strings.Contains(output, "on"), nil
	}

	// If all attempts fail, return a default value
	return false, fmt.Errorf("failed to get microphone mute state with amixer")
}

// GetSoundDevices gets a list of available sound devices
func (e *Environment) GetSoundDevices(ctx context.Context) ([]core.SoundDevice, error) {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using amixer as a fallback
		return e.getSoundDevicesWithAmixer()
	}

	var devices []core.SoundDevice

	// Get output devices
	outputDevices, err := e.getSoundDevicesByType(false)
	if err != nil {
		// Try using amixer as a fallback
		return e.getSoundDevicesWithAmixer()
	}
	devices = append(devices, outputDevices...)

	// Get input devices
	inputDevices, err := e.getSoundDevicesByType(true)
	if err != nil {
		// We already have output devices, so just add some default input devices
		inputDevices, _ = e.getDefaultInputDevices()
		devices = append(devices, inputDevices...)
		return devices, nil
	}
	devices = append(devices, inputDevices...)

	return devices, nil
}

// getSoundDevicesWithAmixer gets a list of sound devices using amixer
func (e *Environment) getSoundDevicesWithAmixer() ([]core.SoundDevice, error) {
	var devices []core.SoundDevice

	// Get a list of controls from amixer
	cmd := "amixer controls"
	output, err := e.runCommand(cmd)
	if err != nil {
		// If amixer fails, return some default devices
		return e.getDefaultSoundDevices()
	}

	// Parse the output to extract device information
	// This is a simplified approach and might not work for all systems
	lines := strings.Split(output, "\n")

	// Track which devices we've already added to avoid duplicates
	addedDevices := make(map[string]bool)

	for _, line := range lines {
		if strings.Contains(line, "Playback") {
			// This is an output device
			name := extractDeviceNameFromAmixer(line)
			if name != "" && !addedDevices[name] {
				addedDevices[name] = true

				// Get volume and mute state
				volume, muted := e.getDeviceVolumeAndMute(name, false)

				device := core.SoundDevice{
					ID:          name,
					Name:        name,
					Description: "Audio output device",
					IsInput:     false,
					IsDefault:   strings.Contains(line, "Master") || strings.Contains(line, "PCM"),
					Volume:      volume,
					Muted:       muted,
				}

				devices = append(devices, device)
			}
		} else if strings.Contains(line, "Capture") {
			// This is an input device
			name := extractDeviceNameFromAmixer(line)
			if name != "" && !addedDevices[name] {
				addedDevices[name] = true

				// Get volume and mute state
				volume, muted := e.getDeviceVolumeAndMute(name, true)

				device := core.SoundDevice{
					ID:          name,
					Name:        name,
					Description: "Audio input device",
					IsInput:     true,
					IsDefault:   strings.Contains(line, "Mic") || strings.Contains(line, "Capture"),
					Volume:      volume,
					Muted:       muted,
				}

				devices = append(devices, device)
			}
		}
	}

	// If we couldn't find any devices, return some default ones
	if len(devices) == 0 {
		return e.getDefaultSoundDevices()
	}

	return devices, nil
}

// getDefaultSoundDevices returns a list of default sound devices
func (e *Environment) getDefaultSoundDevices() ([]core.SoundDevice, error) {
	var devices []core.SoundDevice

	// Add default output device
	outputVolume, outputMuted := e.getDeviceVolumeAndMute("Master", false)
	outputDevice := core.SoundDevice{
		ID:          "default_output",
		Name:        "Default Output",
		Description: "Default audio output device",
		IsInput:     false,
		IsDefault:   true,
		Volume:      outputVolume,
		Muted:       outputMuted,
	}
	devices = append(devices, outputDevice)

	// Add default input devices
	inputDevices, _ := e.getDefaultInputDevices()
	devices = append(devices, inputDevices...)

	return devices, nil
}

// getDefaultInputDevices returns a list of default input devices
func (e *Environment) getDefaultInputDevices() ([]core.SoundDevice, error) {
	var devices []core.SoundDevice

	// Add default microphone
	inputVolume, inputMuted := e.getDeviceVolumeAndMute("Capture", true)
	inputDevice := core.SoundDevice{
		ID:          "default_input",
		Name:        "Default Microphone",
		Description: "Default audio input device",
		IsInput:     true,
		IsDefault:   true,
		Volume:      inputVolume,
		Muted:       inputMuted,
	}
	devices = append(devices, inputDevice)

	return devices, nil
}

// getDeviceVolumeAndMute gets the volume and mute state for a device
func (e *Environment) getDeviceVolumeAndMute(device string, isInput bool) (int, bool) {
	// Get volume
	var volume int = 50 // Default value
	var cmd string

	cmd = fmt.Sprintf("amixer get %s | grep -o '[0-9]*%%' | head -1 | tr -d '%%'", device)
	output, err := e.runCommand(cmd)
	if err == nil && output != "" {
		output = strings.TrimSpace(output)
		vol, err := strconv.Atoi(output)
		if err == nil {
			volume = vol
		}
	}

	// Get mute state
	var muted bool = false // Default value

	cmd = fmt.Sprintf("amixer get %s | grep -o '\\[on\\]\\|\\[off\\]' | head -1", device)
	output, err = e.runCommand(cmd)
	if err == nil && output != "" {
		muted = !strings.Contains(output, "on")
	}

	return volume, muted
}

// extractDeviceNameFromAmixer extracts the device name from an amixer control line
func extractDeviceNameFromAmixer(line string) string {
	// Extract the name from something like "numid=1,iface=MIXER,name='Master Playback Volume'"
	nameStart := strings.Index(line, "name='")
	if nameStart == -1 {
		return ""
	}

	nameStart += 6 // Skip "name='"
	nameEnd := strings.Index(line[nameStart:], "'")
	if nameEnd == -1 {
		return ""
	}

	name := line[nameStart : nameStart+nameEnd]

	// Simplify the name by removing common suffixes
	name = strings.TrimSuffix(name, " Playback Volume")
	name = strings.TrimSuffix(name, " Capture Volume")
	name = strings.TrimSuffix(name, " Playback Switch")
	name = strings.TrimSuffix(name, " Capture Switch")

	return name
}

// getSoundDevicesByType gets a list of sound devices by type (input or output)
func (e *Environment) getSoundDevicesByType(isInput bool) ([]core.SoundDevice, error) {
	var devices []core.SoundDevice
	var cmd string

	if isInput {
		cmd = "pactl list sources"
	} else {
		cmd = "pactl list sinks"
	}

	output, err := e.runCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list %s devices: %w", getDeviceTypeString(isInput), err)
	}

	// Parse the output to extract device information
	// This is a simplified parsing and might need to be improved for more complex setups
	sections := strings.Split(output, "Sink #")
	if isInput {
		sections = strings.Split(output, "Source #")
	}

	for i, section := range sections {
		if i == 0 {
			continue // Skip the header
		}

		lines := strings.Split(section, "\n")
		if len(lines) < 2 {
			continue
		}

		// Extract device ID
		idParts := strings.Fields(lines[0])
		id := ""
		if len(idParts) > 0 {
			id = idParts[0]
		}

		// Extract device name
		name := ""
		description := ""
		isDefault := false
		volume := 0
		muted := false

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Name:") {
				name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			} else if strings.HasPrefix(line, "Description:") {
				description = strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
			} else if strings.HasPrefix(line, "State:") {
				isDefault = strings.Contains(line, "RUNNING")
			} else if strings.HasPrefix(line, "Volume:") {
				vol, err := parseVolumeFromPactl(line)
				if err == nil {
					volume = vol
				}
			} else if strings.HasPrefix(line, "Mute:") {
				muted = strings.Contains(line, "yes")
			}
		}

		device := core.SoundDevice{
			ID:          id,
			Name:        name,
			Description: description,
			IsInput:     isInput,
			IsDefault:   isDefault,
			Volume:      volume,
			Muted:       muted,
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// SetDefaultSoundDevice sets the default sound device
func (e *Environment) SetDefaultSoundDevice(ctx context.Context, deviceID string) error {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try using asoundrc as a fallback (this is a simplified approach)
		return e.setDefaultSoundDeviceWithAsoundrc(deviceID)
	}

	// Check if this is an input or output device
	isInput, err := e.isInputDevice(deviceID)
	if err != nil {
		// If we can't determine the device type, try using asoundrc as a fallback
		return e.setDefaultSoundDeviceWithAsoundrc(deviceID)
	}

	var cmd string
	if isInput {
		cmd = fmt.Sprintf("pactl set-default-source %s", deviceID)
	} else {
		cmd = fmt.Sprintf("pactl set-default-sink %s", deviceID)
	}

	_, err = e.runCommand(cmd)
	if err != nil {
		// Try using asoundrc as a fallback
		return e.setDefaultSoundDeviceWithAsoundrc(deviceID)
	}

	return nil
}

// setDefaultSoundDeviceWithAsoundrc sets the default sound device using .asoundrc
func (e *Environment) setDefaultSoundDeviceWithAsoundrc(deviceID string) error {
	// This is a simplified approach and might not work for all systems
	// In a real implementation, you would need to create or modify the .asoundrc file

	// For now, just return a message that this is not fully implemented
	return fmt.Errorf("setting default sound device without pactl is not fully implemented. Device ID: %s", deviceID)
}

// isInputDevice checks if a device is an input device
func (e *Environment) isInputDevice(deviceID string) (bool, error) {
	// Check if pactl is installed
	_, err := exec.LookPath("pactl")
	if err != nil {
		// Try to infer from the device ID
		return e.inferDeviceTypeFromID(deviceID)
	}

	// Check if the device exists in the list of input devices
	cmd := "pactl list sources short"
	output, err := e.runCommand(cmd)
	if err != nil {
		// Try to infer from the device ID
		return e.inferDeviceTypeFromID(deviceID)
	}

	if strings.Contains(output, deviceID) {
		return true, nil
	}

	// Check if the device exists in the list of output devices
	cmd = "pactl list sinks short"
	output, err = e.runCommand(cmd)
	if err != nil {
		// Try to infer from the device ID
		return e.inferDeviceTypeFromID(deviceID)
	}

	if strings.Contains(output, deviceID) {
		return false, nil
	}

	// If we can't find the device, try to infer from the device ID
	return e.inferDeviceTypeFromID(deviceID)
}

// inferDeviceTypeFromID tries to infer if a device is an input device from its ID
func (e *Environment) inferDeviceTypeFromID(deviceID string) (bool, error) {
	// Common input device identifiers
	inputIdentifiers := []string{
		"mic", "microphone", "input", "capture", "source", "default_input",
	}

	// Check if the device ID contains any input identifiers
	deviceIDLower := strings.ToLower(deviceID)
	for _, identifier := range inputIdentifiers {
		if strings.Contains(deviceIDLower, identifier) {
			return true, nil
		}
	}

	// If it doesn't match any input identifiers, assume it's an output device
	return false, nil
}

// parseVolumeLevel parses a volume level from a string
func parseVolumeLevel(volumeStr string) (int, error) {
	// Remove any % sign
	volumeStr = strings.TrimSuffix(volumeStr, "%")

	// Parse the volume level
	level, err := strconv.Atoi(volumeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid volume level: %s", volumeStr)
	}

	// Ensure level is within valid range
	if level < 0 {
		level = 0
	} else if level > 100 {
		level = 100
	}

	return level, nil
}

// parseVolumeFromPactl parses the volume level from pactl output
func parseVolumeFromPactl(output string) (int, error) {
	// Look for percentage values
	percentIndex := strings.Index(output, "%")
	if percentIndex == -1 {
		return 0, fmt.Errorf("no volume percentage found in output: %s", output)
	}

	// Extract the number before the % sign
	start := percentIndex - 1
	for start >= 0 && (output[start] >= '0' && output[start] <= '9' || output[start] == ' ') {
		start--
	}
	start++

	volumeStr := strings.TrimSpace(output[start:percentIndex])
	volume, err := strconv.Atoi(volumeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse volume: %w", err)
	}

	return volume, nil
}

// getDeviceTypeString returns a string representation of the device type
func getDeviceTypeString(isInput bool) string {
	if isInput {
		return "input"
	}
	return "output"
}

// Note: runCommand method is already defined in appearance.go
