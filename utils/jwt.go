package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(id int) (JWT, error) {
	currentTime := time.Now().Unix()
	expirationTime := currentTime + 3600

	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		return JWT{}, errors.New("Error: JWT_SECRET environment variable not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": currentTime,
		"exp": expirationTime,
		"sub": id,
	})

	signedToken, signErr := token.SignedString([]byte(secret))
	if signErr != nil {
		return JWT{}, signErr
	}

	return JWT{
		Token:   signedToken,
		Expires: expirationTime,
	}, nil
}

func ValidateToken(token string) (jwt.MapClaims, error) {
	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		return nil, errors.New("Error: JWT_SECRET environment variable not found")
	}

	parsedToken, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Error: Unexpected signing method")
		}

		return []byte(secret), nil
	})

	if parseErr != nil {
		return nil, parseErr
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("Error: Invalid signature")
	}

	return claims, nil
}
