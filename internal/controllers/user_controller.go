package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abduls21985/exchange-rate-service/internal/models"
	"github.com/abduls21985/exchange-rate-service/internal/services"
)

type UserController struct {
	Service services.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(service services.UserService) *UserController {
	return &UserController{Service: service}
}

// RegisterUser handles POST /api/register
func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Define a struct to receive the registration payload
	var req struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
		MdaID       string `json:"mda_id"`
		MDA         string `json:"mda"`
		Role        string `json:"role"`
	}

	// Decode the incoming JSON request into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Map the request data to the User model
	user := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Username:    req.Username,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password, // Note: Password will be hashed in the service layer
		MdaID:       req.MdaID,
		MDA:         req.MDA,
		Role:        req.Role,
	}

	// Call the service to register the user
	newUser, err := c.Service.RegisterUser(&user)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Respond with the newly created user
	jsonResponse(w, newUser, http.StatusCreated)
}

// UpdateUser handles PUT /api/users/{id}
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.Service.UpdateUser(&user); err != nil {
		log.Printf("Error updating user: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "User updated successfully"}, http.StatusOK)
}

// InitiatePasswordReset handles POST /api/password-reset/initiate
func (c *UserController) InitiatePasswordReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.Service.InitiatePasswordReset(req.Email)
	if err != nil {
		log.Printf("Error initiating password reset: %v", err)
		http.Error(w, "Failed to initiate password reset", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "Password reset link sent", "email": user.Email}, http.StatusOK)
}

// ResetPassword handles POST /api/password-reset/complete
func (c *UserController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.Service.ResetPassword(req.Token, req.NewPassword); err != nil {
		log.Printf("Error resetting password: %v", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "Password reset successful"}, http.StatusOK)
}

// Helper function to write JSON responses
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
