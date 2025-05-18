package connect

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// DefaultChunkSize is the default size of each chunk (5MB)
	DefaultChunkSize = 5 * 1024 * 1024

	// MaxChunkSize is the maximum size of each chunk (10MB)
	MaxChunkSize = 10 * 1024 * 1024

	// MinChunkSize is the minimum size of each chunk (1MB)
	MinChunkSize = 1 * 1024 * 1024

	// DefaultUploadTimeout is the default timeout for uploads (1 hour)
	DefaultUploadTimeout = 1 * time.Hour

	// DefaultDownloadTimeout is the default timeout for downloads (1 hour)
	DefaultDownloadTimeout = 1 * time.Hour
)

// ChunkInfo represents information about a file chunk
type ChunkInfo struct {
	ChunkID     int    `json:"chunk_id"`
	ChunkSize   int64  `json:"chunk_size"`
	ChunkOffset int64  `json:"chunk_offset"`
	ChunkHash   string `json:"chunk_hash,omitempty"`
}

// UploadInfo represents information about a file upload
type UploadInfo struct {
	UploadID    string      `json:"upload_id"`
	Filename    string      `json:"filename"`
	FileSize    int64       `json:"file_size"`
	ChunkSize   int64       `json:"chunk_size"`
	TotalChunks int         `json:"total_chunks"`
	Chunks      []ChunkInfo `json:"chunks,omitempty"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time,omitempty"`
	Status      string      `json:"status"` // "pending", "in_progress", "completed", "failed"
	TempPath    string      `json:"-"`      // Path to temporary file (not exposed in JSON)
}

// DownloadInfo represents information about a file download
type DownloadInfo struct {
	DownloadID  string      `json:"download_id"`
	Filename    string      `json:"filename"`
	FileSize    int64       `json:"file_size"`
	ChunkSize   int64       `json:"chunk_size"`
	TotalChunks int         `json:"total_chunks"`
	Chunks      []ChunkInfo `json:"chunks,omitempty"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time,omitempty"`
	Status      string      `json:"status"` // "pending", "in_progress", "completed", "failed"
	FilePath    string      `json:"-"`      // Path to the file (not exposed in JSON)
}

// ChunkedTransferManager manages chunked file transfers
type ChunkedTransferManager struct {
	uploadsMutex   sync.RWMutex
	uploads        map[string]*UploadInfo
	downloadsMutex sync.RWMutex
	downloads      map[string]*DownloadInfo
	tempDir        string
	downloadPath   string
	chunkSize      int64
}

// NewChunkedTransferManager creates a new chunked transfer manager
func NewChunkedTransferManager(downloadPath string, chunkSize int64) (*ChunkedTransferManager, error) {
	// Create a temporary directory for uploads
	tempDir, err := os.MkdirTemp("", "lumo-connect-uploads-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Set default download path if not provided
	if downloadPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			downloadPath = filepath.Join(homeDir, "Downloads")
		} else {
			downloadPath = "."
		}
	}

	// Create the download directory if it doesn't exist
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	// Set default chunk size if not provided
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	} else if chunkSize < MinChunkSize {
		chunkSize = MinChunkSize
	} else if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}

	return &ChunkedTransferManager{
		uploads:      make(map[string]*UploadInfo),
		downloads:    make(map[string]*DownloadInfo),
		tempDir:      tempDir,
		downloadPath: downloadPath,
		chunkSize:    chunkSize,
	}, nil
}

// Cleanup cleans up temporary files and directories
func (m *ChunkedTransferManager) Cleanup() error {
	// Remove the temporary directory
	return os.RemoveAll(m.tempDir)
}

