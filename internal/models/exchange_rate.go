// internal/models/exchange_rate.go

package models

import "time"

// ExchangeRateData represents the incoming JSON structure
type ExchangeRateData struct {
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

// Currency represents the currencies table
type Currency struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Code string `gorm:"size:3;uniqueIndex;not null" json:"code"` // ISO 4217 standard length
	Name string `gorm:"size:100" json:"name,omitempty"`
}

// ExchangeRate represents the exchange_rates table
type ExchangeRate struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CurrencyID     uint      `gorm:"not null" json:"currency_id"`
	Currency       Currency  `gorm:"foreignKey:CurrencyID"` // Foreign key relationship
	Rate           float64   `gorm:"not null" json:"rate"`
	Timestamp      time.Time `gorm:"not null" json:"timestamp"`
	BaseCurrencyID uint      `gorm:"not null" json:"base_currency_id"`
	BaseCurrency   Currency  `gorm:"foreignKey:BaseCurrencyID"` // Foreign key relationship
}
