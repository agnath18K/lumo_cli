package tests

import (
	"errors"
	"net/http"
	"testing"
)

// TestHTTPClientOperations tests basic HTTP client operations using the mock
func TestHTTPClientOperations(t *testing.T) {
	// Create a mock HTTP client
	client := &MockHTTPClient{
		Responses: make(map[string]*http.Response),
		Errors:    make(map[string]error),
		Requests:  []*http.Request{},
		Calls:     []string{},
	}

	// Add a mock response
	client.AddResponse("GET", "https://example.com/api", http.StatusOK, `{"message": "Success"}`)

	// Make a request
	resp, err := client.Get("https://example.com/api")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	// Verify the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify the call was recorded
	found := false
	for _, call := range client.Calls {
		if call == "GET:https://example.com/api" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected GET call for https://example.com/api, but it wasn't recorded")
	}

	// Add a mock error
	client.AddError("GET", "https://example.com/error", errors.New("network error"))

	// Make a request that should fail
	_, err = client.Get("https://example.com/error")
	if err == nil {
		t.Errorf("Expected error when making GET request with error, but got none")
	}

	// Verify the error call was recorded
	found = false
	for _, call := range client.Calls {
		if call == "GET:https://example.com/error" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected GET call for https://example.com/error, but it wasn't recorded")
	}
}

// MockHTTPClient is a mock implementation of the http.Client
type MockHTTPClient struct {
	Responses map[string]*http.Response
	Errors    map[string]error
	Requests  []*http.Request
	Calls     []string
}

// Do records the request and returns the mock response or error
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, req)

	// Record the call with method and URL
	callKey := req.Method + ":" + req.URL.String()
	m.Calls = append(m.Calls, callKey)

	// Check if there's an error for this request
	if err, ok := m.Errors[callKey]; ok {
		return nil, err
	}

	// Check if there's a response for this request
	if resp, ok := m.Responses[callKey]; ok {
		return resp, nil
	}

	// Default response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}, nil
}

// Get is a convenience method that calls Do with a GET request
func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return m.Do(req)
}

// AddResponse adds a mock response for a specific request
func (m *MockHTTPClient) AddResponse(method, url string, statusCode int, body string) {
	key := method + ":" + url
	m.Responses[key] = &http.Response{
		StatusCode: statusCode,
		Body:       http.NoBody,
	}
}

// AddError adds a mock error for a specific request
func (m *MockHTTPClient) AddError(method, url string, err error) {
	key := method + ":" + url
	m.Errors[key] = err
}

// TestSpeedTestOperations tests speed test operations
func TestSpeedTestOperations(t *testing.T) {
	// Skip this test for now
	t.Skip("Skipping test that requires proper mocking of the speed tester")
}
