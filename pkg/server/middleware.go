package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/agnath18K/lumo/pkg/auth"
)

// contextKey is a custom type for context keys
type contextKey string

// userContextKey is the key for the username in the request context
const userContextKey contextKey = "username"

// AuthMiddleware is a middleware that checks for a valid JWT token
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request path for debugging
		log.Printf("Request path: %s", r.URL.Path)

		// Skip authentication for certain endpoints
		if !s.config.EnableAuth || isExemptPath(r.URL.Path) {
			log.Printf("Path %s is exempt from authentication", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Authorization header required for path: %s", r.URL.Path)
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header has the correct format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization format, expected 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := s.authenticator.ValidateToken(tokenString)
		if err != nil {
			if err == auth.ErrTokenExpired {
				http.Error(w, "Token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			return
		}

		// Add the username to the request context
		ctx := context.WithValue(r.Context(), userContextKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isExemptPath returns true if the path is exempt from authentication
func isExemptPath(path string) bool {
	// List of paths that don't require authentication
	exemptPaths := []string{
		"/ping",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/api/v1/status",
		// Connect endpoints don't require authentication
		"/api/v1/connect/ws",
		"/api/v1/connect/upload/init",
		"/api/v1/connect/upload/chunk",
		"/api/v1/connect/upload/complete",
		"/api/v1/connect/discover",
		"/api/v1/connect/start-server",
		"/api/v1/connect/connect-to-peer",
		"/api/v1/connect/disconnect",
		"/api/v1/connect/send-file",
	}

	// Check if the path is in the exempt list
	for _, exemptPath := range exemptPaths {
		if path == exemptPath {
			log.Printf("Path %s is in the exempt list", path)
			return true
		}
	}

	// Check if the path is a static file or connect page
	if strings.HasPrefix(path, "/static/") ||
		path == "/" ||
		path == "/index.html" ||
		path == "/favicon.ico" ||
		strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") {
		log.Printf("Path %s is a static file", path)
		return true
	}

	// Check if it's a connect page
	if path == "/connect/" ||
		path == "/connect/index.html" ||
		strings.HasPrefix(path, "/connect/") {
		log.Printf("Path %s is a connect page", path)
		return true
	}

	log.Printf("Path %s is NOT exempt from authentication", path)
	return false
}

// getUsernameFromContext gets the username from the request context
func getUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(userContextKey).(string)
	return username, ok
}
