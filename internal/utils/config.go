package utils

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// InitConfig initializes the application configuration
func InitConfig() error {
	// Set the configuration file name and type
	viper.SetConfigName("config")    // Name of the config file (without extension)
	viper.SetConfigType("yaml")      // Type of the config file
	viper.AddConfigPath("./configs") // Path to the directory containing the config file
	viper.AddConfigPath(".")         // Optionally, add the current directory as a fallback

	// Read environment variables to override values from the config file
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Log the config file being used
	log.Printf("Using config file: %s", viper.ConfigFileUsed())

	// Validate the required configuration fields
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
		"exchange_rate_api_url", // Add required configuration for the exchange rate API URL
		"exchange_rate_app_id",  // Add required configuration for the exchange rate API App ID
	}

	// Iterate over required keys and check if they are set
	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			return fmt.Errorf("missing required configuration: %s", key)
		}
	}
	return nil
}
