package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// InitializeCredentialsStore initializes the credentials store
func (a *Authenticator) InitializeCredentialsStore() error {
	// Check if credentials file exists
	if _, err := os.Stat(a.credentialsPath); os.IsNotExist(err) {
		// Create a new credentials store
		store := CredentialsStore{
			Credentials: []Credentials{},
			UpdatedAt:   time.Now().Format(time.RFC3339),
		}

		// Save the credentials store
		return a.saveCredentialsStore(&store)
	}

	return nil
}

// loadCredentialsStore loads the credentials store from disk
func (a *Authenticator) loadCredentialsStore() (*CredentialsStore, error) {
	// Initialize the credentials store if it doesn't exist
	if err := a.InitializeCredentialsStore(); err != nil {
		return nil, fmt.Errorf("failed to initialize credentials store: %w", err)
	}

	// Read the credentials file
	data, err := os.ReadFile(a.credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Parse the JSON
	var store CredentialsStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return &store, nil
}

// saveCredentialsStore saves the credentials store to disk
func (a *Authenticator) saveCredentialsStore(store *CredentialsStore) error {
	// Update the timestamp
	store.UpdatedAt = time.Now().Format(time.RFC3339)

	// Marshal to JSON
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write to file
	if err := os.WriteFile(a.credentialsPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// AddUser adds a new user to the credentials store
func (a *Authenticator) AddUser(username, password string) error {
	// Load the credentials store
	store, err := a.loadCredentialsStore()
	if err != nil {
		return err
	}

	// Check if user already exists
	for _, cred := range store.Credentials {
		if cred.Username == username {
			return fmt.Errorf("user already exists: %s", username)
		}
	}

	// Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	// Create the new credentials
	now := time.Now().Format(time.RFC3339)
	cred := Credentials{
		Username:     username,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Add to the store
	store.Credentials = append(store.Credentials, cred)

	// Save the store
	return a.saveCredentialsStore(store)
}

// UpdatePassword updates the password for the given user
func (a *Authenticator) UpdatePassword(username, password string) error {
	// Load the credentials store
	store, err := a.loadCredentialsStore()
	if err != nil {
		return err
	}

	// Find the user
	found := false
	for i, cred := range store.Credentials {
		if cred.Username == username {
			// Hash the password
			hash, err := HashPassword(password)
			if err != nil {
				return err
			}

			// Update the credentials
			store.Credentials[i].PasswordHash = hash
			store.Credentials[i].UpdatedAt = time.Now().Format(time.RFC3339)
			found = true
			break
		}
	}

	if !found {
		return ErrUserNotFound
	}

	// Save the store
	return a.saveCredentialsStore(store)
}

// Authenticate authenticates the given username and password
func (a *Authenticator) Authenticate(username, password string) error {
	// Load the credentials store
	store, err := a.loadCredentialsStore()
	if err != nil {
		return err
	}

	// Find the user
	for _, cred := range store.Credentials {
		if cred.Username == username {
			// Verify the password
			if VerifyPassword(password, cred.PasswordHash) {
				return nil
			}
			return ErrInvalidCredentials
		}
	}

	return ErrUserNotFound
}

// HasUsers returns true if there are any users in the credentials store
func (a *Authenticator) HasUsers() (bool, error) {
	// Load the credentials store
	store, err := a.loadCredentialsStore()
	if err != nil {
		return false, err
	}

	return len(store.Credentials) > 0, nil
}

// GetUsers returns a list of usernames in the credentials store
func (a *Authenticator) GetUsers() ([]string, error) {
	// Load the credentials store
	store, err := a.loadCredentialsStore()
	if err != nil {
		return nil, err
	}

	// Extract usernames
	usernames := make([]string, len(store.Credentials))
	for i, cred := range store.Credentials {
		usernames[i] = cred.Username
	}

	return usernames, nil
}

// RemoveUser removes a user from the credentials store
func (a *Authenticator) RemoveUser(username string) error {
	// Load the credentials store
	store, err := a.loadCredentialsStore()
	if err != nil {
		return err
	}

	// Find the user
	found := false
	for i, cred := range store.Credentials {
		if cred.Username == username {
			// Remove the user
			store.Credentials = append(store.Credentials[:i], store.Credentials[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return ErrUserNotFound
	}

	// Save the store
	return a.saveCredentialsStore(store)
}
