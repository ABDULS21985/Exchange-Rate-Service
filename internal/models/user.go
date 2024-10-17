// internal/models/user.go

package models

import "time"

// User represents a user in the system.
type User struct {
	ID               uint   `gorm:"primaryKey"`
	FirstName        string `gorm:"size:100;not null"`
	LastName         string `gorm:"size:100;not null"`
	Email            string `gorm:"size:150;unique;not null"`
	Username         string `gorm:"size:100;unique;not null"`
	PhoneNumber      string `gorm:"size:20"`
	Password         string `gorm:"not null"` // Store hashed password
	ResetToken       string `gorm:"size:255"` // Token for password reset
	ResetTokenExpiry time.Time
	MdaID            string `gorm:"size:50"`
	MDA              string `gorm:"size:150"`
	Role             string `gorm:"size:50"`
	IsOwner          bool   `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
