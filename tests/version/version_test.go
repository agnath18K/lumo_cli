package version_test

import (
	"strings"
	"testing"

	"github.com/agnath18/lumo/pkg/version"
)

// TestGetVersion tests the GetVersion function
func TestGetVersion(t *testing.T) {
	// Get the version string
	versionStr := version.GetVersion()
	
	// Check that the version string is not empty
	if versionStr == "" {
		t.Fatal("Expected non-empty version string")
	}
	
	// Check that the version string contains the version number
	if !strings.Contains(versionStr, version.Version) {
		t.Errorf("Expected version string to contain version number '%s', got: %s", version.Version, versionStr)
	}
	
	// Check that the version string contains the build date
	if !strings.Contains(versionStr, version.BuildDate) {
		t.Errorf("Expected version string to contain build date '%s', got: %s", version.BuildDate, versionStr)
	}
	
	// Check that the version string contains the git commit
	if !strings.Contains(versionStr, version.GitCommit) {
		t.Errorf("Expected version string to contain git commit '%s', got: %s", version.GitCommit, versionStr)
	}
	
	// Check that the version string contains the Go version
	if !strings.Contains(versionStr, version.GoVersion) {
		t.Errorf("Expected version string to contain Go version '%s', got: %s", version.GoVersion, versionStr)
	}
}

// TestGetShortVersion tests the GetShortVersion function
func TestGetShortVersion(t *testing.T) {
	// Get the short version string
	shortVersion := version.GetShortVersion()
	
	// Check that the short version string is not empty
	if shortVersion == "" {
		t.Fatal("Expected non-empty short version string")
	}
	
	// Check that the short version string is equal to the version number
	if shortVersion != version.Version {
		t.Errorf("Expected short version to be '%s', got: %s", version.Version, shortVersion)
	}
}
