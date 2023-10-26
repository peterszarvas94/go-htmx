package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CSRF struct {
	Token string
}

// Generate a new CSRF token
func NewCSRFToken() (*CSRF, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}

	// Encode the random bytes as a base64 string
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	return &CSRF{
		Token: token,
	}, nil
}

// Validate a given token against a stored token
func (c *CSRF) Validate(requestToken string) bool {
	return c.Token == requestToken
}

type Token string

const (
	ACCESS  Token = "access"
	REFRESH Token = "refresh"
)

func NewToken(id int, tokentype Token) (JWT, error) {
	currentTime := time.Now().Unix()
	expirationTime := currentTime + 3600

	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		return JWT{}, errors.New("Error: JWT_SECRET environment variable not found")
	}

	idStr := strconv.Itoa(id)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": currentTime,
		"exp": expirationTime,
		"sub": idStr,
		"typ": tokentype,
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
