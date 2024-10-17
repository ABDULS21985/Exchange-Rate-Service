// package routes

package routes

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/abduls21985/exchange-rate-service/internal/controllers"
	"github.com/abduls21985/exchange-rate-service/internal/repositories"
	"github.com/abduls21985/exchange-rate-service/internal/services"
	"github.com/abduls21985/exchange-rate-service/pkg/middleware"
)

// InitializeRoutes sets up all the routes for the application
func InitializeRoutes(router *mux.Router, db *gorm.DB) {
	// Initialize repositories
	exchangeRateRepo := repositories.NewExchangeRateRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	exchangeRateService := services.NewExchangeRateService(exchangeRateRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userService)

	// Initialize controllers
	exchangeRateController := controllers.NewExchangeRateController(exchangeRateService)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	// User Management Routes
	router.HandleFunc("/api/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/api/login", authController.AuthenticateUser).Methods("POST")

	// API subrouter for protected routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)

	// Exchange Rate Routes
	apiRouter.HandleFunc("/fetch-cbn-exchange-rates", exchangeRateController.GetExchangeRates).Methods("GET")
	apiRouter.HandleFunc("/exchange-rates", exchangeRateController.PostExchangeRates).Methods("POST")
	apiRouter.HandleFunc("/currencies", exchangeRateController.GetCurrencies).Methods("GET")
	apiRouter.HandleFunc("/exchange-rates/historical", exchangeRateController.GetHistoricalExchangeRates)
	apiRouter.HandleFunc("/exchange-rates/convert", exchangeRateController.ConvertCurrency).Methods("POST")
	apiRouter.HandleFunc("/exchange-rates/base-convert", exchangeRateController.ConvertRatesToBaseCurrency).Methods("GET")
	apiRouter.HandleFunc("/exchange-rates/count", exchangeRateController.GetExchangeRateCount).Methods("GET")
	apiRouter.HandleFunc("/convert-rates", exchangeRateController.ConvertMultipleRatesToBaseCurrency).Methods("POST")

	// Health Check Route (public)
	router.HandleFunc("/api/health", exchangeRateController.HealthCheck).Methods("GET")
}
