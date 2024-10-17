package utils

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// InitConfig initializes the application configuration
func InitConfig() error {
	viper.SetConfigName("config")    // Name of the config file (without extension)
	viper.SetConfigType("yaml")      // Type of the config file
	viper.AddConfigPath("./configs") // Path to the directory containing the config file

	viper.AutomaticEnv() // Read environment variables to override values from the config file

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	log.Printf("Using config file: %s", viper.ConfigFileUsed())

	// Optionally, you can validate that required configuration fields are set
	if err := validateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}

	return nil
}

// validateConfig checks if required configurations are present
func validateConfig() error {
	requiredKeys := []string{
		"database.host",
		"database.port",
		"database.user",
		"database.password",
		"database.dbname",
		"database.sslmode",
	}

	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			return fmt.Errorf("missing required configuration: %s", key)
		}
	}
	return nil
}
