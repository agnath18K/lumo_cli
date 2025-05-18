package discovery

import (
	"context"
	"time"
)

// Service represents a discovered service
type Service struct {
	// ID is a unique identifier for the service
	ID string
	// Name is the human-readable name of the service
	Name string
	// Host is the hostname of the service
	Host string
	// IP is the IP address of the service
	IP string
	// Port is the port number of the service
	Port int
	// Info contains additional information about the service
	Info map[string]string
	// LastSeen is the time when the service was last seen
	LastSeen time.Time
}

// Discoverer is the interface for service discovery
type Discoverer interface {
	// Start starts the discovery service
	Start(ctx context.Context) error
	// Stop stops the discovery service
	Stop() error
	// Advertise advertises a service
	Advertise(ctx context.Context, name string, port int, info map[string]string) error
	// StopAdvertising stops advertising a service
	StopAdvertising() error
	// Browse returns a list of discovered services
	Browse(ctx context.Context, serviceType string) ([]Service, error)
	// AddServiceCallback adds a callback function that is called when a service is discovered
	AddServiceCallback(callback func(Service))
}

// NewDiscoverer creates a new discoverer based on the available implementation
func NewDiscoverer() Discoverer {
	// Currently we only support mDNS
	return NewMDNSDiscoverer()
}
