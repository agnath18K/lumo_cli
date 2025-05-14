package system_test

import (
	"strings"
	"testing"

	"github.com/agnath18/lumo/pkg/system"
)

// TestNewHealthChecker tests the creation of a new health checker
func TestNewHealthChecker(t *testing.T) {
	checker := system.NewHealthChecker()
	if checker == nil {
		t.Fatal("Expected health checker to be created, got nil")
	}
}

// TestFormatHealthCheck tests the formatting of health check results
func TestFormatHealthCheck(t *testing.T) {
	// Create a mock health check result
	health := &system.SystemHealth{
		Hostname: "test-host",
		Platform: "test-platform",
		Checks: []system.HealthCheck{
			{
				Component:   "CPU",
				Status:      system.StatusHealthy,
				Value:       "10.0%",
				Description: "CPU usage is 10.0%",
			},
			{
				Component:   "Memory",
				Status:      system.StatusWarning,
				Value:       "85.0% (6.8 GB / 8.0 GB)",
				Description: "Memory usage is 85.0% (6.8 GB used out of 8.0 GB)",
				Advice:      "Memory usage is high, consider freeing up memory",
			},
			{
				Component:   "Disk",
				Status:      system.StatusCritical,
				Value:       "95.0% (950 GB / 1000 GB)",
				Description: "Disk usage is 95.0% (950 GB used out of 1000 GB)",
				Advice:      "Disk space critically low, delete unnecessary files",
			},
		},
		Summary: "System health is CRITICAL: 1 critical issues, 1 warnings",
	}

	// Format the health check
	result := system.FormatHealthCheck(health)

	// Check that the result is not empty
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that the result contains the hostname
	if !strings.Contains(result, "test-host") {
		t.Errorf("Expected result to contain hostname, got: %s", result)
	}

	// Check that the result contains the platform
	if !strings.Contains(result, "test-platform") {
		t.Errorf("Expected result to contain platform, got: %s", result)
	}

	// Check that the result contains the CPU status
	if !strings.Contains(result, "CPU") || !strings.Contains(result, "10.0%") {
		t.Errorf("Expected result to contain CPU status, got: %s", result)
	}

	// Check that the result contains the Memory status
	if !strings.Contains(result, "Memory") || !strings.Contains(result, "85.0%") {
		t.Errorf("Expected result to contain Memory status, got: %s", result)
	}

	// Check that the result contains the Disk status
	if !strings.Contains(result, "Disk") || !strings.Contains(result, "95.0%") {
		t.Errorf("Expected result to contain Disk status, got: %s", result)
	}

	// Check that the result contains the summary
	if !strings.Contains(result, "CRITICAL") || !strings.Contains(result, "1 critical issues") {
		t.Errorf("Expected result to contain summary, got: %s", result)
	}

	// Check that the result contains the advice for Memory
	if !strings.Contains(result, "Memory usage is high") {
		t.Errorf("Expected result to contain Memory advice, got: %s", result)
	}

	// Check that the result contains the advice for Disk
	if !strings.Contains(result, "Disk space critically low") {
		t.Errorf("Expected result to contain Disk advice, got: %s", result)
	}
}
