package system

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// SystemInfo represents system information
type SystemInfo struct {
	Hostname      string `json:"hostname"`
	Platform      string `json:"platform"`
	Architecture  string `json:"architecture"`
	CPUModel      string `json:"cpu_model"`
	CPUCores      int    `json:"cpu_cores"`
	TotalMemory   string `json:"total_memory"`
	TotalDisk     string `json:"total_disk"`
	Uptime        string `json:"uptime"`
	KernelVersion string `json:"kernel_version"`
}

// NetworkInfo represents network information
type NetworkInfo struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name       string `json:"name"`
	IPAddress  string `json:"ip_address"`
	MACAddress string `json:"mac_address"`
	Status     string `json:"status"`
}

// SoftwareInfo represents software information
type SoftwareInfo struct {
	OS           string            `json:"os"`
	GoVersion    string            `json:"go_version"`
	ShellVersion string            `json:"shell_version"`
	PackageInfo  map[string]string `json:"package_info,omitempty"`
}

// SystemReport represents a complete system report
type SystemReport struct {
	Timestamp    time.Time    `json:"timestamp"`
	SystemInfo   SystemInfo   `json:"system_info"`
	NetworkInfo  NetworkInfo  `json:"network_info"`
	SoftwareInfo SoftwareInfo `json:"software_info"`
}

// ReportGenerator handles system report generation
type ReportGenerator struct{}

// NewReportGenerator creates a new report generator
func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{}
}

// GenerateReport generates a comprehensive system report
func (r *ReportGenerator) GenerateReport() (*SystemReport, error) {
	report := &SystemReport{
		Timestamp: time.Now(),
	}

	// Get system information
	systemInfo, err := r.getSystemInfo()
	if err == nil {
		report.SystemInfo = systemInfo
	}

	// Get network information
	networkInfo, err := r.getNetworkInfo()
	if err == nil {
		report.NetworkInfo = networkInfo
	}

	// Get software information
	softwareInfo, err := r.getSoftwareInfo()
	if err == nil {
		report.SoftwareInfo = softwareInfo
	}

	return report, nil
}

// getSystemInfo collects system hardware information
func (r *ReportGenerator) getSystemInfo() (SystemInfo, error) {
	info := SystemInfo{}

	// Get host information
	hostInfo, err := host.Info()
	if err == nil {
		info.Hostname = hostInfo.Hostname
		info.Platform = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
		info.KernelVersion = hostInfo.KernelVersion

		// Format uptime
		uptime := time.Duration(hostInfo.Uptime) * time.Second
		days := int(uptime.Hours() / 24)
		hours := int(uptime.Hours()) % 24
		minutes := int(uptime.Minutes()) % 60

		if days > 0 {
			info.Uptime = fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
		} else {
			info.Uptime = fmt.Sprintf("%d hours, %d minutes", hours, minutes)
		}
	}

	// Get CPU information
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.CPUModel = cpuInfo[0].ModelName
		info.CPUCores = len(cpuInfo)
	}

	// Get architecture
	info.Architecture = runtime.GOARCH

	// Get memory information
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		totalGB := float64(memInfo.Total) / (1024 * 1024 * 1024)
		info.TotalMemory = fmt.Sprintf("%.2f GB", totalGB)
	}

	// Get disk information
	diskInfo, err := disk.Usage("/")
	if err == nil {
		totalGB := float64(diskInfo.Total) / (1024 * 1024 * 1024)
		info.TotalDisk = fmt.Sprintf("%.2f GB", totalGB)
	}

	return info, nil
}

// getNetworkInfo collects network information
func (r *ReportGenerator) getNetworkInfo() (NetworkInfo, error) {
	info := NetworkInfo{
		Interfaces: []NetworkInterface{},
	}

	// Get network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return info, err
	}

	for _, iface := range interfaces {
		// Skip loopback interfaces
		if strings.Contains(iface.Name, "lo") {
			continue
		}

		netInterface := NetworkInterface{
			Name:       iface.Name,
			MACAddress: iface.HardwareAddr,
		}

		// Get IP addresses
		if len(iface.Addrs) > 0 {
			for _, addr := range iface.Addrs {
				if strings.Contains(addr.Addr, "127.0.0.1") {
					continue
				}
				netInterface.IPAddress = addr.Addr
				break
			}
		}

		// Get interface status based on the interface name
		// This is a simplified approach since we can't directly use net.FlagUp
		netInterface.Status = "UNKNOWN"

		// Check if the interface has an IP address
		if netInterface.IPAddress != "" {
			netInterface.Status = "UP"
		} else {
			netInterface.Status = "DOWN"
		}

		info.Interfaces = append(info.Interfaces, netInterface)
	}

	return info, nil
}

