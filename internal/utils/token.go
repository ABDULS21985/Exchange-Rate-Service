package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken generates a random token string
func GenerateRandomToken() string {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "defaultRandomToken" // Fallback in case of an error
	}
	return hex.EncodeToString(bytes)
}
