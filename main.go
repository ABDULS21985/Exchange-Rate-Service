package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/abduls21985/exchange-rate-service/internal/models"
	"github.com/abduls21985/exchange-rate-service/internal/repositories"
	"github.com/abduls21985/exchange-rate-service/internal/routes"
	"github.com/abduls21985/exchange-rate-service/internal/services"
	"github.com/abduls21985/exchange-rate-service/internal/utils"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize configuration
	if err := utils.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize database
	if err := utils.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, _ := utils.DB.DB()
		sqlDB.Close()
	}()

	// Run database migrations using GORM
	if err := utils.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize router
	router := mux.NewRouter()

	// Set up all routes using the routes package
	routes.InitializeRoutes(router, utils.DB)

	// Initialize the ExchangeRateService
	exchangeRateRepo := repositories.NewExchangeRateRepository(utils.DB)
	exchangeRateService := services.NewExchangeRateService(exchangeRateRepo)

	// Add a manual trigger endpoint for fetching exchange rates
	router.HandleFunc("/api/manual-fetch", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Manually triggering exchange rate data fetch...")
		data, err := FetchExchangeRateData()
		if err != nil {
			log.Printf("Error fetching data: %v", err)
			http.Error(w, "Failed to fetch exchange rates", http.StatusInternalServerError)
			return
		}

		// Call the service layer to update the exchange rates
		if err := exchangeRateService.AddExchangeRates(*data); err != nil {
			log.Printf("Error inserting data: %v", err)
			http.Error(w, "Failed to update exchange rates", http.StatusInternalServerError)
			return
		}

		log.Println("Exchange rates updated successfully")
		// Include both the status and the fetched data in the response
		response := map[string]interface{}{
			"status": "Exchange rates updated successfully",
			"data":   data,
		}
		jsonResponse(w, response, http.StatusOK)
	}).Methods("GET")

	// Initialize Cron for daily data synchronization
	c := cron.New()
	c.AddFunc("@daily", func() {
		log.Println("Running scheduled daily data synchronization...")
		data, err := FetchExchangeRateData()
		if err != nil {
			log.Printf("Error fetching data: %v", err)
			return
		}
		// Call the service layer to update the exchange rates
		if err := exchangeRateService.AddExchangeRates(*data); err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
		log.Println("Exchange rates updated successfully")
	})
	c.Start()
	defer c.Stop()

	// Start the server
	serverPort := viper.GetString("server.port")
	log.Printf("Server is running on port %s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// FetchExchangeRateData fetches exchange rate data from the external API
func FetchExchangeRateData() (*models.ExchangeRateData, error) {
	// Get the current date in the format YYYY-MM-DD
	date := time.Now().Format("2006-01-02")

	// Get the API URL and app ID from the configuration
	apiUrl := viper.GetString("EXCHANGE_RATE_API_URL")
	appId := viper.GetString("EXCHANGE_RATE_APP_ID")

	// Construct the URL using the date, API URL, and app ID
	url := fmt.Sprintf(apiUrl+"?app_id=%s", date, appId)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data models.ExchangeRateData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Helper function to write JSON responses
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