// getSoftwareInfo collects software information
func (r *ReportGenerator) getSoftwareInfo() (SoftwareInfo, error) {
	info := SoftwareInfo{
		OS:          runtime.GOOS,
		GoVersion:   runtime.Version(),
		PackageInfo: make(map[string]string),
	}

	// Get shell version
	shellCmd := exec.Command("bash", "--version")
	shellOutput, err := shellCmd.Output()
	if err == nil {
		shellVersionLines := strings.Split(string(shellOutput), "\n")
		if len(shellVersionLines) > 0 {
			info.ShellVersion = strings.TrimSpace(shellVersionLines[0])
		}
	}

	// Get package manager information based on OS
	switch info.OS {
	case "linux":
		// Try apt (Debian/Ubuntu)
		aptCmd := exec.Command("apt", "--version")
		aptOutput, err := aptCmd.Output()
		if err == nil {
			aptVersionLines := strings.Split(string(aptOutput), "\n")
			if len(aptVersionLines) > 0 {
				info.PackageInfo["apt"] = strings.TrimSpace(aptVersionLines[0])
			}
		}

		// Try dnf (Fedora/RHEL)
		dnfCmd := exec.Command("dnf", "--version")
		dnfOutput, err := dnfCmd.Output()
		if err == nil {
			dnfVersionLines := strings.Split(string(dnfOutput), "\n")
			if len(dnfVersionLines) > 0 {
				info.PackageInfo["dnf"] = strings.TrimSpace(dnfVersionLines[0])
			}
		}
	case "darwin":
		// Try brew (macOS)
		brewCmd := exec.Command("brew", "--version")
		brewOutput, err := brewCmd.Output()
		if err == nil {
			brewVersionLines := strings.Split(string(brewOutput), "\n")
			if len(brewVersionLines) > 0 {
				info.PackageInfo["brew"] = strings.TrimSpace(brewVersionLines[0])
			}
		}
	}

	return info, nil
}

// FormatSystemReport formats a system report for display
func FormatSystemReport(report *SystemReport) string {
	var sb strings.Builder
	boxWidth := 60

	// Format header
	headerText := fmt.Sprintf(" System Report (%s) ", report.Timestamp.Format("2006-01-02 15:04:05"))
	sb.WriteString("╭" + padCenter(headerText, boxWidth-2, "─") + "╮\n")

	// Format system information
	sb.WriteString("│ " + padCenter("System Information", boxWidth-4, " ") + " │\n")
	sb.WriteString("├" + strings.Repeat("─", boxWidth-2) + "┤\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Hostname: %s", report.SystemInfo.Hostname), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Platform: %s", report.SystemInfo.Platform), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Architecture: %s", report.SystemInfo.Architecture), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Kernel: %s", report.SystemInfo.KernelVersion), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("CPU: %s", report.SystemInfo.CPUModel), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("CPU Cores: %d", report.SystemInfo.CPUCores), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Memory: %s", report.SystemInfo.TotalMemory), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Disk: %s", report.SystemInfo.TotalDisk), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Uptime: %s", report.SystemInfo.Uptime), boxWidth-4) + " │\n")

	// Format network information
	sb.WriteString("├" + strings.Repeat("─", boxWidth-2) + "┤\n")
	sb.WriteString("│ " + padCenter("Network Information", boxWidth-4, " ") + " │\n")
	sb.WriteString("├" + strings.Repeat("─", boxWidth-2) + "┤\n")

	if len(report.NetworkInfo.Interfaces) == 0 {
		sb.WriteString("│ " + padRight("No network interfaces found", boxWidth-4) + " │\n")
	} else {
		for _, iface := range report.NetworkInfo.Interfaces {
			sb.WriteString("│ " + padRight(fmt.Sprintf("Interface: %s", iface.Name), boxWidth-4) + " │\n")
			sb.WriteString("│   " + padRight(fmt.Sprintf("IP: %s", iface.IPAddress), boxWidth-6) + " │\n")
			sb.WriteString("│   " + padRight(fmt.Sprintf("MAC: %s", iface.MACAddress), boxWidth-6) + " │\n")
			sb.WriteString("│   " + padRight(fmt.Sprintf("Status: %s", iface.Status), boxWidth-6) + " │\n")
		}
	}

	// Format software information
	sb.WriteString("├" + strings.Repeat("─", boxWidth-2) + "┤\n")
	sb.WriteString("│ " + padCenter("Software Information", boxWidth-4, " ") + " │\n")
	sb.WriteString("├" + strings.Repeat("─", boxWidth-2) + "┤\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("OS: %s", report.SoftwareInfo.OS), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Go Version: %s", report.SoftwareInfo.GoVersion), boxWidth-4) + " │\n")
	sb.WriteString("│ " + padRight(fmt.Sprintf("Shell: %s", truncateString(report.SoftwareInfo.ShellVersion, boxWidth-12)), boxWidth-4) + " │\n")

	// Package information
	if len(report.SoftwareInfo.PackageInfo) > 0 {
		sb.WriteString("│ " + padRight("Package Managers:", boxWidth-4) + " │\n")
		for name, version := range report.SoftwareInfo.PackageInfo {
			sb.WriteString("│   " + padRight(fmt.Sprintf("%s: %s", name, truncateString(version, boxWidth-12-len(name))), boxWidth-6) + " │\n")
		}
	}

	sb.WriteString("╰" + strings.Repeat("─", boxWidth-2) + "╯\n")

	return sb.String()
}
