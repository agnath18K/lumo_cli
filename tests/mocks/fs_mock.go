package mocks

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MockFileSystem is a mock implementation of file system operations
type MockFileSystem struct {
	Files       map[string][]byte
	Directories map[string]bool
	Errors      map[string]error
	Calls       []string
}

// NewMockFileSystem creates a new mock file system
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files:       make(map[string][]byte),
		Directories: make(map[string]bool),
		Errors:      make(map[string]error),
		Calls:       []string{},
	}
}

// ReadFile reads a file from the mock file system
func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	m.Calls = append(m.Calls, "ReadFile:"+path)
	
	if err, ok := m.Errors["ReadFile:"+path]; ok {
		return nil, err
	}
	
	if data, ok := m.Files[path]; ok {
		return data, nil
	}
	
	return nil, os.ErrNotExist
}

// WriteFile writes a file to the mock file system
func (m *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	m.Calls = append(m.Calls, "WriteFile:"+path)
	
	if err, ok := m.Errors["WriteFile:"+path]; ok {
		return err
	}
	
	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		m.MkdirAll(dir, 0755)
	}
	
	m.Files[path] = data
	return nil
}

// MkdirAll creates a directory and all parent directories in the mock file system
func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	m.Calls = append(m.Calls, "MkdirAll:"+path)
	
	if err, ok := m.Errors["MkdirAll:"+path]; ok {
		return err
	}
	
	m.Directories[path] = true
	return nil
}

// Remove removes a file or directory from the mock file system
func (m *MockFileSystem) Remove(path string) error {
	m.Calls = append(m.Calls, "Remove:"+path)
	
	if err, ok := m.Errors["Remove:"+path]; ok {
		return err
	}
	
	delete(m.Files, path)
	delete(m.Directories, path)
	return nil
}

// Stat returns file info for a file in the mock file system
func (m *MockFileSystem) Stat(path string) (os.FileInfo, error) {
	m.Calls = append(m.Calls, "Stat:"+path)
	
	if err, ok := m.Errors["Stat:"+path]; ok {
		return nil, err
	}
	
	if _, ok := m.Files[path]; ok {
		return &mockFileInfo{
			name:    filepath.Base(path),
			size:    int64(len(m.Files[path])),
			mode:    0644,
			modTime: time.Now(),
			isDir:   false,
		}, nil
	}
	
	if _, ok := m.Directories[path]; ok {
		return &mockFileInfo{
			name:    filepath.Base(path),
			size:    0,
			mode:    0755 | os.ModeDir,
			modTime: time.Now(),
			isDir:   true,
		}, nil
	}
	
	return nil, os.ErrNotExist
}

// Open opens a file in the mock file system
func (m *MockFileSystem) Open(path string) (io.ReadCloser, error) {
	m.Calls = append(m.Calls, "Open:"+path)
	
	if err, ok := m.Errors["Open:"+path]; ok {
		return nil, err
	}
	
	if data, ok := m.Files[path]; ok {
		return io.NopCloser(strings.NewReader(string(data))), nil
	}
	
	return nil, os.ErrNotExist
}

// Create creates a file in the mock file system
func (m *MockFileSystem) Create(path string) (io.WriteCloser, error) {
	m.Calls = append(m.Calls, "Create:"+path)
	
	if err, ok := m.Errors["Create:"+path]; ok {
		return nil, err
	}
	
	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		m.MkdirAll(dir, 0755)
	}
	
	writer := &mockWriter{
		fs:   m,
		path: path,
		buf:  &strings.Builder{},
	}
	
	return writer, nil
}

// SetError sets an error to be returned for a specific operation and path
func (m *MockFileSystem) SetError(operation, path string, err error) {
	m.Errors[operation+":"+path] = err
}

// Reset clears all files, directories, errors, and calls
func (m *MockFileSystem) Reset() {
	m.Files = make(map[string][]byte)
	m.Directories = make(map[string]bool)
	m.Errors = make(map[string]error)
	m.Calls = []string{}
}

// mockFileInfo implements os.FileInfo for the mock file system
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

// mockWriter implements io.WriteCloser for the mock file system
type mockWriter struct {
	fs   *MockFileSystem
	path string
	buf  *strings.Builder
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *mockWriter) Close() error {
	w.fs.Files[w.path] = []byte(w.buf.String())
	return nil
}

// FileNotFoundError returns a "file not found" error
func FileNotFoundError(path string) error {
	return errors.New("file not found: " + path)
}

// PermissionDeniedError returns a "permission denied" error
func PermissionDeniedError(path string) error {
	return errors.New("permission denied: " + path)
}

// DiskFullError returns a "disk full" error
func DiskFullError() error {
	return errors.New("no space left on device")
}
