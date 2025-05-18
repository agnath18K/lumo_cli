package server

import (
	"encoding/json"
	"net/http"

	"github.com/agnath18K/lumo/pkg/auth"
)

// handleLogin handles the /api/v1/auth/login endpoint
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	if err := s.authenticator.Authenticate(req.Username, req.Password); err != nil {
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidCredentials {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Authentication error", http.StatusInternalServerError)
		}
		return
	}

	// Generate tokens
	token, err := s.authenticator.GenerateToken(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := s.authenticator.GenerateRefreshToken(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Create the response
	resp := LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Username:     req.Username,
		ExpiresIn:    s.config.TokenExpirationHours * 3600, // Convert hours to seconds
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handleRefreshToken handles the /api/v1/auth/refresh endpoint
func (s *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Validate the refresh token
	claims, err := s.authenticator.ValidateToken(req.RefreshToken)
	if err != nil {
		if err == auth.ErrTokenExpired {
			http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		} else {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		}
		return
	}

	// Generate new tokens
	token, err := s.authenticator.GenerateToken(claims.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := s.authenticator.GenerateRefreshToken(claims.Username)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Create the response
	resp := LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Username:     claims.Username,
		ExpiresIn:    s.config.TokenExpirationHours * 3600, // Convert hours to seconds
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handleChangePassword handles the /api/v1/auth/change-password endpoint
func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the username from the context
	username, ok := getUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Current password and new password are required", http.StatusBadRequest)
		return
	}

	// Authenticate with the current password
	if err := s.authenticator.Authenticate(username, req.CurrentPassword); err != nil {
		if err == auth.ErrInvalidCredentials {
			http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		} else {
			http.Error(w, "Authentication error", http.StatusInternalServerError)
		}
		return
	}

	// Update the password
	if err := s.authenticator.UpdatePassword(username, req.NewPassword); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Password updated successfully"}`))
}
