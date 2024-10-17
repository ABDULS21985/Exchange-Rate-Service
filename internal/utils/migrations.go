// internal/utils/migrations.go

package utils

import (
	"log"

	"github.com/abduls21985/exchange-rate-service/internal/models"
)

// RunMigrations uses GORM to auto-migrate database schemas
func RunMigrations() error {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Currency{},
		&models.ExchangeRate{},
	); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully.")
	return nil
}
