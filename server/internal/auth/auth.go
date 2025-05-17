package auth

import (
	"encoding/json"
	"net/http"
)

// HandleLogin authenticates a user and returns a token
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mock authentication - replace with real auth in production
	if creds.Username == "admin" && creds.Password == "password" {
		token := "mock-jwt-token" // In production, generate a real JWT
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

// VerifyToken checks if a token is valid
func VerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// In a real app, validate JWT signature, check expiration, etc.
	w.WriteHeader(http.StatusOK)
}
