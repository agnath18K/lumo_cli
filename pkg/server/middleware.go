package server

import (
	"context"
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
		// Skip authentication for certain endpoints
		if !s.config.EnableAuth || isExemptPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
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
	}

	// Check if the path is in the exempt list
	for _, exemptPath := range exemptPaths {
		if path == exemptPath {
			return true
		}
	}

	// Check if the path is a static file
	if strings.HasPrefix(path, "/static/") ||
		path == "/" ||
		path == "/index.html" ||
		path == "/favicon.ico" ||
		strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") {
		return true
	}

	return false
}

// getUsernameFromContext gets the username from the request context
func getUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(userContextKey).(string)
	return username, ok
}
