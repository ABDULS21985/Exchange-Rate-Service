package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse writes a JSON response to the HTTP response writer
func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
