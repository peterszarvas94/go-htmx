package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken() JWT {
	currentTime := time.Now().Unix()
	expirationTime := currentTime + 3600
	secret, found := os.LookupEnv("JWT_SECRET")

	if !found {
		fmt.Println("Error: JWT_SECRET environment variable not found")
		os.Exit(1)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": currentTime,
		"exp": expirationTime,
		"sub": "user",
	})

	signed_token, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return JWT{
		Token:   signed_token,
		Expires: expirationTime,
	}
}

func ValidateToken(token string) (jwt.MapClaims, error) {
	secret, found := os.LookupEnv("JWT_SECRET")

	if !found {
		return nil, fmt.Errorf("Error: JWT_SECRET environment variable not found")
	}

	parsed_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Error: Unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsed_token.Claims.(jwt.MapClaims)
	if !ok || !parsed_token.Valid {
		return nil, fmt.Errorf("Error: Invalid JWT token")
	}

	return claims, nil
}
