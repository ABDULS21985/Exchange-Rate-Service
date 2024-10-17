# Exchange Rate Service

The **Exchange Rate Service** is a RESTful API built with Go, designed to fetch and manage historical exchange rate data from an external API. This service supports functionalities such as fetching exchange rate data, storing the data in a database, manual and automated data synchronization, and user management for authentication.

## Features

- Fetch historical exchange rates
- Automatic daily synchronization to update the database with the latest exchange rate data.
- Manual trigger for data synchronization through an API endpoint.
- User registration, authentication, and password reset functionality.
- Secure API endpoints protected with JWT-based authentication.

## Technologies Used

- **Golang**: Main programming language.
- **GORM**: ORM library for managing database interactions.
- **Viper**: Configuration management.
- **Gorilla Mux**: HTTP router and URL matcher for building the API.
- **Cron**: For scheduling tasks.
- **PostgreSQL**: Database for persisting exchange rates and user information.

## Project Structure


## Getting Started

### Prerequisites

Make sure you have the following installed on your machine:

- [Go](https://golang.org/dl/) 1.16 or higher
- [PostgreSQL](https://www.postgresql.org/download/)
- [Git](https://git-scm.com/downloads)

### Configuration

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Application-Management-Division/exchange-rate-service.git
   cd exchange-rate-service

# .env

# Database configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password
DATABASE_DBNAME=exchange_rate_db
DATABASE_SSLMODE=disable

# JWT Secret
JWT_SECRET=drkatanga2020

# Exchange Rate API
EXCHANGE_RATE_API_URL=https://openexchangerates.org/api/historical/%s.json?app_id=your_app_id

# Server configuration
SERVER_PORT=8080


database:
  host: ${DATABASE_HOST}
  port: ${DATABASE_PORT}
  user: ${DATABASE_USER}
  password: ${DATABASE_PASSWORD}
  dbname: ${DATABASE_DBNAME}
  sslmode: ${DATABASE_SSLMODE}
jwt_secret: ${JWT_SECRET}
server:
  port: ${SERVER_PORT}
exchange_rate_api_url: ${EXCHANGE_RATE_API_URL}
