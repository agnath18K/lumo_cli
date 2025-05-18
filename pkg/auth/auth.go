package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultTokenExpiration is the default expiration time for JWT tokens (24 hours)
	DefaultTokenExpiration = 24 * time.Hour

	// DefaultRefreshTokenExpiration is the default expiration time for refresh tokens (7 days)
	DefaultRefreshTokenExpiration = 7 * 24 * time.Hour

	// DefaultCredentialsFile is the default file name for storing credentials
	DefaultCredentialsFile = "credentials.json"

	// DefaultBcryptCost is the default cost for bcrypt password hashing
	DefaultBcryptCost = 12
)

var (
	// ErrInvalidCredentials is returned when the provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound is returned when the user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrTokenExpired is returned when the token has expired
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
)

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Credentials represents the user credentials
type Credentials struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CredentialsStore represents the credentials store
type CredentialsStore struct {
	Credentials []Credentials `json:"credentials"`
	UpdatedAt   string        `json:"updated_at"`
}

// Authenticator handles authentication-related functionality
type Authenticator struct {
	jwtSecret         []byte
	credentialsPath   string
	tokenExpiration   time.Duration
	refreshExpiration time.Duration
}

// NewAuthenticator creates a new authenticator instance
func NewAuthenticator(jwtSecret string, credentialsDir string) (*Authenticator, error) {
	// If JWT secret is empty, generate a random one
	var secretBytes []byte
	if jwtSecret == "" {
		// Generate a random 32-byte secret
		secretBytes = make([]byte, 32)
		if _, err := rand.Read(secretBytes); err != nil {
			return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
		}
	} else {
		secretBytes = []byte(jwtSecret)
	}

	// Create credentials directory if it doesn't exist
	if err := os.MkdirAll(credentialsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create credentials directory: %w", err)
	}

	return &Authenticator{
		jwtSecret:         secretBytes,
		credentialsPath:   filepath.Join(credentialsDir, DefaultCredentialsFile),
		tokenExpiration:   DefaultTokenExpiration,
		refreshExpiration: DefaultRefreshTokenExpiration,
	}, nil
}

// GenerateToken generates a JWT token for the given username
func (a *Authenticator) GenerateToken(username string) (string, error) {
	// Create the claims
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lumo",
			Subject:   username,
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	return token.SignedString(a.jwtSecret)
}

// GenerateRefreshToken generates a refresh token for the given username
func (a *Authenticator) GenerateRefreshToken(username string) (string, error) {
	// Create the claims with longer expiration
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.refreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lumo",
			Subject:   username,
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	return token.SignedString(a.jwtSecret)
}

// ValidateToken validates the given token and returns the claims
func (a *Authenticator) ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Get the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	// Hash the password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultBcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies the given password against the hash
func VerifyPassword(password, hash string) bool {
	// Compare the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecureToken generates a secure random token
func GenerateSecureToken(length int) (string, error) {
	// Generate random bytes
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Encode to base64
	return base64.URLEncoding.EncodeToString(b), nil
}
