package system

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// HealthStatus represents the status of a system component
type HealthStatus string

const (
	// StatusHealthy indicates a healthy component
	StatusHealthy HealthStatus = "HEALTHY"
	// StatusWarning indicates a component with potential issues
	StatusWarning HealthStatus = "WARNING"
	// StatusCritical indicates a component with critical issues
	StatusCritical HealthStatus = "CRITICAL"
)

// HealthCheck represents a health check result
type HealthCheck struct {
	Component   string       `json:"component"`
	Status      HealthStatus `json:"status"`
	Value       string       `json:"value"`
	Description string       `json:"description"`
	Threshold   string       `json:"threshold,omitempty"`
	Advice      string       `json:"advice,omitempty"`
}

// SystemHealth represents the overall system health
type SystemHealth struct {
	Timestamp time.Time     `json:"timestamp"`
	Hostname  string        `json:"hostname"`
	Platform  string        `json:"platform"`
	Checks    []HealthCheck `json:"checks"`
	Summary   string        `json:"summary"`
}

// HealthChecker handles system health checks
type HealthChecker struct {
	warningThresholdCPU     float64
	criticalThresholdCPU    float64
	warningThresholdMemory  float64
	criticalThresholdMemory float64
	warningThresholdDisk    float64
	criticalThresholdDisk   float64
}

// NewHealthChecker creates a new health checker with default thresholds
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		warningThresholdCPU:     70.0, // 70% CPU usage is a warning
		criticalThresholdCPU:    90.0, // 90% CPU usage is critical
		warningThresholdMemory:  80.0, // 80% memory usage is a warning
		criticalThresholdMemory: 90.0, // 90% memory usage is critical
		warningThresholdDisk:    85.0, // 85% disk usage is a warning
		criticalThresholdDisk:   95.0, // 95% disk usage is critical
	}
}

// CheckHealth performs a comprehensive system health check
func (h *HealthChecker) CheckHealth() (*SystemHealth, error) {
	// Create a new system health object
	health := &SystemHealth{
		Timestamp: time.Now(),
		Checks:    []HealthCheck{},
	}

	// Get hostname and platform information
	hostInfo, err := host.Info()
	if err == nil {
		health.Hostname = hostInfo.Hostname
		health.Platform = fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.KernelVersion)
	}

	// Check CPU usage
	cpuCheck, err := h.checkCPU()
	if err == nil {
		health.Checks = append(health.Checks, cpuCheck)
	}

	// Check memory usage
	memCheck, err := h.checkMemory()
	if err == nil {
		health.Checks = append(health.Checks, memCheck)
	}

	// Check disk usage
	diskCheck, err := h.checkDisk()
	if err == nil {
		health.Checks = append(health.Checks, diskCheck)
	}

	// Generate summary
	health.Summary = h.generateSummary(health.Checks)

	return health, nil
}

// checkCPU checks CPU usage
func (h *HealthChecker) checkCPU() (HealthCheck, error) {
	check := HealthCheck{
		Component: "CPU",
		Status:    StatusHealthy,
	}

	// Get CPU usage percentage
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return check, err
	}

	if len(percentages) > 0 {
		cpuUsage := percentages[0]
		check.Value = fmt.Sprintf("%.1f%%", cpuUsage)
		check.Description = fmt.Sprintf("CPU usage is %.1f%%", cpuUsage)

		// Set status based on thresholds
		if cpuUsage >= h.criticalThresholdCPU {
			check.Status = StatusCritical
			check.Advice = "Consider closing resource-intensive applications or processes"
		} else if cpuUsage >= h.warningThresholdCPU {
			check.Status = StatusWarning
			check.Advice = "CPU usage is high, monitor for performance issues"
		}

		check.Threshold = fmt.Sprintf("Warning: %.1f%%, Critical: %.1f%%", h.warningThresholdCPU, h.criticalThresholdCPU)
	}

	return check, nil
}

