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

	// Fetch the exchange rates using the service layer
	rates, err := c.Service.FetchExchangeRates(currencyCode, timestamp)
	if err != nil {
		log.Printf("Error fetching exchange rates: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	// Check if the rates slice is empty
	if len(rates) == 0 {
		utils.JSONResponse(w, map[string]string{"error": "No exchange rates found for the given criteria"}, http.StatusNotFound)
		return
	}

	// Aggregate the rates into a map
	ratesMap := make(map[string]float64)
	for _, rate := range rates {
		if rate.Currency.Code != "" {
			ratesMap[rate.Currency.Code] = rate.Rate
		}
	}

	// Format the response data
	responseData := map[string]interface{}{
		"data": map[string]interface{}{
			"timestamp": rates[0].Timestamp.Unix(),  // Assuming all rates have the same timestamp
			"base":      rates[0].BaseCurrency.Code, // Assuming all rates share the same base currency
			"rates":     ratesMap,
		},
		"status": "Exchange rates fetched successfully",
	}

	// Return the formatted response
	utils.JSONResponse(w, responseData, http.StatusOK)
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
func (c *ExchangeRateController) CountExchangeRates(w http.ResponseWriter, r *http.Request) {
	count, err := c.Service.CountExchangeRates()
	if err != nil {
		log.Printf("Error counting exchange rates: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]int{"count": count}, http.StatusOK)
}

// GetHistoricalExchangeRates handles GET /api/exchange-rates/historical
func (c *ExchangeRateController) GetHistoricalExchangeRates(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	currencyCode := query.Get("currency")
	startStr := query.Get("start_date")
	endStr := query.Get("end_date")

	var startDate, endDate int64
	var err error

	if startStr != "" {
		startDate, err = strconv.ParseInt(startStr, 10, 64)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"error": "Invalid start date format"}, http.StatusBadRequest)
			return
		}
	}

	if endStr != "" {
		endDate, err = strconv.ParseInt(endStr, 10, 64)
		if err != nil {
			utils.JSONResponse(w, map[string]string{"error": "Invalid end date format"}, http.StatusBadRequest)
			return
		}
	}

	rates, err := c.Service.GetHistoricalExchangeRates(currencyCode, startDate, endDate)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	response := make(map[string]float64)
	for _, rate := range rates {
		currency, err := c.Service.GetAllCurrencies()
		if err != nil {
			utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
			return
		}
		for _, cur := range currency {
			if cur.ID == rate.CurrencyID {
				response[cur.Code] = rate.Rate
				break
			}
		}
	}

	utils.JSONResponse(w, map[string]interface{}{
		"data": map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
			"base":       "USD",
			"rates":      response,
		},
		"status": "Historical exchange rates fetched successfully",
	}, http.StatusOK)
}

// internal/controllers/exchange_rates_controller.go

// ConvertCurrency handles POST /api/exchange-rates/convert
func (c *ExchangeRateController) ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	var request struct {
		FromCurrency string  `json:"from_currency"`
		ToCurrency   string  `json:"to_currency"`
		Amount       float64 `json:"amount"`
	}

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.JSONResponse(w, map[string]string{"error": "Invalid request payload"}, http.StatusBadRequest)
		return
	}

	// Validate request parameters
	if request.FromCurrency == "" || request.ToCurrency == "" || request.Amount <= 0 {
		utils.JSONResponse(w, map[string]string{"error": "Missing or invalid parameters"}, http.StatusBadRequest)
		return
	}

	// Perform currency conversion
	convertedAmount, err := c.Service.ConvertCurrency(request.FromCurrency, request.ToCurrency, request.Amount)
	if err != nil {
		log.Printf("Error converting currency: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	// Respond with the converted amount
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"from_currency":    request.FromCurrency,
			"to_currency":      request.ToCurrency,
			"original_amount":  request.Amount,
			"converted_amount": convertedAmount,
		},
		"status": "Currency converted successfully",
	}

	utils.JSONResponse(w, response, http.StatusOK)
}

// internal/controllers/exchange_rates_controller.go

// ConvertRatesToBaseCurrency handles GET /api/exchange-rates/base-convert
func (c *ExchangeRateController) ConvertRatesToBaseCurrency(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base")
	if baseCurrency == "" {
		utils.JSONResponse(w, map[string]string{"error": "Base currency is required"}, http.StatusBadRequest)
		return
	}

	// Fetch the exchange rates
	rates, err := c.Service.FetchExchangeRates("", 0)
	if err != nil {
		log.Printf("Error fetching exchange rates: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Failed to fetch exchange rates"}, http.StatusInternalServerError)
		return
	}

	// Convert rates to the specified base currency
	convertedRates, err := c.Service.ConvertToBaseCurrency(rates, baseCurrency)
	if err != nil {
		log.Printf("Error converting to base currency: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Failed to convert exchange rates to base currency"}, http.StatusInternalServerError)
		return
	}

	// Respond with the converted rates
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"base":   baseCurrency,
			"rates":  convertedRates,
			"status": "Exchange rates converted to base currency successfully",
		},
	}

	utils.JSONResponse(w, response, http.StatusOK)
}

// internal/controllers/exchange_rates_controller.go

// GetExchangeRateCount handles GET /api/exchange-rates/count
func (c *ExchangeRateController) GetExchangeRateCount(w http.ResponseWriter, r *http.Request) {
	count, err := c.Service.CountExchangeRates()
	if err != nil {
		log.Printf("Error counting exchange rates: %v", err)
		utils.JSONResponse(w, map[string]string{"error": "Failed to count exchange rates"}, http.StatusInternalServerError)
		return
	}

	// Respond with the count of exchange rates
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"count":  count,
			"status": "Exchange rates counted successfully",
		},
	}

	utils.JSONResponse(w, response, http.StatusOK)
}
