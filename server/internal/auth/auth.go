package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"k8_gui/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret returns the JWT secret from environment variable or a default value
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key_change_in_production"
	}
	return []byte(secret)
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, utils.MsgInvalidRequestBody, http.StatusBadRequest)
		return
	}

	if creds.Username != "admin" || creds.Password != "password" {
		http.Error(w, utils.MsgInvalidCredentials, http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		http.Error(w, utils.MsgFailedGenerateToken, http.StatusInternalServerError)
		return
	}

	fmt.Println("Generated Token:", tokenString)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, utils.MsgAuthorizationHeaderRequired, http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		// No Bearer prefix found, use the header as is
		tokenString = authHeader
	}

	if tokenString == "" {
		http.Error(w, utils.MsgTokenRequired, http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		fmt.Printf("Token parsing error: %v\n", err)
		http.Error(w, utils.MsgInvalidToken, http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		fmt.Printf("Token is invalid\n")
		http.Error(w, utils.MsgInvalidToken, http.StatusUnauthorized)
		return
	}

	fmt.Printf("Token validated successfully for user: %s\n", claims.Username)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "valid", "username": claims.Username})
}

func ValidateJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from "Authorization: Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, utils.MsgMissingAuthorizationHeader, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, utils.MsgInvalidOrExpiredToken, http.StatusUnauthorized)
			return
		}

		// Token is valid, continue to the handler
		next.ServeHTTP(w, r)
	})
}
