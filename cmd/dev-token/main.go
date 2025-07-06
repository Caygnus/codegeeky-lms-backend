package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/dev-token/main.go <jwt_secret>")
		fmt.Println("This script generates a long-lived JWT token for development purposes.")
		fmt.Println("The token will be valid for 1 year from now.")
		os.Exit(1)
	}

	jwtSecret := os.Args[1]

	// Generate a random user ID for development
	userID := generateRandomID()

	// Create claims for the token
	claims := jwt.MapClaims{
		"sub":   userID,
		"aud":   "authenticated",
		"role":  "authenticated",
		"email": "dev@example.com",
		"phone": "",
		"app_metadata": map[string]interface{}{
			"role": "ADMIN",
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().AddDate(1, 0, 0).Unix(), // 1 year from now
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	fmt.Println("=== Development JWT Token ===")
	fmt.Println("This token is valid for 1 year from now.")
	fmt.Println("Use it in your Authorization header as: Bearer <token>")
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(tokenString)
	fmt.Println()
	fmt.Println("User ID:", userID)
	fmt.Println("Email: dev@example.com")
	fmt.Println("Role: ADMIN")
	fmt.Println("Expires:", time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("⚠️  WARNING: This token is for development only!")
	fmt.Println("⚠️  Never use this in production!")
	fmt.Println("⚠️  Keep your JWT secret secure!")
}

func generateRandomID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate random ID: %v", err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}