// checkMemory checks memory usage
func (h *HealthChecker) checkMemory() (HealthCheck, error) {
	check := HealthCheck{
		Component: "Memory",
		Status:    StatusHealthy,
	}

	// Get memory usage
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return check, err
	}

	memUsage := memInfo.UsedPercent
	usedGB := float64(memInfo.Used) / (1024 * 1024 * 1024)
	totalGB := float64(memInfo.Total) / (1024 * 1024 * 1024)

	check.Value = fmt.Sprintf("%.1f%% (%.1f GB / %.1f GB)", memUsage, usedGB, totalGB)
	check.Description = fmt.Sprintf("Memory usage is %.1f%% (%.1f GB used out of %.1f GB)", memUsage, usedGB, totalGB)

	// Set status based on thresholds
	if memUsage >= h.criticalThresholdMemory {
		check.Status = StatusCritical
		check.Advice = "Close unnecessary applications to free up memory"
	} else if memUsage >= h.warningThresholdMemory {
		check.Status = StatusWarning
		check.Advice = "Memory usage is high, consider freeing up memory"
	}

	check.Threshold = fmt.Sprintf("Warning: %.1f%%, Critical: %.1f%%", h.warningThresholdMemory, h.criticalThresholdMemory)

	return check, nil
}

// checkDisk checks disk usage
func (h *HealthChecker) checkDisk() (HealthCheck, error) {
	check := HealthCheck{
		Component: "Disk",
		Status:    StatusHealthy,
	}

	// Get disk usage for root partition
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return check, err
	}

	diskUsage := diskInfo.UsedPercent
	usedGB := float64(diskInfo.Used) / (1024 * 1024 * 1024)
	totalGB := float64(diskInfo.Total) / (1024 * 1024 * 1024)

	check.Value = fmt.Sprintf("%.1f%% (%.1f GB / %.1f GB)", diskUsage, usedGB, totalGB)
	check.Description = fmt.Sprintf("Disk usage is %.1f%% (%.1f GB used out of %.1f GB)", diskUsage, usedGB, totalGB)

	// Set status based on thresholds
	if diskUsage >= h.criticalThresholdDisk {
		check.Status = StatusCritical
		check.Advice = "Disk space critically low, delete unnecessary files"
	} else if diskUsage >= h.warningThresholdDisk {
		check.Status = StatusWarning
		check.Advice = "Disk space is running low, consider cleaning up"
	}

	check.Threshold = fmt.Sprintf("Warning: %.1f%%, Critical: %.1f%%", h.warningThresholdDisk, h.criticalThresholdDisk)

	return check, nil
}

// generateSummary generates a summary of the health checks
func (h *HealthChecker) generateSummary(checks []HealthCheck) string {
	criticalCount := 0
	warningCount := 0
	healthyCount := 0

	for _, check := range checks {
		switch check.Status {
		case StatusCritical:
			criticalCount++
		case StatusWarning:
			warningCount++
		case StatusHealthy:
			healthyCount++
		}
	}

	if criticalCount > 0 {
		return fmt.Sprintf("System health is CRITICAL: %d critical issues, %d warnings", criticalCount, warningCount)
	} else if warningCount > 0 {
		return fmt.Sprintf("System health needs attention: %d warnings", warningCount)
	}
	return "System is healthy"
}

// FormatHealthCheck formats a health check result for display
func FormatHealthCheck(health *SystemHealth) string {
	var sb strings.Builder
	boxWidth := 60

	// Format header
	headerText := fmt.Sprintf(" System Health Check (%s) ", health.Timestamp.Format("2006-01-02 15:04:05"))
	sb.WriteString("╭" + padCenter(headerText, boxWidth-2, "─") + "╮\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Host: %s", health.Hostname), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Platform: %s", health.Platform), boxWidth-4) + " │\n")
	sb.WriteString("├" + strings.Repeat("─", boxWidth-2) + "┤\n")

	// Format checks
	for _, check := range health.Checks {
		var statusSymbol string
		switch check.Status {
		case StatusHealthy:
			statusSymbol = "✅"
		case StatusWarning:
			statusSymbol = "⚠️"
		case StatusCritical:
			statusSymbol = "❌"
		}

		sb.WriteString("│ " + padRight(fmt.Sprintf("%s %s: %s", statusSymbol, check.Component, check.Value), boxWidth-4) + " │\n")
		sb.WriteString("│   " + padRight(truncateString(check.Description, boxWidth-8), boxWidth-6) + " │\n")

		if check.Status != StatusHealthy && check.Advice != "" {
			sb.WriteString("│   " + padRight(fmt.Sprintf("Advice: %s", check.Advice), boxWidth-6) + " │\n")
		}
	}

	// Format summary
	sb.WriteString("├" + strings.Repeat("─", boxWidth-2) + "┤\n")
	sb.WriteString("│ " + padRight(health.Summary, boxWidth-4) + " │\n")
	sb.WriteString("╰" + strings.Repeat("─", boxWidth-2) + "╯\n")

	return sb.String()
}
