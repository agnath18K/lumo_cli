package discovery

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

const (
	// ServiceName is the name of the service to advertise
	ServiceName = "_lumo-connect._tcp"
	// ServiceDomain is the domain to advertise the service on
	ServiceDomain = "local."
	// TTL is the time-to-live for the service advertisement
	TTL = 60
)

// MDNSDiscoverer implements the Discoverer interface using mDNS
type MDNSDiscoverer struct {
	server       *mdns.Server
	entries      map[string]Service
	entriesMutex sync.RWMutex
	callbacks    []func(Service)
	callbackMux  sync.RWMutex
}

// NewMDNSDiscoverer creates a new MDNSDiscoverer
func NewMDNSDiscoverer() *MDNSDiscoverer {
	return &MDNSDiscoverer{
		entries:   make(map[string]Service),
		callbacks: make([]func(Service), 0),
	}
}

// Start starts the discovery service
func (d *MDNSDiscoverer) Start(ctx context.Context) error {
	// Start a goroutine to periodically browse for services
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		// Do an initial browse
		d.browseServices(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				d.browseServices(ctx)
			}
		}
	}()

	return nil
}

// Stop stops the discovery service
func (d *MDNSDiscoverer) Stop() error {
	if d.server != nil {
		d.server.Shutdown()
		d.server = nil
	}
	return nil
}

// Advertise advertises a service
func (d *MDNSDiscoverer) Advertise(ctx context.Context, name string, port int, info map[string]string) error {
	// Stop any existing advertisement
	if d.server != nil {
		d.server.Shutdown()
		d.server = nil
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// Create TXT record
	txt := make([]string, 0, len(info))
	for k, v := range info {
		txt = append(txt, fmt.Sprintf("%s=%s", k, v))
	}

	// Create service
	service, err := mdns.NewMDNSService(
		name,          // Instance name
		ServiceName,   // Service type
		ServiceDomain, // Domain
		hostname,      // Host name
		port,          // Port
		nil,           // IPs (nil = all interfaces)
		txt,           // TXT records
	)
	if err != nil {
		return fmt.Errorf("failed to create mDNS service: %w", err)
	}

	// Create server
	server, err := mdns.NewServer(&mdns.Config{
		Zone: service,
	})
	if err != nil {
		return fmt.Errorf("failed to create mDNS server: %w", err)
	}

	d.server = server
	return nil
}

// StopAdvertising stops advertising a service
func (d *MDNSDiscoverer) StopAdvertising() error {
	if d.server != nil {
		d.server.Shutdown()
		d.server = nil
	}
	return nil
}

// Browse returns a list of discovered services
func (d *MDNSDiscoverer) Browse(ctx context.Context, serviceType string) ([]Service, error) {
	d.browseServices(ctx)

	d.entriesMutex.RLock()
	defer d.entriesMutex.RUnlock()

	services := make([]Service, 0, len(d.entries))
	for _, service := range d.entries {
		services = append(services, service)
	}

	return services, nil
}

// AddServiceCallback adds a callback function that is called when a service is discovered
func (d *MDNSDiscoverer) AddServiceCallback(callback func(Service)) {
	d.callbackMux.Lock()
	defer d.callbackMux.Unlock()
	d.callbacks = append(d.callbacks, callback)
}

// browseServices browses for services
func (d *MDNSDiscoverer) browseServices(ctx context.Context) {
	// Create a channel for results
	entriesCh := make(chan *mdns.ServiceEntry, 10)
	go func() {
		for entry := range entriesCh {
			// Get port from service
			port := entry.Port

			// Parse info from TXT records
			info := make(map[string]string)
			for _, txt := range entry.InfoFields {
				parts := strings.SplitN(txt, "=", 2)
				if len(parts) == 2 {
					info[parts[0]] = parts[1]
				}
			}

			// Create service
			service := Service{
				ID:       entry.Name,
				Name:     entry.Host,
				Host:     entry.Host,
				IP:       entry.AddrV4.String(),
				Port:     port,
				Info:     info,
				LastSeen: time.Now(),
			}

			// Add to entries
			d.entriesMutex.Lock()
			d.entries[service.ID] = service
			d.entriesMutex.Unlock()

			// Call callbacks
			d.callbackMux.RLock()
			for _, callback := range d.callbacks {
				callback(service)
			}
			d.callbackMux.RUnlock()
		}
	}()

	// Start browsing
	params := mdns.DefaultParams(ServiceName)
	params.Entries = entriesCh
	params.Timeout = 5 * time.Second

	// Use a context-aware wrapper for mdns.Query
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create a channel to signal when the query is done
	done := make(chan struct{})
	go func() {
		if err := mdns.Query(params); err != nil {
			log.Printf("Error browsing for services: %v", err)
		}
		close(done)
	}()

	// Wait for either the query to complete or the context to be canceled
	select {
	case <-done:
		// Query completed normally
	case <-queryCtx.Done():
		// Context was canceled
		log.Printf("Service discovery canceled: %v", queryCtx.Err())
	}
	close(entriesCh)
}
