package speedtest

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/agnath18/lumo/pkg/utils"
)

// SpeedTestResult represents the result of a speed test
type SpeedTestResult struct {
	DownloadSpeed float64 // in Mbps
	UploadSpeed   float64 // in Mbps
	Latency       int     // in ms
	ISP           string
	Server        string
	Timestamp     time.Time
}

// SpeedTester handles internet speed testing
type SpeedTester struct {
	client *http.Client
}

// NewSpeedTester creates a new speed tester
func NewSpeedTester() *SpeedTester {
	return &SpeedTester{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RunTest performs a complete speed test (download, upload, and latency)
func (s *SpeedTester) RunTest(ctx context.Context) (*SpeedTestResult, error) {
	// Check if there's an internet connection
	if !utils.CheckInternetConnectivity() {
		return nil, fmt.Errorf("no internet connection detected")
	}

	// Create a result object
	result := &SpeedTestResult{
		Timestamp: time.Now(),
	}

	// Get the best server
	server, err := s.findBestServer()
	if err != nil {
		return nil, fmt.Errorf("failed to find test server: %w", err)
	}
	result.Server = server.Name
	result.ISP = s.detectISP()

	// Measure latency
	latency, err := s.measureLatency(server)
	if err != nil {
		return nil, fmt.Errorf("failed to measure latency: %w", err)
	}
	result.Latency = latency

	// Measure download speed
	downloadSpeed, err := s.measureDownloadSpeed(server)
	if err != nil {
		return nil, fmt.Errorf("failed to measure download speed: %w", err)
	}
	result.DownloadSpeed = downloadSpeed

	// Measure upload speed
	uploadSpeed, err := s.measureUploadSpeed(server)
	if err != nil {
		return nil, fmt.Errorf("failed to measure upload speed: %w", err)
	}
	result.UploadSpeed = uploadSpeed

	return result, nil
}

// RunDownloadTest performs only a download speed test
func (s *SpeedTester) RunDownloadTest(ctx context.Context) (*SpeedTestResult, error) {
	// Check if there's an internet connection
	if !utils.CheckInternetConnectivity() {
		return nil, fmt.Errorf("no internet connection detected")
	}

	// Create a result object
	result := &SpeedTestResult{
		Timestamp: time.Now(),
	}

	// Get the best server
	server, err := s.findBestServer()
	if err != nil {
		return nil, fmt.Errorf("failed to find test server: %w", err)
	}
	result.Server = server.Name
	result.ISP = s.detectISP()

	// Measure download speed
	downloadSpeed, err := s.measureDownloadSpeed(server)
	if err != nil {
		return nil, fmt.Errorf("failed to measure download speed: %w", err)
	}
	result.DownloadSpeed = downloadSpeed

	return result, nil
}

// RunUploadTest performs only an upload speed test
func (s *SpeedTester) RunUploadTest(ctx context.Context) (*SpeedTestResult, error) {
	// Check if there's an internet connection
	if !utils.CheckInternetConnectivity() {
		return nil, fmt.Errorf("no internet connection detected")
	}

	// Create a result object
	result := &SpeedTestResult{
		Timestamp: time.Now(),
	}

	// Get the best server
	server, err := s.findBestServer()
	if err != nil {
		return nil, fmt.Errorf("failed to find test server: %w", err)
	}
	result.Server = server.Name
	result.ISP = s.detectISP()

	// Measure upload speed
	uploadSpeed, err := s.measureUploadSpeed(server)
	if err != nil {
		return nil, fmt.Errorf("failed to measure upload speed: %w", err)
	}
	result.UploadSpeed = uploadSpeed

	return result, nil
}

// FormatResult formats the speed test result as a string
func (s *SpeedTester) FormatResult(result *SpeedTestResult) string {
	var sb strings.Builder

	// Get terminal width for proper formatting
	termWidth := utils.GetTerminalWidth()
	if termWidth < 60 {
		termWidth = 60
	}
	if termWidth > 100 {
		termWidth = 100
	}

	// Create a box with the results
	title := "🚀 Internet Speed Test Results"

	sb.WriteString("╭" + strings.Repeat("─", termWidth-2) + "╮\n")
	sb.WriteString("│ " + utils.PadCenter(title, termWidth-4, " ") + " │\n")
	sb.WriteString("├" + strings.Repeat("─", termWidth-2) + "┤\n")

	// Add download speed
	sb.WriteString("│ " + utils.PadRight("Download:", 12) + " " + utils.PadRight(fmt.Sprintf("%.2f Mbps", result.DownloadSpeed), 12) + " " + utils.PadRight("", termWidth-30) + " │\n")

	// Add upload speed
	sb.WriteString("│ " + utils.PadRight("Upload:", 12) + " " + utils.PadRight(fmt.Sprintf("%.2f Mbps", result.UploadSpeed), 12) + " " + utils.PadRight("", termWidth-30) + " │\n")

	// Add latency if available
	if result.Latency > 0 {
		latencyRating := rateLatency(result.Latency)
		sb.WriteString("│ " + utils.PadRight("Latency:", 12) + " " + utils.PadRight(fmt.Sprintf("%d ms", result.Latency), 12) + " " + utils.PadRight(latencyRating, termWidth-30) + " │\n")
	}

	// Add ISP and server information
	if result.ISP != "" {
		sb.WriteString("│ " + utils.PadRight("ISP:", 12) + " " + utils.PadRight(result.ISP, termWidth-16) + " │\n")
	}
	if result.Server != "" {
		sb.WriteString("│ " + utils.PadRight("Server:", 12) + " " + utils.PadRight(result.Server, termWidth-16) + " │\n")
	}

	// Add timestamp
	sb.WriteString("│ " + utils.PadRight("Time:", 12) + " " + utils.PadRight(result.Timestamp.Format("2006-01-02 15:04:05"), termWidth-16) + " │\n")

	// Add a connection quality rating
	rating := rateConnection(result.DownloadSpeed, result.UploadSpeed, result.Latency)
	sb.WriteString("├" + strings.Repeat("─", termWidth-2) + "┤\n")
	sb.WriteString("│ " + utils.PadCenter(fmt.Sprintf("Connection Quality: %s", rating), termWidth-4, " ") + " │\n")
	sb.WriteString("╰" + strings.Repeat("─", termWidth-2) + "╯\n")

	return sb.String()
}

// Server represents a speed test server
type Server struct {
	Name     string
	URL      string
	Distance float64
}

// findBestServer finds the best server for speed testing
func (s *SpeedTester) findBestServer() (*Server, error) {
	// In a real implementation, this would query a list of servers
	// and select the best one based on ping time and distance
	// For now, we'll return a mock server
	return &Server{
		Name:     "Speedtest.net Server (New York)",
		URL:      "https://speedtest.net",
		Distance: 10.5,
	}, nil
}

// detectISP attempts to detect the user's ISP
func (s *SpeedTester) detectISP() string {
	// In a real implementation, this would query an API to get the ISP
	// For now, we'll return a mock ISP
	return "Example ISP"
}

// measureLatency measures the latency to the server
func (s *SpeedTester) measureLatency(server *Server) (int, error) {
	// In a real implementation, this would send ping requests to the server
	// For now, we'll return a mock latency
	return 25, nil
}

// measureDownloadSpeed measures the download speed
func (s *SpeedTester) measureDownloadSpeed(server *Server) (float64, error) {
	// In a real implementation, this would download files from the server
	// and measure the speed
	// For now, we'll return a mock download speed
	return 95.67, nil
}

// measureUploadSpeed measures the upload speed
func (s *SpeedTester) measureUploadSpeed(server *Server) (float64, error) {
	// In a real implementation, this would upload files to the server
	// and measure the speed
	// For now, we'll return a mock upload speed
	return 35.42, nil
}

// createSpeedBar creates a visual bar representing the speed
func createSpeedBar(speed float64, maxWidth int) string {
	// Determine the number of bars to show based on the speed
	// We'll use a logarithmic scale to make it more visually appealing
	maxBars := maxWidth - 2 // Leave space for brackets

	// Scale: 0-10 Mbps: 1-3 bars, 10-50 Mbps: 4-7 bars, 50-100 Mbps: 8-12 bars, 100+ Mbps: 13+ bars
	var numBars int
	switch {
	case speed < 1:
		numBars = 1
	case speed < 10:
		numBars = int(math.Max(1, math.Min(float64(maxBars/4), 1+speed/3)))
	case speed < 50:
		numBars = int(math.Max(float64(maxBars/4), math.Min(float64(maxBars/2), float64(maxBars/4)+speed/10)))
	case speed < 100:
		numBars = int(math.Max(float64(maxBars/2), math.Min(float64(3*maxBars/4), float64(maxBars/2)+speed/20)))
	default:
		numBars = int(math.Max(float64(3*maxBars/4), math.Min(float64(maxBars), float64(3*maxBars/4)+speed/100)))
	}

	if numBars > maxBars {
		numBars = maxBars
	}

	// Create the bar
	bar := "["
	bar += strings.Repeat("█", numBars)
	bar += strings.Repeat(" ", maxBars-numBars)
	bar += "]"

	return bar
}

// rateLatency rates the latency
func rateLatency(latency int) string {
	switch {
	case latency < 20:
		return "Excellent 🏆"
	case latency < 50:
		return "Very Good 👍"
	case latency < 100:
		return "Good ✓"
	case latency < 150:
		return "Average ⚠️"
	default:
		return "Poor ⚠️"
	}
}

// rateConnection rates the overall connection quality
func rateConnection(downloadSpeed, uploadSpeed float64, latency int) string {
	// Calculate a score based on download, upload, and latency
	// Download has the highest weight, followed by upload and latency
	downloadScore := math.Min(100, downloadSpeed) / 100 * 50     // 50% weight
	uploadScore := math.Min(50, uploadSpeed) / 50 * 30           // 30% weight
	latencyScore := math.Max(0, 200-float64(latency)) / 200 * 20 // 20% weight

	totalScore := downloadScore + uploadScore + latencyScore

	// Rate based on the total score
	switch {
	case totalScore >= 90:
		return "Excellent 🏆"
	case totalScore >= 75:
		return "Very Good 👍"
	case totalScore >= 60:
		return "Good ✓"
	case totalScore >= 40:
		return "Average ⚠️"
	default:
		return "Poor ⚠️"
	}
}
