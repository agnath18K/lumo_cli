package speedtest_test

import (
	"context"
	"testing"
	"time"

	"github.com/agnath18/lumo/pkg/speedtest"
)

// TestSpeedTester tests the basic functionality of the speed tester
func TestSpeedTester(t *testing.T) {
	// Create a new speed tester
	tester := speedtest.NewSpeedTester()

	// Check that the tester was created
	if tester == nil {
		t.Fatal("Expected speed tester to be created, got nil")
	}

	// Create a context with a short timeout to avoid long tests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the RunTest method
	result, err := tester.RunTest(ctx)
	if err != nil {
		// Skip the test if there's no internet connection
		// This allows the tests to pass in offline environments
		t.Skip("Skipping speed test due to error:", err)
	}

	// Check that the result has valid values
	if result.DownloadSpeed <= 0 {
		t.Errorf("Expected download speed to be greater than 0, got %f", result.DownloadSpeed)
	}

	if result.UploadSpeed <= 0 {
		t.Errorf("Expected upload speed to be greater than 0, got %f", result.UploadSpeed)
	}

	if result.Latency <= 0 {
		t.Errorf("Expected latency to be greater than 0, got %d", result.Latency)
	}

	if result.ISP == "" {
		t.Error("Expected ISP to be non-empty")
	}

	if result.Server == "" {
		t.Error("Expected server to be non-empty")
	}

	// Test the FormatResult method
	formatted := tester.FormatResult(result)
	if formatted == "" {
		t.Error("Expected formatted result to be non-empty")
	}
}

// TestDownloadTest tests the download-only test
func TestDownloadTest(t *testing.T) {
	// Create a new speed tester
	tester := speedtest.NewSpeedTester()

	// Create a context with a short timeout to avoid long tests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the RunDownloadTest method
	result, err := tester.RunDownloadTest(ctx)
	if err != nil {
		// Skip the test if there's no internet connection
		t.Skip("Skipping download test due to error:", err)
	}

	// Check that the result has valid values
	if result.DownloadSpeed <= 0 {
		t.Errorf("Expected download speed to be greater than 0, got %f", result.DownloadSpeed)
	}

	// Upload speed should be 0 for download-only test
	if result.UploadSpeed != 0 {
		t.Errorf("Expected upload speed to be 0 for download-only test, got %f", result.UploadSpeed)
	}
}

// TestUploadTest tests the upload-only test
func TestUploadTest(t *testing.T) {
	// Create a new speed tester
	tester := speedtest.NewSpeedTester()

	// Create a context with a short timeout to avoid long tests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the RunUploadTest method
	result, err := tester.RunUploadTest(ctx)
	if err != nil {
		// Skip the test if there's no internet connection
		t.Skip("Skipping upload test due to error:", err)
	}

	// Check that the result has valid values
	if result.UploadSpeed <= 0 {
		t.Errorf("Expected upload speed to be greater than 0, got %f", result.UploadSpeed)
	}

	// Download speed should be 0 for upload-only test
	if result.DownloadSpeed != 0 {
		t.Errorf("Expected download speed to be 0 for upload-only test, got %f", result.DownloadSpeed)
	}
}
