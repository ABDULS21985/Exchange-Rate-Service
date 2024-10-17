package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/abduls21985/exchange-rate-service/internal/models"
	"github.com/abduls21985/exchange-rate-service/internal/services"
	"github.com/abduls21985/exchange-rate-service/internal/utils"
)

// ExchangeRateController handles HTTP requests related to exchange rates
type ExchangeRateController struct {
	Service services.ExchangeRateService
}

// NewExchangeRateController creates a new ExchangeRateController
func NewExchangeRateController(service services.ExchangeRateService) *ExchangeRateController {
	return &ExchangeRateController{Service: service}
}

// GetExchangeRates handles GET /api/exchange-rates
func (c *ExchangeRateController) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	currencyCode := query.Get("currency")
	timestampStr := query.Get("timestamp")

	var timestamp int64
	var err error
	if timestampStr != "" {
		timestamp, err = strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"error": "Invalid timestamp format"}, http.StatusBadRequest)
			return
		}
	}

	rates, err := c.Service.FetchExchangeRates(currencyCode, timestamp)
	if err != nil {
		log.Printf("Error fetching exchange rates: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, rates, http.StatusOK)
}

// PostExchangeRates handles POST /api/exchange-rates
func (c *ExchangeRateController) PostExchangeRates(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data models.ExchangeRateData
	if err := json.Unmarshal(body, &data); err != nil {
		utils.JSONResponse(w, map[string]string{"error": "Invalid JSON format"}, http.StatusBadRequest)
		return
	}

	if err := c.Service.AddExchangeRates(data); err != nil {
		log.Printf("Error adding exchange rates: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{"status": "Exchange rates updated successfully"}, http.StatusCreated)
}

// GetCurrencies handles GET /api/currencies
func (c *ExchangeRateController) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := c.Service.GetAllCurrencies()
	if err != nil {
		log.Printf("Error fetching currencies: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, currencies, http.StatusOK)
}

// HealthCheck handles GET /api/health
func (c *ExchangeRateController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, map[string]string{"status": "OK", "message": "Server is running"}, http.StatusOK)
}

// CountExchangeRates handles GET /api/exchange-rates/count
// func (c *ExchangeRateController) CountExchangeRates(w http.ResponseWriter, r *http.Request) {
//     count, err := c.Service.CountExchangeRates()
//     if err != nil {
//         log.Printf("Error counting exchange rates: %v", err)
//         utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
//         return
//     }

//     utils.JSONResponse(w, map[string]int{"count": count}, http.StatusOK)
// }
