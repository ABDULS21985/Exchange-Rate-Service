// internal/services/exchange_rate_service.go

package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/abduls21985/exchange-rate-service/internal/repositories"

	"github.com/abduls21985/exchange-rate-service/internal/models"
)

type ExchangeRateService interface {
	AddExchangeRates(data models.ExchangeRateData) error
	FetchExchangeRates(currencyCode string, timestamp int64) ([]models.ExchangeRate, error)
	GetAllCurrencies() ([]models.Currency, error)
}

type exchangeRateService struct {
	repo repositories.ExchangeRateRepository
}

func NewExchangeRateService(repo repositories.ExchangeRateRepository) ExchangeRateService {
	return &exchangeRateService{repo: repo}
}

func (s *exchangeRateService) AddExchangeRates(data models.ExchangeRateData) error {
	// Get or create base currency
	baseCurrency, err := s.repo.GetCurrencyByCode(data.Base)
	if err != nil {
		// If not found, create it
		baseCurrency, err = s.repo.CreateCurrency(data.Base, "")
		if err != nil {
			return fmt.Errorf("failed to create base currency: %v", err)
		}
	}

	timestamp := time.Unix(data.Timestamp, 0).UTC()

	for code, rate := range data.Rates {
		// Get or create currency
		currency, err := s.repo.GetCurrencyByCode(code)
		if err != nil {
			// If not found, create it
			currency, err = s.repo.CreateCurrency(code, "")
			if err != nil {
				return fmt.Errorf("failed to create currency %s: %v", code, err)
			}
		}

		exchangeRate := models.ExchangeRate{
			CurrencyID:     currency.ID,
			Rate:           rate,
			Timestamp:      timestamp,
			BaseCurrencyID: baseCurrency.ID,
		}

		if err := s.repo.InsertOrUpdateExchangeRate(&exchangeRate); err != nil {
			return fmt.Errorf("failed to insert/update exchange rate for %s: %v", code, err)
		}
	}

	return nil
}

func (s *exchangeRateService) FetchExchangeRates(currencyCode string, timestamp int64) ([]models.ExchangeRate, error) {
	return s.repo.GetExchangeRates(currencyCode, timestamp)
}

func (s *exchangeRateService) GetAllCurrencies() ([]models.Currency, error) {
	return s.repo.GetAllCurrencies()
}

// CountExchangeRates returns the count of exchange rates
func (s *exchangeRateService) CountExchangeRates() (int, error) {
	// Implement the logic to count exchange rates
	// This is a placeholder implementation
	count := 0 // Replace with actual logic to count exchange rates
	if count < 0 {
		return 0, errors.New("failed to count exchange rates")
	}
	return count, nil
}
