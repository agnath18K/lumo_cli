package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/agnath18K/lumo/pkg/connect"
)

// Global chunked transfer manager
var (
	chunkedTransferManager     *connect.ChunkedTransferManager
	chunkedTransferManagerOnce sync.Once
)

// getChunkedTransferManager returns the global chunked transfer manager
func (s *Server) getChunkedTransferManager() *connect.ChunkedTransferManager {
	chunkedTransferManagerOnce.Do(func() {
		// Get the download path from the config
		homeDir, err := os.UserHomeDir()
		downloadPath := filepath.Join(homeDir, "Downloads")
		if err != nil {
			log.Printf("Error getting user home directory: %v", err)
			downloadPath = "."
		}

		// Create the chunked transfer manager
		manager, err := connect.NewChunkedTransferManager(downloadPath, connect.DefaultChunkSize)
		if err != nil {
			log.Printf("Error creating chunked transfer manager: %v", err)
			return
		}
		chunkedTransferManager = manager
	})
	return chunkedTransferManager
}

// InitUploadRequest represents a request to initialize a file upload
type InitUploadRequest struct {
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
}

// InitUploadResponse represents a response to initialize a file upload
type InitUploadResponse struct {
	Success   bool                `json:"success"`
	Error     string              `json:"error,omitempty"`
	UploadID  string              `json:"upload_id,omitempty"`
	ChunkSize int64               `json:"chunk_size,omitempty"`
	Chunks    []connect.ChunkInfo `json:"chunks,omitempty"`
}

// UploadChunkResponse represents a response to upload a chunk
type UploadChunkResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	ChunkID int    `json:"chunk_id,omitempty"`
}

// CompleteUploadResponse represents a response to complete an upload
type CompleteUploadResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

// handleInitUpload handles the /api/v1/connect/upload/init endpoint
func (s *Server) handleInitUpload(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request InitUploadRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if request.Filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}
	if request.FileSize <= 0 {
		http.Error(w, "File size must be greater than 0", http.StatusBadRequest)
		return
	}

	// Get the chunked transfer manager
	manager := s.getChunkedTransferManager()
	if manager == nil {
		http.Error(w, "Chunked transfer manager not available", http.StatusInternalServerError)
		return
	}

	// Initialize the upload
	uploadInfo, err := manager.InitUpload(request.Filename, request.FileSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize upload: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the response
	response := InitUploadResponse{
		Success:   true,
		UploadID:  uploadInfo.UploadID,
		ChunkSize: uploadInfo.ChunkSize,
		Chunks:    uploadInfo.Chunks,
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleUploadChunk handles the /api/v1/connect/upload/chunk endpoint
func (s *Server) handleUploadChunk(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the upload ID from the query parameters
	uploadID := r.URL.Query().Get("upload_id")
	if uploadID == "" {
		http.Error(w, "Upload ID is required", http.StatusBadRequest)
		return
	}

	// Get the chunk ID from the query parameters
	chunkIDStr := r.URL.Query().Get("chunk_id")
	if chunkIDStr == "" {
		http.Error(w, "Chunk ID is required", http.StatusBadRequest)
		return
	}
	chunkID, err := strconv.Atoi(chunkIDStr)
	if err != nil {
		http.Error(w, "Invalid chunk ID", http.StatusBadRequest)
		return
	}

	// Get the chunked transfer manager
	manager := s.getChunkedTransferManager()
	if manager == nil {
		http.Error(w, "Chunked transfer manager not available", http.StatusInternalServerError)
		return
	}

	// Read the chunk data
	chunkData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read chunk data: %v", err), http.StatusInternalServerError)
		return
	}

	// Upload the chunk
	err = manager.UploadChunk(uploadID, chunkID, chunkData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload chunk: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the response
	response := UploadChunkResponse{
		Success: true,
		ChunkID: chunkID,
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleCompleteUpload handles the /api/v1/connect/upload/complete endpoint
func (s *Server) handleCompleteUpload(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the upload ID from the query parameters
	uploadID := r.URL.Query().Get("upload_id")
	if uploadID == "" {
		http.Error(w, "Upload ID is required", http.StatusBadRequest)
		return
	}

	// Get the chunked transfer manager
	manager := s.getChunkedTransferManager()
	if manager == nil {
		http.Error(w, "Chunked transfer manager not available", http.StatusInternalServerError)
		return
	}

	// Complete the upload
	filePath, err := manager.CompleteUpload(uploadID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to complete upload: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the response
	response := CompleteUploadResponse{
		Success:  true,
		FilePath: filePath,
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