// generateID generates a random ID for uploads and downloads
func generateID() (string, error) {
	// Generate 16 random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// InitUpload initializes a file upload
func (m *ChunkedTransferManager) InitUpload(filename string, fileSize int64) (*UploadInfo, error) {
	// Generate a unique upload ID
	uploadID, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload ID: %w", err)
	}

	// Calculate the number of chunks
	totalChunks := int((fileSize + m.chunkSize - 1) / m.chunkSize)

	// Create a temporary file for the upload
	tempPath := filepath.Join(m.tempDir, uploadID)
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Preallocate the file if possible
	if err := tempFile.Truncate(fileSize); err != nil {
		log.Printf("Warning: Failed to preallocate file: %v", err)
	}

	// Create upload info
	uploadInfo := &UploadInfo{
		UploadID:    uploadID,
		Filename:    filepath.Base(filename),
		FileSize:    fileSize,
		ChunkSize:   m.chunkSize,
		TotalChunks: totalChunks,
		Chunks:      make([]ChunkInfo, totalChunks),
		StartTime:   time.Now(),
		Status:      "pending",
		TempPath:    tempPath,
	}

	// Initialize chunk info
	for i := 0; i < totalChunks; i++ {
		offset := int64(i) * m.chunkSize
		size := m.chunkSize
		if offset+size > fileSize {
			size = fileSize - offset
		}
		uploadInfo.Chunks[i] = ChunkInfo{
			ChunkID:     i,
			ChunkSize:   size,
			ChunkOffset: offset,
		}
	}

	// Store the upload info
	m.uploadsMutex.Lock()
	m.uploads[uploadID] = uploadInfo
	m.uploadsMutex.Unlock()

	return uploadInfo, nil
}

// UploadChunk uploads a chunk of a file
func (m *ChunkedTransferManager) UploadChunk(uploadID string, chunkID int, data []byte) error {
	// Get the upload info
	m.uploadsMutex.RLock()
	uploadInfo, ok := m.uploads[uploadID]
	m.uploadsMutex.RUnlock()
	if !ok {
		return fmt.Errorf("upload not found: %s", uploadID)
	}

	// Check if the chunk ID is valid
	if chunkID < 0 || chunkID >= uploadInfo.TotalChunks {
		return fmt.Errorf("invalid chunk ID: %d", chunkID)
	}

	// Check if the chunk size is valid
	expectedSize := uploadInfo.Chunks[chunkID].ChunkSize
	if int64(len(data)) != expectedSize {
		return fmt.Errorf("invalid chunk size: expected %d, got %d", expectedSize, len(data))
	}

	// Open the temporary file
	file, err := os.OpenFile(uploadInfo.TempPath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open temporary file: %w", err)
	}
	defer file.Close()

	// Write the chunk to the file
	offset := uploadInfo.Chunks[chunkID].ChunkOffset
	if _, err := file.WriteAt(data, offset); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	// Update the upload status
	m.uploadsMutex.Lock()
	uploadInfo.Status = "in_progress"
	uploadInfo.Chunks[chunkID].ChunkHash = "uploaded" // We could calculate a hash here
	m.uploadsMutex.Unlock()

	return nil
}

// CompleteUpload completes a file upload
func (m *ChunkedTransferManager) CompleteUpload(uploadID string) (string, error) {
	// Get the upload info
	m.uploadsMutex.RLock()
	uploadInfo, ok := m.uploads[uploadID]
	m.uploadsMutex.RUnlock()
	if !ok {
		return "", fmt.Errorf("upload not found: %s", uploadID)
	}

	// Check if all chunks have been uploaded
	for _, chunk := range uploadInfo.Chunks {
		if chunk.ChunkHash == "" {
			return "", fmt.Errorf("not all chunks have been uploaded")
		}
	}

	// Create timestamp
	timestamp := time.Now().Format("20060102_150405")

	// Create filename with timestamp
	ext := filepath.Ext(uploadInfo.Filename)
	name := uploadInfo.Filename[:len(uploadInfo.Filename)-len(ext)]
	newFilename := fmt.Sprintf("%s_%s%s", name, timestamp, ext)

	// Create full path
	filePath := filepath.Join(m.downloadPath, newFilename)

	// Move the temporary file to the download directory
	if err := os.Rename(uploadInfo.TempPath, filePath); err != nil {
		// If rename fails (e.g., across different filesystems), try copy
		if err := copyFile(uploadInfo.TempPath, filePath); err != nil {
			return "", fmt.Errorf("failed to move file: %w", err)
		}
		// Remove the temporary file
		os.Remove(uploadInfo.TempPath)
	}

	// Update the upload status
	m.uploadsMutex.Lock()
	uploadInfo.Status = "completed"
	uploadInfo.EndTime = time.Now()
	m.uploadsMutex.Unlock()

	return filePath, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy the file
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
