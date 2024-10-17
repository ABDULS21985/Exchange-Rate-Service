// package repositories

package repositories

import (
	"github.com/abduls21985/exchange-rate-service/internal/models"
	"gorm.io/gorm"
)

// ExchangeRateRepository interface defines the methods for exchange rate operations
type ExchangeRateRepository interface {
	GetCurrencyByCode(code string) (*models.Currency, error)
	CreateCurrency(code string, name string) (*models.Currency, error)
	InsertOrUpdateExchangeRate(rate *models.ExchangeRate) error
	GetExchangeRates(currencyCode string, timestamp int64) ([]models.ExchangeRate, error)
	GetAllCurrencies() ([]models.Currency, error)
}

type exchangeRateRepository struct {
	db *gorm.DB
}

// NewExchangeRateRepository creates a new instance of ExchangeRateRepository
func NewExchangeRateRepository(db *gorm.DB) ExchangeRateRepository {
	return &exchangeRateRepository{db}
}

// GetCurrencyByCode retrieves a currency by its code
func (r *exchangeRateRepository) GetCurrencyByCode(code string) (*models.Currency, error) {
	var currency models.Currency
	err := r.db.Where("code = ?", code).First(&currency).Error
	return &currency, err
}

// CreateCurrency adds a new currency to the database
func (r *exchangeRateRepository) CreateCurrency(code string, name string) (*models.Currency, error) {
	currency := &models.Currency{Code: code, Name: name}
	err := r.db.Create(currency).Error
	return currency, err
}

// InsertOrUpdateExchangeRate inserts a new exchange rate or updates the existing one
func (r *exchangeRateRepository) InsertOrUpdateExchangeRate(rate *models.ExchangeRate) error {
	// GORM's `Save` method will update if the record already exists
	return r.db.Save(rate).Error
}

// GetExchangeRates retrieves exchange rates based on currency code and timestamp
func (r *exchangeRateRepository) GetExchangeRates(currencyCode string, timestamp int64) ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate

	// Start the query
	query := r.db.Joins("JOIN currencies AS c1 ON exchange_rates.currency_id = c1.id").
		Joins("JOIN currencies AS c2 ON exchange_rates.base_currency_id = c2.id").
		Preload("Currency").Preload("BaseCurrency")

	// Apply filters if provided
	if currencyCode != "" {
		query = query.Where("c1.code = ?", currencyCode)
	}
	if timestamp != 0 {
		query = query.Where("exchange_rates.timestamp = ?", timestamp)
	}

	// Execute the query
	err := query.Find(&rates).Error
	return rates, err
}

// GetAllCurrencies retrieves all currencies
func (r *exchangeRateRepository) GetAllCurrencies() ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.Order("code ASC").Find(&currencies).Error
	return currencies, err
}
