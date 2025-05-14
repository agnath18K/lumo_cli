package system_test

import (
	"strings"
	"testing"

	"github.com/agnath18/lumo/pkg/system"
)

// TestNewReportGenerator tests the creation of a new report generator
func TestNewReportGenerator(t *testing.T) {
	generator := system.NewReportGenerator()
	if generator == nil {
		t.Fatal("Expected report generator to be created, got nil")
	}
}

// TestFormatSystemReport tests the formatting of system reports
func TestFormatSystemReport(t *testing.T) {
	// Create a mock system report
	report := &system.SystemReport{
		SystemInfo: system.SystemInfo{
			Hostname:      "test-host",
			Platform:      "test-platform",
			Architecture:  "x86_64",
			CPUModel:      "Test CPU",
			CPUCores:      4,
			TotalMemory:   "8.0 GB",
			TotalDisk:     "500.0 GB",
			Uptime:        "2 days, 3 hours, 45 minutes",
			KernelVersion: "5.10.0",
		},
		NetworkInfo: system.NetworkInfo{
			Interfaces: []system.NetworkInterface{
				{
					Name:       "eth0",
					IPAddress:  "192.168.1.100",
					MACAddress: "00:11:22:33:44:55",
					Status:     "UP",
				},
				{
					Name:       "wlan0",
					IPAddress:  "192.168.1.101",
					MACAddress: "AA:BB:CC:DD:EE:FF",
					Status:     "UP",
				},
			},
		},
		SoftwareInfo: system.SoftwareInfo{
			OS:           "linux",
			GoVersion:    "go1.20.1",
			ShellVersion: "GNU bash, version 5.1.16",
			PackageInfo: map[string]string{
				"apt": "apt 2.4.5 (amd64)",
			},
		},
	}

	// Format the system report
	result := system.FormatSystemReport(report)

	// Check that the result is not empty
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that the result contains the system information
	if !strings.Contains(result, "test-host") {
		t.Errorf("Expected result to contain hostname, got: %s", result)
	}
	if !strings.Contains(result, "test-platform") {
		t.Errorf("Expected result to contain platform, got: %s", result)
	}
	if !strings.Contains(result, "x86_64") {
		t.Errorf("Expected result to contain architecture, got: %s", result)
	}
	if !strings.Contains(result, "Test CPU") {
		t.Errorf("Expected result to contain CPU model, got: %s", result)
	}
	if !strings.Contains(result, "4") {
		t.Errorf("Expected result to contain CPU cores, got: %s", result)
	}
	if !strings.Contains(result, "8.0 GB") {
		t.Errorf("Expected result to contain total memory, got: %s", result)
	}
	if !strings.Contains(result, "500.0 GB") {
		t.Errorf("Expected result to contain total disk, got: %s", result)
	}
	if !strings.Contains(result, "2 days, 3 hours, 45 minutes") {
		t.Errorf("Expected result to contain uptime, got: %s", result)
	}
	if !strings.Contains(result, "5.10.0") {
		t.Errorf("Expected result to contain kernel version, got: %s", result)
	}

	// Check that the result contains the network information
	if !strings.Contains(result, "eth0") {
		t.Errorf("Expected result to contain eth0 interface, got: %s", result)
	}
	if !strings.Contains(result, "192.168.1.100") {
		t.Errorf("Expected result to contain eth0 IP, got: %s", result)
	}
	if !strings.Contains(result, "00:11:22:33:44:55") {
		t.Errorf("Expected result to contain eth0 MAC, got: %s", result)
	}
	if !strings.Contains(result, "wlan0") {
		t.Errorf("Expected result to contain wlan0 interface, got: %s", result)
	}
	if !strings.Contains(result, "192.168.1.101") {
		t.Errorf("Expected result to contain wlan0 IP, got: %s", result)
	}
	if !strings.Contains(result, "AA:BB:CC:DD:EE:FF") {
		t.Errorf("Expected result to contain wlan0 MAC, got: %s", result)
	}

	// Check that the result contains the software information
	if !strings.Contains(result, "linux") {
		t.Errorf("Expected result to contain OS, got: %s", result)
	}
	if !strings.Contains(result, "go1.20.1") {
		t.Errorf("Expected result to contain Go version, got: %s", result)
	}
	if !strings.Contains(result, "GNU bash") {
		t.Errorf("Expected result to contain shell version, got: %s", result)
	}
	if !strings.Contains(result, "apt") {
		t.Errorf("Expected result to contain apt package info, got: %s", result)
	}
}
