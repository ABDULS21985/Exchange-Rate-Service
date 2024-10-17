package services

import (
	"fmt"
	"time"

	"github.com/abduls21985/exchange-rate-service/internal/models"
	"github.com/abduls21985/exchange-rate-service/internal/repositories"
	"github.com/abduls21985/exchange-rate-service/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// UserService interface defines user-related operations
type UserService interface {
	RegisterUser(user *models.User) (*models.User, error)
	AuthenticateUser(username, password string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	InitiatePasswordReset(email string) (*models.User, error)
	ResetPassword(token, newPassword string) error
}

type userService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// RegisterUser registers a new user
func (s *userService) RegisterUser(user *models.User) (*models.User, error) {
	// Check if the username or email already exists
	existingUser, err := s.repo.FindUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already taken")
	}

	existingUser, err = s.repo.FindUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = string(hashedPassword)

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save the user to the database
	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// AuthenticateUser authenticates a user by username and password
func (s *userService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.repo.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

// UpdateUser updates user details
func (s *userService) UpdateUser(user *models.User) error {
	return s.repo.UpdateUser(user)
}

// InitiatePasswordReset initiates a password reset
func (s *userService) InitiatePasswordReset(email string) (*models.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Generate a reset token and expiry time
	user.ResetToken = utils.GenerateRandomToken()
	user.ResetTokenExpiry = time.Now().Add(1 * time.Hour)

	// Update the user with reset token details
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user with reset token: %v", err)
	}

	return user, nil
}

// ResetPassword resets the user's password if the token is valid
func (s *userService) ResetPassword(token, newPassword string) error {
	user, err := s.repo.FindUserByResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	if time.Now().After(user.ResetTokenExpiry) {
		return fmt.Errorf("reset token has expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}

	return s.repo.UpdateUser(user)
}

// GetUserByUsername retrieves a user by username
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.repo.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}
