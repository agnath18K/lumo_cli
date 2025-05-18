package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed web
var webFS embed.FS

// GetWebFileSystem returns a http.FileSystem for the embedded web files
func GetWebFileSystem() http.FileSystem {
	// Get the embedded web/static directory
	webStatic, err := fs.Sub(webFS, "web/static")
	if err != nil {
		// If there's an error, return an empty file system
		return http.FS(embed.FS{})
	}
	return http.FS(webStatic)
}

// HasEmbeddedWebFiles returns true if the embedded web files exist
func HasEmbeddedWebFiles() bool {
	// Try to open the index.html file
	f, err := webFS.Open("web/static/index.html")
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error opening embedded web file: %v\n", err)
		return false
	}
	f.Close()
	return true
}
