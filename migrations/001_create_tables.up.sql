-- migrations/001_create_tables.up.sql

-- Table to store currency information
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(50)
);

-- Table to store exchange rates
CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER REFERENCES currencies(id) ON DELETE CASCADE,
    rate DECIMAL(18,6) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    base_currency_id INTEGER REFERENCES currencies(id) ON DELETE CASCADE,
    UNIQUE (currency_id, timestamp)
);

-- Index for faster queries on exchange_rates
CREATE INDEX IF NOT EXISTS idx_exchange_rates_currency_timestamp ON exchange_rates (currency_id, timestamp);
