package desktop

import (
	"fmt"
	"sync"

	"github.com/agnath18K/lumo/dbus/common"
	"github.com/agnath18K/lumo/internal/core"
)

// Factory implements the core.DesktopFactory interface
type Factory struct {
	// environments is a map of desktop environment names to their implementations
	environments map[string]core.DesktopEnvironment
	// mutex protects the environments map
	mutex sync.RWMutex
}

// NewFactory creates a new desktop environment factory
func NewFactory() *Factory {
	return &Factory{
		environments: make(map[string]core.DesktopEnvironment),
	}
}

// RegisterEnvironment registers a desktop environment with the factory
func (f *Factory) RegisterEnvironment(env core.DesktopEnvironment) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.environments[env.Name()] = env
}

// DetectEnvironment detects the current desktop environment
func (f *Factory) DetectEnvironment() (core.DesktopEnvironment, error) {
	// Detect the current desktop environment
	desktopName := common.DetectDesktopEnvironment()

	// Try to get the environment by name
	env, err := f.GetEnvironment(desktopName)
	if err == nil && env.IsAvailable() {
		return env, nil
	}

	// If the detected environment is not available, try to find any available environment
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	for _, env := range f.environments {
		if env.IsAvailable() {
			return env, nil
		}
	}

	return nil, fmt.Errorf("no available desktop environment found")
}

// GetEnvironment gets a specific desktop environment by name
func (f *Factory) GetEnvironment(name string) (core.DesktopEnvironment, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	// Try exact match first
	if env, ok := f.environments[name]; ok {
		return env, nil
	}

	// Try case-insensitive match
	for envName, env := range f.environments {
		if envName == name {
			return env, nil
		}
	}

	return nil, fmt.Errorf("desktop environment not found: %s", name)
}

// ListAvailableEnvironments lists all available desktop environments
func (f *Factory) ListAvailableEnvironments() []string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	var available []string
	for name, env := range f.environments {
		if env.IsAvailable() {
			available = append(available, name)
		}
	}

	return available
}
