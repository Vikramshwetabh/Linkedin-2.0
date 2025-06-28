package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a simple token (in production, use JWT)
func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// WriteJSONResponse writes a JSON response to the HTTP writer
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse writes an error response to the HTTP writer
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"success": false,
		"error":   message,
	}
	WriteJSONResponse(w, statusCode, response)
}

// WriteSuccessResponse writes a success response to the HTTP writer
func WriteSuccessResponse(w http.ResponseWriter, data interface{}) {
	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	WriteJSONResponse(w, http.StatusOK, response)
}

// GetUserIDFromToken extracts user ID from token (simplified - in production use JWT)
func GetUserIDFromToken(token string) (int, error) {
	// This is a simplified implementation
	// In production, you would decode a JWT token
	// For now, we'll use a simple token format: "user_id:timestamp"
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}

	var userID int
	_, err := fmt.Sscanf(parts[0], "%d", &userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// CreateToken creates a token for a user
func CreateToken(userID int) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d:%d", userID, timestamp)
}
