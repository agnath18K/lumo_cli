package mocks

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// MockHTTPClient is a mock implementation of the http.Client
type MockHTTPClient struct {
	Responses map[string]*http.Response
	Errors    map[string]error
	Requests  []*http.Request
	Calls     []string
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		Responses: make(map[string]*http.Response),
		Errors:    make(map[string]error),
		Requests:  []*http.Request{},
		Calls:     []string{},
	}
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
		Body:       io.NopCloser(strings.NewReader("{}")),
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

// Post is a convenience method that calls Do with a POST request
func (m *MockHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return m.Do(req)
}

// AddResponse adds a mock response for a specific request
func (m *MockHTTPClient) AddResponse(method, url string, statusCode int, body string) {
	key := method + ":" + url
	m.Responses[key] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// AddError adds a mock error for a specific request
func (m *MockHTTPClient) AddError(method, url string, err error) {
	key := method + ":" + url
	m.Errors[key] = err
}

// Reset clears all responses, errors, requests, and calls
func (m *MockHTTPClient) Reset() {
	m.Responses = make(map[string]*http.Response)
	m.Errors = make(map[string]error)
	m.Requests = []*http.Request{}
	m.Calls = []string{}
}

// MockSpeedTester is a mock implementation of the speed tester
type MockSpeedTester struct {
	DownloadSpeed float64
	UploadSpeed   float64
	Latency       int
	ISP           string
	Server        string
	Error         error
	Calls         []string
}

// NewMockSpeedTester creates a new mock speed tester with default values
func NewMockSpeedTester() *MockSpeedTester {
	return &MockSpeedTester{
		DownloadSpeed: 95.67,
		UploadSpeed:   25.34,
		Latency:       25,
		ISP:           "Mock ISP",
		Server:        "Mock Server",
		Calls:         []string{},
	}
}

// RunTest records the call and returns the mock result or error
func (m *MockSpeedTester) RunTest(ctx context.Context) (*SpeedTestResult, error) {
	m.Calls = append(m.Calls, "RunTest")

	if m.Error != nil {
		return nil, m.Error
	}

	return &SpeedTestResult{
		DownloadSpeed: m.DownloadSpeed,
		UploadSpeed:   m.UploadSpeed,
		Latency:       m.Latency,
		ISP:           m.ISP,
		Server:        m.Server,
		Timestamp:     time.Now(),
	}, nil
}

// RunDownloadTest records the call and returns the mock result or error
func (m *MockSpeedTester) RunDownloadTest(ctx context.Context) (*SpeedTestResult, error) {
	m.Calls = append(m.Calls, "RunDownloadTest")

	if m.Error != nil {
		return nil, m.Error
	}

	return &SpeedTestResult{
		DownloadSpeed: m.DownloadSpeed,
		ISP:           m.ISP,
		Server:        m.Server,
		Timestamp:     time.Now(),
	}, nil
}

// RunUploadTest records the call and returns the mock result or error
func (m *MockSpeedTester) RunUploadTest(ctx context.Context) (*SpeedTestResult, error) {
	m.Calls = append(m.Calls, "RunUploadTest")

	if m.Error != nil {
		return nil, m.Error
	}

	return &SpeedTestResult{
		UploadSpeed: m.UploadSpeed,
		ISP:         m.ISP,
		Server:      m.Server,
		Timestamp:   time.Now(),
	}, nil
}

// RunLatencyTest records the call and returns the mock result or error
func (m *MockSpeedTester) RunLatencyTest(ctx context.Context) (*SpeedTestResult, error) {
	m.Calls = append(m.Calls, "RunLatencyTest")

	if m.Error != nil {
		return nil, m.Error
	}

	return &SpeedTestResult{
		Latency:   m.Latency,
		ISP:       m.ISP,
		Server:    m.Server,
		Timestamp: time.Now(),
	}, nil
}

// SpeedTestResult represents the result of a speed test
type SpeedTestResult struct {
	DownloadSpeed float64
	UploadSpeed   float64
	Latency       int
	ISP           string
	Server        string
	Timestamp     time.Time
}

// NetworkError returns a network error
func NetworkError(message string) error {
	return errors.New("network error: " + message)
}

// TimeoutError returns a timeout error
func TimeoutError() error {
	return errors.New("timeout error: operation timed out")
}
