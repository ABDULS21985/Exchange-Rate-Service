// package repositories

package repositories

import (
	"github.com/abduls21985/exchange-rate-service/internal/models"
	"gorm.io/gorm"
)

// UserRepository interface defines the methods for user operations
type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByResetToken(token string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// CreateUser inserts a new user into the database
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// FindUserByUsername retrieves a user by their username
func (r *userRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// UpdateUser updates user details in the database
func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// FindUserByEmail retrieves a user by their email
func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindUserByResetToken retrieves a user by their reset token
func (r *userRepository) FindUserByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("reset_token = ?", token).First(&user).Error
	return &user, err
}
