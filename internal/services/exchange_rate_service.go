// internal/services/exchange_rate_service.go

package services

import (
	"fmt"
	"time"

	"github.com/abduls21985/exchange-rate-service/internal/repositories"

	"github.com/abduls21985/exchange-rate-service/internal/models"
)

type ExchangeRateService interface {
	AddExchangeRates(data models.ExchangeRateData) error
	FetchExchangeRates(currencyCode string, timestamp int64) ([]models.ExchangeRate, error)
	GetAllCurrencies() ([]models.Currency, error)
	CountExchangeRates() (int, error)
	GetHistoricalExchangeRates(currencyCode string, startDate, endDate int64) ([]models.ExchangeRate, error)
	ConvertCurrency(fromCurrency, toCurrency string, amount float64) (float64, error)
	ConvertToBaseCurrency(rates []models.ExchangeRate, baseCurrency string) ([]models.ExchangeRate, error)
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

// internal/services/exchange_rate_service.go

func (s *exchangeRateService) GetHistoricalExchangeRates(currencyCode string, startDate, endDate int64) ([]models.ExchangeRate, error) {
	return s.repo.GetHistoricalExchangeRates(currencyCode, startDate, endDate)
}

// internal/services/exchange_rate_service.go

func (s *exchangeRateService) ConvertCurrency(fromCurrency, toCurrency string, amount float64) (float64, error) {
	// Get exchange rate for the source currency
	fromRate, err := s.repo.GetExchangeRateByCurrency(fromCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate for %s: %v", fromCurrency, err)
	}

	// Get exchange rate for the target currency
	toRate, err := s.repo.GetExchangeRateByCurrency(toCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate for %s: %v", toCurrency, err)
	}

	// Perform the conversion
	convertedAmount := (amount / fromRate.Rate) * toRate.Rate
	return convertedAmount, nil
}

// internal/services/exchange_rate_service.go

func (s *exchangeRateService) ConvertToBaseCurrency(rates []models.ExchangeRate, baseCurrency string) ([]models.ExchangeRate, error) {
	// Get the exchange rate for the specified base currency
	baseRate, err := s.repo.GetExchangeRateByCurrency(baseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base currency rate for %s: %v", baseCurrency, err)
	}

	// Convert each rate in the list to the new base currency
	for i := range rates {
		if rates[i].CurrencyID != baseRate.CurrencyID {
			rates[i].Rate /= baseRate.Rate
		}
	}

	return rates, nil
}

// internal/services/exchange_rate_service.go

func (s *exchangeRateService) CountExchangeRates() (int, error) {
	// Use the repository to get the count of exchange rates
	count, err := s.repo.CountExchangeRates()
	if err != nil {
		return 0, fmt.Errorf("failed to count exchange rates: %v", err)
	}
	return count, nil
}
