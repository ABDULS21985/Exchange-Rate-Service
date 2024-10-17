package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abduls21985/exchange-rate-service/internal/services"
)

type AuthController struct {
	AuthService services.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// AuthenticateUser handles POST /api/login
func (c *AuthController) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.AuthService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := c.AuthService.GenerateJWT(user.Username)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "Login successful", "token": token}, http.StatusOK)
}
