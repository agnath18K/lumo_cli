package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ChunkedClient is a client for chunked file transfers
type ChunkedClient struct {
	baseURL     string
	downloadDir string
	chunkSize   int64
	httpClient  *http.Client
}

// NewChunkedClient creates a new chunked client
func NewChunkedClient(baseURL, downloadDir string, chunkSize int64) *ChunkedClient {
	// Set default values
	if baseURL == "" {
		baseURL = "http://localhost:7531"
	}
	if downloadDir == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			downloadDir = filepath.Join(homeDir, "Downloads")
		} else {
			downloadDir = "."
		}
	}
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	} else if chunkSize < MinChunkSize {
		chunkSize = MinChunkSize
	} else if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}

	return &ChunkedClient{
		baseURL:     baseURL,
		downloadDir: downloadDir,
		chunkSize:   chunkSize,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // 30 second timeout for regular requests
		},
	}
}

// UploadFile uploads a file using chunked transfer
func (c *ChunkedClient) UploadFile(filePath string, progressCallback func(int)) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular file")
	}

	// Get base filename
	filename := filepath.Base(filePath)

	// Format file size
	sizeStr := formatFileSize(fileInfo.Size())
	fmt.Printf("\033[1;32m📤 Uploading file: %s (%s)...\033[0m\n", filename, sizeStr)

	// Initialize the upload
	uploadInfo, err := c.initUpload(filename, fileInfo.Size())
	if err != nil {
		return "", fmt.Errorf("failed to initialize upload: %w", err)
	}

	// Calculate total chunks
	totalChunks := uploadInfo.TotalChunks

	// Show progress bar
	fmt.Printf("\033[1;32m[                    ] 0%%\033[0m")
	fmt.Printf("\r")

	// Upload each chunk
	buffer := make([]byte, uploadInfo.ChunkSize)
	for i := 0; i < totalChunks; i++ {
		// Calculate the chunk size
		chunkSize := uploadInfo.ChunkSize
		if i == totalChunks-1 {
			// Last chunk might be smaller
			chunkSize = fileInfo.Size() - int64(i)*uploadInfo.ChunkSize
		}

		// Seek to the correct position
		if _, err := file.Seek(int64(i)*uploadInfo.ChunkSize, 0); err != nil {
			return "", fmt.Errorf("failed to seek file: %w", err)
		}

		// Read the chunk
		n, err := io.ReadFull(file, buffer[:chunkSize])
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return "", fmt.Errorf("failed to read chunk: %w", err)
		}

		// Upload the chunk
		if err := c.uploadChunk(uploadInfo.UploadID, i, buffer[:n]); err != nil {
			return "", fmt.Errorf("failed to upload chunk %d: %w", i, err)
		}

		// Update progress
		progress := (i + 1) * 100 / totalChunks
		if progressCallback != nil {
			progressCallback(progress)
		}

		// Update progress bar
		bars := progress / 5
		spaces := 20 - bars
		fmt.Printf("\033[1;32m[%s%s] %d%%\033[0m", strings.Repeat("=", bars), strings.Repeat(" ", spaces), progress)
		fmt.Printf("\r")
	}

	// Complete the upload
	filePath, err = c.completeUpload(uploadInfo.UploadID)
	if err != nil {
		return "", fmt.Errorf("failed to complete upload: %w", err)
	}

	// Update progress bar to 100%
	fmt.Printf("\033[1;32m[====================] 100%%\033[0m\n")
	fmt.Printf("\033[1;32m📤 File uploaded successfully!\033[0m\n")

	return filePath, nil
}

// initUpload initializes a file upload
func (c *ChunkedClient) initUpload(filename string, fileSize int64) (*UploadInfo, error) {
	// Create the request body
	reqBody := map[string]interface{}{
		"filename":  filename,
		"file_size": fileSize,
	}

	// Convert the request body to JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", c.baseURL+"/api/v1/connect/upload/init", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the content type
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// Read the response body
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned error: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	var respBody struct {
		Success   bool        `json:"success"`
		Error     string      `json:"error,omitempty"`
		UploadID  string      `json:"upload_id,omitempty"`
		ChunkSize int64       `json:"chunk_size,omitempty"`
		Chunks    []ChunkInfo `json:"chunks,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if !respBody.Success {
		return nil, fmt.Errorf("server returned error: %s", respBody.Error)
	}

	// Create the upload info
	uploadInfo := &UploadInfo{
		UploadID:    respBody.UploadID,
		Filename:    filename,
		FileSize:    fileSize,
		ChunkSize:   respBody.ChunkSize,
		TotalChunks: len(respBody.Chunks),
		Chunks:      respBody.Chunks,
		StartTime:   time.Now(),
		Status:      "pending",
	}

	return uploadInfo, nil
}

// uploadChunk uploads a chunk of a file
func (c *ChunkedClient) uploadChunk(uploadID string, chunkID int, data []byte) error {
	// Create the URL with query parameters
	url := fmt.Sprintf("%s/api/v1/connect/upload/chunk?upload_id=%s&chunk_id=%d", c.baseURL, uploadID, chunkID)

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the content type
	req.Header.Set("Content-Type", "application/octet-stream")

	// Create a client with a longer timeout for chunk uploads
	client := &http.Client{
		Timeout: 5 * time.Minute, // 5 minute timeout for chunk uploads
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// Read the response body
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	var respBody struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
		ChunkID int    `json:"chunk_id,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if !respBody.Success {
		return fmt.Errorf("server returned error: %s", respBody.Error)
	}

	return nil
}

// completeUpload completes a file upload
func (c *ChunkedClient) completeUpload(uploadID string) (string, error) {
	// Create the URL with query parameters
	url := fmt.Sprintf("%s/api/v1/connect/upload/complete?upload_id=%s", c.baseURL, uploadID)

	// Create the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// Read the response body
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server returned error: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	var respBody struct {
		Success  bool   `json:"success"`
		Error    string `json:"error,omitempty"`
		FilePath string `json:"file_path,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if !respBody.Success {
		return "", fmt.Errorf("server returned error: %s", respBody.Error)
	}

	return respBody.FilePath, nil
}
