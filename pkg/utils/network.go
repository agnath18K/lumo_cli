package utils

import (
	"fmt"
	"net"
	"strconv"
)

// IsPortAvailable checks if a port is available for use
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// FindAvailablePort finds an available port starting from the given port
// If the given port is available, it returns that port
// Otherwise, it increments the port number until it finds an available port
// It will try up to maxAttempts times before giving up
func FindAvailablePort(startPort int, maxAttempts int) (int, error) {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		port := startPort + attempt
		if IsPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("could not find an available port after %d attempts", maxAttempts)
}

// GetPortRangeMessage returns a message suggesting alternative ports based on the component
func GetPortRangeMessage(component string) string {
	switch component {
	case "server":
		return "Try using a port in the range 7500-7599 for the Lumo server."
	case "connect":
		return "Try using a port in the range 8000-8099 for Lumo connect."
	default:
		return "Try using a different port."
	}
}
