package services

import (
	"fmt"
	"os"
	"time"

	"github.com/abduls21985/exchange-rate-service/internal/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// AuthService interface defines authentication-related operations
type AuthService interface {
	AuthenticateUser(username, password string) (*models.User, error)
	GenerateJWT(username string) (string, error)
}

type authService struct {
	userRepo UserService
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo UserService) AuthService {
	return &authService{userRepo}
}

// AuthenticateUser verifies the username and password for login
func (s *authService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

// GenerateJWT generates a JWT token for the authenticated user
func (s *authService) GenerateJWT(username string) (string, error) {
	// Get the secret key from environment variables
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT secret key not configured")
	}

	// Define the token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create JWT claims
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: expirationTime.Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}

	return tokenString, nil
}
